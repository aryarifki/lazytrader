import { KongloDualBox } from "./KongloShared";
import type { KongloReportData } from "./kongloTypes";

export default function KongloSummaryTab({
  reportData,
  period,
}: {
  reportData: KongloReportData | null;
  period: string;
}) {
  const top5 = reportData?.top5_summary;

  return (
    <div id="tab-summary" className="space-y-4 animate-fadeIn">
      <div className="flex items-center gap-2 text-xs font-bold text-[var(--md-sys-color-on-surface)]">
        <span className="w-2.5 h-2.5 rounded-full bg-[var(--md-sys-color-primary)]" />
        <span className="text-sm">Top Gainer &amp; Loser</span>
        <span className="text-[var(--md-sys-color-on-surface-variant)] text-xs font-normal font-mono">
          {period}
        </span>
      </div>

      <KongloDualBox
        gainerTitle="Top Gainer"
        loserTitle="Top Loser"
        badgeMode="group"
        top5PctUp={top5?.top5_pct_up}
        top5McUp={top5?.top5_mc_up}
        top5PctDown={top5?.top5_pct_down}
        top5McDown={top5?.top5_mc_down}
      />
    </div>
  );
}
