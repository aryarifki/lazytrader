import { useState, useEffect, useCallback } from "react";
import { useSearchParams } from "react-router-dom";
import { Icon } from "@iconify/react";
import {
  getMmdTickers, getMmdDetail, getMmdAnalysis, getMmdT5,
  getMmdBrokerDrilldown, getStockOwnershipProfile, scanMmdExplorer,
} from "@/lib/api/terminal";
import type { StockOwnership } from "@/features/ownership/stock/stockTypes";
import { useAppStore } from "@/store/useAppStore";
import MmdHeroStats from "./MmdHeroStats";
import MmdChartCard from "./MmdChartCard";
import MmdBrokerTable from "./MmdBrokerTable";
import MmdQuickModeCard from "./MmdQuickModeCard";
import MmdBandarResultTab from "./MmdBandarResultTab";
import MmdBandarDetailTab from "./MmdBandarDetailTab";
import MmdBandarOtherTab from "./MmdBandarOtherTab";
import MmdBigPlayerView from "./MmdBigPlayerView";
import MmdExplorerView from "./MmdExplorerView";
import MmdParameterDrawer from "./MmdParameterDrawer";
import MmdBrokerModal from "./MmdBrokerModal";
import type {
  BandarTabMode, ExplorerCategory, MainViewMode, MmdAnalysis,
  MmdBrokerDrilldown, MmdDetailResponse, MmdExplorerRow, MmdT5Item, MmdTicker, MobileBrokerTab,
} from "./mmdTypes";

const NAV_TABS = [
  { id: "main" as const, icon: "ph:gauge-duotone", label: "MMD Pro" },
  { id: "bigPlayer" as const, icon: "ph:crown-duotone", label: "Big Player" },
  { id: "explorer" as const, icon: "ph:magnifying-glass-duotone", label: "Explorer" },
];

const BANDAR_TABS = [
  { id: "resume" as const, label: "Result" },
  { id: "detail" as const, label: "Detail" },
  { id: "other" as const, label: "Other" },
];

