import type {
  AdminReportMode,
  IhsgResponse,
  MarketReportData,
  PeriodReportData,
  ScreenerRow,
} from "./foreignProTypes";

export function computeMarketReportData(screenerData: ScreenerRow[]): MarketReportData {
  const withTodayFF = screenerData
    .map((r) => {
      const arr = r.dailyFFValue || [];
      const todayVal = arr.length > 0 ? arr[arr.length - 1] : 0.0;
      return { ticker: r.ticker, val: todayVal, changePct: r.changePct };
    })
    .filter((r) => typeof r.val === "number");

  const topForeignBuy = [...withTodayFF].sort((a, b) => b.val - a.val).slice(0, 5);
  const topForeignSell = [...withTodayFF].sort((a, b) => a.val - b.val).slice(0, 5);

  const liquidToday = screenerData.filter((r) => {
    const arr = r.dailyValue || [];
    const v = arr.length > 0 ? arr[arr.length - 1] : 0.0;
    return v > 2.5 && typeof r.changePct === "number";
  });

  const topGainer = [...liquidToday]
    .sort((a, b) => b.changePct - a.changePct)
    .slice(0, 5)
    .map((r) => ({ ticker: r.ticker, val: r.changePct }));

  const topLoser = [...liquidToday]
    .sort((a, b) => a.changePct - b.changePct)
    .slice(0, 5)
    .map((r) => ({ ticker: r.ticker, val: r.changePct }));

  return { topForeignBuy, topForeignSell, topGainer, topLoser };
}

export function computePeriodReportData(
  screenerData: ScreenerRow[],
  periodReportWindow: number | null
): PeriodReportData {
  if (!periodReportWindow) return { topBuy: [], topSell: [] };
  const w = periodReportWindow;
  const mapped = screenerData.map((r) => {
    const accum = r.summary?.[`ffAccum${w}d`] || r.accumVal || 0.0;
    return { ticker: r.ticker, val: accum, changePct: r.changePct };
  });
  const topBuy = [...mapped].sort((a, b) => b.val - a.val).slice(0, 5);
  const topSell = [...mapped].sort((a, b) => a.val - b.val).slice(0, 5);
  return { topBuy, topSell };
}

export function maskTicker(tk: string): string {
  if (tk.length <= 2) return tk;
  const chars = tk.split("");
  chars[0] = "*";
  if (chars.length >= 3) chars[chars.length - 2] = "*";
  return chars.join("");
}

export function formatIndoDate(dStr?: string): string {
  if (!dStr) return "";
  try {
    const parts = dStr.split("-");
    if (parts.length === 3) {
      const year = parseInt(parts[0], 10);
      const month = parseInt(parts[1], 10) - 1;
      const day = parseInt(parts[2], 10);
      const dt = new Date(year, month, day);
      const days = ["Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"];
      const months = ["Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"];
      return `${days[dt.getDay()]}, ${day < 10 ? `0${day}` : day} ${months[month]} ${year}`;
    }
  } catch {
    // ignore malformed date
  }
  return dStr;
}

