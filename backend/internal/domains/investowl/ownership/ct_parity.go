package ownership

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

var idMonthFull = [12]string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

var typeFull = map[string]string{
	"ID": "Individual",
	"CP": "Corporate",
	"MF": "Mutual Fund",
	"SC": "Securities Company",
	"IS": "Insurance",
	"PF": "Pension Fund",
	"FD": "Foundation",
	"IB": "Financial Institution",
}

var (
	investorTypeOptionsLocal   = []string{"ID", "CP", "IB", "OTHER"}
	investorTypeOptionsForeign = []string{"CP", "IB", "ID", "MF", "OTHER"}
	investorSigCodesLocal      = map[string]bool{"ID": true, "CP": true, "IB": true}
	investorSigCodesForeign    = map[string]bool{"CP": true, "IB": true, "ID": true, "MF": true}

	ownershipRanges = []string{">10", "5-10", "1-5", "<1"}
	sidRanges       = []string{"<1", "1-3", "3-10", ">10"}
	valueRanges     = []string{"<5M", "5-100M", ">100M"}
	mktcapTiers     = []string{"micro", "small", "midsmall", "mid", "big"}
	retailRanges    = []string{"<1", "1-3", "3-5"}
	portfolioRanges = []string{"<3", "3-10", ">10"}
)

var mktcapLabels = map[string]string{
	"micro":    "Micro Cap (<300 M)",
	"small":    "Small Cap (300 M - 1 T)",
	"midsmall": "Mid-Small Cap (1 T - 5 T)",
	"mid":      "Mid Cap (5 T - 50 T)",
	"big":      "Big Cap (>50 T)",
}

var stockAdvancedOrder = []string{
	"naikPorsi", "aktifBeli", "buyback", "mayoritas", "baru",
	"retailTurun", "hiddenNaik", "foreignNaik", "masukFlat",
	"foreign3Bulan", "retail3Bulan", "sid3Turun", "sid3Naik", "retailRendah",
}

var stockAdvancedLabels = map[string]string{
	"naikPorsi":     "Big Player Menambah Porsi Bulan Ini",
	"aktifBeli":     "Pembelian Aktif oleh Big Player",
	"buyback":       "Emiten Buyback Bulan Ini",
	"mayoritas":     "Pemegang Saham Mayoritas Menambah Porsi",
	"baru":          "Masuknya Big Player Baru Bulan Ini",
	"retailTurun":   "Kepemilikan Retail Menurun",
	"hiddenNaik":    "Hidden Big Player Meningkat",
	"foreignNaik":   "Kepemilikan Foreign Meningkat",
	"masukFlat":     "Big Player Masuk (Harga Belum Naik Signifikan)",
	"foreign3Bulan": "Foreign Naik 3 Bulan Beruntun",
	"retail3Bulan":  "Retail Keluar 3 Bulan Beruntun",
	"sid3Turun":     "SID Turun 3 Bulan Beruntun",
	"sid3Naik":      "SID Naik 3 Bulan Beruntun",
	"retailRendah":  "Kepemilikan Retail Rendah",
}

var investorAdvancedOrder = []string{
	"portofolioNaik", "portofolioTurun", "investorAktif", "investorNambah",
}

var investorAdvancedLabels = map[string]string{
	"portofolioNaik":  "Portofolio Investor Naik Bulan Ini",
	"portofolioTurun": "Portofolio Investor Turun Bulan Ini",
	"investorAktif":   "Investor Aktif Bertransaksi Sebulan Terakhir",
	"investorNambah":  "Investor Menambah Porsi Sebulan Terakhir",
}

var conditionLabels = func() map[string]string {
	m := map[string]string{
		"bertambah": "Net Buy / Bertambah",
		"berkurang": "Net Sell / Berkurang",
		"sidNaik":   "SID Bertambah",
		"sidTurun":  "SID Berkurang",
	}
	for k, v := range stockAdvancedLabels {
		m[k] = v
	}
	for k, v := range investorAdvancedLabels {
		m[k] = v
	}
	return m
}()

var (
	rePersero  = regexp.MustCompile(`\(PERSERO\)`)
	reTbk      = regexp.MustCompile(`\bTBK\b`)
	rePT       = regexp.MustCompile(`\bPT\b`)
	rePunct    = regexp.MustCompile(`[.,]`)
)

func normalizeCompany(value string) string {
	s := strings.ToUpper(strings.TrimSpace(value))
	s = rePersero.ReplaceAllString(s, " ")
	s = reTbk.ReplaceAllString(s, " ")
	s = rePT.ReplaceAllString(s, " ")
	s = rePunct.ReplaceAllString(s, " ")
	return strings.Join(strings.Fields(s), " ")
}

func parseFullDate(value string) (time.Time, bool) {
	parts := strings.Fields(strings.TrimSpace(value))
	if len(parts) < 3 {
		return time.Time{}, false
	}
	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, false
	}
	var month int
	lowM := strings.ToLower(parts[1])
	for idx, nm := range idMonthFull {
		if strings.ToLower(nm) == lowM {
			month = idx + 1
			break
		}
	}
	if month == 0 {
		return time.Time{}, false
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, false
	}
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return time.Time{}, false
	}
	return t, true
}

