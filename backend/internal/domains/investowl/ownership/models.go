package ownership

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var monthNamesID = [12]string{
	"Jan", "Feb", "Mar", "Apr", "Mei", "Jun",
	"Jul", "Agu", "Sep", "Okt", "Nov", "Des",
}

var monthlyColors = map[string]string{
	"foreign_named":    "#3A7BD5",
	"foreign_hidden":   "#2A5BA0",
	"retail":           "#1DB87A",
	"local_non_retail": "#8C6A00",
}

var snapshotColors = map[string]string{
	"controlling": "#8C5000",
	"substantial": "#3A7BD5",
	"public":      "#1DB87A",
	"treasury":    "#9F8D83",
	"management":  "#2A5BA0",
	"other":       "#6B7280",
}

var ctCompositionColors = map[string]string{
	"foreign_named":    "#3A7BD5",
	"foreign_hidden":   "#2A5BA0",
	"retail":           "#1DB87A",
	"local_non_retail": "#2A2E3D",
}

const (
	publicFilterSQL   = "(name ILIKE 'Masyarakat%' OR name ILIKE 'Publik%')"
	treasuryFilterSQL = "name ILIKE '%Treasury%'"
)

// roundHalfToEven matches CPython built-in round(v, digits) rounding semantics via decimal formatting for digits >= 0.
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
	// strconv.FormatFloat uses Ryū/dtoa shortest/exact decimal rounding matching CPython _Py_dg_dtoa.
	return roundFromFormat(v, digits)
}

func roundEvenPtr(v *float64, digits int) *float64 {
	if v == nil {
		return nil
	}
	r := roundEven(*v, digits)
	return &r
}

func pctChange(now, before *float64) *float64 {
	if now == nil || before == nil || *before <= 0 {
		return nil
	}
	r := roundEven(((*now-*before)/(*before))*100.0, 2)
	return &r
}

func periodLabel(d *time.Time) string {
	if d == nil {
		return "-"
	}
	m := int(d.Month())
	if m < 1 || m > 12 {
		return d.Format("2006-01-02")
	}
	return fmt.Sprintf("%s %d", monthNamesID[m-1], d.Year())
}

func classifyHolder(name string, pct float64, isControlling bool) string {
	clean := strings.ToLower(strings.TrimSpace(name))
	if isControlling {
		return "controlling"
	}
	if strings.Contains(clean, "treasury") {
		return "treasury"
	}
	if strings.HasPrefix(clean, "masyarakat") || strings.HasPrefix(clean, "publik") {
		return "public"
	}
	if pct >= 5.0 {
		return "substantial"
	}
	return "management"
}

// MarketCapProvenance mirrors ownership_sources.market_cap_provenance.
type MarketCapProvenance struct {
	Value       *float64 `json:"value"`
	Source      *string  `json:"source"`
	Basis       string   `json:"basis"`
	AsOf        *string  `json:"as_of"`
	PriceAsOf   *string  `json:"price_as_of"`
	SharesAsOf  *string  `json:"shares_as_of"`
	Independent bool     `json:"independent"`
}

func marketCapProvenance(
	priceClose *float64,
	priceDate *string,
	priceSource *string,
	listedShares *float64,
	sharesAsOf *string,
	sharesSource *string,
	ctMarketCapB *float64,
) MarketCapProvenance {
	if priceClose != nil && *priceClose != 0 && listedShares != nil && *listedShares != 0 {
		val := *priceClose * *listedShares
		pSrc := "unknown"
		if priceSource != nil && *priceSource != "" {
			pSrc = *priceSource
		}
		sSrc := "unknown"
		if sharesSource != nil && *sharesSource != "" {
			sSrc = *sharesSource
		}
		src := fmt.Sprintf("%s x %s", pSrc, sSrc)
		ind := priceSource != nil && *priceSource != "" &&
			sharesSource != nil && *sharesSource != "" &&
			!strings.Contains(*priceSource+*sharesSource, "ct_baseline")
		return MarketCapProvenance{
			Value:       &val,
			Source:      &src,
			Basis:       "computed",
			AsOf:        priceDate,
			PriceAsOf:   priceDate,
			SharesAsOf:  sharesAsOf,
			Independent: ind,
		}
	}
	if ctMarketCapB != nil {
		val := *ctMarketCapB * 1_000_000_000.0
		src := "ct_baseline.ticker_meta"
		return MarketCapProvenance{
			Value:       &val,
			Source:      &src,
			Basis:       "reported",
			AsOf:        nil,
			PriceAsOf:   nil,
			SharesAsOf:  nil,
			Independent: false,
		}
	}
	return MarketCapProvenance{
		Value:       nil,
		Source:      nil,
		Basis:       "unavailable",
		AsOf:        nil,
		PriceAsOf:   priceDate,
		SharesAsOf:  sharesAsOf,
		Independent: false,
	}
}

