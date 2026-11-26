package bandarmology

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
)

type brokerDirEntry struct {
	Name  string
	Group string
	Type  string
}

type documentedAffEntry struct {
	Badge        string
	RelationType string
	Entity       string
	Relation     string
	Confidence   string
}

var brokerDirectory = map[string]brokerDirEntry{
	"II": {"PT Danatama Makmur Sekuritas", "Danatama Group", "Local Institution"},
	"SQ": {"PT BCA Sekuritas", "BCA / Djarum Group", "Local Institution"},
	"CC": {"PT Mandiri Sekuritas", "Bank Mandiri Group (BUMN)", "Local Institution"},
	"NI": {"PT BNI Sekuritas", "BNI Group (BUMN)", "Local Institution"},
	"OD": {"PT BRI Danareksa Sekuritas", "Bank BRI Group (BUMN)", "Local Institution"},
	"DX": {"PT Bahana Sekuritas", "IFG / BUMN Group", "Local Institution"},
	"EP": {"PT MNC Sekuritas", "MNC Group (Hary Tanoesoedibjo)", "Local Institution"},
	"LG": {"PT Trimegah Sekuritas Indonesia Tbk", "Boy Thohir / Garibaldi Group", "Local Institution"},
	"KI": {"PT Ciptadana Sekuritas Asia", "Lippo Group (Keluarga Riady)", "Local Institution"},
	"AZ": {"PT Sucor Sekuritas", "Sucorinvest Group", "Local Institution"},
	"HP": {"PT Henan Putihrai Sekuritas", "Henan / Tamaris Group", "Local Institution"},
	"TP": {"PT OCBC Sekuritas Indonesia", "OCBC Bank Group", "Local Institution"},
	"IF": {"PT Samuel Sekuritas Indonesia", "Samuel Group", "Local Institution"},
	"PD": {"PT Indo Premier Sekuritas", "Indo Premier (IPOT)", "Retail / Institution"},
	"YP": {"PT Mirae Asset Sekuritas Indonesia", "Mirae Asset Group (Korea)", "Retail-Dominant"},
	"XL": {"PT Stockbit Sekuritas Digital", "Stockbit / Bibit", "Retail-Dominant"},
	"XC": {"PT Ajaib Sekuritas Asia", "Ajaib Group", "Retail-Dominant"},
	"AK": {"PT UBS Sekuritas Indonesia", "UBS AG (Swiss)", "Foreign Smart Money"},
	"BK": {"PT J.P. Morgan Sekuritas Indonesia", "J.P. Morgan Chase & Co (US)", "Foreign Smart Money"},
	"MS": {"PT Morgan Stanley Sekuritas Indonesia", "Morgan Stanley (US)", "Foreign Smart Money"},
	"KZ": {"PT CLSA Sekuritas Indonesia", "CITIC / CLSA", "Foreign Smart Money"},
	"CS": {"PT Credit Suisse Sekuritas Indonesia", "UBS / Credit Suisse", "Foreign Smart Money"},
	"GR": {"PT RHB Sekuritas Indonesia", "RHB Bank Group (Malaysia)", "Foreign Smart Money"},
	"ZP": {"PT Maybank Sekuritas Indonesia", "Maybank Group (Malaysia)", "Foreign Smart Money"},
	"YU": {"PT CGS International Sekuritas Indonesia", "CGS International", "Foreign Institution"},
	"RX": {"PT Macquarie Sekuritas Indonesia", "Macquarie Group (Australia)", "Foreign Smart Money"},
	"GA": {"PT Surya Fajar Sekuritas", "Surya Fajar Capital", "Local Speculative / Boutique"},
	"MG": {"PT Semesta Indovest Sekuritas", "Semesta Indovest", "Speculative / Momentum"},
}

type affKey struct{ ticker, broker string }

