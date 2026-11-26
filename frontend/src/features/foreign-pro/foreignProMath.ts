import type {
  AkumulasiResult,
  IhsgMacroMetrics,
  IhsgResponse,
  StockSeriesResponse,
} from "./foreignProTypes";

export function roundNum(n: number, d: number): number {
  const factor = Math.pow(10, d);
  return Math.round(n * factor) / factor;
}

export function calcWindowCumulative(
  dailyArr: number[] | undefined,
  start: number,
  end: number
): number[] {
  if (!dailyArr || !dailyArr.length) return [];
  const windowSlice = dailyArr.slice(start, end);
  let running = 0;
  return windowSlice.map((v) => {
    running += v || 0;
    return Math.round(running * 100) / 100;
  });
}

export function calcCumLnrSlice(
  base: number[] | undefined,
  start: number,
  end: number
): number[] {
  if (!base || !base.length) return [];
  const rawSlice = base.slice(start, end);
  const startVal = rawSlice[0] || 0;
  return rawSlice.map((v) => Math.round((v - startVal) * 100) / 100);
}

export function calculateSMA(series: number[], period: number): Array<number | null> {
  const res: Array<number | null> = [];
  for (let i = 0; i < series.length; i++) {
    if (i + 1 < period) {
      res.push(null);
    } else {
      const seg = series.slice(i + 1 - period, i + 1);
      res.push(seg.reduce((a, b) => a + b, 0) / period);
    }
  }
  return res;
}

export function calculateAkumulasi(
  stockData: StockSeriesResponse | null,
  akumDays: number
): AkumulasiResult {
  const an = stockData?.analytics;
  if (!an?.ff_daily_val || !an.closes || !an.opens || !an.highs || !an.lows) {
    return { avgPrice: null, totalFlow: 0, isAcc: true };
  }
  const n = an.closes.length;
  const w = Math.min(akumDays, n);
  if (w < 1) return { avgPrice: null, totalFlow: 0, isAcc: true };

  const sliceCloses = an.closes.slice(-w);
  const sliceOpens = an.opens.slice(-w);
  const sliceHighs = an.highs.slice(-w);
  const sliceLows = an.lows.slice(-w);
  const sliceFF = an.ff_daily_val.slice(-w);

  let sumShares = 0;
  let sumMoney = 0;
  let totalFlow = 0;

  for (let i = 0; i < w; i++) {
    const avgP = (sliceOpens[i] + sliceHighs[i] + sliceLows[i] + sliceCloses[i]) / 4.0;
    const ffM = sliceFF[i] || 0;
    totalFlow += ffM;
    const shares = avgP > 0 ? (ffM * 1e9) / avgP : 0;
    sumShares += shares;
    sumMoney += ffM * 1e9;
  }

  const isAcc = totalFlow >= 0;
  const avgPrice =
    Math.abs(sumShares) > 0 ? Math.round(Math.abs(sumMoney) / Math.abs(sumShares)) : null;
  return { avgPrice, totalFlow, isAcc };
}

export function calcIhsgMacroMetrics(
  ihsgData: IhsgResponse | null,
  flipKF: boolean,
  flipXL: boolean,
  flipSR: boolean
): IhsgMacroMetrics {
  if (!ihsgData?.points || !ihsgData.points.length) {
    return { ff: "+6.0 T", fco: "+4.2 T", kf: "-2.8 T", xl: "-1.1 T", sr: "+0.8 T" };
  }
  const windowPts = ihsgData.points.slice(-120);
  const sumKey = (met: string, isFlip = false) => {
    let s = 0;
    windowPts.forEach((p) => {
      const val = p.ihsg?.[met] || 0;
      s += isFlip ? -val : val;
    });
    const tVal = s / 1000.0;
    return `${tVal >= 0 ? "+" : ""}${tVal.toFixed(1)} T`;
  };
  return {
    ff: sumKey("ff"),
    fco: sumKey("fco"),
    kf: sumKey("kf", flipKF),
    xl: sumKey("xl", flipXL),
    sr: sumKey("sr", flipSR),
  };
}
