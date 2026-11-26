export type FFMainTab = "stock" | "ihsg" | "screener";
export type ScreenerSubTab = "basic" | "custom" | "manual" | "admin";
export type AdminReportMode = "sesi1" | "end";
export type HistTopMode = "ff" | "fco";
export type UniverseKey = "ff" | "fco" | "kf" | "xl" | "semiritel";

export interface MasterTickerItem {
  ticker: string;
  name?: string;
  company_name?: string;
  last_close?: number;
}

export interface SpecialCandle {
  date: string;
  type: "PURPLE" | "RED" | "YELLOW" | "ORANGE" | string;
}

export interface ForeignPowerInfo {
  score: number;
  badge: string;
  color: string;
}

export interface CorrelationStats {
  ff?: number;
  fco?: number;
  kf?: number;
  xl?: number;
  semiritel?: number;
}

export interface ConsecutiveInflowStats {
  "50d"?: number | string;
  "100d"?: number | string;
  "200d"?: number | string;
}

export interface InflowHargaNaikEntry {
  count?: number;
  pct?: number;
}

export interface SignalHorizonStats {
  events?: number;
  horizons?: {
    "5"?: { pct?: number };
    "20"?: { pct?: number };
  };
}

export interface StockAnalyticsStatistics {
  correlations?: CorrelationStats;
  consecutive_inflow?: ConsecutiveInflowStats;
  inflow_harga_naik?: {
    "50d"?: InflowHargaNaikEntry;
    "100d"?: InflowHargaNaikEntry;
    "200d"?: InflowHargaNaikEntry;
  };
  mark_down?: {
    "50d"?: number;
    "100d"?: number;
    "200d"?: number;
  };
  signal_stats?: {
    massive_acc?: SignalHorizonStats;
  };
}

export interface StockAnalytics {
  dates?: string[];
  closes?: number[];
  opens?: number[];
  highs?: number[];
  lows?: number[];
  values_miliar?: number[];
  ff_daily_val?: number[];
  fco_daily_val?: number[];
  kf_daily_val?: number[];
  xl_daily_val?: number[];
  semiritel_daily_val?: number[];
  cum_lnr?: number[];
  partisipasi_series?: number[];
  kepemilikan_asing_total?: number[];
  kepemilikan_foreign_hidden?: number[];
  kepemilikan_ritel?: number[];
  correlation_20?: Array<number | null>;
  correlation_kf_20?: Array<number | null>;
  special_candles?: SpecialCandle[];
  latest_fp?: ForeignPowerInfo;
  last_price?: number;
  price_change?: number;
  statistics?: StockAnalyticsStatistics;
  highlights?: string[];
}

export interface StockSeriesResponse {
  ticker?: string;
  badges?: string[];
  analytics?: StockAnalytics;
}

export interface IhsgPointMetric {
  ff?: number;
  fco?: number;
  kf?: number;
  xl?: number;
  sr?: number;
  [key: string]: number | undefined;
}

export interface IhsgPoint {
  date: string;
  close: number;
  ihsg?: IhsgPointMetric;
  konglo?: IhsgPointMetric;
  fundamental?: IhsgPointMetric;
  [key: string]: unknown;
}

export interface IhsgSummary {
  date?: string;
  close?: number;
  changePct?: number;
  ffValue?: number;
  status?: string;
}

export interface IhsgResponse {
  points?: IhsgPoint[];
  summary?: IhsgSummary;
}

export interface ScreenerRow {
  ticker: string;
  name: string;
  close: number;
  changePct: number;
  fpBadge: string;
  accumVal: number;
  avgTransValue10d: number;
  badges?: string[];
  dailyFFValue?: number[];
  dailyFCoValue?: number[];
  dailyKFValue?: number[];
  dailyXLValue?: number[];
  dailySemiRitelValue?: number[];
  dailyClose?: number[];
  dailyValue?: number[];
  summary?: Record<string, number | undefined>;
}

export interface ScreenerResponse {
  screener?: ScreenerRow[];
}

export interface MarketReportItem {
  ticker: string;
  val: number;
  changePct?: number;
}

export interface MarketReportData {
  topForeignBuy: MarketReportItem[];
  topForeignSell: MarketReportItem[];
  topGainer: MarketReportItem[];
  topLoser: MarketReportItem[];
}

export interface PeriodReportData {
  topBuy: MarketReportItem[];
  topSell: MarketReportItem[];
}

export interface AkumulasiResult {
  avgPrice: number | null;
  totalFlow: number;
  isAcc: boolean;
}

export interface IhsgMacroMetrics {
  ff: string;
  fco: string;
  kf: string;
  xl: string;
  sr: string;
}

export interface StockChartInput {
  stockData: StockSeriesResponse | null;
  isLightMode: boolean;
  effectiveWindowStart: number;
  effectiveWindowEnd: number;
  datesSlice: string[];
  opensSlice: number[];
  closesSlice: number[];
  lowsSlice: number[];
  highsSlice: number[];
  cumFFSlice: number[];
  cumFCoSlice: number[];
  cumKFSlice: number[];
  cumXLSlice: number[];
  cumSRSlice: number[];
  cumLNRSlice: number[];
  showFF: boolean;
  showFCo: boolean;
  showKF: boolean;
  showXL: boolean;
  showSemiRitel: boolean;
  showMA: boolean;
  maPeriod: number;
  maTarget: string;
  additionalOverlay: string;
  showPBNyangkut: boolean;
  histTopMode: HistTopMode;
  histBottomMode: string;
}

export interface IhsgChartInput {
  ihsgData: IhsgResponse | null;
  isLightMode: boolean;
  catIhsg: boolean;
  catKonglo: boolean;
  catFund: boolean;
  metFF: boolean;
  metFCo: boolean;
  metKF: boolean;
  metXL: boolean;
  metSR: boolean;
  flipKF: boolean;
  flipXL: boolean;
  flipSR: boolean;
}
