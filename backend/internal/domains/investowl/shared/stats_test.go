package shared

import (
	"math"
	"math/rand/v2"
	"testing"
)

func almostEq(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestDistributionsVsSciPy(t *testing.T) {
	if got := StudentTSF(2.1, 15); !almostEq(got, 0.02652762807602138, 1e-10) {
		t.Fatalf("StudentTSF(2.1, 15) = %.15f, want 0.02652762807602138", got)
	}
	if got := StudentTSF(-1.5, 10); !almostEq(got, 0.9177463367772799, 1e-10) {
		t.Fatalf("StudentTSF(-1.5, 10) = %.15f, want 0.9177463367772799", got)
	}
	if got := FSF(3.2, 3, 40); !almostEq(got, 0.03342118072289294, 1e-10) {
		t.Fatalf("FSF(3.2, 3, 40) = %.15f, want 0.03342118072289294", got)
	}
}

func TestADFullerVsStatsmodels(t *testing.T) {
	rw := []float64{
		0.304717079754, -0.735267026486, 0.01518416932, 0.955748885712, -0.995286302942,
		-2.297465809805, -2.169625406637, -2.485867998981, -2.502669156485, -3.355713084059,
		-2.476315109196, -1.698523173767, -1.632492476206, -0.505251269238, -0.037741926986,
		-0.897034389869, -0.528283605786, -1.487166206615, -0.608715905308, -0.658641816294,
		-0.84350417984, -1.524433724244, -0.30189238557, -0.456421867638, -0.884749689801,
		-1.23688324029, -0.704574054736, -0.339129990372, 0.073602621224, 0.504423624232,
		2.646071225102, 2.239656208717, 1.727413479646, 0.913640751398, 1.529620173974,
		2.658592466694, 2.54464500904, 1.704488532077, 0.880007316386, 1.530600104211,
		2.273854275414, 2.817008543719, 2.15149883643, 2.383660159497, 2.500345968638,
		2.719034565367, 3.590463343315, 3.81405889209, 4.492972455162, 4.560551524651,
		4.849670923341, 5.480959149179, 4.023803329323, 3.704132112966, 3.233759458673,
		2.59488161043, 2.319739359203, 3.814680670438, 2.948849554744, 3.917127909336,
		2.23425813772, 1.899373107734, 2.062126172839, 2.648348504199, 3.359575083991,
		4.152922319191, 3.804197246943, 3.341845454278, 4.199821335536, 4.008517010654,
		2.732830687316, 1.599543473313, 0.680091187311, 1.177251931365, 1.319677667435,
		2.010163021503, 1.582910375167, 1.741450066243, 2.367040460211, 2.05769392049,
	}
	p, ok := ADFullerPValue(rw)
	if !ok {
		t.Fatalf("ADFullerPValue failed")
	}
	if !almostEq(p, 0.44754374269602637, 1e-6) {
		t.Fatalf("ADFullerPValue(rw) = %.12f, want 0.44754374269602637", p)
	}
}

func TestGrangerAndStatsHelpers(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 7))
	n := 80
	x := make([]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = rng.NormFloat64()
		if i > 0 {
			y[i] = 0.65*x[i-1] + 0.3*rng.NormFloat64()
		}
	}
	pvals, ok := GrangerSSRFTest(y, x, 5)
	if !ok {
		t.Fatalf("GrangerSSRFTest failed")
	}
	if pvals[1] >= 0.01 {
		t.Fatalf("expected strong Granger lag-1 causality p < 0.01, got %f", pvals[1])
	}

	if q := LinearPercentile([]float64{10, 20, 30, 40}, 25); !almostEq(q, 17.5, 1e-12) {
		t.Fatalf("LinearPercentile = %v, want 17.5", q)
	}

	qvs := BenjaminiHochberg([]float64{0.01, 0.04, 0.20})
	if !almostEq(qvs[0], 0.03, 1e-9) || !almostEq(qvs[1], 0.06, 1e-9) || !almostEq(qvs[2], 0.20, 1e-9) {
		t.Fatalf("BenjaminiHochberg = %v", qvs)
	}

	se := NeweyWestSE([]float64{0.01, 0.02, -0.01, 0.03, 0.015}, 2)
	if math.IsNaN(se) || se <= 0 {
		t.Fatalf("invalid NeweyWestSE: %v", se)
	}
}
