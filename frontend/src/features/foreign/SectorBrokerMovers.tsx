import { Icon } from "@iconify/react";
import { fmtRadarFlow } from "@/lib/format";
import type { BrokerMover } from "./sectorTypes";

interface SectorBrokerMoversProps {
  saviorBrokers?: BrokerMover[];
  culpritBrokers?: BrokerMover[];
}

export default function SectorBrokerMovers({
  saviorBrokers,
  culpritBrokers,
}: SectorBrokerMoversProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-3.5">
      {/* Savior Brokers (Top Net Buy Asing) */}
      <div className="p-4 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-emerald-500/30 space-y-2.5">
        <div className="flex items-center justify-between text-xs font-bold text-emerald-600 dark:text-emerald-400">
          <span className="flex items-center gap-1.5">
            <Icon icon="ph:arrow-fat-line-up-fill" width="16" />
            Top Pembeli Asing (Savior Brokers - 14 Hari)
          </span>
          <span className="text-[10px] font-normal text-[var(--md-sys-color-on-surface-variant)]">
            Akumulasi Terbanyak
          </span>
        </div>
        <div className="space-y-1.5">
          {saviorBrokers && saviorBrokers.length > 0 ? (
            saviorBrokers.map((b) => (
              <div
                key={b.broker_code}
                className="flex items-center justify-between p-2 rounded-lg bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/50 text-xs font-mono"
              >
                <div className="flex items-center gap-2">
                  <span className="font-extrabold text-emerald-500 bg-emerald-500/10 px-2 py-0.5 rounded-md">
                    {b.broker_code}
                  </span>
                  <span className="text-[var(--md-sys-color-on-surface)] font-sans truncate max-w-[140px]">
                    {b.broker_name}
                  </span>
                </div>
                <span className="font-bold text-emerald-500">+{fmtRadarFlow(b.net_value)}</span>
              </div>
            ))
          ) : (
            <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] block text-center py-2">
              Belum ada akumulasi asing dominan.
            </span>
          )}
        </div>
      </div>

      {/* Culprit Brokers (Top Net Sell Asing) */}
      <div className="p-4 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-rose-500/30 space-y-2.5">
        <div className="flex items-center justify-between text-xs font-bold text-rose-600 dark:text-rose-400">
          <span className="flex items-center gap-1.5">
            <Icon icon="ph:arrow-fat-line-down-fill" width="16" />
            Top Penjual Asing (Culprit Brokers - 14 Hari)
          </span>
          <span className="text-[10px] font-normal text-[var(--md-sys-color-on-surface-variant)]">
            Distribusi Terbanyak
          </span>
        </div>
        <div className="space-y-1.5">
          {culpritBrokers && culpritBrokers.length > 0 ? (
            culpritBrokers.map((b) => (
              <div
                key={b.broker_code}
                className="flex items-center justify-between p-2 rounded-lg bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/50 text-xs font-mono"
              >
                <div className="flex items-center gap-2">
                  <span className="font-extrabold text-rose-500 bg-rose-500/10 px-2 py-0.5 rounded-md">
                    {b.broker_code}
                  </span>
                  <span className="text-[var(--md-sys-color-on-surface)] font-sans truncate max-w-[140px]">
                    {b.broker_name}
                  </span>
                </div>
                <span className="font-bold text-rose-500">{fmtRadarFlow(b.net_value)}</span>
              </div>
            ))
          ) : (
            <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] block text-center py-2">
              Belum ada distribusi asing dominan.
            </span>
          )}
        </div>
      </div>
    </div>
  );
}
