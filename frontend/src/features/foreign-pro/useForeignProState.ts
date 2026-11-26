import { useState, useEffect, useMemo, useCallback } from "react";
import { useSearchParams } from "react-router-dom";
import {
  getForeignFlowSeries,
  getForeignFlowScreener,
  getForeignFlowIhsg,
  getMmdTickers,
} from "@/lib/api/terminal";
import { useAppStore } from "@/store/useAppStore";
import {
  calcCumLnrSlice,
  calcIhsgMacroMetrics,
  calcWindowCumulative,
  calculateAkumulasi,
} from "./foreignProMath";
import { filterScreenerRows } from "./foreignProScreenerMath";
import {
  buildAdminReportText,
  computeMarketReportData,
  computePeriodReportData,
  getAccumHeaderLabel,
} from "./foreignProReportMath";
import type {
  AdminReportMode,
  FFMainTab,
  HistTopMode,
  IhsgResponse,
  MasterTickerItem,
  ScreenerResponse,
  ScreenerRow,
  ScreenerSubTab,
  StockSeriesResponse,
} from "./foreignProTypes";

export function useForeignProState() {
  const [searchParams] = useSearchParams();
  const queryTicker = searchParams.get("ticker");
  const activeTicker = useAppStore((s) => s.activeTicker);
  const setActiveTicker = useAppStore((s) => s.setActiveTicker);
  const theme = useAppStore((s) => s.theme);
  const isLightMode = theme === "light";

  const [activeMainTab, setActiveMainTab] = useState<FFMainTab>("stock");
  const [screenerSubTab, setScreenerSubTab] = useState<ScreenerSubTab>("basic");
  const [parameterDrawerOpen, setParameterDrawerOpen] = useState(false);

  const [ticker, setTicker] = useState<string>(queryTicker?.toUpperCase() || activeTicker || "BBCA");
  const [tickerSearch, setTickerSearch] = useState<string>(ticker);
  const [suggestOpen, setSuggestOpen] = useState(false);
  const [allTickers, setAllTickers] = useState<MasterTickerItem[]>([]);

  useEffect(() => {
    if (queryTicker && queryTicker.toUpperCase() !== ticker) {
      const up = queryTicker.toUpperCase();
      setTicker(up);
      setTickerSearch(up);
      setActiveTicker(up);
    }
  }, [queryTicker, ticker, setActiveTicker]);

  const [showFF, setShowFF] = useState(true);
  const [showFCo, setShowFCo] = useState(false);
  const [showKF, setShowKF] = useState(true);
  const [showXL, setShowXL] = useState(false);
  const [showSemiRitel, setShowSemiRitel] = useState(false);
  const [showMA, setShowMA] = useState(false);
  const [maPeriod, setMaPeriod] = useState(20);
  const [maTarget, setMaTarget] = useState("ff");
  const [additionalOverlay, setAdditionalOverlay] = useState("none");
  const [statType, setStatType] = useState("indicatorCorrelation");
  const [signalStat, setSignalStat] = useState("inflow_streak");
  const [akumDays, setAkumDays] = useState(20);
  const [showPBNyangkut, setShowPBNyangkut] = useState(false);
  const [notifOpen, setNotifOpen] = useState(true);
  const [histTopMode, setHistTopMode] = useState<HistTopMode>("ff");
  const [histBottomMode, setHistBottomMode] = useState("value");
  const [zoomDays, setZoomDays] = useState(100);
  const [windowOffset] = useState(0);

  const [stockData, setStockData] = useState<StockSeriesResponse | null>(null);
  const [ihsgData, setIhsgData] = useState<IhsgResponse | null>(null);
  const [screenerData, setScreenerData] = useState<ScreenerRow[]>([]);
  const [loading, setLoading] = useState(true);

  const [universe, setUniverse] = useState("ff");
  const [windowLookback, setWindowLookback] = useState(40);
  const [fpFilter] = useState("ALL");
  const [transBucket, setTransBucket] = useState("all");
  const [priceClass, setPriceClass] = useState("All");
  const [screenerSearch, setScreenerSearch] = useState("");
  const [stateOnly, setStateOnly] = useState(false);
  const [hfpOnly, setHfpOnly] = useState(false);
  const [mfpOnly, setMfpOnly] = useState(false);
  const [customCategory, setCustomCategory] = useState("top");
  const [activeStrategyId, setActiveStrategyId] = useState("basicHFP");
  const [strategyDesc, setStrategyDesc] = useState(
    "Saham dengan skor Foreign Power TINGGI (badge HFP, skor > 30)."
  );
  const [activeReportType, setActiveReportType] = useState<string | null>(null);
  const [periodReportWindow, setPeriodReportWindow] = useState<number | null>(null);
  const [adminMode, setAdminMode] = useState<AdminReportMode>("end");
  const [adminCopied, setAdminCopied] = useState(false);

  const [catIhsg, setCatIhsg] = useState(true);
  const [catKonglo, setCatKonglo] = useState(false);
  const [catFund, setCatFund] = useState(false);
  const [metFF, setMetFF] = useState(true);
  const [metFCo, setMetFCo] = useState(false);
  const [metKF, setMetKF] = useState(true);
  const [metXL, setMetXL] = useState(false);
  const [metSR, setMetSR] = useState(false);
  const [flipKF, setFlipKF] = useState(false);
  const [flipXL, setFlipXL] = useState(false);
  const [flipSR, setFlipSR] = useState(false);

  useEffect(() => {
    getMmdTickers<{ tickers?: MasterTickerItem[] }>()
      .then((res) => {
        if (res.tickers) setAllTickers(res.tickers);
      })
      .catch((err) => console.error("Error loading master tickers:", err));
  }, []);

  const loadStockSeries = useCallback((tk: string) => {
    setLoading(true);
    getForeignFlowSeries<StockSeriesResponse>(tk, 250)
      .then((res) => {
        setStockData(res);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error loading FF series:", err);
        setLoading(false);
      });
  }, []);

  useEffect(() => {
    if (activeMainTab === "stock") loadStockSeries(ticker);
  }, [ticker, activeMainTab, loadStockSeries]);

  const loadIhsg = useCallback(() => {
    getForeignFlowIhsg<IhsgResponse>()
      .then((res) => setIhsgData(res))
      .catch((err) => console.error("Error loading IHSG data:", err));
  }, []);

  useEffect(() => {
    if (!ihsgData) loadIhsg();
  }, [ihsgData, loadIhsg]);

  const loadScreener = useCallback(() => {
    setLoading(true);
    getForeignFlowScreener<ScreenerResponse>(universe, windowLookback, fpFilter, transBucket, undefined, 1000)
      .then((res) => {
        setScreenerData(res.screener || []);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error loading FF screener:", err);
        setLoading(false);
      });
  }, [universe, windowLookback, fpFilter, transBucket]);

  useEffect(() => {
    if (activeMainTab === "screener") loadScreener();
  }, [activeMainTab, loadScreener]);

  const handleSelectTicker = useCallback(
    (tk: string) => {
      const up = tk.toUpperCase();
      setTicker(up);
      setTickerSearch(up);
      setActiveTicker(up);
      setSuggestOpen(false);
    },
    [setActiveTicker]
  );

  const handleRefresh = useCallback(() => {
    if (activeMainTab === "ihsg") loadIhsg();
    else if (activeMainTab === "screener") loadScreener();
    else loadStockSeries(ticker);
  }, [activeMainTab, loadIhsg, loadScreener, loadStockSeries, ticker]);

  const suggestions = useMemo(() => {
    if (!tickerSearch) return [];
    const term = tickerSearch.toUpperCase();
    return allTickers.filter((t) => t.ticker.toUpperCase().includes(term)).slice(0, 8);
  }, [allTickers, tickerSearch]);

  const ihsgMacroMetrics = useMemo(
    () => calcIhsgMacroMetrics(ihsgData, flipKF, flipXL, flipSR),
    [ihsgData, flipKF, flipXL, flipSR]
  );

  const an = stockData?.analytics;
  const totalStockPoints = an?.dates?.length || 0;
  const effectiveWindowEnd = Math.max(zoomDays, totalStockPoints - windowOffset);
  const effectiveWindowStart = Math.max(0, effectiveWindowEnd - zoomDays);

  const datesSlice = an?.dates?.slice(effectiveWindowStart, effectiveWindowEnd) || [];
  const closesSlice = an?.closes?.slice(effectiveWindowStart, effectiveWindowEnd) || [];
  const opensSlice = an?.opens?.slice(effectiveWindowStart, effectiveWindowEnd) || [];
  const highsSlice = an?.highs?.slice(effectiveWindowStart, effectiveWindowEnd) || [];
  const lowsSlice = an?.lows?.slice(effectiveWindowStart, effectiveWindowEnd) || [];

  const cumFFSlice = useMemo(() => calcWindowCumulative(an?.ff_daily_val, effectiveWindowStart, effectiveWindowEnd), [an?.ff_daily_val, effectiveWindowStart, effectiveWindowEnd]);
  const cumFCoSlice = useMemo(() => calcWindowCumulative(an?.fco_daily_val, effectiveWindowStart, effectiveWindowEnd), [an?.fco_daily_val, effectiveWindowStart, effectiveWindowEnd]);
  const cumKFSlice = useMemo(() => calcWindowCumulative(an?.kf_daily_val, effectiveWindowStart, effectiveWindowEnd), [an?.kf_daily_val, effectiveWindowStart, effectiveWindowEnd]);
  const cumXLSlice = useMemo(() => calcWindowCumulative(an?.xl_daily_val, effectiveWindowStart, effectiveWindowEnd), [an?.xl_daily_val, effectiveWindowStart, effectiveWindowEnd]);
  const cumSRSlice = useMemo(() => calcWindowCumulative(an?.semiritel_daily_val, effectiveWindowStart, effectiveWindowEnd), [an?.semiritel_daily_val, effectiveWindowStart, effectiveWindowEnd]);
  const cumLNRSlice = useMemo(() => calcCumLnrSlice(an?.cum_lnr, effectiveWindowStart, effectiveWindowEnd), [an?.cum_lnr, effectiveWindowStart, effectiveWindowEnd]);

  const startDateText = datesSlice.length > 0 ? datesSlice[0] : "";
  const ffValNum = cumFFSlice.length ? cumFFSlice[cumFFSlice.length - 1] : 0;
  const fcoValNum = cumFCoSlice.length ? cumFCoSlice[cumFCoSlice.length - 1] : 0;
  const kfValNum = cumKFSlice.length ? cumKFSlice[cumKFSlice.length - 1] : 0;
  const xlValNum = cumXLSlice.length ? cumXLSlice[cumXLSlice.length - 1] : 0;
  const srValNum = cumSRSlice.length ? cumSRSlice[cumSRSlice.length - 1] : 0;

  const valFFText = `${ffValNum >= 0 ? "+" : ""}${ffValNum.toFixed(1)} M ${startDateText ? `(${startDateText})` : ""}`;
  const valFCoText = `${fcoValNum >= 0 ? "+" : ""}${fcoValNum.toFixed(1)} M`;
  const valKFText = `${kfValNum >= 0 ? "+" : ""}${kfValNum.toFixed(1)} M`;
  const valXLText = `${xlValNum >= 0 ? "+" : ""}${xlValNum.toFixed(1)} M`;
  const valSRText = `${srValNum >= 0 ? "+" : ""}${srValNum.toFixed(1)} M`;

  const latestPartisipasi = an?.partisipasi_series?.length ? an.partisipasi_series[an.partisipasi_series.length - 1] : 0;
  const latestKat = an?.kepemilikan_asing_total?.length ? an.kepemilikan_asing_total[an.kepemilikan_asing_total.length - 1] : 0;
  const latestKfh = an?.kepemilikan_foreign_hidden?.length ? an.kepemilikan_foreign_hidden[an.kepemilikan_foreign_hidden.length - 1] : 0;
  const latestKret = an?.kepemilikan_ritel?.length ? an.kepemilikan_ritel[an.kepemilikan_ritel.length - 1] : 0;
  const latestCorr = an?.correlation_20?.length && an.correlation_20[an.correlation_20.length - 1] !== null ? (an.correlation_20[an.correlation_20.length - 1] as number) : 0;
  const latestCorrKF = an?.correlation_kf_20?.length && an.correlation_kf_20[an.correlation_kf_20.length - 1] !== null ? (an.correlation_kf_20[an.correlation_kf_20.length - 1] as number) : 0;
  const latestValM = an?.values_miliar?.length ? an.values_miliar[an.values_miliar.length - 1] : 0;

  const akumulasiResult = useMemo(() => calculateAkumulasi(stockData, akumDays), [stockData, akumDays]);

  const filteredScreenerList = useMemo(
    () =>
      filterScreenerRows({
        screenerData,
        screenerSearch,
        screenerSubTab,
        activeStrategyId,
        universe,
        hfpOnly,
        mfpOnly,
        stateOnly,
        priceClass,
      }),
    [screenerData, screenerSearch, screenerSubTab, activeStrategyId, universe, hfpOnly, mfpOnly, stateOnly, priceClass]
  );

  const marketReportData = useMemo(() => computeMarketReportData(screenerData), [screenerData]);
  const periodReportData = useMemo(() => computePeriodReportData(screenerData, periodReportWindow), [screenerData, periodReportWindow]);
  const adminReportText = useMemo(
    () => buildAdminReportText({ adminMode, ihsgData, marketReportData, screenerData }),
    [adminMode, ihsgData, marketReportData, screenerData]
  );
  const accumHeaderLabel = useMemo(() => getAccumHeaderLabel(universe, windowLookback), [universe, windowLookback]);
  const latestFp = an?.latest_fp || { score: 50.0, badge: "MFP", color: "#FF6820" };

  return {
    isLightMode, activeMainTab, setActiveMainTab, screenerSubTab, setScreenerSubTab,
    parameterDrawerOpen, setParameterDrawerOpen, ticker, tickerSearch, setTickerSearch,
    suggestOpen, setSuggestOpen, suggestions, handleSelectTicker, handleRefresh, loadStockSeries,
    showFF, setShowFF, showFCo, setShowFCo, showKF, setShowKF, showXL, setShowXL,
    showSemiRitel, setShowSemiRitel, showMA, setShowMA, maPeriod, setMaPeriod, maTarget, setMaTarget,
    additionalOverlay, setAdditionalOverlay, statType, setStatType, signalStat, setSignalStat,
    akumDays, setAkumDays, showPBNyangkut, setShowPBNyangkut, notifOpen, setNotifOpen,
    histTopMode, setHistTopMode, histBottomMode, setHistBottomMode, zoomDays, setZoomDays, windowOffset,
    stockData, ihsgData, loading, ihsgMacroMetrics, totalStockPoints, effectiveWindowStart, effectiveWindowEnd,
    datesSlice, closesSlice, opensSlice, highsSlice, lowsSlice,
    cumFFSlice, cumFCoSlice, cumKFSlice, cumXLSlice, cumSRSlice, cumLNRSlice,
    valFFText, valFCoText, valKFText, valXLText, valSRText,
    latestPartisipasi, latestKat, latestKfh, latestKret, latestCorr, latestCorrKF, latestValM,
    akumulasiResult, latestFp,
    universe, setUniverse, windowLookback, setWindowLookback, transBucket, setTransBucket,
    priceClass, setPriceClass, screenerSearch, setScreenerSearch, stateOnly, setStateOnly,
    hfpOnly, setHfpOnly, mfpOnly, setMfpOnly, customCategory, setCustomCategory,
    activeStrategyId, setActiveStrategyId, strategyDesc, setStrategyDesc,
    activeReportType, setActiveReportType, periodReportWindow, setPeriodReportWindow,
    adminMode, setAdminMode, adminCopied, setAdminCopied,
    filteredScreenerList, marketReportData, periodReportData, adminReportText, accumHeaderLabel,
    catIhsg, setCatIhsg, catKonglo, setCatKonglo, catFund, setCatFund,
    metFF, setMetFF, metFCo, setMetFCo, metKF, setMetKF, metXL, setMetXL, metSR, setMetSR,
    flipKF, setFlipKF, flipXL, setFlipXL, flipSR, setFlipSR,
  };
}
