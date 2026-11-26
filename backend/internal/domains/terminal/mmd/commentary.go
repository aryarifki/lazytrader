package mmd

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func generateCommentaries(
	in AnalysisInput,
	sortedAll []DailyCandle,
	dailyResumes []sliceResume,
	buyers, sellers []brokerNet,
	topBuyers []TopBrokerRow,
	pengendali PengendaliStats,
	actionCls, retailAction Badge,
	aiScore int,
	korelasiRaw, vr, avgDailyVal float64,
	totNetVol, totalVol int64,
	lastClose, priceChange float64,
	today time.Time,
) ([]string, []string) {
	_ = totalVol
	var ipoDate *string
	if in.Meta != nil {
		ipoDate = in.Meta.IPODate
	}
	maxOther := math.Max(pengendali.SemiRitel.Main, pengendali.Institusi.Main)

	candidates := []string{
		commentaryIPO(in.Ticker, ipoDate, today),
		commentaryRentangPeriode(len(in.Candles)),
		commentarySection2BigPlayer(in.T5Records),
		commentarySection3Partisipasi(pengendali.Ritel.Main, maxOther),
		commentarySection3Point5(pengendali.Ritel.Main, avgDailyVal/1e9),
		commentaryAkumulasiBerturutTurut(sortedAll, dailyResumes),
		commentaryLiveCandleRise(sortedAll, priceChange),
		commentaryVolumeRatio(vr),
		commentaryRataRataTransaksi(avgDailyVal, totNetVol),
		commentaryRitelMinim(buyers, sellers),
		commentaryKoreksiMarkdownRetail(sortedAll),
		commentaryLonjakanTransaksi(sortedAll),
		commentaryRitelBeliPuncak(sortedAll, dailyResumes),
		commentarySidewaysAkumulasi(sortedAll),
		commentaryKinerjaGroup1Hari(in.GroupPerf),
		commentaryKinerjaGroup1Minggu(in.GroupPerf),
		commentaryBrokerBaru(sortedAll),
		commentaryMarkupPeriode(priceChange, actionCls.Cls, len(in.Candles)),
		commentaryPeriodeHarga(priceChange, actionCls.Cls, actionCls.Label, len(in.Candles)),
		commentaryPeriodeTransaksi(sortedAll),
		commentaryAverageAkumulasiBandar(totNetVol, len(in.Candles), lastClose),
		commentaryKorelasi(korelasiRaw),
		commentaryKorelasiParuh(sortedAll, dailyResumes),
		commentaryVolatilitas(sortedAll),
		commentaryConclusion(actionCls.Label, retailAction.Label, aiScore),
	}

	active := make([]string, 0, len(candidates))
	for _, msg := range candidates {
		if msg != "" && msg != "-" {
			active = append(active, msg)
		}
	}
	if len(active) == 0 {
		active = []string{
			fmt.Sprintf("Pergerakan harga saham %s dalam periode ini mencatatkan perubahan %+.2f%%.", in.Ticker, priceChange),
			fmt.Sprintf("Akumulasi bandar berada pada status %s dengan nilai transaksi rata-rata Rp %.2f Miliar.", actionCls.Label, avgDailyVal/1e9),
			fmt.Sprintf("Dominasi broker pengendali memegang kendali partisipasi pasar sebesar %s%%.", formatPyFloat(pengendali.Ratio)),
			commentaryConclusion(actionCls.Label, retailAction.Label, aiScore),
		}
	}

	topBuyerBroker := "-"
	if len(topBuyers) > 0 {
		topBuyerBroker = topBuyers[0].Broker
	}
	bigPlayerSummary := []string{
		fmt.Sprintf("Investor institusi dan pemegang saham mayoritas terpantau aktif melakukan transaksi dengan dominasi broker %s di sisi beli.", topBuyerBroker),
		commentarySection2BigPlayer(in.T5Records),
	}
	return active, bigPlayerSummary
}

func commentaryIPO(ticker string, ipoDate *string, today time.Time) string {
	if ipoDate == nil || *ipoDate == "" {
		return "-"
	}
	ipoDt, err := time.Parse("2006-01-02", *ipoDate)
	if err != nil {
		return "-"
	}
	months := (today.Year()-ipoDt.Year())*12 + int(today.Month()-ipoDt.Month())
	if today.Day() < ipoDt.Day() {
		months--
	}
	if months < 6 {
		return fmt.Sprintf("Analisa belum ideal karena IPO saham %s belum sampai 6 bulan.", ticker)
	}
	if months < 12 {
		return fmt.Sprintf("Ada kemungkinan analisa kurang akurat karena IPO saham %s belum sampai 12 bulan.", ticker)
	}
	return "-"
}

