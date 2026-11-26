package shared

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConglomerateGroup represents a conglomerate group definition loaded from ownership_staging or ct_baseline.
type ConglomerateGroup struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	HubEntity    string   `json:"hub_entity"`
	LogoURL      *string  `json:"logo_url"`
	CatalogCode  *string  `json:"catalog_code"`
	KnownTickers []string `json:"known_tickers"`
	GroupName    string   `json:"group_name"`
	TickerCount  int      `json:"ticker_count"`
	Source       string   `json:"source"`
}

var groupDescriptions = map[string][2]string{
	"barito group":    {"Konglomerasi energi, petrokimia, dan tambang.", "Prajogo Pangestu"},
	"salim group":     {"Konglomerasi konsumen, agribisnis, ritel, dan infrastruktur.", "Anthoni Salim"},
	"astra group":     {"Konglomerasi otomotif, alat berat, agribisnis, dan jasa keuangan.", "PT Astra International Tbk"},
	"djarum group":    {"Konglomerasi perbankan, menara telekomunikasi, dan ritel digital.", "Robert Budi Hartono & Michael Hartono"},
	"sinarmas group":  {"Konglomerasi pulp & paper, agribisnis, energi, dan properti.", "Widjaja Family"},
	"sinarmas land":   {"Konglomerasi properti dan real estat.", "Widjaja Family"},
	"lippo group":     {"Konglomerasi properti, kesehatan, ritel, dan telekomunikasi.", "Riady Family"},
	"mnc group":       {"Konglomerasi media, hiburan, dan jasa keuangan.", "Hary Tanoesoedibjo"},
	"bakrie group":    {"Konglomerasi sumber daya alam, energi, infrastruktur, dan media.", "Bakrie & Brothers"},
	"thohir group":    {"Konglomerasi energi, pertambangan, dan investasi.", "Garibaldi Thohir"},
	"emtek group":     {"Konglomerasi media, teknologi digital, dan kesehatan.", "Eddy Kusnadi Sariaatmadja"},
	"ciputra group":   {"Konglomerasi properti, konstruksi, dan real estat.", "Ciputra Family"},
	"rajawali group":  {"Konglomerasi properti, agribisnis, dan infrastruktur.", "Peter Sondakh"},
	"panin group":     {"Konglomerasi jasa keuangan dan perbankan.", "Panin Family"},
	"saratoga group":  {"Konglomerasi investasi, energi, dan infrastruktur.", "Saratoga Investama Sedaya"},
	"tancorp group":   {"Konglomerasi konsumen, distribusi, dan investasi.", "Tancorp"},
	"aguan group":     {"Konglomerasi properti, infrastruktur, dan konsumen.", "Sugianto Kusuma"},
	"triputra group":  {"Konglomerasi agribisnis, energi, dan manufaktur.", "Triputra Group"},
	"ltk group":       {"Konglomerasi energi dan pertambangan.", "Low Tuck Kwong"},
	"mayapada group":  {"Konglomerasi kesehatan, properti, dan jasa keuangan.", "Dato Sri Tahir"},
	"alfa group":      {"Konglomerasi ritel, distribusi, dan konsumen.", "Alfa Group"},
	"ssi group":       {"Konglomerasi media, hiburan, dan investasi.", "SSI Group"},
	"haji isam group": {"Konglomerasi energi dan pertambangan.", "Haji Isam"},
	"hashim group":    {"Konglomerasi teknologi, telekomunikasi, dan investasi.", "Hashim Djojohadikusumo"},
	"northstar group": {"Konglomerasi jasa keuangan dan investasi.", "Northstar Group"},
	"harita group":    {"Konglomerasi pertambangan nikel dan bauksit.", "Lim Family"},
	"ct corp group":   {"Konglomerasi media, perbankan, dan konsumen.", "Chairul Tanjung"},
	"kresna group":    {"Konglomerasi jasa keuangan dan investasi.", "Kresna Group"},
	"hapsoro group":   {"Konglomerasi energi, pertambangan, dan infrastruktur.", "Hapsoro Group"},
}

