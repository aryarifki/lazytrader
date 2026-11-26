package shared

import (
	"math"
	"sort"
	"strings"
)

// NormCDF returns the standard normal cumulative distribution function P(Z <= z).
func NormCDF(z float64) float64 {
	return 0.5 * math.Erfc(-z/math.Sqrt2)
}

func logBeta(a, b float64) float64 {
	la, _ := math.Lgamma(a)
	lb, _ := math.Lgamma(b)
	lab, _ := math.Lgamma(a + b)
	return la + lb - lab
}

func betacf(a, b, x float64) float64 {
	const maxIter = 200
	const eps = 3.0e-14
	const fpMin = 1.0e-300

	qab := a + b
	qap := a + 1.0
	qam := a - 1.0
	c := 1.0
	d := 1.0 - qab*x/qap
	if math.Abs(d) < fpMin {
		d = fpMin
	}
	d = 1.0 / d
	h := d

	for m := 1; m <= maxIter; m++ {
		m2 := float64(2 * m)
		mf := float64(m)
		aa := mf * (b - mf) * x / ((qam + m2) * (a + m2))
		d = 1.0 + aa*d
		if math.Abs(d) < fpMin {
			d = fpMin
		}
		c = 1.0 + aa/c
		if math.Abs(c) < fpMin {
			c = fpMin
		}
		d = 1.0 / d
		h *= d * c

		aa = -(a + mf) * (qab + mf) * x / ((a + m2) * (qap + m2))
		d = 1.0 + aa*d
		if math.Abs(d) < fpMin {
			d = fpMin
		}
		c = 1.0 + aa/c
		if math.Abs(c) < fpMin {
			c = fpMin
		}
		d = 1.0 / d
		del := d * c
		h *= del
		if math.Abs(del-1.0) <= eps {
			break
		}
	}
	return h
}

// IncompleteBeta computes the regularized incomplete beta function I_x(a, b).
func IncompleteBeta(a, b, x float64) float64 {
	if x <= 0.0 {
		return 0.0
	}
	if x >= 1.0 {
		return 1.0
	}
	bt := math.Exp(logBeta(a, b)*-1.0 + a*math.Log(x) + b*math.Log(1.0-x))
	if x < (a+1.0)/(a+b+2.0) {
		return bt * betacf(a, b, x) / a
	}
	return 1.0 - bt*betacf(b, a, 1.0-x)/b
}

// StudentTSF computes the upper tail survival function P(T > t) for Student's t with df degrees of freedom.
func StudentTSF(t, df float64) float64 {
	if math.IsNaN(t) || math.IsNaN(df) || df <= 0 {
		return math.NaN()
	}
	x := df / (df + t*t)
	halfIx := 0.5 * IncompleteBeta(df/2.0, 0.5, x)
	if t >= 0 {
		return halfIx
	}
	return 1.0 - halfIx
}

// FSF computes the upper tail survival function P(F > f) for F(df1, df2).
func FSF(f, df1, df2 float64) float64 {
	if math.IsNaN(f) || math.IsNaN(df1) || math.IsNaN(df2) || df1 <= 0 || df2 <= 0 {
		return math.NaN()
	}
	if f <= 0 {
		return 1.0
	}
	x := df2 / (df2 + df1*f)
	return IncompleteBeta(df2/2.0, df1/2.0, x)
}

// LinearPercentile computes the NumPy default linear-interpolated percentile (q in [0, 100]).
func LinearPercentile(vals []float64, q float64) float64 {
	if len(vals) == 0 {
		return 0.0
	}
	sorted := make([]float64, len(vals))
	copy(sorted, vals)
	sort.Float64s(sorted)
	if len(sorted) == 1 || q <= 0 {
		return sorted[0]
	}
	if q >= 100 {
		return sorted[len(sorted)-1]
	}
	idx := (q / 100.0) * float64(len(sorted)-1)
	lo := int(math.Floor(idx))
	hi := int(math.Ceil(idx))
	if lo == hi {
		return sorted[lo]
	}
	frac := idx - float64(lo)
	return sorted[lo]*(1.0-frac) + sorted[hi]*frac
}

// MeanStd computes population mean and standard deviation (ddof=0, matching NumPy default).
func MeanStd(vals []float64) (mean, std float64) {
	if len(vals) == 0 {
		return 0, 0
	}
	var sum float64
	for _, v := range vals {
		sum += v
	}
	mean = sum / float64(len(vals))
	var sq float64
	for _, v := range vals {
		d := v - mean
		sq += d * d
	}
	std = math.Sqrt(sq / float64(len(vals)))
	return mean, std
}

// PearsonCorr computes Pearson correlation (zero std -> 0.0) matching np.corrcoef.
func PearsonCorr(x, y []float64) float64 {
	if len(x) < 3 || len(x) != len(y) {
		return 0.0
	}
	mx, sx := MeanStd(x)
	my, sy := MeanStd(y)
	if sx <= 0 || sy <= 0 {
		return 0.0
	}
	var cov float64
	for i := range x {
		cov += (x[i] - mx) * (y[i] - my)
	}
	corr := (cov / float64(len(x))) / (sx * sy)
	if math.IsNaN(corr) || math.IsInf(corr, 0) {
		return 0.0
	}
	return Round(corr, 3)
}

