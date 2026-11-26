export interface GroupTopStock {
  ticker: string;
  name: string;
  mtd_pct: number | null;
}

export interface GroupSummary {
  id: string;
  name: string;
  description: string;
  hub_entity: string;
  logo_url?: string | null;
  catalog_code?: string | null;
  ticker_count: number;
  covered_tickers: number;
  avg_mtd_pct: number | null;
  advancers: number;
  decliners: number;
  public_pct_avg: number | null;
  controlling_pct_avg: number | null;
  holder_rows: number;
  monthly_covered: number;
  sid_total_covered: number | null;
  sid_total: number | null;
  top_stocks: GroupTopStock[];
  market_cap: number | null;
  market_cap_b: number | null;
  market_cap_available: boolean;
  market_cap_unit: string;
  ct_baseline: boolean;
}

export interface GroupsResponse {
  as_of: string | null;
  coverage: {
    mapped_groups: number;
    mapped_tickers: number;
    total_tickers: number;
    market_cap_available: boolean;
    t5_available: boolean;
    ct_baseline: boolean;
    note: string;
  };
  top_groups: GroupSummary[];
  groups: GroupSummary[];
  notes: string[];
}

export interface Constituent {
  ticker: string;
  name: string;
  sector: string;
  badges?: string[];
  badge?: string | null;
  conglomerate?: string | null;
  konglo?: string | null;
  price: number | null;
  price_date: string | null;
  change_1d_pct: number | null;
  mtd_pct: number | null;
  holder_count: number;
  controlling_pct: number;
  public_pct: number;
  max_holder_pct: number;
  free_float_est: number;
  market_cap_b: number | null;
  big_player_pct: number | null;
  foreign_total_pct: number | null;
  foreign_hidden_pct: number | null;
  retail_pct: number | null;
  sid: number | null;
  monthly: {
    period_label: string;
    total_sid: number;
    retail_pct: number;
    foreign_named_pct: number;
    foreign_hidden_pct: number;
    local_non_retail_pct: number;
  } | null;
}

export interface GroupHolding {
  ticker: string;
  pct: number;
  is_controlling: boolean;
}

export interface GroupTopHolding {
  name: string;
  ticker_count: number;
  total_pct: number;
  holdings: GroupHolding[];
}

export interface T5Transaction {
  ticker: string;
  date: string | null;
  broker: string | null;
  name: string | null;
  action: string | null;
  change_lot: number | null;
  value_m: number | null;
  lot_after: number | null;
  pct_juli: number | null;
}

export interface T5Summary {
  count: number;
  buy_count: number;
  sell_count: number;
  buy_value_m: number;
  sell_value_m: number;
  net_value_m: number;
}

export interface GroupDetail {
  group: GroupSummary;
  constituents: Constituent[];
  top_holdings: GroupTopHolding[];
  transactions: T5Transaction[];
  t5_summary: T5Summary;
  data_quality: {
    monthly_composition: boolean;
    market_cap: boolean;
    t5_transactions: boolean;
    sid_history: boolean;
    ct_baseline: boolean;
    notes: string[];
  };
}

export interface NetworkNode {
  id: string;
  label: string;
  name: string;
  type: "hub" | "owner" | "emiten";
  val?: number;
}

export interface NetworkLink {
  source: string;
  target: string;
  pct: number;
  is_controlling: boolean;
  label: string;
}

export interface NetworkData {
  group_id: string;
  group_name: string;
  hub_entity: string;
  nodes: NetworkNode[];
  links: NetworkLink[];
  stats: {
    total_nodes: number;
    total_emitens: number;
    total_links: number;
    hub_entity: string;
  };
}

export type GroupView = "jejaring" | "matriks" | "emiten";
