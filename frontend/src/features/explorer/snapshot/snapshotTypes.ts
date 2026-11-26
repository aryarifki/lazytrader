export interface Condition {
  id: string;
  label: string;
  available: boolean;
  category?: string;
}

export interface ConditionGroup {
  id: string;
  label: string;
}

export interface SortOption {
  id: string;
  label: string;
}

export interface UnsupportedCondition {
  id: string;
  label: string;
  reason: string;
}

export interface StockRow {
  ticker: string;
  name: string;
  sector: string;
  badges?: string[];
  badge?: string | null;
  conglomerate?: string | null;
  konglo?: string | null;
  holder_count: number;
  total_pct: number;
  controlling_pct: number;
  public_pct: number;
  treasury_pct: number;
  substantial_pct: number;
  max_holder_pct: number;
  controlling_count: number;
  public_count: number;
  price: number | null;
  price_date: string | null;
  mtd_pct: number | null;
  foreign_5d: number | null;
  market_cap_b: number | null;
}

export interface InvestorRow {
  name: string;
  ticker_count: number;
  total_pct: number;
  max_pct: number;
  controlling_count: number;
  substantial_count: number;
  sector_count: number;
  sample_tickers: string[];
}

export interface ExplorerResponse {
  mode: "stocks" | "investors";
  condition: string;
  condition_label: string;
  count: number;
  total: number;
  limit: number;
  offset: number;
  results: StockRow[] | InvestorRow[];
  conditions: Condition[];
  condition_groups: ConditionGroup[];
  sorts: SortOption[];
  sectors: string[];
  unsupported_conditions: UnsupportedCondition[];
  data_quality: {
    source: string;
    periodized_holdings: boolean;
    t5_transactions: boolean;
    market_cap?: boolean;
    notes: string[];
  };
}

export type Mode = "stocks" | "investors";
export type SortOrder = "asc" | "desc";
export type MtdFilter = "all" | "up" | "down";
export type FlowFilter = "all" | "inflow" | "outflow";

export interface Filters {
  sector: string;
  minHolders: string;
  maxHolders: string;
  minControlling: string;
  maxControlling: string;
  minPublic: string;
  maxPublic: string;
  minTopHolder: string;
  maxTopHolder: string;
  minMarketCap: string;
  maxMarketCap: string;
  mtd: MtdFilter;
  foreignFlow: FlowFilter;
  hasTreasury: boolean;
  minTickers: string;
  maxTickers: string;
  minMaxPct: string;
  minSectorCount: string;
}

export interface ActiveFilterItem {
  key: string;
  label: string;
  clear: () => void;
}
