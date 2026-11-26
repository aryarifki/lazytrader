import { Card, DataNote, SectionTitle } from "@/components/research/ui";
import type { ExplorerResponse } from "./ctTypes";
import { QualityRow } from "./QualityRow";

export function CTExplorerQuality({
  dataQuality,
}: {
  dataQuality: ExplorerResponse["data_quality"] | undefined;
}) {
  return (
    <Card className="p-4 sm:p-5">
      <SectionTitle
        icon="ph:database-duotone"
        title="Kualitas data"
        subtitle="Sumber dan batas interpretasi"
      />
      <div className="mt-3 grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_280px] gap-3">
        <DataNote>
          <ul className="list-disc pl-4 space-y-1">
            {(dataQuality?.notes ?? [
              "Explorer dihitung server-side dari baseline CT.",
              "G/L memakai peak-capital dan tidak masuk jalur sinyal.",
            ]).map((note) => (
              <li key={note}>{note}</li>
            ))}
          </ul>
        </DataNote>
        <div className="grid grid-cols-2 gap-2 text-[10px] font-mono content-start">
          <QualityRow label="CT baseline" ok={Boolean(dataQuality?.ct_baseline)} />
          <QualityRow label="Periodized" ok={Boolean(dataQuality?.periodized_holdings)} />
          <QualityRow label="T5" ok={Boolean(dataQuality?.t5_transactions)} />
          <QualityRow label="G/L" ok={Boolean(dataQuality?.gl)} />
          <QualityRow label="Market cap" ok={Boolean(dataQuality?.market_cap)} />
          <QualityRow label="Bukan sinyal" ok />
        </div>
      </div>
    </Card>
  );
}