// --- Core Database Models ---

type TickerMeta struct {
	Ticker string `json:"ticker"`
	Name   string `json:"name"`
	Sector string `json:"sector"`
}

type PriceSnapshot struct {
	Date            *string  `json:"date"`
	Close           float64  `json:"close"`
	PreviousDate    *string  `json:"previous_date"`
	PreviousClose   *float64 `json:"previous_close"`
	Change1DPct     *float64 `json:"change_1d_pct"`
	MonthStartDate  *string  `json:"month_start_date"`
	MonthStartClose *float64 `json:"month_start_close"`
	MTDPct          *float64 `json:"mtd_pct"`
	Source          string   `json:"source"`
}

type CTPriceSnapshot struct {
	Date            *string  `json:"date"`
	Close           float64  `json:"close"`
	PreviousClose   *float64 `json:"previous_close"`
	Change1DPct     *float64 `json:"change_1d_pct"`
	MonthStartClose *float64 `json:"month_start_close"`
	MTDPct          *float64 `json:"mtd_pct"`
	Source          string   `json:"source"`
}

type HolderAggregate struct {
	HolderCount             int     `json:"holder_count"`
	TotalPct                float64 `json:"total_pct"`
	ControllingPct          float64 `json:"controlling_pct"`
	EffectiveControllingPct float64 `json:"effective_controlling_pct"`
	PublicPct               float64 `json:"public_pct"`
	TreasuryPct             float64 `json:"treasury_pct"`
	SubstantialPct          float64 `json:"substantial_pct"`
	ManagementPct           float64 `json:"management_pct"`
	MaxHolderPct            float64 `json:"max_holder_pct"`
	ControllingCount        int     `json:"controlling_count"`
	PublicCount             int     `json:"public_count"`
	TreasuryCount           int     `json:"treasury_count"`
	SubstantialCount        int     `json:"substantial_count"`
	FreeFloatEst            float64 `json:"free_float_est"`
}

type MonthlySnapshot struct {
	PeriodDate           *string  `json:"period_date"`
	PeriodLabel          string   `json:"period_label"`
	TotalSID             int      `json:"total_sid"`
	SIDChangeMoM         int      `json:"sid_change_mom"`
	SIDChangePct         float64  `json:"sid_change_pct"`
	ForeignNamedPct      float64  `json:"foreign_named_pct"`
	ForeignHiddenPct     float64  `json:"foreign_hidden_pct"`
	RetailPct            float64  `json:"retail_pct"`
	LocalNonRetailPct    float64  `json:"local_non_retail_pct"`
	CurrentPrice         *float64 `json:"current_price"`
	AvgAccumulationPrice *float64 `json:"avg_accumulation_price"`
	UnderwaterPct        *float64 `json:"underwater_pct"`
	EstimatedNyangkutPct *float64 `json:"estimated_nyangkut_pct"`
	RARScore             *float64 `json:"rar_score"`
	PBIScore             *float64 `json:"pbi_score"`
	IsHeroBadge          bool     `json:"is_hero_badge"`
}

type MonthlyHistoryItem struct {
	PeriodDate           *string  `json:"period_date"`
	PeriodLabel          string   `json:"period_label"`
	TotalSID             int      `json:"total_sid"`
	SIDChangeMoM         int      `json:"sid_change_mom"`
	SIDChangePct         float64  `json:"sid_change_pct"`
	ForeignTotalPct      float64  `json:"foreign_total_pct"`
	ForeignNamedPct      float64  `json:"foreign_named_pct"`
	ForeignHiddenPct     float64  `json:"foreign_hidden_pct"`
	RetailPct            float64  `json:"retail_pct"`
	LocalNonRetailPct    float64  `json:"local_non_retail_pct"`
	CurrentPrice         *float64 `json:"current_price"`
	UnderwaterPct        *float64 `json:"underwater_pct"`
	EstimatedNyangkutPct *float64 `json:"estimated_nyangkut_pct"`
	IsHeroBadge          bool     `json:"is_hero_badge"`
}

type MonthlySliceBrief struct {
	PeriodLabel       string  `json:"period_label"`
	TotalSID          int     `json:"total_sid"`
	RetailPct         float64 `json:"retail_pct"`
	ForeignNamedPct   float64 `json:"foreign_named_pct"`
	ForeignHiddenPct  float64 `json:"foreign_hidden_pct"`
	LocalNonRetailPct float64 `json:"local_non_retail_pct"`
}