// NeweyWestSE mirrors idx_bandarmology/analysis.py _newey_west_se.
func NeweyWestSE(vals []float64, maxLags int) float64 {
	clean := make([]float64, 0, len(vals))
	for _, v := range vals {
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			clean = append(clean, v)
		}
	}
	n := len(clean)
	if n < 2 {
		return math.NaN()
	}
	var sum float64
	for _, v := range clean {
		sum += v
	}
	mean := sum / float64(n)
	demeaned := make([]float64, n)
	var dot0 float64
	for i, v := range clean {
		d := v - mean
		demeaned[i] = d
		dot0 += d * d
	}
	gamma0 := dot0 / float64(n-1)
	if gamma0 <= 0 {
		return math.NaN()
	}
	lrv := gamma0
	maxL := maxLags
	if maxL > n-2 {
		maxL = n - 2
	}
	if maxL < 0 {
		maxL = 0
	}
	for lag := 1; lag <= maxL; lag++ {
		w := 1.0 - float64(lag)/(float64(maxL)+1.0)
		var dotK float64
		for i := lag; i < n; i++ {
			dotK += demeaned[i] * demeaned[i-lag]
		}
		gammaK := dotK / float64(n-1)
		lrv += 2.0 * w * gammaK
	}
	lrv = math.Max(lrv, math.Max(gamma0*0.1, 1e-12))
	return math.Sqrt(lrv / float64(n))
}

// BenjaminiHochberg computes BH FDR-adjusted q-values preserving NaNs.
func BenjaminiHochberg(pvals []float64) []float64 {
	out := make([]float64, len(pvals))
	type pair struct {
		idx int
		p   float64
	}
	valid := make([]pair, 0, len(pvals))
	for i, p := range pvals {
		out[i] = math.NaN()
		if !math.IsNaN(p) && !math.IsInf(p, 0) {
			valid = append(valid, pair{idx: i, p: p})
		}
	}
	m := len(valid)
	if m == 0 {
		return out
	}
	sort.SliceStable(valid, func(i, j int) bool { return valid[i].p < valid[j].p })
	qSorted := make([]float64, m)
	for i, v := range valid {
		qSorted[i] = v.p * float64(m) / float64(i+1)
	}
	for i := m - 2; i >= 0; i-- {
		if qSorted[i+1] < qSorted[i] {
			qSorted[i] = qSorted[i+1]
		}
	}
	for i, q := range qSorted {
		if q < 0 {
			q = 0
		} else if q > 1 {
			q = 1
		}
		out[valid[i].idx] = q
	}
	return out
}

// HolmBonferroniLags applies step-down Holm correction across tested lags.
func HolmBonferroniLags(pvals map[int]float64) (bestLag int, minAdjP, rawMinP float64, ok bool) {
	type lagP struct {
		lag int
		p   float64
	}
	items := make([]lagP, 0, len(pvals))
	for lag, p := range pvals {
		if !math.IsNaN(p) && !math.IsInf(p, 0) {
			items = append(items, lagP{lag: lag, p: p})
		}
	}
	if len(items) == 0 {
		return 0, 0, 0, false
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].p != items[j].p {
			return items[i].p < items[j].p
		}
		return items[i].lag < items[j].lag
	})
	m := len(items)
	prevAdj := 0.0
	adj0 := 1.0
	for i, it := range items {
		adj := math.Min(1.0, float64(m-i)*it.p)
		if prevAdj > adj {
			adj = prevAdj
		}
		prevAdj = adj
		if i == 0 {
			adj0 = adj
		}
	}
	return items[0].lag, adj0, items[0].p, true
}

// ActivityInput holds minimal per-broker daily activity needed for alpha scanning.
type ActivityInput struct {
	Date       string
	Ticker     string
	BrokerCode string
	NetValue   float64
}

// PriceInput holds minimal daily close needed for forward returns & causality.
type PriceInput struct {
	Date   string
	Ticker string
	Close  float64
}

// AlphaScanRow mirrors one row returned by analysis.broker_alpha_scan.
type AlphaScanRow struct {
	Ticker          string
	BrokerCode      string
	NEvents         int
	MeanFwdReturn   float64
	MedianFwdReturn float64
	WinRate         float64
	AvgNetValue     float64
	TotalNetValue   float64
	SEHAC           float64
	TStat           float64
	PValueOneSided  float64
	PValueFDR       float64
	Significant     bool
	Status          string
}

