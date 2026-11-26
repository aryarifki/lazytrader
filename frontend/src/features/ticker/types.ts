export interface AffiliatedBrokerDetected {
  broker_code: string;
  broker_name: string;
  badge: string;
  relation: string;
  days_active: number;
  cumulative_net: number;
  latest_date: string;
  latest_net: number;
}

export interface PriceChartPoint {
  date: string;
  close: number;
  volume?: number;
  bandar_vwap?: number | null;
}

export interface SignalOverlayPoint {
  date: string;
  signal: string;
  score: number;
}

export interface BrokerSummaryRow {
  side: string;
  broker: string;
  type: string;
  net: number;
  spark: string;
  is_affiliated?: boolean;
  affiliation?: { relation?: string };
}

export interface PricePerformanceRow {
  period: string;
  value: number;
}

export interface RetailBrokerRow {
  broker: string;
  net: number;
  buy: number;
  sell: number;
}

export interface RetailSummary {
  sentiment: string;
  sentiment_tone: "positive" | "negative" | "primary" | string;
  sentiment_desc: string;
  smart_cum_net: number;
  retail_cum_net: number;
  top_retail_brokers?: RetailBrokerRow[];
}

export interface SmartDailyRow {
  date: string;
  smart_net: number;
  cumulative_net: number;
}

export interface ProfileFlowBroker {
  broker_code: string;
  participant_type: string;
  net: number;
}

export interface ProfileFlowRow {
  label: string;
  description?: string;
  net: number;
  top_brokers?: ProfileFlowBroker[];
}

export interface ProfileBrokerDetailRow {
  profile: string;
  profile_key?: string;
  broker: string;
  type: string;
  buy: number;
  sell: number;
  net: number;
  freq: number;
  days?: number;
  avg_value_tx?: number;
}

export interface TickerDetailData {
  error?: string;
  ticker: string;
  company?: { name?: string; group?: string };
  window_start?: string;
  analysis_date?: string;
  conviction_score?: number | null;
  conviction?: { score?: number; broker_note?: string };
  conviction_breakdown?: { broker_note?: string };
  score?: number;
  broker_note?: string;
  signal?: string | null;
  signal_score?: number | null;
  ret_5d?: number | null;
  foreign_5d?: number | null;
  top_buyer?: { broker?: string; net?: number | null };
  smart_cumulative?: number | null;
  smart_daily?: SmartDailyRow[];
  alerts?: string[];
  verdict?: string;
  affiliated_brokers_detected?: AffiliatedBrokerDetected[];
  bandar_pnl_pct?: number | null;
  bandar_vwap?: number | null;
  price_chart?: PriceChartPoint[];
  signal_overlay?: SignalOverlayPoint[];
  broker_summary?: BrokerSummaryRow[];
  price_performance?: PricePerformanceRow[];
  retail_summary?: RetailSummary | null;
  profile_flow?: ProfileFlowRow[];
  profile_broker_detail?: ProfileBrokerDetailRow[];
}

export interface DistributionEdge {
  buyer_code: string;
  buyer_type: string;
  seller_code: string;
  seller_type: string;
  matched_value: number;
}

export interface BrokerComparisonSummaryRow {
  buy_broker?: string;
  buy_type?: string;
  buy_value?: number | null;
  buy_lot?: number | null;
  buy_avg?: number | null;
  sell_broker?: string;
  sell_type?: string;
  sell_value?: number | null;
  sell_lot?: number | null;
  sell_avg?: number | null;
}

export interface DetailedBrokerRow {
  broker: string;
  type: string;
  subtype?: string;
  buy?: number;
  sell?: number;
  net: number;
  freq?: number;
}

export interface BrokerFlowData {
  all_codes?: string[];
  default_codes?: string[];
  compare_chart?: Array<Record<string, string | number | null>>;
  distribution?: {
    buyers?: unknown[];
    sellers?: unknown[];
    edges?: DistributionEdge[];
    dist_start?: string;
    dist_end?: string;
  };
  summary?: BrokerComparisonSummaryRow[];
  detail_rows?: DetailedBrokerRow[];
  profile_flow?: ProfileFlowRow[];
  profile_broker_detail?: ProfileBrokerDetailRow[];
}

export interface CausalityData {
  granger_test?: {
    is_significant: boolean;
    min_p_value: number;
    best_lag: number;
  } | null;
  participant_causality?: Array<{
    participant: string;
    lag: number;
    p_value: number;
    is_significant: boolean;
  }>;
  top_brokers?: Array<{
    code: string;
    lag: number;
    p_value: number;
    is_significant: boolean;
  }>;
}

export interface BrokerScanValidationRow {
  ticker: string;
  broker_code: string;
  n_events: number;
  win_rate: number;
  mean_fwd_return: number;
  median_fwd_return: number;
  total_net_value: number;
  p_value_one_sided?: number | null;
  significant: boolean;
}

export interface EventStudyChartPoint {
  day: string;
  median: number;
  range?: [number, number];
  [key: string]: string | number | [number, number] | undefined;
}

export interface EventStudyTableRow {
  ticker: string;
  signal_date: string;
  signal: string;
  t_plus_1d?: number | null;
  t_plus_3d?: number | null;
  t_plus_5d?: number | null;
  t_plus_10d?: number | null;
}

export interface ValidationV2Data {
  broker_scan?: {
    ticker?: BrokerScanValidationRow[];
    all?: BrokerScanValidationRow[];
  };
  event_study?: {
    chart?: EventStudyChartPoint[];
    paths?: Array<{ id: string; data: Record<string, number> }>;
    table?: EventStudyTableRow[];
  };
}

export interface RawFlowRow {
  date: string;
  signal: string;
  score: number;
  foreign_net: number;
  local_net: number;
  total_value: number;
}

export interface RawActivityRow {
  date: string;
  broker: string;
  type: string;
  buy: number;
  sell: number;
  net: number;
  freq: number;
}

export interface RawTablesData {
  flow?: RawFlowRow[];
  activity?: RawActivityRow[];
}
