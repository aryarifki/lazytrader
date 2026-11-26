export interface InvestorSearchItem {
  name: string;
  lf: string;
  lf_label: string;
  ticker_count: number;
  total_pct: number;
  total: number | null;
  in_index: boolean;
}

export interface InvestorSearchResponse {
  query: string;
  count: number;
  total_matches: number;
  results: InvestorSearchItem[];
  source: string;
  data_quality: {
    ct_baseline: boolean;
    includes_names_missing_from_index: boolean;
    notes: string[];
  };
}

export interface InvestorHolding {
  ticker: string;
  issuer: string;
  badges?: string[];
  badge?: string | null;
  conglomerate?: string | null;
  konglo?: string | null;
  lf: string;
  sl_mei: number | null;
  pct_mei: number | null;
  sl_jun: number | null;
  pct_jun: number | null;
  status: "new" | "out" | "both";
  price: number | null;
  market_cap_b: number | null;
  value_m: number | null;
  delta_lot: number | null;
  delta_pct: number | null;
  is_treasury: boolean;
}

export interface InvestorTransaction {
  ticker: string;
  date: string;
  broker: string;
  name: string;
  action: string;
  change_lot: number | null;
  value_m: number | null;
  lot_after: number | null;
  pct_juli: number | null;
}

export interface InvestorDetail {
  name: string;
  lf: string;
  lf_label: string;
  summary: {
    ticker_count: number;
    total_value_m: number | null;
    gl_rp_m: number | null;
    gl_pct: number | null;
    buy_value_m: number;
    buy_count: number;
    sell_value_m: number;
    sell_count: number;
    net_value_m: number;
    peak_basis_m: number | null;
  };
  period: {
    last_update_date: string | null;
    month_prev: string;
    month_curr: string;
    gl_month: string;
  };
  holdings: InvestorHolding[];
  transactions: InvestorTransaction[];
  t5_summary: {
    count: number;
    buy_count?: number;
    sell_count?: number;
    buy_value_m?: number;
    sell_value_m?: number;
    net_value_m?: number;
  };
  data_quality: {
    periodized_holdings: boolean;
    t5_transactions: boolean;
    gl: boolean;
    market_cap: boolean;
    ct_baseline: boolean;
    notes: string[];
  };
}
