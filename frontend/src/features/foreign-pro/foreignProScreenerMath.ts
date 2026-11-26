import type { ScreenerRow, ScreenerSubTab } from "./foreignProTypes";

export function getUniverseArray(row: ScreenerRow, universe: string): number[] {
  if (universe === "fco") return row.dailyFCoValue || [];
  if (universe === "kf") return row.dailyKFValue || [];
  if (universe === "xl") return row.dailyXLValue || [];
  if (universe === "semiritel") return row.dailySemiRitelValue || [];
  return row.dailyFFValue || [];
}

export function consecutiveInflow(arr: number[], n: number): boolean {
  if (!arr || arr.length < n) return false;
  const tail = arr.slice(-n);
  return tail.every((v) => typeof v === "number" && v > 0);
}

export function inflowCountAtLeast(arr: number[], windowN: number, minCount: number): boolean {
  if (!arr || arr.length < windowN) return false;
  const tail = arr.slice(-windowN);
  const count = tail.filter((v) => typeof v === "number" && v > 0).length;
  return count >= minCount;
}

export function firstInflowInWindow(arr: number[], windowN: number): boolean {
  if (!arr || arr.length < windowN) return false;
  const tail = arr.slice(-windowN);
  const today = tail[tail.length - 1];
  if (!(typeof today === "number" && today > 0)) return false;
  const rest = tail.slice(0, -1);
  return rest.every((v) => !(typeof v === "number" && v > 0));
}

export function isTodayMaxInWindow(arr: number[], n: number): boolean {
  if (!arr || arr.length < n) return false;
  const tail = arr.slice(-n);
  const today = tail[tail.length - 1];
  const valid = tail.filter((v) => typeof v === "number");
  if (typeof today !== "number" || valid.length === 0) return false;
  return today >= Math.max(...valid);
}

export function isLastRollingSumMax(arr: number[], sumWindow: number, totalWindow: number): boolean {
  if (!arr || arr.length < totalWindow) return false;
  const tail = arr.slice(-totalWindow);
  const sums: Array<number | null> = [];
  for (let i = 0; i <= tail.length - sumWindow; i++) {
    const seg = tail.slice(i, i + sumWindow);
    if (seg.some((v) => typeof v !== "number")) {
      sums.push(null);
      continue;
    }
    sums.push(seg.reduce((a, b) => a + b, 0));
  }
  const validSums = sums.filter((s): s is number => typeof s === "number");
  const lastSum = sums[sums.length - 1];
  if (typeof lastSum !== "number" || validSums.length === 0) return false;
  return lastSum >= Math.max(...validSums);
}

export function isLastPeriodSumPositiveAndMax(
  arr: number[],
  periodLen: number,
  numPeriods: number
): boolean {
  const totalWindow = periodLen * numPeriods;
  if (!arr || arr.length < totalWindow) return false;
  const tail = arr.slice(-totalWindow);
  const periods: Array<number | null> = [];
  for (let p = 0; p < numPeriods; p++) {
    const seg = tail.slice(p * periodLen, (p + 1) * periodLen);
    if (seg.some((v) => typeof v !== "number")) {
      periods.push(null);
      continue;
    }
    periods.push(seg.reduce((a, b) => a + b, 0));
  }
  const lastPeriod = periods[periods.length - 1];
  if (typeof lastPeriod !== "number" || lastPeriod <= 0) return false;
  const others = periods.slice(0, -1);
  if (others.some((v) => typeof v !== "number")) return false;
  return lastPeriod > Math.max(...(others as number[]));
}

function avgArr(arr: number[]): number | null {
  if (!arr.length) return null;
  return arr.reduce((a, b) => a + b, 0) / arr.length;
}

export function isBearishMarkUp(row: ScreenerRow, universe: string): boolean {
  const closeArr = row.dailyClose || [];
  const flowArr = getUniverseArray(row, universe);
  const valueArr = row.dailyValue || [];
  const n = closeArr.length;
  if (n < 21 || flowArr.length < 21 || valueArr.length < 21) return false;

  const close2dAgo = closeArr[n - 3];
  const ma20 = avgArr(closeArr.slice(-20));
  const ma5 = avgArr(closeArr.slice(-5));
  if (typeof close2dAgo !== "number" || ma20 === null || ma5 === null) return false;
  if (!(close2dAgo < ma20 && close2dAgo < ma5)) return false;

  const accum20 = flowArr.slice(-20).reduce((a, b) => a + b, 0);
  const accum5 = flowArr.slice(-5).reduce((a, b) => a + b, 0);
  if (!(accum20 < 0 && accum5 < 0)) return false;

  const todayChg = row.changePct;
  if (!(typeof todayChg === "number" && todayChg > 0)) return false;

  const todayFlow = flowArr[n - 1];
  const todayValue = valueArr[n - 1];
  if (typeof todayFlow !== "number" || typeof todayValue !== "number" || todayValue === 0) {
    return false;
  }
  const todayRatio = todayFlow / todayValue;

  const flowTail = flowArr.slice(-20);
  const valueTail = valueArr.slice(-20);
  const accumDayRatios: number[] = [];
  for (let i = 0; i < flowTail.length; i++) {
    const f = flowTail[i];
    const v = valueTail[i];
    if (typeof f === "number" && typeof v === "number" && f > 0 && v !== 0) {
      accumDayRatios.push(f / v);
    }
  }
  if (accumDayRatios.length === 0) return false;
  const avgAccumRatio = avgArr(accumDayRatios);
  return avgAccumRatio !== null && todayRatio < avgAccumRatio;
}

