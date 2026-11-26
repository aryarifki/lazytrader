import { getUniverseArray } from "./foreignProScreenerMath";
import type { MarketReportData, PeriodReportData, ScreenerRow } from "./foreignProTypes";

interface ForeignProScreenerResultsProps {
  activeReportType: string | null;
  marketReportData: MarketReportData;
  periodReportWindow: number | null;
  periodReportData: PeriodReportData;
  filteredScreenerList: ScreenerRow[];
  universe: string;
  windowLookback: number;
  accumHeaderLabel: string;
  isLightMode: boolean;
  onSelectStock: (ticker: string) => void;
}

function SparklineSvg({
  arr,
  universe,
  isLightMode,
}: {
  arr: number[];
  universe: string;
  isLightMode: boolean;
}) {
  const vals = (arr || []).slice(-10);
  if (!vals.length) return null;
  const w = 90;
  const h = 24;
  const gap = 1.5;
  const barW = w / vals.length - gap;
  const maxAbs = Math.max(...vals.map((v) => Math.abs(v)), 0.0001);
  const mid = h / 2;
  const midStroke = isLightMode ? "rgba(30, 26, 23, 0.18)" : "rgba(255,255,255,0.15)";

  return (
    <svg className="inline-block" width={w} height={h} viewBox={`0 0 ${w} ${h}`}>
      {vals.map((v, i) => {
        const bh = Math.max((Math.abs(v) / maxAbs) * (h / 2 - 1), 0.5);
        const x = i * (barW + gap);
        const y = v >= 0 ? mid - bh : mid;
        const color =
          universe === "fco"
            ? "#EE2E31"
            : universe === "kf"
              ? "#6EC6FF"
              : universe === "xl"
                ? "#A78BFA"
                : universe === "semiritel"
                  ? "#F0B429"
                  : v >= 0
                    ? "#FF6820"
                    : "#E5483D";
        return (
          <rect
            key={i}
            x={x.toFixed(1)}
            y={y.toFixed(1)}
            width={barW.toFixed(1)}
            height={bh.toFixed(1)}
            fill={color}
            opacity={0.85}
          />
        );
      })}
      <line x1={0} y1={mid} x2={w} y2={mid} stroke={midStroke} strokeWidth={1} />
    </svg>
  );
}

function fpBadgeClass(fpBadge: string): string {
  if (fpBadge === "HFP") {
    return "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border border-[var(--md-sys-color-primary)]/40";
  }
  if (fpBadge === "MFP") {
    return "bg-teal-500/15 text-teal-700 dark:text-[#2FD4C6] border border-teal-500/40";
  }
  return "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)]/50";
}

