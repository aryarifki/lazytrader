export type IpoMainView = "analisis" | "list" | "screener";
export type IpoRightTab = "tx" | "info";
export type IpoAnalysisMode = "after" | "periodic";

export interface IpoApiResponse<T> {
  success: boolean;
  data: T;
}

export interface IpoQuickPick {
  ticker: string;
  company_name: string;
  listing_date: string;
  ipo_price: number;
}

export interface IpoUnderwriterCount {
  broker: string;
  count: number;
}

export interface IpoMasterTicker {
  ticker: string;
  company_name: string;
  listing_date: string | null;
  ipo_price: number;
}

export interface IpoMeta {
  total_ipos: number;
  d0_date: string | null;
  earliest_date: string | null;
  latest_date: string | null;
  available_dates: string[];
  quick_picks: IpoQuickPick[];
  years: number[];
  underwriters: IpoUnderwriterCount[];
  tickers: IpoMasterTicker[];
}

export interface IpoCandle {
  date: string;
  open: number | null;
  high: number | null;
  low: number | null;
  close: number | null;
  prev: number | null;
  tavg: number | null;
  volume: number | null;
  net_volume: number | null;
}

export interface IpoDonutSlice {
  key: string;
  label: string;
  color: string;
  value: number;
  pct?: number;
  pct_ipo?: number;
  pct_ring?: number;
  top_brokers?: string;
  sub_brokers?: string;
}

export interface IpoDonutSide {
  total: number;
  slices: IpoDonutSlice[];
}

export interface IpoPieModePair {
  buyer: IpoDonutSide;
  seller: IpoDonutSide;
}

export interface IpoBuyerSellerPieCard {
  net: IpoPieModePair;
  gross: IpoPieModePair;
}

export interface IpoDistSide {
  available: boolean;
  empty?: boolean;
  ipo_lot?: number;
  slices: IpoDonutSlice[];
  total: number;
}

export interface IpoAfterIpoDistCard {
  enabled: boolean;
  sell: IpoDistSide;
  buy: IpoDistSide;
}

export interface IpoTbsRow {
  broker: string;
  net: number;
  value_jt: number | null;
  category_label: string;
  style_key: string;
  is_uw: boolean;
}

export interface IpoTopBuyerSellerCard {
  buyers: IpoTbsRow[];
  sellers: IpoTbsRow[];
  acc_dist: number[];
  value_price: number | null;
}

export interface IpoUnderwriterRow {
  broker: string;
  role: string;
  net: number;
  pct_net: number | null;
  trans: number;
  pct_trans: number | null;
}

export interface IpoUnderwriterCard {
  has_uw: boolean;
  total_net: number;
  total_pct_net: number | null;
  total_trans: number;
  total_pct_trans: number | null;
  lot_ipo: number | null;
  pct_of_ipo: number | null;
  status: string;
  rows: IpoUnderwriterRow[];
}

export interface IpoCategoryRow {
  key: string;
  label: string;
  bandar_brokers: string[];
  net: number;
  pct_net: number | null;
  trans: number;
  pct_trans: number | null;
}

export interface IpoCategoryTotal {
  net: number;
  pct_net: number | null;
  trans: number;
  pct_trans: number | null;
}

export interface IpoCategoryCard {
  rows: IpoCategoryRow[];
  total: IpoCategoryTotal;
}

export interface IpoActionLabel {
  label: string;
  cls: string;
}

export interface IpoRetailRow {
  broker: string;
  net: number;
  pct_net: number | null;
  trans: number;
  pct_trans: number | null;
  is_auto_bandar: boolean;
  is_active_bandar: boolean;
}

export interface IpoRetailTotalExclAuto {
  sum_trans: number;
  sum_net: number;
  pct_trans: number | null;
  pct_net: number | null;
}

