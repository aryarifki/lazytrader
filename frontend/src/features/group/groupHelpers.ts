import { fmtNumber, toneFor } from "@/components/research/format";

export const VIEWS = [
  { id: "jejaring", label: "Jejaring", icon: "ph:share-network-duotone" },
  { id: "matriks", label: "Matriks", icon: "ph:table-duotone" },
  { id: "emiten", label: "Emiten", icon: "ph:buildings-duotone" },
] as const;

export const fetcher = async (url: string) => {
  const response = await fetch(url, { credentials: "include" });
  if (!response.ok) {
    let detail = `HTTP ${response.status}`;
    try {
      const body = await response.json();
      detail = body?.detail ?? detail;
    } catch {
      // keep status
    }
    throw new Error(detail);
  }
  return response.json();
};

export function fmtMarketCapB(value: number | null | undefined): string {
  if (value == null || Number.isNaN(value)) return "—";
  if (value >= 1_000) return `Rp ${(value / 1_000).toFixed(2)} T`;
  return `Rp ${fmtNumber(value, 0)} M`;
}

export function fmtLot(value: number | null | undefined): string {
  if (value == null || Number.isNaN(value)) return "—";
  const abs = Math.abs(value);
  const sign = value < 0 ? "-" : "";
  if (abs >= 1_000_000) return `${sign}${(abs / 1_000_000).toFixed(1)} jt`;
  if (abs >= 1_000) return `${sign}${(abs / 1_000).toFixed(0)}k`;
  return `${sign}${fmtNumber(abs)}`;
}

export function pctClass(value: number | null | undefined): string {
  const tone = toneFor(value);
  if (tone === "up") return "text-emerald-600 dark:text-emerald-400";
  if (tone === "down") return "text-rose-600 dark:text-rose-400";
  return "text-[var(--md-sys-color-on-surface-variant)]";
}