func addMonth(d time.Time) time.Time {
	if d.Month() == time.December {
		return time.Date(d.Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	return time.Date(d.Year(), d.Month()+1, 1, 0, 0, 0, 0, time.UTC)
}

func glMonthLabel(meta map[string]any) string {
	raw, _ := meta["date_jun"].(string)
	d, ok := parseFullDate(raw)
	if !ok {
		return strings.ReplaceAll(strings.TrimSpace(raw), "30 ", "")
	}
	nextM := addMonth(d)
	return fmt.Sprintf("%s %d", idMonthFull[nextM.Month()-1], nextM.Year())
}

// --- CT Context Structures ---

type ctContextHolder struct {
	Ticker string
	Name   string
	LF     string
	SLMei  *float64
	PctMei *float64
	SLJun  *float64
	PctJun *float64
	Status string
}

type ctContextInvHolding struct {
	Name      string
	Ticker    string
	RawIssuer *string
	LF        string
	SLMei     *float64
	PctMei    *float64
	SLJun     *float64
	PctJun    *float64
	Status    string
}

type ctContextPct struct {
	PctAsing    *float64
	PctAsingLt1 *float64
	PctRitel    *float64
}

type ctContextTicker struct {
	Ticker     string
	Issuer     string
	Konglo     *string
	Badge      *string
	Badges     []string
	FreeFloat  *float64
	MktcapB    *float64
	SecNum     *float64
	Price      *float64
	SLRibuLot  *float64
	PctSL      *float64
	Rows       []ctContextHolder
	PctHistory map[string]ctContextPct
	SID        map[string]int
}

type ctContextOHLC struct {
	Date  *string
	Close *float64
}

type ctContextT5 struct {
	Ticker    string
	SortOrder int
	Date      *string
	DayNum    *int
	MonthNum  *int
	Broker    *string
	Name      string
	Action    string
	ChangeLot *float64
	ValueM    *float64
	LotAfter  *float64
	PctJuli   *float64
}

type ctContext struct {
	Meta             map[string]any
	MonthKeys        []string
	MonthLabels      map[string]string
	MonthPrev        string
	MonthCurr        string
	LastUpdateDate   any
	GLMonth          string
	NowDate          time.Time
	T5AnchorDate     time.Time
	GLYear           int
	GLMonthNum       int
	Year             int
	Tickers          map[string]*ctContextTicker
	TickerOrder      []string
	SIDPeriods       []string
	OHLC             map[string]ctContextOHLC
	T5ByTicker       map[string][]ctContextT5
	T5ByName         map[string][]ctContextT5
	T5NameOrder      []string
	InvestorHoldings map[string][]ctContextInvHolding
	InvestorOrder    []string
	InvestorIndex    map[string]float64
	IndexNameOrder   []string
}

func rowDeltaPct(pctMei, pctJun *float64, status string) float64 {
	st := strings.ToLower(strings.TrimSpace(status))
	pj := derefFloat(pctJun)
	pm := derefFloat(pctMei)
	if st == "new" {
		return pj
	}
	if st == "out" {
		return -pm
	}
	return pj - pm
}

func bucketPct(absVal float64) string {
	if absVal > 10 {
		return ">10"
	}
	if absVal >= 5 {
		return "5-10"
	}
	if absVal >= 1 {
		return "1-5"
	}
	return "<1"
}

func bucketSIDPct(absVal float64) string {
	if absVal > 10 {
		return ">10"
	}
	if absVal >= 3 {
		return "3-10"
	}
	if absVal >= 1 {
		return "1-3"
	}
	return "<1"
}

func bucketValue(absVal float64) string {
	if absVal > 100 {
		return ">100M"
	}
	if absVal >= 5 {
		return "5-100M"
	}
	return "<5M"
}

func bucketMktcap(val *float64) string {
	if val == nil || *val <= 0 {
		return ""
	}
	v := *val
	if v < 300 {
		return "micro"
	}
	if v < 1000 {
		return "small"
	}
	if v < 5000 {
		return "midsmall"
	}
	if v < 50000 {
		return "mid"
	}
	return "big"
}

func bucketRetail(v float64) string {
	if v < 1 {
		return "<1"
	}
	if v < 3 {
		return "1-3"
	}
	return "3-5"
}

func bucketPortfolio(absPct float64) string {
	if absPct < 3 {
		return "<3"
	}
	if absPct < 10 {
		return "3-10"
	}
	return ">10"
}

func passMktcap(entry *ctContextTicker, tier string) bool {
	b := bucketMktcap(entry.MktcapB)
	return b != "" && b == tier
}

func investorLFHoldings(holdings []ctContextInvHolding) string {
	for _, r := range holdings {
		if strings.TrimSpace(r.LF) != "" {
			return strings.TrimSpace(r.LF)
		}
	}
	return ""
}

func lfLabel(lf string) string {
	parts := strings.Fields(strings.TrimSpace(lf))
	if len(parts) < 2 {
		return lf
	}
	full, ok := typeFull[parts[1]]
	if !ok {
		full = parts[1]
	}
	return parts[0] + " " + full
}

func isSubsetSet(sub []string, sup map[string]bool) bool {
	for _, s := range sub {
		if !sup[s] {
			return false
		}
	}
	return true
}

func passInvestorType(lf, origin string, selectedTypes map[string]bool) bool {
	localAll := isSubsetSet(investorTypeOptionsLocal, selectedTypes)
	foreignAll := isSubsetSet(investorTypeOptionsForeign, selectedTypes)
	parts := strings.Fields(strings.TrimSpace(lf))
	if len(parts) < 2 {
		if origin == "Local" {
			return localAll
		}
		return foreignAll
	}
	typeOrigin, code := parts[0], parts[1]
	if typeOrigin != origin {
		return false
	}
	sig := investorSigCodesForeign
	if typeOrigin == "Local" {
		sig = investorSigCodesLocal
	}
	cat := "OTHER"
	if sig[code] {
		cat = code
	}
	if typeOrigin == "Local" && !localAll && !selectedTypes[cat] {
		return false
	}
	if typeOrigin == "Foreign" && !foreignAll && !selectedTypes[cat] {
		return false
	}
	return true
}

func t5Date(tx ctContextT5, year int) (time.Time, bool) {
	if tx.MonthNum != nil && *tx.MonthNum != 0 && tx.DayNum != nil && *tx.DayNum != 0 {
		m, d := *tx.MonthNum, *tx.DayNum
		t := time.Date(year, time.Month(m), d, 0, 0, 0, 0, time.UTC)
		if t.Year() == year && int(t.Month()) == m && t.Day() == d {
			return t, true
		}
		return time.Time{}, false
	}
	if tx.Date != nil {
		return parseFullDate(*tx.Date)
	}
	return time.Time{}, false
}

func loadCTContext(ctx context.Context, db *pgxpool.Pool) (*ctContext, error) {
	if !shared.CTAvailable(ctx, db) {
		return nil, nil
	}
	snap := ctSnapshotMeta(ctx, db)
	meta, _ := snap["meta"].(map[string]any)
	if meta == nil {
		meta = map[string]any{}
	}
	pctMonths, _ := meta["pct_months"].([]any)
	monthKeys := []string{}
	monthLabels := map[string]string{}
	for _, raw := range pctMonths {
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		k, _ := m["key"].(string)
		if k == "" {
			continue
		}
		lbl, _ := m["label"].(string)
		if lbl == "" {
			lbl = k
		}
		monthKeys = append(monthKeys, k)
		monthLabels[k] = lbl
	}

	nowDate := time.Now().UTC()
	nowDate = time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), 0, 0, 0, 0, time.UTC)
	if lu, ok := meta["last_update_date"].(string); ok {
		if d, ok2 := parseFullDate(lu); ok2 {
			nowDate = d
		} else if dj, ok3 := meta["date_jun"].(string); ok3 {
			if d2, ok4 := parseFullDate(dj); ok4 {
				nowDate = d2
			}
		}
	} else if dj, ok3 := meta["date_jun"].(string); ok3 {
		if d2, ok4 := parseFullDate(dj); ok4 {
			nowDate = d2
		}
	}
	var disclosureDate *time.Time
	if dj, ok := meta["date_jun"].(string); ok {
		if d, ok2 := parseFullDate(dj); ok2 {
			disclosureDate = &d
		}
	}

	tickers := map[string]*ctContextTicker{}
	tickerOrder := []string{}

	tmRows, err := db.Query(ctx, `
		SELECT ticker, issuer, konglo, badge, free_float::double precision,
		       mktcap_b::double precision, sec_num::double precision, price::double precision,
		       sl_ribu_lot::double precision, pct_sl::double precision,
		       COALESCE(badges, ARRAY[]::text[]), (badges IS NOT NULL) AS has_badges
		FROM ct_baseline.ticker_meta
	`)
	if err != nil {
		return nil, err
	}
	defer tmRows.Close()
	for tmRows.Next() {
		var t string
		var issuer, konglo, badge *string
		var ff, mc, secNum, price, sl, pctSL *float64
		var badges []string
		var hasBadges bool
		if scanErr := tmRows.Scan(&t, &issuer, &konglo, &badge, &ff, &mc, &secNum, &price, &sl, &pctSL, &badges, &hasBadges); scanErr != nil {
			return nil, scanErr
		}
		tUp := strings.ToUpper(t)
		iss := tUp
		if issuer != nil && *issuer != "" {
			iss = *issuer
		}
		var kongloPtr *string
		if konglo != nil && *konglo != "" && *konglo != "-" {
			k := *konglo
			kongloPtr = &k
		}
		var bList []string
		if hasBadges {
			if badges == nil {
				bList = []string{}
			} else {
				bList = badges
			}
		} else if badge != nil && *badge != "" {
			bList = []string{*badge}
		} else {
			bList = []string{}
		}
		entry := &ctContextTicker{
			Ticker:     tUp,
			Issuer:     iss,
			Konglo:     kongloPtr,
			Badge:      badge,
			Badges:     bList,
			FreeFloat:  roundEvenPtr(ff, 4),
			MktcapB:    roundEvenPtr(mc, 4),
			SecNum:     roundEvenPtr(secNum, 4),
			Price:      roundEvenPtr(price, 4),
			SLRibuLot:  roundEvenPtr(sl, 4),
			PctSL:      roundEvenPtr(pctSL, 4),
			Rows:       []ctContextHolder{},
			PctHistory: map[string]ctContextPct{},
			SID:        map[string]int{},
		}
		if _, exists := tickers[tUp]; !exists {
			tickerOrder = append(tickerOrder, tUp)
		}
		tickers[tUp] = entry
	}
	if err := tmRows.Err(); err != nil {
		return nil, err
	}

	hRows, err := db.Query(ctx, `
		SELECT ticker, name, lf, sl_mei::double precision, pct_mei::double precision,
		       sl_jun::double precision, pct_jun::double precision, status
		FROM ct_baseline.holders
		ORDER BY ticker, COALESCE(pct_jun, pct_mei) DESC NULLS LAST, name
	`)
	if err != nil {
		return nil, err
	}
	defer hRows.Close()
	for hRows.Next() {
		var t string
		var name, lf, status *string
		var slMei, pctMei, slJun, pctJun *float64
		if scanErr := hRows.Scan(&t, &name, &lf, &slMei, &pctMei, &slJun, &pctJun, &status); scanErr != nil {
			return nil, scanErr
		}
		tUp := strings.ToUpper(t)
		entry, ok := tickers[tUp]
		if !ok {
			continue
		}
		st := "both"
		if status != nil && strings.TrimSpace(*status) != "" {
			st = strings.ToLower(strings.TrimSpace(*status))
		}
		entry.Rows = append(entry.Rows, ctContextHolder{
			Ticker: tUp,
			Name:   cleanStr(name),
			LF:     cleanStr(lf),
			SLMei:  roundEvenPtr(slMei, 4),
			PctMei: roundEvenPtr(pctMei, 4),
			SLJun:  roundEvenPtr(slJun, 4),
			PctJun: roundEvenPtr(pctJun, 4),
			Status: st,
		})
	}
	if err := hRows.Err(); err != nil {
		return nil, err
	}

	pRows, err := db.Query(ctx, `
		SELECT ticker, month_key, month_index, pct_asing::double precision,
		       pct_asing_lt1::double precision, pct_ritel::double precision
		FROM ct_baseline.pct_history
		ORDER BY ticker, month_index ASC NULLS LAST, month_key ASC
	`)
	if err != nil {
		return nil, err
	}
	defer pRows.Close()
	pctKeySet := map[string]bool{}
	for pRows.Next() {
		var t, mk string
		var mi *int
		var pa, palt1, pr *float64
		if scanErr := pRows.Scan(&t, &mk, &mi, &pa, &palt1, &pr); scanErr != nil {
			return nil, scanErr
		}
		tUp := strings.ToUpper(t)
		entry, ok := tickers[tUp]
		if !ok {
			continue
		}
		pctKeySet[mk] = true
		entry.PctHistory[mk] = ctContextPct{
			PctAsing:    roundEvenPtr(pa, 4),
			PctAsingLt1: roundEvenPtr(palt1, 4),
			PctRitel:    roundEvenPtr(pr, 4),
		}
	}
	if err := pRows.Err(); err != nil {
		return nil, err
	}

	sRows, err := db.Query(ctx, `
		SELECT ticker, period, sid
		FROM ct_baseline.sid_history
		ORDER BY ticker, period ASC
	`)
	if err != nil {
		return nil, err
	}
	defer sRows.Close()
	sidPeriodSet := map[string]bool{}
	for sRows.Next() {
		var t, period string
		var sid int
		if scanErr := sRows.Scan(&t, &period, &sid); scanErr != nil {
			return nil, scanErr
		}
		tUp := strings.ToUpper(t)
		entry, ok := tickers[tUp]
		if !ok {
			continue
		}
		entry.SID[period] = sid
		sidPeriodSet[period] = true
	}
	if err := sRows.Err(); err != nil {
		return nil, err
	}
	sidPeriods := make([]string, 0, len(sidPeriodSet))
	for p := range sidPeriodSet {
		sidPeriods = append(sidPeriods, p)
	}
	sort.Strings(sidPeriods)

	oRows, err := db.Query(ctx, `SELECT ticker, date, close::double precision FROM ct_baseline.ohlc_latest`)
	if err != nil {
		return nil, err
	}
	defer oRows.Close()
	ohlc := map[string]ctContextOHLC{}
	for oRows.Next() {
		var t string
		var d *time.Time
		var cl *float64
		if scanErr := oRows.Scan(&t, &d, &cl); scanErr != nil {
			return nil, scanErr
		}
		ohlc[strings.ToUpper(t)] = ctContextOHLC{
			Date:  isoDatePtr(d),
			Close: roundEvenPtr(cl, 4),
		}
	}
	if err := oRows.Err(); err != nil {
		return nil, err
	}

	t5Rows, err := db.Query(ctx, `
		SELECT ticker, sort_order, date_label, day_num, month_num, broker,
		       name, action, change_lot::double precision, value_m::double precision,
		       lot_after::double precision, pct_juli::double precision
		FROM ct_baseline.t5_transactions
		ORDER BY month_num DESC NULLS LAST, day_num DESC NULLS LAST,
		         ticker ASC, sort_order ASC
	`)
	if err != nil {
		return nil, err
	}
	defer t5Rows.Close()
	t5ByTicker := make(map[string][]ctContextT5, len(tickers))
	for _, tk := range tickerOrder {
		t5ByTicker[tk] = []ctContextT5{}
	}
	t5ByName := map[string][]ctContextT5{}
	t5NameOrder := []string{}
	for t5Rows.Next() {
		var t string
		var sortOrd *int
		var dateLbl, broker, name, action *string
		var dayNum, monthNum *int
		var chgLot, valM, lotAfter, pctJuli *float64
		if scanErr := t5Rows.Scan(
			&t, &sortOrd, &dateLbl, &dayNum, &monthNum, &broker,
			&name, &action, &chgLot, &valM, &lotAfter, &pctJuli,
		); scanErr != nil {
			return nil, scanErr
		}
		tUp := strings.ToUpper(t)
		item := ctContextT5{
			Ticker:    tUp,
			SortOrder: derefInt(sortOrd),
			Date:      dateLbl,
			DayNum:    dayNum,
			MonthNum:  monthNum,
			Broker:    broker,
			Name:      cleanStr(name),
			Action:    strings.ToUpper(cleanStr(action)),
			ChangeLot: roundEvenPtr(chgLot, 4),
			ValueM:    roundEvenPtr(valM, 4),
			LotAfter:  roundEvenPtr(lotAfter, 4),
			PctJuli:   roundEvenPtr(pctJuli, 4),
		}
		t5ByTicker[tUp] = append(t5ByTicker[tUp], item)
		if item.Name != "" {
			if _, exists := t5ByName[item.Name]; !exists {
				t5NameOrder = append(t5NameOrder, item.Name)
			}
			t5ByName[item.Name] = append(t5ByName[item.Name], item)
		}
	}
	if err := t5Rows.Err(); err != nil {
		return nil, err
	}

	ihRows, err := db.Query(ctx, `
		SELECT name, ticker, issuer, lf, sl_mei::double precision, pct_mei::double precision,
		       sl_jun::double precision, pct_jun::double precision, status
		FROM ct_baseline.investor_holdings
		ORDER BY ctid
	`)
	if err != nil {
		return nil, err
	}
	defer ihRows.Close()
	invHoldings := map[string][]ctContextInvHolding{}
	invOrder := []string{}
	for ihRows.Next() {
		var namePtr, issuerPtr, lfPtr, statusPtr *string
		var ticker string
		var slMei, pctMei, slJun, pctJun *float64
		if scanErr := ihRows.Scan(&namePtr, &ticker, &issuerPtr, &lfPtr, &slMei, &pctMei, &slJun, &pctJun, &statusPtr); scanErr != nil {
			return nil, scanErr
		}
		nm := cleanStr(namePtr)
		if nm == "" {
			continue
		}
		st := "both"
		if statusPtr != nil && strings.TrimSpace(*statusPtr) != "" {
			st = strings.ToLower(strings.TrimSpace(*statusPtr))
		}
		if _, exists := invHoldings[nm]; !exists {
			invOrder = append(invOrder, nm)
		}
		invHoldings[nm] = append(invHoldings[nm], ctContextInvHolding{
			Name:      nm,
			Ticker:    strings.ToUpper(ticker),
			RawIssuer: issuerPtr,
			LF:        cleanStr(lfPtr),
			SLMei:     roundEvenPtr(slMei, 4),
			PctMei:    roundEvenPtr(pctMei, 4),
			SLJun:     roundEvenPtr(slJun, 4),
			PctJun:    roundEvenPtr(pctJun, 4),
			Status:    st,
		})
	}
	if err := ihRows.Err(); err != nil {
		return nil, err
	}

	iiRows, err := db.Query(ctx, `SELECT name, total::double precision FROM ct_baseline.investor_index`)
	if err != nil {
		return nil, err
	}
	defer iiRows.Close()
	invIndex := map[string]float64{}
	indexNameOrder := []string{}
	for iiRows.Next() {
		var nmPtr *string
		var tot *float64
		if scanErr := iiRows.Scan(&nmPtr, &tot); scanErr != nil {
			return nil, scanErr
		}
		nm := cleanStr(nmPtr)
		if nm != "" {
			if _, exists := invIndex[nm]; !exists {
				indexNameOrder = append(indexNameOrder, nm)
			}
			invIndex[nm] = derefFloat(tot)
		}
	}
	if err := iiRows.Err(); err != nil {
		return nil, err
	}

	if len(monthKeys) == 0 {
		for k := range pctKeySet {
			monthKeys = append(monthKeys, k)
		}
		sort.Strings(monthKeys)
	}
	if len(monthLabels) == 0 {
		for _, k := range monthKeys {
			monthLabels[k] = k
		}
	}
	monthPrev, _ := meta["month_prev"].(string)
	monthPrev = strings.TrimSpace(monthPrev)
	if monthPrev == "" && len(monthKeys) >= 2 {
		monthPrev = monthLabels[monthKeys[len(monthKeys)-2]]
	}
	monthCurr, _ := meta["month_curr"].(string)
	monthCurr = strings.TrimSpace(monthCurr)
	if monthCurr == "" && len(monthKeys) >= 1 {
		monthCurr = monthLabels[monthKeys[len(monthKeys)-1]]
	}

	var maxT5Date *time.Time
	for _, list := range t5ByTicker {
		for _, tx := range list {
			d, ok := t5Date(tx, nowDate.Year())
			if ok && (maxT5Date == nil || d.After(*maxT5Date)) {
				dCopy := d
				maxT5Date = &dCopy
			}
		}
	}
	t5AnchorDate := nowDate
	if maxT5Date != nil && maxT5Date.Before(nowDate) {
		t5AnchorDate = *maxT5Date
	}
	glRefDate := t5AnchorDate
	if disclosureDate != nil {
		glRefDate = addMonth(*disclosureDate)
	}

	return &ctContext{
		Meta:             meta,
		MonthKeys:        monthKeys,
		MonthLabels:      monthLabels,
		MonthPrev:        monthPrev,
		MonthCurr:        monthCurr,
		LastUpdateDate:   meta["last_update_date"],
		GLMonth:          glMonthLabel(meta),
		NowDate:          nowDate,
		T5AnchorDate:     t5AnchorDate,
		GLYear:           glRefDate.Year(),
		GLMonthNum:       int(glRefDate.Month()),
		Year:             nowDate.Year(),
		Tickers:          tickers,
		TickerOrder:      tickerOrder,
		SIDPeriods:       sidPeriods,
		OHLC:             ohlc,
		T5ByTicker:       t5ByTicker,
		T5ByName:         t5ByName,
		T5NameOrder:      t5NameOrder,
		InvestorHoldings: invHoldings,
		InvestorOrder:    invOrder,
		InvestorIndex:    invIndex,
		IndexNameOrder:   indexNameOrder,
	}, nil
}

