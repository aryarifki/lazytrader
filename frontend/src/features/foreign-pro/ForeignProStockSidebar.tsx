import ForeignProStatisticsBox from "./ForeignProStatisticsBox";
import type {
  AkumulasiResult,
  ForeignPowerInfo,
  MasterTickerItem,
  StockSeriesResponse,
} from "./foreignProTypes";

export interface ForeignProStockSidebarProps {
  ticker: string;
  tickerSearch: string;
  setTickerSearch: (val: string) => void;
  suggestOpen: boolean;
  setSuggestOpen: (open: boolean) => void;
  suggestions: MasterTickerItem[];
  onSelectTicker: (tk: string) => void;
  onRefresh: () => void;
  stockData: StockSeriesResponse | null;
  latestFp: ForeignPowerInfo;
  showFF: boolean;
  setShowFF: (v: boolean) => void;
  showFCo: boolean;
  setShowFCo: (v: boolean) => void;
  showKF: boolean;
  setShowKF: (v: boolean) => void;
  showXL: boolean;
  setShowXL: (v: boolean) => void;
  showSemiRitel: boolean;
  setShowSemiRitel: (v: boolean) => void;
  valFFText: string;
  valFCoText: string;
  valKFText: string;
  valXLText: string;
  valSRText: string;
  showMA: boolean;
  setShowMA: (v: boolean) => void;
  maPeriod: number;
  setMaPeriod: (v: number) => void;
  maTarget: string;
  setMaTarget: (v: string) => void;
  additionalOverlay: string;
  setAdditionalOverlay: (v: string) => void;
  akumDays: number;
  setAkumDays: (v: number) => void;
  akumulasiResult: AkumulasiResult;
  statType: string;
  setStatType: (v: string) => void;
  signalStat: string;
  setSignalStat: (v: string) => void;
}

