package ownership

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aryarifki/lazytrader/backend/internal/core/httpx"
)

// SearchCTInvestors handles GET /api/v1/ownership/investors/ct.
func (h *Handler) SearchCTInvestors(w http.ResponseWriter, r *http.Request) {
	qMap := r.URL.Query()
	queryStr := qMap.Get("q")
	limit := 12
	if raw := strings.TrimSpace(qMap.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 100")
			return
		}
		limit = n
	}

	ck := ownershipCacheKey("ownership_ct_investors_v1", fmt.Sprintf("%s|%d", queryStr, limit))
	if h.tryCache(w, r, ck, 1800) {
		return
	}
	res, err := searchCTInvestors(r.Context(), h.DB, queryStr, limit)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeCache(w, r, ck, res, 1800)
}

// GetCTInvestorDetail handles GET /api/v1/ownership/investors/ct/detail.
func (h *Handler) GetCTInvestorDetail(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) < 2 {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "name is required and must be at least 2 characters")
		return
	}

	ck := ownershipCacheKey("ownership_ct_investor_detail_v1", name)
	if h.tryCache(w, r, ck, 1800) {
		return
	}
	res, err := ctInvestorDetail(r.Context(), h.DB, strings.TrimSpace(name))
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if res == nil {
		httpx.WriteError(w, http.StatusNotFound, "Investor tidak ditemukan di baseline CT")
		return
	}
	h.storeCache(w, r, ck, res, 1800)
}

// GetCTExplorer handles GET /api/v1/ownership/explorer/ct.
func (h *Handler) GetCTExplorer(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	mode := q.Get("mode")
	if mode == "" {
		mode = "stocks"
	}
	if mode != "stocks" && mode != "investors" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "mode must be stocks or investors")
		return
	}
	category := q.Get("category")
	if category == "" {
		category = "kepemilikan"
	}
	if category != "kepemilikan" && category != "transaksi" && category != "advanced" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "category must be kepemilikan, transaksi, or advanced")
		return
	}
	condition := q.Get("condition")
	if condition == "" {
		condition = "bertambah"
	}
	mktcapTier := q.Get("mktcap_tier")
	if mktcapTier == "" {
		mktcapTier = "big"
	}
	if !containsStrSlice(mktcapTiers, mktcapTier) {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid mktcap_tier")
		return
	}
	transaksiRange := q.Get("transaksi_range")
	if transaksiRange == "" {
		transaksiRange = "<5M"
	}
	if !containsStrSlice(valueRanges, transaksiRange) {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "invalid transaksi_range")
		return
	}
	origin := q.Get("investor_origin")
	if origin == "" {
		origin = "Local"
	}
	if origin != "Local" && origin != "Foreign" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "investor_origin must be Local or Foreign")
		return
	}
	sortKey := q.Get("sort")
	if sortKey == "" {
		sortKey = "default"
	}
	order := q.Get("order")
	if order == "" {
		order = "desc"
	}
	if order != "asc" && order != "desc" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "order must be asc or desc")
		return
	}
	limit := 50
	if raw := strings.TrimSpace(q.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 200 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "limit must be between 1 and 200")
			return
		}
		limit = n
	}
	offset := 0
	if raw := strings.TrimSpace(q.Get("offset")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 {
			httpx.WriteError(w, http.StatusUnprocessableEntity, "offset must be >= 0")
			return
		}
		offset = n
	}

	optStr := func(key string) *string {
		if !q.Has(key) {
			return nil
		}
		v := q.Get(key)
		return &v
	}

	ck := ownershipCacheKey("ownership_ct_explorer_v1", canonicalQueryString(q))
	if h.tryCache(w, r, ck, 900) {
		return
	}

	res, err := ctExplorer(r.Context(), h.DB, CTExplorerParams{
		Mode:            mode,
		Category:        category,
		Condition:       condition,
		Ranges:          optStr("ranges"),
		MktcapTier:      mktcapTier,
		TransaksiRange:  transaksiRange,
		InvestorOrigin:  origin,
		InvestorTypes:   optStr("investor_types"),
		RetailRanges:    optStr("retail_ranges"),
		PortfolioRanges: optStr("portfolio_ranges"),
		Sort:            sortKey,
		Order:           order,
		Limit:           limit,
		Offset:          offset,
	})
	if err != nil {
		if errors.Is(err, errCTValidation) {
			httpx.WriteError(w, http.StatusBadRequest, strings.TrimPrefix(err.Error(), errCTValidation.Error()+": "))
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.storeCache(w, r, ck, res, 900)
}
