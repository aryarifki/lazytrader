package flags

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFeatureFlags_ParseAndOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "feature_flags.yaml")
	content := `
# Comment line
SCORING_THRESHOLD_FIX: true
SIGNAL_COOLDOWN_DAYS: 5
EXODUS_THRESHOLD_MODE: "hybrid_p5_p15" # inline comment
AI_GATE_MODE: shadow
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	st := NewStore(path)
	if v, ok := st.Get("SCORING_THRESHOLD_FIX"); !ok || !v.AsBool() {
		t.Fatalf("expected SCORING_THRESHOLD_FIX=true, got %+v ok=%v", v, ok)
	}
	if v, ok := st.Get("SIGNAL_COOLDOWN_DAYS"); !ok || v.AsInt(0) != 5 {
		t.Fatalf("expected SIGNAL_COOLDOWN_DAYS=5, got %+v ok=%v", v, ok)
	}
	if v, ok := st.Get("EXODUS_THRESHOLD_MODE"); !ok || v.Raw != "hybrid_p5_p15" {
		t.Fatalf("expected hybrid_p5_p15, got %q", v.Raw)
	}

	t.Setenv("FLAG_EXODUS_THRESHOLD_MODE", "absolute")
	if v, ok := st.Get("EXODUS_THRESHOLD_MODE"); !ok || v.Raw != "absolute" {
		t.Fatalf("expected env override absolute, got %q", v.Raw)
	}
}
