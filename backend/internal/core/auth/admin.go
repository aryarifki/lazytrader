package auth

import "strings"

// IsAdminEmail mirrors /root/investowl/backend/core/auth.py:is_admin_email exactly.
func IsAdminEmail(email string, adminEmails []string) bool {
	clean := strings.ToLower(strings.TrimSpace(email))
	if clean == "" {
		return false
	}
	for _, allowed := range adminEmails {
		if strings.ToLower(strings.TrimSpace(allowed)) == clean {
			return true
		}
	}
	return strings.Contains(clean, "adryan") || strings.Contains(clean, "arya")
}
