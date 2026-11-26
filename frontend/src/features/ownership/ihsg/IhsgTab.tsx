import { useMemo } from "react";
import useSWR from "swr";
import { Icon } from "@iconify/react";

import {
  Badge,
  EmptyState,
  ErrorBlock,
  LoadingBlock,
  StatCard,
} from "@/components/research/ui";
import {
  fmtCompact,
  fmtNumber,
  fmtSignedPct,
  toneFor,
} from "@/components/research/format";
import { ownershipFetcher } from "../ownershipFetcher";
import { IhsgBreadthFlow } from "./IhsgBreadthFlow";
import { IhsgCompositionDonuts } from "./IhsgCompositionDonuts";
import { IhsgHeaderCard } from "./IhsgHeaderCard";
import { IhsgMutationTable } from "./IhsgMutationTable";
import { IhsgSectorsCoverage } from "./IhsgSectorsCoverage";
import { IhsgTop20Table } from "./IhsgTop20Table";
import type { MarketOverview } from "./ihsgTypes";

export default function IhsgTab() {
  const { data, error, isLoading } = useSWR<MarketOverview>(
    "/api/v1/ownership/market?history_limit=260",
    ownershipFetcher
  );

  const topSectors = useMemo(() => (data?.sectors ?? []).slice(0, 14), [data]);
  const coveragePct = data
    ? Math.min(100, (data.ownership_coverage.covered_tickers / Math.max(1, data.ownership_coverage.total_tickers)) * 100)
    : 0;

  const comp = data?.composition;

  if (isLoading) return <LoadingBlock label="Memuat IHSG dan rotasi sektor..." />;
  if (error) return <ErrorBlock message={`Gagal memuat IHSG: ${error.message}`} />;
  if (!data) return <EmptyState title="Data IHSG tidak tersedia" />;

  return (
    <div className="space-y-5">
      <IhsgHeaderCard data={data} />

      {/* Metric Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-2.5">
        <StatCard
          label="3 bulan"
          value={fmtSignedPct(data.index.return_3m_pct)}
          sub="Return harga"
          icon="ph:chart-line-up-duotone"
          tone={toneFor(data.index.return_3m_pct)}
        />
        <StatCard
          label="6 bulan"
          value={fmtSignedPct(data.index.return_6m_pct)}
          sub="Return harga"
          icon="ph:chart-line-duotone"
          tone={toneFor(data.index.return_6m_pct)}
        />
        <StatCard
          label="52W high / low"
          value={`${fmtNumber(data.index.high_52w, 0)} / ${fmtNumber(data.index.low_52w, 0)}`}
          sub="Batas 52 minggu"
          icon="ph:arrows-down-up-duotone"
        />
        <StatCard
          label="Market cap"
          value={data.index.market_cap ? `Rp ${fmtCompact(data.index.market_cap)}` : "—"}
          sub="Nilai pasar indeks"
          icon="ph:coins-duotone"
          tone="primary"
        />
      </div>

      <IhsgBreadthFlow data={data} />

      {/* TAHAP A: Macro Ownership Composition (3 Donut Cards) */}
      {comp && comp.donuts && (
        <div className="space-y-4">
          <div className="flex items-center justify-between gap-3">
            <div>
              <h3 className="text-sm font-bold tracking-tight text-[var(--md-sys-color-on-surface)] flex items-center gap-2">
                <Icon icon="ph:chart-pie-slice-duotone" className="text-[var(--md-sys-color-primary)] text-base" />
                Komposisi Kepemilikan Makro IHSG
              </h3>
              <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-0.5">
                Struktur kepemilikan saham agregat KSEI di Bursa Efek Indonesia
              </p>
            </div>
            <Badge tone="primary">Agregat Makro KSEI</Badge>
          </div>

          <IhsgCompositionDonuts comp={comp} />

          {comp.summary_comparison && comp.summary_comparison.length > 0 && (
            <IhsgMutationTable rows={comp.summary_comparison} />
          )}

          {comp.top20 && comp.top20.length > 0 && (
            <IhsgTop20Table top20={comp.top20} />
          )}
        </div>
      )}

      <IhsgSectorsCoverage
        data={data}
        topSectors={topSectors}
        coveragePct={coveragePct}
      />
    </div>
  );
}
