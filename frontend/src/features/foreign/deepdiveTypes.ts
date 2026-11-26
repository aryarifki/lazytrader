import type { CostBasisTimeSeries } from "./charts/ForeignCostBasisChart";
import type { CounterpartyData } from "./charts/CounterpartyFlowChart";
import type { HeatmapData } from "./charts/BrokerHeatmap";
import type { KseiHistoryItem } from "./charts/KseiRealityChart";

export interface BrokerEntry {
  broker_code: string;
  broker_name: string;
  net_value: number;
  net_lot: number;
  buy_value: number;
  sell_value: number;
  avg_buy_price: number;
  avg_sell_price: number;
  floating_pnl_pct: number;
  active_days: number;
}

export interface T5Entry {
  date: string;
  investor_name: string;
  broker: string;
  action: string;
  change_lot: number;
  value_m: number;
  price: number;
}

export interface CompanyInfo {
  name?: string;
  group?: string;
}

export interface CostBasisInfo {
  avg_buy_price?: number;
  floating_pnl_pct?: number;
  total_net_value?: number;
  position_status?: string;
  inventory_retention_pct?: number;
  latest_close?: number;
  top_holders?: BrokerEntry[];
  top_sellers?: BrokerEntry[];
}

export interface ImpactDivergenceInfo {
  divergence_tone?: string;
  streak_days?: number;
  divergence_label?: string;
  summary_text?: string;
  crossing_warning?: boolean;
  max_single_day_share_pct?: number;
  foreign_power_badge?: string;
  foreign_power_score?: number | null;
  avg_participation_pct?: number;
  foreign_price_corr?: number;
  buy_days?: number;
  sell_days?: number;
  consistency_pct?: number;
}

export interface CounterpartyInfo extends CounterpartyData {
  matchup_tone?: string;
  matchup_verdict?: string;
  totals?: {
    foreign?: number;
    local_inst?: number;
    retail?: number;
  };
}

export interface KseiRealityInfo {
  latest_pct_asing?: number | null;
  alignment_status?: string;
  delta_1m_pp?: number | null;
  history?: KseiHistoryItem[] | null;
  latest_period?: string;
  pct_asing_gt1?: number;
  pct_asing_lt1?: number;
  latest_pct_ritel_pure?: number;
  t5_recent?: T5Entry[];
}

export interface ForeignFeaturesInfo {
  foreign_hhi?: number;
  foreign_zscore?: number;
}

export interface ForeignModelsInfo {
  broker_heatmap?: HeatmapData | null;
}

export interface ForeignFlowAnalyticsResponse {
  latest_date?: string;
  company?: CompanyInfo;
  cost_basis?: CostBasisInfo;
  impact_divergence?: ImpactDivergenceInfo;
  counterparty?: CounterpartyInfo;
  ksei_reality?: KseiRealityInfo;
  features?: ForeignFeaturesInfo;
  models?: ForeignModelsInfo;
  timeseries?: CostBasisTimeSeries | null;
}

export interface UniverseAllResponse {
  tickers?: string[];
}