func commentaryRentangPeriode(totalDays int) string {
	if totalDays > 25 {
		return "Rentang periode analisa lebih dari 1 bulan, gunakan periode analisa yang lebih pendek untuk menghasilkan analisa yang lebih akurat dan komprehensif."
	}
	return "-"
}

func commentarySection2BigPlayer(t5 []T5Record) string {
	if len(t5) > 0 {
		recent := t5[0]
		act := ""
		if recent.Action != nil {
			act = strings.ToUpper(*recent.Action)
		}
		inv := "Pemegang Saham >5%"
		if recent.Investor != nil && *recent.Investor != "" {
			inv = *recent.Investor
		}
		absLot := int64(recent.Lot)
		if absLot < 0 {
			absLot = -absLot
		}
		lotStr := formatCommaInt(absLot)
		if act == "BUY" {
			return fmt.Sprintf("Terdeteksi penambahan kepemilikan oleh %s sebesar %s lot pada transaksi terakhir.", inv, lotStr)
		}
		if act == "SELL" {
			return fmt.Sprintf("Terdeteksi pengurangan kepemilikan oleh %s sebesar %s lot pada transaksi terakhir.", inv, lotStr)
		}
	}
	return "Kepemilikan Big Player stabil tanpa pergerakan aksi jual agresif dalam beberapa bulan terakhir."
}

func commentarySection3Partisipasi(retailMain, maxOther float64) string {
	if retailMain > maxOther {
		return "Transaksi broker Ritel terlalu tinggi, waspadai kemungkinan Bandar menggunakan beberapa broker Ritel."
	}
	return "-"
}

func commentarySection3Point5(retailMain, avgDailyValB float64) string {
	if retailMain < 3.0 && avgDailyValB > 15.0 {
		return "Investor Ritel tidak banyak bertransaksi atau mengenal saham ini."
	}
	return "-"
}

func commentaryAkumulasiBerturutTurut(sortedAll []DailyCandle, dailyResumes []sliceResume) string {
	if len(sortedAll) < 3 {
		return "-"
	}
	streak := 0
	start := maxInt(0, len(sortedAll)-7)
	for i := len(sortedAll) - 1; i >= start; i-- {
		cls := dailyResumes[i].Action.Cls
		if cls == "normal-acc" || cls == "big-acc" || cls == "massive-acc" {
			streak++
		} else {
			break
		}
	}
	if streak >= 3 {
		return fmt.Sprintf("Terdeteksi akumulasi berturut-turut oleh Bandar selama %d hari bursa terakhir.", streak)
	}
	return "-"
}

func commentaryLiveCandleRise(sortedAll []DailyCandle, priceChange float64) string {
	if len(sortedAll) >= 3 && priceChange > 3.0 {
		b3 := computeSliceResume(sortedAll[len(sortedAll)-3:])
		if b3.Action.Cls == "big-acc" || b3.Action.Cls == "massive-acc" {
			return "Kenaikan harga didukung oleh akumulasi Bandar yang solid dalam 3 hari bursa terakhir."
		}
	}
	return "-"
}

func commentaryVolumeRatio(vr float64) string {
	if vr > 10.0 {
		return fmt.Sprintf("Average Volume Ratio %.1f kali jauh di atas batas wajar, saham ini kemungkinan memiliki tujuan special dan akumulasi atau distribusi Bandar sulit untuk di deteksi.", vr)
	}
	if vr >= 8.0 && vr <= 10.0 {
		return "Volume Ratio saham ini di atas batas wajar, ada kemungkinan ini adalah saham REPO atau memiliki tujuan special lainnya."
	}
	return "-"
}

func commentaryRataRataTransaksi(avgDailyVal float64, totNetVol int64) string {
	avgM := avgDailyVal / 1e9
	if avgM < 1.0 {
		return "Rata-rata transaksi harian dibawah 1 M, pergerakan Bandar sangat minim dan masuknya Ritel dalam jumlah besar dapat langsung dideteksi oleh Bandar."
	}
	if avgM < 5.0 {
		return "Rata-rata transaksi harian rendah, pergerakan Bandar umumnya sulit dideteksi atau memberikan indikasi yang kurang akurat."
	}
	if totNetVol < 2000 {
		return "Net volume pada periode ini sangat kecil (di bawah 2.000 lot), pergerakan Bandar sulit dideteksi."
	}
	if avgM < 7.0 {
		return fmt.Sprintf("Rata-rata transaksi harian hanya sebesar %.1f M, berhati-hati jika ingin membeli saham ini dalam jumlah besar.", avgM)
	}
	return "-"
}