var groupAliasMap = map[string]string{
	"barito":          "prajogo-pangestu",
	"barito-group":    "prajogo-pangestu",
	"prajogo":         "prajogo-pangestu",
	"djarum-group":    "djarum",
	"salim-group":     "salim",
	"astra-group":     "astra",
	"bakrie-group":    "bakrie",
	"sinarmas-group":  "sinarmas",
	"lippo-group":     "lippo",
	"emtek-group":     "emtek",
	"panin-group":     "panin",
	"saratoga-group":  "saratoga",
	"harita-group":    "harita",
	"mnc-group":       "mnc",
	"thohir-group":    "thohir",
	"haji-isam-group": "haji-isam",
	"danantara-group": "danantara",
	"triputra":        "tp-rachmat",
	"triputra-group":  "tp-rachmat",
	"bayan":           "low-tuck-kwong",
	"bayan-group":     "low-tuck-kwong",
	"ltk":             "low-tuck-kwong",
	"ltk-group":       "low-tuck-kwong",
}

var nonAlnumRe = regexp.MustCompile(`[^a-z0-9]+`)

// Slug converts a group name to its canonical URL slug matching ct_baseline._slug.
func Slug(value string) string {
	clean := strings.Trim(nonAlnumRe.ReplaceAllString(strings.ToLower(value), "-"), "-")
	if clean == "" {
		return "group"
	}
	return clean
}

// CTAvailable returns true when ct_baseline.ticker_meta exists in PostgreSQL.
func CTAvailable(ctx context.Context, db *pgxpool.Pool) bool {
	if db == nil {
		return false
	}
	var reg *string
	if err := db.QueryRow(ctx, "SELECT to_regclass('ct_baseline.ticker_meta')").Scan(&reg); err != nil {
		return false
	}
	return reg != nil && *reg != ""
}

func groupDescHub(name string) (string, string) {
	if pair, ok := groupDescriptions[strings.ToLower(name)]; ok {
		return pair[0], pair[1]
	}
	return "Konglomerasi " + name + ".", strings.TrimSuffix(name, " Group")
}

