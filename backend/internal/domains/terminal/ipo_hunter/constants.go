package ipohunter

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// Exact broker classifications and thresholds matching Creative Trader IPO Hunter.
var (
	RetailBrokers = []string{"XL", "YP", "PD", "XC", "NI", "KK", "EP", "DH"}
	RetailBrokerSet = map[string]struct{}{
		"XL": {}, "YP": {}, "PD": {}, "XC": {}, "NI": {}, "KK": {}, "EP": {}, "DH": {},
	}
	SemiBrokers = map[string]struct{}{
		"AT": {}, "AZ": {}, "CP": {}, "DR": {}, "GR": {}, "OD": {}, "SQ": {},
	}
	InstitusiBrokers = map[string]struct{}{
		"CC": {}, "AK": {}, "ZP": {}, "YU": {}, "BK": {}, "MG": {}, "LG": {},
		"BB": {}, "BQ": {}, "HD": {}, "KZ": {}, "TP": {}, "YB": {}, "XA": {},
	}
	NonRitelBrokers = map[string]struct{}{
		"IF": {}, "RX": {}, "DX": {}, "KI": {}, "AG": {}, "AI": {},
		"YJ": {}, "HP": {}, "SS": {}, "AO": {}, "DP": {}, "IN": {},
	}
	RepoBrokers = map[string]struct{}{
		"LS": {}, "BQ": {}, "TP": {}, "AG": {}, "RF": {}, "IN": {}, "HD": {}, "SH": {},
		"FS": {}, "II": {}, "FZ": {}, "AN": {}, "PO": {}, "SA": {}, "PS": {}, "AP": {},
	}
	MNCGroup = map[string]struct{}{
		"BMTR": {}, "MNCN": {}, "BHIT": {}, "BCAP": {}, "IPTV": {}, "KPIG": {}, "BABP": {}, "MSIN": {},
	}
	SinarmasGroup = map[string]struct{}{
		"DSSA": {}, "INKP": {}, "TKIM": {}, "SMMA": {}, "BSIM": {}, "SMAR": {}, "MORA": {},
	}
	RetailBandarCoef = map[string]float64{
		"XL": 3.5, "YP": 1.0, "PD": 1.0, "XC": 1.0,
		"NI": 0.5, "KK": 0.5, "EP": 0.5, "DH": 0.5,
	}
	RetailBandarTransThreshold = map[string]float64{
		"XL": 35.0, "YP": 20.0, "PD": 20.0, "XC": 20.0,
		"NI": 15.5, "KK": 15.5, "EP": 15.5, "DH": 15.5,
	}
	BandarLevels = []struct {
		MinVal float64
		Label  string
		Cls    string
	}{
		{35.0, "Massive Acc", "massive-acc"},
		{25.0, "Big Acc", "big-acc"},
		{15.0, "Normal Acc", "normal-acc"},
		{7.5, "Small Acc", "small-acc"},
		{-7.5, "Undetected", "undetected"},
		{-15.0, "Small Dist", "small-dist"},
		{-25.0, "Normal Dist", "normal-dist"},
		{-35.0, "Big Dist", "big-dist"},
		{math.Inf(-1), "Massive Dist", "massive-dist"},
	}
	PieKeys = []string{"retail", "bandar", "semi", "institusi", "nonRitel", "repo", "uw"}
	PieMeta = map[string][2]string{
		"retail":    {"Retail", "#0579CC"},
		"bandar":    {"Ritel dipakai Bandar", "#8B95A7"},
		"semi":      {"Semi Ritel", "#255186"},
		"institusi": {"Institusi", "#CC3939"},
		"nonRitel":  {"Non Ritel", "#A8504A"},
		"repo":      {"Repo", "#6B2B29"},
		"uw":        {"Underwriter", "#E24B3D"},
		"ipo":       {"IPO", "#CBD5E1"},
		"others":    {"Lain-lain", "#C4A574"},
	}
	ScrPresetDefaults = map[string]struct {
		Col string
		Dir int
	}{
		"topPerf":  {Col: "chg", Dir: -1},
		"topLoser": {Col: "chg", Dir: 1},
		"volSurge": {Col: "volRatio", Dir: -1},
	}
)

// JSRound matches JavaScript Math.round(x) (half towards +infinity).
func JSRound(x float64) int64 {
	return int64(math.Floor(x + 0.5))
}

// ParseISODate parses YYYY-MM-DD (or prefix) into UTC midnight time.
func ParseISODate(s string) (time.Time, bool) {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 10 {
		return time.Time{}, false
	}
	t, err := time.Parse("2006-01-02", trimmed[:10])
	if err != nil {
		return time.Time{}, false
	}
	return t.UTC(), true
}

// FormatDateISO returns YYYY-MM-DD in UTC.
func FormatDateISO(t time.Time) string {
	return t.UTC().Format("2006-01-02")
}

// DateOrdinal converts a UTC date to Python date.toordinal() integer.
func DateOrdinal(t time.Time) int {
	u := time.Date(t.UTC().Year(), t.UTC().Month(), t.UTC().Day(), 0, 0, 0, 0, time.UTC)
	return int(u.Unix()/86400) + 719163
}

