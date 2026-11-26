package forensic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/flags"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
)

// Service computes Foreign Exodus Radar with Redis+L1 caching.
type Service struct {
	Repo  *Repository
	Cache *cache.Store
	Flags *flags.Store
}

// GetForeignExodusRadar computes or fetches cached Foreign Exodus Radar for `days`.
func (s *Service) GetForeignExodusRadar(ctx context.Context, days int) (*ExodusRadarResponse, error) {
	cacheKey := fmt.Sprintf("forensic:foreign_exodus:%d", days)
	if s.Cache != nil {
		if raw, ok := s.Cache.Get(ctx, cacheKey); ok {
			var cached ExodusRadarResponse
			if json.Unmarshal(raw, &cached) == nil {
				return &cached, nil
			}
		}
	}

	thresholdMode := "absolute"
	if s.Flags != nil {
		thresholdMode = s.Flags.GetString("EXODUS_THRESHOLD_MODE", "absolute")
	}
	lookbackLimit := ExodusLookbackLimit(days, thresholdMode)

	timeline, err := s.Repo.FetchMacroTimeline(ctx, lookbackLimit)
	if err != nil {
		return nil, err
	}

	// Forward-fill then backward-fill 0/null ihsg_close
	lastValid := 0.0
	for i := range timeline {
		if timeline[i].IHSGClose > 0 {
			lastValid = timeline[i].IHSGClose
		} else if lastValid > 0 {
			timeline[i].IHSGClose = lastValid
		}
	}
	lastValid = 0.0
	for i := len(timeline) - 1; i >= 0; i-- {
		if timeline[i].IHSGClose > 0 {
			lastValid = timeline[i].IHSGClose
		} else if lastValid > 0 {
			timeline[i].IHSGClose = lastValid
		}
	}

	ihsgCurrent := 0.0
	ihsgPrev := 0.0
	if n := len(timeline); n > 0 {
		ihsgCurrent = timeline[n-1].IHSGClose
		if n > 1 {
			ihsgPrev = timeline[n-2].IHSGClose
		} else {
			ihsgPrev = ihsgCurrent
		}
	}
	ihsgChangePct := 0.0
	if ihsgPrev > 0 {
		ihsgChangePct = shared.Round(((ihsgCurrent-ihsgPrev)/ihsgPrev)*100.0, 2)
	}

	marketState := EvaluateMarketState(timeline, s.Flags)

	if len(timeline) > days {
		start := len(timeline) - days
		baseCum := 0.0
		if start > 0 {
			baseCum = timeline[start-1].CumulativeFlow
		}
		timeline = append([]TimelinePoint(nil), timeline[start:]...)
		for i := range timeline {
			timeline[i].CumulativeFlow -= baseCum
		}
	}

	culprits, err := s.Repo.FetchCulpritBrokers(ctx)
	if err != nil {
		return nil, err
	}

	res := &ExodusRadarResponse{
		MarketState:    marketState,
		IHSGCurrent:    ihsgCurrent,
		IHSGChangePct:  ihsgChangePct,
		Timeline:       timeline,
		CulpritBrokers: culprits,
	}

	if s.Cache != nil {
		if payload, err := json.Marshal(res); err == nil {
			s.Cache.Set(ctx, cacheKey, payload, 300*time.Second)
		}
	}
	return res, nil
}

// Refresh purges forensic cache keys matching routers/forensic.py.
func (s *Service) Refresh(ctx context.Context) {
	if s.Cache != nil {
		s.Cache.SafePurge(ctx, "forensic:*")
	}
}
