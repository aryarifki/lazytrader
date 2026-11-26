export const CONSTANT_CONDITION_LABELS: Record<string, string> = {
  bertambah: "Net Buy / Bertambah",
  berkurang: "Net Sell / Berkurang",
  sidNaik: "SID Bertambah",
  sidTurun: "SID Berkurang",
  naikPorsi: "Big Player Menambah Porsi Bulan Ini",
  aktifBeli: "Pembelian Aktif oleh Big Player",
  buyback: "Emiten Buyback Bulan Ini",
  mayoritas: "Pemegang Saham Mayoritas Menambah Porsi",
  baru: "Masuknya Big Player Baru Bulan Ini",
  retailTurun: "Kepemilikan Retail Menurun",
  hiddenNaik: "Hidden Big Player Meningkat",
  foreignNaik: "Kepemilikan Foreign Meningkat",
  masukFlat: "Big Player Masuk (Harga Belum Naik Signifikan)",
  foreign3Bulan: "Foreign Naik 3 Bulan Beruntun",
  retail3Bulan: "Retail Keluar 3 Bulan Beruntun",
  sid3Turun: "SID Turun 3 Bulan Beruntun",
  sid3Naik: "SID Naik 3 Bulan Beruntun",
  retailRendah: "Kepemilikan Retail Rendah",
  portofolioNaik: "Portofolio Investor Naik Bulan Ini",
  portofolioTurun: "Portofolio Investor Turun Bulan Ini",
  investorAktif: "Investor Aktif Bertransaksi Sebulan Terakhir",
  investorNambah: "Investor Menambah Porsi Sebulan Terakhir",
};

export const STOCK_ADVANCED_IDS = [
  "naikPorsi",
  "aktifBeli",
  "buyback",
  "mayoritas",
  "baru",
  "retailTurun",
  "hiddenNaik",
  "foreignNaik",
  "masukFlat",
  "foreign3Bulan",
  "retail3Bulan",
  "sid3Turun",
  "sid3Naik",
  "retailRendah",
];

export const INVESTOR_ADVANCED_IDS = [
  "portofolioNaik",
  "portofolioTurun",
  "investorAktif",
  "investorNambah",
];

export const RETAIL_RANGES = ["<1", "1-3", "3-5"];
export const PORTFOLIO_RANGES = ["<3", "3-10", ">10"];
export const OWNERSHIP_RANGE_OPTIONS = [
  { id: ">10", label: ">10%" },
  { id: "5-10", label: "5% - 10%" },
  { id: "1-5", label: "1% - 5%" },
  { id: "<1", label: "<1%" },
];
export const SID_RANGE_OPTIONS = [
  { id: "<1", label: "<1%" },
  { id: "1-3", label: "1% - 3%" },
  { id: "3-10", label: "3% - 10%" },
  { id: ">10", label: ">10%" },
];
export const VALUE_RANGE_OPTIONS = [
  { id: "<5M", label: "< 5M" },
  { id: "5-100M", label: "5M - 100M" },
  { id: ">100M", label: "> 100M" },
];
export const RETAIL_RANGE_OPTIONS = [
  { id: "<1", label: "<1%" },
  { id: "1-3", label: "1% - 3%" },
  { id: "3-5", label: "3% - 5%" },
];
export const PORTFOLIO_RANGE_OPTIONS = [
  { id: "<3", label: "<3%" },
  { id: "3-10", label: "3% - 10%" },
  { id: ">10", label: ">10%" },
];
export const MKTCAP_TIERS = [
  { id: "big", label: "Big Cap", title: ">50T" },
  { id: "mid", label: "Mid Cap", title: "5T-50T" },
  { id: "midsmall", label: "Mid-Small Cap", title: "1T-5T" },
  { id: "small", label: "Small Cap", title: "300M-1T" },
  { id: "micro", label: "Micro Cap", title: "<300M" },
];
export const LOCAL_TYPES = [
  { id: "ID", label: "Individual" },
  { id: "CP", label: "Corporate" },
  { id: "IB", label: "Financial Institution" },
  { id: "OTHER", label: "Others" },
];
export const FOREIGN_TYPES = [
  { id: "CP", label: "Corporate" },
  { id: "IB", label: "Financial Institution" },
  { id: "ID", label: "Individual" },
  { id: "MF", label: "Mutual Fund" },
  { id: "OTHER", label: "Others" },
];
