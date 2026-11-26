package shared

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	defaultWatchlist = []string{"BBCA", "BBRI", "BMRI", "BBNI", "TLKM", "ASII", "UNVR", "GOTO", "BREN", "ANTM"}

	idx30 = []string{
		"ADRO", "AMMN", "AMRT", "ANTM", "ARTO", "ASII", "BBCA", "BBNI",
		"BBRI", "BBTN", "BMRI", "BRMS", "BRPT", "BSDE", "BUKA", "CPIN",
		"CTRA", "ESSA", "GGRM", "GOTO", "HRUM", "ICBP", "INCO", "INDF",
		"INKP", "ISAT", "ITMG", "KLBF", "MAPI", "MBMA", "MDKA", "MEDC",
		"PGAS", "PTBA", "SMGR", "TLKM", "TOWR", "UNTR", "UNVR",
	}

	lq45 = []string{
		"ADRO", "AMRT", "ANTM", "ARTO", "ASII", "BBCA", "BBNI", "BBRI",
		"BBTN", "BMRI", "BRPT", "BSDE", "BUKA", "CPIN", "CTRA", "ERAA",
		"ESSA", "EXCL", "GGRM", "GOTO", "HEAL", "HRUM", "ICBP", "INCO",
		"INDF", "INKP", "INTP", "ISAT", "ITMG", "JPFA", "JSMR", "KLBF",
		"MAPI", "MBMA", "MDKA", "MEDC", "MIKA", "MYOR", "PGAS", "PTBA",
		"SMGR", "TLKM", "TOWR", "UNTR", "UNVR",
	}

	idx80 = []string{
		"AADI", "ACES", "ADMR", "ADRO", "AKRA", "AMMN", "AMRT", "ANTM",
		"ARTO", "ASII", "BBCA", "BBNI", "BBRI", "BBTN", "BFIN", "BKSL",
		"BMRI", "BRMS", "BRPT", "BSDE", "BUKA", "BUMI", "CBDK", "CMRY",
		"CPIN", "CTRA", "CUAN", "DEWA", "DSNG", "ELSA", "EMTK", "ENRG",
		"ERAA", "ESSA", "EXCL", "GGRM", "GOTO", "HEAL", "HRTA", "HRUM",
		"ICBP", "INCO", "INDF", "INDY", "INKP", "ISAT", "ITMG", "JPFA",
		"JSMR", "KIJA", "KLBF", "KPIG", "LSIP", "MAPA", "MAPI", "MBMA",
		"MDKA", "MEDC", "MIKA", "MYOR", "NCKL", "PGAS", "PGEO", "PNLF",
		"PTBA", "PTRO", "PWON", "RAJA", "RATU", "SCMA", "SMGR", "SMRA",
		"SSIA", "TAPG", "TLKM", "TOWR", "TPIA", "UNTR", "UNVR", "WIFI",
	}
)

func configuredWatchlist() []string {
	if raw := strings.TrimSpace(os.Getenv("WATCHLIST")); raw != "" {
		return uniqueSorted(strings.Split(raw, ","))
	}
	return uniqueSorted(defaultWatchlist)
}

func uniqueSorted(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, it := range items {
		s := strings.ToUpper(strings.TrimSpace(it))
		if s == "" {
			continue
		}
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}

// GetMasterTickers returns sorted tickers from public.tickers.
func GetMasterTickers(ctx context.Context, db *pgxpool.Pool, activeOnly bool) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database pool is nil")
	}
	q := "SELECT ticker FROM tickers"
	if activeOnly {
		q += " WHERE is_active = TRUE"
	}
	q += " ORDER BY ticker"
	rows, err := db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]string, 0, 128)
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err == nil {
			if s := strings.TrimSpace(t); s != "" {
				out = append(out, s)
			}
		}
	}
	return out, rows.Err()
}

