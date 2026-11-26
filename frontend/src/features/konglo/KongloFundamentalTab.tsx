import { KongloDualBox, KongloTickerLink } from "./KongloShared";
import { fmtFreq, fmtM, fmtMc, fmtNum, fmtT, type KongloReportData } from "./kongloTypes";

export default function KongloFundamentalTab({
  reportData,
  period,
}: {
  reportData: KongloReportData | null;
  period: string;
}) {
  const fundaSummary = reportData?.fundamental_summary;
  const fundaReport = reportData?.fundamental_report;
  const missingFunda = reportData?.missing_funda;

  const totalMcChg = fundaSummary?.total_mc_chg ?? 12.88;
  const mcPositive = (fundaSummary?.total_mc_chg ?? 0) >= 0;
  const mcPct = fundaSummary?.mc_pct ?? 0.41;

  return (
    <div id="tab-fundamental" className="space-y-4 animate-fadeIn">
      <div className="flex items-center gap-2 text-xs font-bold text-[var(--md-sys-color-on-surface)]">
        <span className="w-2.5 h-2.5 rounded-full bg-[var(--md-sys-color-primary)]" />
        <span className="text-sm">Funda Report</span>
        <span className="text-[var(--md-sys-color-on-surface-variant)] text-xs font-normal font-mono">
          {period}
        </span>
      </div>

      {/* ── Fundamental 3 Stat Cards ── */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-3.5">
        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl p-4 shadow-2xs">
          <div className="text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] tracking-widest mb-1.5 font-mono">
            TOTAL SAHAM FUNDAMENTAL
          </div>
          <div className="text-2xl md:text-3xl font-extrabold text-[var(--md-sys-color-on-surface)] font-mono">
            {fundaSummary?.total_stocks ?? 50}
          </div>
          <div className="text-xs text-[var(--md-sys-color-on-surface-variant)] mt-1">
            {missingFunda?.length ?? 0} tidak ada data
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl p-4 shadow-2xs">
          <div className="text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] tracking-widest mb-1.5 font-mono">
            NAIK / TURUN
          </div>
          <div className="text-2xl md:text-3xl font-extrabold flex items-baseline gap-2 font-mono">
            <span className="text-emerald-600 dark:text-emerald-400">
              {fundaSummary?.rising ?? 13}
            </span>
            <span className="text-[var(--md-sys-color-on-surface-variant)] text-lg font-normal">
              /
            </span>
            <span className="text-rose-600 dark:text-rose-400">
              {fundaSummary?.falling ?? 33}
            </span>
          </div>
          <div className="text-xs text-[var(--md-sys-color-on-surface-variant)] mt-1">
            {fundaSummary?.flat ?? 4} saham flat
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl p-4 shadow-2xs">
          <div className="text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] tracking-widest mb-1.5 font-mono">
            PERUBAHAN MC SAHAM FUNDAMENTAL
          </div>
          <div
            className={`text-2xl md:text-3xl font-extrabold font-mono ${
              mcPositive
                ? "text-emerald-600 dark:text-emerald-400"
                : "text-rose-600 dark:text-rose-400"
            }`}
          >
            {mcPositive ? "+" : "-"}Rp {fmtT(totalMcChg)} T
          </div>
          <div className="text-xs text-[var(--md-sys-color-on-surface-variant)] mt-1 font-mono">
            {mcPct >= 0 ? "+" : ""}
            {mcPct.toFixed(2)}% dari MC sebelumnya
          </div>
        </div>
      </div>

      {/* ── Fundamental Dual Box: Top Gainer & Top Loser ── */}
      <KongloDualBox
        gainerTitle="Top Gainer Fundamental"
        loserTitle="Top Loser Fundamental"
        badgeMode="funda"
        top5PctUp={fundaSummary?.top5_pct_up}
        top5McUp={fundaSummary?.top5_mc_up}
        top5PctDown={fundaSummary?.top5_pct_down}
        top5McDown={fundaSummary?.top5_mc_down}
      />

      {/* ── Full 50 Fundamental Stocks Table ── */}
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl overflow-hidden shadow-2xs">
        <div className="px-4 py-3 border-b border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-high)]/60 flex items-center justify-between">
          <span className="text-xs font-bold tracking-wide text-[var(--md-sys-color-on-surface)] uppercase">
            DAFTAR LENGKAP 50 SAHAM FUNDAMENTAL
          </span>
          <span className="text-xs text-[var(--md-sys-color-on-surface-variant)] font-mono">
            {fundaReport?.length ?? 0} Saham
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead>
              <tr className="border-b border-[var(--md-sys-color-outline-variant)]/40 bg-[var(--md-sys-color-surface-container-high)]/30 text-[var(--md-sys-color-on-surface-variant)] font-mono uppercase text-[10px]">
                <th className="py-2.5 px-3 w-8 font-bold">#</th>
                <th className="py-2.5 px-3 font-bold">Saham</th>
                <th className="py-2.5 px-3 font-bold">Harga</th>
                <th className="py-2.5 px-3 text-right font-bold">Change (%)</th>
                <th className="py-2.5 px-3 text-right font-bold">Value (M)</th>
                <th className="py-2.5 px-3 text-right font-bold">Freq</th>
                <th className="py-2.5 px-3 text-right font-bold">Perubahan MC</th>
              </tr>
            </thead>
            <tbody>
              {fundaReport?.map((s, idx) => (
                <tr
                  key={s.ticker}
                  className="border-b border-[var(--md-sys-color-outline-variant)]/20 hover:bg-[var(--md-sys-color-surface-container-high)]/40 transition-colors"
                >
                  <td className="py-2 px-3 text-[var(--md-sys-color-on-surface-variant)] font-mono">
                    {idx + 1}
                  </td>
                  <td className="py-2 px-3">
                    <KongloTickerLink ticker={s.ticker} className="text-sm" />
                  </td>
                  <td className="py-2 px-3 text-[var(--md-sys-color-on-surface-variant)] font-mono">
                    {fmtNum(s.prev)} → {fmtNum(s.close)}
                  </td>
                  <td
                    className={`py-2 px-3 text-right font-semibold font-mono ${
                      s.pct >= 0
                        ? "text-emerald-600 dark:text-emerald-400"
                        : "text-rose-600 dark:text-rose-400"
                    }`}
                  >
                    {s.pct > 0 ? `+${s.pct.toFixed(1)}%` : `${s.pct.toFixed(1)}%`}
                  </td>
                  <td className="py-2 px-3 text-right text-[var(--md-sys-color-on-surface-variant)] font-mono">
                    {fmtM(s.val)}
                  </td>
                  <td className="py-2 px-3 text-right text-[var(--md-sys-color-on-surface-variant)] font-mono">
                    {fmtFreq(s.freq)}
                  </td>
                  <td
                    className={`py-2 px-3 text-right font-semibold font-mono ${
                      s.mc_chg >= 0
                        ? "text-emerald-600 dark:text-emerald-400"
                        : "text-rose-600 dark:text-rose-400"
                    }`}
                  >
                    {fmtMc(s.mc_chg)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* ── Fundamental Missing Banner if any ── */}
      {missingFunda && missingFunda.length > 0 && (
        <div
          id="missingBannerFunda"
          className="p-3.5 rounded-xl bg-amber-500/10 border border-amber-500/30 text-xs space-y-1"
        >
          <div className="font-bold text-amber-700 dark:text-amber-400 flex items-center gap-1.5 font-mono">
            <span>⚠️</span>
            <span>{missingFunda.length} saham fundamental tidak ada datanya</span>
          </div>
          <div className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
            {missingFunda.join(" · ")}
          </div>
        </div>
      )}
    </div>
  );
}
