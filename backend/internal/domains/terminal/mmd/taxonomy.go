package mmd

import "math"

// BandarLevel defines a ratio threshold tier and corresponding badge.
type BandarLevel struct {
	Min   float64
	Badge Badge
}

// BandarLevels matches mmd_engine.BANDAR_LEVELS in descending order.
var BandarLevels = []BandarLevel{
	{Min: 35.0, Badge: Badge{Label: "Massive Acc", Cls: "massive-acc", Color: "#2FD98A"}},
	{Min: 25.0, Badge: Badge{Label: "Big Acc", Cls: "big-acc", Color: "#2FD98A"}},
	{Min: 15.0, Badge: Badge{Label: "Normal Acc", Cls: "normal-acc", Color: "#4ADE80"}},
	{Min: 7.5, Badge: Badge{Label: "Small Acc", Cls: "small-acc", Color: "#86EFAC"}},
	{Min: -7.5, Badge: Badge{Label: "Undetected", Cls: "undetected", Color: "#94A3B8"}},
	{Min: -15.0, Badge: Badge{Label: "Small Dist", Cls: "small-dist", Color: "#FCA5A5"}},
	{Min: -25.0, Badge: Badge{Label: "Normal Dist", Cls: "normal-dist", Color: "#F87171"}},
	{Min: -35.0, Badge: Badge{Label: "Big Dist", Cls: "big-dist", Color: "#EF4444"}},
	{Min: math.Inf(-1), Badge: Badge{Label: "Massive Dist", Cls: "massive-dist", Color: "#DC2626"}},
}

// CoreRetailList defines the deterministic display order in retail_rows.
var CoreRetailList = []string{"YP", "PD", "XC", "NI", "KK", "XL", "EP", "CC", "GR"}

// EffectiveRetailBrokers matches sorted(RETAIL_BROKERS) from mmd_engine.py.
var EffectiveRetailBrokers = []string{"CC", "EP", "GR", "KK", "NI", "PD", "XC", "XL", "YP"}

// RetailBrokers matches mmd_engine.RETAIL_BROKERS.
var RetailBrokers = map[string]bool{
	"CC": true, "EP": true, "GR": true, "KK": true, "NI": true,
	"PD": true, "XC": true, "XL": true, "YP": true,
}

// SemiRetailBrokers matches mmd_engine.SEMI_RITEL_BROKERS.
var SemiRetailBrokers = map[string]bool{
	"AN": true, "AT": true, "CP": true, "HD": true, "HP": true,
	"LG": true, "MG": true, "OD": true, "SF": true, "SQ": true,
}

// SemiInstitusiBrokers matches mmd_engine.SEMI_INSTITUSI_BROKERS.
var SemiInstitusiBrokers = map[string]bool{
	"AZ": true,
}

// InstitusiBrokers matches mmd_engine.INSTITUSI_BROKERS (canonical 73 + AZ).
var InstitusiBrokers = map[string]bool{
	"AD": true, "AF": true, "AG": true, "AH": true, "AI": true, "AK": true, "AO": true, "AP": true, "AR": true, "AZ": true,
	"BB": true, "BF": true, "BK": true, "BQ": true, "BR": true, "BS": true, "CD": true, "DD": true, "DH": true, "DP": true,
	"DR": true, "DU": true, "DX": true, "EL": true, "ES": true, "FO": true, "FS": true, "FZ": true, "GA": true, "GI": true,
	"GW": true, "IC": true, "ID": true, "IF": true, "IH": true, "II": true, "IN": true, "IP": true, "IT": true, "IU": true,
	"JB": true, "KI": true, "KZ": true, "LS": true, "MI": true, "MU": true, "OK": true, "PC": true, "PF": true, "PG": true,
	"PI": true, "PO": true, "PP": true, "PS": true, "QA": true, "RB": true, "RF": true, "RG": true, "RO": true, "RS": true,
	"RX": true, "SA": true, "SH": true, "SS": true, "TF": true, "TP": true, "TS": true, "XA": true, "YB": true, "YJ": true,
	"YO": true, "YU": true, "ZP": true, "ZR": true,
}

// AI Score constants matching mmd_engine.py.
const (
	AIScoreMax        = 98
	AIScoreMin        = 6
	AIScoreCorrMult   = 3.5
	AIScoreTotalScale = 1.5
)

