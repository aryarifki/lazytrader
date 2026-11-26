import type { IpoScreenerRow } from "./ipoTypes";

export type ScreenerMetricKey = Exclude<keyof IpoScreenerRow, "kode" | "company_name" | "listing_date">;

export interface ScreenerMetricDef {
  k: ScreenerMetricKey;
  label: string;
  kind: "num" | "pct" | "x";
  desc: string;
}

export const SCR_METRICS: ScreenerMetricDef[] = [
  {
    k: "ratrans",
    label: "Average Trans (M)",
    kind: "num",
    desc: "Rata-rata nilai transaksi harian (dalam miliar Rp)",
  },
  {
    k: "uw",
    label: "Underwriter Buyback",
    kind: "pct",
    desc: "Net beli semua broker underwriter ÷ total volume saham (%)",
  },
  {
    k: "netRetail",
    label: "Retail Sell",
    kind: "pct",
    desc: "Net jual broker Retail ÷ total volume (%). Positif = retail net sell (jual lebih besar dari beli); negatif = retail net buy",
  },
  { k: "netSemi", label: "Net Semi Retail", kind: "pct", desc: "Net broker Semi Retail ÷ total volume (%)" },
  { k: "netInst", label: "Net Institusi", kind: "pct", desc: "Net broker Institusi ÷ total volume (%)" },
  { k: "netNonRitel", label: "Net Non Retail", kind: "pct", desc: "Net broker Non Retail ÷ total volume (%)" },
  { k: "netRepo", label: "Net Repo + Others", kind: "pct", desc: "Net broker Repo + Others ÷ total volume (%)" },
  { k: "perfIpo", label: "IPO vs Close", kind: "pct", desc: "Close terakhir periode vs harga IPO (%)" },
  {
    k: "perfAvg",
    label: "Average vs Close",
    kind: "pct",
    desc: "Close terakhir vs harga rata-rata (VWAP) periode (%)",
  },
  {
    k: "chg",
    label: "Change",
    kind: "pct",
    desc: "Perubahan harga periode: Close akhir vs Prev hari pertama (Mode After IPO: vs harga IPO)",
  },
  {
    k: "volRatio",
    label: "Volume Melonjak",
    kind: "x",
    desc: "Rata-rata volume harian periode ÷ rata-rata volume harian sebelum periode (maks. 20 hari data sebelumnya). Mode After IPO: volume hari terakhir ÷ rata-rata volume hari-hari sebelumnya",
  },
];

export interface ScreenerPreset {
  id: string;
  label: string;
  col: ScreenerMetricKey;
  dir: number;
  desc: string;
}

export const SCR_PRESETS: ScreenerPreset[] = [
  {
    id: "topPerf",
    label: "Top Performer",
    col: "chg",
    dir: -1,
    desc: "Kenaikan harga terbesar di periode (Change tertinggi dulu)",
  },
  {
    id: "topLoser",
    label: "Top Loser",
    col: "chg",
    dir: 1,
    desc: "Penurunan harga terbesar di periode (Change terendah dulu)",
  },
  {
    id: "volSurge",
    label: "Volume Melonjak",
    col: "volRatio",
    dir: -1,
    desc: "Volume harian rata-rata periode dibanding sebelum periode (kali lipat, terbesar dulu)",
  },
];

export const SCR_PRIMARY: ScreenerMetricKey[] = [
  "uw",
  "netRetail",
  "netSemi",
  "netInst",
  "netNonRitel",
  "netRepo",
  "perfIpo",
  "perfAvg",
];

export function fmtScreenerVal(k: string, v: number | null | undefined): string {
  if (v == null || Number.isNaN(v)) return "—";
  const m = SCR_METRICS.find((x) => x.k === k);
  if (m?.kind === "num") {
    return v.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  }
  if (m?.kind === "x") return `${v.toFixed(1)}x`;
  return `${v > 0 ? "+" : ""}${v.toFixed(1)}%`;
}
