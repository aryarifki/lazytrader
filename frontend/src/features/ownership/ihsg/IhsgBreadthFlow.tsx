import { Badge, Card, MiniBars, SectionTitle } from "@/components/research/ui";
import { fmtCompact, fmtDate, toneFor } from "@/components/research/format";
import { pctClass } from "./ihsgHelpers";
import type { MarketOverview } from "./ihsgTypes";

export function IhsgBreadthFlow({ data }: { data: MarketOverview }) {
  return (
    <div className="grid grid-cols-1 xl:grid-cols-[1.05fr_0.95fr] gap-4">
      <Card className="p-4 sm:p-5">
        <SectionTitle
          icon="ph:chart-bar-duotone"
          title="Breadth pasar"
          subtitle="Perubahan harga harian seluruh emiten"
          action={<Badge>{data.breadth.total} emiten</Badge>}
        />
        <div className="mt-4 grid grid-cols-1 sm:grid-cols-[1fr_180px] gap-5 items-center">
          <MiniBars
            items={[
              { label: "Naik", value: data.breadth.advancers, color: "#10B981" },
              { label: "Turun", value: data.breadth.decliners, color: "#F43F5E" },
              { label: "Stagnan", value: data.breadth.unchanged, color: "#9CA3AF" },
            ]}
          />
          <div className="rounded-xl border border-[var(--md-sys-color-outline-variant)]/60 bg-[var(--md-sys-color-surface-container-high)]/45 p-3.5 text-center">
            <div className="text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
              Rasio naik
            </div>
            <div className="text-2xl font-extrabold text-emerald-600 dark:text-emerald-400 mt-1">
              {data.breadth.total > 0
                ? `${((data.breadth.advancers / data.breadth.total) * 100).toFixed(1)}%`
                : "—"}
            </div>
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] mt-1">
              dari {data.breadth.total} emiten
            </div>
          </div>
        </div>
      </Card>

      <Card className="p-4 sm:p-5">
        <SectionTitle
          icon="ph:globe-hemisphere-west-duotone"
          title="Foreign flow pasar"
          subtitle={`Arus neto asing bursa (100% BEI) · ${fmtDate(data.foreign_flow.date)}`}
        />
        <div className="mt-4 grid grid-cols-3 gap-2.5">
          <FlowCard label="1 hari" value={data.foreign_flow["1d"]} />
          <FlowCard label="1 minggu" value={data.foreign_flow["1w"]} />
          <FlowCard label="1 bulan" value={data.foreign_flow["1m"]} />
        </div>
        <p className="mt-3 text-[10px] leading-relaxed text-[var(--md-sys-color-on-surface-variant)]">
          Nilai pasar agregat seluruh emiten (foreign_buy − foreign_sell) × VWAP dari BEI.
          Nilai positif berarti net foreign inflow ke bursa.
        </p>
      </Card>
    </div>
  );
}

function FlowCard({ label, value }: { label: string; value: number | null }) {
  const tone = toneFor(value);
  return (
    <div className="rounded-xl border border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-high)]/45 p-3">
      <div className="text-[8px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
        {label}
      </div>
      <div className={`mt-1 text-sm font-extrabold font-mono ${pctClass(value)}`}>
        {value === null ? "—" : `Rp ${fmtCompact(value)}`}
      </div>
      <div className="mt-0.5 text-[9px] text-[var(--md-sys-color-on-surface-variant)]">
        {tone === "up" ? "net buy" : tone === "down" ? "net sell" : "flat"}
      </div>
    </div>
  );
}