var (
	resumeWeekWTD1WWeight = map[string]int{
		"massive-acc": 10, "big-acc": 6, "normal-acc": 3,
		"small-acc": 0, "undetected": 0, "small-dist": 0,
		"normal-dist": -3, "big-dist": -6, "massive-dist": -10,
	}
	resumeWeek2WWeight = map[string]int{
		"massive-acc": 8, "big-acc": 5, "normal-acc": 2,
		"small-acc": 0, "undetected": 0, "small-dist": 0,
		"normal-dist": -2, "big-dist": -5, "massive-dist": -8,
	}
	resumeDailyWeight = map[string]int{
		"massive-acc": 5, "big-acc": 2, "normal-acc": 1,
		"small-acc": 0, "undetected": 0, "small-dist": 0,
		"normal-dist": -1, "big-dist": -2, "massive-dist": -5,
	}
	resumeDaily30Weight = map[string]float64{
		"massive-acc": 1.5, "big-acc": 1.0,
		"small-acc": 0.0, "undetected": 0.0, "small-dist": 0.0, "normal-acc": 0.0, "normal-dist": 0.0,
		"big-dist": -1.0, "massive-dist": -1.5,
	}
	resumePeriodWeight = map[string]int{
		"massive-acc": 20, "big-acc": 12, "normal-acc": 7,
		"small-acc": 0, "undetected": 0, "small-dist": 0,
		"normal-dist": -7, "big-dist": -12, "massive-dist": -20,
	}
	resumeMonthCurrWeight = map[string]int{
		"massive-acc": 10, "big-acc": 6, "normal-acc": 3,
		"small-acc": 0, "undetected": 0, "small-dist": 0,
		"normal-dist": -3, "big-dist": -6, "massive-dist": -10,
	}
	resumeMonthPrevWeight = map[string]int{
		"massive-acc": 8, "big-acc": 5, "normal-acc": 2,
		"small-acc": 0, "undetected": 0, "small-dist": 0,
		"normal-dist": -2, "big-dist": -5, "massive-dist": -8,
	}
	corrIdealRange = map[string][2]float64{
		"massive-acc":  {2.5, 7.0},
		"big-acc":      {1.0, 4.5},
		"normal-acc":   {0.0, 2.5},
		"normal-dist":  {-2.5, 0.0},
		"big-dist":     {-4.5, -1.0},
		"massive-dist": {-7.0, -2.5},
	}
	corrNeutralTiers = map[string]bool{
		"small-acc": true, "undetected": true, "small-dist": true,
	}
)

type sidewaysPeriod struct {
	days int
	band float64
	acc  map[string]int
	dist map[string]int
}

var sidewaysPeriods = []sidewaysPeriod{
	{days: 5, band: 1.0, acc: map[string]int{"massive-acc": 15, "big-acc": 10}, dist: map[string]int{"massive-dist": -15, "big-dist": -10}},
	{days: 7, band: 1.25, acc: map[string]int{"massive-acc": 17, "big-acc": 12}, dist: map[string]int{"massive-dist": -17, "big-dist": -12}},
	{days: 10, band: 1.4, acc: map[string]int{"massive-acc": 20, "big-acc": 15}, dist: map[string]int{"massive-dist": -20, "big-dist": -15}},
}

// ClassifyBandar maps an optional ratio percentage to a Bandar badge.
func ClassifyBandar(ratioPct *float64) Badge {
	if ratioPct == nil || math.IsNaN(*ratioPct) {
		return Badge{Label: "—", Cls: "undetected", Color: "#94A3B8"}
	}
	return ClassifyBandarVal(*ratioPct)
}

// ClassifyBandarVal maps a float64 ratio percentage to a Bandar badge.
func ClassifyBandarVal(ratioPct float64) Badge {
	if math.IsNaN(ratioPct) {
		return Badge{Label: "—", Cls: "undetected", Color: "#94A3B8"}
	}
	for _, lvl := range BandarLevels {
		if ratioPct >= lvl.Min {
			return lvl.Badge
		}
	}
	return BandarLevels[len(BandarLevels)-1].Badge
}

