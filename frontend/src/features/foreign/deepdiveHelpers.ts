export const POPULAR_TICKERS = ["BBCA", "BBRI", "BMRI", "BBNI", "TLKM", "ASII", "ADRO", "BREN", "AMMN", "BRPT"];

export const WINDOW_OPTIONS = [
  { value: 20, label: "20 Hari" },
  { value: 30, label: "30 Hari" },
  { value: 60, label: "60 Hari" },
  { value: 90, label: "90 Hari" },
  { value: 180, label: "180 Hari" },
];

export function formatIdrCompact(val: number): string {
  const abs = Math.abs(val);
  const sign = val < 0 ? "-" : val > 0 ? "+" : "";
  if (abs >= 1e12) return `${sign}Rp ${(abs / 1e12).toFixed(2)} T`;
  if (abs >= 1e9) return `${sign}Rp ${(abs / 1e9).toFixed(1)} M`;
  if (abs >= 1e6) return `${sign}Rp ${(abs / 1e6).toFixed(0)} Jt`;
  return `${sign}Rp ${abs.toLocaleString("id-ID")}`;
}

export function getToneColor(tone?: string): string {
  if (tone === "positive") return "var(--color-positive)";
  if (tone === "negative") return "var(--color-negative)";
  if (tone === "warning") return "#f59e0b";
  return "var(--md-sys-color-outline)";
}
