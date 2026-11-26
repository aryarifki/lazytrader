package nyangkut

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var monthNamesID = [12]string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

var monthNamesShort = [12]string{
	"Jan", "Feb", "Mar", "Apr", "Mei", "Jun",
	"Jul", "Agu", "Sep", "Okt", "Nov", "Des",
}

// OwnershipCategory mirrors OwnershipCategory schema in InvestOwl schemas.py.
type OwnershipCategory struct {
	Category string  `json:"category"`
	Pct      float64 `json:"pct"`
	Color    string  `json:"color"`
	Desc     string  `json:"desc"`
}

// VolumeNode mirrors VolumeNode schema in InvestOwl schemas.py.
type VolumeNode struct {
	Price          float64 `json:"price"`
	VolumePct      float64 `json:"volumePct"`
	IsNyangkutZone bool    `json:"isNyangkutZone"`
}

// OverviewStats mirrors NyangkutOverviewStats schema.
type OverviewStats struct {
	MarketCap           string `json:"marketCap"`
	TotalSahamScripless string `json:"totalSahamScripless"`
	ScriplessPctText    string `json:"scriplessPctText"`
	FreeFloat           string `json:"freeFloat"`
	JumlahPemegangSaham string `json:"jumlahPemegangSaham"`
	AvgSahamRetail      string `json:"avgSahamRetail"`
	Konglo              string `json:"konglo"`
}

// TrendPoint mirrors NyangkutTrendPoint schema.
type TrendPoint struct {
	Month string  `json:"month"`
	Val   float64 `json:"val"`
}

// TrendCard mirrors NyangkutTrendCard schema.
type TrendCard struct {
	Label  string       `json:"label"`
	Value  string       `json:"value"`
	Change string       `json:"change"`
	Color  string       `json:"color"`
	Points []TrendPoint `json:"points"`
}

// StockDataResponse mirrors NyangkutStockDataResponse schema.
type StockDataResponse struct {
	Ticker               string              `json:"ticker"`
	Name                 string              `json:"name"`
	CurrentPrice         float64             `json:"currentPrice"`
	AvgAccumulationPrice float64             `json:"avgAccumulationPrice"`
	UnderwaterPct        float64             `json:"underwaterPct"`
	EstimatedNyangkutPct float64             `json:"estimatedNyangkutPct"`
	IsEstimated          bool                `json:"isEstimated"`
	TotalSid             int64               `json:"totalSid"`
	SidChangeMoM         int64               `json:"sidChangeMoM"`
	SidChangePct         float64             `json:"sidChangePct"`
	RarScore             float64             `json:"rarScore"`
	PbiScore             float64             `json:"pbiScore"`
	IsHeroBadge          bool                `json:"isHeroBadge"`
	Ownership            []OwnershipCategory `json:"ownership"`
	VolumeNodes          []VolumeNode        `json:"volumeNodes"`
	PeriodDate           string              `json:"periodDate"`
	UpdateLagLabel       string              `json:"updateLagLabel"`
	SparklineSid         []int64             `json:"sparklineSid"`
	OverviewStats        OverviewStats       `json:"overviewStats"`
	TrendCards           []TrendCard         `json:"trendCards"`
}

// LeaderboardHeroItem mirrors LeaderboardHeroItem schema.
type LeaderboardHeroItem struct {
	Rank                 int     `json:"rank"`
	Ticker               string  `json:"ticker"`
	Name                 string  `json:"name"`
	PbiScore             float64 `json:"pbiScore"`
	RarScore             float64 `json:"rarScore"`
	PriceDropPct         float64 `json:"priceDropPct"`
	SidGrowthMoM         int64   `json:"sidGrowthMoM"`
	CurrentPrice         float64 `json:"currentPrice"`
	AvgAccumulationPrice float64 `json:"avgAccumulationPrice"`
}

// SummaryStatsResponse mirrors get_nyangkut_summary_stats payload.
type SummaryStatsResponse struct {
	PeriodDate           string  `json:"periodDate"`
	TotalStocksAnalyzed  int64   `json:"totalStocksAnalyzed"`
	AvgRetailNyangkutPct float64 `json:"avgRetailNyangkutPct"`
	TotalTrackedSid      int64   `json:"totalTrackedSid"`
	HeroStocksCount      int64   `json:"heroStocksCount"`
}

// TickerSearchResult mirrors TickerSearchResult schema.
type TickerSearchResult struct {
	Ticker        string   `json:"ticker"`
	Name          string   `json:"name"`
	Sector        string   `json:"sector"`
	IsHeroBadge   bool     `json:"isHeroBadge"`
	UnderwaterPct *float64 `json:"underwaterPct"`
}

// roundEven matches Python round() exact-IEEE754 decimal banker's rounding and normalizes -0.0 to 0.0.
func roundEven(v float64, digits int) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0.0
	}
	if digits < 0 {
		pow := math.Pow10(-digits)
		r := math.RoundToEven(v/pow) * pow
		if r == 0 {
			return 0.0
		}
		return r
	}
	r, err := strconv.ParseFloat(strconv.FormatFloat(v, 'f', digits, 64), 64)
	if err != nil || r == 0 {
		return 0.0
	}
	return r
}

// formatIntDot formats an integer with '.' thousands separators matching Python f"{n:,}".replace(",", ".").
func formatIntDot(n int64) string {
	neg := n < 0
	if neg {
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		if neg {
			return "-" + s
		}
		return s
	}
	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	rem := len(s) % 3
	if rem > 0 {
		b.WriteString(s[:rem])
		if len(s) > rem {
			b.WriteByte('.')
		}
	}
	for i := rem; i < len(s); i += 3 {
		if i > rem {
			b.WriteByte('.')
		}
		b.WriteString(s[i : i+3])
	}
	return b.String()
}