// LoadConglomerateGroups loads groups preserving SQL order, falling back to ct_baseline when staging is absent.
func LoadConglomerateGroups(ctx context.Context, db *pgxpool.Pool) ([]ConglomerateGroup, map[string]ConglomerateGroup) {
	ordered := []ConglomerateGroup{}
	byID := make(map[string]ConglomerateGroup)
	if db == nil {
		return ordered, byID
	}

	rows, err := db.Query(ctx, `
		SELECT v.group_name, count(*) AS ticker_count,
		       array_agg(v.ticker ORDER BY v.ticker) AS tickers,
		       COALESCE(s.logo_url, c.logo_url) AS logo_url,
		       COALESCE(s.catalog_code, c.catalog_code) AS catalog_code
		FROM ownership_staging.v_conglomerate_members v
		LEFT JOIN ownership_staging.stockbit_conglomerate_catalogs s ON s.slug = v.group_slug
		LEFT JOIN ownership_staging.conglomerate_groups c ON c.slug = v.group_slug
		GROUP BY v.group_name, COALESCE(s.logo_url, c.logo_url), COALESCE(s.catalog_code, c.catalog_code)
		ORDER BY count(*) DESC, v.group_name ASC
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			var cnt int
			var tickers []string
			var logoURL, catalogCode *string
			if scanErr := rows.Scan(&name, &cnt, &tickers, &logoURL, &catalogCode); scanErr != nil {
				continue
			}
			if tickers == nil {
				tickers = []string{}
			}
			base := Slug(name)
			groupID := base
			suffix := 2
			for _, exists := byID[groupID]; exists; _, exists = byID[groupID] {
				groupID = fmt.Sprintf("%s-%d", base, suffix)
				suffix++
			}
			desc, hub := groupDescHub(name)
			g := ConglomerateGroup{
				ID:           groupID,
				Name:         name,
				Description:  desc,
				HubEntity:    hub,
				LogoURL:      logoURL,
				CatalogCode:  catalogCode,
				KnownTickers: tickers,
				GroupName:    name,
				TickerCount:  len(tickers),
				Source:       "ownership_staging.v_conglomerate_members",
			}
			byID[groupID] = g
			ordered = append(ordered, g)
		}
		if len(ordered) > 0 {
			return ordered, byID
		}
	}

	if !CTAvailable(ctx, db) {
		return ordered, byID
	}
	fbRows, fbErr := db.Query(ctx, `
		SELECT g.group_name, COALESCE(g.ticker_count, 0)::int,
		       array_agg(m.ticker ORDER BY m.ticker) AS tickers
		FROM ct_baseline.groups g
		JOIN ct_baseline.group_members m ON m.group_name = g.group_name
		GROUP BY g.group_name, g.ticker_count
		ORDER BY g.ticker_count DESC, g.group_name ASC
	`)
	if fbErr != nil {
		return ordered, byID
	}
	defer fbRows.Close()
	for fbRows.Next() {
		var name string
		var cnt int
		var tickers []string
		if scanErr := fbRows.Scan(&name, &cnt, &tickers); scanErr != nil {
			continue
		}
		if tickers == nil {
			tickers = []string{}
		}
		base := Slug(name)
		groupID := base
		suffix := 2
		for _, exists := byID[groupID]; exists; _, exists = byID[groupID] {
			groupID = fmt.Sprintf("%s-%d", base, suffix)
			suffix++
		}
		desc, hub := groupDescHub(name)
		if cnt == 0 {
			cnt = len(tickers)
		}
		g := ConglomerateGroup{
			ID:           groupID,
			Name:         name,
			Description:  desc,
			HubEntity:    hub,
			KnownTickers: tickers,
			GroupName:    name,
			TickerCount:  cnt,
			Source:       "conglomerate_static",
		}
		byID[groupID] = g
		ordered = append(ordered, g)
	}
	return ordered, byID
}

// GetConglomerateGroup resolves groupID using direct match, -group suffix trimming/addition, and alias map.
func GetConglomerateGroup(ctx context.Context, db *pgxpool.Pool, groupID string) (*ConglomerateGroup, bool) {
	_, byID := LoadConglomerateGroups(ctx, db)
	key := strings.ToLower(strings.TrimSpace(groupID))
	if g, ok := byID[key]; ok {
		return &g, true
	}
	if strings.HasSuffix(key, "-group") {
		if g, ok := byID[key[:len(key)-6]]; ok {
			return &g, true
		}
	}
	if g, ok := byID[key+"-group"]; ok {
		return &g, true
	}
	if target, ok := groupAliasMap[key]; ok {
		if g, exists := byID[target]; exists {
			return &g, true
		}
	}
	return nil, false
}

var staticConglomerateGroupsPreset = []ConglomerateGroup{
	{
		ID:           "barito",
		Name:         "Barito Pacific / Prajogo Pangestu",
		Description:  "Konglomerasi petrokimia, energi baru terbarukan, dan tambang.",
		HubEntity:    "Prajogo Pangestu",
		KnownTickers: []string{"BREN", "BRPT", "CDIA", "CUAN", "STRK", "TPIA"},
		TickerCount:  6,
		Source:       "conglomerate_static",
	},
	{
		ID:           "salim",
		Name:         "Salim Group / Anthoni Salim",
		Description:  "Konglomerasi konsumen primer, agribisnis, ritel, data center, dan infrastruktur.",
		HubEntity:    "Anthoni Salim",
		KnownTickers: []string{"INDF", "ICBP", "SIMP", "LSIP", "DNET", "DCII", "EMTK", "FAST", "MEDC"},
		TickerCount:  9,
		Source:       "conglomerate_static",
	},
	{
		ID:           "astra",
		Name:         "Astra International Group",
		Description:  "Grup otomotif, alat berat, agribisnis, dan jasa keuangan terbesar Indonesia.",
		HubEntity:    "PT Astra International Tbk",
		KnownTickers: []string{"ASII", "AALI", "ASGR", "AUTO", "UNTR", "HEAL"},
		TickerCount:  6,
		Source:       "conglomerate_static",
	},
	{
		ID:           "djarum",
		Name:         "Djarum / Hartono Brothers",
		Description:  "Konglomerasi perbankan swasta terbesar, menara telekomunikasi, dan e-commerce.",
		HubEntity:    "Robert Budi Hartono & Michael Hartono",
		KnownTickers: []string{"BBCA", "TOWR", "SUPR", "RANC", "BELI"},
		TickerCount:  5,
		Source:       "conglomerate_static",
	},
	{
		ID:           "sinarmas",
		Name:         "Sinar Mas / Widjaja Family",
		Description:  "Konglomerasi pulp & paper, agribisnis, energi, jasa keuangan, dan properti.",
		HubEntity:    "Sinar Mas / Widjaja Family",
		KnownTickers: []string{"DSSA", "GEMS", "BSIM", "LIFE", "SMMA", "INKP", "TKIM", "BSDE"},
		TickerCount:  8,
		Source:       "conglomerate_static",
	},
	{
		ID:           "adaro",
		Name:         "Adaro Energy / Boy Thohir",
		Description:  "Konglomerasi energi batubara, green power, kimia, dan investasi pertambangan.",
		HubEntity:    "Garibaldi Thohir",
		KnownTickers: []string{"ADRO", "AADI", "MDKA", "ESSA", "PALM", "TRIM"},
		TickerCount:  6,
		Source:       "conglomerate_static",
	},
	{
		ID:           "lippo",
		Name:         "Lippo Group / Mochtar Riady",
		Description:  "Konglomerasi properti, rumah sakit, ritel modern, dan telekomunikasi.",
		HubEntity:    "PT Inti Anugerah Pratama (Riady Family)",
		KnownTickers: []string{"LPKR", "MLPL", "LPPS", "LPGI", "LPLI", "SILO"},
		TickerCount:  6,
		Source:       "conglomerate_static",
	},
	{
		ID:           "mnc",
		Name:         "MNC Group / Hary Tanoesoedibjo",
		Description:  "Konglomerasi media terintegrasi, hospitality/resort, dan jasa keuangan.",
		HubEntity:    "Hary Tanoesoedibjo",
		KnownTickers: []string{"BMTR", "MNCN", "BCAP", "BABP", "KPIG", "MSKY"},
		TickerCount:  6,
		Source:       "conglomerate_static",
	},
	{
		ID:           "bakrie",
		Name:         "Bakrie Group",
		Description:  "Konglomerasi sumber daya alam, energi, infrastruktur, EV, dan media.",
		HubEntity:    "Bakrie & Brothers",
		KnownTickers: []string{"BNBR", "BRMS", "ENRG", "VKTR", "VIVA", "BTEL", "ALII", "UNSP"},
		TickerCount:  8,
		Source:       "conglomerate_static",
	},
	{
		ID:           "harita",
		Name:         "Harita Group / Lim Family",
		Description:  "Konglomerasi pertambangan nikel terintegrasi, smelter HPAL, dan bauksit.",
		HubEntity:    "PT Harita Jayaraya",
		KnownTickers: []string{"NCKL", "CITA"},
		TickerCount:  2,
		Source:       "conglomerate_static",
	},
	{
		ID:           "medco",
		Name:         "Medco Energi / Panigoro Family",
		Description:  "Konglomerasi migas, tambang tembaga-emas, dan pembangkit listrik.",
		HubEntity:    "Hilmi Panigoro / Medco Daya",
		KnownTickers: []string{"MEDC", "AMMN"},
		TickerCount:  2,
		Source:       "conglomerate_static",
	},
	{
		ID:           "emtek",
		Name:         "Elang Mahkota Teknologi (Emtek)",
		Description:  "Konglomerasi media siaran, teknologi digital, dan jaringan layanan kesehatan.",
		HubEntity:    "Eddy Kusnadi Sariaatmadja",
		KnownTickers: []string{"EMTK", "SCMA", "BUKA", "SAME"},
		TickerCount:  4,
		Source:       "conglomerate_static",
	},
}

// StaticConglomerateGroups returns the 12 preset CONGLOMERATE_GROUPS fallback in deterministic order.
func StaticConglomerateGroups() ([]ConglomerateGroup, map[string]ConglomerateGroup) {
	ordered := make([]ConglomerateGroup, len(staticConglomerateGroupsPreset))
	byID := make(map[string]ConglomerateGroup, len(staticConglomerateGroupsPreset))
	for i, g := range staticConglomerateGroupsPreset {
		tickers := make([]string, len(g.KnownTickers))
		copy(tickers, g.KnownTickers)
		gCopy := g
		gCopy.KnownTickers = tickers
		gCopy.GroupName = g.Name
		ordered[i] = gCopy
		byID[g.ID] = gCopy
	}
	return ordered, byID
}
