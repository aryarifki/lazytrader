import { toneFor } from "@/components/research/format";

export function pctClass(value: number | null | undefined): string {
  const tone = toneFor(value);
  if (tone === "up") return "text-emerald-600 dark:text-emerald-400";
  if (tone === "down") return "text-rose-600 dark:text-rose-400";
  return "text-[var(--md-sys-color-on-surface-variant)]";
}
