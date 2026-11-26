export interface BenchmarkInfo {
  code: string;
  name: string;
  close: number;
  ret_1d: number;
  value: number;
  market_state?: "NORMAL" | "WARNING" | "CRASH" | "BOTTOMING" | "RECOVERY_TESTING" | "CONFIRMED_BOTTOM";
}

export interface SectorItem {
  code: string;
  name: string;
  db_sector: string;
  icon: string;
  color: string;
  close: number;
  ret_1d: number;
  ret_5d: number;
  ret_20d: number;
  value: number;
  value_share_pct: number;
  foreign_net_1d: number;
  foreign_net_5d: number;
  foreign_net_20d: number;
  foreign_state: "INFLOW SURGE" | "HIDDEN DISTRIBUTION" | "SECTOR EXODUS" | "BOTTOMING" | "NORMAL";
  foreign_state_label: string;
  rs_ratio: number;
  rs_momentum: number;
  quadrant: "Leading" | "Weakening" | "Lagging" | "Improving";
  top_stocks: string[];
}

export interface SectorSummaryResponse {
  date: string;
  benchmark: BenchmarkInfo;
  macro_market_state?: "NORMAL" | "WARNING" | "CRASH" | "BOTTOMING" | "RECOVERY_TESTING" | "CONFIRMED_BOTTOM";
  total_sector_value: number;
  total_foreign_net_1d?: number;
  sectors: SectorItem[];
}

export interface SectorHistoryPoint {
  date: string;
  close: number;
  value: number;
  foreign_net: number;
  cumulative_foreign_flow: number;
  rs_ratio: number;
}

export interface BrokerMover {
  broker_code: string;
  broker_name: string;
  net_value: number;
}

export interface SectorHistoryResponse {
  sector: string;
  name: string;
  count: number;
  data: SectorHistoryPoint[];
  savior_brokers: BrokerMover[];
  culprit_brokers: BrokerMover[];
}
