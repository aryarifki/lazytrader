export interface IndexPoint {
  date: string;
  close: number;
  volume: number;
  value: number;
  market_cap: number | null;
}

export interface SectorRow {
  code: string;
  name: string;
  close: number;
  change_1d_pct: number | null;
  change_1w_pct: number | null;
  change_1m_pct: number | null;
  date: string;
}

export interface IhsgDonutItem {
  label: string;
  value_trn?: number;
  count?: number;
  pct: number;
  color: string;
}

export interface IhsgSummaryRow {
  category: string;
  prev_trn: number | null;
  curr_trn: number | null;
  diff_trn: number | null;
  pct: number | null;
}

export interface IhsgTop20Item {
  rank: number;
  ticker: string;
  issuer: string;
  foreign_named_trn: number | null;
  foreign_named_pct: number | null;
  foreign_hidden_trn: number | null;
  foreign_hidden_pct: number | null;
  local_nonret_trn: number | null;
  local_nonret_pct: number | null;
  retail_trn: number | null;
  retail_pct: number | null;
  free_float: number | null;
  badges?: string[];
  badge?: string | null;
  conglomerate?: string | null;
  konglo?: string | null;
}

export interface IhsgComposition {
  months: Record<string, any>;
  summary_comparison: IhsgSummaryRow[];
  donuts: {
    scripless_vs_warkat: IhsgDonutItem[];
    investor_type: IhsgDonutItem[];
    foreign_distribution: IhsgDonutItem[];
  };
  top20: IhsgTop20Item[];
  total_mktcap_trn: number;
  warkat_trn: number;
  n_tickers: number;
  n_foreign_ge50: number;
}

export interface MarketOverview {
  index: {
    code: string;
    name: string;
    date: string;
    close: number;
    change_1d_pct: number | null;
    return_1w_pct: number | null;
    return_1m_pct: number | null;
    return_3m_pct: number | null;
    return_6m_pct: number | null;
    return_1y_pct: number | null;
    return_ytd_pct: number | null;
    high_52w: number | null;
    low_52w: number | null;
    market_cap: number | null;
    value: number | null;
    volume: number | null;
  };
  history: IndexPoint[];
  sectors: SectorRow[];
  breadth: {
    advancers: number;
    decliners: number;
    unchanged: number;
    total: number;
  };
  foreign_flow: {
    date: string | null;
    "1d": number | null;
    "1w": number | null;
    "1m": number | null;
  };
  composition?: IhsgComposition | null;
  ownership_coverage: {
    covered_tickers: number;
    total_tickers: number;
    period_date: string | null;
    composition_available: boolean;
    top20_available: boolean;
    note: string;
  };
  data_quality: {
    notes: string[];
  };
}
