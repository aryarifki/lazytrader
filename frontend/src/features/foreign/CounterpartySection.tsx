import React, { Suspense } from "react";
import { Icon } from "@iconify/react";
import { formatIdrCompact } from "./deepdiveHelpers";
import type { CounterpartyInfo, ImpactDivergenceInfo } from "./deepdiveTypes";
import type { HeatmapData } from "./charts/BrokerHeatmap";

const CounterpartyFlowChart = React.lazy(() => import("./charts/CounterpartyFlowChart"));
const BrokerHeatmap = React.lazy(() => import("./charts/BrokerHeatmap"));

interface CounterpartySectionProps {
  cp?: CounterpartyInfo;
  imp?: ImpactDivergenceInfo;
  hhiValue: number;
  brokerHeatmap?: HeatmapData | null;
}

export default function CounterpartySection({
  cp,
  imp,
  hhiValue,
  brokerHeatmap,
}: CounterpartySectionProps) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-5">
      <div className="bg-[var(--md-sys-color-surface-container)] p-4 sm:p-5 rounded-2xl border border-[var(--md-sys-color-outline-variant)] shadow-xs flex flex-col justify-between">
        <div>
          <div className="flex flex-wrap items-center justify-between gap-2 pb-3 mb-3 border-b border-[var(--md-sys-color-outline-variant)]/60">
            <div className="flex items-center gap-2">
              <Icon icon="ph:users-three-duotone" className="text-[var(--md-sys-color-primary)]" width="20" />
              <h3 className="text-xs sm:text-sm font-extrabold text-[var(--md-sys-color-on-surface)] uppercase tracking-tight">
                Tarik-Ulur: Asing vs Institusi Lokal vs Ritel
              </h3>
            </div>
          </div>
          <div className="grid grid-cols-3 gap-2 mb-3 text-[11px]">
            <div className="p-2 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]/50">
              <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">Asing</div>
              <div className={`font-mono font-extrabold ${(cp?.totals?.foreign || 0) >= 0 ? "text-emerald-500" : "text-rose-500"}`}>
                {formatIdrCompact(cp?.totals?.foreign || 0)}
              </div>
            </div>
            <div className="p-2 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]/50">
              <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">Inst. Lokal</div>
              <div className={`font-mono font-extrabold ${(cp?.totals?.local_inst || 0) >= 0 ? "text-indigo-400" : "text-rose-500"}`}>
                {formatIdrCompact(cp?.totals?.local_inst || 0)}
              </div>
            </div>
            <div className="p-2 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]/50">
              <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">Ritel</div>
              <div className={`font-mono font-extrabold ${(cp?.totals?.retail || 0) >= 0 ? "text-amber-500" : "text-emerald-500"}`}>
                {formatIdrCompact(cp?.totals?.retail || 0)}
              </div>
            </div>
          </div>
          <div className="w-full h-[250px] relative z-0">
            <Suspense fallback={<div className="h-full w-full" />}>
              <CounterpartyFlowChart data={cp} />
            </Suspense>
          </div>
        </div>
        <p className="mt-3 pt-2.5 border-t border-[var(--md-sys-color-outline-variant)]/50 text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
          Akumulasi paling sehat terjadi ketika kurva <span className="text-emerald-500 font-bold">Asing</span> menanjak sementara <span className="text-rose-500 font-bold">Ritel</span> menurun (distribusi barang dari ritel ke big money).
        </p>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container)] p-4 sm:p-5 rounded-2xl border border-[var(--md-sys-color-outline-variant)] shadow-xs flex flex-col justify-between">
        <div>
          <div className="flex items-center justify-between pb-3 mb-3 border-b border-[var(--md-sys-color-outline-variant)]/60">
            <div className="flex items-center gap-2">
              <Icon icon="ph:squares-four-bold" className="text-[var(--md-sys-color-primary)]" width="20" />
              <h3 className="text-xs sm:text-sm font-extrabold text-[var(--md-sys-color-on-surface)] uppercase tracking-tight">
                Heatmap Broker Asing &amp; Konsistensi Arus
              </h3>
            </div>
            <span className="text-[10px] font-mono font-semibold text-[var(--md-sys-color-on-surface-variant)]">
              HHI {hhiValue.toFixed(3)} · {imp?.buy_days || 0}B / {imp?.sell_days || 0}S ({imp?.consistency_pct || 0}%)
            </span>
          </div>
          <div className="w-full h-[305px] relative z-0">
            <Suspense fallback={<div className="h-full w-full" />}>
              <BrokerHeatmap heatmapData={brokerHeatmap} />
            </Suspense>
          </div>
        </div>
        <p className="mt-3 pt-2.5 border-t border-[var(--md-sys-color-outline-variant)]/50 text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
          Menunjukkan konsistensi harian top 15 broker asing. HHI &gt; 0,4 menandakan akumulasi/distribusi terpusat pada sedikit broker jangkar.
        </p>
      </div>
    </div>
  );
}