export function isBullishMarkDown(row: ScreenerRow, universe: string): boolean {
  const closeArr = row.dailyClose || [];
  const flowArr = getUniverseArray(row, universe);
  const valueArr = row.dailyValue || [];
  const n = closeArr.length;
  if (n < 21 || flowArr.length < 21 || valueArr.length < 21) return false;

  const close2dAgo = closeArr[n - 3];
  const ma20 = avgArr(closeArr.slice(-20));
  const ma5 = avgArr(closeArr.slice(-5));
  if (typeof close2dAgo !== "number" || ma20 === null || ma5 === null) return false;
  if (!(close2dAgo > ma20 && close2dAgo > ma5)) return false;

  const accum20 = flowArr.slice(-20).reduce((a, b) => a + b, 0);
  const accum5 = flowArr.slice(-5).reduce((a, b) => a + b, 0);
  if (!(accum20 > 0 && accum5 > 0)) return false;

  const todayChg = row.changePct;
  if (!(typeof todayChg === "number" && todayChg < 0)) return false;

  const todayFlow = flowArr[n - 1];
  const todayValue = valueArr[n - 1];
  if (typeof todayFlow !== "number" || typeof todayValue !== "number" || todayValue === 0) {
    return false;
  }
  const todayRatio = todayFlow / todayValue;

  const flowTail = flowArr.slice(-20);
  const valueTail = valueArr.slice(-20);
  const distDayRatios: number[] = [];
  for (let i = 0; i < flowTail.length; i++) {
    const f = flowTail[i];
    const v = valueTail[i];
    if (typeof f === "number" && typeof v === "number" && f < 0 && v !== 0) {
      distDayRatios.push(f / v);
    }
  }
  if (distDayRatios.length === 0) return false;
  const avgDistRatio = avgArr(distDayRatios);
  return avgDistRatio !== null && todayRatio > avgDistRatio;
}

export function filterScreenerRows(params: {
  screenerData: ScreenerRow[];
  screenerSearch: string;
  screenerSubTab: ScreenerSubTab;
  activeStrategyId: string;
  universe: string;
  hfpOnly: boolean;
  mfpOnly: boolean;
  stateOnly: boolean;
  priceClass: string;
}): ScreenerRow[] {
  const {
    screenerData,
    screenerSearch,
    screenerSubTab,
    activeStrategyId,
    universe,
    hfpOnly,
    mfpOnly,
    stateOnly,
    priceClass,
  } = params;

  let list = [...screenerData];

  if (screenerSearch) {
    const q = screenerSearch.toUpperCase();
    list = list.filter((s) => s.ticker.includes(q) || s.name.toUpperCase().includes(q));
  }

  if (screenerSubTab === "basic") {
    if (activeStrategyId === "basicHFP") {
      list = list.filter((s) => s.fpBadge === "HFP");
    } else if (activeStrategyId === "basicInflow50") {
      list = list.filter((s) => {
        const arr = s.dailyFFValue || [];
        return arr.length > 0 && arr[arr.length - 1] > 50.0;
      });
    }
  } else if (screenerSubTab === "custom") {
    if (activeStrategyId === "inflow3m") {
      list = list.filter((s) => {
        const arr = getUniverseArray(s, universe);
        return (
          arr.length >= 21 &&
          isTodayMaxInWindow(arr, 21) &&
          (s.fpBadge === "HFP" || s.fpBadge === "MFP")
        );
      });
    } else if (activeStrategyId === "consec3") {
      list = list.filter((s) => consecutiveInflow(getUniverseArray(s, universe), 3));
    } else if (activeStrategyId === "consec5") {
      list = list.filter((s) => consecutiveInflow(getUniverseArray(s, universe), 5));
    } else if (activeStrategyId === "inflow3of5") {
      list = list.filter((s) => inflowCountAtLeast(getUniverseArray(s, universe), 5, 3));
    } else if (activeStrategyId === "firstInflow10") {
      list = list.filter((s) => firstInflowInWindow(getUniverseArray(s, universe), 10));
    } else if (activeStrategyId === "maxInflow5min3") {
      list = list.filter(
        (s) =>
          isTodayMaxInWindow(getUniverseArray(s, universe), 5) &&
          inflowCountAtLeast(getUniverseArray(s, universe), 5, 3)
      );
    } else if (activeStrategyId === "top2day20") {
      list = list.filter((s) => isLastRollingSumMax(getUniverseArray(s, universe), 2, 20));
    } else if (activeStrategyId === "weeklyInflowMax") {
      list = list.filter((s) =>
        isLastPeriodSumPositiveAndMax(getUniverseArray(s, universe), 5, 4)
      );
    } else if (activeStrategyId === "bearishMarkUp") {
      list = list.filter((s) => isBearishMarkUp(s, universe));
    } else if (activeStrategyId === "bullishMarkDown") {
      list = list.filter((s) => isBullishMarkDown(s, universe));
    }
  } else if (screenerSubTab === "manual") {
    if (hfpOnly) list = list.filter((s) => s.fpBadge === "HFP");
    if (mfpOnly) list = list.filter((s) => s.fpBadge === "MFP");
    if (stateOnly) {
      if (universe === "ff" || universe === "fco") {
        list = list.filter((s) => s.accumVal > 0);
      } else {
        list = list.filter((s) => s.accumVal < 0);
      }
    }
    if (priceClass === "Up") list = list.filter((s) => s.changePct > 0);
    else if (priceClass === "Down") list = list.filter((s) => s.changePct < 0);
    else if (priceClass === "Sideways") list = list.filter((s) => s.changePct === 0);
  }

  return list;
}