export interface IpoRetailActionCard {
  retail_ratio: number | null;
  retail_action: IpoActionLabel;
  retail_bandar_auto_result: string | null;
  retail_bandar_broker: string | null;
  retail_bandar_brokers: string[];
  retail_bandar_reasons: Record<string, string[]>;
  effective_retail_list: string[];
  rows: IpoRetailRow[];
  total_excl_auto: IpoRetailTotalExclAuto;
}

export interface IpoGapRange {
  from: string;
  to: string;
}

export interface IpoInfoCard {
  ticker: string;
  company_name: string;
  listing_date: string | null;
  ipo_price: number;
  last_close: number | null;
  vwap: number | null;
  ipo_vs_close_pct: number | null;
  avg_vs_close_pct: number | null;
  lot_offer: number | null;
  lot_total: number | null;
  pct_ipo: string | null;
  nilai_miliar: number | null;
  penjatahan: number | null;
  papan: string | null;
  jumlah_sid: number | null;
  trading_days_since_ipo: number | null;
  holiday_warn: boolean;
  max_holiday_year: number;
  gap: IpoGapRange | null;
  no_d0: boolean;
  in_d0: boolean;
  d0_date: string | null;
  n_later: number;
}

export interface IpoAnalysisPeriod {
  start: string | null;
  end: string | null;
  day_count: number;
  price_change: number | null;
  avg_daily_value: number | null;
  total_volume: number;
  tot_net_volume: number;
  tot_buy: number;
  tot_sell: number;
}

export interface IpoAnalysisCards {
  top_buyer_seller: IpoTopBuyerSellerCard;
  underwriter_card: IpoUnderwriterCard;
  category_card: IpoCategoryCard;
  after_ipo_distribution: IpoAfterIpoDistCard;
  buyer_seller_pie: IpoBuyerSellerPieCard;
  retail_action_card: IpoRetailActionCard;
  info_card: IpoInfoCard;
}

export interface IpoAnalysisNoDataMeta {
  ticker: string;
  company_name: string;
  listing_date: string | null;
  ipo_price: number;
}

export interface IpoAnalysisData {
  ticker: string;
  mode: IpoAnalysisMode;
  has_data: boolean;
  master?: IpoAnalysisNoDataMeta;
  period?: IpoAnalysisPeriod;
  candles: IpoCandle[];
  available_ticker_dates?: string[];
  cards?: IpoAnalysisCards;
}

export interface IpoListUnderwriterChip {
  broker: string;
  role: string;
  category: string;
}

export interface IpoListRow {
  ticker: string;
  company_name: string;
  listing_date: string | null;
  ipo_price: number;
  nilai_ipo: number | null;
  nilai_miliar: number | null;
  papan: string | null;
  underwriters: IpoListUnderwriterChip[];
  jumlah_sid: number | null;
  penjatahan: number | null;
  day_returns: Array<number | null>;
  perf_2w: number | null;
  perf_1m: number | null;
  perf_3m: number | null;
}

export interface IpoKinerjaSummary {
  rata_rata: Record<string, number | null>;
  kemungkinan_naik: Record<string, number | null>;
}

export interface IpoListData {
  total: number;
  total_all: number;
  rows: IpoListRow[];
  kinerja_summary: IpoKinerjaSummary;
}

export interface IpoScreenerCondition {
  k: string;
  op: string;
  v: string;
}

export interface IpoScreenerRow {
  kode: string;
  company_name: string;
  listing_date: string | null;
  uw: number | null;
  netRetail: number | null;
  netSemi: number | null;
  netInst: number | null;
  netNonRitel: number | null;
  netRepo: number | null;
  perfIpo: number | null;
  perfAvg: number | null;
  chg: number | null;
  ratrans: number | null;
  volRatio: number | null;
}

export interface IpoScreenerData {
  rows: IpoScreenerRow[];
  total: number;
  n_days: number;
  lo: string | null;
  hi: string | null;
  after: boolean;
  active_col: string;
  sort_key: string;
  sort_dir: number;
}
