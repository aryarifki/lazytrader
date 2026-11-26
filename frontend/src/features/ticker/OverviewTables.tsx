import { Icon } from "@iconify/react";
import { fmtRp, fmtPct, signedColor, getProfileTextColor } from "@/lib/format";
import type { BrokerSummaryRow, PricePerformanceRow, ProfileBrokerDetailRow } from "./types";

export function TopBrokerSummaryCard({ rows }: { rows: BrokerSummaryRow[] }) {
  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-5 shadow-xs">
      <h3 className="text-sm font-bold text-[var(--md-sys-color-on-surface)] mb-1 tracking-tight">Top Broker Summary</h3>
      <p className="text-[10px] font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider mb-3">
        Net buy/sell tanggal analisis
      </p>
      <div className="overflow-x-auto scrollbar-thin">
        <table className="w-full text-xs">
          <thead>
            <tr className="text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)] text-[10px] font-mono uppercase tracking-wider font-bold">
              <th className="text-left py-2">Side</th>
              <th className="text-left py-2">Broker</th>
              <th className="text-left py-2">Type</th>
              <th className="text-right py-2">Net</th>
              <th className="text-left py-2 pl-3">5D</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row, i) => (
              <tr
                key={i}
                className="border-b border-[var(--md-sys-color-outline-variant)]/60 hover:bg-[var(--md-sys-color-surface-container-high)] transition-colors"
              >
                <td
                  className="py-2.5 font-bold"
                  style={{ color: row.side === "Buy" ? "var(--color-positive)" : "var(--color-negative)" }}
                >
                  {row.side}
                </td>
                <td className="py-2.5 text-[var(--md-sys-color-on-surface)] font-extrabold font-mono">
                  <div className="flex items-center gap-1.5">
                    <span>{row.broker}</span>
                    {row.is_affiliated && (
                      <span
                        title={row.affiliation?.relation || "Broker terafiliasi dengan emiten"}
                        className="inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded-sm text-[9px] font-bold bg-amber-500 text-black uppercase tracking-wider shadow-2xs cursor-help"
                      >
                        <Icon icon="ph:buildings-bold" width="10" />
                        <span>Afiliasi</span>
                      </span>
                    )}
                  </div>
                </td>
                <td className="py-2.5 text-[var(--md-sys-color-on-surface-variant)] font-medium">{row.type}</td>
                <td className="py-2.5 text-right font-mono font-bold" style={{ color: signedColor(row.net) }}>
                  {fmtRp(row.net)}
                </td>
                <td className="py-2.5 pl-3 text-[var(--md-sys-color-on-surface-variant)] font-mono font-bold tracking-widest">
                  {row.spark}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export function PricePerformanceCard({ rows }: { rows: PricePerformanceRow[] }) {
  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-5 shadow-xs">
      <h3 className="text-sm font-bold text-[var(--md-sys-color-on-surface)] mb-2 tracking-tight">Price Performance</h3>
      <div className="overflow-x-auto">
        <table className="w-full text-xs">
          <thead>
            <tr className="text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)] text-[10px] font-mono uppercase tracking-wider font-bold">
              <th className="text-left py-2">Period</th>
              <th className="text-right py-2">Return</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row, i) => (
              <tr
                key={i}
                className="border-b border-[var(--md-sys-color-outline-variant)]/60 hover:bg-[var(--md-sys-color-surface-container-high)] transition-colors"
              >
                <td className="py-2 text-[var(--md-sys-color-on-surface)] font-medium">{row.period}</td>
                <td className="py-2 text-right font-mono font-bold" style={{ color: signedColor(row.value) }}>
                  {fmtPct(row.value)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export function ProfileBrokerDetailCard({ rows }: { rows: ProfileBrokerDetailRow[] }) {
  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-5 shadow-xs">
      <h3 className="text-sm sm:text-base font-bold text-[var(--md-sys-color-on-surface)] mb-4 tracking-tight">
        Broker Detail by Profile Breakdown
      </h3>
      {rows.length === 0 ? (
        <p className="text-xs font-mono font-medium text-[var(--md-sys-color-on-surface-variant)]">
          Tidak ada rincian broker untuk periode ini.
        </p>
      ) : (
        <div className="overflow-x-auto scrollbar-thin">
          <table className="w-full text-xs">
            <thead>
              <tr className="text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)] text-[10px] font-mono uppercase tracking-wider font-bold">
                <th className="text-left py-2 pr-3">Profile</th>
                <th className="text-left py-2 pr-3">Broker</th>
                <th className="text-left py-2 pr-3">Type</th>
                <th className="text-right py-2 pr-3">Buy</th>
                <th className="text-right py-2 pr-3">Sell</th>
                <th className="text-right py-2 pr-3">Net</th>
                <th className="text-right py-2 pr-3">Freq</th>
                <th className="text-right py-2 pr-3">Days</th>
                <th className="text-right py-2">Avg/Tx</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row, i) => (
                <tr
                  key={i}
                  className="border-b border-[var(--md-sys-color-outline-variant)]/60 hover:bg-[var(--md-sys-color-surface-container-high)] transition-colors"
                >
                  <td className="py-2.5 pr-3 font-bold uppercase" style={{ color: getProfileTextColor(row.profile) }}>
                    {row.profile}
                  </td>
                  <td className="py-2.5 pr-3 text-[var(--md-sys-color-on-surface)] font-extrabold font-mono">{row.broker}</td>
                  <td className="py-2.5 pr-3 text-[var(--md-sys-color-on-surface-variant)] font-medium">{row.type}</td>
                  <td className="py-2.5 pr-3 text-right font-mono font-bold" style={{ color: "var(--color-positive)" }}>
                    {fmtRp(row.buy)}
                  </td>
                  <td className="py-2.5 pr-3 text-right font-mono font-bold" style={{ color: "var(--color-negative)" }}>
                    {fmtRp(row.sell)}
                  </td>
                  <td className="py-2.5 pr-3 text-right font-mono font-bold" style={{ color: signedColor(row.net) }}>
                    {fmtRp(row.net)}
                  </td>
                  <td className="py-2.5 pr-3 text-right font-mono font-bold text-[var(--md-sys-color-on-surface-variant)]">
                    {row.freq?.toLocaleString("id-ID") || 0}
                  </td>
                  <td className="py-2.5 pr-3 text-right font-mono font-bold text-[var(--md-sys-color-on-surface-variant)]">
                    {row.days || 0}
                  </td>
                  <td className="py-2.5 text-right font-mono font-bold text-[var(--md-sys-color-on-surface-variant)]">
                    {fmtRp(row.avg_value_tx)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
