import { Icon } from "@iconify/react";
import { fmtRp, fmtRadarFlow, signedColor } from "@/lib/format";
import type { SectorHistoryPoint } from "./sectorTypes";

interface SectorHistoryTableProps {
  sectorCode: string;
  historyTimeline: SectorHistoryPoint[];
  isHistoryLoading: boolean;
}

export default function SectorHistoryTable({
  sectorCode,
  historyTimeline,
  isHistoryLoading,
}: SectorHistoryTableProps) {
  return (
    <div className="space-y-2">
      <span className="text-xs font-semibold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
        <Icon icon="ph:clock-counter-clockwise-duotone" width="14" className="text-[var(--md-sys-color-primary)]" />
        Histori Harian Indeks, Transaksi &amp; Arus Asing ({sectorCode} vs COMPOSITE)
      </span>

      {isHistoryLoading ? (
        <div className="py-6 text-center text-xs text-[var(--md-sys-color-on-surface-variant)] flex items-center justify-center gap-2">
          <Icon icon="ph:spinner-gap-bold" className="animate-spin text-[var(--md-sys-color-primary)]" width="16" />
          <span>Memuat deret waktu sektor...</span>
        </div>
      ) : historyTimeline.length > 0 ? (
        <div className="overflow-x-auto rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-low)]">
          <table className="w-full text-left text-xs">
            <thead>
              <tr className="border-b border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] font-semibold text-[11px]">
                <th className="py-2 px-3">Tanggal</th>
                <th className="py-2 px-3 text-right">Close Indeks</th>
                <th className="py-2 px-3 text-right">Turnover</th>
                <th className="py-2 px-3 text-right">Foreign Net</th>
                <th className="py-2 px-3 text-right">RS Ratio</th>
                <th className="py-2 px-3 text-center">RRG Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/50">
              {[...historyTimeline].reverse().slice(0, 10).map((pt) => {
                const isStrong = pt.rs_ratio >= 100;
                return (
                  <tr key={pt.date} className="hover:bg-[var(--md-sys-color-surface-container)] transition-colors">
                    <td className="py-2 px-3 font-medium text-[var(--md-sys-color-on-surface)] font-mono text-[11px]">
                      {pt.date}
                    </td>
                    <td className="py-2 px-3 text-right font-mono font-semibold text-[var(--md-sys-color-on-surface)]">
                      {pt.close.toLocaleString("id-ID", { minimumFractionDigits: 1, maximumFractionDigits: 1 })}
                    </td>
                    <td className="py-2 px-3 text-right font-mono text-[var(--md-sys-color-on-surface-variant)]">
                      {fmtRp(pt.value)}
                    </td>
                    <td
                      className="py-2 px-3 text-right font-mono font-bold"
                      style={{ color: signedColor(pt.foreign_net) }}
                    >
                      {pt.foreign_net > 0 ? `+${fmtRadarFlow(pt.foreign_net)}` : fmtRadarFlow(pt.foreign_net)}
                    </td>
                    <td className="py-2 px-3 text-right font-mono font-bold text-[var(--md-sys-color-on-surface)]">
                      {pt.rs_ratio}
                    </td>
                    <td className="py-2 px-3 text-center">
                      <span
                        className={`text-[9px] font-bold px-2 py-0.5 rounded-full ${
                          isStrong
                            ? "bg-emerald-500/15 text-emerald-600 dark:text-emerald-400"
                            : "bg-rose-500/15 text-rose-600 dark:text-rose-400"
                        }`}
                      >
                        {isStrong ? "OUTPERFORM" : "UNDERPERFORM"}
                      </span>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      ) : null}
    </div>
  );
}