// TradingDays counts inclusive Mon-Fri non-holiday dates between a and b.
func TradingDays(a, b time.Time, holidays map[int]struct{}) int {
	startOrd := DateOrdinal(a)
	endOrd := DateOrdinal(b)
	if endOrd < startOrd {
		return 0
	}
	cur := time.Date(a.UTC().Year(), a.UTC().Month(), a.UTC().Day(), 0, 0, 0, 0, time.UTC)
	n := 0
	for ord := startOrd; ord <= endOrd; ord++ {
		wd := cur.Weekday()
		if wd != time.Saturday && wd != time.Sunday {
			if _, isHol := holidays[ord]; !isHol {
				n++
			}
		}
		cur = cur.AddDate(0, 0, 1)
	}
	return n
}

// WeekdaysBetween counts trading days clamped to at least 1.
func WeekdaysBetween(a, b time.Time, holidays map[int]struct{}) int {
	n := TradingDays(a, b, holidays)
	if n < 1 {
		return 1
	}
	return n
}

// UWSetOf returns a set of underwriter broker codes for a master IPO record.
func UWSetOf(master *MasterRecord) map[string]struct{} {
	out := make(map[string]struct{})
	if master == nil {
		return out
	}
	for _, u := range master.Underwriters {
		out[u.Broker] = struct{}{}
	}
	return out
}

// EffectiveRetailBrokers preserves RetailBrokers order while excluding MNC EP, Sinarmas DH, and underwriters.
func EffectiveRetailBrokers(ticker string, uwSet map[string]struct{}) []string {
	out := make([]string, 0, len(RetailBrokers))
	_, isMNC := MNCGroup[ticker]
	_, isSinarmas := SinarmasGroup[ticker]
	for _, b := range RetailBrokers {
		if b == "EP" && isMNC {
			continue
		}
		if b == "DH" && isSinarmas {
			continue
		}
		if _, isUW := uwSet[b]; isUW {
			continue
		}
		out = append(out, b)
	}
	return out
}

// CategoryOf maps a broker to its base category with UW > Group > Retail/Semi/Institusi/NonRitel > Zombie priority.
func CategoryOf(broker, ticker string, uwSet map[string]struct{}) string {
	if _, ok := uwSet[broker]; ok {
		return "uw"
	}
	if broker == "EP" {
		if _, ok := MNCGroup[ticker]; ok {
			return "institusi"
		}
	}
	if broker == "DH" {
		if _, ok := SinarmasGroup[ticker]; ok {
			return "institusi"
		}
	}
	if _, ok := RetailBrokerSet[broker]; ok {
		return "retail"
	}
	if _, ok := SemiBrokers[broker]; ok {
		return "semi"
	}
	if _, ok := InstitusiBrokers[broker]; ok {
		return "institusi"
	}
	if _, ok := NonRitelBrokers[broker]; ok {
		return "nonRitel"
	}
	return "zombie"
}

// ClassifyRetailAction classifies net retail ratio percentage into label and CSS class.
func ClassifyRetailAction(ratio *float64) ActionLabel {
	if ratio == nil || math.IsNaN(*ratio) {
		return ActionLabel{Label: "—", Cls: "undetected"}
	}
	r := *ratio
	if r > 0 {
		if r < 3 {
			return ActionLabel{Label: "Do Nothing", Cls: "undetected"}
		}
		if r < 7 {
			return ActionLabel{Label: "Small Buy", Cls: "retail-buy"}
		}
		if r < 15 {
			return ActionLabel{Label: "Normal Buy", Cls: "retail-buy"}
		}
		if r < 30 {
			return ActionLabel{Label: "Big Buy", Cls: "retail-buy"}
		}
		return ActionLabel{Label: "Massive Buy", Cls: "retail-buy"}
	}
	if r <= -30 {
		return ActionLabel{Label: "Massive Sell", Cls: "retail-sell"}
	}
	if r <= -15 {
		return ActionLabel{Label: "Big Sell", Cls: "retail-sell"}
	}
	if r <= -7 {
		return ActionLabel{Label: "Normal Sell", Cls: "retail-sell"}
	}
	if r <= -3 {
		return ActionLabel{Label: "Small Sell", Cls: "retail-sell"}
	}
	return ActionLabel{Label: "Do Nothing", Cls: "undetected"}
}

// ClassifyBandar classifies adjusted bandar accumulation ratio into label and CSS class.
func ClassifyBandar(ratioPct *float64) ActionLabel {
	if ratioPct == nil || math.IsNaN(*ratioPct) {
		return ActionLabel{Label: "—", Cls: "undetected"}
	}
	r := *ratioPct
	for _, lvl := range BandarLevels {
		if r >= lvl.MinVal {
			return ActionLabel{Label: lvl.Label, Cls: lvl.Cls}
		}
	}
	last := BandarLevels[len(BandarLevels)-1]
	return ActionLabel{Label: last.Label, Cls: last.Cls}
}

// VolumeRatioWeight applies the piecewise CT Volume Ratio damping weight.
func VolumeRatioWeight(vr float64) float64 {
	if vr <= 3.5 {
		return 1.0
	}
	if vr <= 10.0 {
		return 0.7 + (10.0-vr)/6.5*0.3
	}
	if vr <= 15.0 {
		return 0.5 + (15.0-vr)/5.0*0.2
	}
	if vr <= 30.0 {
		return 0.25 + (30.0-vr)/15.0*0.25
	}
	return 0.25
}

// FormatIntComma formats a signed integer with comma thousands separators matching Python f"{n:,}".
func FormatIntComma(n int64) string {
	neg := n < 0
	var digits string
	if neg {
		digits = strconv.FormatInt(-n, 10)
	} else {
		digits = strconv.FormatInt(n, 10)
	}
	l := len(digits)
	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	for i := 0; i < l; i++ {
		if i > 0 && (l-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(digits[i])
	}
	return b.String()
}