// GetUniverse resolves a universe mode into a sorted ticker list matching universe.py.
func GetUniverse(ctx context.Context, db *pgxpool.Pool, mode string, customList []string) ([]string, error) {
	m := strings.ToLower(strings.TrimSpace(mode))
	if m == "" {
		m = "watchlist"
	}
	switch m {
	case "watchlist":
		return configuredWatchlist(), nil
	case "idx30":
		return uniqueSorted(idx30), nil
	case "lq45":
		return uniqueSorted(lq45), nil
	case "idx80":
		return uniqueSorted(idx80), nil
	case "custom":
		if len(customList) == 0 {
			return nil, fmt.Errorf("custom_list is required when mode='custom'")
		}
		return uniqueSorted(customList), nil
	case "all":
		if tickers, err := GetMasterTickers(ctx, db, true); err == nil && len(tickers) > 0 {
			return tickers, nil
		}
		return uniqueSorted(append(append([]string{}, idx80...), configuredWatchlist()...)), nil
	case "liquid":
		base := uniqueSorted(append(append([]string{}, idx80...), configuredWatchlist()...))
		if master, err := GetMasterTickers(ctx, db, true); err == nil && len(master) > 0 {
			mSet := make(map[string]struct{}, len(master))
			for _, t := range master {
				mSet[strings.ToUpper(t)] = struct{}{}
			}
			filtered := make([]string, 0, len(base))
			for _, t := range base {
				if _, ok := mSet[t]; ok {
					filtered = append(filtered, t)
				}
			}
			return filtered, nil
		}
		return base, nil
	default:
		return nil, fmt.Errorf("unknown universe mode: %s", m)
	}
}