func commentaryRitelMinim(buyers, sellers []brokerNet) string {
	if len(buyers) == 0 && len(sellers) == 0 {
		return "-"
	}
	for i := 0; i < 7 && i < len(buyers); i++ {
		if RetailBrokers[buyers[i].Broker] {
			return "-"
		}
	}
	for i := 0; i < 7 && i < len(sellers); i++ {
		if RetailBrokers[sellers[i].Broker] {
			return "-"
		}
	}
	return "Pergerakan investor Ritel sangat minim pada periode ini."
}

func commentaryKoreksiMarkdownRetail(sortedAll []DailyCandle) string {
	if len(sortedAll) >= 5 {
		b5 := computeSliceResume(sortedAll[len(sortedAll)-5:])
		cls := b5.Action.Cls
		if b5.PriceChange < -15.0 && (cls == "normal-acc" || cls == "big-acc" || cls == "massive-acc" || cls == "undetected") {
			return "Terjadi penurunan harga signifikan namun disertai aksi bertahan atau akumulasi Bandar."
		}
	}
	return "-"
}

func commentaryLonjakanTransaksi(sortedAll []DailyCandle) string {
	n := len(sortedAll)
	if n >= 10 {
		var s5, s10 float64
		for _, c := range sortedAll[n-5:] {
			if c.Volume != nil {
				s5 += float64(*c.Volume)
			}
		}
		for _, c := range sortedAll[n-10 : n-5] {
			if c.Volume != nil {
				s10 += float64(*c.Volume)
			}
		}
		v5 := s5 / 5.0
		v10 := s10 / 5.0
		if v10 > 0 && (v5/v10) >= 2.0 {
			b5 := computeSliceResume(sortedAll[n-5:])
			switch b5.Action.Cls {
			case "big-acc", "massive-acc":
				return "Bandar semakin agresif mengakumulasi dalam 5 hari terakhir, volume transaksi meningkat drastis."
			case "undetected":
				return "Ada peningkatan volume transaksi yang cukup signifikan dalam 5 hari terakhir, namun belum terdeteksi aksi akumulasi maupun distribusi Bandar yang jelas."
			case "big-dist", "massive-dist":
				return "Terjadi kenaikan volume transaksi yang signifikan dalam 5 hari terakhir dan disertai distribusi Bandar."
			}
		}
	}
	return "-"
}

func commentaryRitelBeliPuncak(sortedAll []DailyCandle, dailyResumes []sliceResume) string {
	n := len(sortedAll)
	if n >= 10 {
		sl := sortedAll[n-10:]
		maxPrevHigh := math.Inf(-1)
		for i := 0; i < 9; i++ {
			h := highOrClose(sl[i])
			if h > maxPrevHigh {
				maxPrevHigh = h
			}
		}
		if highOrClose(sl[9]) > maxPrevHigh {
			cEnd := closeVal(sl[9])
			cStart := closeVal(sl[5])
			roc5 := 0.0
			if cStart > 0 {
				roc5 = ((cEnd - cStart) / cStart) * 100.0
			}
			rCls := dailyResumes[n-1].RetailAction.Cls
			if roc5 > 5.0 && (rCls == "normal-buy" || rCls == "aggressive-buy") {
				return "Investor Ritel terpantau melakukan pembelian agresif di area puncak harga, waspada potensi pembalikan arah."
			}
		}
	}
	return "-"
}

func commentarySidewaysAkumulasi(sortedAll []DailyCandle) string {
	for nDays := 4; nDays <= 10; nDays++ {
		if len(sortedAll) >= nDays {
			bn := computeSliceResume(sortedAll[len(sortedAll)-nDays:])
			if math.Abs(bn.PriceChange) <= 1.5 && (bn.Action.Cls == "big-acc" || bn.Action.Cls == "massive-acc") {
				return fmt.Sprintf("Harga bergerak sideways dalam %d hari terakhir sementara Bandar terdeteksi melakukan akumulasi tersembunyi.", nDays)
			}
		}
	}
	return "-"
}

