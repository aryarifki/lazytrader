export interface ActiveSignal {
  id: string;
  ticker: string;
  company_name: string;
  group_name: string;
  signal_date: string;
  entry_price: number;
  tp_price: number;
  sl_price: number;
  target_gain_pct: number;
  stop_risk_pct: number;
  risk_reward_ratio: number;
  status: string;
  days_held: number;
  confirmation_count?: number;
  last_confirmed_date?: string | null;
  composite_score: number;
  wyckoff_phase: string;
  wyckoff_name?: string;
  smart_money_notes?: string;
  sector_notes?: string;
  persistence_5d?: number;
  turnover_idr?: number;
  ml_win_prob?: number | null;
  cmf: number | null;
  mfi: number | null;
  rsi: number | null;
  vcp_grade: string;
  gate_notes: string;
  features: Record<string, unknown>;
}

export interface DailySignalRow {
  ticker: string;
  date: string;
  ml_label: string;
  ml_win_prob: number | null;
  composite_score: number;
  scores: {
    technical: number;
    smart_money: number;
    sector: number;
    fundamental: number;
  };
  features: {
    wyckoff_phase?: string;
    wyckoff_name?: string;
    smart_money_notes?: string;
    sector_notes?: string;
    persistence_5d?: number;
    turnover_idr?: number;
    rr_ratio?: number;
    wp?: string;
    cmf_v?: number;
    mfi_v?: number;
    rsi_v?: number;
    vcp_grade?: string;
    [key: string]: unknown;
  };
  wyckoff_phase?: string;
  wyckoff_name?: string;
  smart_money_notes?: string;
  sector_notes?: string;
  persistence_5d?: number;
  turnover_idr?: number;
  rr_ratio?: number;
  gate_notes: string;
}

export interface ConfluenceRow {
  ticker: string;
  company_name: string;
  group_name: string;
  current_price: number;
  signal: string;
  conviction_score: number;
  bandar_avg_price: number;
  bandar_price_distance_pct: number;
  foreign_5d: number;
  smart_foreign_net: number;
  gorengan_net: number;
  foreign_hmm_state: string;
  free_float_est: number;
  top_buyer: string;
  top_seller: string;
  quick_tags: string[];
  total_value?: number;
}

export type ConfluenceSortKey = "conviction" | "distance" | "foreign" | "gorengan";

export interface BacktestStats {
  total_signals: number;
  total_closed: number;
  total_open: number;
  win_count: number;
  loss_count: number;
  win_rate: number;
  avg_win_pct: number;
  avg_loss_pct: number;
  avg_days_held: number;
  profit_factor: number;
  latest_audit_date: string | null;
}

export interface BacktestSignalItem {
  id: string;
  ticker: string;
  signal_date: string;
  entry_price: number;
  tp_price: number;
  sl_price: number;
  exit_price: number | null;
  exit_date: string | null;
  status: "WIN" | "LOSS" | "OPEN" | string;
  days_held: number;
  pnl_pct: number | null;
  composite_score: number;
  sm_score: number;
  sec_score: number;
}

export interface HistoryResponse {
  status_filter: string;
  count: number;
  data: BacktestSignalItem[];
}
