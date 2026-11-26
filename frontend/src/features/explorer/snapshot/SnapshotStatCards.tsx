import { StatCard } from "@/components/research/ui";
import { fmtNumber } from "@/components/research/format";

export function SnapshotStatCards({
  total,
  shown,
  activeFilterCount,
  source,
}: {
  total: number;
  shown: number;
  activeFilterCount: number;
  source: string | undefined;
}) {
  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-2.5">
      <StatCard
        label="Hasil"
        value={fmtNumber(total)}
        sub={`${fmtNumber(shown)} ditampilkan`}
        icon="ph:list-magnifying-glass-duotone"
        tone="primary"
      />
      <StatCard
        label="Filter aktif"
        value={activeFilterCount}
        sub={activeFilterCount ? "Klik chip untuk menghapus" : "Belum ada filter tambahan"}
        icon="ph:funnel-duotone"
      />
      <StatCard
        label="Sumber"
        value={source?.split("+")[0]?.trim() ?? "—"}
        sub="Tabel read-only"
        icon="ph:database-duotone"
      />
      <StatCard
        label="Delta / T5"
        value="Belum tersedia"
        sub="Butuh periodized holdings"
        icon="ph:warning-duotone"
        tone="down"
      />
    </div>
  );
}