func commentaryKinerjaGroup1Hari(gp *GroupPerformance) string {
	if gp == nil {
		return "-"
	}
	for _, g := range []*GroupPerfBlock{gp.Conglomerate, gp.Sector} {
		if g != nil && len(g.Members) >= 3 {
			var sum float64
			for _, m := range g.Members {
				sum += m.TodayPct
			}
			avgR := sum / float64(len(g.Members))
			if math.Abs(avgR) >= 2.0 {
				return fmt.Sprintf("Rata-rata kinerja saham se-%s %s mencatatkan pergerakan %+.1f%% hari ini.", g.Type, g.GroupName, avgR)
			}
		}
	}
	return "-"
}

func commentaryKinerjaGroup1Minggu(gp *GroupPerformance) string {
	if gp == nil {
		return "-"
	}
	for _, g := range []*GroupPerfBlock{gp.Conglomerate, gp.Sector} {
		if g == nil || len(g.Members) == 0 {
			continue
		}
		var sumW1 float64
		countW1 := 0
		for _, m := range g.Members {
			if m.W1Pct != nil {
				sumW1 += *m.W1Pct
				countW1++
			}
		}
		if countW1 >= 3 {
			avgR := sumW1 / float64(countW1)
			if math.Abs(avgR) >= 3.0 {
				return fmt.Sprintf("Dalam 1 minggu terakhir, rata-rata saham se-%s bergerak %+.1f%%.", g.Type, avgR)
			}
		} else if len(g.Members) >= 3 {
			var sumM1 float64
			for _, m := range g.Members {
				sumM1 += m.M1Pct
			}
			avgR := sumM1 / float64(len(g.Members))
			if math.Abs(avgR) >= 5.0 {
				return fmt.Sprintf("Dalam 1 bulan terakhir, rata-rata saham se-%s bergerak %+.1f%%.", g.Type, avgR)
			}
		}
	}
	return "-"
}

func commentaryBrokerBaru(sortedAll []DailyCandle) string {
	n := len(sortedAll)
	if n >= 21 {
		lastC := sortedAll[n-1]
		hist := sortedAll[n-21 : n-1]
		histBrokers := make(map[string]bool, 64)
		for _, c := range hist {
			for b := range c.BuyBrokers {
				histBrokers[b] = true
			}
			for b := range c.SellBrokers {
				histBrokers[b] = true
			}
		}
		for _, b := range sortedJSONBKeys(lastC.BuyBrokers) {
			if !RetailBrokers[b] && !histBrokers[b] {
				return fmt.Sprintf("Terdeteksi kemunculan broker baru %s di sisi beli yang sebelumnya tidak aktif dalam 20 hari terakhir.", b)
			}
		}
	}
	return "-"
}

func commentaryMarkupPeriode(priceChange float64, actionCls string, totalDays int) string {
	if 2 <= totalDays && totalDays <= 5 && priceChange > 4.0 && (actionCls == "big-acc" || actionCls == "massive-acc") {
		return "Terjadi markup harga terencana oleh Bandar dalam rentang periode ini."
	}
	return "-"
}

func commentaryPeriodeHarga(priceChange float64, actionCls, actionLabel string, totalDays int) string {
	pctStr := fmt.Sprintf("%.1f%%", math.Abs(priceChange))
	if totalDays > 5 && priceChange > 3.0 && (actionCls == "normal-dist" || actionCls == "big-dist" || actionCls == "massive-dist") {
		return fmt.Sprintf("Harga naik %s pada periode ini, namun Result periode menunjukkan %s. Kenaikan harga disertai distribusi Bandar, waspada potensi pelemahan.", pctStr, actionLabel)
	}
	if priceChange < -3.0 && (actionCls == "normal-acc" || actionCls == "big-acc" || actionCls == "massive-acc") {
		return fmt.Sprintf("Harga turun %s pada periode ini, namun Result periode menunjukkan %s. Penurunan harga disertai akumulasi Bandar.", pctStr, actionLabel)
	}
	return "-"
}

func commentaryPeriodeTransaksi(sortedAll []DailyCandle) string {
	n := len(sortedAll)
	if n >= 10 {
		mid := n / 2
		var vol1, vol2 int64
		for _, c := range sortedAll[:mid] {
			if c.Volume != nil {
				vol1 += *c.Volume
			}
		}
		for _, c := range sortedAll[mid:] {
			if c.Volume != nil {
				vol2 += *c.Volume
			}
		}
		if vol1 > 0 && float64(vol2)/float64(vol1) >= 1.5 {
			return "Aktivitas transaksi harian menunjukkan peningkatan likuiditas yang stabil."
		}
	}
	return "-"
}

