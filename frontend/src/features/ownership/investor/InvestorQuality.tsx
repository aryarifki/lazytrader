import { Card, DataNote, SectionTitle } from "@/components/research/ui";
import { fmtDate } from "@/components/research/format";
import { QualityRow } from "../QualityRow";
import type { InvestorDetail } from "./investorTypes";

export function InvestorQuality({ detail }: { detail: InvestorDetail }) {
  return (
    <Card className="p-4 sm:p-5">
      <SectionTitle
        icon="ph:database-duotone"
        title="Kualitas data"
        subtitle="Batas interpretasi baseline CT"
      />
      <div className="mt-3 space-y-2.5">
        <DataNote>
          <ul className="list-disc pl-4 space-y-1">
            {detail.data_quality.notes.map((note) => (
              <li key={note}>{note}</li>
            ))}
          </ul>
        </DataNote>
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-2 text-[10px] font-mono">
          <QualityRow label="Periodized holdings" ok={detail.data_quality.periodized_holdings} />
          <QualityRow label="Transaksi T5" ok={detail.data_quality.t5_transactions} />
          <QualityRow label="G/L peak-capital" ok={detail.data_quality.gl} />
          <QualityRow label="Market cap" ok={detail.data_quality.market_cap} />
        </div>
        <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
          Tanggal harga: {detail.period.last_update_date ? fmtDate(detail.period.last_update_date) : "—"}
        </div>
      </div>
    </Card>
  );
}
