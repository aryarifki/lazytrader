import { KongloGroupBadge, KongloRankBadge, KongloTickerLink } from "./KongloShared";
import { fmtFreq, fmtM, fmtMc, fmtNum, type KongloReportData } from "./kongloTypes";

export type KongloRankingTabId = "pct-up" | "mc-up" | "pct-down" | "mc-down";

const TAB_TITLES: Record<KongloRankingTabId, string> = {
  "pct-up": "Top 15 — % kenaikan terbesar",
  "mc-up": "Top 15 — kenaikan market cap terbesar",
  "pct-down": "Top 15 — % penurunan terbesar",
  "mc-down": "Top 15 — penurunan market cap terbesar",
};

export default function KongloRankingsTab({
  activeTab,
  reportData,
  period,
}: {
  activeTab: KongloRankingTabId;
  reportData: KongloReportData | null;
  period: string;
}) {
  const isUp = activeTab === "pct-up" || activeTab === "mc-up";
  const tables = reportData?.top15_tables;
  const rows =
    activeTab === "pct-up"
      ? tables?.pct_up
      : activeTab === "mc-up"
        ? tables?.mc_up
        : activeTab === "pct-down"
          ? tables?.pct_down
          : tables?.mc_down;

  const mcColumnHeader =
    activeTab === "mc-up"
      ? "Kenaikan Market Cap"
      : activeTab === "mc-down"
        ? "Penurunan Market Cap"
        : "Perubahan Market Cap";

  return (
    <div
      id={`tab-${activeTab}`}
      className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl overflow-hidden shadow-2xs animate-fadeIn"
    >
      <div className="px-4 py-3 border-b border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-high)]/60 flex items-center justify-between gap-2">
        <div className="flex items-center gap-2">
          <span
            className={`w-2 h-2 rounded-full ${isUp ? "bg-emerald-500" : "bg-rose-500"}`}
          />
          <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)]">
            {TAB_TITLES[activeTab]}
          </span>
        </div>
        <span className="text-xs text-[var(--md-sys-color-on-surface-variant)] font-mono">
          {period}
        </span>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left text-xs">
          <thead>
            <tr className="border-b border-[var(--md-sys-color-outline-variant)]/40 bg-[var(--md-sys-color-surface-container-high)]/30 text-[var(--md-sys-color-on-surface-variant)] font-mono uppercase text-[10px]">
              <th className="py-2.5 px-3 w-9 font-bold">#</th>
              <th className="py-2.5 px-3 font-bold">Saham</th>
              <th className="py-2.5 px-3 font-bold">Konglomerasi</th>
              <th className="py-2.5 px-3 text-right font-bold">Change (%)</th>
              <th className="py-2.5 px-3 text-right font-bold">Value (M)</th>
              <th className="py-2.5 px-3 text-right font-bold">Freq</th>
              <th className="py-2.5 px-3 text-right font-bold">{mcColumnHeader}</th>
            </tr>
          </thead>
          <tbody>
            {rows?.map((s, idx) => (
              <tr
                key={s.ticker}
                className="border-b border-[var(--md-sys-color-outline-variant)]/20 hover:bg-[var(--md-sys-color-surface-container-high)]/40 transition-colors"
              >
                <td className="py-2 px-3">
                  <KongloRankBadge index={idx} direction={isUp ? "up" : "down"} />
                </td>
                <td className="py-2 px-3">
                  <div className="flex items-center gap-2">
                    <KongloTickerLink ticker={s.ticker} className="text-sm" />
                    <span className="text-xs text-[var(--md-sys-color-on-surface-variant)] font-mono">
                      · {fmtNum(s.prev)} → {fmtNum(s.close)}
                    </span>
                  </div>
                </td>
                <td className="py-2 px-3">
                  <KongloGroupBadge group={s.group} size="sm" />
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
  );
}