var documentedAffiliations = map[affKey]documentedAffEntry{
	{"BULL", "II"}: {"AFILIASI PENGENDALI & DIREKSI", "INTERLOCKING_BOARD_AND_OFFICE", "PT Danatama Makmur Sekuritas", "Berbagi Komisaris Utama (Halim Jusuf) dan Direktur Utama (Henry Jusuf, mantan Direktur Danatama) dengan BULL, serta beralamat di kantor yang sama (Danatama Square II). Entitas terafiliasi Danatama memiliki porsi saham ~5.6% di BULL.", "CONFIRMED"},
	{"BBCA", "SQ"}: {"ANAK USAHA EMITEN", "DIRECT_SUBSIDIARY", "PT BCA Sekuritas", "PT BCA Sekuritas adalah entitas anak terkonsolidasi langsung dari PT Bank Central Asia Tbk (Grup Djarum).", "CONFIRMED"},
	{"BMRI", "CC"}: {"ANAK USAHA BUMN", "DIRECT_SUBSIDIARY", "PT Mandiri Sekuritas", "PT Mandiri Sekuritas adalah anak perusahaan langsung dari PT Bank Mandiri (Persero) Tbk.", "CONFIRMED"},
	{"BBNI", "NI"}: {"ANAK USAHA BUMN", "DIRECT_SUBSIDIARY", "PT BNI Sekuritas", "PT BNI Sekuritas adalah anak perusahaan terkonsolidasi dari PT Bank Negara Indonesia (Persero) Tbk.", "CONFIRMED"},
	{"BBRI", "OD"}: {"ANAK USAHA BUMN", "DIRECT_SUBSIDIARY", "PT BRI Danareksa Sekuritas", "PT BRI Danareksa Sekuritas adalah anak usaha terafiliasi dari PT Bank Rakyat Indonesia (Persero) Tbk.", "CONFIRMED"},
	{"ADRO", "LG"}: {"AFILIASI PENGENDALI", "COMMON_CONTROLLER", "PT Trimegah Sekuritas Indonesia Tbk", "Garibaldi 'Boy' Thohir adalah pengendali dan pemilik saham signifikan pada PT Trimegah Sekuritas Indonesia Tbk.", "HIGH"},
	{"ESSA", "LG"}: {"AFILIASI PENGENDALI", "COMMON_CONTROLLER", "PT Trimegah Sekuritas Indonesia Tbk", "Garibaldi 'Boy' Thohir terafiliasi sebagai pemegang saham pengendali pada Trimegah Sekuritas dan emiten ESSA.", "HIGH"},
	{"MBMA", "LG"}: {"AFILIASI PENGENDALI", "COMMON_CONTROLLER", "PT Trimegah Sekuritas Indonesia Tbk", "Afiliasi ekosistem Merdeka / Thohir Group dengan PT Trimegah Sekuritas Indonesia Tbk.", "HIGH"},
	{"BHIT", "EP"}: {"ANAK USAHA KONGLOMERASI", "CONGLOMERATE_MEMBER", "PT MNC Sekuritas", "PT MNC Sekuritas berada di bawah naungan PT MNC Asia Holding Tbk (Grup MNC / Hary Tanoesoedibjo).", "CONFIRMED"},
	{"KPIG", "EP"}: {"AFILIASI KONGLOMERASI", "CONGLOMERATE_MEMBER", "PT MNC Sekuritas", "Entitas grup sekuritas dari ekosistem MNC Group (Hary Tanoesoedibjo).", "CONFIRMED"},
	{"MNCN", "EP"}: {"AFILIASI KONGLOMERASI", "CONGLOMERATE_MEMBER", "PT MNC Sekuritas", "Entitas sekuritas terafiliasi dari ekosistem media & investasi MNC Group.", "CONFIRMED"},
	{"BMTR", "EP"}: {"AFILIASI KONGLOMERASI", "CONGLOMERATE_MEMBER", "PT MNC Sekuritas", "Entitas sekuritas terafiliasi di bawah bendera konglomerasi MNC Group.", "CONFIRMED"},
	{"LPKR", "KI"}: {"AFILIASI KONGLOMERASI", "CONGLOMERATE_MEMBER", "PT Ciptadana Sekuritas Asia", "PT Ciptadana Sekuritas Asia terafiliasi dengan keluarga Riady / Lippo Group.", "HIGH"},
	{"SILO", "KI"}: {"AFILIASI KONGLOMERASI", "CONGLOMERATE_MEMBER", "PT Ciptadana Sekuritas Asia", "Ciptadana Sekuritas terafiliasi di bawah ekosistem konglomerasi Lippo Group.", "HIGH"},
}

const affiliationDisclaimer = "Observational research lead berdasarkan data keterbukaan publik, bukan tuduhan insider trading."

// GetBrokerInfo mirrors affiliation_service.get_broker_info.
func GetBrokerInfo(brokerCode string) map[string]string {
	code := strings.ToUpper(strings.TrimSpace(brokerCode))
	if entry, ok := brokerDirectory[code]; ok {
		return map[string]string{
			"code":  code,
			"name":  entry.Name,
			"group": entry.Group,
			"type":  entry.Type,
		}
	}
	return map[string]string{
		"code":  code,
		"name":  fmt.Sprintf("Broker %s", code),
		"group": "Unknown",
		"type":  "Other Broker",
	}
}

