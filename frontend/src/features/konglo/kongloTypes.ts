export type KongloTab =
  | "summary"
  | "pct-up"
  | "mc-up"
  | "pct-down"
  | "mc-down"
  | "konglo"
  | "fundamental";

export interface KongloTabItem {
  id: KongloTab;
  label: string;
  icon: string;
}

export const KONGLO_TABS: KongloTabItem[] = [
  { id: "summary", label: "Top Gainer & Loser", icon: "ph:trophy-duotone" },
  { id: "pct-up", label: "% Kenaikan Tertinggi", icon: "ph:trend-up-duotone" },
  { id: "mc-up", label: "Kenaikan MC Terbesar", icon: "ph:chart-bar-duotone" },
  { id: "pct-down", label: "% Penurunan Terbesar", icon: "ph:trend-down-duotone" },
  { id: "mc-down", label: "Penurunan MC Terbesar", icon: "ph:chart-line-down-duotone" },
  { id: "konglo", label: "Konglo Report", icon: "ph:buildings-duotone" },
  { id: "fundamental", label: "Funda Report", icon: "ph:file-text-duotone" },
];

export interface GroupColorStyle {
  bg: string;
  color: string;
  border: string;
}

// Exact Group Colors from legacy Konglo Report
export const GROUP_COLORS: Record<string, GroupColorStyle> = {
  "Barito Group":    { bg: "rgba(251,146,60,.15)",  color: "#FB923C", border: "rgba(251,146,60,.3)" },
  "Bakrie Group":    { bg: "rgba(239,68,68,.15)",   color: "#F87171", border: "rgba(239,68,68,.3)" },
  "Salim Group":     { bg: "rgba(59,130,246,.15)",  color: "#60A5FA", border: "rgba(59,130,246,.3)" },
  "Astra Group":     { bg: "rgba(16,185,129,.15)",  color: "#34D399", border: "rgba(16,185,129,.3)" },
  "Sinarmas Group":  { bg: "rgba(245,158,11,.15)",  color: "#FBBF24", border: "rgba(245,158,11,.3)" },
  "Djarum Group":    { bg: "rgba(139,92,246,.15)",  color: "#A78BFA", border: "rgba(139,92,246,.3)" },
  "Lippo Group":     { bg: "rgba(236,72,153,.15)",  color: "#F472B6", border: "rgba(236,72,153,.3)" },
  "MNC Group":       { bg: "rgba(20,184,166,.15)",  color: "#2DD4BF", border: "rgba(20,184,166,.3)" },
  "Emtek Group":     { bg: "rgba(99,102,241,.15)",  color: "#818CF8", border: "rgba(99,102,241,.3)" },
  "LTK Group":       { bg: "rgba(234,179,8,.15)",   color: "#EAB308", border: "rgba(234,179,8,.3)" },
  "Saratoga Group":  { bg: "rgba(6,182,212,.15)",   color: "#22D3EE", border: "rgba(6,182,212,.3)" },
  "Aguan Group":     { bg: "rgba(168,85,247,.15)",  color: "#C084FC", border: "rgba(168,85,247,.3)" },
  "Thohir Group":    { bg: "rgba(244,63,94,.15)",   color: "#FB7185", border: "rgba(244,63,94,.3)" },
  "Triputra Group":  { bg: "rgba(132,204,22,.15)",  color: "#A3E635", border: "rgba(132,204,22,.3)" },
  "Hapsoro Group":   { bg: "rgba(14,165,233,.15)",  color: "#38BDF8", border: "rgba(14,165,233,.3)" },
  "Northstar Group": { bg: "rgba(251,191,36,.15)",  color: "#FCD34D", border: "rgba(251,191,36,.3)" },
  "Kresna Group":    { bg: "rgba(217,70,239,.15)",  color: "#E879F9", border: "rgba(217,70,239,.3)" },
  "Ciputra Group":   { bg: "rgba(52,211,153,.15)",  color: "#6EE7B7", border: "rgba(52,211,153,.3)" },
  "Tancorp Group":   { bg: "rgba(251,113,133,.15)", color: "#FDA4AF", border: "rgba(251,113,133,.3)" },
  "Rajawali Group":  { bg: "rgba(196,181,253,.15)", color: "#C4B5FD", border: "rgba(196,181,253,.3)" },
  "Mayapada Group":  { bg: "rgba(253,186,116,.15)", color: "#FDB374", border: "rgba(253,186,116,.3)" },
  "SSI Group":       { bg: "rgba(110,231,183,.15)", color: "#6EE7B7", border: "rgba(110,231,183,.3)" },
  "Haji Isam Group": { bg: "rgba(147,197,253,.15)", color: "#93C5FD", border: "rgba(147,197,253,.3)" },
  "Sinarmas Land":   { bg: "rgba(165,243,252,.15)", color: "#A5F3FC", border: "rgba(165,243,252,.3)" },
  "Hashim Group":    { bg: "rgba(240,171,252,.15)", color: "#F0ABFC", border: "rgba(240,171,252,.3)" },
  "Harita Group":    { bg: "rgba(254,215,170,.15)", color: "#FED7AA", border: "rgba(254,215,170,.3)" },
  "Panin Group":     { bg: "rgba(167,243,208,.15)", color: "#A7F3D0", border: "rgba(167,243,208,.3)" },
  "CT Corp Group":   { bg: "rgba(253,164,175,.15)", color: "#FCA5A5", border: "rgba(253,164,175,.3)" },
  "Alfa Group":      { bg: "rgba(199,210,254,.15)", color: "#C7D2FE", border: "rgba(199,210,254,.3)" },
};

