import { Icon } from "@iconify/react";

import { Badge } from "@/components/research/ui";
import { fmtNumber, fmtPct } from "@/components/research/format";
import { MomPill } from "./MomPill";
import type { MonthlyComparisonRow } from "./stockTypes";

export function StockMonthlyMomTable({
  rows,
}: {
  rows: MonthlyComparisonRow[];
}) {
  return (
    <div className="mt-4 pt-3.5 border-t border-[var(--md-sys-color-outline-variant)]/60">
      <div className="flex items-center justify-between mb-2">
        <span className="text-[11px] font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
          <Icon icon="ph:clock-counter-clockwise-duotone" width="13" className="text-primary-500" />
          Rincian Perbandingan Tiap Bulan (MoM)
        </span>
        <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
          Urutan: Terbaru ➔ Terlama
        </span>
      </div>

      <div className="overflow-x-auto rounded-lg border border-[var(--md-sys-color-outline-variant)]/40">
        <table className="w-full text-left text-[11px] border-collapse">
          <thead>
            <tr className="bg-[var(--md-sys-color-surface-container)] text-[10px] font-mono uppercase text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/40">
              <th className="py-2 px-2 font-semibold">Bulan</th>
              <th className="py-2 px-2 font-semibold text-right">Asing Total</th>
              <th className="py-2 px-2 font-semibold text-right">Foreign Hidden</th>
              <th className="py-2 px-2 font-semibold text-right">Ritel</th>
              <th className="py-2 px-2 font-semibold text-right">Lokal Non-Ritel</th>
              <th className="py-2 px-2 font-semibold text-right">SID Investor</th>
              <th className="py-2 px-2 font-semibold text-center">Arus CT</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/30">
            {rows.map((row) => (
              <tr
                key={row.period_label}
                className={`hover:bg-[var(--md-sys-color-surface-container-high)]/60 transition-colors ${
                  row.isLatest ? "bg-primary-500/5 font-medium" : ""
                }`}
              >
                <td className="py-2 px-2.5 whitespace-nowrap">
                  <div className="flex items-center gap-1.5">
                    <span className="font-mono text-[11px] font-bold text-[var(--md-sys-color-on-surface)]">
                      {row.period_label}
                    </span>
                    {row.isLatest ? (
                      <span className="text-[8px] font-bold uppercase tracking-wider px-1 py-0.2 rounded bg-primary-500/20 text-primary-400">
                        Terbaru
                      </span>
                    ) : null}
                  </div>
                </td>
                <td className="py-2 px-2.5 text-right whitespace-nowrap font-mono">
                  <div>{fmtPct(row.currForeign)}</div>
                  <MomPill delta={row.foreignDelta} />
                </td>
                <td className="py-2 px-2.5 text-right whitespace-nowrap font-mono">
                  <div>{fmtPct(row.currForeignHidden)}</div>
                  <MomPill delta={row.foreignHiddenDelta} />
                </td>
                <td className="py-2 px-2.5 text-right whitespace-nowrap font-mono">
                  <div>{fmtPct(row.currRetail)}</div>
                  <MomPill delta={row.retailDelta} reverse />
                </td>
                <td className="py-2 px-2.5 text-right whitespace-nowrap font-mono">
                  <div>{fmtPct(row.currLocalNonRetail)}</div>
                  <MomPill delta={row.localNonRetailDelta} />
                </td>
                <td className="py-2 px-2.5 text-right whitespace-nowrap font-mono">
                  <div>{row.currSid != null ? fmtNumber(row.currSid) : "—"}</div>
                  {row.sidDelta != null ? (
                    <div className="text-[9px] flex items-center justify-end gap-0.5">
                      {row.sidDelta > 0 ? (
                        <span className="text-emerald-500 inline-flex items-center">
                          <Icon icon="ph:trend-up-bold" width="9" />
                          +{fmtNumber(row.sidDelta)} ({row.sidDeltaPct != null ? `${row.sidDeltaPct > 0 ? "+" : ""}${row.sidDeltaPct.toFixed(1)}%` : "Naik"})
                        </span>
                      ) : row.sidDelta < 0 ? (
                        <span className="text-rose-500 inline-flex items-center">
                          <Icon icon="ph:trend-down-bold" width="9" />
                          {fmtNumber(row.sidDelta)} ({row.sidDeltaPct != null ? `${row.sidDeltaPct.toFixed(1)}%` : "Turun"})
                        </span>
                      ) : (
                        <span className="text-[var(--md-sys-color-on-surface-variant)]">0 (Tetap)</span>
                      )}
                    </div>
                  ) : (
                    <span className="text-[9px] text-[var(--md-sys-color-on-surface-variant)]">—</span>
                  )}
                </td>
                <td className="py-2 px-2.5 text-center whitespace-nowrap">
                  <Badge tone={row.flowTone}>{row.flowSignal}</Badge>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
