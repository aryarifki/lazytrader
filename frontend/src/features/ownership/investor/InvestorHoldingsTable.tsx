import { Link } from "react-router-dom";

import { Badge, Card, SectionTitle } from "@/components/research/ui";
import { fmtPct } from "@/components/research/format";
import { StockBadges } from "@/components/research/StockBadges";
import { fmtDelta, fmtLot, fmtValueM } from "./investorHelpers";
import type { InvestorDetail } from "./investorTypes";

export function InvestorHoldingsTable({ detail }: { detail: InvestorDetail }) {
  return (
    <Card className="overflow-hidden">
      <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <SectionTitle
          icon="ph:briefcase-duotone"
          title="Portofolio Saham"
          subtitle={`${detail.period.month_prev} → ${detail.period.month_curr} · nilai memakai harga ${detail.period.last_update_date ?? "terakhir"}`}
          action={
            <div className="flex items-center gap-1.5">
              <Badge tone="neutral">CT Aliased Baseline</Badge>
              <Badge>{detail.holdings.length} baris</Badge>
            </div>
          }
        />
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left text-[11px]">
          <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
            <tr>
              <th className="px-3 py-2.5 font-bold">#</th>
              <th className="px-3 py-2.5 font-bold">Kode / Emiten</th>
              <th className="px-3 py-2.5 font-bold text-right">(Lot) {detail.period.month_prev}</th>
              <th className="px-3 py-2.5 font-bold text-right">% {detail.period.month_prev}</th>
              <th className="px-3 py-2.5 font-bold text-right">(Lot) {detail.period.month_curr}</th>
              <th className="px-3 py-2.5 font-bold text-right">% {detail.period.month_curr}</th>
              <th className="px-3 py-2.5 font-bold text-right">Val (M)</th>
              <th className="px-3 py-2.5 font-bold text-right">Change</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
            {detail.holdings.map((row, index) => (
              <tr
                key={row.ticker}
                className={`hover:bg-[var(--md-sys-color-surface-container-high)]/40 ${
                  row.status === "out" ? "opacity-55" : ""
                }`}
              >
                <td className="px-3 py-2.5 font-mono text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
                  {index + 1}
                </td>
                <td className="px-3 py-2.5">
                  <div className="flex items-center gap-1.5 flex-wrap">
                    <Link
                      to={`/ownership?tab=saham&ticker=${encodeURIComponent(row.ticker)}`}
                      className="font-mono font-bold text-[var(--md-sys-color-primary)] hover:underline"
                    >
                      {row.ticker}
                    </Link>
                    <StockBadges
                      badges={row.badges}
                      badge={row.badge}
                      conglomerate={row.conglomerate || row.konglo}
                      size="sm"
                    />
                  </div>
                  <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] max-w-[260px] line-clamp-1">
                    {row.issuer}
                  </div>
                  {row.is_treasury ? (
                    <div className="mt-1">
                      <Badge tone="warning">Treasury Stock</Badge>
                    </div>
                  ) : null}
                </td>
                <td className="px-3 py-2.5 text-right font-mono">
                  {row.status === "new" ? "—" : fmtLot(row.sl_mei)}
                </td>
                <td className="px-3 py-2.5 text-right font-mono">
                  {row.status === "new" ? "—" : fmtPct(row.pct_mei)}
                </td>
                <td className="px-3 py-2.5 text-right font-mono">
                  {row.status === "out" ? (
                    <Badge tone="down">OUT</Badge>
                  ) : (
                    fmtLot(row.sl_jun)
                  )}
                </td>
                <td className="px-3 py-2.5 text-right font-mono font-bold">
                  {row.status === "out" ? "—" : fmtPct(row.pct_jun)}
                </td>
                <td className="px-3 py-2.5 text-right font-mono">
                  {row.status === "out" ? "—" : fmtValueM(row.value_m)}
                </td>
                <td className="px-3 py-2.5 text-right font-mono font-bold">
                  {fmtDelta(row)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="p-4 border-t border-[var(--md-sys-color-outline-variant)]/60">
        <div className="flex flex-wrap items-center gap-2 text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
          <span>Tanggal Last Update: {detail.period.last_update_date ?? "—"}</span>
          <span className="opacity-40">•</span>
          <span>Saham Scrip tidak ditampilkan</span>
        </div>
      </div>
    </Card>
  );
}
