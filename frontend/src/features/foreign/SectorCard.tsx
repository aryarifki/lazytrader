import { Link } from "react-router-dom";
import { Icon } from "@iconify/react";
import { fmtRp, fmtRadarFlow, signedColor } from "@/lib/format";
import { getForeignStateBadge } from "./sectorHelpers";
import type { SectorItem } from "./sectorTypes";

interface SectorCardProps {
  sector: SectorItem;
  timeframe: "1d" | "5d" | "20d";
  isSelected: boolean;
  onSelect: () => void;
}

export default function SectorCard({
  sector,
  timeframe,
  isSelected,
  onSelect,
}: SectorCardProps) {
  const ret =
    timeframe === "5d" ? sector.ret_5d : timeframe === "20d" ? sector.ret_20d : sector.ret_1d;
  const foreignFlow =
    timeframe === "5d"
      ? sector.foreign_net_5d
      : timeframe === "20d"
      ? sector.foreign_net_20d
      : sector.foreign_net_1d;

  const isUp = ret > 0;
  const isDown = ret < 0;

  const quadrantBadgeColor =
    sector.quadrant === "Leading"
      ? "bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/30"
      : sector.quadrant === "Improving"
      ? "bg-sky-500/15 text-sky-600 dark:text-sky-400 border-sky-500/30"
      : sector.quadrant === "Weakening"
      ? "bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/30"
      : "bg-rose-500/15 text-rose-600 dark:text-rose-400 border-rose-500/30";

  const fBadge = getForeignStateBadge(sector.foreign_state);

  return (
    <div
      onClick={onSelect}
      className={`p-4 rounded-2xl border transition-all cursor-pointer flex flex-col justify-between ${
        isSelected
          ? "bg-[var(--md-sys-color-surface-container-high)] border-[var(--md-sys-color-primary)] shadow-sm ring-2 ring-[var(--md-sys-color-primary)]/50 scale-[1.01]"
          : "bg-[var(--md-sys-color-surface-container)] border-[var(--md-sys-color-outline-variant)] hover:border-[var(--md-sys-color-primary)]/40 hover:bg-[var(--md-sys-color-surface-container-high)]"
      }`}
    >
      <div>
        {/* Sector Header */}
        <div className="flex items-start justify-between gap-2 mb-2">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0 border border-[var(--md-sys-color-outline-variant)]">
              <Icon icon={sector.icon} width="18" />
            </div>
            <div>
              <h3
                className="text-xs font-extrabold text-[var(--md-sys-color-on-surface)] truncate max-w-[160px]"
                title={sector.name}
              >
                {sector.name}
              </h3>
              <span className="text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                {sector.code}
              </span>
            </div>
          </div>

          <div className="flex flex-col items-end gap-1">
            {/* Foreign State Badge */}
            <span
              className={`text-[9px] font-bold px-1.5 py-0.5 rounded-md border flex items-center gap-1 ${fBadge.className}`}
              title={sector.foreign_state_label}
            >
              <Icon icon={fBadge.icon} width="11" />
              <span>{fBadge.label}</span>
            </span>

            {/* RRG Quadrant Badge */}
            <span className={`text-[9px] font-semibold px-1.5 py-0.2 rounded-md border ${quadrantBadgeColor}`}>
              {sector.quadrant}
            </span>
          </div>
        </div>

        {/* Price & Foreign Flow Display */}
        <div className="grid grid-cols-2 gap-2 mt-3 pt-2.5 border-t border-[var(--md-sys-color-outline-variant)]/50">
          {/* Price & Return */}
          <div className="flex flex-col">
            <span className="text-[9px] text-[var(--md-sys-color-on-surface-variant)]">
              Indeks ({timeframe.toUpperCase()})
            </span>
            <div className="flex items-baseline gap-1.5 mt-0.5">
              <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] font-mono">
                {sector.close.toLocaleString("id-ID", { minimumFractionDigits: 1, maximumFractionDigits: 1 })}
              </span>
              <span
                className={`text-[11px] font-extrabold ${
                  isUp
                    ? "text-emerald-500 dark:text-emerald-400"
                    : isDown
                    ? "text-rose-500 dark:text-rose-400"
                    : "text-[var(--md-sys-color-on-surface-variant)]"
                }`}
              >
                {ret > 0 ? `+${ret}%` : `${ret}%`}
              </span>
            </div>
          </div>

          {/* Foreign Net Flow */}
          <div className="flex flex-col items-end">
            <span className="text-[9px] text-[var(--md-sys-color-on-surface-variant)]">
              Net Foreign ({timeframe.toUpperCase()})
            </span>
            <span
              className="text-xs font-extrabold font-mono mt-0.5"
              style={{ color: signedColor(foreignFlow) }}
            >
              {foreignFlow > 0 ? `+${fmtRadarFlow(foreignFlow)}` : fmtRadarFlow(foreignFlow)}
            </span>
          </div>
        </div>

        {/* Turnover Share Progress Bar */}
        <div className="mt-3 space-y-1">
          <div className="flex items-center justify-between text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
            <span>Turnover: {fmtRp(sector.value)}</span>
            <span className="font-semibold text-[var(--md-sys-color-on-surface)]">
              {sector.value_share_pct}% pasar
            </span>
          </div>
          <div className="h-1.5 w-full bg-[var(--md-sys-color-surface-container-highest)] rounded-full overflow-hidden">
            <div
              style={{ width: `${Math.min(sector.value_share_pct * 3.5, 100)}%` }}
              className="h-full bg-[var(--md-sys-color-primary)] rounded-full transition-all"
            />
          </div>
        </div>
      </div>

      {/* Top Constituent Stocks Chips & Accordion Indicator */}
      <div className="mt-3.5 pt-2.5 border-t border-[var(--md-sys-color-outline-variant)]/50 space-y-2">
        <div className="flex items-center justify-between text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
          <span className="font-semibold">Top Penggerak Sektor:</span>
          <span className="text-[9px]">Klik emiten</span>
        </div>
        <div className="flex flex-wrap gap-1">
          {sector.top_stocks.map((ticker) => (
            <Link
              key={ticker}
              to={`/${ticker}`}
              onClick={(e) => e.stopPropagation()}
              className="text-[10px] font-mono font-bold px-2 py-0.5 rounded-md bg-[var(--md-sys-color-surface-container-low)] text-[var(--md-sys-color-primary)] hover:bg-[var(--md-sys-color-primary)] hover:text-[var(--md-sys-color-on-primary)] border border-[var(--md-sys-color-outline-variant)] transition-all"
            >
              {ticker}
            </Link>
          ))}
        </div>

        {/* Dynamic Action Trigger Bar */}
        <div className="flex items-center justify-between pt-1 text-[11px] font-bold text-[var(--md-sys-color-primary)] border-t border-[var(--md-sys-color-outline-variant)]/30">
          <span className="flex items-center gap-1">
            <Icon icon={isSelected ? "ph:caret-up-bold" : "ph:caret-down-bold"} width="14" />
            <span>{isSelected ? "Tutup Radar Sektor" : "Buka Radar & Dinamika Asing"}</span>
          </span>
          {isSelected && (
            <span className="flex h-2 w-2 relative">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-[var(--md-sys-color-primary)] opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-[var(--md-sys-color-primary)]"></span>
            </span>
          )}
        </div>
      </div>
    </div>
  );
}
