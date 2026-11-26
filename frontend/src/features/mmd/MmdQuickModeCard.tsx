import { Icon } from "@iconify/react";
import { bandarBadgeClass } from "./mmdBadges";
import type { MmdAnalysis } from "./mmdTypes";

interface MmdQuickModeCardProps {
  analysisData: MmdAnalysis | null;
}

export default function MmdQuickModeCard({ analysisData }: MmdQuickModeCardProps) {
  const periods = [
    ...(analysisData?.previous_results?.weekly || []),
    ...(analysisData?.previous_results?.monthly || []),
  ];

  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-amber-500/70 dark:border-[#FFC24D] rounded-2xl p-3.5 shadow-lg space-y-2">
      <div className="flex items-center justify-between pb-1.5 border-b border-[var(--md-sys-color-outline-variant)]">
        <span className="text-xs font-bold text-amber-700 dark:text-[#FFC24D] tracking-wide flex items-center gap-1.5">
          <Icon icon="ph:lightning-duotone" width={16} />
          QUICK MODE — RINGKASAN PERIODE
        </span>
        <span className="text-[10px] font-mono text-[var(--md-sys-color-on-surface)]">6 PERIODE</span>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left text-xs font-mono-nums">
          <thead>
            <tr className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]">
              <th className="py-1.5 px-1 text-center">No</th>
              <th className="py-1.5 px-2">Result</th>
              <th className="py-1.5 px-2">Periode</th>
              <th className="py-1.5 px-2 text-right">Change</th>
              <th className="py-1.5 px-2 text-center">Buyer x Seller</th>
              <th className="py-1.5 px-2 text-center">Broker Pengendali</th>
              <th className="py-1.5 px-2 text-center">Retail Action</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
            {periods.map((item, idx) => (
              <tr key={idx} className="hover:bg-[var(--md-sys-color-surface-container-high)] transition">
                <td className="py-1.5 px-1 text-center text-[var(--md-sys-color-outline)] font-bold">
                  {idx + 1}
                </td>
                <td className="py-1.5 px-2">
                  <span
                    className={`${bandarBadgeClass(item.badge?.cls)} text-[10px] py-0.5 px-2 min-w-[75px]`}
                  >
                    {item.badge?.label || "—"}
                  </span>
                </td>
                <td className="py-1.5 px-2 font-bold text-[var(--md-sys-color-on-surface)]">
                  <div>{item.label}</div>
                  <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-normal">
                    {item.range || item.date}
                  </div>
                </td>
                <td
                  className={`py-1.5 px-2 text-right font-bold ${
                    item.change_pct >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                  }`}
                >
                  {item.change_pct > 0 ? `+${item.change_pct}%` : `${item.change_pct}%`}
                </td>
                <td className="py-1.5 px-2 text-center">
                  <span className={`${bandarBadgeClass(item.badge?.cls)} text-[9px] py-0.5 px-1.5`}>
                    {item.badge?.label?.split(" ")?.[0] || "—"}
                  </span>
                </td>
                <td className="py-1.5 px-2 text-center">
                  <span
                    className={`${bandarBadgeClass(item.pengendali_action?.cls)} text-[9px] py-0.5 px-1.5`}
                  >
                    {item.pengendali_action?.label?.split(" ")?.[0] || "—"}
                  </span>
                </td>
                <td className="py-1.5 px-2 text-center">
                  <span className={`${bandarBadgeClass(item.retail_action?.cls)} text-[9px] py-0.5 px-1.5`}>
                    {item.retail_action?.label || "Do Nothing"}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
