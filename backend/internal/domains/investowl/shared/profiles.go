package shared

import (
	"math"
	"strconv"
	"strings"
)

// NormalizeTicker uppercases, trims, and strips the .JK suffix.
func NormalizeTicker(ticker string) string {
	s := strings.ToUpper(ticker)
	s = strings.ReplaceAll(s, ".JK", "")
	return strings.TrimSpace(s)
}

// Round performs exact decimal half-to-even rounding matching CPython round(v, decimals).
func Round(v float64, decimals int) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	if decimals < 0 {
		pow := math.Pow10(-decimals)
		res := math.RoundToEven(v/pow) * pow
		if res == 0 {
			return 0
		}
		return res
	}
	res, err := strconv.ParseFloat(strconv.FormatFloat(v, 'f', decimals, 64), 64)
	if err != nil || res == 0 {
		return 0
	}
	return res
}

// Finite converts a nullable float64 to a finite float64, falling back if nil/NaN/Inf.
func Finite(v *float64, fallback float64) float64 {
	if v == nil || math.IsNaN(*v) || math.IsInf(*v, 0) {
		return fallback
	}
	return *v
}

// FinitePtr returns nil if v is nil, NaN, or Inf; otherwise returns a clean float64 pointer.
func FinitePtr(v *float64) *float64 {
	if v == nil || math.IsNaN(*v) || math.IsInf(*v, 0) {
		return nil
	}
	cp := *v
	if cp == 0 {
		cp = 0
	}
	return &cp
}

var brokerProfiles = map[string]string{
	"AK": "smart_foreign", "BK": "smart_foreign", "MS": "smart_foreign",
	"GR": "smart_foreign", "LG": "smart_foreign", "KZ": "smart_foreign",
	"CS": "smart_foreign", "DX": "smart_foreign",
	"SQ": "bandar_gorengan", "MG": "bandar_gorengan", "EP": "bandar_gorengan",
	"DR": "bandar_gorengan", "BZ": "bandar_gorengan",
	"XA": "retail", "AZ": "retail", "KI": "retail", "YO": "retail", "ZP": "retail",
	"CC": "local_institutional", "NI": "local_institutional", "OD": "local_institutional",
	"TP": "local_institutional", "IF": "local_institutional",
	"YU": "market_maker", "RX": "market_maker", "PD": "market_maker",
}

// ProfileMetadata mirrors idx_bandarmology/analysis.py _PROFILE_META.
type ProfileMetadata struct {
	Label       string
	Description string
}

// ProfileMeta maps behavioral profile keys to their labels and descriptions.
var ProfileMeta = map[string]ProfileMetadata{
	"smart_foreign":       {Label: "Foreign Smart Money", Description: "Directional foreign institutions with higher conviction"},
	"local_institutional": {Label: "Local Institutions", Description: "Local funds and institution-like accounts"},
	"market_maker":        {Label: "Market Makers", Description: "Often active on both sides; net position matters"},
	"bandar_gorengan":     {Label: "Speculative Operators", Description: "Potential pump-and-dump style participants"},
	"retail":              {Label: "Retail-Dominant", Description: "Retail-heavy platforms, often late or contrarian"},
	"lainnya":             {Label: "Other Brokers", Description: "Brokers outside the defined behavioral profiles"},
}

// BrokerProfileOf classifies a broker code into its behavioral profile.
func BrokerProfileOf(code string) string {
	if p, ok := brokerProfiles[strings.ToUpper(strings.TrimSpace(code))]; ok {
		return p
	}
	return "lainnya"
}

// IsSmartProfile reports whether a profile is in SMART_PROFILES.
func IsSmartProfile(profile string) bool {
	return profile == "smart_foreign" || profile == "local_institutional"
}

// BrokerNameMap maps major broker codes to institutional display names (routers/forensic.py).
var BrokerNameMap = map[string]string{
	"ZP": "Credit Suisse",
	"AK": "Morgan Stanley",
	"BK": "Deutsche Bank",
	"KZ": "CLSA Sekuritas",
	"RX": "UBS Sekuritas",
	"TP": "JPMorgan Sekuritas",
	"YU": "CGS-CIMB Sekuritas",
	"CS": "Credit Suisse",
	"CG": "Citigroup Sekuritas",
	"ML": "Merrill Lynch / BofA",
	"MS": "Morgan Stanley",
	"DB": "Deutsche Bank",
	"CC": "Mandiri Sekuritas",
	"NI": "BNI Sekuritas",
	"OD": "BRI Danareksa",
	"PD": "Indo Premier Sekuritas",
	"YP": "Mirae Asset Sekuritas",
	"XC": "Ajaib Sekuritas",
	"SQ": "BCA Sekuritas",
	"MG": "Semesta Indovest",
	"CP": "KB Valbury Sekuritas",
	"EP": "MNC Sekuritas",
	"LG": "Trimegah Sekuritas",
	"GR": "Panin Sekuritas",
	"AZ": "Sucor Sekuritas",
	"DR": "RHB Sekuritas",
	"KI": "Ciptadana Sekuritas",
	"AI": "UOB Kay Hian",
	"AG": "Kiwoom Sekuritas",
	"XA": "Woori Korindo",
}

// AccSignals and DistSignals define canonical signal buckets.
var (
	AccSignals  = map[string]bool{"STRONG_ACCUMULATION": true, "ACCUMULATION": true, "NET_BUY": true, "AKUMULASI_KUAT": true, "AKUMULASI": true}
	DistSignals = map[string]bool{"STRONG_DISTRIBUTION": true, "DISTRIBUTION": true, "NET_SELL": true, "DISTRIBUSI_KUAT": true, "DISTRIBUSI": true}
)

var signalFormatMap = map[string]string{
	"AKUMULASI_KUAT":      "Strong Accumulation",
	"AKUMULASI":           "Accumulation",
	"DISTRIBUSI_KUAT":     "Strong Distribution",
	"DISTRIBUSI":          "Distribution",
	"NETRAL":              "Neutral",
	"STRONG_ACCUMULATION": "Strong Accumulation",
	"ACCUMULATION":        "Accumulation",
	"NET_BUY":             "Net Buy",
	"STRONG_DISTRIBUTION": "Strong Distribution",
	"DISTRIBUTION":        "Distribution",
	"NET_SELL":            "Net Sell",
	"NEUTRAL":             "Neutral",
}

// FmtSignal normalizes internal signal identifiers to human-readable labels.
func FmtSignal(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "-"
	}
	if mapped, ok := signalFormatMap[s]; ok {
		return mapped
	}
	return titleWords(strings.ReplaceAll(s, "_", " "))
}

// TitleWords formats underscore-split signal names like Python str.title().
func TitleWords(s string) string {
	return titleWords(s)
}

func titleWords(s string) string {
	parts := strings.Fields(strings.ToLower(s))
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

// ParticipantLabel translates Indonesian participant types to English labels.
func ParticipantLabel(raw string) string {
	switch strings.TrimSpace(raw) {
	case "Asing":
		return "FOREIGN"
	case "Lokal":
		return "LOCAL"
	case "Pemerintah":
		return "GOV"
	case "":
		return "-"
	default:
		return raw
	}
}