// BrokerAlphaScan mirrors analysis.broker_alpha_scan grouped by (ticker, broker_code).
func BrokerAlphaScan(activities []ActivityInput, prices []PriceInput, horizon, minEvents int, minNetValue float64) []AlphaScanRow {
	if len(activities) == 0 || len(prices) == 0 || horizon <= 0 {
		return nil
	}
	pxByTicker := make(map[string][]PriceInput)
	for _, p := range prices {
		if p.Close > 0 {
			tk := strings.ToUpper(strings.TrimSpace(p.Ticker))
			pxByTicker[tk] = append(pxByTicker[tk], PriceInput{Date: p.Date, Ticker: tk, Close: p.Close})
		}
	}
	fwdRet := make(map[string]map[string]float64, len(pxByTicker))
	for tk, list := range pxByTicker {
		sort.Slice(list, func(i, j int) bool { return list[i].Date < list[j].Date })
		m := make(map[string]float64, len(list))
		for i := 0; i+horizon < len(list); i++ {
			m[list[i].Date] = list[i+horizon].Close/list[i].Close - 1.0
		}
		fwdRet[tk] = m
	}

	type tbKey struct {
		ticker string
		broker string
	}
	type ev struct {
		date   string
		netVal float64
		ret    float64
	}
	grouped := make(map[tbKey][]ev)
	for _, a := range activities {
		if !(a.NetValue > minNetValue) {
			continue
		}
		tk := strings.ToUpper(strings.TrimSpace(a.Ticker))
		m, ok := fwdRet[tk]
		if !ok {
			continue
		}
		r, ok := m[a.Date]
		if !ok {
			continue
		}
		k := tbKey{ticker: tk, broker: a.BrokerCode}
		grouped[k] = append(grouped[k], ev{date: a.Date, netVal: a.NetValue, ret: r})
	}

	minReq := minEvents
	if minReq < 5 {
		minReq = 5
	}

	out := make([]AlphaScanRow, 0, len(grouped))
	byTickerIdx := make(map[string][]int)

	for k, evs := range grouped {
		if len(evs) < minEvents {
			continue
		}
		sort.Slice(evs, func(i, j int) bool { return evs[i].date < evs[j].date })
		n := len(evs)
		vals := make([]float64, n)
		var sumRet, sumNet float64
		wins := 0
		for i, e := range evs {
			vals[i] = e.ret
			sumRet += e.ret
			sumNet += e.netVal
			if e.ret > 0 {
				wins++
			}
		}
		meanRet := sumRet / float64(n)
		sortedVals := append([]float64(nil), vals...)
		sort.Float64s(sortedVals)
		var medianRet float64
		if n%2 == 1 {
			medianRet = sortedVals[n/2]
		} else {
			medianRet = 0.5 * (sortedVals[n/2-1] + sortedVals[n/2])
		}
		winRate := float64(wins) / float64(n)

		maxL := horizon - 1
		if maxL > n-2 {
			maxL = n - 2
		}
		if maxL < 0 {
			maxL = 0
		}
		seHAC := NeweyWestSE(vals, maxL)
		tStat := math.NaN()
		pVal := math.NaN()
		if !math.IsNaN(seHAC) && seHAC > 0 {
			tStat = meanRet / seHAC
			if n >= 2 {
				pVal = StudentTSF(tStat, float64(n-1))
			}
		}
		status := "need >=5 events"
		if n >= 5 {
			status = "ok"
		}
		idx := len(out)
		out = append(out, AlphaScanRow{
			Ticker:          k.ticker,
			BrokerCode:      k.broker,
			NEvents:         n,
			MeanFwdReturn:   meanRet,
			MedianFwdReturn: medianRet,
			WinRate:         winRate,
			AvgNetValue:     sumNet / float64(n),
			TotalNetValue:   sumNet,
			SEHAC:           seHAC,
			TStat:           tStat,
			PValueOneSided:  pVal,
			PValueFDR:       math.NaN(),
			Status:          status,
		})
		byTickerIdx[k.ticker] = append(byTickerIdx[k.ticker], idx)
	}

	if len(out) == 0 {
		return out
	}

	for _, idxs := range byTickerIdx {
		pvs := make([]float64, len(idxs))
		for i, id := range idxs {
			pvs[i] = out[id].PValueOneSided
		}
		qvs := BenjaminiHochberg(pvs)
		for i, id := range idxs {
			q := qvs[i]
			out[id].PValueFDR = q
			out[id].Significant = out[id].NEvents >= minReq && !math.IsNaN(q) && q < 0.05 && out[id].MeanFwdReturn > 0
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Significant != b.Significant {
			return a.Significant
		}
		pa, pb := a.PValueFDR, b.PValueFDR
		if math.IsNaN(pa) != math.IsNaN(pb) {
			return !math.IsNaN(pa)
		}
		if !math.IsNaN(pa) && pa != pb {
			return pa < pb
		}
		oa, ob := a.PValueOneSided, b.PValueOneSided
		if math.IsNaN(oa) != math.IsNaN(ob) {
			return !math.IsNaN(oa)
		}
		if !math.IsNaN(oa) && oa != ob {
			return oa < ob
		}
		return a.MeanFwdReturn > b.MeanFwdReturn
	})
	return out
}
