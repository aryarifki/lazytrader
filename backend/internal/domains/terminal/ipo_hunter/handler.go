package ipohunter

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler exposes the 4 IPO Hunter Pro HTTP endpoints.
type Handler struct {
	Repo  *Repository
	Cache *cache.Store
}

// Routes mounts /api/v1/terminal/ipo-hunter/*.
func Routes(db *pgxpool.Pool, store *cache.Store) chi.Router {
	h := &Handler{Repo: &Repository{DB: db}, Cache: store}
	r := chi.NewRouter()
	r.Get("/meta", h.GetMeta)
	r.Get("/list", h.GetList)
	r.Get("/analysis/{ticker}", h.GetAnalysis)
	r.Post("/screener", h.PostScreener)
	return r
}

func (h *Handler) GetMeta(w http.ResponseWriter, r *http.Request) {
	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, "ipo_hunter_meta", 120*time.Second, func(ctx context.Context) (any, int, string, error) {
		masters, err := h.Repo.ListMasters(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		dates, err := h.Repo.ListDistinctDailyDates(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		withData, err := h.Repo.ListDistinctDailyTickers(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		d0Date, err := h.Repo.GetAnyD0Date(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		return map[string]any{
			"success": true,
			"data":    BuildMetaData(masters, dates, withData, d0Date),
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func parseBoolQuery(raw string, defVal bool, name string) (bool, error) {
	s := strings.TrimSpace(strings.ToLower(raw))
	if s == "" {
		return defVal, nil
	}
	switch s {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("Input should be a valid boolean for '%s'", name)
	}
}

func parseOptionalIntQuery(raw, name string) (*int, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("Input should be a valid integer for '%s'", name)
	}
	return &v, nil
}

func parseIntQuery(raw string, defVal int, name string) (int, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return defVal, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("Input should be a valid integer for '%s'", name)
	}
	return v, nil
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	leadOnly, err := parseBoolQuery(q.Get("lead_only"), false, "lead_only")
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	yearPtr, err := parseOptionalIntQuery(q.Get("year"), "year")
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	sortDir, err := parseIntQuery(q.Get("dir"), -1, "dir")
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	sortBy := strings.TrimSpace(q.Get("sort"))
	if sortBy == "" {
		sortBy = "date"
	}
	params := ListParams{
		UW:        q.Get("uw"),
		LeadOnly:  leadOnly,
		ValBucket: q.Get("val_bucket"),
		Year:      yearPtr,
		Query:     q.Get("q"),
		SortBy:    sortBy,
		SortDir:   sortDir,
	}
	yearStr := "None"
	if yearPtr != nil {
		yearStr = strconv.Itoa(*yearPtr)
	}
	cacheKey := fmt.Sprintf(
		"ipo_hunter_list_%s_%t_%s_%s_%s_%s_%d",
		params.UW, params.LeadOnly, params.ValBucket, yearStr, params.Query, params.SortBy, params.SortDir,
	)

	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 120*time.Second, func(ctx context.Context) (any, int, string, error) {
		masters, err := h.Repo.ListMasters(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		return map[string]any{
			"success": true,
			"data":    BuildIPOListData(masters, params),
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) GetAnalysis(w http.ResponseWriter, r *http.Request) {
	tk := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "ticker")))
	q := r.URL.Query()
	mode := strings.TrimSpace(q.Get("mode"))
	if mode == "" {
		mode = "after"
	}
	if mode != "after" && mode != "periodic" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "String should match pattern '^(after|periodic)$'")
		return
	}
	autoRetail, err := parseBoolQuery(q.Get("auto_retail"), true, "auto_retail")
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	dateFromRaw := strings.TrimSpace(q.Get("date_from"))
	dateToRaw := strings.TrimSpace(q.Get("date_to"))
	retailOverride := q.Get("retail_override")

	var reqFrom *time.Time
	if dateFromRaw != "" {
		dt, ok := ParseISODate(dateFromRaw)
		if !ok {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "Invalid date_from format, expected YYYY-MM-DD")
			return
		}
		reqFrom = &dt
	}
	var reqTo *time.Time
	if dateToRaw != "" {
		dt, ok := ParseISODate(dateToRaw)
		if !ok {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "Invalid date_to format, expected YYYY-MM-DD")
			return
		}
		reqTo = &dt
	}

	cacheKey := fmt.Sprintf(
		"ipo_hunter_analysis_%s_%s_%s_%s_%t_%s",
		tk, mode, dateFromRaw, dateToRaw, autoRetail, retailOverride,
	)

	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 120*time.Second, func(ctx context.Context) (any, int, string, error) {
		master, err := h.Repo.GetMasterByTicker(ctx, tk)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		if master == nil {
			return nil, http.StatusNotFound, fmt.Sprintf("Kode %s tidak ada di daftar IPO (S.csv). Pilih kode emiten dari IPO List.", tk), nil
		}

		sortedDaily, err := h.Repo.ListDailyByTicker(ctx, tk)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		d0Row, err := h.Repo.GetD0ByTicker(ctx, tk)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		holidays, maxHolYr, err := h.Repo.ListHolidays(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}

		hiDate := reqTo
		if hiDate == nil && len(sortedDaily) > 0 {
			dt := sortedDaily[len(sortedDaily)-1].Date
			hiDate = &dt
		}

		var afterInfo *AfterInfo
		var analysisDays []DayRow
		var chartDays []DayRow

		if mode == "after" {
			ai := BuildAfterDays(tk, master, d0Row, sortedDaily, hiDate, holidays)
			afterInfo = &ai
			analysisDays = ai.Days

			var loChart *time.Time
			if len(sortedDaily) > 0 {
				dt := sortedDaily[0].Date
				loChart = &dt
				if master.ListingDate != nil {
					for i := range sortedDaily {
						if !sortedDaily[i].Date.Before(*master.ListingDate) {
							d := sortedDaily[i].Date
							loChart = &d
							break
						}
					}
				}
			}
			chartDays = make([]DayRow, 0, len(sortedDaily))
			for _, d := range sortedDaily {
				if (loChart == nil || !d.Date.Before(*loChart)) && (hiDate == nil || !d.Date.After(*hiDate)) {
					chartDays = append(chartDays, d)
				}
			}
		} else {
			loDate := reqFrom
			if loDate == nil && len(sortedDaily) > 0 {
				dt := sortedDaily[0].Date
				loDate = &dt
			}
			if loDate != nil && hiDate != nil && loDate.After(*hiDate) {
				loDate, hiDate = hiDate, loDate
			}
			analysisDays = make([]DayRow, 0, len(sortedDaily))
			for _, d := range sortedDaily {
				if (loDate == nil || !d.Date.Before(*loDate)) && (hiDate == nil || !d.Date.After(*hiDate)) {
					analysisDays = append(analysisDays, d)
				}
			}
			chartDays = analysisDays
		}

		agg := BuildAggregation(tk, analysisDays)
		if agg == nil {
			cname := tk
			if master.CompanyName != nil && *master.CompanyName != "" {
				cname = *master.CompanyName
			}
			var ldStr *string
			if master.ListingDate != nil {
				s := FormatDateISO(*master.ListingDate)
				ldStr = &s
			}
			return map[string]any{
				"success": true,
				"data": AnalysisNoDataResponse{
					Ticker:  tk,
					Mode:    mode,
					HasData: false,
					Master: AnalysisNoDataMeta{
						Ticker:      tk,
						CompanyName: cname,
						ListingDate: ldStr,
						IPOPrice:    master.IPOPrice,
					},
					Candles: []CandleItem{},
				},
			}, http.StatusOK, "", nil
		}

		uwSet := UWSetOf(master)
		dayCount := 0
		for i := range analysisDays {
			dc := analysisDays[i].DayCount
			if dc <= 0 {
				dc = 1
			}
			dayCount += dc
		}
		decision := ComputeDecision(agg, dayCount, uwSet, autoRetail, retailOverride)
		cards := BuildAnalysisCards(tk, mode, master, agg, decision, analysisDays, afterInfo, holidays, maxHolYr)

		candles := make([]CandleItem, 0, len(chartDays))
		for i := range chartDays {
			d := &chartDays[i]
			if d.Close == nil {
				continue
			}
			candles = append(candles, CandleItem{
				Date:      FormatDateISO(d.Date),
				Open:      d.Open,
				High:      d.High,
				Low:       d.Low,
				Close:     d.Close,
				Prev:      d.Prev,
				Tavg:      d.Tavg,
				Volume:    d.Volume,
				NetVolume: d.NetVolume,
			})
		}

		availDates := make([]string, 0, len(sortedDaily))
		for i := range sortedDaily {
			availDates = append(availDates, FormatDateISO(sortedDaily[i].Date))
		}

		var periodStart, periodEnd *string
		if mode == "after" && afterInfo != nil && len(analysisDays) > 0 {
			startDt := analysisDays[0].Date
			if afterInfo.ListingDate != nil {
				startDt = *afterInfo.ListingDate
			} else if afterInfo.D0Date != nil {
				startDt = *afterInfo.D0Date
			}
			s := FormatDateISO(startDt)
			periodStart = &s
		} else if len(analysisDays) > 0 {
			s := FormatDateISO(analysisDays[0].Date)
			periodStart = &s
		}
		if mode == "after" && afterInfo != nil && afterInfo.LastDate != nil {
			s := FormatDateISO(*afterInfo.LastDate)
			periodEnd = &s
		} else if len(analysisDays) > 0 {
			s := FormatDateISO(analysisDays[len(analysisDays)-1].Date)
			periodEnd = &s
		}

		return map[string]any{
			"success": true,
			"data": AnalysisResponse{
				Ticker:  tk,
				Mode:    mode,
				HasData: true,
				Period: AnalysisPeriod{
					Start:         periodStart,
					End:           periodEnd,
					DayCount:      dayCount,
					PriceChange:   agg.PriceChange,
					AvgDailyValue: agg.AvgDailyValue,
					TotalVolume:   agg.TotalVolume,
					TotNetVolume:  agg.TotNetVolume,
					TotBuy:        agg.TotBuy,
					TotSell:       agg.TotSell,
				},
				Candles:              candles,
				AvailableTickerDates: availDates,
				Cards:                cards,
			},
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}

func (h *Handler) PostScreener(w http.ResponseWriter, r *http.Request) {
	defaultQuick := 5
	req := ScreenerRequest{
		Quick:  &defaultQuick,
		Chg:    "all",
		Metric: "topPerf",
		Conds:  []ScreenerCondition{},
	}
	rawBody, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Unable to read request body")
		return
	}
	if len(strings.TrimSpace(string(rawBody))) > 0 {
		if err := json.Unmarshal(rawBody, &req); err != nil {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "Invalid JSON payload: "+err.Error())
			return
		}
	}
	if req.Chg == "" {
		req.Chg = "all"
	}
	if req.Metric == "" {
		req.Metric = "topPerf"
	}
	for i := range req.Conds {
		if req.Conds[i].Op == "" {
			req.Conds[i].Op = ">"
		}
	}

	var parsedFrom *time.Time
	if req.DateFrom != nil && strings.TrimSpace(*req.DateFrom) != "" {
		dt, ok := ParseISODate(*req.DateFrom)
		if !ok {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "Invalid date_from format")
			return
		}
		parsedFrom = &dt
	}
	var parsedTo *time.Time
	if req.DateTo != nil && strings.TrimSpace(*req.DateTo) != "" {
		dt, ok := ParseISODate(*req.DateTo)
		if !ok {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "Invalid date_to format")
			return
		}
		parsedTo = &dt
	}

	canonBytes, _ := json.Marshal(req)
	sum := sha256.Sum256(canonBytes)
	cacheKey := "ipo_hunter_screener_" + hex.EncodeToString(sum[:16])

	body, status, detail, err := shared.GetOrComputeCached(r.Context(), h.Cache, cacheKey, 60*time.Second, func(ctx context.Context) (any, int, string, error) {
		allDates, err := h.Repo.ListDistinctDailyDates(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		if len(allDates) == 0 {
			return map[string]any{
				"success": true,
				"data": map[string]any{
					"rows":   []ScreenerRow{},
					"n_days": 0,
					"lo":     nil,
					"hi":     nil,
				},
			}, http.StatusOK, "", nil
		}

		hiDate := allDates[len(allDates)-1]
		if parsedTo != nil {
			hiDate = *parsedTo
		}
		var loDate time.Time
		if parsedFrom != nil {
			loDate = *parsedFrom
		} else {
			qn := 5
			if req.Quick != nil && *req.Quick != 0 {
				qn = *req.Quick
			}
			idxHi := len(allDates) - 1
			for i, d := range allDates {
				if d.Equal(hiDate) {
					idxHi = i
					break
				}
			}
			loDate = allDates[max(0, idxHi-qn+1)]
		}
		if loDate.After(hiDate) {
			loDate, hiDate = hiDate, loDate
		}

		nDays := 0
		for _, d := range allDates {
			if !d.Before(loDate) && !d.After(hiDate) {
				nDays++
			}
		}

		masters, err := h.Repo.ListMasters(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}
		holidays, _, err := h.Repo.ListHolidays(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err.Error(), err
		}

		masterByTicker := make(map[string]*MasterRecord, len(masters))
		filteredMasters := make([]*MasterRecord, 0, len(masters))
		var tickerFilter []string
		for i := range masters {
			m := &masters[i]
			if req.Year != nil && (m.ListingDate == nil || m.ListingDate.UTC().Year() != *req.Year) {
				continue
			}
			masterByTicker[m.Ticker] = m
			filteredMasters = append(filteredMasters, m)
			if req.Year != nil {
				tickerFilter = append(tickerFilter, m.Ticker)
			}
		}

		rowByTicker := make(map[string]ScreenerRow, len(filteredMasters))
		if req.After {
			d0Map, err := h.Repo.ListAllD0Map(ctx)
			if err != nil {
				return nil, http.StatusInternalServerError, err.Error(), err
			}
			hiCopy := hiDate
			err = h.Repo.StreamDailyAfterIPO(ctx, hiDate, tickerFilter, func(tk string, laterDays []DayRow) {
				m, ok := masterByTicker[tk]
				if !ok {
					return
				}
				if row := ComputeScreenerRow(tk, m, d0Map[tk], laterDays, nil, true, &hiCopy, holidays); row != nil {
					rowByTicker[tk] = *row
				}
			})
			if err != nil {
				return nil, http.StatusInternalServerError, err.Error(), err
			}
			// Ensure tickers with a D0 baseline but zero later rows <= hiDate are still scored.
			for _, m := range filteredMasters {
				if _, already := rowByTicker[m.Ticker]; already {
					continue
				}
				if d0Row := d0Map[m.Ticker]; d0Row != nil {
					if row := ComputeScreenerRow(m.Ticker, m, d0Row, nil, nil, true, &hiCopy, holidays); row != nil {
						rowByTicker[m.Ticker] = *row
					}
				}
			}
		} else {
			baseVols, err := h.Repo.ListBaseVolumesBefore(ctx, loDate, tickerFilter)
			if err != nil {
				return nil, http.StatusInternalServerError, err.Error(), err
			}
			err = h.Repo.StreamDailyPeriodic(ctx, loDate, hiDate, tickerFilter, func(tk string, periodDays []DayRow) {
				m, ok := masterByTicker[tk]
				if !ok {
					return
				}
				if row := ComputeScreenerRow(tk, m, nil, periodDays, baseVols[tk], false, nil, holidays); row != nil {
					rowByTicker[tk] = *row
				}
			})
			if err != nil {
				return nil, http.StatusInternalServerError, err.Error(), err
			}
		}

		// Reassemble in exact master order before filtering & stable sort.
		candidates := make([]ScreenerRow, 0, len(rowByTicker))
		for _, m := range filteredMasters {
			if row, ok := rowByTicker[m.Ticker]; ok {
				candidates = append(candidates, row)
			}
		}

		ordered, activeCol, sk, sdir := FilterAndSortScreenerRows(candidates, req.Chg, req.Conds, req.Metric, req.SortKey, req.Dir)
		loStr := FormatDateISO(loDate)
		hiStr := FormatDateISO(hiDate)

		return map[string]any{
			"success": true,
			"data": ScreenerResponse{
				Rows:      ordered,
				Total:     len(ordered),
				NDays:     nDays,
				Lo:        &loStr,
				Hi:        &hiStr,
				After:     req.After,
				ActiveCol: activeCol,
				SortKey:   sk,
				SortDir:   sdir,
			},
		}, http.StatusOK, "", nil
	})
	if err != nil || (status != 0 && status != http.StatusOK) {
		httpx.WriteError(w, status, detail)
		return
	}
	shared.WriteCachedTerminalJSON(w, body)
}