type investorGLResult struct {
	Modal      float64
	ModalAwal  float64
	PeakBasis  float64
	MarketNow  float64
	GL         float64
	GLPct      *float64
	ShareCount float64
}

type glTx struct {
	Date    time.Time
	ValueRp float64
	IsBuy   bool
}

func investorGL(c *ctContext, name string, holdings []ctContextInvHolding) *investorGLResult {
	nameLower := strings.ToLower(strings.TrimSpace(name))
	var modal, modalAwal, sharesNow, marketNow float64
	hasData := false
	monthTxs := []glTx{}
	seenTickers := map[string]bool{}

	for _, row := range holdings {
		tk := strings.ToUpper(row.Ticker)
		seenTickers[tk] = true
		entry := c.Tickers[tk]
		var priceLast *float64
		if entry != nil {
			priceLast = entry.Price
		}
		latest := c.OHLC[tk]
		var priceNow *float64
		if latest.Close != nil && *latest.Close >= 1 {
			priceNow = latest.Close
		} else {
			priceNow = priceLast
		}
		if priceLast == nil || *priceLast == 0 || priceNow == nil {
			continue
		}

		baseShares := derefFloat(row.SLJun)
		if strings.ToLower(strings.TrimSpace(row.Status)) == "out" {
			baseShares = 0.0
		}
		baseModal := baseShares * *priceLast
		modal += baseModal
		modalAwal += baseModal
		shares := baseShares

		for _, tx := range c.T5ByTicker[tk] {
			if strings.ToLower(strings.TrimSpace(tx.Name)) != nameLower {
				continue
			}
			txDate, ok := t5Date(tx, c.Year)
			if !ok || int(txDate.Month()) != c.GLMonthNum || txDate.Year() != c.GLYear {
				continue
			}
			chgShares := derefFloat(tx.ChangeLot) * 100.0
			shares += chgShares
			valRp := derefFloat(tx.ValueM) * 1_000_000_000.0
			isBuy := chgShares >= 0
			if isBuy {
				modal += valRp
			} else {
				modal -= valRp
			}
			monthTxs = append(monthTxs, glTx{Date: txDate, ValueRp: valRp, IsBuy: isBuy})
		}
		sharesNow += shares
		marketNow += shares * *priceNow
		hasData = true
	}

	invT5Rows := c.T5ByName[name]
	if len(invT5Rows) == 0 {
		for _, cand := range c.T5NameOrder {
			if strings.ToLower(cand) == nameLower {
				invT5Rows = c.T5ByName[cand]
				break
			}
		}
	}
	extraOrder := []string{}
	type extraTx struct {
		tx   ctContextT5
		date time.Time
	}
	extraByTicker := map[string][]extraTx{}
	for _, tx := range invT5Rows {
		tk := strings.ToUpper(tx.Ticker)
		if tk == "" || seenTickers[tk] {
			continue
		}
		txDate, ok := t5Date(tx, c.Year)
		if !ok || int(txDate.Month()) != c.GLMonthNum || txDate.Year() != c.GLYear {
			continue
		}
		if _, exists := extraByTicker[tk]; !exists {
			extraOrder = append(extraOrder, tk)
		}
		extraByTicker[tk] = append(extraByTicker[tk], extraTx{tx: tx, date: txDate})
	}
	for _, tk := range extraOrder {
		txs := extraByTicker[tk]
		entry := c.Tickers[tk]
		latest := c.OHLC[tk]
		var priceNow *float64
		if latest.Close != nil && *latest.Close >= 1 {
			priceNow = latest.Close
		} else if entry != nil {
			priceNow = entry.Price
		}
		if priceNow == nil {
			continue
		}
		shares := 0.0
		for _, item := range txs {
			chgShares := derefFloat(item.tx.ChangeLot) * 100.0
			shares += chgShares
			valRp := derefFloat(item.tx.ValueM) * 1_000_000_000.0
			isBuy := chgShares >= 0
			if isBuy {
				modal += valRp
			} else {
				modal -= valRp
			}
			monthTxs = append(monthTxs, glTx{Date: item.date, ValueRp: valRp, IsBuy: isBuy})
		}
		sharesNow += shares
		marketNow += shares * *priceNow
		hasData = true
	}

	sort.SliceStable(monthTxs, func(i, j int) bool {
		return monthTxs[i].Date.Before(monthTxs[j].Date)
	})
	var cash, minCash float64
	for _, tx := range monthTxs {
		if tx.IsBuy {
			cash -= tx.ValueRp
		} else {
			cash += tx.ValueRp
		}
		if cash < minCash {
			minCash = cash
		}
	}

	peakBasis := modalAwal + math.Abs(minCash)
	if !hasData || peakBasis <= 0 {
		return nil
	}
	gl := marketNow - modal
	var glPct *float64
	if peakBasis > 0 {
		v := gl / peakBasis * 100.0
		glPct = &v
	}
	return &investorGLResult{
		Modal:      modal,
		ModalAwal:  modalAwal,
		PeakBasis:  peakBasis,
		MarketNow:  marketNow,
		GL:         gl,
		GLPct:      glPct,
		ShareCount: sharesNow,
	}
}