export function buildAdminReportText(params: {
  adminMode: AdminReportMode;
  ihsgData: IhsgResponse | null;
  marketReportData: MarketReportData;
  screenerData: ScreenerRow[];
}): string {
  const { adminMode, ihsgData, marketReportData, screenerData } = params;
  const ihsg = ihsgData?.summary;
  const dateFormatted = formatIndoDate(ihsg?.date) || "Senin, 02 Nov 2026";
  const move =
    (ihsg?.changePct || 0) >= 0
      ? `kenaikan +${Math.abs(ihsg?.changePct || 0).toFixed(2)}%`
      : `penurunan -${Math.abs(ihsg?.changePct || 0).toFixed(2)}%`;
  const ffVal = ihsg?.ffValue ?? -356.27;
  const driver =
    ffVal >= 0
      ? `pembelian Investor Asing sebesar ${Math.abs(ffVal).toFixed(0)} Miliar Rupiah`
      : `penjualan Investor Asing sebesar ${Math.abs(ffVal).toFixed(0)} Miliar Rupiah`;

  const foreignBuyYesterday: Array<{ ticker: string; chg: number; buyVal: number }> = [];
  const rebound3d: Array<{ ticker: string; chg: number }> = [];
  const volumeSpikes: Array<{ ticker: string; multiplier: number }> = [];

  screenerData.forEach((r) => {
    const tk = r.ticker;
    if (!tk || tk === "IHSG_Aggregate") return;
    const dailyFF = r.dailyFFValue || [];
    const dailyClose = r.dailyClose || [];
    const dailyVal = r.dailyValue || [];
    const chg = r.changePct || 0;
    const avg10d = r.avgTransValue10d || 0;

    const yesterdayFF = dailyFF.length >= 2 ? dailyFF[dailyFF.length - 2] : null;
    if (typeof yesterdayFF === "number" && yesterdayFF > 0.5 && chg > 0) {
      foreignBuyYesterday.push({ ticker: tk, chg, buyVal: yesterdayFF });
    }
    if (
      dailyClose.length >= 4 &&
      typeof dailyClose[dailyClose.length - 2] === "number" &&
      typeof dailyClose[dailyClose.length - 3] === "number" &&
      typeof dailyClose[dailyClose.length - 4] === "number" &&
      dailyClose[dailyClose.length - 2] < dailyClose[dailyClose.length - 3] &&
      dailyClose[dailyClose.length - 3] < dailyClose[dailyClose.length - 4] &&
      chg > 0
    ) {
      rebound3d.push({ ticker: tk, chg });
    }
    if (dailyVal.length >= 1 && avg10d > 1.0) {
      const todayVal = dailyVal[dailyVal.length - 1];
      if (typeof todayVal === "number") {
        const mult = todayVal / avg10d;
        if (mult >= 1.5) volumeSpikes.push({ ticker: tk, multiplier: mult });
      }
    }
  });

  foreignBuyYesterday.sort((a, b) => b.buyVal - a.buyVal);
  rebound3d.sort((a, b) => b.chg - a.chg);
  volumeSpikes.sort((a, b) => b.multiplier - a.multiplier);

  const lines: string[] = [];
  if (adminMode === "sesi1") {
    lines.push(`Market Report Sesi 1 by InvestOwl ID | ${dateFormatted} (12:00 WIB)`);
    lines.push("");
    lines.push(`IHSG mengalami ${move} dimotori oleh ${driver}.`);
    lines.push("");
    lines.push("Berikut report penutupan sesi pertama :");
    lines.push("");
    lines.push("Top Foreign Buy");
    marketReportData.topForeignBuy.forEach((r) => lines.push(`${maskTicker(r.ticker)}, ${r.val.toFixed(1)}M`));
    lines.push("");
    lines.push("Top Local Buy");
    marketReportData.topForeignSell.forEach((r) => lines.push(`${maskTicker(r.ticker)}, ${r.val.toFixed(1)}M`));
    lines.push("");
    lines.push("Top Gainer");
    marketReportData.topGainer.forEach((r) =>
      lines.push(`${r.ticker}, ${r.val > 0 ? `+${r.val.toFixed(1)}%` : `${r.val.toFixed(1)}%`}`)
    );
    lines.push("");
    lines.push("Top Loser");
    marketReportData.topLoser.forEach((r) => lines.push(`${r.ticker}, ${r.val.toFixed(1)}%`));

    if (foreignBuyYesterday.length > 0) {
      lines.push("");
      lines.push("Naik setelah Investor Asing Beli kemarin:");
      foreignBuyYesterday.slice(0, 3).forEach((item) => {
        lines.push(`${item.ticker} ${item.chg > 0 ? `+${item.chg.toFixed(1)}%` : `${item.chg.toFixed(1)}%`} / ${item.buyVal.toFixed(1)}M`);
      });
    }
    if (rebound3d.length > 0) {
      lines.push("");
      lines.push("Naik setelah turun 3 hari:");
      rebound3d.slice(0, 3).forEach((item) => {
        lines.push(`${item.ticker} ${item.chg > 0 ? `+${item.chg.toFixed(1)}%` : `${item.chg.toFixed(1)}%`}`);
      });
    }
    if (volumeSpikes.length > 0) {
      lines.push("");
      lines.push("Lonjakan volume:");
      volumeSpikes.slice(0, 3).forEach((item) => {
        lines.push(`${item.ticker}  ${item.multiplier.toFixed(1)}x`);
      });
    }
  } else {
    lines.push(`Market & Foreign Flow Report Penutupan Sesi Kedua by Investowl ID | ${dateFormatted} (16:00 WIB)`);
    lines.push("");
    lines.push(`IHSG mengalami ${move} dimotori oleh ${driver}.`);
    lines.push("");
    lines.push("Berikut report penutupan sesi kedua :");
    lines.push("");
    lines.push("Top Foreign Buy");
    marketReportData.topForeignBuy.forEach((r) =>
      lines.push(`${r.ticker}, ${r.val.toFixed(1)}M / ${(r.changePct ?? 0) >= 0 ? `+${r.changePct}%` : `${r.changePct}%`}`)
    );
    lines.push("");
    lines.push("Top Local Buy");
    marketReportData.topForeignSell.forEach((r) =>
      lines.push(`${r.ticker}, ${r.val.toFixed(1)}M / ${(r.changePct ?? 0) >= 0 ? `+${r.changePct}%` : `${r.changePct}%`}`)
    );
    lines.push("");
    lines.push("Top Gainer");
    marketReportData.topGainer.forEach((r) =>
      lines.push(`${r.ticker}, ${r.val > 0 ? `+${r.val.toFixed(1)}%` : `${r.val.toFixed(1)}%`}`)
    );
    lines.push("");
    lines.push("Top Loser");
    marketReportData.topLoser.forEach((r) => lines.push(`${r.ticker}, ${r.val.toFixed(1)}%`));
  }
  return lines.join("\n");
}

export function getAccumHeaderLabel(universe: string, windowLookback: number): string {
  if (universe === "fco") return "F & Co Acc / Dist";
  if (universe === "kf") return "Retail Buy / Sell";
  if (universe === "xl") return "XL Buy / Sell";
  if (universe === "semiritel") return "Semi Ritel Buy / Sell";
  return `Foreign Acc / Dist (${windowLookback}D)`;
}
