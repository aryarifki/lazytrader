import { lazy, Suspense } from "react";
import type { EChartsOption } from "echarts";
import { Badge, Card, EmptyState, ErrorBlock, LoadingBlock, SectionTitle } from "@/components/research/ui";
import type { NetworkData } from "./groupTypes";

const ReactECharts = lazy(() => import("echarts-for-react"));

interface GroupNetworkViewProps {
  network: NetworkData | undefined;
  networkError: Error | undefined;
  networkLoading: boolean;
  minPct: number;
  graphOption: EChartsOption;
}

export function GroupNetworkView({
  network,
  networkError,
  networkLoading,
  minPct,
  graphOption,
}: GroupNetworkViewProps) {
  return (
    <Card className="p-4 sm:p-5">
      <SectionTitle
        icon="ph:share-network-duotone"
        title="Jejaring kepemilikan"
        subtitle={
          network
            ? `${network.stats.total_nodes} simpul · ${network.stats.total_links} relasi · min ${minPct}%`
            : "Memuat jejaring..."
        }
        action={
          network ? <Badge tone="primary">{network.stats.total_emitens} emiten</Badge> : null
        }
      />
      <div className="mt-3">
        {networkLoading ? <LoadingBlock label="Memuat graph kepemilikan..." /> : null}
        {networkError ? <ErrorBlock message={`Gagal memuat graph: ${networkError.message}`} /> : null}
        {network && network.nodes.length > 0 ? (
          <div className="rounded-xl border border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-low)]/40 overflow-hidden">
            <Suspense fallback={<LoadingBlock label="Memuat graph kepemilikan..." />}>
              <ReactECharts option={graphOption} style={{ height: 540, width: "100%" }} notMerge lazyUpdate />
            </Suspense>
          </div>
        ) : !networkLoading && !networkError ? (
          <EmptyState
            icon="ph:share-network-duotone"
            title="Tidak ada relasi pada ambang ini"
            description="Turunkan minimum porsi atau pilih grup lain untuk melihat jejaring."
          />
        ) : null}
      </div>
      <p className="mt-3 text-[10px] leading-relaxed text-[var(--md-sys-color-on-surface-variant)]">
        Node emiten dan pemegang saham berasal dari snapshot company_shareholders. Garis
        menunjukkan porsi kepemilikan; garis tebal/oranye menandai status pengendali.
      </p>
    </Card>
  );
}
