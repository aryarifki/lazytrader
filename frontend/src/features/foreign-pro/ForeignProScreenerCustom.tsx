interface ForeignProScreenerCustomProps {
  universe: string;
  setUniverse: (u: string) => void;
  customCategory: string;
  setCustomCategory: (c: string) => void;
  activeStrategyId: string;
  setActiveStrategyId: (id: string) => void;
  strategyDesc: string;
  setStrategyDesc: (desc: string) => void;
}

const UNIVERSES = [
  { key: "ff", label: "Foreign Flow" },
  { key: "fco", label: "F&Co" },
  { key: "kf", label: "Retail" },
  { key: "xl", label: "XL" },
  { key: "semiritel", label: "SR" },
];

const CATEGORIES = [
  { key: "top", label: "Inflow Terbesar" },
  { key: "consistency", label: "Inflow Konsisten" },
  { key: "divergence", label: "Divergence Harga vs Flow" },
  { key: "reversal", label: "Mark Up, Mark Down & Reversal" },
];

const STRATEGIES_BY_CATEGORY: Record<
  string,
  Array<{ id: string; label: string; desc: string }>
> = {
  top: [
    {
      id: "inflow3m",
      label: "Inflow Terbesar 3 Bulan",
      desc: "Hari ini mencatat Daily Inflow TERBESAR dalam 3 bulan terakhir (60 hari trading), dengan Foreign Power Medium/High dan Avg Transaksi 10D > 35 M.",
    },
    {
      id: "maxInflow5min3",
      label: "Inflow Terbesar 5 Hari (min 3x)",
      desc: "Hari ini adalah Daily Inflow TERBESAR dalam 5 hari terakhir, DAN minimal 3 dari 5 hari itu ada inflow, dengan Foreign Power Medium/High dan Avg Transaksi 10D > 35 M.",
    },
    {
      id: "top2day20",
      label: "Inflow 2 Hari Terbesar 20 Hari",
      desc: "Jumlah 2 hari terakhir adalah yang TERBESAR dibanding semua pasangan 2-hari lain dalam 20 hari terakhir.",
    },
    {
      id: "weeklyInflowMax",
      label: "Inflow Mingguan Terbesar 1 Bulan",
      desc: "Minggu ini (blok 5 hari terakhir) positif DAN lebih besar dari 3 minggu lainnya.",
    },
  ],
  consistency: [
    {
      id: "consec3",
      label: "3 Hari Berturut-turut Inflow",
      desc: "Net BELI (inflow) 3 hari berturut-turut terakhir, dengan Foreign Power Medium/High dan Avg Transaksi 10D > 35 M.",
    },
    {
      id: "consec5",
      label: "5 Hari Berturut-turut Inflow",
      desc: "Net BELI (inflow) 5 hari berturut-turut terakhir, dengan Foreign Power Medium/High dan Avg Transaksi 10D > 35 M.",
    },
    {
      id: "inflow3of5",
      label: "Minimal 3x Inflow dalam 5 Hari",
      desc: "Minimal 3 dari 5 hari terakhir mencatat net BELI (inflow, tidak harus berturut-turut).",
    },
  ],
  reversal: [
    {
      id: "firstInflow10",
      label: "Inflow Pertama dalam 10 Hari",
      desc: "Hari ini net BELI (inflow), padahal 9 hari sebelumnya SAMA SEKALI tidak ada inflow (sinyal reversal).",
    },
    {
      id: "bearishMarkUp",
      label: "Bearish Mark Up",
      desc: "Harga hari ini NAIK, tapi akumulasi 20 hari & 5 hari NEGATIF (net distribusi tersembunyi).",
    },
    {
      id: "bullishMarkDown",
      label: "Bullish Mark Down",
      desc: "Harga hari ini TURUN, tapi akumulasi 20 hari & 5 hari POSITIF (net akumulasi tersembunyi).",
    },
  ],
};

export default function ForeignProScreenerCustom({
  universe,
  setUniverse,
  customCategory,
  setCustomCategory,
  activeStrategyId,
  setActiveStrategyId,
  strategyDesc,
  setStrategyDesc,
}: ForeignProScreenerCustomProps) {
  const activeBtn =
    "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]";
  const idleBtn =
    "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]";
  const strategies = STRATEGIES_BY_CATEGORY[customCategory] || [];

  return (
    <div
      id="customPanel"
      className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-3 rounded-xl space-y-3"
    >
      <div className="space-y-1.5">
        <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
          Universe:
        </div>
        <div id="customUniverseBtnRow" className="flex items-center gap-2 flex-wrap">
          {UNIVERSES.map((u) => (
            <button
              key={u.key}
              type="button"
              onClick={() => setUniverse(u.key)}
              className={`px-3 py-1 rounded text-xs font-bold transition cursor-pointer ${
                universe === u.key ? activeBtn : idleBtn
              }`}
            >
              {u.label}
            </button>
          ))}
        </div>
      </div>

      <div className="space-y-1.5">
        <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
          Category
        </div>
        <div id="customCategoryBtnRow" className="flex items-center gap-2 flex-wrap">
          {CATEGORIES.map((c) => (
            <button
              key={c.key}
              type="button"
              onClick={() => setCustomCategory(c.key)}
              className={`px-3 py-1 rounded text-xs font-bold transition cursor-pointer ${
                customCategory === c.key
                  ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border border-[var(--md-sys-color-primary)]/40"
                  : idleBtn
              }`}
            >
              {c.label}
            </button>
          ))}
        </div>
      </div>

      <div
        id="customSearchCategories"
        className="pt-2 border-t border-[var(--md-sys-color-outline-variant)]/40 flex items-center gap-2 flex-wrap"
      >
        {strategies.map((st) => (
          <button
            key={st.id}
            type="button"
            onClick={() => {
              setActiveStrategyId(st.id);
              setStrategyDesc(st.desc);
            }}
            className={`px-3 py-1 rounded text-xs font-bold transition cursor-pointer ${
              activeStrategyId === st.id ? activeBtn : idleBtn
            }`}
          >
            {st.label}
          </button>
        ))}
      </div>

      <div
        id="customDesc"
        className="text-xs text-[var(--md-sys-color-on-surface-variant)] bg-[var(--md-sys-color-surface-container-high)] p-2.5 rounded border border-[var(--md-sys-color-outline-variant)]/40"
      >
        {strategyDesc}
      </div>
    </div>
  );
}
