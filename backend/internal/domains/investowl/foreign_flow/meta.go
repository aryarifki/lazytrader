package foreign_flow

// SectorMetaEntry defines display & DB mapping metadata for an IDX sector index.
type SectorMetaEntry struct {
	Code     string
	Name     string
	DBSector string
	Icon     string
	Color    string
}

// SectorOrder preserves deterministic iteration order of SECTOR_META in routers/foreign_flow.py.
var SectorOrder = []string{
	"IDXENERGY", "IDXBASIC", "IDXINDUST", "IDXNONCYC", "IDXCYCLIC",
	"IDXHEALTH", "IDXFINANCE", "IDXPROPERT", "IDXTECHNO", "IDXINFRA", "IDXTRANS",
}

// SectorMeta mirrors SECTOR_META from routers/foreign_flow.py.
var SectorMeta = map[string]SectorMetaEntry{
	"IDXENERGY":  {Code: "IDXENERGY", Name: "Energi (Energy)", DBSector: "Energi", Icon: "ph:lightning-duotone", Color: "amber"},
	"IDXBASIC":   {Code: "IDXBASIC", Name: "Barang Baku (Basic Materials)", DBSector: "Barang Baku", Icon: "ph:cube-duotone", Color: "emerald"},
	"IDXINDUST":  {Code: "IDXINDUST", Name: "Perindustrian (Industrials)", DBSector: "Perindustrian", Icon: "ph:factory-duotone", Color: "indigo"},
	"IDXNONCYC":  {Code: "IDXNONCYC", Name: "Konsumer Primer (Consumer Non-Cyclicals)", DBSector: "Barang Konsumen Primer", Icon: "ph:shopping-cart-duotone", Color: "teal"},
	"IDXCYCLIC":  {Code: "IDXCYCLIC", Name: "Konsumer Non-Primer (Consumer Cyclicals)", DBSector: "Barang Konsumen Non-Primer", Icon: "ph:storefront-duotone", Color: "sky"},
	"IDXHEALTH":  {Code: "IDXHEALTH", Name: "Kesehatan (Healthcare)", DBSector: "Kesehatan", Icon: "ph:first-aid-duotone", Color: "rose"},
	"IDXFINANCE": {Code: "IDXFINANCE", Name: "Keuangan (Financials)", DBSector: "Keuangan", Icon: "ph:bank-duotone", Color: "blue"},
	"IDXPROPERT": {Code: "IDXPROPERT", Name: "Properti & Real Estat (Real Estate)", DBSector: "Properti & Real Estat", Icon: "ph:buildings-duotone", Color: "purple"},
	"IDXTECHNO":  {Code: "IDXTECHNO", Name: "Teknologi (Technology)", DBSector: "Teknologi", Icon: "ph:cpu-duotone", Color: "violet"},
	"IDXINFRA":   {Code: "IDXINFRA", Name: "Infrastruktur (Infrastructure)", DBSector: "Infrastruktur", Icon: "ph:broadcast-duotone", Color: "orange"},
	"IDXTRANS":   {Code: "IDXTRANS", Name: "Transportasi & Logistik (Logistics)", DBSector: "Transportasi & Logistik", Icon: "ph:airplane-tilt-duotone", Color: "cyan"},
}

// DefaultFallbackStocks mirrors DEFAULT_FALLBACK_STOCKS in routers/foreign_flow.py.
var DefaultFallbackStocks = map[string][]string{
	"IDXENERGY":  {"ADRO", "BREN", "PTBA", "MEDC", "AKRA"},
	"IDXBASIC":   {"BRPT", "TPIA", "INKP", "MDKA", "ANTM"},
	"IDXINDUST":  {"UNTR", "ASII", "ASGR", "MARK", "HEXA"},
	"IDXNONCYC":  {"ICBP", "INDF", "UNVR", "CPIN", "AMRT"},
	"IDXCYCLIC":  {"ASII", "ACES", "MAPI", "AUTO", "ERAA"},
	"IDXHEALTH":  {"KLBF", "MIKA", "HEAL", "SILO", "SIDO"},
	"IDXFINANCE": {"BBCA", "BBRI", "BMRI", "BBNI", "BRIS"},
	"IDXPROPERT": {"BSDE", "CTRA", "PWON", "SMRA", "PANI"},
	"IDXTECHNO":  {"GOTO", "EMTK", "DCII", "BUKA", "WIFI"},
	"IDXINFRA":   {"TLKM", "JSMR", "TOWR", "TBIG", "ISAT"},
	"IDXTRANS":   {"BIRD", "SMDR", "ASSA", "TMAS", "GIAA"},
}

const (
	thrExodus1D           = -250_000_000_000.0
	thrExodus5D           = -500_000_000_000.0
	thrHiddenDist5D       = -150_000_000_000.0
	thrInflowSurge1D      = 80_000_000_000.0
	thrInflowSurgeCombo1D = 40_000_000_000.0
	thrInflowSurgeCombo5D = 60_000_000_000.0
	thrBottoming1D        = 15_000_000_000.0
	thrRetHiddenDist5D    = -0.8
	thrRetHiddenDistPeak5 = -0.5
	thrRetBottoming20D    = -3.0
	thrRetBottoming5D     = -2.0
)

// EvaluateSectorForeignState mirrors routers/foreign_flow.py _evaluate_sector_foreign_state.
func EvaluateSectorForeignState(net1D, net5D, ret5D, ret20D float64, macroState string) (string, string) {
	if net1D <= thrExodus1D || (net5D <= thrExodus5D && net1D < 0) {
		if ret5D >= thrRetHiddenDist5D {
			return "HIDDEN DISTRIBUTION", "Distribusi Asing Tersembunyi (Harga Ditahan)"
		}
		return "SECTOR EXODUS", "Eksodus Modal Asing Masif (Tekanan Jual Kuat)"
	}

	if net5D <= thrHiddenDist5D && net1D <= 0 && ret5D >= thrRetHiddenDistPeak5 {
		return "HIDDEN DISTRIBUTION", "Distribusi Asing Tersembunyi di Puncak"
	}

	if (net1D >= thrInflowSurgeCombo1D && net5D >= thrInflowSurgeCombo5D) || (net1D >= thrInflowSurge1D) {
		return "INFLOW SURGE", "Serbuan Akumulasi Modal Asing (Inflow Kuat)"
	}

	if (ret20D <= thrRetBottoming20D || ret5D <= thrRetBottoming5D) && net1D >= thrBottoming1D {
		if (macroState == "CRASH" || macroState == "RECOVERY_TESTING") && net1D < thrInflowSurgeCombo1D {
			return "NORMAL", "Rebound Spekulatif Sektoral (Pasar Makro Masih Tertekan)"
		}
		return "BOTTOMING", "Kapitulasi Selesai: Asing Mulai Serok Bawah"
	}

	if net1D > 0 {
		return "NORMAL", "Aliran Asing Stabil (Net Buy Moderat)"
	}
	return "NORMAL", "Aliran Asing Normal (Netral)"
}