export default function ForeignProStockSidebar(props: ForeignProStockSidebarProps) {
  const {
    ticker,
    tickerSearch,
    setTickerSearch,
    suggestOpen,
    setSuggestOpen,
    suggestions,
    onSelectTicker,
    onRefresh,
    stockData,
    latestFp,
    showFF,
    setShowFF,
    showFCo,
    setShowFCo,
    showKF,
    setShowKF,
    showXL,
    setShowXL,
    showSemiRitel,
    setShowSemiRitel,
    valFFText,
    valFCoText,
    valKFText,
    valXLText,
    valSRText,
    showMA,
    setShowMA,
    maPeriod,
    setMaPeriod,
    maTarget,
    setMaTarget,
    additionalOverlay,
    setAdditionalOverlay,
    akumDays,
    setAkumDays,
    akumulasiResult,
    statType,
    setStatType,
    signalStat,
    setSignalStat,
  } = props;

  const an = stockData?.analytics;
  const priceChange = an?.price_change || 0;

  const mainIndicators = [
    { id: "chkFF", valId: "valFF", label: "Foreign Flow", checked: showFF, set: setShowFF, color: "#FF6820", text: valFFText },
    { id: "chkFCoV2", valId: "valFCoV2", label: "Foreign & Co", checked: showFCo, set: setShowFCo, color: "#D62839", text: valFCoText },
  ];

  const additionalIndicators = [
    { id: "chkKF", valId: "valKF", label: "Korban Flow", checked: showKF, set: setShowKF, color: "#2F7676", text: valKFText },
    { id: "chkXL", valId: "valXL", label: "Broker XL", checked: showXL, set: setShowXL, color: "#3A7BD5", text: valXLText },
    { id: "chkSemiRitel", valId: "valSemiRitel", label: "Semi Ritel", checked: showSemiRitel, set: setShowSemiRitel, color: "#1DB87A", text: valSRText },
  ];

  return (
    <div className="hidden lg:block lg:col-span-3 space-y-3">
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-xl p-3 shadow-md space-y-3 font-mono">
        <div className="flex items-center justify-between pb-2 border-b border-[var(--md-sys-color-outline-variant)]/40">
          <div className="flex items-center gap-1.5 flex-wrap">
            <span className="w-2.5 h-2.5 rounded-full bg-[var(--md-sys-color-primary)]" />
            <span id="tickerName" className="text-base font-black text-[var(--md-sys-color-on-surface)]">
              {ticker}
            </span>
            <div id="tickerBadges" className="flex items-center gap-1">
              <span className="px-1.5 py-0.5 rounded text-[10px] font-bold bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-primary-container)]">
                {latestFp.badge}
              </span>
              {stockData?.badges?.map((b) => (
                <span
                  key={b}
                  className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                >
                  {b}
                </span>
              ))}
            </div>
          </div>

          <div className="text-right">
            <div id="lastPrice" className="text-sm font-bold text-[var(--md-sys-color-on-surface)]">
              Rp {an?.last_price?.toLocaleString("id-ID") || "—"}
            </div>
            <div
              id="lastChange"
              className={`text-[10px] font-bold ${
                priceChange >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
              }`}
            >
              {priceChange > 0 ? `+${priceChange}%` : `${priceChange}%`}
            </div>
          </div>
        </div>

        <div className="relative flex items-center gap-1.5">
          <input
            type="text"
            id="tickerSearch"
            value={tickerSearch}
            onChange={(e) => {
              setTickerSearch(e.target.value.toUpperCase());
              setSuggestOpen(true);
            }}
            onKeyDown={(e) => {
              if (e.key === "Enter" && tickerSearch) onSelectTicker(tickerSearch);
              else if (e.key === "Escape") setSuggestOpen(false);
            }}
            onFocus={() => setSuggestOpen(true)}
            onBlur={() => setTimeout(() => setSuggestOpen(false), 150)}
            placeholder="Cari Ticker..."
            className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/60 rounded px-2.5 py-1 text-xs text-[var(--md-sys-color-on-surface)] uppercase focus:outline-none focus:border-[var(--md-sys-color-primary)]"
          />
          <button
            id="btnRefreshData"
            type="button"
            onClick={onRefresh}
            title="Refresh data"
            className="p-1.5 bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)]/50 rounded text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] cursor-pointer"
          >
            ↻
          </button>

          {suggestOpen && suggestions.length > 0 && (
            <div
              id="tickerSuggestions"
              className="absolute top-full left-0 mt-1 w-full bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)]/70 rounded-md shadow-2xl z-50 max-h-48 overflow-y-auto"
            >
              {suggestions.map((t) => (
                <div
                  key={t.ticker}
                  onMouseDown={(e) => {
                    e.preventDefault();
                    onSelectTicker(t.ticker);
                  }}
                  className="px-2.5 py-1.5 text-xs text-[var(--md-sys-color-on-surface)] hover:bg-[var(--md-sys-color-primary)] hover:text-[var(--md-sys-color-on-primary)] cursor-pointer flex justify-between"
                >
                  <span>{t.ticker}</span>
                  <span className="opacity-70 text-[10px]">Rp {t.last_close}</span>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="space-y-1.5 pt-1">
          <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
            Main Indicator
          </div>
          {mainIndicators.map((ind) => (
            <label key={ind.id} className="flex items-center justify-between text-xs cursor-pointer p-1 rounded hover:bg-[var(--md-sys-color-surface-container-high)]">
              <div className="flex items-center gap-2">
                <input type="checkbox" id={ind.id} checked={ind.checked} onChange={(e) => ind.set(e.target.checked)} style={{ accentColor: ind.color }} />
                <span className="text-[var(--md-sys-color-on-surface)]">{ind.label}</span>
              </div>
              <span id={ind.valId} className="font-bold text-[10px] font-mono text-right" style={{ color: ind.color }}>
                {ind.text}
              </span>
            </label>
          ))}
        </div>

        <div className="space-y-1.5 pt-1 border-t border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
            Additional Indicator
          </div>
          {additionalIndicators.map((ind) => (
            <label key={ind.id} className="flex items-center justify-between text-xs cursor-pointer p-1 rounded hover:bg-[var(--md-sys-color-surface-container-high)]">
              <div className="flex items-center gap-2">
                <input type="checkbox" id={ind.id} checked={ind.checked} onChange={(e) => ind.set(e.target.checked)} style={{ accentColor: ind.color }} />
                <span className="text-[var(--md-sys-color-on-surface)]">{ind.label}</span>
              </div>
              <span id={ind.valId} className="font-bold text-[10px] font-mono text-right" style={{ color: ind.color }}>
                {ind.text}
              </span>
            </label>
          ))}
        </div>

        <div className="space-y-1.5 pt-1 border-t border-[var(--md-sys-color-outline-variant)]/40">
          <label className="flex items-center gap-2 text-xs cursor-pointer">
            <input type="checkbox" id="chkMA" checked={showMA} onChange={(e) => setShowMA(e.target.checked)} className="accent-[var(--md-sys-color-primary)]" />
            <span className="text-[var(--md-sys-color-on-surface)] font-bold">Moving Average</span>
          </label>
          {showMA && (
            <div className="grid grid-cols-2 gap-1.5 text-xs pt-1">
              <select id="maPeriodSelect" value={maPeriod} onChange={(e) => setMaPeriod(Number(e.target.value))} className="bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/50 rounded p-1 text-[var(--md-sys-color-on-surface)] text-[11px]">
                <option value={10}>10 hari</option>
                <option value={20}>20 hari</option>
                <option value={50}>50 hari</option>
                <option value={100}>100 hari</option>
              </select>
              <select id="maTargetSelect" value={maTarget} onChange={(e) => setMaTarget(e.target.value)} className="bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/50 rounded p-1 text-[var(--md-sys-color-on-surface)] text-[11px]">
                <option value="ff">Foreign Flow</option>
                <option value="fco">Foreign &amp; Co</option>
                <option value="kf">Korban Flow</option>
                <option value="xl">XL</option>
                <option value="semiritel">Semi Ritel</option>
              </select>
            </div>
          )}
        </div>

        <div className="pt-1 border-t border-[var(--md-sys-color-outline-variant)]/40 space-y-1">
          <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">Add Indicator</div>
          <select id="additionalOverlaySelect" value={additionalOverlay} onChange={(e) => setAdditionalOverlay(e.target.value)} className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/50 rounded p-1 text-xs text-[var(--md-sys-color-on-surface)]">
            <option value="none">Add Indicator</option>
            <option value="localNonRitel">Local Non Ritel</option>
            <option value="kepRitel">% Kep Ritel</option>
            <option value="kepHiddenF">% Kep Hidden F</option>
            <option value="kepForeign">% Kep Foreign</option>
            <option value="korelasiKF">Korelasi Korban Flow</option>
            <option value="fcoV2">F&amp;Co v2 (Asing+Lokal)</option>
          </select>
        </div>

        <div className="pt-1 border-t border-[var(--md-sys-color-outline-variant)]/40 space-y-1.5">
          <div className="flex items-center justify-between text-xs">
            <span className="text-[var(--md-sys-color-on-surface-variant)] text-[11px]">Average Acc:</span>
            <div className="flex items-center gap-1">
              <input type="number" id="akumDaysInput" min={1} max={100} value={akumDays} onChange={(e) => setAkumDays(Math.max(1, parseInt(e.target.value, 10) || 1))} className="w-12 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/50 rounded px-1 py-0.5 text-center text-xs text-[var(--md-sys-color-on-surface)]" />
              <span className="text-[var(--md-sys-color-on-surface-variant)] text-[10px]">hari</span>
            </div>
          </div>

          <div id="akumCombined" className="bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/40 rounded p-2 text-center text-xs">
            <span className="text-[var(--md-sys-color-on-surface)] font-bold">
              Price : Rp {akumulasiResult.avgPrice?.toLocaleString("id-ID") || "—"}{" "}
              <span className={`text-[10px] ${akumulasiResult.isAcc ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"}`}>
                ({akumulasiResult.isAcc ? "Acc" : "Dist"})
              </span>
            </span>
            <span className="text-[var(--md-sys-color-on-surface-variant)] mx-1.5">|</span>
            <span className={`font-bold ${akumulasiResult.totalFlow >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"}`}>
              {akumulasiResult.totalFlow >= 0 ? "In +" : "Out "}
              {akumulasiResult.totalFlow.toFixed(1)} M
            </span>
          </div>
        </div>

        <ForeignProStatisticsBox
          statType={statType}
          setStatType={setStatType}
          signalStat={signalStat}
          setSignalStat={setSignalStat}
          stats={an?.statistics}
        />
      </div>
    </div>
  );
}
