export type MarketState =
  | "NORMAL"
  | "WARNING"
  | "CRASH"
  | "BOTTOMING"
  | "RECOVERY_TESTING"
  | "CONFIRMED_BOTTOM";

export interface TimelinePoint {
  date: string;
  ihsg_close: number;
  daily_flow: number;
  cumulative_flow: number;
  market_state?: MarketState;
}

export interface CulpritBroker {
  broker_code: string;
  broker_name: string;
  net_sell_idr: number;
}

export interface RadarData {
  market_state: MarketState;
  ihsg_current: number;
  ihsg_change_pct: number;
  timeline: TimelinePoint[];
  culprit_brokers: CulpritBroker[];
}

export interface TickerCardItem {
  ticker: string;
  signal?: string | null;
  close?: number | null;
  ret_5d?: number | null;
  foreign_net_5d?: number | null;
  spark?: number[];
}

export interface WatchlistDetailResponse {
  ticker: string;
  signal_raw?: string | null;
  signal?: string | null;
  close?: number | null;
  ret_5d?: number | null;
  foreign_5d?: number | null;
  price_chart?: Array<{ close: number }>;
  error?: string;
}
