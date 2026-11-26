package foreign_flow

import (
	"math"
	"sort"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
)

// ComputeForeignHHI mirrors services/foreign_analytics.py compute_foreign_hhi.
func ComputeForeignHHI(brokerNetValues []float64) float64 {
	if len(brokerNetValues) == 0 {
		return 0.0
	}
	absVals := make([]float64, len(brokerNetValues))
	var total float64
	for i, v := range brokerNetValues {
		av := math.Abs(v)
		absVals[i] = av
		total += av
	}
	if total == 0 {
		return 0.0
	}
	var hhi float64
	for _, av := range absVals {
		share := av / total
		hhi += share * share
	}
	return hhi
}

// ComputeBrokerHeatmap mirrors services/foreign_analytics.py compute_broker_heatmap
// including pandas pivot_table (mean on collisions, fill_value=0) and sample std (ddof=1).
func ComputeBrokerHeatmap(rows []MicroRow, topN int) BrokerHeatmap {
	if len(rows) == 0 {
		return BrokerHeatmap{
			XDates:     []string{},
			YBrokers:   []string{},
			MatrixData: [][]any{},
		}
	}

	absTotals := make(map[string]float64)
	for _, r := range rows {
		absTotals[r.BrokerCode] += math.Abs(r.NetValue)
	}
	type bAbs struct {
		code string
		val  float64
	}
	ranked := make([]bAbs, 0, len(absTotals))
	for c, v := range absTotals {
		ranked = append(ranked, bAbs{code: c, val: v})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].val != ranked[j].val {
			return ranked[i].val > ranked[j].val
		}
		return ranked[i].code < ranked[j].code
	})
	if len(ranked) > topN {
		ranked = ranked[:topN]
	}

	topSet := make(map[string]bool, len(ranked))
	yBrokers := make([]string, len(ranked))
	for i, b := range ranked {
		topSet[b.code] = true
		yBrokers[i] = b.code
	}
	sort.Strings(yBrokers)

	dateSet := make(map[string]struct{})
	type cell struct {
		sum   float64
		count int
	}
	cells := make(map[string]map[string]cell)
	for _, r := range rows {
		if !topSet[r.BrokerCode] {
			continue
		}
		dateSet[r.Date] = struct{}{}
		byBroker, ok := cells[r.Date]
		if !ok {
			byBroker = make(map[string]cell)
			cells[r.Date] = byBroker
		}
		c := byBroker[r.BrokerCode]
		c.sum += r.NetValue
		c.count++
		byBroker[r.BrokerCode] = c
	}

	xDates := make([]string, 0, len(dateSet))
	for d := range dateSet {
		xDates = append(xDates, d)
	}
	sort.Strings(xDates)

	nD := len(xDates)
	nB := len(yBrokers)
	mat := make([][]float64, nD)
	for i, d := range xDates {
		mat[i] = make([]float64, nB)
		byBroker := cells[d]
		for j, b := range yBrokers {
			if c, ok := byBroker[b]; ok && c.count > 0 {
				mat[i][j] = c.sum / float64(c.count)
			}
		}
	}

	// Column-wise Z-score matching pandas pivot.mean() / pivot.std() (ddof=1)
	zMat := make([][]float64, nD)
	for i := range zMat {
		zMat[i] = make([]float64, nB)
	}
	for j := 0; j < nB; j++ {
		if nD <= 1 {
			continue
		}
		var sum float64
		for i := 0; i < nD; i++ {
			sum += mat[i][j]
		}
		mean := sum / float64(nD)
		var ss float64
		for i := 0; i < nD; i++ {
			d := mat[i][j] - mean
			ss += d * d
		}
		std := math.Sqrt(ss / float64(nD-1))
		if std == 0 || math.IsNaN(std) || math.IsInf(std, 0) {
			continue
		}
		for i := 0; i < nD; i++ {
			z := (mat[i][j] - mean) / std
			if math.IsNaN(z) || math.IsInf(z, 0) {
				z = 0
			} else if z < -1 {
				z = -1
			} else if z > 1 {
				z = 1
			}
			zMat[i][j] = shared.Round(z, 3)
		}
	}

	matrixData := make([][]any, 0, nD*nB)
	for i := 0; i < nD; i++ {
		for j := 0; j < nB; j++ {
			matrixData = append(matrixData, []any{i, j, zMat[i][j]})
		}
	}

	return BrokerHeatmap{
		XDates:     xDates,
		YBrokers:   yBrokers,
		MatrixData: matrixData,
	}
}