// ClassifyRetailAction maps retail net ratio to a retail action badge.
func ClassifyRetailAction(retailRatio float64) Badge {
	if retailRatio >= 15.0 {
		return Badge{Label: "Aggressive Buy", Cls: "aggressive-buy", Color: "#2FD98A"}
	}
	if retailRatio >= 7.5 {
		return Badge{Label: "Buy", Cls: "normal-buy", Color: "#4ADE80"}
	}
	if retailRatio > -7.5 {
		return Badge{Label: "Do Nothing", Cls: "undetected", Color: "#94A3B8"}
	}
	if retailRatio > -15.0 {
		return Badge{Label: "Sell", Cls: "normal-sell", Color: "#F87171"}
	}
	return Badge{Label: "Panic Sell", Cls: "panic-sell", Color: "#EF4444"}
}

// ClassifyBrokerCategory returns the category label used by broker-drilldown.
func ClassifyBrokerCategory(broker string) string {
	if RetailBrokers[broker] {
		return "Retail"
	}
	if SemiRetailBrokers[broker] {
		return "Semi-Ritel"
	}
	if SemiInstitusiBrokers[broker] {
		return "Semi-Institusi"
	}
	return "Institusi"
}

// CategoryLabelForBroker mirrors mmd_engine.category_label_for_broker verbatim.
func CategoryLabelForBroker(broker string, rankIdx int, side, resumeLabel string, activeRetailBandar *string) string {
	if activeRetailBandar != nil && broker == *activeRetailBandar {
		return "Pengendali"
	}
	if RetailBrokers[broker] {
		return "Retail"
	}
	if SemiInstitusiBrokers[broker] {
		if rankIdx == 0 {
			return "Pengendali"
		}
		return "Semi-Institusi"
	}
	if InstitusiBrokers[broker] {
		if rankIdx == 0 {
			return "Pengendali"
		}
		if rankIdx <= 1 && (resumeLabel == "Big Acc" || resumeLabel == "Massive Acc" || resumeLabel == "Big Dist" || resumeLabel == "Massive Dist") {
			return "Pengendali"
		}
		return "Institusi"
	}
	if SemiRetailBrokers[broker] {
		return "Semi-Ritel"
	}
	if rankIdx == 0 {
		return "Pengendali"
	}
	return "Others"
}

// VolumeRatioWeight matches mmd_engine.volume_ratio_weight (FMA-guarded float64 steps).
func VolumeRatioWeight(vr float64) float64 {
	if vr <= 3.5 {
		return 1.0
	}
	if vr <= 10.0 {
		return 0.7 + float64(float64(10.0-vr)/6.5)*0.3
	}
	if vr <= 15.0 {
		return 0.5 + float64(float64(15.0-vr)/5.0)*0.2
	}
	if vr <= 30.0 {
		return 0.25 + float64(float64(30.0-vr)/15.0)*0.25
	}
	return 0.25
}

// AIScoreCorrDayPoint mirrors mmd_engine.ai_score_corr_day_point verbatim.
func AIScoreCorrDayPoint(cls string, pct float64) float64 {
	if math.IsNaN(pct) || cls == "" {
		return 0.0
	}
	if corrNeutralTiers[cls] {
		if -1.5 <= pct && pct <= 1.5 {
			return 5.0
		}
		return 1.0
	}
	rng, ok := corrIdealRange[cls]
	if !ok {
		return 0.0
	}
	lo, hi := rng[0], rng[1]
	if lo <= pct && pct <= hi {
		return 5.0
	}

	nearBound := hi
	farBound := lo
	if math.Abs(lo) < math.Abs(hi) {
		nearBound = lo
		farBound = hi
	}
	width := math.Abs(farBound - nearBound)
	refNear := math.Abs(nearBound)
	sameDirection := pct <= 0
	if farBound > 0 {
		sameDirection = pct >= 0
	}

	if sameDirection {
		if math.Abs(pct) < refNear {
			frac := 0.0
			if refNear > 0 {
				frac = math.Abs(pct) / refNear
			}
			return 1.0 + float64(frac*4.0)
		}
		overshoot := math.Abs(pct) - math.Abs(farBound)
		frac := 0.0
		if width > 0 {
			frac = math.Max(0.0, 1.0-float64(overshoot/width))
		}
		return 1.0 + float64(frac*4.0)
	}

	if refNear == 0 {
		return -5.0
	}
	frac := math.Min(1.0, math.Abs(pct)/refNear)
	return -1.0 - float64(frac*4.0)
}
