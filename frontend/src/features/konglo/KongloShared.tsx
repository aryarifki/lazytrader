import { Link } from "react-router-dom";
import {
  RANK_COLORS_DOWN,
  RANK_COLORS_UP,
  fmtM,
  fmtMc,
  getGroupColor,
  type KongloStockRow,
} from "./kongloTypes";

export function KongloTickerLink({
  ticker,
  className = "text-xs",
}: {
  ticker: string;
  className?: string;
}) {
  return (
    <Link
      to={`/${encodeURIComponent(ticker)}`}
      className={`font-bold text-[var(--md-sys-color-on-surface)] hover:text-[var(--md-sys-color-primary)] hover:underline transition-colors font-mono ${className}`}
      title={`Buka analisis ${ticker}`}
    >
      {ticker}
    </Link>
  );
}

export function KongloGroupBadge({
  group,
  size = "xs",
  onClick,
  title,
}: {
  group?: string;
  size?: "xs" | "sm";
  onClick?: () => void;
  title?: string;
}) {
  if (!group) return null;
  const col = getGroupColor(group);
  const sizeClass =
    size === "sm"
      ? "text-[10px] px-2.5 py-0.5 font-medium"
      : "text-[9px] px-2 py-0.5 font-medium";

  const pillStyle = {
    backgroundColor: col.bg,
    color: col.color,
    border: `1px solid ${col.border}`,
  };

  if (onClick) {
    return (
      <span className="inline-flex rounded-full bg-[#1E1A17]/90 dark:bg-transparent">
        <button
          type="button"
          onClick={onClick}
          title={title}
          className="text-xs font-semibold px-2.5 py-0.5 rounded-full transition hover:opacity-85 cursor-pointer"
          style={pillStyle}
        >
          {group}
        </button>
      </span>
    );
  }

  return (
    <span className="inline-flex rounded-full bg-[#1E1A17]/90 dark:bg-transparent">
      <span className={`rounded-full ${sizeClass}`} style={pillStyle}>
        {group}
      </span>
    </span>
  );
}

export function KongloRankBadge({
  index,
  direction,
}: {
  index: number;
  direction: "up" | "down";
}) {
  const rankColor =
    direction === "up"
      ? RANK_COLORS_UP[index] || "#70E8B8"
      : RANK_COLORS_DOWN[index] || "#E88080";

  return (
    <div
      className="w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold text-white font-mono shadow-2xs"
      style={{ backgroundColor: rankColor }}
    >
      <span className="drop-shadow-[0_1px_1px_rgba(0,0,0,0.45)]">{index + 1}</span>
    </div>
  );
}

interface SubTableProps {
  title: string;
  metricHeader: string;
  rows?: KongloStockRow[];
  direction: "up" | "down";
  metricMode: "pct" | "mc";
  badgeMode: "group" | "funda";
}

function DualSubTable({
  title,
  metricHeader,
  rows,
  direction,
  metricMode,
  badgeMode,
}: SubTableProps) {
  const isUp = direction === "up";
  const toneClass = isUp
    ? "text-emerald-600 dark:text-emerald-400"
    : "text-rose-600 dark:text-rose-400";

  return (
    <div>
      <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider mb-2 font-mono">
        {title}
      </div>
      <table className="w-full text-xs text-left border-collapse">
        <thead>
          <tr className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/40 font-mono uppercase">
            <th className="py-1 px-1.5 w-7 font-bold">#</th>
            <th className="py-1 px-1.5 font-bold">
              {badgeMode === "group" ? "Kode · Konglo" : "Kode"}
            </th>
            <th className="py-1 px-1.5 text-right font-bold">{metricHeader}</th>
            <th className="py-1 px-1.5 text-right font-bold">Value (M)</th>
          </tr>
        </thead>
        <tbody>
          {rows?.map((s, idx) => (
            <tr
              key={s.ticker}
              className="border-b border-[var(--md-sys-color-outline-variant)]/20 hover:bg-[var(--md-sys-color-surface-container-high)]/40 transition-colors"
            >
              <td className="py-2 px-1.5">
                <KongloRankBadge index={idx} direction={direction} />
              </td>
              <td className="py-2 px-1.5">
                <span className="inline-flex items-center gap-1.5 flex-wrap">
                  <KongloTickerLink ticker={s.ticker} />
                  {badgeMode === "group" ? (
                    <KongloGroupBadge group={s.group} size="xs" />
                  ) : (
                    <span className="text-[9px] px-2 py-0.5 rounded-full font-medium bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)]">
                      Funda
                    </span>
                  )}
                </span>
              </td>
              <td className={`py-2 px-1.5 text-right font-semibold font-mono ${toneClass}`}>
                {metricMode === "pct"
                  ? isUp
                    ? `+${s.pct.toFixed(1)}%`
                    : `${s.pct.toFixed(1)}%`
                  : fmtMc(s.mc_chg, !isUp)}
              </td>
              <td className="py-2 px-1.5 text-right text-[var(--md-sys-color-on-surface-variant)] font-mono">
                {fmtM(s.val)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export function KongloDualBox({
  gainerTitle,
  loserTitle,
  badgeMode,
  top5PctUp,
  top5McUp,
  top5PctDown,
  top5McDown,
}: {
  gainerTitle: string;
  loserTitle: string;
  badgeMode: "group" | "funda";
  top5PctUp?: KongloStockRow[];
  top5McUp?: KongloStockRow[];
  top5PctDown?: KongloStockRow[];
  top5McDown?: KongloStockRow[];
}) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl overflow-hidden shadow-2xs">
        <div className="px-4 py-3 border-b border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-high)]/60 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-emerald-500" />
            <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)]">
              {gainerTitle}
            </span>
          </div>
          <span className="text-[10px] text-emerald-600 dark:text-emerald-400 font-bold font-mono">
            5 Terbaik
          </span>
        </div>

        <div className="p-4 space-y-4">
          <DualSubTable
            title="% Kenaikan Tertinggi"
            metricHeader="Change %"
            rows={top5PctUp}
            direction="up"
            metricMode="pct"
            badgeMode={badgeMode}
          />
          <div className="pt-3 border-t border-[var(--md-sys-color-outline-variant)]/30">
            <DualSubTable
              title="% Kenaikan MC Terbesar"
              metricHeader="Kenaikan MC"
              rows={top5McUp}
              direction="up"
              metricMode="mc"
              badgeMode={badgeMode}
            />
          </div>
        </div>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl overflow-hidden shadow-2xs">
        <div className="px-4 py-3 border-b border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-high)]/60 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-rose-500" />
            <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)]">
              {loserTitle}
            </span>
          </div>
          <span className="text-[10px] text-rose-600 dark:text-rose-400 font-bold font-mono">
            5 Terendah
          </span>
        </div>

        <div className="p-4 space-y-4">
          <DualSubTable
            title="% Penurunan Terbesar"
            metricHeader="Change %"
            rows={top5PctDown}
            direction="down"
            metricMode="pct"
            badgeMode={badgeMode}
          />
          <div className="pt-3 border-t border-[var(--md-sys-color-outline-variant)]/30">
            <DualSubTable
              title="% Penurunan MC Terbesar"
              metricHeader="Penurunan MC"
              rows={top5McDown}
              direction="down"
              metricMode="mc"
              badgeMode={badgeMode}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
