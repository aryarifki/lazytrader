package shared

import (
	"math"
	"sort"
	"time"
)

type olsFit struct {
	ssr     float64
	dfResid float64
	tvalues []float64
}

// solveOLS fits y ~ X using column-scaled Householder QR decomposition.
func solveOLS(y []float64, xCols [][]float64) (olsFit, bool) {
	n := len(y)
	p := len(xCols)
	if n == 0 || p == 0 || n <= p {
		return olsFit{}, false
	}

	scales := make([]float64, p)
	A := make([][]float64, p)
	for j := 0; j < p; j++ {
		if len(xCols[j]) != n {
			return olsFit{}, false
		}
		var norm2 float64
		col := make([]float64, n)
		for i := 0; i < n; i++ {
			v := xCols[j][i]
			col[i] = v
			norm2 += v * v
		}
		s := math.Sqrt(norm2)
		if s == 0 || math.IsNaN(s) || math.IsInf(s, 0) {
			return olsFit{}, false
		}
		scales[j] = s
		invS := 1.0 / s
		for i := 0; i < n; i++ {
			col[i] *= invS
		}
		A[j] = col
	}

	qty := make([]float64, n)
	copy(qty, y)

	for k := 0; k < p; k++ {
		var norm2 float64
		for i := k; i < n; i++ {
			norm2 += A[k][i] * A[k][i]
		}
		norm := math.Sqrt(norm2)
		if norm < 1e-15 {
			return olsFit{}, false
		}
		alpha := -math.Copysign(norm, A[k][k])
		if alpha == 0 {
			alpha = norm
		}
		v := make([]float64, n-k)
		for i := k; i < n; i++ {
			v[i-k] = A[k][i]
		}
		v[0] -= alpha
		var vNorm2 float64
		for _, val := range v {
			vNorm2 += val * val
		}
		if vNorm2 == 0 {
			continue
		}
		twoOverV := 2.0 / vNorm2
		A[k][k] = alpha
		for i := k + 1; i < n; i++ {
			A[k][i] = 0
		}
		for j := k + 1; j < p; j++ {
			var dot float64
			for i := k; i < n; i++ {
				dot += v[i-k] * A[j][i]
			}
			factor := twoOverV * dot
			for i := k; i < n; i++ {
				A[j][i] -= factor * v[i-k]
			}
		}
		var dotY float64
		for i := k; i < n; i++ {
			dotY += v[i-k] * qty[i]
		}
		factorY := twoOverV * dotY
		for i := k; i < n; i++ {
			qty[i] -= factorY * v[i-k]
		}
	}

	for j := 0; j < p; j++ {
		if math.Abs(A[j][j]) < 1e-12 {
			return olsFit{}, false
		}
	}

	z := make([]float64, p)
	for i := p - 1; i >= 0; i-- {
		sum := qty[i]
		for j := i + 1; j < p; j++ {
			sum -= A[j][i] * z[j]
		}
		z[i] = sum / A[i][i]
	}

	beta := make([]float64, p)
	for j := 0; j < p; j++ {
		beta[j] = z[j] / scales[j]
	}

	var ssr float64
	for i := 0; i < n; i++ {
		pred := 0.0
		for j := 0; j < p; j++ {
			pred += beta[j] * xCols[j][i]
		}
		res := y[i] - pred
		ssr += res * res
	}

	dfResid := float64(n - p)
	s2 := ssr / dfResid
	tvals := make([]float64, p)

	// R^{-1} columns for covariance diagonal
	rInv := make([][]float64, p)
	for colIdx := 0; colIdx < p; colIdx++ {
		col := make([]float64, p)
		for i := colIdx; i >= 0; i-- {
			rhs := 0.0
			if i == colIdx {
				rhs = 1.0
			}
			for j := i + 1; j <= colIdx; j++ {
				rhs -= A[j][i] * col[j]
			}
			col[i] = rhs / A[i][i]
		}
		rInv[colIdx] = col
	}

	for j := 0; j < p; j++ {
		var rowNorm2 float64
		for colIdx := j; colIdx < p; colIdx++ {
			val := rInv[colIdx][j]
			rowNorm2 += val * val
		}
		varB := s2 * rowNorm2 / (scales[j] * scales[j])
		if varB > 0 {
			tvals[j] = beta[j] / math.Sqrt(varB)
		}
	}

	return olsFit{ssr: ssr, dfResid: dfResid, tvalues: tvals}, true
}

