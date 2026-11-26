export function getWyckoffBadge(phase: string): string {
  const p = (phase || "").toUpperCase();
  if (p === "E") return "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/30";
  if (p === "D") return "bg-cyan-500/15 text-cyan-700 dark:text-cyan-400 border-cyan-500/30";
  if (p === "C") return "bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/30";
  if (p === "B") return "bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/30";
  if (p === "A") return "bg-rose-500/15 text-rose-700 dark:text-rose-400 border-rose-500/30";
  return "bg-neutral-500/15 text-neutral-600 dark:text-neutral-400 border-neutral-500/30";
}

export function getScoreBadge(score: number): string {
  if (score >= 70) return "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/30";
  if (score >= 50) return "bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/30";
  if (score >= 40) return "bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/30";
  return "bg-rose-500/15 text-rose-700 dark:text-rose-400 border-rose-500/30";
}

export function getBandarSignalBadge(sig: string): string {
  const s = (sig || "").toUpperCase();
  if (s.includes("BIG ACCUMULATION")) return "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/30";
  if (s.includes("ACCUMULATION") || s.includes("NET BUY"))
    return "bg-teal-500/15 text-teal-700 dark:text-teal-400 border-teal-500/30";
  if (s.includes("BIG DISTRIBUTION")) return "bg-rose-500/15 text-rose-700 dark:text-rose-400 border-rose-500/30";
  if (s.includes("DISTRIBUTION") || s.includes("NET SELL"))
    return "bg-orange-500/15 text-orange-700 dark:text-orange-400 border-orange-500/30";
  return "bg-neutral-500/15 text-neutral-600 dark:text-neutral-400 border-neutral-500/30";
}

export function getHmmBadge(state: string): string {
  const s = (state || "").toLowerCase();
  if (s.includes("strong accum")) return "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/30";
  if (s.includes("accum")) return "bg-cyan-500/15 text-cyan-700 dark:text-cyan-400 border-cyan-500/30";
  if (s.includes("strong dist")) return "bg-rose-500/15 text-rose-700 dark:text-rose-400 border-rose-500/30";
  if (s.includes("dist")) return "bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/30";
  return "bg-neutral-500/15 text-neutral-600 dark:text-neutral-400 border-neutral-500/30";
}
