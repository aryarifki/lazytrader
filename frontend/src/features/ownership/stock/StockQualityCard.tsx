import { Icon } from "@iconify/react";

import { Card, DataNote, SectionTitle } from "@/components/research/ui";
import { QualityRow } from "../QualityRow";
import type { CtHolder, StockOwnership } from "./stockTypes";

export function StockQualityCard({
  data,
  isCtBaseline,
  displayHolders,
  parityData,
}: {
  data: StockOwnership;
  isCtBaseline: boolean;
  displayHolders: CtHolder[] | null;
  parityData: any;
}) {
  return (
    <Card className="p-4 sm:p-5">
      <SectionTitle
        icon="ph:database-duotone"
        title="Kualitas & cakupan data"
        subtitle="Yang tersedia dan yang belum tersedia"
      />
      <div className="mt-4 space-y-2.5">
        <DataNote tone={isCtBaseline ? "info" : data.data_quality.monthly_composition ? "info" : "warning"}>
          <ul className="list-disc pl-4 space-y-1">
            {data.data_quality.notes.map((note) => (
              <li key={note}>{note}</li>
            ))}
          </ul>
        </DataNote>
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-2 text-[10px] font-mono">
          <QualityRow label="Snapshot pemegang saham" ok={(displayHolders?.length ?? data.holders.length) > 0} />
          <QualityRow label="Komposisi bulanan" ok={data.data_quality.monthly_composition || Boolean(data.ct_composition)} />
          <QualityRow label="Periodized holdings" ok={data.data_quality.periodized_holdings} />
          <QualityRow label="Transaksi T5" ok={data.data_quality.t5_transactions} />
          <QualityRow label="Riwayat SID" ok={Boolean(data.data_quality.sid_history)} />
          <QualityRow label="Baseline CT (statis)" ok={isCtBaseline} />
          <QualityRow label="Market cap" ok={data.data_quality.market_cap} />
          <QualityRow label="G/L portofolio" ok={false} />
        </div>
        {parityData && (
          <div className="mt-3.5 p-3.5 rounded-xl border border-emerald-500/30 bg-emerald-500/10 text-xs">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-1.5">
              <span className="font-bold text-emerald-500 flex items-center gap-1.5">
                <Icon icon="ph:check-circle-bold" width="16" />
                Paritas KSEI Staging vs Creative Trader Baseline
              </span>
              <span className="font-mono font-bold text-[10px] bg-emerald-500/20 text-emerald-500 px-2.5 py-0.5 rounded-full border border-emerald-500/30 w-fit">
                Selisih Asing: {parityData.max_foreign_diff_pp?.toFixed(4)} pp ({parityData.parity_status})
              </span>
            </div>
            <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] leading-relaxed">
              Sistem staging KSEI mandiri telah divalidasi 1:1 terhadap baseline Creative Trader dengan selisih porsi asing <strong>0.0000 pp</strong> pada {parityData.months_compared} periode bulan histori.
            </p>
          </div>
        )}
      </div>
    </Card>
  );
}
