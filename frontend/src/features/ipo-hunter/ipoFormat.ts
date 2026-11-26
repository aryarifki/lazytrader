export function fmtLot(n: number | null | undefined): string {
  if (n == null || Number.isNaN(n)) return "—";
  return Math.round(n).toLocaleString("en-US");
}

export function fmtSignedLot(n: number | null | undefined): string {
  if (n == null || Number.isNaN(n)) return "—";
  const r = Math.round(n);
  return `${r >= 0 ? "+" : ""}${r.toLocaleString("en-US")}`;
}

export function fmtPct(v: number | null | undefined): string {
  if (v == null || Number.isNaN(v)) return "—";
  return `${v >= 0 ? "+" : ""}${v.toFixed(1)}%`;
}

export function fmtSignedPct0(p: number | null | undefined): string {
  if (p == null || Number.isNaN(p)) return "—";
  return `${p > 0 ? "+" : ""}${p}%`;
}

export function fmtInt(n: number | null | undefined): string {
  if (n == null || Number.isNaN(n)) return "—";
  return Math.round(n).toLocaleString("en-US");
}

export function tbsColorClass(styleKey: string): string {
  switch (styleKey) {
    case "uw":
      return "text-rose-600 dark:text-rose-400 font-extrabold";
    case "institusi":
      return "text-rose-500 dark:text-rose-400 font-bold";
    case "retail":
      return "text-sky-600 dark:text-sky-400 font-bold";
    case "nonRitel":
      return "text-amber-600 dark:text-amber-400 font-bold";
    case "semi":
      return "text-indigo-600 dark:text-indigo-400 font-bold";
    default:
      return "font-bold";
  }
}

export function retailBadgeClass(cls: string | undefined): string {
  if (cls === "retail-buy") {
    return "bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/30";
  }
  if (cls === "retail-sell") {
    return "bg-rose-500/15 text-rose-600 dark:text-rose-400 border-rose-500/30";
  }
  return "bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]";
}

export function uwChipClass(cat: string): string {
  switch (cat) {
    case "retail":
      return "bg-sky-500/15 text-sky-600 dark:text-sky-400 border-sky-500/30";
    case "institusi":
      return "bg-rose-500/15 text-rose-600 dark:text-rose-400 border-rose-500/30";
    case "semi":
      return "bg-indigo-500/15 text-indigo-600 dark:text-indigo-400 border-indigo-500/30";
    case "nonRitel":
      return "bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/30";
    default:
      return "bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]";
  }
}
