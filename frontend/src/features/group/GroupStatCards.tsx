import { StatCard } from "@/components/research/ui";
import { fmtPct, fmtSignedPct, toneFor } from "@/components/research/format";
import { fmtMarketCapB } from "./groupHelpers";
import type { GroupSummary } from "./groupTypes";

interface GroupStatCardsProps {
  group: GroupSummary | undefined;
}

export function GroupStatCards({ group }: GroupStatCardsProps) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-2.5">
      <StatCard
        label="Rata-rata MTD"
        value={fmtSignedPct(group?.avg_mtd_pct)}
        sub={`${group?.advancers ?? 0} naik · ${group?.decliners ?? 0} turun`}
        icon="ph:trend-up-duotone"
        tone={toneFor(group?.avg_mtd_pct)}
      />
      <StatCard
        label="Emiten tercakup"
        value={`${group?.covered_tickers ?? 0}/${group?.ticker_count ?? 0}`}
        sub="Punya baris pemegang saham"
        icon="ph:buildings-duotone"
        tone="primary"
      />
      <StatCard
        label="Rata-rata publik"
        value={fmtPct(group?.public_pct_avg)}
        sub="Snapshot pemegang saham"
        icon="ph:users-three-duotone"
      />
      <StatCard
        label="Rata-rata pengendali"
        value={fmtPct(group?.controlling_pct_avg)}
        sub="Berdasarkan flag controlling"
        icon="ph:crown-duotone"
      />
      <StatCard
        label="Market cap grup"
        value={fmtMarketCapB(group?.market_cap_b)}
        sub={group?.market_cap_available ? "Baseline CT / harga x shares" : "Belum tersedia"}
        icon="ph:currency-circle-dollar-duotone"
        tone="primary"
      />
    </div>
  );
}
