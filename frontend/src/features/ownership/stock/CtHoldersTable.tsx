import { Badge } from "@/components/research/ui";
import { fmtPct } from "@/components/research/format";
import { fmtLot, HolderStatus } from "./stockHelpers";
import type { CtHolder } from "./stockTypes";

export function CtHoldersTable({ holders }: { holders: CtHolder[] }) {
  return (
    <table className="w-full min-w-[820px] text-left text-[11px]">
      <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
        <tr>
          <th className="px-4 py-2.5 font-bold">#</th>
          <th className="px-4 py-2.5 font-bold">Nama</th>
          <th className="px-4 py-2.5 font-bold">Tipe</th>
          <th className="px-4 py-2.5 font-bold text-right">Mei (lot)</th>
          <th className="px-4 py-2.5 font-bold text-right">Mei %</th>
          <th className="px-4 py-2.5 font-bold text-right">Jun (lot)</th>
          <th className="px-4 py-2.5 font-bold text-right">Jun %</th>
          <th className="px-4 py-2.5 font-bold">Status</th>
        </tr>
      </thead>
      <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
        {holders.map((holder, index) => (
          <tr
            key={`${holder.name}-${index}`}
            className="hover:bg-[var(--md-sys-color-surface-container-high)]/40"
          >
            <td className="px-4 py-2.5 font-mono text-[var(--md-sys-color-on-surface-variant)]">
              {index + 1}
            </td>
            <td className="px-4 py-2.5 font-medium text-[var(--md-sys-color-on-surface)] max-w-[320px]">
              <span className="line-clamp-2">{holder.name}</span>
            </td>
            <td className="px-4 py-2.5">
              <Badge tone="neutral">{holder.lf ?? "—"}</Badge>
            </td>
            <td className="px-4 py-2.5 text-right font-mono whitespace-nowrap">
              {fmtLot(holder.sl_mei)}
            </td>
            <td className="px-4 py-2.5 text-right font-mono font-bold whitespace-nowrap">
              {fmtPct(holder.pct_mei)}
            </td>
            <td className="px-4 py-2.5 text-right font-mono whitespace-nowrap">
              {fmtLot(holder.sl_jun)}
            </td>
            <td className="px-4 py-2.5 text-right font-mono font-bold whitespace-nowrap">
              {fmtPct(holder.pct_jun)}
            </td>
            <td className="px-4 py-2.5">
              <HolderStatus status={holder.status} />
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