export const DEFAULT_GROUP_COLOR: GroupColorStyle = {
  bg: "rgba(58,123,213,.12)",
  color: "#60A5FA",
  border: "rgba(58,123,213,.25)",
};

export function getGroupColor(group?: string): GroupColorStyle {
  if (!group) return DEFAULT_GROUP_COLOR;
  return GROUP_COLORS[group] || DEFAULT_GROUP_COLOR;
}

export const RANK_COLORS_UP = [
  "#1a6644", "#1a6644", "#1a6644", "#1DB87A", "#1DB87A",
  "#25D48E", "#25D48E", "#25D48E", "#45DCA0", "#45DCA0",
  "#45DCA0", "#70E8B8", "#70E8B8", "#70E8B8", "#70E8B8",
];

export const RANK_COLORS_DOWN = [
  "#7A1A1A", "#7A1A1A", "#7A1A1A", "#B83030", "#B83030",
  "#D44A4A", "#D44A4A", "#D44A4A", "#E06060", "#E06060",
  "#E06060", "#E88080", "#E88080", "#E88080", "#E88080",
];

export interface KongloStockRow {
  ticker: string;
  group?: string;
  pct: number;
  val: number;
  mc_chg: number;
  prev?: number;
  close?: number;
  freq?: number;
}

export interface KongloGroupTop2 {
  ticker: string;
  pct: number;
}

export interface KongloTop5Group {
  group: string;
  top2?: KongloGroupTop2[];
  total_val: number;
  mc_chg: number;
  avg_pct: number;
}

export interface KongloGroupStats {
  group: string;
  total_mc: number;
  total_members: number;
  mc_chg: number;
  mc_pct: number;
  total_val: number;
  total_freq: number;
  members?: KongloStockRow[];
}

export interface KongloMissingItem {
  ticker: string;
  group: string;
}

export interface KongloOverview {
  total_konglo_stocks?: number;
  rising?: number;
  falling?: number;
  flat?: number;
  total_mc_chg?: number;
}

export interface KongloHighlightCard {
  group?: string;
  val?: number;
  pct?: number;
  count?: number;
}

export interface KongloCards {
  overview?: KongloOverview;
  best_mc?: KongloHighlightCard;
  worst_mc?: KongloHighlightCard;
  best_pct?: KongloHighlightCard;
  worst_pct?: KongloHighlightCard;
}

export interface KongloFundamentalSummary {
  total_stocks?: number;
  rising?: number;
  falling?: number;
  flat?: number;
  total_mc_chg?: number;
  mc_pct?: number;
  top5_pct_up?: KongloStockRow[];
  top5_mc_up?: KongloStockRow[];
  top5_pct_down?: KongloStockRow[];
  top5_mc_down?: KongloStockRow[];
}

export interface KongloReportMetadata {
  date_display?: string;
  time_display?: string;
  last_updated_str?: string;
}

export interface KongloReportSummary {
  stocks_up?: number;
  stocks_down?: number;
  stocks_flat?: number;
  total_active_stocks?: number;
  total_stocks?: number;
  ihsg_close?: string | number;
  ihsg_pct?: string | number;
}

export interface KongloReportData {
  metadata?: KongloReportMetadata;
  summary?: KongloReportSummary;
  konglo_cards?: KongloCards;
  top5_summary?: {
    top5_pct_up?: KongloStockRow[];
    top5_mc_up?: KongloStockRow[];
    top5_pct_down?: KongloStockRow[];
    top5_mc_down?: KongloStockRow[];
  };
  top15_tables?: {
    pct_up?: KongloStockRow[];
    mc_up?: KongloStockRow[];
    pct_down?: KongloStockRow[];
    mc_down?: KongloStockRow[];
  };
  top5_konglo?: KongloTop5Group[];
  all_group_stats?: KongloGroupStats[];
  fundamental_summary?: KongloFundamentalSummary;
  fundamental_report?: KongloStockRow[];
  missing_funda?: string[];
  missing_konglo?: KongloMissingItem[];
}

export interface KongloReportApiResponse {
  report?: KongloReportData;
}

export function fmtNum(n?: number | null): string {
  return n ? Number(n).toLocaleString("id-ID") : "0";
}

export function fmtM(n?: number | null): string {
  return n ? `${(n / 1e9).toFixed(1)} M` : "-";
}

export function fmtT(n?: number | null): string {
  const val = Number(n ?? 0);
  return Math.abs(val) >= 10
    ? Math.round(Math.abs(val)).toLocaleString("id-ID")
    : Math.abs(val).toFixed(2);
}

export function fmtMc(n?: number | null, forceNegativeSign = false): string {
  const val = Number(n ?? 0);
  if (forceNegativeSign || val < 0) {
    return `-Rp ${fmtT(val)} T`;
  }
  return `+Rp ${fmtT(val)} T`;
}

export function fmtFreq(n?: number | null): string {
  return n ? Math.round(n).toLocaleString("id-ID") : "-";
}
