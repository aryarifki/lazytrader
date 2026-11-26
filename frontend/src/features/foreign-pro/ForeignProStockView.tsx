import { useMemo } from "react";
import ReactECharts from "echarts-for-react";
import ForeignProStockHero from "./ForeignProStockHero";
import ForeignProStockSidebar from "./ForeignProStockSidebar";
import ForeignProStockChartFooter from "./ForeignProStockChartFooter";
import { buildStockChartOption } from "./foreignProStockChartOption";
import type {
  AkumulasiResult,
  ForeignPowerInfo,
  HistTopMode,
  MasterTickerItem,
  StockSeriesResponse,
} from "./foreignProTypes";

interface ForeignProStockViewProps {
  ticker: string;
  tickerSearch: string;
  setTickerSearch: (val: string) => void;
  suggestOpen: boolean;
  setSuggestOpen: (open: boolean) => void;
  suggestions: MasterTickerItem[];
  onSelectTicker: (tk: string) => void;
  onRefresh: () => void;
  loading: boolean;
  isLightMode: boolean;
  stockData: StockSeriesResponse | null;
  latestFp: ForeignPowerInfo;
  akumDays: number;
  setAkumDays: (v: number) => void;
  akumulasiResult: AkumulasiResult;
  valFFText: string;
  valFCoText: string;
  valKFText: string;
  valXLText: string;
  valSRText: string;
  latestPartisipasi: number;
  latestKat: number;
  latestKfh: number;
  latestKret: number;
  latestCorr: number;
  latestCorrKF: number;
  latestValM: number;
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
  showMA: boolean;
  setShowMA: (v: boolean) => void;
  maPeriod: number;
  setMaPeriod: (v: number) => void;
  maTarget: string;
  setMaTarget: (v: string) => void;
  additionalOverlay: string;
  setAdditionalOverlay: (v: string) => void;
  statType: string;
  setStatType: (v: string) => void;
  signalStat: string;
  setSignalStat: (v: string) => void;
  showPBNyangkut: boolean;
  setShowPBNyangkut: (v: boolean) => void;
  notifOpen: boolean;
  setNotifOpen: (v: boolean) => void;
  histTopMode: HistTopMode;
  setHistTopMode: (v: HistTopMode) => void;
  histBottomMode: string;
  setHistBottomMode: (v: string) => void;
  zoomDays: number;
  setZoomDays: (v: number) => void;
  windowOffset: number;
  totalStockPoints: number;
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
}

