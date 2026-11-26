import type React from "react";
import { Icon } from "@iconify/react";
import { fmtRp, fmtRadarFlow, signedColor } from "@/lib/format";
import { getForeignStateBadge } from "./sectorHelpers";
import SectorHistoryChart from "./SectorHistoryChart";
import SectorBrokerMovers from "./SectorBrokerMovers";
import SectorHistoryTable from "./SectorHistoryTable";
import type { SectorHistoryPoint, SectorHistoryResponse, SectorItem } from "./sectorTypes";

interface SectorRadarDrawerProps {
  drawerRef: React.RefObject<HTMLDivElement | null>;
  activeSectorDetail: SectorItem;
  historyData?: SectorHistoryResponse;
  historyTimeline: SectorHistoryPoint[];
  isHistoryLoading: boolean;
  fn1d: number | null;
  fn5d: number | null;
  isNetOutflow: boolean;
  gradientOffset: number;
  onClose: () => void;
}

export default function SectorRadarDrawer({
  drawerRef,
  activeSectorDetail,
  historyData,
  historyTimeline,
  isHistoryLoading,
  fn1d,
  fn5d,
  isNetOutflow,
  gradientOffset,
  onClose,
}: SectorRadarDrawerProps) {
  const badge = getForeignStateBadge(activeSectorDetail.foreign_state);

  return (
    <div
      ref={drawerRef}
      className="col-span-1 md:col-span-2 lg:col-span-3 w-full rounded-[22px] bg-[var(--md-sys-color-surface-container)] border-2 border-[var(--md-sys-color-primary)]/60 p-4 sm:p-6 shadow-lg space-y-5 transition-all duration-300 ease-out animate-in fade-in-50 slide-in-from-top-3"
    >
      {/* Header Drawer */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-3.5 border-b border-[var(--md-sys-color-outline-variant)]/60 gap-3">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0 border border-[var(--md-sys-color-outline-variant)]">
            <Icon icon={activeSectorDetail.icon} width="22" />
          </div>
          <div>
            <h3 className="text-sm sm:text-base font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-2">
              <span>{activeSectorDetail.name}</span>
              <span className="text-[10px] font-mono px-2 py-0.5 rounded-full bg-[var(--md-sys-color-secondary-container)] text-[var(--md-sys-color-on-secondary-container)]">
                {activeSectorDetail.code}
              </span>
            </h3>
            <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
              {activeSectorDetail.foreign_state_label} • Relative Strength RRG: {activeSectorDetail.quadrant}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <span
            className={`text-xs font-bold px-2.5 py-1 rounded-lg border flex items-center gap-1.5 ${badge.className}`}
          >
            <Icon icon={badge.icon} width="13" />
            <span>Radar: {activeSectorDetail.foreign_state}</span>
          </span>
          <button
            onClick={onClose}
            className="text-xs text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] p-1 rounded-lg hover:bg-[var(--md-sys-color-surface-container-high)] transition-colors cursor-pointer"
            title="Tutup Detail"
          >
            <Icon icon="ph:x-bold" width="14" />
          </button>
        </div>
      </div>

      {/* Quick Metrics of Active Sector */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-2.5">
        <div className="p-3 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]">
          <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] block">RS Ratio vs IHSG</span>
          <span className="text-sm sm:text-base font-bold text-[var(--md-sys-color-on-surface)] font-mono">
            {activeSectorDetail.rs_ratio}
          </span>
          <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] block mt-0.5">
            {activeSectorDetail.rs_ratio >= 100 ? "Outperforming (>100)" : "Underperforming (<100)"}
          </span>
        </div>

        <div className="p-3 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]">
          <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] block">RS Momentum (Smooth)</span>
          <span className="text-sm sm:text-base font-bold text-[var(--md-sys-color-on-surface)] font-mono">
            {activeSectorDetail.rs_momentum}
          </span>
          <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] block mt-0.5">
            {activeSectorDetail.rs_momentum >= 100 ? "Akselerasi Positif" : "Perlambatan Momentum"}
          </span>
        </div>

        <div className="p-3 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]">
          <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] block">Net Foreign (1D / 5D)</span>
          <div className="flex flex-col">
            <span className="text-xs sm:text-sm font-bold font-mono" style={{ color: signedColor(fn1d) }}>
              1D: {fn1d != null ? (fn1d > 0 ? `+${fmtRadarFlow(fn1d)}` : fmtRadarFlow(fn1d)) : "-"}
            </span>
            <span className="text-[10px] font-mono" style={{ color: signedColor(fn5d) }}>
              5D: {fn5d != null ? (fn5d > 0 ? `+${fmtRadarFlow(fn5d)}` : fmtRadarFlow(fn5d)) : "-"}
            </span>
          </div>
        </div>

        <div className="p-3 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]">
          <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] block">Turnover Harian</span>
          <span className="text-sm sm:text-base font-bold text-[var(--md-sys-color-primary)] font-mono">
            {fmtRp(activeSectorDetail.value)}
          </span>
          <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] block mt-0.5">
            Porsi {activeSectorDetail.value_share_pct}% total bursa
          </span>
        </div>
      </div>

      {/* DUAL-AXIS RECHARTS TIMELINE: INDEKS SEKTOR vs CUMULATIVE FOREIGN FLOW */}
      <SectorHistoryChart
        sectorCode={activeSectorDetail.code}
        historyTimeline={historyTimeline}
        isHistoryLoading={isHistoryLoading}
        isNetOutflow={isNetOutflow}
        gradientOffset={gradientOffset}
      />

      {/* TOP INSTITUTIONAL MOVERS (SAVIOR vs CULPRIT BROKERS) */}
      <SectorBrokerMovers
        saviorBrokers={historyData?.savior_brokers}
        culpritBrokers={historyData?.culprit_brokers}
      />

      {/* Historical Table / Points */}
      <SectorHistoryTable
        sectorCode={activeSectorDetail.code}
        historyTimeline={historyTimeline}
        isHistoryLoading={isHistoryLoading}
      />
    </div>
  );
}