func periodPayload(c *ctContext) map[string]any {
	return map[string]any{
		"last_update_date": c.LastUpdateDate,
		"month_prev":       c.MonthPrev,
		"month_curr":       c.MonthCurr,
		"gl_month":         c.GLMonth,
		"month_labels":     c.MonthLabels,
		"month_keys":       c.MonthKeys,
	}
}

func searchCTInvestors(ctx context.Context, db *pgxpool.Pool, query string, limit int) (map[string]any, error) {
	c, err := loadCTContext(ctx, db)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return map[string]any{
			"query":   query,
			"count":   0,
			"results": []any{},
			"source":  "unavailable",
		}, nil
	}

	term := strings.ToLower(strings.TrimSpace(query))
	seen := map[string]bool{}
	names := make([]string, 0, len(c.IndexNameOrder)+len(c.InvestorOrder))
	for _, n := range c.IndexNameOrder {
		if !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}
	for _, n := range c.InvestorOrder {
		if !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}

	type searchItem struct {
		Name        string
		LF          string
		LFLabel     string
		TickerCount int
		TotalPct    float64
		Total       *float64
		InIndex     bool
	}
	results := []searchItem{}
	for _, name := range names {
		if term != "" && !strings.Contains(strings.ToLower(name), term) {
			continue
		}
		holdings := c.InvestorHoldings[name]
		lf := investorLFHoldings(holdings)
		totalPct := 0.0
		for _, r := range holdings {
			if r.PctJun != nil {
				totalPct += *r.PctJun
			} else {
				totalPct += derefFloat(r.PctMei)
			}
		}
		var totalPtr *float64
		if tot, ok := c.InvestorIndex[name]; ok {
			tCopy := tot
			totalPtr = &tCopy
		}
		results = append(results, searchItem{
			Name:        name,
			LF:          lf,
			LFLabel:     lfLabel(lf),
			TickerCount: len(holdings),
			TotalPct:    roundEven(totalPct, 4),
			Total:       totalPtr,
			InIndex:     totalPtr != nil,
		})
	}

	sort.SliceStable(results, func(i, j int) bool {
		a, b := results[i], results[j]
		aNull, bNull := a.Total == nil, b.Total == nil
		if aNull != bNull {
			return !aNull
		}
		aNeg := -derefFloat(a.Total)
		bNeg := -derefFloat(b.Total)
		if aNeg != bNeg {
			return aNeg < bNeg
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	if limit < 1 {
		limit = 1
	} else if limit > 100 {
		limit = 100
	}
	end := limit
	if len(results) < end {
		end = len(results)
	}
	limited := results[:end]
	hasMissing := false
	outRows := make([]map[string]any, 0, len(limited))
	for _, item := range limited {
		if !item.InIndex {
			hasMissing = true
		}
		outRows = append(outRows, map[string]any{
			"name":         item.Name,
			"lf":           item.LF,
			"lf_label":     item.LFLabel,
			"ticker_count": item.TickerCount,
			"total_pct":    item.TotalPct,
			"total":        item.Total,
			"in_index":     item.InIndex,
		})
	}

	return map[string]any{
		"query":         query,
		"count":         len(outRows),
		"total_matches": len(results),
		"results":       outRows,
		"source":        "ct_baseline",
		"period":        periodPayload(c),
		"data_quality": map[string]any{
			"ct_baseline":                      true,
			"includes_names_missing_from_index": hasMissing,
			"notes": []string{
				"Autocomplete menggabungkan investor_index dan investor_holdings.",
				"Nama yang tidak ada di investor_index tetap dapat ditemukan.",
			},
		},
	}, nil
}

func ctInvestorDetail(ctx context.Context, db *pgxpool.Pool, name string) (map[string]any, error) {
	c, err := loadCTContext(ctx, db)
	if err != nil || c == nil {
		return nil, err
	}

	var canonicalName string
	var holdings []ctContextInvHolding
	if rows, ok := c.InvestorHoldings[name]; ok {
		canonicalName = name
		holdings = rows
	} else {
		low := strings.ToLower(name)
		for _, cand := range c.InvestorOrder {
			if strings.ToLower(cand) == low {
				canonicalName = cand
				holdings = c.InvestorHoldings[cand]
				break
			}
		}
	}
	if canonicalName == "" {
		return nil, nil
	}

	var totalValueM, buyValueM, sellValueM float64
	var buyCount, sellCount int
	enriched := make([]map[string]any, 0, len(holdings))

	for _, row := range holdings {
		tk := strings.ToUpper(row.Ticker)
		meta := c.Tickers[tk]
		var price float64
		if meta != nil && meta.Price != nil {
			price = *meta.Price
		}
		slMei := derefFloat(row.SLMei)
		slJun := derefFloat(row.SLJun)
		diff := slJun - slMei
		var valM *float64
		if slJun != 0 && price != 0 {
			v := slJun * price / 1_000_000_000.0
			valM = &v
			totalValueM += v
		}
		if diff > 0 {
			buyCount++
			if price != 0 {
				buyValueM += diff * price / 1_000_000_000.0
			}
		} else if diff < 0 {
			sellCount++
			if price != 0 {
				sellValueM += math.Abs(diff) * price / 1_000_000_000.0
			}
		}

		issuerDisplay := tk
		if row.RawIssuer != nil && *row.RawIssuer != "" {
			issuerDisplay = *row.RawIssuer
		} else if meta != nil && meta.Issuer != "" {
			issuerDisplay = meta.Issuer
		}
		rawIssuerStr := ""
		if row.RawIssuer != nil {
			rawIssuerStr = *row.RawIssuer
		}
		var priceOut any
		if price != 0 {
			priceOut = price
		}
		var mktcapBOut *float64
		if meta != nil {
			mktcapBOut = meta.MktcapB
		}

		deltaPct := derefFloat(row.PctJun) - derefFloat(row.PctMei)
		enriched = append(enriched, map[string]any{
			"name":         row.Name,
			"ticker":       tk,
			"issuer":       issuerDisplay,
			"lf":           row.LF,
			"sl_mei":       row.SLMei,
			"pct_mei":      row.PctMei,
			"sl_jun":       row.SLJun,
			"pct_jun":      row.PctJun,
			"status":       row.Status,
			"price":        priceOut,
			"market_cap_b": mktcapBOut,
			"value_m":      roundEvenPtr(valM, 4),
			"delta_lot":    roundEven(diff/100.0, 4),
			"delta_shares": roundEven(diff, 4),
			"delta_pct":    roundEven(deltaPct, 4),
			"is_treasury":  normalizeCompany(row.Name) == normalizeCompany(rawIssuerStr),
		})
	}

	gl := investorGL(c, canonicalName, holdings)
	invT5 := append([]ctContextT5(nil), c.T5ByName[canonicalName]...)
	if len(invT5) == 0 {
		low := strings.ToLower(canonicalName)
		for _, cand := range c.T5NameOrder {
			if strings.ToLower(cand) == low {
				invT5 = append(invT5, c.T5ByName[cand]...)
				break
			}
		}
	}
	sort.SliceStable(invT5, func(i, j int) bool {
		return strings.TrimSpace(invT5[i].Ticker) < strings.TrimSpace(invT5[j].Ticker)
	})
	sort.SliceStable(invT5, func(i, j int) bool {
		di, okI := t5Date(invT5[i], c.Year)
		dj, okJ := t5Date(invT5[j], c.Year)
		if !okI {
			di = time.Time{}
		}
		if !okJ {
			dj = time.Time{}
		}
		return di.After(dj)
	})

	transactions := make([]map[string]any, 0, len(invT5))
	for _, tx := range invT5 {
		transactions = append(transactions, map[string]any{
			"ticker":     tx.Ticker,
			"date":       tx.Date,
			"broker":     tx.Broker,
			"name":       tx.Name,
			"action":     tx.Action,
			"change_lot": roundEvenPtr(tx.ChangeLot, 4),
			"value_m":    roundEvenPtr(tx.ValueM, 4),
			"lot_after":  roundEvenPtr(tx.LotAfter, 4),
			"pct_juli":   roundEvenPtr(tx.PctJuli, 4),
		})
	}

	lf := investorLFHoldings(holdings)
	netValueM := buyValueM - sellValueM

	var glRpM, glPct, peakBasisM *float64
	if gl != nil {
		v := roundEven(gl.GL/1_000_000_000.0, 4)
		glRpM = &v
		glPct = roundEvenPtr(gl.GLPct, 4)
		pb := roundEven(gl.PeakBasis/1_000_000_000.0, 4)
		peakBasisM = &pb
	}

	return map[string]any{
		"name":     canonicalName,
		"lf":       lf,
		"lf_label": lfLabel(lf),
		"summary": map[string]any{
			"ticker_count": len(holdings),
			"total_value_m": roundEven(totalValueM, 4),
			"gl_rp_m":       glRpM,
			"gl_pct":        glPct,
			"buy_value_m":   roundEven(buyValueM, 4),
			"buy_count":     buyCount,
			"sell_value_m":  roundEven(sellValueM, 4),
			"sell_count":    sellCount,
			"net_value_m":   roundEven(netValueM, 4),
			"peak_basis_m":  peakBasisM,
		},
		"period":       periodPayload(c),
		"holdings":     enriched,
		"transactions": transactions,
		"t5_summary":   ctT5MapSummary(transactions),
		"data_quality": map[string]any{
			"periodized_holdings": true,
			"t5_transactions":     len(transactions) > 0,
			"gl":                  gl != nil,
			"market_cap":          true,
			"ct_baseline":         true,
			"notes": []string{
				"Portofolio, T5, dan G/L berasal dari baseline CT yang diisolasi di schema ct_baseline.",
				"G/L memakai basis peak-capital: modal awal + defisit kas terdalam pada transaksi bulan berjalan.",
				"Nilai transaksi bukan harga eksekusi; tidak ada biaya transaksi, pajak, atau dividen.",
			},
		},
	}, nil
}

func cleanStr(p *string) string {
	if p == nil {
		return ""
	}
	return strings.TrimSpace(*p)
}