var hardcodedIndices = map[string][]string{
	"idx_bumn":          {"ADHI", "ANTM", "BBNI", "BBRI", "BBTN", "BMRI", "BRIS", "ELSA", "JSMR", "MTEL", "PGAS", "PGEO", "PTBA", "PTPP", "SMGR", "TINS", "TLKM", "WIKA", "WSKT"},
	"idx_high_dividend": {"ADRO", "AMRT", "ANTM", "ASII", "BBNI", "BBRI", "BBCA", "BMRI", "BNGA", "BRPT", "EXCL", "HEXA", "HMSP", "INDF", "ITMG", "KLBF", "PTBA", "TLKM", "UNTR"},
	"esg_kehati":        {"AALI", "ADHI", "ASII", "BBCA", "BBNI", "BBRI", "BBTN", "BMRI", "BSDE", "INDF", "JSMR", "KLBF", "PGAS", "PTBA", "SMGR", "TLKM", "UNTR", "UNVR", "WIKA"},
	"bisnis-27":         {"ADRO", "AKRA", "AMRT", "ANTM", "ASII", "BBCA", "BBNI", "BBRI", "BMRI", "BRPT", "CPIN", "CTRA", "EXCL", "INKP", "ITMG", "JSMR", "KLBF", "MAPI", "MIKA", "PGAS", "PTBA", "SMGR", "TLKM", "TOWR", "UNTR"},
	"idx_smc":           {"ABMM", "ACES", "AGII", "AKRA", "AMFG", "ARNA", "ASSA", "AUTO", "BIRD", "BNGA", "BSDE", "CLEO", "CTRA", "DRMA", "DSNG", "ELSA", "ENRG", "ESSA", "HEAL", "HRUM", "IMAS", "INDY", "JPFA", "KEEN", "LSIP", "MARK", "MBAP", "MCOL", "MEDC", "MIKA", "MYOR", "NISP", "PANR", "PNLF", "PTRO", "RAJA", "SGER", "SIDO", "SMSM", "SSIA", "TAPG", "TOTL"},
	"idxenergy":         {"ADMR", "ADRO", "AKRA", "APEX", "BIPI", "BUMI", "BYAN", "CUAN", "DEWA", "DOID", "DSSA", "ELSA", "ENRG", "ESSA", "GEMS", "GTBO", "HRUM", "INDY", "ITMG", "KKGI", "KOPI", "MBAP", "MCOL", "MEDC", "PGAS", "PGEO", "PTBA", "PTRO", "RAJA", "RMKE", "SGER", "SMMT", "TEBE", "TOBA", "WINS"},
	"idxtechno":         {"AWAN", "BELI", "BUKA", "DCII", "DIVA", "EDGE", "ELIT", "EMTK", "GOTO", "KCI", "KIOS", "LUCK", "MCAS", "MLPT", "MTDL", "NFCX", "PTSN", "TECH", "TFAS", "WIFI", "ZATA"},
	"idxfinance":        {"ADMF", "AGRO", "AMOR", "ARTO", "BBCA", "BBHI", "BBKP", "BBNI", "BBRI", "BBTN", "BDMN", "BFIN", "BJBR", "BJTM", "BMRI", "BNBA", "BNGA", "BNII", "BNLI", "BRIS", "BTPN", "CFIN", "MEGA", "NISP", "NOBU", "PNBN", "PNBS", "PNLF", "POLA", "SDRA"},
	"idxhealth":         {"CARE", "DGNS", "HEAL", "IRRA", "KAEF", "KLBF", "MIKA", "OMNI", "PEHA", "PRDA", "RSGK", "SAME", "SIDO", "SILO", "SOHO"},
	"idxpropert":        {"APLN", "ASRI", "BAPA", "BEST", "BKSL", "BSDE", "CITY", "CTRA", "DILD", "DMAS", "DUTI", "GWNG", "KIJA", "LPCK", "LPKR", "MDLN", "MTLA", "PWON", "SMRA", "SSIA"},
	"idxinfra":          {"ADHI", "BALI", "BICC", "CASS", "CMNP", "EXCL", "FREN", "ISAT", "JSMR", "KEEN", "META", "MTEL", "PGAS", "PORT", "PTPP", "TBAL", "TBIG", "TLKM", "TOWR", "WIKA", "WSKT"},
	"idxindust":         {"ASII", "AUTO", "BIMA", "BMSR", "BUDI", "HEXA", "IMAS", "KBLI", "KOBX", "LION", "MARK", "SCCO", "SMSM", "UNTR"},
	"idxbasic":          {"ADMG", "AGII", "AMMN", "ANTM", "ARCI", "AVIA", "BRMS", "BRPT", "CPIN", "ESSA", "FASW", "INCO", "INKP", "INTP", "JPFA", "LTLS", "MDKA", "NCKL", "SMCB", "SMGR", "TPIA"},
	"idxcyclic":         {"ACES", "AMRT", "AUTO", "CSAP", "ERAA", "LPPF", "MAPI", "MAPA", "MSKY", "RALS"},
	"idxnoncyc":         {"AALI", "CINT", "CLEO", "CMRY", "CPIN", "DSNG", "GGRM", "HMSP", "ICBP", "INDF", "JPFA", "KINO", "LSIP", "MAIN", "MYOR", "ROTI", "SSMS", "STTP", "TBLA", "ULTJ", "UNVR"},
	"idxtrans":          {"ASSA", "BIRD", "BPTR", "CASS", "CMPP", "GIAA", "HAIS", "HITS", "IPCC", "KJEN", "NELY", "PORT", "PSSI", "SAPX", "TMAS", "TRUK"},
}

// GetDynamicUniverse mirrors bandarmology.py get_dynamic_universe with hardcoded indices and fallbacks.
func GetDynamicUniverse(ctx context.Context, db *pgxpool.Pool, mode string) []string {
	m := strings.ToLower(strings.TrimSpace(mode))
	switch m {
	case "watchlist", "idx30", "lq45", "idx80", "all", "liquid":
		if tickers, err := GetUniverse(ctx, db, m, nil); err == nil && len(tickers) > 0 {
			return tickers
		}
	}
	if list, ok := hardcodedIndices[m]; ok {
		return uniqueSorted(list)
	}
	return []string{"ANTM", "BBCA", "BBRI", "BMRI", "GOTO", "TLKM"}
}
