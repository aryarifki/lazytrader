import { Badge, Card, Donut, EmptyState, SectionTitle } from "@/components/research/ui";
import type { Composition, CompositionItem } from "./stockTypes";

export function StockCompositionCard({
  displayComposition,
  sortedComposition,
}: {
  displayComposition: Composition | null;
  sortedComposition: CompositionItem[];
}) {
  const largest = sortedComposition[0];

  return (
    <Card className="p-4 sm:p-5">
      <SectionTitle
        icon="ph:chart-donut-duotone"
        title="Komposisi kepemilikan"
        subtitle={displayComposition?.period_label ?? "Belum tersedia"}
        action={
          <Badge
            tone={
              displayComposition?.basis === "staging_normalized"
                ? "up"
                : displayComposition?.basis === "ct_baseline_pct_history"
                  ? "neutral"
                  : displayComposition?.basis === "monthly_ownership_metrics"
                    ? "primary"
                    : "warning"
            }
          >
            {displayComposition?.basis === "staging_normalized"
              ? "KSEI / IDX (Normalized)"
              : displayComposition?.basis === "ct_baseline_pct_history"
                ? "Baseline CT (Quarantined)"
                : displayComposition?.basis === "monthly_ownership_metrics"
                  ? "KSEI Bulanan"
                  : "Snapshot"}
          </Badge>
        }
      />
      <div className="mt-4">
        {sortedComposition.length > 0 ? (
          <Donut
            items={sortedComposition}
            centerLabel={largest?.label ?? "Total"}
            centerValue={`${(largest?.pct ?? 0).toFixed(1)}%`}
          />
        ) : (
          <EmptyState
            icon="ph:chart-donut-duotone"
            title="Komposisi belum tersedia"
            description="Tidak ada baris pemegang saham yang dapat diklasifikasikan."
          />
        )}
      </div>
      <p className="mt-3 text-[10px] leading-relaxed text-[var(--md-sys-color-on-surface-variant)]">
        {displayComposition?.note ?? "Komposisi kepemilikan belum tersedia."}
      </p>
    </Card>
  );
}