// MacKinnonP computes statsmodels MacKinnon (1994) p-value for regression="c", N=1.
func MacKinnonP(testStat float64) float64 {
	const (
		maxStat  = 2.74
		minStat  = -18.83
		starStat = -1.61
	)
	if testStat > maxStat {
		return 1.0
	}
	if testStat < minStat {
		return 0.0
	}
	var poly float64
	if testStat <= starStat {
		// smallp coefficients: [2.1659, 1.4412, 0.038269]
		poly = (0.038269*testStat+1.4412)*testStat + 2.1659
	} else {
		// largep coefficients: [1.7339, 0.93202, -0.12745, -0.010368]
		poly = ((-0.010368*testStat-0.12745)*testStat+0.93202)*testStat + 1.7339
	}
	return NormCDF(poly)
}

// ADFullerPValue replicates statsmodels.tsa.stattools.adfuller(x, regression="c", autolag="AIC").
func ADFullerPValue(x []float64) (float64, bool) {
	nobs := len(x)
	if nobs < 10 {
		return 1.0, false
	}
	minX, maxX := x[0], x[0]
	for _, v := range x[1:] {
		if v < minX {
			minX = v
		}
		if v > maxX {
			maxX = v
		}
	}
	if minX == maxX {
		return 1.0, false
	}

	ntrend := 1 // regression = "c"
	maxLag := int(math.Ceil(12.0 * math.Pow(float64(nobs)/100.0, 0.25)))
	limitLag := nobs/2 - ntrend - 1
	if limitLag < maxLag {
		maxLag = limitLag
	}
	if maxLag < 0 {
		return 1.0, false
	}

	xdiff := make([]float64, nobs-1)
	for i := 0; i < nobs-1; i++ {
		xdiff[i] = x[i+1] - x[i]
	}

	// _autolag compares AIC over identical trimmed nEff = len(xdiff) - maxLag observations.
	nEff := len(xdiff) - maxLag
	if nEff <= maxLag+2 {
		return 1.0, false
	}
	yAutolag := xdiff[maxLag:]
	onesAuto := make([]float64, nEff)
	levelAuto := make([]float64, nEff)
	for i := 0; i < nEff; i++ {
		onesAuto[i] = 1.0
		levelAuto[i] = x[nobs-nEff-1+i]
	}
	lagColsAuto := make([][]float64, 0, maxLag+2)
	lagColsAuto = append(lagColsAuto, onesAuto, levelAuto)
	for lag := 1; lag <= maxLag; lag++ {
		col := make([]float64, nEff)
		for i := 0; i < nEff; i++ {
			col[i] = xdiff[maxLag-lag+i]
		}
		lagColsAuto = append(lagColsAuto, col)
	}

	bestAIC := math.Inf(1)
	bestDiffLag := 0
	for dl := 0; dl <= maxLag; dl++ {
		kCols := 2 + dl
		fit, ok := solveOLS(yAutolag, lagColsAuto[:kCols])
		if !ok || fit.ssr <= 0 {
			continue
		}
		// Match statsmodels.OLS.aic ordering (constant term -n/2*(ln(2pi)+1) cancels)
		aic := float64(nEff)*math.Log(fit.ssr/float64(nEff)) + 2.0*float64(kCols)
		if aic < bestAIC {
			bestAIC = aic
			bestDiffLag = dl
		}
	}

	// Re-fit OLS on trimmed rows matching bestDiffLag: nFinal = len(xdiff) - bestDiffLag
	nFinal := len(xdiff) - bestDiffLag
	yFinal := xdiff[bestDiffLag:]
	levelFinal := make([]float64, nFinal)
	onesFinal := make([]float64, nFinal)
	for i := 0; i < nFinal; i++ {
		onesFinal[i] = 1.0
		levelFinal[i] = x[nobs-nFinal-1+i]
	}
	colsFinal := make([][]float64, 0, bestDiffLag+2)
	colsFinal = append(colsFinal, levelFinal) // index 0 is level -> adfstat = tvalues[0]
	for lag := 1; lag <= bestDiffLag; lag++ {
		col := make([]float64, nFinal)
		for i := 0; i < nFinal; i++ {
			col[i] = xdiff[bestDiffLag-lag+i]
		}
		colsFinal = append(colsFinal, col)
	}
	colsFinal = append(colsFinal, onesFinal)

	finalFit, ok := solveOLS(yFinal, colsFinal)
	if !ok {
		return 1.0, false
	}
	return MacKinnonP(finalFit.tvalues[0]), true
}

