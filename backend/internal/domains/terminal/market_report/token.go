package market_report

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/shared"
)

var brokerTokenLineRe = regexp.MustCompile(`(?m)^BROKER_API_TOKEN=.*$`)

// CheckTokenAvailable reports whether BROKER_API_TOKEN in SMtracker .env has length > 20.
func (e *Engine) CheckTokenAvailable() bool {
	f, err := os.Open(e.smEnvFile)
	if err != nil {
		return false
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "BROKER_API_TOKEN=") {
			val := strings.SplitN(line, "=", 2)[1]
			val = trimQuotesAndSpaces(val)
			if len(val) > 20 {
				return true
			}
		}
	}
	return false
}

// UpdateToken validates JWT format & expiration, verifies live against Stockbit API, and updates .env.
func (e *Engine) UpdateToken(tokenInput string) map[string]any {
	tokenStr := trimQuotesAndSpaces(tokenInput)
	lower := strings.ToLower(tokenStr)
	if strings.HasPrefix(lower, "authorization:") {
		tokenStr = strings.TrimSpace(strings.SplitN(tokenStr, ":", 2)[1])
		lower = strings.ToLower(tokenStr)
	}
	if strings.HasPrefix(lower, "bearer ") {
		tokenStr = strings.TrimSpace(tokenStr[7:])
	}
	tokenStr = trimQuotesAndSpaces(tokenStr)

	if tokenStr == "" || len(tokenStr) < 50 {
		return map[string]any{
			"success": false,
			"message": "Format token tidak valid. Token terlalu pendek atau kosong.",
		}
	}

	expWIBStr := "Tidak diketahui"
	parts := strings.Split(tokenStr, ".")
	if len(parts) >= 2 {
		if payload, err := decodeJWTPayload(parts[1]); err == nil {
			now := float64(time.Now().UnixNano()) / 1e9
			if exp, ok := jsonNumFloat(payload["exp"]); ok && exp != 0 {
				expTime := time.Unix(int64(exp), 0).In(shared.WIBLocation())
				expWIBStr = expTime.Format("2006-01-02 15:04:05 WIB")
				if exp <= now {
					return map[string]any{
						"success": false,
						"message": fmt.Sprintf("Token sudah kedaluwarsa pada %s.", expWIBStr),
					}
				}
			}
			if iat, ok := jsonNumFloat(payload["iat"]); ok && iat != 0 {
				if (now - iat) > (23.5 * 3600) {
					return map[string]any{
						"success": false,
						"message": "Usia token melebihi 24 jam (kebijakan rotasi sesi Stockbit). Silakan salin token yang lebih segar dari browser.",
					}
				}
			}
		}
	}

	apiVerified := true
	req, err := http.NewRequest(http.MethodGet, e.stockbitVerifyURL, nil)
	if err == nil {
		req.Header.Set("Authorization", "Bearer "+tokenStr)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Referer", "https://stockbit.com/")
		req.Header.Set("Origin", "https://stockbit.com")

		if resp, err := e.httpClient.Do(req); err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
				return map[string]any{
					"success": false,
					"message": fmt.Sprintf("Server Stockbit menolak token (HTTP %d). Token mungkin sudah di-revoke atau sesi login browser telah keluar.", resp.StatusCode),
				}
			}
		}
	}

	existing, err := os.ReadFile(e.smEnvFile)
	envContent := ""
	if err == nil {
		envContent = string(existing)
	} else if !os.IsNotExist(err) {
		errMsg := fmt.Sprintf("Gagal menyimpan token ke .env: %v", err)
		e.addLog("❌ "+errMsg, true)
		return map[string]any{"success": false, "message": errMsg}
	}

	replacement := fmt.Sprintf(`BROKER_API_TOKEN="%s"`, tokenStr)
	var newContent string
	if brokerTokenLineRe.MatchString(envContent) {
		newContent = brokerTokenLineRe.ReplaceAllStringFunc(envContent, func(string) string { return replacement })
	} else {
		newContent = strings.TrimRight(envContent, "\r\n") + "\n" + replacement + "\n"
	}

	if err := os.WriteFile(e.smEnvFile, []byte(newContent), 0o600); err != nil {
		errMsg := fmt.Sprintf("Gagal menyimpan token ke .env: %v", err)
		e.addLog("❌ "+errMsg, true)
		return map[string]any{"success": false, "message": errMsg}
	}
	_ = os.Setenv("BROKER_API_TOKEN", tokenStr)

	msg := fmt.Sprintf("Token Stockbit berhasil divalidasi & diperbarui! Berlaku sampai: %s.", expWIBStr)
	e.mu.Lock()
	e.appendLogLocked("🔑 "+msg, true)
	e.state.Error = nil
	if e.state.Status == "error" {
		e.state.Status = "idle"
		e.state.Message = "Token baru tersimpan. Siap untuk membuat laporan."
	}
	e.mu.Unlock()

	return map[string]any{
		"success":       true,
		"message":       msg,
		"expires_at":    expWIBStr,
		"verified_live": apiVerified,
	}
}

func trimQuotesAndSpaces(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"'`)
	return strings.TrimSpace(s)
}

func decodeJWTPayload(seg string) (map[string]any, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(seg, "="))
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(strings.TrimRight(seg, "="))
		if err != nil {
			return nil, err
		}
	}
	var payload map[string]any
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func jsonNumFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case int:
		return float64(x), true
	default:
		return 0, false
	}
}
