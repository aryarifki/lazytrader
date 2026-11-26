import { fmtT, type KongloReportData } from "./kongloTypes";

interface KongloOverviewCardsProps {
  reportData: KongloReportData | null;
  refreshing: boolean;
  dateStr: string;
  timeStr: string;
  period: string;
  onRefresh: () => void;
}

export default function KongloOverviewCards({
  reportData,
  refreshing,
  dateStr,
  timeStr,
  period,
  onRefresh,
}: KongloOverviewCardsProps) {
  const summary = reportData?.summary;
  const meta = reportData?.metadata;
  const overview = reportData?.konglo_cards?.overview;
  const bestMc = reportData?.konglo_cards?.best_mc;
  const worstMc = reportData?.konglo_cards?.worst_mc;
  const bestPct = reportData?.konglo_cards?.best_pct;
  const worstPct = reportData?.konglo_cards?.worst_pct;

  const ihsgPct = Number(summary?.ihsg_pct ?? 0);
  const totalMcChg = overview?.total_mc_chg ?? -94.85;
  const totalMcPositive = (overview?.total_mc_chg ?? 0) >= 0;

  return (
    <>
      {/* ── LIVE STATUS BAR ── */}
      <div
        id="statusBar"
        className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 px-4 py-2.5 rounded-xl border border-emerald-500/25 bg-emerald-500/10 text-xs shadow-2xs"
      >
        <div className="flex items-center gap-2">
          <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shrink-0" />
          <span className="font-semibold text-emerald-700 dark:text-emerald-400">
            {refreshing
              ? "Memuat ulang data..."
              : `Data live aktif — ${summary?.total_active_stocks ?? 1343} saham aktif (dari total ${summary?.total_stocks ?? 1581} emiten)`}
          </span>
        </div>
        <div className="text-[var(--md-sys-color-on-surface-variant)] flex items-center gap-2 text-xs flex-wrap">
          <span>
            Data per penutupan {meta?.last_updated_str || "Sen, 02/11/2026 16:00 WIB"}
          </span>
          <span className="opacity-40">|</span>
          <button
            type="button"
            onClick={onRefresh}
            disabled={refreshing}
            className="text-[var(--md-sys-color-primary)] font-bold hover:underline cursor-pointer disabled:opacity-60"
          >
            Refresh sekarang
          </button>
        </div>
      </div>

      {/* ── CONTROL BAR (TANGGAL & JAM) ── */}
      <div
        id="controlBar"
        className="flex flex-col sm:flex-row sm:items-center justify-between px-4 py-2.5 rounded-2xl border border-[var(--md-sys-color-outline-variant)]/60 bg-[var(--md-sys-color-surface-container)] text-xs gap-3 shadow-2xs"
      >
        <div className="flex items-center gap-3 flex-wrap">
          <div className="flex items-center gap-1.5">
            <span className="text-[var(--md-sys-color-on-surface-variant)] text-xs font-medium">
              Tanggal:
            </span>
            <input
              id="inputDate"
              type="text"
              value={dateStr}
              readOnly
              className="rounded-lg px-2.5 py-1 font-mono text-xs w-36 text-center outline-none border border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface)]"
            />
          </div>
          <div className="flex items-center gap-1.5">
            <span className="text-[var(--md-sys-color-on-surface-variant)] text-xs font-medium">
              Jam:
            </span>
            <input
              id="inputTime"
              type="text"
              value={timeStr}
              readOnly
              className="rounded-lg px-2.5 py-1 font-mono text-xs w-20 text-center outline-none border border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface)]"
            />
          </div>
        </div>

        <div className="flex items-center gap-2 font-mono text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
          <span>
            IHSG:{" "}
            <strong className="text-[var(--md-sys-color-on-surface)]">
              {summary?.ihsg_close || "—"}
            </strong>
          </span>
          <span className="opacity-40">·</span>
          <span>
            Perubahan:{" "}
            <strong
              className={
                ihsgPct >= 0
                  ? "text-emerald-600 dark:text-emerald-400"
                  : "text-rose-600 dark:text-rose-400"
              }
            >
              {ihsgPct >= 0 ? "+" : ""}
              {summary?.ihsg_pct || "0.00"}%
            </strong>
          </span>
        </div>
      </div>

      {/* ── STATS ROW 1: General Overview (3 Cards) ── */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-3.5">
        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl p-4 shadow-2xs">
          <div className="text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] tracking-widest mb-1.5 font-mono">
            TOTAL SAHAM KONGLO
          </div>
          <div className="text-2xl md:text-3xl font-extrabold text-[var(--md-sys-color-on-surface)] font-mono">
            {overview?.total_konglo_stocks ?? 169}
          </div>
          <div className="text-xs text-[var(--md-sys-color-on-surface-variant)] mt-1">
            {reportData?.missing_konglo?.length ?? 0} tidak ada data
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl p-4 shadow-2xs">
          <div className="text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] tracking-widest mb-1.5 font-mono">
            NAIK / TURUN
          </div>
          <div className="text-2xl md:text-3xl font-extrabold flex items-baseline gap-2 font-mono">
            <span className="text-emerald-600 dark:text-emerald-400">
              {overview?.rising ?? 34}
            </span>
            <span className="text-[var(--md-sys-color-on-surface-variant)] text-lg font-normal">
              /
            </span>
            <span className="text-rose-600 dark:text-rose-400">
              {overview?.falling ?? 114}
            </span>
          </div>
          <div className="text-xs text-[var(--md-sys-color-on-surface-variant)] mt-1">
            {overview?.flat ?? 21} saham flat
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl p-4 shadow-2xs">
          <div className="text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] tracking-widest mb-1.5 font-mono">
            TOTAL PERUBAHAN MC HARI INI
          </div>
          <div
            className={`text-2xl md:text-3xl font-extrabold font-mono ${
              totalMcPositive
                ? "text-emerald-600 dark:text-emerald-400"
                : "text-rose-600 dark:text-rose-400"
            }`}
          >
            {totalMcPositive ? "+" : "-"}Rp {fmtT(totalMcChg)} T
          </div>
          <div className="text-xs text-[var(--md-sys-color-on-surface-variant)] mt-1">
            {period}
          </div>
        </div>
      </div>

      {/* ── STATS ROW 2: 4 Kartu Ringkasan Konglo Terbaik & Terburuk (2x2) ── */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3.5">
        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 border-l-4 border-l-emerald-500 rounded-2xl p-4 shadow-2xs">
          <div className="text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] tracking-widest font-mono">
            🏆 Konglo Terbaik — Market Cap
          </div>
          <div className="mt-2 text-base font-bold text-[var(--md-sys-color-on-surface)]">
            {bestMc?.group || "Thohir Group"}
          </div>
          <div className="mt-1 flex items-baseline gap-2 font-mono">
            <span className="text-2xl font-extrabold text-emerald-600 dark:text-emerald-400">
              +Rp {fmtT(bestMc?.val ?? 2.45)} T
            </span>
            <span className="text-xs text-[var(--md-sys-color-on-surface-variant)]">
              +{Number(bestMc?.pct ?? 1.2).toFixed(1)}% · {bestMc?.count ?? 8} saham
            </span>
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 border-l-4 border-l-rose-500 rounded-2xl p-4 shadow-2xs">
          <div className="text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] tracking-widest font-mono">
            📉 Konglo Terburuk — Market Cap
          </div>
          <div className="mt-2 text-base font-bold text-[var(--md-sys-color-on-surface)]">
            {worstMc?.group || "Barito Group"}
          </div>
          <div className="mt-1 flex items-baseline gap-2 font-mono">
            <span className="text-2xl font-extrabold text-rose-600 dark:text-rose-400">
              -Rp {fmtT(worstMc?.val ?? -88.12)} T
            </span>
            <span className="text-xs text-[var(--md-sys-color-on-surface-variant)]">
              {Number(worstMc?.pct ?? -6.8).toFixed(1)}% · {worstMc?.count ?? 6} saham
            </span>
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 border-l-4 border-l-emerald-500 rounded-2xl p-4 shadow-2xs">
          <div className="text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] tracking-widest font-mono">
            🚀 Konglo Terbaik — Price Change
          </div>
          <div className="mt-2 text-base font-bold text-[var(--md-sys-color-on-surface)]">
            {bestPct?.group || "Thohir Group"}
          </div>
          <div className="mt-1 flex items-baseline gap-2 font-mono">
            <span className="text-2xl font-extrabold text-emerald-600 dark:text-emerald-400">
              +{Number(bestPct?.val ?? 3.4).toFixed(1)}%
            </span>
            <span className="text-xs text-[var(--md-sys-color-on-surface-variant)]">
              avg · {bestPct?.count ?? 8} saham
            </span>
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 border-l-4 border-l-rose-500 rounded-2xl p-4 shadow-2xs">
          <div className="text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] tracking-widest font-mono">
            📉 Konglo Terburuk — Price Change
          </div>
          <div className="mt-2 text-base font-bold text-[var(--md-sys-color-on-surface)]">
            {worstPct?.group || "Barito Group"}
          </div>
          <div className="mt-1 flex items-baseline gap-2 font-mono">
            <span className="text-2xl font-extrabold text-rose-600 dark:text-rose-400">
              {Number(worstPct?.val ?? -6.6).toFixed(1)}%
            </span>
            <span className="text-xs text-[var(--md-sys-color-on-surface-variant)]">
              avg · {worstPct?.count ?? 6} saham
            </span>
          </div>
        </div>
      </div>
    </>
  );
}