// MakeStationary applies ADF test (p > 0.05 -> first difference), matching analysis.py _make_stationary.
func MakeStationary(series []float64) ([]float64, []bool) {
	n := len(series)
	valid := make([]bool, n)
	out := make([]float64, n)
	copy(out, series)
	for i := range valid {
		valid[i] = !math.IsNaN(series[i]) && !math.IsInf(series[i], 0)
	}
	clean := make([]float64, 0, n)
	for i, v := range series {
		if valid[i] {
			clean = append(clean, v)
		}
	}
	_, std := MeanStd(clean)
	if std == 0 || len(clean) < 10 {
		return out, valid
	}
	pVal, ok := ADFullerPValue(clean)
	if ok && pVal > 0.05 {
		diff := make([]float64, n)
		diffValid := make([]bool, n)
		diffValid[0] = false
		for i := 1; i < n; i++ {
			if valid[i] && valid[i-1] {
				diff[i] = series[i] - series[i-1]
				diffValid[i] = true
			}
		}
		return diff, diffValid
	}
	return out, valid
}

// GrangerSSRFTest computes statsmodels grangercausalitytests ssr_ftest p-values for lags 1..maxLag.
func GrangerSSRFTest(y, x []float64, maxLag int) (map[int]float64, bool) {
	n := len(y)
	if len(x) != n || maxLag <= 0 || n <= 3*maxLag+1 {
		return nil, false
	}
	pvals := make(map[int]float64, maxLag)
	for lag := 1; lag <= maxLag; lag++ {
		m := n - lag
		if m <= 2*lag+1 {
			return nil, false
		}
		target := y[lag:]
		ones := make([]float64, m)
		for i := 0; i < m; i++ {
			ones[i] = 1.0
		}
		ownCols := make([][]float64, 0, lag+1)
		jointCols := make([][]float64, 0, 2*lag+1)
		for l := 1; l <= lag; l++ {
			yc := make([]float64, m)
			for i := 0; i < m; i++ {
				yc[i] = y[lag-l+i]
			}
			ownCols = append(ownCols, yc)
			jointCols = append(jointCols, yc)
		}
		for l := 1; l <= lag; l++ {
			xc := make([]float64, m)
			for i := 0; i < m; i++ {
				xc[i] = x[lag-l+i]
			}
			jointCols = append(jointCols, xc)
		}
		ownCols = append(ownCols, ones)
		jointCols = append(jointCols, ones)

		fitOwn, okOwn := solveOLS(target, ownCols)
		fitJoint, okJoint := solveOLS(target, jointCols)
		if !okOwn || !okJoint || fitJoint.ssr <= 0 || fitJoint.dfResid <= 0 {
			return nil, false
		}
		my, _ := MeanStd(target)
		var tss float64
		for _, v := range target {
			d := v - my
			tss += d * d
		}
		if tss == 0 || (fitJoint.ssr/tss) < 2.220446049250313e-16 {
			return nil, false
		}
		fgc := ((fitOwn.ssr - fitJoint.ssr) / fitJoint.ssr) / float64(lag) * fitJoint.dfResid
		p := FSF(fgc, float64(lag), fitJoint.dfResid)
		if math.IsNaN(p) {
			return nil, false
		}
		pvals[lag] = p
	}
	return pvals, true
}

