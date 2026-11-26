package broker

import (
	"context"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
)

var (
	accKeywords  = []string{"akumulasi", "accumulation", "bullish"}
	distKeywords = []string{"distribusi", "distribution", "bearish"}
)

func containsAny(s string, kws []string) bool {
	low := strings.ToLower(s)
	for _, kw := range kws {
		if strings.Contains(low, kw) {
			return true
		}
	}
	return false
}

// Service implements aggregation and business logic for /api/broker-flow/*.
type Service struct {
	Repo *Repository
}

// BuildSummary computes BrokerFlowSummary matching routers/broker.py:get_broker_flow_summary.
func (s *Service) BuildSummary(ctx context.Context, ticker string, days int) (*BrokerFlowSummary, error) {
	now := time.Now()
	cutoff := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -days)
	rows, err := s.Repo.GetSummaryRows(ctx, ticker, cutoff)
	if err != nil || len(rows) == 0 {
		return nil, err
	}

	var foreignSum, localSum, govSum, totalSum float64
	accDays, distDays := 0, 0
	for _, r := range rows {
		foreignSum += shared.Finite(r.ForeignNetBroker, 0)
		localSum += shared.Finite(r.LocalNetBroker, 0)
		govSum += shared.Finite(r.GovNetBroker, 0)
		totalSum += shared.Finite(r.TotalValue, 0)
		if r.BandarSignal != nil {
			if containsAny(*r.BandarSignal, accKeywords) {
				accDays++
			}
			if containsAny(*r.BandarSignal, distKeywords) {
				distDays++
			}
		}
	}

	totalParties := math.Abs(foreignSum) + math.Abs(localSum) + math.Abs(govSum)
	foreignDom := 0.0
	if totalParties > 0 {
		foreignDom = math.Abs(foreignSum) / totalParties * 100.0
	}

	latest := rows[len(rows)-1]
	// Preserve Python truthiness: float(x) if latest.bandar_signal_score else None
	var latestScore *float64
	if latest.BandarSignalScore != nil && *latest.BandarSignalScore != 0 {
		v := *latest.BandarSignalScore
		latestScore = &v
	}
	latestDate := latest.Date

	return &BrokerFlowSummary{
		Ticker:                  ticker,
		PeriodDays:              days,
		TradingDays:             len(rows),
		ForeignNetBrokerSum:     shared.Round(foreignSum, 2),
		LocalNetBrokerSum:       shared.Round(localSum, 2),
		GovNetBrokerSum:         shared.Round(govSum, 2),
		TotalValueSum:           shared.Round(totalSum, 2),
		ForeignDominancePct:     shared.Round(foreignDom, 2),
		LatestBandarSignal:      latest.BandarSignal,
		LatestBandarSignalScore: latestScore,
		LatestForeignSignal:     latest.ForeignSignal,
		AccumulationDays:        accDays,
		DistributionDays:        distDays,
		LatestDate:              &latestDate,
	}, nil
}

// GetWhalesList returns classified brokers sorted by (is_smart DESC, broker_code ASC).
func (s *Service) GetWhalesList(ctx context.Context) ([]BrokerWhaleListItem, error) {
	codes, err := s.Repo.GetDistinctBrokers(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]BrokerWhaleListItem, 0, len(codes))
	for _, c := range codes {
		prof := shared.BrokerProfileOf(c)
		meta, ok := shared.ProfileMeta[prof]
		if !ok {
			meta = shared.ProfileMetadata{Label: "Lainnya", Description: "Profil broker lainnya"}
		}
		out = append(out, BrokerWhaleListItem{
			BrokerCode:   c,
			Profile:      prof,
			ProfileLabel: meta.Label,
			ProfileDesc:  meta.Description,
			IsSmart:      shared.IsSmartProfile(prof),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].IsSmart != out[j].IsSmart {
			return out[i].IsSmart
		}
		return out[i].BrokerCode < out[j].BrokerCode
	})
	return out, nil
}

// GetWhaleDetail returns accumulation/distribution breakdown for a single broker.
func (s *Service) GetWhaleDetail(ctx context.Context, code string, days, limit int) (*BrokerWhalesDetailResponse, error) {
	dates, rows, err := s.Repo.GetWhaleActivity(ctx, code, days)
	if err != nil || len(dates) == 0 {
		return nil, err
	}

	minDate, maxDate := dates[0], dates[0]
	for _, d := range dates[1:] {
		if d.Before(minDate) {
			minDate = d
		}
		if d.After(maxDate) {
			maxDate = d
		}
	}

	accs := make([]WhaleStockItem, 0)
	dists := make([]WhaleStockItem, 0)
	var totalBuy, totalSell, totalNet float64

	for _, r := range rows {
		totalBuy += r.TotalBuy
		totalSell += r.TotalSell
		totalNet += r.TotalNet

		var bAvg, sAvg *float64
		if r.BuyAvgPrice != nil {
			v := shared.Round(*r.BuyAvgPrice, 2)
			bAvg = &v
		}
		if r.SellAvgPrice != nil {
			v := shared.Round(*r.SellAvgPrice, 2)
			sAvg = &v
		}

		item := WhaleStockItem{
			Ticker:       r.Ticker,
			NetValue:     shared.Round(r.TotalNet, 2),
			BuyValue:     shared.Round(r.TotalBuy, 2),
			SellValue:    shared.Round(r.TotalSell, 2),
			NetLot:       shared.Round(r.NetLot, 2),
			BuyAvgPrice:  bAvg,
			SellAvgPrice: sAvg,
		}
		if r.TotalNet > 0 {
			accs = append(accs, item)
		} else if r.TotalNet < 0 {
			dists = append(dists, item)
		}
	}

	sort.Slice(accs, func(i, j int) bool { return accs[i].NetValue > accs[j].NetValue })
	sort.Slice(dists, func(i, j int) bool { return dists[i].NetValue < dists[j].NetValue })

	if len(accs) > limit {
		accs = accs[:limit]
	}
	if len(dists) > limit {
		dists = dists[:limit]
	}

	prof := shared.BrokerProfileOf(code)
	meta, ok := shared.ProfileMeta[prof]
	if !ok {
		meta = shared.ProfileMetadata{Label: "Lainnya", Description: "Profil broker lainnya"}
	}

	return &BrokerWhalesDetailResponse{
		BrokerCode:       code,
		Profile:          prof,
		ProfileLabel:     meta.Label,
		ProfileDesc:      meta.Description,
		IsSmart:          shared.IsSmartProfile(prof),
		PeriodDays:       len(dates),
		StartDate:        minDate.Format("2006-01-02"),
		EndDate:          maxDate.Format("2006-01-02"),
		TotalBuyValue:    shared.Round(totalBuy, 2),
		TotalSellValue:   shared.Round(totalSell, 2),
		TotalNetValue:    shared.Round(totalNet, 2),
		TopAccumulations: accs,
		TopDistributions: dists,
	}, nil
}
