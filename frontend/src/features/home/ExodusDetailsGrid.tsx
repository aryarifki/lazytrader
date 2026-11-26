import { Icon } from "@iconify/react";
import { fmtRadarFlow, signedColor } from "@/lib/format";
import type { CulpritBroker, MarketState, TimelinePoint } from "./types";

interface ExodusDetailsGridProps {
  culpritBrokers: CulpritBroker[];
  timeline: TimelinePoint[];
}

function renderStateBadge(state?: MarketState) {
  switch (state) {
    case "CRASH":
      return (
        <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-rose-500/10 text-rose-500 border border-rose-500/20">
          CRASH
        </span>
      );
    case "BOTTOMING":
      return (
        <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-sky-500/10 text-sky-500 border border-sky-500/20">
          BOTTOM
        </span>
      );
    case "RECOVERY_TESTING":
      return (
        <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-amber-500/10 text-amber-500 border border-amber-500/20 animate-pulse">
          TEST
        </span>
      );
    case "CONFIRMED_BOTTOM":
      return (
        <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-emerald-600/10 text-emerald-600 border border-emerald-600/20">
          CONFIRM
        </span>
      );
    case "WARNING":
      return (
        <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-amber-500/10 text-amber-500 border border-amber-500/20">
          WARN
        </span>
      );
    default:
      return (
        <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)]">
          NORM
        </span>
      );
  }
}

export default function ExodusDetailsGrid({ culpritBrokers, timeline }: ExodusDetailsGridProps) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-5 pt-2">
      <div className="lg:col-span-5 bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 sm:p-5 shadow-2xs space-y-3">
        <div className="flex items-center justify-between border-b border-[var(--md-sys-color-outline-variant)] pb-3">
          <div className="flex items-center gap-2">
            <Icon icon="ph:skull-duotone" width="18" className="text-rose-500" />
            <h4 className="text-xs font-extrabold uppercase tracking-wider text-[var(--md-sys-color-on-surface)]">
              Top 5 Culprit Brokers (14D)
            </h4>
          </div>
          <span className="text-[10px] font-mono font-medium text-[var(--md-sys-color-on-surface-variant)]">
            Net Sell Asing
          </span>
        </div>

        {culpritBrokers.length === 0 ? (
          <div className="py-8 text-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
            Tidak ada data broker asing yang terdeteksi net sell masif pada periode ini.
          </div>
        ) : (
          <div className="space-y-2">
            {culpritBrokers.map((broker, idx) => (
              <div
                key={broker.broker_code}
                className="flex items-center justify-between p-2.5 rounded-xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 text-xs font-mono hover:bg-[var(--md-sys-color-surface-container-high)] transition-colors"
              >
                <div className="flex items-center gap-2.5 min-w-0">
                  <span className="w-5 h-5 rounded-md bg-rose-500/10 text-rose-500 font-extrabold flex items-center justify-center text-[10px] shrink-0">
                    #{idx + 1}
                  </span>
                  <div className="min-w-0 truncate">
                    <span className="font-extrabold text-[var(--md-sys-color-on-surface)] mr-2">
                      {broker.broker_code}
                    </span>
                    <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] truncate" title={broker.broker_name}>
                      {broker.broker_name}
                    </span>
                  </div>
                </div>

                <div className="font-bold text-rose-600 dark:text-rose-400 tabular-nums shrink-0 ml-2">
                  {fmtRadarFlow(broker.net_sell_idr)}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="lg:col-span-7 bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 sm:p-5 shadow-2xs space-y-3">
        <div className="flex items-center justify-between border-b border-[var(--md-sys-color-outline-variant)] pb-3">
          <div className="flex items-center gap-2">
            <Icon icon="ph:table-duotone" width="18" className="text-[var(--md-sys-color-primary)]" />
            <h4 className="text-xs font-extrabold uppercase tracking-wider text-[var(--md-sys-color-on-surface)]">
              Timeline Arus Asing Harian
            </h4>
          </div>
          <span className="text-[10px] font-mono font-medium text-[var(--md-sys-color-on-surface-variant)]">
            {timeline.length} Hari Perdagangan
          </span>
        </div>

        <div className="max-h-[220px] overflow-y-auto scrollbar-thin scrollbar-thumb-[var(--md-sys-color-outline-variant)]">
          <table className="w-full text-left text-xs whitespace-nowrap font-mono">
            <thead className="sticky top-0 bg-[var(--md-sys-color-surface-container-high)] text-[10px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]">
              <tr>
                <th className="py-2 px-3">Tanggal</th>
                <th className="py-2 px-2 text-center">Status</th>
                <th className="py-2 px-3 text-right">IHSG</th>
                <th className="py-2 px-3 text-right">Net Flow Asing</th>
                <th className="py-2 px-3 text-right">Kumulatif</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40 text-[11px]">
              {[...timeline].reverse().map((row) => (
                <tr key={row.date} className="hover:bg-[var(--md-sys-color-surface-container)] transition-colors">
                  <td className="py-2 px-3 font-medium text-[var(--md-sys-color-on-surface)]">{row.date}</td>
                  <td className="py-2 px-2 text-center">{renderStateBadge(row.market_state)}</td>
                  <td className="py-2 px-3 text-right font-semibold text-[#38bdf8]">
                    {row.ihsg_close.toLocaleString("id-ID")}
                  </td>
                  <td className="py-2 px-3 text-right font-bold tabular-nums" style={{ color: signedColor(row.daily_flow) }}>
                    {row.daily_flow > 0 ? "+" : ""}
                    {fmtRadarFlow(row.daily_flow)}
                  </td>
                  <td
                    className="py-2 px-3 text-right font-bold tabular-nums"
                    style={{ color: signedColor(row.cumulative_flow) }}
                  >
                    {row.cumulative_flow > 0 ? "+" : ""}
                    {fmtRadarFlow(row.cumulative_flow)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
