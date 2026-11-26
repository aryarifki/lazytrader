import { Icon } from "@iconify/react";

import { Badge, Card, DataNote, SectionTitle } from "@/components/research/ui";
import { fmtNumber, fmtSignedPct } from "@/components/research/format";
import { pctClass } from "./ihsgHelpers";
import type { MarketOverview, SectorRow } from "./ihsgTypes";

export function IhsgSectorsCoverage({
  data,
  topSectors,
  coveragePct,
}: {
  data: MarketOverview;
  topSectors: SectorRow[];
  coveragePct: number;
}) {
  return (
    <div className="grid grid-cols-1 xl:grid-cols-[1.15fr_0.85fr] gap-4">
      <Card className="overflow-hidden">
        <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)]/60">
          <SectionTitle
            icon="ph:squares-four-duotone"
            title="Indeks & sektor"
            subtitle="Diurutkan berdasarkan return 1 bulan"
            action={<Badge>{data.sectors.length} indeks</Badge>}
          />
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-[11px]">
            <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
              <tr>
                <th className="px-4 py-2.5 font-bold">Kode</th>
                <th className="px-4 py-2.5 font-bold">Nama</th>
                <th className="px-4 py-2.5 font-bold text-right">Close</th>
                <th className="px-4 py-2.5 font-bold text-right">1D</th>
                <th className="px-4 py-2.5 font-bold text-right">1W</th>
                <th className="px-4 py-2.5 font-bold text-right">1M</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
              {topSectors.map((sector) => (
                <tr key={sector.code} className="hover:bg-[var(--md-sys-color-surface-container-high)]/40">
                  <td className="px-4 py-2.5 font-mono font-bold text-[var(--md-sys-color-primary)]">
                    {sector.code}
                  </td>
                  <td className="px-4 py-2.5 text-[10px] text-[var(--md-sys-color-on-surface-variant)] max-w-[220px] truncate">
                    {sector.name}
                  </td>
                  <td className="px-4 py-2.5 text-right font-mono">{fmtNumber(sector.close, 2)}</td>
                  <td className={`px-4 py-2.5 text-right font-mono font-bold ${pctClass(sector.change_1d_pct)}`}>
                    {fmtSignedPct(sector.change_1d_pct)}
                  </td>
                  <td className={`px-4 py-2.5 text-right font-mono font-bold ${pctClass(sector.change_1w_pct)}`}>
                    {fmtSignedPct(sector.change_1w_pct)}
                  </td>
                  <td className={`px-4 py-2.5 text-right font-mono font-bold ${pctClass(sector.change_1m_pct)}`}>
                    {fmtSignedPct(sector.change_1m_pct)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>

      <div className="space-y-4">
        <Card className="p-4 sm:p-5">
          <SectionTitle
            icon="ph:chart-donut-duotone"
            title="Cakupan Data Saham"
            subtitle="Cakupan emiten KSEI"
          />
          <div className="mt-4">
            <div className="flex items-end justify-between gap-3">
              <div>
                <div className="text-2xl font-extrabold text-[var(--md-sys-color-primary)]">
                  {data.ownership_coverage.covered_tickers}
                </div>
                <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
                  dari {data.ownership_coverage.total_tickers} emiten
                </div>
              </div>
              <Badge tone={data.ownership_coverage.composition_available ? "primary" : "warning"}>
                {data.ownership_coverage.composition_available ? "Aktif" : "Terbatas"}
              </Badge>
            </div>
            <div className="mt-3 h-2 rounded-full bg-[var(--md-sys-color-surface-container-highest)] overflow-hidden">
              <div
                className="h-full rounded-full bg-[var(--md-sys-color-primary)]"
                style={{ width: `${coveragePct}%` }}
              />
            </div>
            <p className="mt-3 text-[10px] leading-relaxed text-[var(--md-sys-color-on-surface-variant)]">
              {data.ownership_coverage.note}
            </p>
          </div>
        </Card>

        <Card className="p-4 sm:p-5">
          <SectionTitle
            icon="ph:database-duotone"
            title="Sumber & batas"
            subtitle="Transparansi data"
          />
          <div className="mt-3 space-y-2">
            <DataNote tone="info">
              <ul className="list-disc pl-4 space-y-1">
                {data.data_quality.notes.map((note) => (
                  <li key={note}>{note}</li>
                ))}
              </ul>
            </DataNote>
            <div className="grid grid-cols-2 gap-2 text-[10px] font-mono">
              <SourceRow label="IHSG" ok />
              <SourceRow label="Breadth" ok />
              <SourceRow label="Foreign flow" ok />
              <SourceRow label="Top 20" ok={data.ownership_coverage.top20_available} />
              <SourceRow label="Komposisi IHSG" ok={data.ownership_coverage.composition_available} />
              <SourceRow label="Market cap emiten" ok />
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}

function SourceRow({ label, ok }: { label: string; ok: boolean }) {
  return (
    <div className="flex items-center gap-1.5 rounded-md border border-[var(--md-sys-color-outline-variant)]/50 px-2 py-1.5">
      <Icon
        icon={ok ? "ph:check-circle-duotone" : "ph:minus-circle-duotone"}
        width="13"
        className={ok ? "text-emerald-500" : "text-[var(--md-sys-color-on-surface-variant)]"}
      />
      <span className="text-[var(--md-sys-color-on-surface-variant)]">{label}</span>
    </div>
  );
}
