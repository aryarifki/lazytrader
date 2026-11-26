import type React from "react";

export type Mode = "stocks" | "investors";
export type Category = "kepemilikan" | "transaksi" | "advanced";
export type SortOrder = "asc" | "desc";

export interface ConditionOption {
  id: string;
  label: string;
}

export interface ExplorerFilters {
  ownership_ranges: string[];
  sid_ranges: string[];
  value_ranges: string[];
  mktcap_tiers: Array<{ id: string; label: string }>;
  retail_ranges: string[];
  portfolio_ranges: string[];
  investor_origins: string[];
  investor_types: Record<string, string[]>;
}

export interface ExplorerPeriod {
  last_update_date: string | null;
  month_prev: string;
  month_curr: string;
  gl_month: string;
  month_labels: Record<string, string>;
  month_keys: string[];
}

export interface ExplorerMatch {
  ticker: string;
  issuer?: string;
  delta: number;
}

export interface ExplorerTickerValue {
  ticker: string;
  net_value_m: number;
}

export interface ExplorerRow {
  ticker?: string;
  issuer?: string;
  name?: string;
  badges?: string[];
  badge?: string | null;
  conglomerate?: string | null;
  konglo?: string | null;
  lf?: string;
  lf_label?: string;
  market_cap_b?: number | null;
  price?: number | null;
  close?: number | null;
  mtd_pct?: number | null;
  net_pct?: number | null;
  holder_count?: number | null;
  sid_change?: number | null;
  sid?: number | null;
  sid_prev?: number | null;
  period_prev?: string;
  period_last?: string;
  net_value_m?: number | null;
  count?: number | null;
  pct_change?: number | null;
  latest?: number | null;
  holder_name?: string;
  best?: number | null;
  match_count?: number | null;
  matches?: ExplorerMatch[];
  ticker_count?: number | null;
  tickers?: string[] | ExplorerTickerValue[];
  abs_value_m?: number | null;
  gl_rp_m?: number | null;
  pct?: number | null;
}

export interface ExplorerResponse {
  mode: Mode;
  category: Category;
  condition: string;
  condition_label: string;
  count: number;
  total: number;
  limit: number;
  offset: number;
  results: ExplorerRow[];
  conditions: ConditionOption[];
  sorts: Array<{ id: string; label: string }>;
  filters: ExplorerFilters;
  period: ExplorerPeriod;
  source: string;
  data_quality: {
    ct_baseline: boolean;
    periodized_holdings: boolean;
    t5_transactions: boolean;
    gl: boolean;
    market_cap: boolean;
    notes: string[];
  };
}

export interface Column {
  key: string;
  label: string;
  align?: "right";
  sort?: string;
  render: (row: ExplorerRow, index: number) => React.ReactNode;
}
