import { Badge, Card, SectionTitle } from "@/components/research/ui";
import { fmtNumber } from "@/components/research/format";
import type { IhsgSummaryRow } from "./ihsgTypes";

export function IhsgMutationTable({ rows }: { rows: IhsgSummaryRow[] }) {
  return (
    <Card className="p-4 sm:p-5">
      <SectionTitle
        icon="ph:table-duotone"
        title="Mutasi Makro Kepemilikan IHSG"
        subtitle="Perbandingan agregat bulanan Scriptless KSEI (Agustus vs September)"
        action={<Badge tone="primary">Per September 2026</Badge>}
      />
      <div className="mt-4 overflow-x-auto">
        <table className="w-full text-left text-[11px]">
          <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
            <tr>
              <th className="px-4 py-2.5 font-bold">Kategori Kepemilikan</th>
              <th className="px-4 py-2.5 font-bold text-right">Agustus (Triliun)</th>
              <th className="px-4 py-2.5 font-bold text-right">September (Triliun)</th>
              <th className="px-4 py-2.5 font-bold text-right">Perubahan (Δ Triliun)</th>
              <th className="px-4 py-2.5 font-bold text-right">Porsi (%)</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40 font-mono">
            {rows.map((row) => {
              const diff = row.diff_trn ?? 0;
              const diffClass =
                diff > 0
                  ? "text-emerald-500 font-bold"
                  : diff < 0
                    ? "text-rose-500 font-bold"
                    : "text-[var(--md-sys-color-on-surface-variant)]";
              return (
                <tr
                  key={row.category}
                  className="hover:bg-[var(--md-sys-color-surface-container-high)]/40 transition-colors"
                >
                  <td className="px-4 py-2.5 font-sans font-medium text-[var(--md-sys-color-on-surface)]">
                    {row.category}
                  </td>
                  <td className="px-4 py-2.5 text-right">
                    {row.prev_trn != null ? `Rp ${fmtNumber(row.prev_trn, 2)} T` : "—"}
                  </td>
                  <td className="px-4 py-2.5 text-right font-bold text-[var(--md-sys-color-on-surface)]">
                    {row.curr_trn != null ? `Rp ${fmtNumber(row.curr_trn, 2)} T` : "—"}
                  </td>
                  <td className={`px-4 py-2.5 text-right ${diffClass}`}>
                    {diff > 0 ? `+${fmtNumber(diff, 2)}` : fmtNumber(diff, 2)} T
                  </td>
                  <td className="px-4 py-2.5 text-right font-bold text-[var(--md-sys-color-primary)]">
                    {row.pct != null ? `${row.pct.toFixed(2)}%` : "—"}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </Card>
  );
}
