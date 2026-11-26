const BADGE_BASE =
  "inline-flex items-center justify-center font-extrabold tracking-[0.3px] rounded-md text-center transition-all";

const BADGE_STYLES: Record<string, string> = {
  "massive-acc":
    "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border border-emerald-500/35 shadow-[0_0_12px_-3px_rgba(47,217,138,0.22)]",
  "big-acc":
    "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border border-emerald-500/35 shadow-[0_0_12px_-3px_rgba(47,217,138,0.22)]",
  "normal-acc":
    "bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border border-emerald-500/25",
  "small-acc":
    "bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border border-emerald-500/25",
  undetected:
    "bg-slate-500/12 text-slate-700 dark:text-slate-300 border border-slate-500/25",
  "small-dist":
    "bg-rose-500/10 text-rose-700 dark:text-rose-300 border border-rose-500/25",
  "normal-dist":
    "bg-rose-500/10 text-rose-700 dark:text-rose-300 border border-rose-500/25",
  "big-dist":
    "bg-rose-500/15 text-rose-700 dark:text-rose-400 border border-rose-500/35 shadow-[0_0_12px_-3px_rgba(239,68,68,0.22)]",
  "massive-dist":
    "bg-rose-500/15 text-rose-700 dark:text-rose-400 border border-rose-500/35 shadow-[0_0_12px_-3px_rgba(239,68,68,0.22)]",
  "aggressive-buy":
    "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border border-emerald-500/50",
  "normal-buy":
    "bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border border-emerald-500/40",
  "normal-sell":
    "bg-rose-500/10 text-rose-700 dark:text-rose-300 border border-rose-500/40",
  "panic-sell":
    "bg-rose-500/15 text-rose-700 dark:text-rose-400 border border-rose-500/50",
};

export function bandarBadgeClass(cls?: string): string {
  const key = (cls || "undetected").toLowerCase();
  const variant = BADGE_STYLES[key] || BADGE_STYLES.undetected;
  return `${BADGE_BASE} ${variant}`;
}

export function brokerCategoryClass(kategori?: string): string {
  const base = "inline-block px-1.5 py-0.5 rounded text-[10px]";
  if (kategori === "Retail") {
    return `${base} font-bold text-teal-700 dark:text-[#2FD4C6] bg-teal-500/12 border border-teal-500/25`;
  }
  if (kategori === "Institusi") {
    return `${base} font-bold text-amber-700 dark:text-[#E0954D] bg-amber-500/12 border border-amber-500/25`;
  }
  if (kategori === "Pengendali") {
    return `${base} font-extrabold text-purple-700 dark:text-purple-400 bg-purple-500/14 border border-purple-500/30`;
  }
  if (kategori === "Semi Ritel" || !kategori) {
    return `${base} font-bold text-indigo-700 dark:text-[#818CF8] bg-indigo-500/12 border border-indigo-500/25`;
  }
  return `${base} font-bold text-indigo-700 dark:text-[#818CF8] bg-indigo-500/12 border border-indigo-500/25`;
}