export default function MmdProPage() {
  const [searchParams] = useSearchParams();
  const queryTicker = searchParams.get("ticker");
  const activeTicker = useAppStore((s) => s.activeTicker);
  const setActiveTicker = useAppStore((s) => s.setActiveTicker);

  const [activeView, setActiveView] = useState<MainViewMode>("main");
  const [activeBandarTab, setActiveBandarTab] = useState<BandarTabMode>("resume");
  const [mobileBrokerTab, setMobileBrokerTab] = useState<MobileBrokerTab>("buyer");
  const [parameterDrawerOpen, setParameterDrawerOpen] = useState(false);
  const [quickModeActive, setQuickModeActive] = useState(false);
  const [showT5Markers, setShowT5Markers] = useState(true);
  const [showTooltip, setShowTooltip] = useState(true);

  const [retailBandarAutoEnabled, setRetailBandarAutoEnabled] = useState(true);
  const [retailBandarManualOverride, setRetailBandarManualOverride] = useState("");

  const [tickers, setTickers] = useState<MmdTicker[]>([]);
  const [selectedTicker, setSelectedTicker] = useState<string>(queryTicker?.toUpperCase() || activeTicker || "BBCA");
  const [tickerInput, setTickerInput] = useState<string>(selectedTicker);
  const [suggestOpen, setSuggestOpen] = useState(false);
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");

  const [candleData, setCandleData] = useState<MmdDetailResponse | null>(null);
  const [analysisData, setAnalysisData] = useState<MmdAnalysis | null>(null);
  const [t5Data, setT5Data] = useState<MmdT5Item[]>([]);
  const [ownership, setOwnership] = useState<StockOwnership | null>(null);
  const [loading, setLoading] = useState(true);

  const [brokerModalOpen, setBrokerModalOpen] = useState(false);
  const [brokerModalLoading, setBrokerModalLoading] = useState(false);
  const [brokerModalData, setBrokerModalData] = useState<MmdBrokerDrilldown | null>(null);

  const [explorerCat, setExplorerCat] = useState<ExplorerCategory>("normal");
  const [explorerCriteria, setExplorerCriteria] = useState("ALL");
  const [explorerVolumeBucket, setExplorerVolumeBucket] = useState(0);
  const [explorerResults, setExplorerResults] = useState<MmdExplorerRow[]>([]);
  const [explorerLoading, setExplorerLoading] = useState(false);

  useEffect(() => {
    if (queryTicker && queryTicker.toUpperCase() !== selectedTicker) {
      const up = queryTicker.toUpperCase();
      setSelectedTicker(up);
      setTickerInput(up);
      setActiveTicker(up);
      setDateFrom("");
      setDateTo("");
    }
  }, [queryTicker, selectedTicker, setActiveTicker]);

  useEffect(() => {
    let cancelled = false;
    getMmdTickers<{ tickers?: MmdTicker[] }>()
      .then((res) => { if (!cancelled && res.tickers?.length) setTickers(res.tickers); })
      .catch((err) => console.error("Error loading tickers:", err));
    return () => { cancelled = true; };
  }, []);

  useEffect(() => {
    if (!selectedTicker) return;
    let cancelled = false;
    setLoading(true);
    Promise.all([
      getMmdDetail<MmdDetailResponse>(selectedTicker, 60),
      getMmdAnalysis<{ analysis?: MmdAnalysis }>(selectedTicker, dateFrom || undefined, dateTo || undefined, retailBandarAutoEnabled, retailBandarManualOverride || undefined),
      getMmdT5<{ transactions?: MmdT5Item[]; t5_records?: MmdT5Item[] }>(selectedTicker),
      getStockOwnershipProfile(selectedTicker),
    ])
      .then(([detailRes, analysisRes, t5Res, ownershipRes]) => {
        if (cancelled) return;
        setCandleData(detailRes);
        setAnalysisData(analysisRes?.analysis || null);
        setT5Data(t5Res?.transactions || t5Res?.t5_records || []);
        setOwnership(ownershipRes);
        if (detailRes?.candles?.length) {
          const dates = detailRes.candles.map((c) => c.date_iso);
          if (!dateFrom || !dates.includes(dateFrom)) setDateFrom(dates[0]);
          if (!dateTo || !dates.includes(dateTo)) setDateTo(dates[dates.length - 1]);
        }
        setLoading(false);
      })
      .catch((err) => {
        if (cancelled) return;
        console.error("Error loading MMD data:", err);
        setLoading(false);
      });
    return () => { cancelled = true; };
  }, [selectedTicker, dateFrom, dateTo, retailBandarAutoEnabled, retailBandarManualOverride]);

  const executeScan = useCallback(
    (cat = explorerCat, crit = explorerCriteria, vol = explorerVolumeBucket) => {
      setExplorerLoading(true);
      scanMmdExplorer<{ results?: MmdExplorerRow[] }>(cat, crit, vol)
        .then((res) => { setExplorerResults(res.results || []); setExplorerLoading(false); })
        .catch((err) => { console.error("Error scanning explorer:", err); setExplorerLoading(false); });
    },
    [explorerCat, explorerCriteria, explorerVolumeBucket]
  );

  useEffect(() => { if (activeView === "explorer") executeScan(); }, [activeView, executeScan]);

  const handleSelectTicker = useCallback((tk: string) => {
    const upper = tk.trim().toUpperCase();
    if (!upper) return;
    setSelectedTicker(upper);
    setTickerInput(upper);
    setActiveTicker(upper);
    setDateFrom("");
    setDateTo("");
    setSuggestOpen(false);
  }, [setActiveTicker]);

  const handleOpenBrokerModal = useCallback(async (brokerCode: string) => {
    if (!brokerCode || brokerCode === "—") return;
    setBrokerModalOpen(true);
    setBrokerModalLoading(true);
    try {
      const res = await getMmdBrokerDrilldown<MmdBrokerDrilldown>(selectedTicker, brokerCode, dateFrom || undefined, dateTo || undefined);
      setBrokerModalData(res);
    } catch (err) {
      console.error("Error loading broker drilldown:", err);
      setBrokerModalData(null);
    } finally {
      setBrokerModalLoading(false);
    }
  }, [selectedTicker, dateFrom, dateTo]);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between gap-2 border-b border-[var(--md-sys-color-outline-variant)] pb-3 flex-wrap sm:flex-nowrap">
        <div className="flex items-center bg-[var(--md-sys-color-surface-container)] p-1 rounded-xl border border-[var(--md-sys-color-outline-variant)] shadow-2xs">
          {NAV_TABS.map((item) => (
            <button
              key={item.id} type="button" onClick={() => setActiveView(item.id)}
              className={`px-3 py-1.5 rounded-lg text-xs font-bold transition flex items-center gap-1.5 cursor-pointer ${
                activeView === item.id
                  ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] shadow-2xs font-extrabold"
                  : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
              }`}
            >
              <Icon icon={item.icon} width={16} />
              <span>{item.label}</span>
            </button>
          ))}
        </div>

        <div className="flex items-center gap-2">
          {loading && activeView !== "explorer" && (
            <span className="text-[10px] font-mono text-[var(--md-sys-color-primary)] animate-pulse">Syncing...</span>
          )}
          <button
            type="button" onClick={() => setQuickModeActive(!quickModeActive)} title="Toggle Quick Mode (6 Periode)"
            className={`flex items-center gap-1 px-2.5 py-1.5 rounded-lg text-xs font-mono font-bold border transition-all cursor-pointer ${
              quickModeActive
                ? "bg-amber-400 text-stone-950 border-amber-500 dark:bg-[#FFC24D] dark:text-black dark:border-[#FFC24D] shadow-xs"
                : "bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)] hover:text-[var(--md-sys-color-on-surface)]"
            }`}
          >
            <Icon icon="ph:lightning-duotone" width={14} />
            <span className="hidden xs:inline">Quick Mode</span>
          </button>
          <button
            type="button" onClick={() => setParameterDrawerOpen(true)} title="Atur Periode & Parameter Model"
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)] cursor-pointer shadow-2xs transition"
          >
            <Icon icon="ph:sliders-horizontal-duotone" width={16} />
            <span>Parameter</span>
          </button>
        </div>
      </div>

      {activeView === "main" && (
        <div className="space-y-4">
          <MmdHeroStats
            tickerInput={tickerInput} setTickerInput={setTickerInput}
            suggestOpen={suggestOpen} setSuggestOpen={setSuggestOpen}
            tickers={tickers} onSelectTicker={handleSelectTicker}
            analysisData={analysisData} dateFrom={dateFrom} dateTo={dateTo}
          />
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
            <div className="lg:col-span-7 space-y-4">
              <MmdChartCard
                candleData={candleData} t5Data={t5Data}
                showT5Markers={showT5Markers} setShowT5Markers={setShowT5Markers}
                showTooltip={showTooltip} setShowTooltip={setShowTooltip}
              />
              <MmdBrokerTable
                analysisData={analysisData} mobileBrokerTab={mobileBrokerTab}
                setMobileBrokerTab={setMobileBrokerTab} onOpenBroker={handleOpenBrokerModal}
              />
              <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-3.5 shadow-2xs space-y-2">
                <div className="text-xs font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-[var(--md-sys-color-primary)]" />
                  Rangkuman Pergerakan Big Player
                </div>
                <div className="space-y-1.5 text-xs text-[var(--md-sys-color-on-surface-variant)] bg-[var(--md-sys-color-surface-container-low)] p-3 rounded-xl border border-[var(--md-sys-color-outline-variant)]/40">
                  {analysisData?.big_player_summary?.length ? (
                    analysisData.big_player_summary.map((txt, idx) => (
                      <p key={idx} className="leading-relaxed flex items-start gap-2">
                        <span className="text-[var(--md-sys-color-primary)] font-bold">›</span>
                        <span>{txt}</span>
                      </p>
                    ))
                  ) : (
                    <p>Memuat ringkasan big player...</p>
                  )}
                </div>
              </div>
            </div>

            <div className="lg:col-span-5 space-y-4">
              {quickModeActive && <MmdQuickModeCard analysisData={analysisData} />}
              <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl overflow-hidden shadow-2xs">
                <div className="bg-[var(--md-sys-color-surface-container-high)] px-3.5 pt-2.5 pb-0 border-b border-[var(--md-sys-color-outline-variant)] flex items-center justify-between">
                  <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] uppercase tracking-wider">
                    Bandarmologi Analysis
                  </span>
                  <div className="flex items-center gap-1">
                    {BANDAR_TABS.map((tab) => (
                      <button
                        key={tab.id} type="button" onClick={() => setActiveBandarTab(tab.id)}
                        className={`px-3 py-1.5 text-xs font-bold rounded-t-xl transition border-b-2 cursor-pointer ${
                          activeBandarTab === tab.id
                            ? "border-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-surface)] bg-[var(--md-sys-color-surface-container)]"
                            : "border-transparent text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                        }`}
                      >
                        {tab.label}
                      </button>
                    ))}
                  </div>
                </div>
                <div className="p-3.5">
                  {activeBandarTab === "resume" && (
                    <MmdBandarResultTab analysisData={analysisData} selectedTicker={selectedTicker} onSelectTicker={handleSelectTicker} onOpenBroker={handleOpenBrokerModal} />
                  )}
                  {activeBandarTab === "detail" && (
                    <MmdBandarDetailTab
                      analysisData={analysisData}
                      retailBandarAutoEnabled={retailBandarAutoEnabled} setRetailBandarAutoEnabled={setRetailBandarAutoEnabled}
                      retailBandarManualOverride={retailBandarManualOverride} setRetailBandarManualOverride={setRetailBandarManualOverride}
                    />
                  )}
                  {activeBandarTab === "other" && (
                    <MmdBandarOtherTab analysisData={analysisData} ownership={ownership} selectedTicker={selectedTicker} />
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {activeView === "bigPlayer" && (
        <MmdBigPlayerView
          selectedTicker={selectedTicker} candleData={candleData} t5Data={t5Data}
          analysisData={analysisData} ownership={ownership} showT5Markers={showT5Markers} showTooltip={showTooltip}
        />
      )}

      {activeView === "explorer" && (
        <MmdExplorerView
          explorerCat={explorerCat} setExplorerCat={setExplorerCat}
          explorerCriteria={explorerCriteria} setExplorerCriteria={setExplorerCriteria}
          explorerVolumeBucket={explorerVolumeBucket} setExplorerVolumeBucket={setExplorerVolumeBucket}
          explorerResults={explorerResults} explorerLoading={explorerLoading}
          onScan={executeScan}
          onSelectTicker={(tk) => { handleSelectTicker(tk); setActiveView("main"); }}
        />
      )}

      <MmdParameterDrawer
        open={parameterDrawerOpen} onClose={() => setParameterDrawerOpen(false)} candleData={candleData}
        dateFrom={dateFrom} setDateFrom={setDateFrom} dateTo={dateTo} setDateTo={setDateTo}
        showT5Markers={showT5Markers} setShowT5Markers={setShowT5Markers}
        showTooltip={showTooltip} setShowTooltip={setShowTooltip}
        quickModeActive={quickModeActive} setQuickModeActive={setQuickModeActive}
        retailBandarAutoEnabled={retailBandarAutoEnabled} setRetailBandarAutoEnabled={setRetailBandarAutoEnabled}
      />
      <MmdBrokerModal open={brokerModalOpen} loading={brokerModalLoading} data={brokerModalData} onClose={() => setBrokerModalOpen(false)} />
    </div>
  );
}
