interface ForeignProScreenerBasicProps {
  activeStrategyId: string;
  setActiveStrategyId: (id: string) => void;
  strategyDesc: string;
  setStrategyDesc: (desc: string) => void;
  activeReportType: string | null;
  setActiveReportType: (type: string | null) => void;
  periodReportWindow: number | null;
  setPeriodReportWindow: (w: number | null) => void;
}

const PERIOD_REPORTS = [
  { w: 5, label: "1 Minggu" },
  { w: 10, label: "2 Minggu" },
  { w: 20, label: "1 Bulan" },
  { w: 40, label: "2 Bulan" },
  { w: 60, label: "3 Bulan" },
  { w: 120, label: "6 Bulan" },
];

export default function ForeignProScreenerBasic({
  activeStrategyId,
  setActiveStrategyId,
  strategyDesc,
  setStrategyDesc,
  activeReportType,
  setActiveReportType,
  periodReportWindow,
  setPeriodReportWindow,
}: ForeignProScreenerBasicProps) {
  const activeBtn =
    "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]";
  const idleBtn =
    "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]";

  return (
    <div
      id="basicPanel"
      className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-3 rounded-xl space-y-3"
    >
      <div className="space-y-1.5">
        <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
          Basic Search
        </div>
        <div id="basicFixedBtnRow" className="flex items-center gap-2 flex-wrap">
          <button
            type="button"
            onClick={() => {
              setActiveStrategyId("basicHFP");
              setStrategyDesc("Saham dengan skor Foreign Power TINGGI (badge HFP, skor > 30).");
              setActiveReportType(null);
              setPeriodReportWindow(null);
            }}
            className={`px-3 py-1 rounded text-xs font-bold transition cursor-pointer ${
              activeStrategyId === "basicHFP" && !activeReportType && !periodReportWindow
                ? activeBtn
                : idleBtn
            }`}
          >
            HFP
          </button>
          <button
            type="button"
            onClick={() => {
              setActiveStrategyId("basicInflow50");
              setStrategyDesc("Foreign Flow (net beli asing) HARI INI lebih besar dari Rp 50 miliar.");
              setActiveReportType(null);
              setPeriodReportWindow(null);
            }}
            className={`px-3 py-1 rounded text-xs font-bold transition cursor-pointer ${
              activeStrategyId === "basicInflow50" && !activeReportType && !periodReportWindow
                ? activeBtn
                : idleBtn
            }`}
          >
            Inflow Hari Ini &gt; 50 M
          </button>
        </div>
      </div>

      <div className="space-y-1.5">
        <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
          Market Report Hari Ini
        </div>
        <div id="basicReportBtnRow" className="flex items-center gap-2 flex-wrap">
          <button
            type="button"
            onClick={() => {
              setActiveReportType("topBuySell");
              setPeriodReportWindow(null);
              setStrategyDesc(
                "Laporan 4 panel: Top 5 Foreign Buy & Top 5 Foreign Sell hari ini (Rp M), Top 5 Gainer & Top 5 Loser hari ini."
              );
            }}
            className={`px-3 py-1 rounded text-xs font-bold transition cursor-pointer ${
              activeReportType === "topBuySell" ? activeBtn : idleBtn
            }`}
          >
            Top Buy/Sell &amp; Gainer/Loser
          </button>
          {PERIOD_REPORTS.map((p) => (
            <button
              key={p.w}
              type="button"
              onClick={() => {
                setPeriodReportWindow(p.w);
                setActiveReportType("periodReport");
                setStrategyDesc(
                  `Laporan akumulasi periode ${p.label}: Top 5 Foreign Buy & Top 5 Foreign Sell.`
                );
              }}
              className={`px-3 py-1 rounded text-xs font-bold transition cursor-pointer ${
                periodReportWindow === p.w ? activeBtn : idleBtn
              }`}
            >
              {p.label}
            </button>
          ))}
        </div>
      </div>

      <div
        id="basicDesc"
        className="text-xs text-[var(--md-sys-color-on-surface-variant)] bg-[var(--md-sys-color-surface-container-high)] p-2.5 rounded border border-[var(--md-sys-color-outline-variant)]/40"
      >
        {strategyDesc}
      </div>
    </div>
  );
}
