import { Icon } from "@iconify/react";
import type { StockOwnership } from "@/features/ownership/stock/stockTypes";
import { bandarBadgeClass } from "./mmdBadges";
import type { MmdAnalysis, MmdPeriodItem } from "./mmdTypes";

interface MmdBandarOtherTabProps {
  analysisData: MmdAnalysis | null;
  ownership: StockOwnership | null;
  selectedTicker: string;
}

function PeriodRow({ label, items }: { label: string; items: MmdPeriodItem[] }) {
  return (
    <tr className="border-b border-[var(--md-sys-color-outline-variant)]/30 last:border-b-0">
      <td className="py-2 pr-2 font-bold text-[10px] text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider align-top w-14">
        {label}
      </td>
      {items.map((item, idx) => (
        <td
          key={idx}
          className="py-1.5 px-2 align-top text-center border-l border-[var(--md-sys-color-outline-variant)]/20"
        >
          <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface)]">{item.label}</div>
          <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] truncate max-w-[90px] mx-auto">
            {item.range || item.date}
          </div>
          <div className="my-1">
            <span className={`${bandarBadgeClass(item.badge?.cls)} text-[9px] py-0.5 px-1.5`}>
              {item.badge?.label || "—"}
            </span>
          </div>
          <div
            className={`text-[10px] font-bold ${
              item.change_pct >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
            }`}
          >
            {item.change_pct > 0 ? `+${item.change_pct}%` : `${item.change_pct}%`}
          </div>
        </td>
      ))}
    </tr>
  );
}

function formatMarketCap(ownership: StockOwnership | null): string {
  const mcapB = ownership?.profile?.market_cap_b;
  if (mcapB != null && mcapB > 0) {
    return `Rp ${(mcapB / 1000).toFixed(1)} T`;
  }
  const mcap = ownership?.profile?.market_cap;
  if (mcap != null && mcap > 0) {
    return `Rp ${(mcap / 1e12).toFixed(1)} T`;
  }
  return "—";
}

export default function MmdBandarOtherTab({
  analysisData,
  ownership,
  selectedTicker,
}: MmdBandarOtherTabProps) {
  const freeFloat = ownership?.profile?.free_float_pct ?? ownership?.summary?.free_float_est;
  const marketCapStr = formatMarketCap(ownership);

  return (
    <div className="space-y-3 font-mono-nums">
      <div className="bg-[var(--md-sys-color-surface-container-low)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40 space-y-2.5">
        <div className="flex items-center justify-between pb-1 border-b border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-xs font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
            <Icon
              icon="ph:clock-counter-clockwise-duotone"
              width={15}
              className="text-[var(--md-sys-color-primary)]"
            />
            <span>PREVIOUS RESULTS</span>
          </div>
          <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
            6 Periode Historis
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-xs text-left border-collapse font-mono-nums">
            <tbody>
              <PeriodRow label="Weekly" items={analysisData?.previous_results?.weekly || []} />
              <PeriodRow label="Monthly" items={analysisData?.previous_results?.monthly || []} />
            </tbody>
          </table>
        </div>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container-low)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40 space-y-2">
        <div className="text-xs font-bold text-[var(--md-sys-color-on-surface)] pb-1 border-b border-[var(--md-sys-color-outline-variant)]/40">
          PROFIL EMITEN &amp; KEPEMILIKAN KSEI &gt;1%
        </div>
        <div className="space-y-1.5 text-xs">
          <div className="flex justify-between">
            <span className="text-[var(--md-sys-color-on-surface-variant)]">Nama Emiten:</span>
            <span className="font-bold text-[var(--md-sys-color-on-surface)] truncate max-w-[180px]">
              {ownership?.name || `${selectedTicker} Tbk.`}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-[var(--md-sys-color-on-surface-variant)]">Free Float:</span>
            <span className="font-bold text-[var(--color-positive)]">
              {freeFloat != null ? `${Number(freeFloat).toFixed(1)}%` : "—"}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-[var(--md-sys-color-on-surface-variant)]">Market Cap:</span>
            <span className="font-bold text-[var(--md-sys-color-on-surface)]">{marketCapStr}</span>
          </div>
        </div>
      </div>
    </div>
  );
}
