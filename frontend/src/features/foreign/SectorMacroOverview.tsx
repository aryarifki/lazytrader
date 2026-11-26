import { Icon } from "@iconify/react";
import { fmtRadarFlow, signedColor } from "@/lib/format";
import type { BenchmarkInfo } from "./sectorTypes";

interface SectorMacroOverviewProps {
  date: string;
  benchmark: BenchmarkInfo;
  macroMarketState?: BenchmarkInfo["market_state"];
  computedTotalForeign1d: number;
  selectedQuadrant: string;
  onSelectQuadrant: (q: string) => void;
  quadrantCounts: Record<string, number>;
  foreignStateCounts: {
    surge: number;
    warning: number;
    exodus: number;
    bottoming: number;
  };
  isRefreshing: boolean;
  onRefresh: () => void;
}

export default function SectorMacroOverview({
  date,
  benchmark,
  macroMarketState,
  computedTotalForeign1d,
  selectedQuadrant,
  onSelectQuadrant,
  quadrantCounts,
  foreignStateCounts,
  isRefreshing,
  onRefresh,
}: SectorMacroOverviewProps) {
  const effectiveMarketState = macroMarketState || benchmark.market_state;

  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 sm:p-6 shadow-xs space-y-4">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-3.5 border-b border-[var(--md-sys-color-outline-variant)]/60 gap-3">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 sm:w-11 sm:h-11 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0 border border-[var(--md-sys-color-outline-variant)] shadow-2xs">
            <Icon icon="ph:chart-polar-duotone" width="24" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h2 className="text-base sm:text-lg font-extrabold text-[var(--md-sys-color-on-surface)] tracking-tight">
                Peta Rotasi Sektor &amp; Foreign Exodus Radar
              </h2>
              <div className="flex items-center gap-1.5">
                <span className="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)] font-mono">
                  {date}
                </span>
                <button
                  onClick={onRefresh}
                  disabled={isRefreshing}
                  title="Sinkronkan data rotasi & arus asing bursa"
                  className="p-1 rounded-full text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-primary)] hover:bg-[var(--md-sys-color-surface-container-high)] transition-all cursor-pointer"
                >
                  <Icon
                    icon="ph:arrows-clockwise-bold"
                    className={isRefreshing ? "animate-spin text-[var(--md-sys-color-primary)]" : ""}
                    width="14"
                  />
                </button>
              </div>
            </div>
            <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-0.5">
              Relative Strength (RRG) 11 sektor IDX dipadukan radar deteksi akumulasi vs eksodus modal asing
            </p>
          </div>
        </div>

        {/* Macro Mini Badges */}
        <div className="flex items-center gap-2.5 self-start sm:self-auto flex-wrap">
          {/* IHSG Benchmark */}
          <div className="flex items-center gap-2.5 p-2 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]">
            <div className="flex flex-col">
              <div className="flex items-center gap-1.5">
                <span className="text-[9px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
                  IHSG (COMPOSITE)
                </span>
                {effectiveMarketState && (
                  <span
                    className={`text-[8px] font-extrabold px-1.5 py-0.2 rounded ${
                      effectiveMarketState === "CRASH"
                        ? "bg-rose-500/15 text-rose-600 dark:text-rose-400 border border-rose-500/30"
                        : effectiveMarketState === "WARNING"
                        ? "bg-amber-500/15 text-amber-600 dark:text-amber-400 border border-amber-500/30"
                        : effectiveMarketState === "BOTTOMING"
                        ? "bg-sky-500/15 text-sky-600 dark:text-sky-400 border border-sky-500/30"
                        : effectiveMarketState === "RECOVERY_TESTING"
                        ? "bg-amber-500/15 text-amber-600 dark:text-amber-400 border border-amber-500/30 animate-pulse"
                        : effectiveMarketState === "CONFIRMED_BOTTOM"
                        ? "bg-emerald-600/15 text-emerald-700 dark:text-emerald-400 border border-emerald-600/30"
                        : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)]"
                    }`}
                  >
                    {effectiveMarketState}
                  </span>
                )}
              </div>
              <span className="text-xs font-extrabold text-[var(--md-sys-color-on-surface)] font-mono">
                {benchmark.close.toLocaleString("id-ID", { minimumFractionDigits: 1, maximumFractionDigits: 1 })}
              </span>
            </div>
            <span
              className={`text-[10px] font-bold px-1.5 py-0.5 rounded-md ${
                benchmark.ret_1d >= 0
                  ? "bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30"
                  : "bg-rose-500/15 text-rose-600 dark:text-rose-400 border border-rose-500/30"
              }`}
            >
              {benchmark.ret_1d > 0 ? `+${benchmark.ret_1d}%` : `${benchmark.ret_1d}%`}
            </span>
          </div>

          {/* Total Foreign Net 1D (Guaranteed Always Rendered) */}
          <div className="flex items-center gap-2.5 p-2 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]">
            <div className="flex flex-col">
              <span className="text-[9px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
                Total Foreign Net (1D)
              </span>
              <div className="flex items-center gap-1.5 mt-0.5">
                <span
                  className="text-xs sm:text-sm font-extrabold font-mono"
                  style={{ color: signedColor(computedTotalForeign1d) }}
                >
                  {computedTotalForeign1d > 0
                    ? `+${fmtRadarFlow(computedTotalForeign1d)}`
                    : fmtRadarFlow(computedTotalForeign1d)}
                </span>
                <span
                  className={`text-[9px] font-bold px-1.5 py-0.2 rounded-md ${
                    computedTotalForeign1d > 0
                      ? "bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30"
                      : computedTotalForeign1d < 0
                      ? "bg-rose-500/15 text-rose-600 dark:text-rose-400 border border-rose-500/30"
                      : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)]"
                  }`}
                >
                  {computedTotalForeign1d > 0 ? "NET BUY" : computedTotalForeign1d < 0 ? "NET SELL" : "NETRAL"}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 4 QUADRANTS RRG SUMMARY CARDS */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-2.5">
        {/* LEADING */}
        <button
          onClick={() => onSelectQuadrant(selectedQuadrant === "leading" ? "all" : "leading")}
          className={`p-3 rounded-xl border text-left transition-all cursor-pointer ${
            selectedQuadrant === "leading"
              ? "bg-emerald-500/15 border-emerald-500 shadow-xs"
              : "bg-[var(--md-sys-color-surface-container-low)] border-[var(--md-sys-color-outline-variant)] hover:border-emerald-500/50"
          }`}
        >
          <div className="flex items-center justify-between text-xs font-semibold mb-1">
            <span className="flex items-center gap-1.5 text-emerald-600 dark:text-emerald-400">
              <span className="w-2.5 h-2.5 rounded-full bg-emerald-500" />
              Leading
            </span>
            <span className="px-1.5 py-0.2 rounded-full text-[10px] font-extrabold bg-emerald-500/20 text-emerald-600 dark:text-emerald-400">
              {quadrantCounts["Leading"] || 0} Sektor
            </span>
          </div>
          <p className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] line-clamp-1">
            Outperforming &amp; Momentum Kuat
          </p>
        </button>

        {/* IMPROVING */}
        <button
          onClick={() => onSelectQuadrant(selectedQuadrant === "improving" ? "all" : "improving")}
          className={`p-3 rounded-xl border text-left transition-all cursor-pointer ${
            selectedQuadrant === "improving"
              ? "bg-sky-500/15 border-sky-500 shadow-xs"
              : "bg-[var(--md-sys-color-surface-container-low)] border-[var(--md-sys-color-outline-variant)] hover:border-sky-500/50"
          }`}
        >
          <div className="flex items-center justify-between text-xs font-semibold mb-1">
            <span className="flex items-center gap-1.5 text-sky-600 dark:text-sky-400">
              <span className="w-2.5 h-2.5 rounded-full bg-sky-500" />
              Improving
            </span>
            <span className="px-1.5 py-0.2 rounded-full text-[10px] font-extrabold bg-sky-500/20 text-sky-600 dark:text-sky-400">
              {quadrantCounts["Improving"] || 0} Sektor
            </span>
          </div>
          <p className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] line-clamp-1">
            Momentum Bangkit Menuju Outperform
          </p>
        </button>

        {/* WEAKENING */}
        <button
          onClick={() => onSelectQuadrant(selectedQuadrant === "weakening" ? "all" : "weakening")}
          className={`p-3 rounded-xl border text-left transition-all cursor-pointer ${
            selectedQuadrant === "weakening"
              ? "bg-amber-500/15 border-amber-500 shadow-xs"
              : "bg-[var(--md-sys-color-surface-container-low)] border-[var(--md-sys-color-outline-variant)] hover:border-amber-500/50"
          }`}
        >
          <div className="flex items-center justify-between text-xs font-semibold mb-1">
            <span className="flex items-center gap-1.5 text-amber-600 dark:text-amber-400">
              <span className="w-2.5 h-2.5 rounded-full bg-amber-500" />
              Weakening
            </span>
            <span className="px-1.5 py-0.2 rounded-full text-[10px] font-extrabold bg-amber-500/20 text-amber-600 dark:text-amber-400">
              {quadrantCounts["Weakening"] || 0} Sektor
            </span>
          </div>
          <p className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] line-clamp-1">
            Diatas IHSG tapi Momentum Melambat
          </p>
        </button>

        {/* LAGGING */}
        <button
          onClick={() => onSelectQuadrant(selectedQuadrant === "lagging" ? "all" : "lagging")}
          className={`p-3 rounded-xl border text-left transition-all cursor-pointer ${
            selectedQuadrant === "lagging"
              ? "bg-rose-500/15 border-rose-500 shadow-xs"
              : "bg-[var(--md-sys-color-surface-container-low)] border-[var(--md-sys-color-outline-variant)] hover:border-rose-500/50"
          }`}
        >
          <div className="flex items-center justify-between text-xs font-semibold mb-1">
            <span className="flex items-center gap-1.5 text-rose-600 dark:text-rose-400">
              <span className="w-2.5 h-2.5 rounded-full bg-rose-500" />
              Lagging
            </span>
            <span className="px-1.5 py-0.2 rounded-full text-[10px] font-extrabold bg-rose-500/20 text-rose-600 dark:text-rose-400">
              {quadrantCounts["Lagging"] || 0} Sektor
            </span>
          </div>
          <p className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] line-clamp-1">
            Underperforming &amp; Aliran Lemah
          </p>
        </button>
      </div>

      {/* FOREIGN ACTION RADAR MINI TICKER */}
      <div className="flex flex-wrap items-center gap-2 pt-2 border-t border-[var(--md-sys-color-outline-variant)]/40 text-[11px]">
        <span className="text-[var(--md-sys-color-on-surface-variant)] font-semibold flex items-center gap-1">
          <Icon icon="ph:radar-duotone" width="14" className="text-[var(--md-sys-color-primary)]" />
          Distribusi Status Asing Sektoral:
        </span>
        <span className="px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 font-bold">
          {foreignStateCounts.surge} Inflow Surge
        </span>
        <span className="px-2 py-0.5 rounded-full bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20 font-bold">
          {foreignStateCounts.warning} Hidden Distrib
        </span>
        <span className="px-2 py-0.5 rounded-full bg-rose-500/10 text-rose-600 dark:text-rose-400 border border-rose-500/20 font-bold">
          {foreignStateCounts.exodus} Sector Exodus
        </span>
        <span className="px-2 py-0.5 rounded-full bg-sky-500/10 text-sky-600 dark:text-sky-400 border border-sky-500/20 font-bold">
          {foreignStateCounts.bottoming} Bottoming
        </span>
      </div>
    </div>
  );
}
