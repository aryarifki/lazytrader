import type { Filters, SortOption } from "./snapshotTypes";

export const EMPTY_FILTERS: Filters = {
  sector: "",
  minHolders: "",
  maxHolders: "",
  minControlling: "",
  maxControlling: "",
  minPublic: "",
  maxPublic: "",
  minTopHolder: "",
  maxTopHolder: "",
  minMarketCap: "",
  maxMarketCap: "",
  mtd: "all",
  foreignFlow: "all",
  hasTreasury: false,
  minTickers: "",
  maxTickers: "",
  minMaxPct: "",
  minSectorCount: "",
};

export const PAGE_SIZES = [25, 50, 100, 200];

export const STOCK_SORT_FALLBACK: SortOption[] = [
  { id: "default", label: "Default (top holder)" },
  { id: "ticker", label: "Ticker" },
  { id: "name", label: "Nama emiten" },
  { id: "holders", label: "Jumlah pemegang" },
  { id: "controlling", label: "Porsi pengendali" },
  { id: "public", label: "Porsi publik" },
  { id: "top_holder", label: "Porsi top holder" },
  { id: "substantial", label: "Porsi substantial" },
  { id: "treasury", label: "Porsi treasury" },
  { id: "market_cap", label: "Market cap" },
  { id: "mtd", label: "Perubahan MTD" },
  { id: "foreign_5d", label: "Foreign flow 5D" },
];

export const INVESTOR_SORT_FALLBACK: SortOption[] = [
  { id: "default", label: "Default (jumlah emiten)" },
  { id: "name", label: "Nama investor" },
  { id: "tickers", label: "Jumlah emiten" },
  { id: "total_pct", label: "Total porsi" },
  { id: "max_pct", label: "Porsi maksimum" },
  { id: "controlling", label: "Jumlah pengendali" },
  { id: "substantial", label: "Jumlah substantial" },
  { id: "sectors", label: "Jumlah sektor" },
];

export const STOCK_PRESETS = [
  { id: "controlling_high", label: "Pengendali >50%", icon: "ph:crown-duotone" },
  { id: "public_low", label: "Publik <10%", icon: "ph:users-duotone" },
  { id: "many_holders", label: "Pemegang ≥20", icon: "ph:users-three-duotone" },
  { id: "treasury", label: "Treasury", icon: "ph:bank-duotone" },
  { id: "cross_holding", label: "Cross holding", icon: "ph:share-network-duotone" },
  { id: "foreign_inflow_5d", label: "Foreign inflow", icon: "ph:arrows-down-up-duotone" },
  { id: "price_up_mtd", label: "Harga naik MTD", icon: "ph:trend-up-duotone" },
];

export const INVESTOR_PRESETS = [
  { id: "multi_holding", label: "≥5 emiten", icon: "ph:buildings-duotone" },
  { id: "controlling", label: "Punya pengendali", icon: "ph:crown-duotone" },
  { id: "substantial", label: "Porsi ≥5%", icon: "ph:chart-pie-slice-duotone" },
  { id: "cross_sector", label: "Lintas sektor", icon: "ph:share-network-duotone" },
];
