import { Badge } from "@/components/research/ui";
import { categoryLabel, fmtPct } from "@/components/research/format";
import type { Holder } from "./stockTypes";

export function SnapshotHoldersTable({ holders }: { holders: Holder[] }) {
  return (
    <table className="w-full text-left text-[11px]">
      <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
        <tr>
          <th className="px-4 py-2.5 font-bold">#</th>
          <th className="px-4 py-2.5 font-bold">Nama</th>
          <th className="px-4 py-2.5 font-bold">Kategori</th>
          <th className="px-4 py-2.5 font-bold text-right">Porsi</th>
          <th className="px-4 py-2.5 font-bold">Status</th>
        </tr>
      </thead>
      <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
        {holders.map((holder, index) => (
          <tr key={`${holder.name}-${index}`} className="hover:bg-[var(--md-sys-color-surface-container-high)]/40">
            <td className="px-4 py-2.5 font-mono text-[var(--md-sys-color-on-surface-variant)]">
              {index + 1}
            </td>
            <td className="px-4 py-2.5 font-medium text-[var(--md-sys-color-on-surface)] max-w-[360px]">
              <span className="line-clamp-2">{holder.name}</span>
            </td>
            <td className="px-4 py-2.5">
              <Badge tone={holder.category === "controlling" ? "warning" : "neutral"}>
                {categoryLabel(holder.category)}
              </Badge>
            </td>
            <td className="px-4 py-2.5 text-right font-mono font-bold text-[var(--md-sys-color-on-surface)]">
              {fmtPct(holder.pct, 3)}
            </td>
            <td className="px-4 py-2.5">
              {holder.is_controlling ? (
                <Badge tone="primary">Pengendali</Badge>
              ) : (
                <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">—</span>
              )}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
