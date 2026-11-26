import type { SectorItem } from "./sectorTypes";

export const SECTOR_STORAGE_KEY = "investowl_sector_rotation_selected";
export const TIMEFRAME_STORAGE_KEY = "investowl_sector_rotation_timeframe";
export const QUADRANT_STORAGE_KEY = "investowl_sector_rotation_quadrant";

export function formatShortDate(dateStr: string): string {
  if (!dateStr) return "";
  try {
    const parts = dateStr.split("-");
    if (parts.length === 3) {
      const months = ["Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"];
      return `${parseInt(parts[2], 10)} ${months[parseInt(parts[1], 10) - 1]}`;
    }
    return dateStr;
  } catch {
    return dateStr;
  }
}

export function getForeignStateBadge(state: SectorItem["foreign_state"]) {
  switch (state) {
    case "INFLOW SURGE":
      return {
        label: "INFLOW SURGE",
        icon: "ph:trend-up-bold",
        className: "bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/30",
        dotColor: "bg-emerald-500",
      };
    case "SECTOR EXODUS":
      return {
        label: "SECTOR EXODUS",
        icon: "ph:warning-octagon-fill",
        className: "bg-rose-500/15 text-rose-600 dark:text-rose-400 border-rose-500/30 animate-pulse",
        dotColor: "bg-rose-500",
      };
    case "HIDDEN DISTRIBUTION":
      return {
        label: "HIDDEN DISTRIB",
        icon: "ph:warning-circle-fill",
        className: "bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/30",
        dotColor: "bg-amber-500",
      };
    case "BOTTOMING":
      return {
        label: "BOTTOMING",
        icon: "ph:arrow-circle-up-fill",
        className: "bg-sky-500/15 text-sky-600 dark:text-sky-400 border-sky-500/30",
        dotColor: "bg-sky-500",
      };
    default:
      return {
        label: "FLOW STABIL",
        icon: "ph:shield-check-duotone",
        className:
          "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]",
        dotColor: "bg-slate-400",
      };
  }
}
