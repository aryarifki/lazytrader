import { Badge, Card, StatCard } from "@/components/research/ui";
import { fmtSignedPct, toneFor } from "@/components/research/format";
import { fmtRpM } from "./investorHelpers";
import type { InvestorDetail } from "./investorTypes";

export function InvestorSummary({ detail }: { detail: InvestorDetail }) {
  const tone = toneFor(detail.summary.gl_rp_m);

  return (
    <>
      <Card className="p-4 sm:p-5">
        <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-3">
          <div className="min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <Badge tone="primary">CT</Badge>
              <Badge>{detail.lf_label || detail.lf || "Tipe investor —"}</Badge>
              <Badge tone="neutral">{detail.holdings.length} saham</Badge>
            </div>
            <h2 className="mt-2 text-lg sm:text-2xl font-extrabold tracking-tight">
              {detail.name}
            </h2>
            <p className="mt-0.5 text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
              Portofolio per {detail.period.month_curr} · harga {detail.period.last_update_date ?? "—"}
            </p>
          </div>
          <div className="text-right shrink-0">
            <div className="text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
              Total Nilai Portofolio
            </div>
            <div className="text-xl font-extrabold text-[var(--md-sys-color-primary)]">
              {fmtRpM(detail.summary.total_value_m)}
            </div>
          </div>
        </div>
      </Card>

      <div className="grid grid-cols-2 lg:grid-cols-5 gap-2.5">
        <StatCard
          label={`G/L ${detail.period.gl_month}`}
          value={fmtRpM(detail.summary.gl_rp_m, true)}
          sub={
            detail.summary.gl_pct !== null
              ? `${fmtSignedPct(detail.summary.gl_pct)} · basis peak-capital`
              : "Tidak tersedia"
          }
          icon="ph:chart-line-up-duotone"
          tone={tone}
        />
        <StatCard
          label={`Buy ${detail.period.month_curr}`}
          value={fmtRpM(detail.summary.buy_value_m)}
          sub={`${detail.summary.buy_count} saham`}
          icon="ph:trend-up-duotone"
          tone="up"
        />
        <StatCard
          label={`Sell ${detail.period.month_curr}`}
          value={fmtRpM(detail.summary.sell_value_m)}
          sub={`${detail.summary.sell_count} saham`}
          icon="ph:trend-down-duotone"
          tone="down"
        />
        <StatCard
          label={`Net ${detail.period.month_curr}`}
          value={fmtRpM(detail.summary.net_value_m, true)}
          sub="Buy - Sell"
          icon="ph:arrows-left-right-duotone"
          tone={toneFor(detail.summary.net_value_m)}
        />
        <StatCard
          label="Jumlah saham"
          value={detail.summary.ticker_count}
          sub="Emiten tercatat >1%"
          icon="ph:buildings-duotone"
        />
      </div>
    </>
  );
}
