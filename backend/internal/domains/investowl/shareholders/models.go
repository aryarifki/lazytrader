package shareholders

import (
	"math"
	"strconv"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
)

// HoldingItem mirrors each entry in entity holdings arrays.
type HoldingItem struct {
	Ticker        string  `json:"ticker"`
	Pct           float64 `json:"pct"`
	IsControlling bool    `json:"is_controlling"`
}

// EntityHoldingRow mirrors top-cross-holdings & search result items.
type EntityHoldingRow struct {
	Name        string        `json:"name"`
	TickerCount int           `json:"ticker_count"`
	TotalPct    float64       `json:"total_pct"`
	Holdings    []HoldingItem `json:"holdings"`
}

// SearchResponse mirrors GET /api/bandar/shareholders/search.
type SearchResponse struct {
	Query   string             `json:"query"`
	Count   int                `json:"count"`
	Results []EntityHoldingRow `json:"results"`
}

// GroupListItem mirrors each group entry in GET /api/bandar/shareholders/groups.
type GroupListItem struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	HubEntity     string   `json:"hub_entity"`
	LogoURL       *string  `json:"logo_url"`
	CatalogCode   *string  `json:"catalog_code"`
	SampleTickers []string `json:"sample_tickers"`
	TickerCount   int      `json:"ticker_count"`
	Source        string   `json:"source"`
}

// GroupsResponse mirrors GET /api/bandar/shareholders/groups.
type GroupsResponse struct {
	Groups []GroupListItem `json:"groups"`
}

// GraphNode mirrors each node in network-graph.
type GraphNode struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Sector        string `json:"sector"`
	IsControlling bool   `json:"is_controlling"`
	Degree        int    `json:"degree"`
	Val           int    `json:"val"`
}

// GraphLink mirrors each edge in network-graph.
type GraphLink struct {
	Source        string  `json:"source"`
	Target        string  `json:"target"`
	Pct           float64 `json:"pct"`
	IsControlling bool    `json:"is_controlling"`
	Label         string  `json:"label"`
}

// GraphStats mirrors stats object in network-graph.
type GraphStats struct {
	TotalNodes   int    `json:"total_nodes"`
	TotalEmitens int    `json:"total_emitens"`
	TotalLinks   int    `json:"total_links"`
	HubEntity    string `json:"hub_entity"`
}

// AvailGroup mirrors available_groups items in network-graph.
type AvailGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	TickerCount int    `json:"ticker_count"`
}

// NetworkGraphResponse mirrors GET /api/bandar/shareholders/network-graph.
type NetworkGraphResponse struct {
	GroupID         string       `json:"group_id"`
	GroupName       string       `json:"group_name"`
	HubEntity       string       `json:"hub_entity"`
	Nodes           []GraphNode  `json:"nodes"`
	Links           []GraphLink  `json:"links"`
	Stats           GraphStats   `json:"stats"`
	AvailableGroups []AvailGroup `json:"available_groups"`
}

// ShareholderItem mirrors each row in GET /api/bandar/shareholders/{ticker}.
type ShareholderItem struct {
	Name          string  `json:"name"`
	Pct           float64 `json:"pct"`
	IsControlling bool    `json:"is_controlling"`
	Category      string  `json:"category"`
}

// ShareholderSummary mirrors the summary object in GET /api/bandar/shareholders/{ticker}.
type ShareholderSummary struct {
	ControllingPct      float64  `json:"controlling_pct"`
	OtherSubstantialPct float64  `json:"other_substantial_pct"`
	PublicPct           float64  `json:"public_pct"`
	TreasuryPct         float64  `json:"treasury_pct"`
	ManagementPct       float64  `json:"management_pct"`
	FreeFloatEst        float64  `json:"free_float_est"`
	ControllingEntities []string `json:"controlling_entities"`
}

// StockShareholdersResponse mirrors GET /api/bandar/shareholders/{ticker}.
type StockShareholdersResponse struct {
	Ticker       string             `json:"ticker"`
	Count        int                `json:"count"`
	Shareholders []ShareholderItem  `json:"shareholders"`
	Summary      ShareholderSummary `json:"summary"`
}

var staticConglomerateGroupsOrdered = []shared.ConglomerateGroup{
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

var knownNameToTicker = map[string]string{
	"pt barito pacific tbk":              "BRPT",
	"pt chandra asri pacific tbk":        "TPIA",
	"pt chandra asri petrochemical tbk":  "TPIA",
	"pt petrindo jaya kreasi tbk":        "CUAN",
	"pt barito renewables energy tbk":    "BREN",
	"pt indofood sukses makmur tbk":      "INDF",
	"pt indofood cbp sukses makmur tbk":  "ICBP",
	"pt astra international tbk":         "ASII",
	"pt sarana menara nusantara tbk":     "TOWR",
	"pt djarum":                          "BBCA",
	"pt dwimuria investama andalan":      "BBCA",
	"pt dian swastatika sentosa tbk":     "DSSA",
	"pt adaro energy indonesia tbk":      "ADRO",
	"pt adaro andalan indonesia tbk":     "AADI",
	"pt merdeka copper gold tbk":         "MDKA",
	"pt elang mahkota teknologi tbk":     "EMTK",
	"pt surya citra media tbk":           "SCMA",
	"pt global mediacom tbk":             "BMTR",
	"pt bakrie & brothers tbk":           "BNBR",
	"pt bakrie and brothers tbk":         "BNBR",
	"pt harita jayaraya":                 "NCKL",
	"pt medco energi internasional tbk":  "MEDC",
	"pt multipolar tbk":                  "MLPL",
	"pt lippo karawaci tbk":              "LPKR",
}

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
