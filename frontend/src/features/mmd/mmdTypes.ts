export type MainViewMode = "main" | "bigPlayer" | "explorer";
export type BandarTabMode = "resume" | "detail" | "other";
export type ExplorerCategory = "normal" | "advanced" | "bigPlayer";
export type MobileBrokerTab = "buyer" | "seller" | "accdist";

export interface MmdTicker {
  ticker: string;
  company_name?: string | null;
  name?: string | null;
  last_date?: string | null;
  last_close?: number | null;
  last_net_volume?: number | null;
  earliest_date?: string | null;
  latest_date?: string | null;
  total_days?: number | null;
  last_volume?: number | null;
  last_tavg?: number | null;
}

export interface MmdCandle {
  date_iso: string;
  date_str?: string | null;
  file_name?: string | null;
  open: number | null;
  high: number | null;
  low: number | null;
  close: number | null;
  volume: number | null;
  prev?: number | null;
  tavg?: number | null;
  net_volume?: number | null;
  buy_brokers?: Record<string, number>;
  sell_brokers?: Record<string, number>;
}

export interface MmdBandarAccPoint {
  date: string;
  net_volume: number;
  cumulative_bandar: number;
  close?: number | null;
}

export interface MmdDetailResponse {
  success?: boolean;
  ticker: string;
  meta?: MmdTicker | null;
  total_candles?: number;
  candles: MmdCandle[];
  bandar_accumulation_series: MmdBandarAccPoint[];
}

export interface MmdBadge {
  label: string;
  cls: string;
  color: string;
}

export interface MmdBrokerItem {
  broker: string;
  net: number;
  volume?: number;
  volume_fmt: string;
  value_jt: number;
  kategori: string;
}

export interface MmdTop10Row {
  rank: number;
  buyer_broker: string;
  buyer_net: number;
  buyer_net_fmt: string;
  buyer_val_jt: number;
  buyer_kategori: string;
  seller_broker: string;
  seller_net: number;
  seller_net_fmt: string;
  seller_val_jt: number;
  seller_kategori: string;
  acc_dist: number;
  acc_dist_fmt: string;
}

export interface MmdBandarRow {
  label: string;
  buy: number;
  sell: number;
  nilai: number;
  ratio: number;
  badge: MmdBadge;
  is_action?: boolean;
}

export interface MmdRetailRow {
  broker: string;
  trans: number;
  pct_trans: number;
  net: number;
  pct_net: number;
  is_retail_bandar: boolean;
}

export interface MmdPeriodItem {
  label: string;
  date: string;
  badge: MmdBadge;
  change_pct: number;
  pengendali_action: MmdBadge;
  retail_action: MmdBadge;
  range: string;
}

export interface MmdPreviousResults {
  weekly: MmdPeriodItem[];
  monthly: MmdPeriodItem[];
}

export interface MmdGroupMember {
  ticker: string;
  last_close: number;
  today_pct: number;
  w1_pct?: number | null;
  m1_pct: number;
  month_pct?: number;
  avg_val: number;
}

export interface MmdGroupInfo {
  group_name: string;
  title: string;
  type: string;
  members: MmdGroupMember[];
}

export interface MmdGroupPerf {
  conglomerate?: MmdGroupInfo | null;
  sector?: MmdGroupInfo | null;
}

export interface MmdPengendaliGroupStat {
  main: number;
  sub: number;
}

export interface MmdPengendaliStats {
  institusi: MmdPengendaliGroupStat;
  semi_ritel: MmdPengendaliGroupStat;
  ritel: MmdPengendaliGroupStat;
  non_ritel: MmdPengendaliGroupStat;
  repo: { val: number };
  lainnya: { val: number };
  ratio: number;
  action: MmdBadge;
}

export interface MmdAIScoreComponents {
  c1: number;
  c2: number;
  c3: number;
  c4: number;
  c5: number;
  c6: number;
  c7: number;
  c8: number;
  c9: number;
  c10: number;
  c_corr: number;
  c_mark: number;
  c_daily30: number;
  c_sideways: number;
  total_raw: number;
  total_scaled: number;
}

export interface MmdT5Item {
  date: string;
  investor?: string | null;
  broker?: string | null;
  action?: string | null;
  lot: number;
  value_m: number;
  price: number;
  pct: number;
}

export interface MmdAnalysis {
  ticker: string;
  total_days: number;
  date_from: string;
  date_to: string;
  last_price: number;
  price_change: number;
  total_buy: number;
  total_sell: number;
  total_volume: number;
  tot_net_volume: number;
  avg_daily_value: number;
  buyers: MmdBrokerItem[];
  sellers: MmdBrokerItem[];
  acc_dist_rows: number[];
  top10_combined: MmdTop10Row[];
  bandar_table: MmdBandarRow[];
  bandar_action: MmdBadge;
  retail_rows: MmdRetailRow[];
  retail_ratio: number;
  retail_action: MmdBadge;
  retail_bandar_auto_result?: string | null;
  active_retail_bandar?: string | null;
  retail_bandar_auto_enabled: boolean;
  retail_bandar_manual_override: string;
  effective_retail_brokers: string[];
  pengendali_stats: MmdPengendaliStats;
  resume: MmdBadge;
  ai_score: number;
  korelasi: number;
  ai_score_components?: MmdAIScoreComponents;
  previous_results: MmdPreviousResults;
  ai_commentary?: string[];
  commentary?: string[];
  big_player_summary?: string[];
  t5_records?: MmdT5Item[];
  group_performance?: MmdGroupPerf | null;
}

export interface MmdBrokerDrilldownDay {
  date: string;
  buy_lot: number;
  buy_val: number;
  buy_avg: number;
  sell_lot: number;
  sell_val: number;
  sell_avg: number;
  net_lot: number;
  net_val: number;
}

export interface MmdBrokerDrilldownSummary {
  total_buy_lot: number;
  total_buy_value_m: number;
  total_buy_value_b: number;
  avg_buy_price: number;
  total_sell_lot: number;
  total_sell_value_m: number;
  total_sell_value_b: number;
  avg_sell_price: number;
  net_lot: number;
  net_value_m: number;
  net_value_b: number;
  days_active: number;
}

export interface MmdBrokerDrilldown {
  success?: boolean;
  ticker: string;
  broker_code: string;
  broker_name?: string;
  category?: string;
  date_from?: string | null;
  date_to?: string | null;
  summary?: MmdBrokerDrilldownSummary;
  daily_breakdown?: MmdBrokerDrilldownDay[];
  daily?: MmdBrokerDrilldownDay[];
  transactions?: MmdBrokerDrilldownDay[];
}

export interface MmdExplorerRow {
  ticker: string;
  company_name?: string | null;
  last_close: number;
  last_price?: number;
  last_volume: number;
  price_change: number;
  change_pct?: number;
  ratio?: number;
  badge: MmdBadge;
  detail: string;
}