type CompositionItem struct {
	Key         string  `json:"key"`
	Label       string  `json:"label"`
	Pct         float64 `json:"pct"`
	Color       string  `json:"color"`
	Description string  `json:"description"`
}

type CompositionBlock struct {
	Basis       string            `json:"basis"`
	PeriodDate  *string           `json:"period_date"`
	PeriodLabel string            `json:"period_label"`
	TotalPct    float64           `json:"total_pct"`
	Items       []CompositionItem `json:"items"`
	Note        string            `json:"note"`
}

type HolderRow struct {
	Name          string  `json:"name"`
	Pct           float64 `json:"pct"`
	IsControlling bool    `json:"is_controlling"`
	Category      string  `json:"category"`
}

type CompanyShares struct {
	AsOf            *string `json:"as_of"`
	ListedShares    float64 `json:"listed_shares"`
	TradeableShares float64 `json:"tradeable_shares"`
	Source          string  `json:"source"`
}

type IndexHistoryPoint struct {
	Date      *string  `json:"date"`
	Close     float64  `json:"close"`
	Volume    int64    `json:"volume"`
	Value     int64    `json:"value"`
	MarketCap *float64 `json:"market_cap"`
}

// --- CT Baseline Repo Models ---

type CTTickerMeta struct {
	Ticker     string         `json:"ticker"`
	Issuer     *string        `json:"issuer"`
	Konglo     *string        `json:"konglo"`
	Badge      *string        `json:"badge"`
	Badges     []string       `json:"badges"`
	FreeFloat  *float64       `json:"free_float"`
	MktcapB    *float64       `json:"mktcap_b"`
	SecNum     *float64       `json:"sec_num"`
	Price      *float64       `json:"price"`
	SLRibuLot  *float64       `json:"sl_ribu_lot"`
	PctSL      *float64       `json:"pct_sl"`
	Pie        map[string]any `json:"pie"`
}

type CTHolderRow struct {
	Name   string   `json:"name"`
	LF     *string  `json:"lf"`
	SLMei  *float64 `json:"sl_mei"`
	PctMei *float64 `json:"pct_mei"`
	SLJun  *float64 `json:"sl_jun"`
	PctJun *float64 `json:"pct_jun"`
	Status *string  `json:"status"`
}

type CTPctPoint struct {
	MonthKey    string   `json:"month_key"`
	MonthIndex  *int     `json:"month_index"`
	PctAsing    *float64 `json:"pct_asing"`
	PctAsingLt1 *float64 `json:"pct_asing_lt1"`
	PctRitel    *float64 `json:"pct_ritel"`
}

type CTHolderAgg struct {
	HolderCount int     `json:"holder_count"`
	TotalPct    float64 `json:"total_pct"`
	MaxPct      float64 `json:"max_pct"`
}

type CTSIDPoint struct {
	Period    string   `json:"period"`
	SID       int      `json:"sid"`
	Change    *int     `json:"change"`
	ChangePct *float64 `json:"change_pct"`
}

type CTT5Row struct {
	Ticker    string   `json:"ticker"`
	Date      *string  `json:"date"`
	Broker    *string  `json:"broker"`
	Name      *string  `json:"name"`
	Action    *string  `json:"action"`
	ChangeLot *float64 `json:"change_lot"`
	ValueM    *float64 `json:"value_m"`
	LotAfter  *float64 `json:"lot_after"`
	PctJuli   *float64 `json:"pct_juli"`
}

type CTT5Summary struct {
	Count      int     `json:"count"`
	BuyCount   int     `json:"buy_count"`
	SellCount  int     `json:"sell_count"`
	BuyValueM  float64 `json:"buy_value_m"`
	SellValueM float64 `json:"sell_value_m"`
	NetValueM  float64 `json:"net_value_m"`
}

type CTHistoryItem struct {
	PeriodLabel       string   `json:"period_label"`
	PeriodDate        *string  `json:"period_date"`
	TotalSID          *int     `json:"total_sid"`
	SIDChangeMoM      *int     `json:"sid_change_mom"`
	SIDChangePct      *float64 `json:"sid_change_pct"`
	ForeignTotalPct   float64  `json:"foreign_total_pct"`
	ForeignNamedPct   float64  `json:"foreign_named_pct"`
	ForeignHiddenPct  float64  `json:"foreign_hidden_pct"`
	RetailPct         float64  `json:"retail_pct"`
	LocalNonRetailPct float64  `json:"local_non_retail_pct"`
}