func commentaryAverageAkumulasiBandar(totNetVol int64, totalDays int, lastClose float64) string {
	if totalDays > 0 && totNetVol > 0 {
		avgLot := float64(totNetVol) / float64(totalDays)
		valJt := float64(float64(avgLot*100.0)*lastClose) / 1e6
		if valJt >= 50.0 {
			return fmt.Sprintf("Rata-rata akumulasi Bandar pada periode ini berada pada kisaran Rp %s Juta per hari.", formatCommaFloat(valJt, 1))
		}
	}
	return "-"
}

func commentaryKorelasi(korelasiRaw float64) string {
	if korelasiRaw >= 5.0 {
		return "Korelasi pergerakan harga dengan akumulasi Bandar tergolong sangat ideal dan konsisten (Ideal Alignment)."
	}
	if korelasiRaw <= -3.0 {
		return "Pergerakan harga dan akumulasi Bandar menunjukkan divergensi negatif (anomali arah)."
	}
	return "-"
}

func commentaryKorelasiParuh(sortedAll []DailyCandle, dailyResumes []sliceResume) string {
	n := len(sortedAll)
	if n >= 20 {
		dayPt := func(idx int) float64 {
			d := sortedAll[idx]
			pCurr := closeVal(d)
			var pPrev float64
			if idx > 0 {
				pPrev = closeVal(sortedAll[idx-1])
			} else {
				pPrev = candlePrevOpenClose(d)
			}
			pctDay := 0.0
			if pPrev > 0 {
				pctDay = ((pCurr - pPrev) / pPrev) * 100.0
			}
			return AIScoreCorrDayPoint(dailyResumes[idx].Action.Cls, pctDay)
		}
		var sum1, sum2 float64
		for idx := n - 20; idx < n-10; idx++ {
			sum1 += dayPt(idx)
		}
		for idx := n - 10; idx < n; idx++ {
			sum2 += dayPt(idx)
		}
		s1 := sum1 / 5.0
		s2 := sum2 / 5.0
		if s1 >= 6.5 && s2 <= 3.75 {
			return "Korelasi pergerakan Bandar dengan harga mulai memburuk dalam 10 hari terakhir — sebelumnya sangat searah, sekarang tidak lagi, waspada perubahan arah."
		}
		if s1 <= 1.75 && s2 >= 3.75 {
			return "Korelasi pergerakan Bandar dengan harga mulai menguat dalam 10 hari terakhir — sebelumnya tidak searah, sekarang mulai searah dengan aksi akumulasi/distribusi Bandar."
		}
	}
	return "-"
}

func commentaryVolatilitas(sortedAll []DailyCandle) string {
	n := len(sortedAll)
	if n >= 10 {
		var sumRange float64
		count := 0
		for _, c := range sortedAll[n-10:] {
			h := highOrClose(c)
			l := lowOrClose(c)
			cl := 1.0
			if c.Close != nil && *c.Close != 0 {
				cl = *c.Close
			}
			if cl > 0 {
				sumRange += ((h - l) / cl) * 100.0
				count++
			}
		}
		if count > 0 && (sumRange/float64(count)) >= 7.5 {
			return "Volatilitas harga saham ini sangat tinggi dalam 10 hari terakhir, pergerakan harga intraday fluktuatif."
		}
	}
	return "-"
}

func commentaryConclusion(actionLabel, retailLabel string, aiScore int) string {
	zone := "Zona Netral / Konsolidasi"
	if aiScore >= 70 {
		zone = "Zona Akumulasi Kuat"
	} else if aiScore <= 40 {
		zone = "Zona Distribusi / Waspada"
	}
	return fmt.Sprintf("Kesimpulan: Status Bandar berada pada posisi %s didampingi aksi ritel %s. AI Score mencerminkan %s (%d/100).", actionLabel, retailLabel, zone, aiScore)
}

func highOrClose(c DailyCandle) float64 {
	if c.High != nil && *c.High != 0 {
		return *c.High
	}
	return closeVal(c)
}

func lowOrClose(c DailyCandle) float64 {
	if c.Low != nil && *c.Low != 0 {
		return *c.Low
	}
	return closeVal(c)
}