export default function ForeignProScreenerResults({
  activeReportType,
  marketReportData,
  periodReportWindow,
  periodReportData,
  filteredScreenerList,
  universe,
  windowLookback,
  accumHeaderLabel,
  isLightMode,
  onSelectStock,
}: ForeignProScreenerResultsProps) {
  return (
    <>
      {activeReportType === "topBuySell" && (
        <div id="reportPanel" className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-3 rounded-xl space-y-2">
            <div className="text-xs font-bold text-[var(--color-positive)] border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1">
              TOP FOREIGN BUY (HARI INI)
            </div>
            <div id="reportBuyList" className="space-y-1.5 text-xs">
              {marketReportData.topForeignBuy.map((r) => (
                <div key={r.ticker} className="flex justify-between items-center">
                  <span className="font-bold text-[var(--md-sys-color-on-surface)]">{r.ticker}</span>
                  <span className="text-[var(--color-positive)] font-bold">+{r.val.toFixed(1)} M</span>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-3 rounded-xl space-y-2">
            <div className="text-xs font-bold text-[var(--color-negative)] border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1">
              TOP FOREIGN SELL (HARI INI)
            </div>
            <div id="reportSellList" className="space-y-1.5 text-xs">
              {marketReportData.topForeignSell.map((r) => (
                <div key={r.ticker} className="flex justify-between items-center">
                  <span className="font-bold text-[var(--md-sys-color-on-surface)]">{r.ticker}</span>
                  <span className="text-[var(--color-negative)] font-bold">{r.val.toFixed(1)} M</span>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-3 rounded-xl space-y-2">
            <div className="text-xs font-bold text-teal-700 dark:text-[#2FD4C6] border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1">
              TOP GAINER (&gt; 2.5 M)
            </div>
            <div id="reportGainerList" className="space-y-1.5 text-xs">
              {marketReportData.topGainer.map((r) => (
                <div key={r.ticker} className="flex justify-between items-center">
                  <span className="font-bold text-[var(--md-sys-color-on-surface)]">{r.ticker}</span>
                  <span className="text-teal-700 dark:text-[#2FD4C6] font-bold">+{r.val.toFixed(1)}%</span>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-3 rounded-xl space-y-2">
            <div className="text-xs font-bold text-[var(--color-negative)] border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1">
              TOP LOSER (&gt; 2.5 M)
            </div>
            <div id="reportLoserList" className="space-y-1.5 text-xs">
              {marketReportData.topLoser.map((r) => (
                <div key={r.ticker} className="flex justify-between items-center">
                  <span className="font-bold text-[var(--md-sys-color-on-surface)]">{r.ticker}</span>
                  <span className="text-[var(--color-negative)] font-bold">{r.val.toFixed(1)}%</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {activeReportType === "periodReport" && (
        <div id="periodReportPanel" className="grid grid-cols-1 md:grid-cols-2 gap-3">
          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-3 rounded-xl space-y-2">
            <div
              id="periodBuyTitle"
              className="text-xs font-bold text-[var(--color-positive)] border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1"
            >
              TOP FOREIGN BUY ({periodReportWindow} HARI)
            </div>
            <div id="periodBuyList" className="space-y-1.5 text-xs">
              {periodReportData.topBuy.map((r) => (
                <div key={r.ticker} className="flex justify-between items-center">
                  <span className="font-bold text-[var(--md-sys-color-on-surface)]">{r.ticker}</span>
                  <span className="text-[var(--color-positive)] font-bold">+{r.val.toFixed(1)} M</span>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-3 rounded-xl space-y-2">
            <div
              id="periodSellTitle"
              className="text-xs font-bold text-[var(--color-negative)] border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1"
            >
              TOP FOREIGN SELL ({periodReportWindow} HARI)
            </div>
            <div id="periodSellList" className="space-y-1.5 text-xs">
              {periodReportData.topSell.map((r) => (
                <div key={r.ticker} className="flex justify-between items-center">
                  <span className="font-bold text-[var(--md-sys-color-on-surface)]">{r.ticker}</span>
                  <span className="text-[var(--color-negative)] font-bold">{r.val.toFixed(1)} M</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* SCREENER RESULTS TABLE */}
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-xl overflow-hidden shadow-md">
        <div className="bg-[var(--md-sys-color-surface-container-high)] px-3 py-2 border-b border-[var(--md-sys-color-outline-variant)]/60 flex items-center justify-between">
          <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] tracking-wide">
            HASIL SCREENER — {universe.toUpperCase()} ({windowLookback} HARI)
          </span>
          <span id="resultCount" className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
            {filteredScreenerList.length} Saham Terfilter
          </span>
        </div>

        {/* Mobile Card List */}
        <div className="md:hidden divide-y divide-[var(--md-sys-color-outline-variant)]/40 p-2 space-y-2">
          {filteredScreenerList.slice(0, 60).map((s) => (
            <div
              key={s.ticker}
              onClick={() => onSelectStock(s.ticker)}
              className="bg-[var(--md-sys-color-surface-container-low)] p-3 rounded-xl border border-[var(--md-sys-color-outline-variant)]/50 active:bg-[var(--md-sys-color-surface-container-highest)] cursor-pointer transition space-y-2 font-mono-nums"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className="font-extrabold text-base text-[var(--md-sys-color-primary)]">
                    {s.ticker}
                  </span>
                  <span className={`px-2 py-0.5 rounded text-[10px] font-bold ${fpBadgeClass(s.fpBadge)}`}>
                    {s.fpBadge}
                  </span>
                </div>
                <div className="text-right">
                  <span className="font-bold text-xs text-[var(--md-sys-color-on-surface)]">
                    Rp {s.close?.toLocaleString("id-ID")}
                  </span>
                  <span
                    className={`block text-[11px] font-bold ${
                      s.changePct >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                    }`}
                  >
                    {s.changePct > 0 ? `+${s.changePct}%` : `${s.changePct}%`}
                  </span>
                </div>
              </div>

              <div className="flex items-center justify-between text-xs text-[var(--md-sys-color-on-surface-variant)] pt-1 border-t border-[var(--md-sys-color-outline-variant)]/30">
                <span className="truncate max-w-[170px] text-[11px]">{s.name}</span>
                <span
                  className={`font-bold ${
                    s.accumVal >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                  }`}
                >
                  Flow: {s.accumVal > 0 ? `+${s.accumVal.toFixed(1)}` : s.accumVal.toFixed(1)} M
                </span>
              </div>
            </div>
          ))}
        </div>

        {/* Desktop Table */}
        <div className="hidden md:block overflow-x-auto">
          <table className="w-full text-left text-xs font-mono">
            <thead>
              <tr className="bg-[var(--md-sys-color-surface-container-low)] text-[10px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/40">
                <th className="py-2 px-3">KODE</th>
                <th className="py-2 px-3">NAMA EMITEN</th>
                <th className="py-2 px-2 text-right">HARGA</th>
                <th className="py-2 px-2 text-right">CHG %</th>
                <th className="py-2 px-2 text-center">FP</th>
                <th id="accumHeader" className="py-2 px-3 text-right">
                  {accumHeaderLabel}
                </th>
                <th className="py-2 px-3 text-right">AVG TRANS</th>
                <th className="py-2 px-3 text-center">TREND (10D)</th>
                <th className="py-2 px-3 text-center">AKSI</th>
              </tr>
            </thead>
            <tbody id="tbody" className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
              {filteredScreenerList.slice(0, 100).map((s) => (
                <tr key={s.ticker} className="hover:bg-[var(--md-sys-color-surface-container-high)]/50 transition">
                  <td className="py-2 px-3 font-bold text-[var(--md-sys-color-on-surface)]">
                    <div className="flex items-center gap-1.5">
                      <span>{s.ticker}</span>
                      {s.badges?.map((b) => (
                        <span
                          key={b}
                          className="text-[9px] bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] px-1 py-0.5 rounded font-sans font-bold"
                        >
                          {b}
                        </span>
                      ))}
                    </div>
                  </td>
                  <td className="py-2 px-3 text-[var(--md-sys-color-on-surface-variant)] truncate max-w-[200px]">
                    {s.name}
                  </td>
                  <td className="py-2 px-2 text-right font-bold text-[var(--md-sys-color-on-surface)]">
                    Rp {s.close?.toLocaleString("id-ID")}
                  </td>
                  <td
                    className={`py-2 px-2 text-right font-bold ${
                      s.changePct >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                    }`}
                  >
                    {s.changePct > 0 ? `+${s.changePct}%` : `${s.changePct}%`}
                  </td>
                  <td className="py-2 px-2 text-center">
                    <span className={`px-1.5 py-0.5 rounded text-[10px] font-bold ${fpBadgeClass(s.fpBadge)}`}>
                      {s.fpBadge}
                    </span>
                  </td>
                  <td
                    className={`py-2 px-3 text-right font-bold ${
                      s.accumVal >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                    }`}
                  >
                    {s.accumVal > 0 ? `+${s.accumVal.toFixed(1)}` : s.accumVal.toFixed(1)} M
                  </td>
                  <td className="py-2 px-3 text-right text-[var(--md-sys-color-on-surface-variant)]">
                    Rp {s.avgTransValue10d.toFixed(1)} M
                  </td>
                  <td className="py-2 px-3 text-center">
                    <SparklineSvg
                      arr={getUniverseArray(s, universe)}
                      universe={universe}
                      isLightMode={isLightMode}
                    />
                  </td>
                  <td className="py-2 px-3 text-center">
                    <button
                      type="button"
                      onClick={() => onSelectStock(s.ticker)}
                      className="bg-[var(--md-sys-color-primary)] hover:opacity-90 text-[var(--md-sys-color-on-primary)] px-2.5 py-1 rounded text-[10px] font-bold font-sans transition cursor-pointer"
                    >
                      Chart
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </>
  );
}
