import { Link } from "react-router-dom";

import { Badge, Card, SectionTitle } from "@/components/research/ui";
import { fmtNumber } from "@/components/research/format";
import { StockBadges } from "@/components/research/StockBadges";
import type { IhsgTop20Item } from "./ihsgTypes";

export function IhsgTop20Table({ top20 }: { top20: IhsgTop20Item[] }) {
  return (
    <Card className="overflow-hidden">
      <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <SectionTitle
          icon="ph:trophy-duotone"
          title="Top 20 Emiten Saham Terbesar Kepemilikan IHSG"
          subtitle="20 emiten dengan bobot nilai kepemilikan saham terbesar di BEI beserta label dan konglomerasi"
          action={<Badge tone="primary">{top20.length} Emiten</Badge>}
        />
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left text-[11px]">
          <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
            <tr>
              <th className="px-3 py-2.5 font-bold text-center w-10">#</th>
              <th className="px-3 py-2.5 font-bold">Saham</th>
              <th className="px-3 py-2.5 font-bold text-right">Foreign Named</th>
              <th className="px-3 py-2.5 font-bold text-right">Foreign Hidden</th>
              <th className="px-3 py-2.5 font-bold text-right">Local Non-Retail</th>
              <th className="px-3 py-2.5 font-bold text-right">Retail</th>
              <th className="px-3 py-2.5 font-bold text-right">Free Float</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40 font-mono">
            {top20.map((emiten) => (
              <tr
                key={emiten.ticker}
                className="hover:bg-[var(--md-sys-color-surface-container-high)]/40 transition-colors"
              >
                <td className="px-3 py-2.5 text-center font-bold text-[var(--md-sys-color-on-surface-variant)]">
                  {emiten.rank}
                </td>
                <td className="px-3 py-2.5">
                  <div className="flex items-center gap-1.5 flex-wrap">
                    <Link
                      to={`/ownership?tab=saham&ticker=${encodeURIComponent(emiten.ticker)}`}
                      className="font-mono font-bold text-[var(--md-sys-color-primary)] hover:underline"
                    >
                      {emiten.ticker}
                    </Link>
                    <StockBadges
                      badges={emiten.badges}
                      badge={emiten.badge}
                      conglomerate={emiten.conglomerate || emiten.konglo}
                      size="sm"
                    />
                  </div>
                  <div className="text-[10px] font-sans text-[var(--md-sys-color-on-surface-variant)] max-w-[220px] truncate">
                    {emiten.issuer}
                  </div>
                </td>
                <td className="px-3 py-2.5 text-right">
                  <div className="text-[var(--md-sys-color-on-surface)]">
                    {emiten.foreign_named_trn != null ? `Rp ${fmtNumber(emiten.foreign_named_trn, 1)} T` : "—"}
                  </div>
                  <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-bold">
                    {emiten.foreign_named_pct != null ? `${emiten.foreign_named_pct.toFixed(2)}%` : "—"}
                  </div>
                </td>
                <td className="px-3 py-2.5 text-right">
                  <div className="text-[var(--md-sys-color-on-surface)]">
                    {emiten.foreign_hidden_trn != null ? `Rp ${fmtNumber(emiten.foreign_hidden_trn, 1)} T` : "—"}
                  </div>
                  <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-bold">
                    {emiten.foreign_hidden_pct != null ? `${emiten.foreign_hidden_pct.toFixed(2)}%` : "—"}
                  </div>
                </td>
                <td className="px-3 py-2.5 text-right">
                  <div className="text-[var(--md-sys-color-on-surface)]">
                    {emiten.local_nonret_trn != null ? `Rp ${fmtNumber(emiten.local_nonret_trn, 1)} T` : "—"}
                  </div>
                  <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-bold">
                    {emiten.local_nonret_pct != null ? `${emiten.local_nonret_pct.toFixed(2)}%` : "—"}
                  </div>
                </td>
                <td className="px-3 py-2.5 text-right">
                  <div className="text-[var(--md-sys-color-on-surface)]">
                    {emiten.retail_trn != null ? `Rp ${fmtNumber(emiten.retail_trn, 1)} T` : "—"}
                  </div>
                  <div className="text-[9px] text-amber-500 font-bold">
                    {emiten.retail_pct != null ? `${emiten.retail_pct.toFixed(2)}%` : "—"}
                  </div>
                </td>
                <td className="px-3 py-2.5 text-right font-bold text-[var(--md-sys-color-primary)]">
                  {emiten.free_float != null ? `${(emiten.free_float * 100).toFixed(2)}%` : "—"}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Card>
  );
}