// FlowInput holds daily aggregate foreign net flow.
type FlowInput struct {
	Date             string
	ForeignNetBroker float64
}

// CausalityResult mirrors analysis.causality_foreign_vs_price output.
type CausalityResult struct {
	Ticker       string  `json:"ticker"`
	BestLag      int     `json:"best_lag"`
	MinPValue    float64 `json:"min_p_value"`
	RawMinPValue float64 `json:"raw_min_p_value"`
	IsSignificant bool   `json:"is_significant"`
	NObs         int     `json:"n_obs"`
}

// CausalityForeignVsPrice replicates analysis.causality_foreign_vs_price.
func CausalityForeignVsPrice(ticker string, prices []PriceInput, flows []FlowInput, maxLags, lookbackDays int) *CausalityResult {
	if len(prices) == 0 || len(flows) == 0 {
		return nil
	}
	px := make([]PriceInput, 0, len(prices))
	for _, p := range prices {
		if p.Close > 0 {
			px = append(px, p)
		}
	}
	sort.Slice(px, func(i, j int) bool { return px[i].Date < px[j].Date })
	if lookbackDays > 0 && len(px) > lookbackDays {
		px = px[len(px)-lookbackDays:]
	}
	if len(px) == 0 {
		return nil
	}

	fl := make([]FlowInput, len(flows))
	copy(fl, flows)
	sort.Slice(fl, func(i, j int) bool { return fl[i].Date < fl[j].Date })
	if lookbackDays > 0 && len(fl) > lookbackDays {
		fl = fl[len(fl)-lookbackDays:]
	}
	flowByDate := make(map[string]float64, len(fl))
	for _, f := range fl {
		flowByDate[f.Date] += f.ForeignNetBroker
	}

	rets := make([]float64, len(px))
	retValid := make([]bool, len(px))
	fSeries := make([]float64, len(px))
	for i := range px {
		if i > 0 {
			rets[i] = px[i].Close/px[i-1].Close - 1.0
			retValid[i] = true
		}
		fSeries[i] = flowByDate[px[i].Date]
	}

	statFlow, flowValid := MakeStationary(fSeries)
	y := make([]float64, 0, len(px))
	x := make([]float64, 0, len(px))
	for i := range px {
		if retValid[i] && flowValid[i] {
			y = append(y, rets[i])
			x = append(x, statFlow[i])
		}
	}
	_, stdY := MeanStd(y)
	_, stdX := MeanStd(x)
	if len(y) < maxLags*3+2 || stdY == 0 || stdX == 0 {
		return nil
	}
	pvals, ok := GrangerSSRFTest(y, x, maxLags)
	if !ok {
		return nil
	}
	bestLag, minAdjP, rawMinP, ok := HolmBonferroniLags(pvals)
	if !ok {
		return nil
	}
	return &CausalityResult{
		Ticker:        ticker,
		BestLag:       bestLag,
		MinPValue:     minAdjP,
		RawMinPValue:  rawMinP,
		IsSignificant: minAdjP < 0.05,
		NObs:          len(y),
	}
}

// LookbackStartDate matches Python (date.today() - timedelta(days=int(lookbackDays * 1.6))).
func LookbackStartDate(lookbackDays int) time.Time {
	daysBack := int(float64(lookbackDays) * 1.6)
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -daysBack)
}
