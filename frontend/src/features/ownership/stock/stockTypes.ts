export interface SearchResult {
  ticker: string;
  name: string;
  sector: string;
}

export interface PriceSnapshot {
  date: string | null;
  close: number;
  previous_close: number | null;
  change_1d_pct: number | null;
  month_start_close: number | null;
  mtd_pct: number | null;
}

export interface CompositionItem {
  key: string;
  label: string;
  pct: number;
  color: string;
  description: string;
}

export interface Composition {
  basis:
    | "monthly_ownership_metrics"
    | "company_shareholders_snapshot"
    | "ct_baseline_pct_history"
    | "staging_normalized";
  period_label: string;
  period_date: string | null;
  total_pct: number;
  items: CompositionItem[];
  note: string;
}

export interface StagingCompositionItem {
  label: string;
  kind: string;
  shares: number | null;
  pct: number | null;
  color: string | null;
}

export interface StagingComposition {
  report_date: string;
  items: StagingCompositionItem[];
  source: string;
}

export interface Holder {
  name: string;
  pct: number;
  is_controlling: boolean;
  category: string;
}

export interface StagingHolderItem {
  name: string;
  investor_type: string | null;
  classification: string | null;
  location: string | null;
  nationality: string | null;
  pct: number;
  total_shares: number | null;
  previous_pct: number | null;
  previous_total_shares: number | null;
  change_shares: number;
  change_pct: number;
  status: string;
}

export interface StagingHolders {
  report_date: string;
  holders: StagingHolderItem[];
  source: string;
}

export interface StagingMovementItem {
  date: string | null;
  name: string;
  action: string;
  change_shares: number | null;
  change_pct: number | null;
  current_shares: number | null;
  current_pct: number | null;
  price: number | null;
  value: number | null;
  value_m: number | null;
  broker: string | null;
  source_label: string | null;
}

export interface StagingMovementsSummary {
  count: number;
  buy_count: number;
  sell_count: number;
  buy_value_m: number;
  sell_value_m: number;
  net_value_m: number;
}

export interface StagingUboSubholding {
  name: string;
  pct_to?: number;
}

export interface StagingUboItem {
  name: string;
  pct: number;
  previous_pct: number | null;
  subholdings: StagingUboSubholding[];
  report_date: string | null;
}

export interface CtHolder {
  name: string;
  lf: string | null;
  sl_mei: number | null;
  pct_mei: number | null;
  sl_jun: number | null;
  pct_jun: number | null;
  status: string | null;
}

export interface T5Transaction {
  ticker?: string;
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

export interface SidPoint {
  period: string;
  sid: number;
  change: number | null;
  change_pct: number | null;
}

export interface MonthlyPoint {
  period_label: string;
  period_date: string | null;
  total_sid: number | null;
  sid_change_mom: number | null;
  sid_change_pct: number | null;
  foreign_total_pct?: number | null;
  foreign_named_pct: number;
  foreign_hidden_pct: number;
  retail_pct: number;
  local_non_retail_pct: number;
}

export interface StockOwnership {
  ticker: string;
  name: string;
  sector: string;
  badges?: string[];
  badge?: string | null;
  conglomerate?: string | null;
  price: PriceSnapshot | null;
  profile: {
    market_cap: number | null;
    market_cap_b?: number | null;
    shares_outstanding: number | null;
    free_float_pct: number | null;
    ct_free_float_pct?: number | null;
    total_sid: number | null;
    smtracker_total_sid?: number | null;
    avg_retail_shares: number | null;
    conglomerate?: string | null;
    badge?: string | null;
    badges?: string[];
    scripless_lot?: number | null;
    pct_scrip?: number | null;
    source?: string;
  };
  summary: {
    holder_count: number;
    total_pct: number;
    controlling_pct: number;
    public_pct: number;
    treasury_pct: number;
    substantial_pct: number;
    management_pct: number;
    max_holder_pct: number;
    free_float_est: number;
  };
  composition: Composition;
  staging_composition?: StagingComposition | null;
  history: MonthlyPoint[];
  ct_composition?: Composition | null;
  ct_history?: MonthlyPoint[] | null;
  holders: Holder[];
  staging_holders?: StagingHolders | null;
  ct_holders?: CtHolder[] | null;
  ct_sid_history?: SidPoint[] | null;
  transactions?: T5Transaction[] | null;
  staging_movements?: StagingMovementItem[];
  staging_movements_summary?: StagingMovementsSummary | null;
  staging_ubo?: StagingUboItem[];
  t5_summary?: T5Summary | null;
  data_quality: {
    monthly_composition: boolean;
    periodized_holdings: boolean;
    t5_transactions: boolean;
    sid_history?: boolean;
    ct_baseline?: boolean;
    market_cap: boolean;
    notes: string[];
  };
  as_of: {
    price: string | null;
    monthly: string | null;
    shareholders: string;
    staging_holders?: string | null;
    staging_movements?: string | null;
    ct_baseline?: string | null;
  };
}

export interface MonthlyComparisonRow {
  period_label: string;
  isLatest: boolean;
  currForeign: number;
  currForeignNamed: number;
  currForeignHidden: number;
  foreignDelta: number | null;
  foreignHiddenDelta: number | null;
  currRetail: number;
  retailDelta: number | null;
  currLocalNonRetail: number;
  localNonRetailDelta: number | null;
  currSid: number | null;
  sidDelta: number | null;
  sidDeltaPct: number | null;
  flowSignal: string;
  flowTone: "up" | "down" | "primary" | "warning" | "neutral";
}