// GetBrokerAffiliation mirrors affiliation_service.get_broker_affiliation.
func GetBrokerAffiliation(ctx context.Context, repo *Repository, ticker, brokerCode string) map[string]any {
	tk := strings.ToUpper(strings.TrimSpace(ticker))
	code := strings.ToUpper(strings.TrimSpace(brokerCode))
	bInfo := GetBrokerInfo(code)

	if doc, ok := documentedAffiliations[affKey{tk, code}]; ok {
		return map[string]any{
			"is_affiliated": true,
			"broker_code":   code,
			"broker_name":   bInfo["name"],
			"broker_group":  bInfo["group"],
			"badge":         doc.Badge,
			"relation_type": doc.RelationType,
			"relation":      doc.Relation,
			"confidence":    doc.Confidence,
			"disclaimer":    affiliationDisclaimer,
		}
	}

	if repo != nil {
		if dyn := repo.CheckDynamicBoardMatch(ctx, tk, code); dyn != nil {
			return map[string]any{
				"is_affiliated": true,
				"broker_code":   code,
				"broker_name":   bInfo["name"],
				"broker_group":  bInfo["group"],
				"badge":         dyn["badge"],
				"relation_type": dyn["relation_type"],
				"relation":      dyn["relation"],
				"confidence":    dyn["confidence"],
				"disclaimer":    affiliationDisclaimer,
			}
		}
	}
	return nil
}

// EnrichBrokerSummaryAffiliations mirrors affiliation_service.enrich_broker_summary_affiliations.
func EnrichBrokerSummaryAffiliations(ctx context.Context, repo *Repository, ticker string, rows []map[string]any) ([]map[string]any, []map[string]any) {
	enriched := make([]map[string]any, 0, len(rows))
	detected := make([]map[string]any, 0)
	for _, r := range rows {
		code := ""
		if v, ok := r["broker"].(string); ok && v != "" {
			code = strings.ToUpper(strings.TrimSpace(v))
		} else if v, ok := r["broker_code"].(string); ok {
			code = strings.ToUpper(strings.TrimSpace(v))
		}
		bInfo := GetBrokerInfo(code)
		aff := GetBrokerAffiliation(ctx, repo, ticker, code)

		nr := make(map[string]any, len(r)+3)
		for k, v := range r {
			nr[k] = v
		}
		nr["broker_name"] = bInfo["name"]
		if aff != nil {
			nr["is_affiliated"] = true
			nr["affiliation"] = aff
			detected = append(detected, map[string]any{
				"broker_code": code,
				"broker_name": aff["broker_name"],
				"side":        r["side"],
				"net":         r["net"],
				"badge":       aff["badge"],
				"relation":    aff["relation"],
				"confidence":  aff["confidence"],
			})
		} else {
			nr["is_affiliated"] = false
			nr["affiliation"] = nil
		}
		enriched = append(enriched, nr)
	}
	return enriched, detected
}

// DetectWindowAffiliations mirrors affiliation_service.detect_window_affiliations.
func DetectWindowAffiliations(ctx context.Context, repo *Repository, ticker string, acts []ActivityRow) []map[string]any {
	if len(acts) == 0 {
		return []map[string]any{}
	}
	type grpAgg struct {
		cumNet     float64
		dates      map[string]struct{}
		latestDate string
		latestNet  float64
	}
	byCode := make(map[string]*grpAgg)
	for _, a := range acts {
		code := strings.ToUpper(strings.TrimSpace(a.BrokerCode))
		if code == "" {
			continue
		}
		dStr := fmtDateISO(a.Date)
		g, ok := byCode[code]
		if !ok {
			g = &grpAgg{dates: make(map[string]struct{})}
			byCode[code] = g
		}
		g.cumNet += a.NetValue
		g.dates[dStr] = struct{}{}
		if dStr >= g.latestDate {
			g.latestDate = dStr
			g.latestNet = a.NetValue
		}
	}

	detected := make([]map[string]any, 0)
	for code, g := range byCode {
		aff := GetBrokerAffiliation(ctx, repo, ticker, code)
		if aff == nil {
			continue
		}
		detected = append(detected, map[string]any{
			"broker_code":    code,
			"broker_name":    aff["broker_name"],
			"broker_group":   aff["broker_group"],
			"badge":          aff["badge"],
			"relation_type":  aff["relation_type"],
			"relation":       aff["relation"],
			"confidence":     aff["confidence"],
			"cumulative_net": g.cumNet,
			"days_active":    len(g.dates),
			"latest_net":     g.latestNet,
			"latest_date":    g.latestDate,
			"disclaimer":     aff["disclaimer"],
		})
	}
	sort.SliceStable(detected, func(i, j int) bool {
		return math.Abs(detected[i]["cumulative_net"].(float64)) > math.Abs(detected[j]["cumulative_net"].(float64))
	})
	return detected
}