export default function ForeignProStockView(props: ForeignProStockViewProps) {
  const {
    ticker,
    loading,
    isLightMode,
    stockData,
    latestFp,
    akumDays,
    akumulasiResult,
    valFFText,
    latestPartisipasi,
    latestKat,
    latestKfh,
    latestKret,
    latestCorr,
    latestCorrKF,
    latestValM,
    showFF,
    showFCo,
    showKF,
    showXL,
    showSemiRitel,
    showMA,
    maPeriod,
    maTarget,
    additionalOverlay,
    showPBNyangkut,
    setShowPBNyangkut,
    notifOpen,
    setNotifOpen,
    histTopMode,
    setHistTopMode,
    histBottomMode,
    setHistBottomMode,
    zoomDays,
    setZoomDays,
    windowOffset,
    totalStockPoints,
    effectiveWindowStart,
    effectiveWindowEnd,
    datesSlice,
    opensSlice,
    closesSlice,
    lowsSlice,
    highsSlice,
    cumFFSlice,
    cumFCoSlice,
    cumKFSlice,
    cumXLSlice,
    cumSRSlice,
    cumLNRSlice,
  } = props;

  const chartOption = useMemo(
    () =>
      buildStockChartOption({
        stockData,
        isLightMode,
        effectiveWindowStart,
        effectiveWindowEnd,
        datesSlice,
        opensSlice,
        closesSlice,
        lowsSlice,
        highsSlice,
        cumFFSlice,
        cumFCoSlice,
        cumKFSlice,
        cumXLSlice,
        cumSRSlice,
        cumLNRSlice,
        showFF,
        showFCo,
        showKF,
        showXL,
        showSemiRitel,
        showMA,
        maPeriod,
        maTarget,
        additionalOverlay,
        showPBNyangkut,
        histTopMode,
        histBottomMode,
      }),
    [
      stockData,
      isLightMode,
      effectiveWindowStart,
      effectiveWindowEnd,
      datesSlice,
      opensSlice,
      closesSlice,
      lowsSlice,
      highsSlice,
      cumFFSlice,
      cumFCoSlice,
      cumKFSlice,
      cumXLSlice,
      cumSRSlice,
      cumLNRSlice,
      showFF,
      showFCo,
      showKF,
      showXL,
      showSemiRitel,
      showMA,
      maPeriod,
      maTarget,
      additionalOverlay,
      showPBNyangkut,
      histTopMode,
      histBottomMode,
    ]
  );

  const highlights = stockData?.analytics?.highlights || [];

  return (
    <div className="space-y-4">
      <ForeignProStockHero
        ticker={ticker}
        stockData={stockData}
        latestFp={latestFp}
        akumDays={akumDays}
        akumulasiResult={akumulasiResult}
        valFFText={valFFText}
        latestCorr={latestCorr}
      />

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
        <ForeignProStockSidebar {...props} />

        <div className="lg:col-span-9 space-y-3">
          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-xl p-3 shadow-md relative">
            <div className="flex items-center justify-between pb-2 border-b border-[var(--md-sys-color-outline-variant)]/40 flex-wrap gap-2">
              <div className="flex items-center gap-2">
                <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] tracking-wide">
                  {ticker} — CANDLESTICK &amp; CUMULATIVE FOREIGN FLOW
                </span>
                <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] bg-[var(--md-sys-color-surface-container-high)] px-2 py-0.5 rounded font-mono">
                  {zoomDays} Hari
                </span>
              </div>

              <div
                id="fpBadge"
                className="flex items-center gap-1.5 px-3 py-1 rounded-lg border font-mono font-bold text-xs bg-[var(--md-sys-color-surface-container-high)]"
                style={{ borderColor: latestFp.color, color: latestFp.color }}
              >
                <span id="fpBadgeLetter">{latestFp.badge}</span>
                <span className="text-[var(--md-sys-color-outline-variant)]">|</span>
                <span id="fpBadgeValue">{latestFp.score} pts</span>
              </div>
            </div>

            <div id="priceChart" className="h-[600px] w-full">
              {loading ? (
                <div
                  id="loadingVeil"
                  className="h-full flex items-center justify-center text-xs text-[var(--md-sys-color-on-surface-variant)] font-mono"
                >
                  Memuat grafik Foreign Flow Pro...
                </div>
              ) : (
                <ReactECharts option={chartOption} style={{ height: "100%", width: "100%" }} notMerge={true} />
              )}
            </div>

            <ForeignProStockChartFooter
              latestPartisipasi={latestPartisipasi}
              latestKat={latestKat}
              latestKfh={latestKfh}
              latestKret={latestKret}
              latestCorr={latestCorr}
              latestCorrKF={latestCorrKF}
              latestValM={latestValM}
              additionalOverlay={additionalOverlay}
              histTopMode={histTopMode}
              setHistTopMode={setHistTopMode}
              histBottomMode={histBottomMode}
              setHistBottomMode={setHistBottomMode}
              showPBNyangkut={showPBNyangkut}
              setShowPBNyangkut={setShowPBNyangkut}
              zoomDays={zoomDays}
              setZoomDays={setZoomDays}
              windowOffset={windowOffset}
              totalStockPoints={totalStockPoints}
              notifOpen={notifOpen}
              setNotifOpen={setNotifOpen}
              highlights={highlights}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
