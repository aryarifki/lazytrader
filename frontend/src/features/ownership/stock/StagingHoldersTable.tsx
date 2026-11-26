import { Badge } from "@/components/research/ui";
import { fmtNumber, fmtPct } from "@/components/research/format";
import { fmtLot, HolderStatus } from "./stockHelpers";
import type { StagingHolderItem } from "./stockTypes";

export function StagingHoldersTable({ holders }: { holders: StagingHolderItem[] }) {
  return (
    <table className="w-full min-w-[850px] text-left text-[11px]">
      <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
        <tr>
          <th className="px-4 py-2.5 font-bold">#</th>
          <th className="px-4 py-2.5 font-bold">Nama Pemegang Saham</th>
          <th className="px-4 py-2.5 font-bold">Tipe & Domisili</th>
          <th className="px-4 py-2.5 font-bold text-right">Jumlah Saham</th>
          <th className="px-4 py-2.5 font-bold text-right">Bulan Lalu %</th>
          <th className="px-4 py-2.5 font-bold text-right">Saat Ini %</th>
          <th className="px-4 py-2.5 font-bold text-right">Perubahan</th>
          <th className="px-4 py-2.5 font-bold">Status</th>
        </tr>
      </thead>
      <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
        {holders.map((holder, index) => {
          const isUp = holder.change_pct > 0;
          const isDown = holder.change_pct < 0;
          return (
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
                <span className="inline-flex items-center gap-1">
                  <Badge tone="neutral">
                    {holder.investor_type ?? "—"}
                  </Badge>
                  <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
                    {holder.location === "L" ? "Lokal" : holder.location === "A" || holder.location === "F" ? "Asing" : ""}
                  </span>
                </span>
              </td>
              <td className="px-4 py-2.5 text-right font-mono whitespace-nowrap">
                {holder.total_shares != null ? (
                  <>
                    <div className="font-semibold text-[var(--md-sys-color-on-surface)]">
                      {fmtLot(Math.round(holder.total_shares / 100))} lot
                    </div>
                    <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
                      {fmtNumber(holder.total_shares)} lbr
                    </div>
                  </>
                ) : (
                  "—"
                )}
              </td>
              <td className="px-4 py-2.5 text-right font-mono text-[var(--md-sys-color-on-surface-variant)] whitespace-nowrap">
                {holder.previous_pct != null ? fmtPct(holder.previous_pct) : "—"}
              </td>
              <td className="px-4 py-2.5 text-right font-mono font-bold text-[var(--md-sys-color-on-surface)] whitespace-nowrap">
                {fmtPct(holder.pct)}
              </td>
              <td
                className={`px-4 py-2.5 text-right font-mono font-bold whitespace-nowrap ${
                  isUp
                    ? "text-emerald-500"
                    : isDown
                      ? "text-rose-500"
                      : "text-[var(--md-sys-color-on-surface-variant)]"
                }`}
              >
                {isUp ? "+" : ""}{fmtPct(holder.change_pct)}
                {holder.change_shares !== 0 && (
                  <div className="text-[9px] font-normal">
                    {holder.change_shares > 0 ? "+" : ""}{fmtLot(Math.round(holder.change_shares / 100))} lot
                  </div>
                )}
              </td>
              <td className="px-4 py-2.5">
                <HolderStatus status={holder.status} />
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}
