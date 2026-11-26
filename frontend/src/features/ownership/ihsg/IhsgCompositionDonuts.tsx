import { Card, Donut, SectionTitle } from "@/components/research/ui";
import { fmtNumber } from "@/components/research/format";
import type { IhsgComposition } from "./ihsgTypes";

export function IhsgCompositionDonuts({ comp }: { comp: IhsgComposition }) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      {/* Donut 1: Scriptless vs Warkat */}
      <Card className="p-4 sm:p-5 flex flex-col justify-between">
        <div>
          <SectionTitle
            icon="ph:scroll-duotone"
            title="Scriptless vs Warkat"
            subtitle="Proporsi saham tercatat di KSEI"
          />
          <div className="mt-4 flex justify-center">
            <Donut
              items={comp.donuts.scripless_vs_warkat.map((i) => ({
                label: `${i.label} (${fmtNumber(i.value_trn, 1)} T)`,
                pct: i.pct,
                color: i.color,
              }))}
              centerLabel="Market Cap"
              centerValue={`Rp ${fmtNumber(comp.total_mktcap_trn, 1)} T`}
              size={150}
            />
          </div>
        </div>
        <div className="mt-3 pt-3 border-t border-[var(--md-sys-color-outline-variant)]/40 text-[10px] text-[var(--md-sys-color-on-surface-variant)] flex justify-between">
          <span>Warkat (Non-KSEI):</span>
          <span className="font-mono font-bold text-[var(--md-sys-color-on-surface)]">
            Rp {fmtNumber(comp.warkat_trn, 1)} T
          </span>
        </div>
      </Card>

      {/* Donut 2: Investor Type Composition */}
      <Card className="p-4 sm:p-5 flex flex-col justify-between">
        <div>
          <SectionTitle
            icon="ph:users-three-duotone"
            title="Tipe Kepemilikan Investor"
            subtitle="Asing, Lokal Institusi, dan Ritel"
          />
          <div className="mt-4 flex justify-center">
            <Donut
              items={comp.donuts.investor_type.map((i) => ({
                label: `${i.label} (${fmtNumber(i.value_trn, 1)} T)`,
                pct: i.pct,
                color: i.color,
              }))}
              centerLabel="Scriptless"
              centerValue={`Rp ${fmtNumber(comp.months?.jun?.total_trn ?? 6864.8, 1)} T`}
              size={150}
            />
          </div>
        </div>
        <div className="mt-3 pt-3 border-t border-[var(--md-sys-color-outline-variant)]/40 text-[10px] text-[var(--md-sys-color-on-surface-variant)] flex justify-between">
          <span>Ritel Nasional:</span>
          <span className="font-mono font-bold text-amber-500">
            {comp.donuts.investor_type.find((i) => i.label === "Retail")?.pct ?? 12.22}%
          </span>
        </div>
      </Card>

      {/* Donut 3: Foreign Ownership Distribution */}
      <Card className="p-4 sm:p-5 flex flex-col justify-between">
        <div>
          <SectionTitle
            icon="ph:globe-hemisphere-east-duotone"
            title="Distribusi Saham Asing"
            subtitle="Konsentrasi porsi kepemilikan asing"
          />
          <div className="mt-4 flex justify-center">
            <Donut
              items={comp.donuts.foreign_distribution.map((i) => ({
                label: `${i.label} (${fmtNumber(i.count)} emiten)`,
                pct: i.pct,
                color: i.color,
              }))}
              centerLabel="Total Emiten"
              centerValue={`${fmtNumber(comp.n_tickers)} Saham`}
              size={150}
            />
          </div>
        </div>
        <div className="mt-3 pt-3 border-t border-[var(--md-sys-color-outline-variant)]/40 text-[10px] text-[var(--md-sys-color-on-surface-variant)] flex justify-between">
          <span>Mayoritas Asing (≥50%):</span>
          <span className="font-mono font-bold text-purple-400">
            {fmtNumber(comp.n_foreign_ge50)} Emiten
          </span>
        </div>
      </Card>
    </div>
  );
}
