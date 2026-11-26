import { buildStockOverlaySeries } from "./foreignProStockChartSeries";
import type { StockChartInput } from "./foreignProTypes";

export function buildStockChartOption(input: StockChartInput): Record<string, unknown> {
  const {
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
    showFF,
    showFCo,
    showKF,
    showXL,
    showSemiRitel,
    additionalOverlay,
    showPBNyangkut,
    histTopMode,
    histBottomMode,
  } = input;

  if (!stockData?.analytics) return {};
  const an = stockData.analytics;

  const colors = {
    tooltipBg: isLightMode ? "rgba(255, 255, 255, 0.96)" : "rgba(22, 19, 32, 0.95)",
    tooltipBorder: isLightMode ? "rgba(132, 116, 106, 0.35)" : "rgba(255,255,255,0.15)",
    tooltipText: isLightMode ? "#1E1A17" : "#F2F2F2",
    pointerBg: isLightMode ? "rgba(237, 223, 211, 0.96)" : "rgba(30, 24, 48, 0.95)",
    pointerBorder: isLightMode ? "rgba(132, 116, 106, 0.4)" : "rgba(255,255,255,0.2)",
    legendText: isLightMode ? "#4E453F" : "#A39DB5",
    axisLine: isLightMode ? "rgba(30, 26, 23, 0.18)" : "rgba(255,255,255,0.12)",
    axisTick: isLightMode ? "rgba(30, 26, 23, 0.2)" : "rgba(255,255,255,0.2)",
    xAxisLabel: isLightMode ? "#4E453F" : "#8E88A0",
    yAxisLabel: isLightMode ? "#4E453F" : "#A39DB5",
    yAxisSubLabel: isLightMode ? "#4E453F" : "#7A7290",
    splitLinePrimary: isLightMode ? "rgba(30, 26, 23, 0.08)" : "rgba(255,255,255,0.05)",
    splitLineSub: isLightMode ? "rgba(30, 26, 23, 0.08)" : "rgba(255,255,255,0.06)",
    candleDefault: isLightMode ? "#12233F" : "#FFFFFF",
    localNonRitel: isLightMode ? "#1E1A17" : "#FFFFFF",
    histValueBar: isLightMode ? "rgba(30, 26, 23, 0.25)" : "rgba(255,255,255,0.4)",
    rightFlowAxis: isLightMode ? "#8C5000" : "#FF6820",
    rightPctAxis: isLightMode ? "#0D9488" : "#2FD4C6",
  };

  const candleData = datesSlice.map((dt, idx) => {
    const sig = an.special_candles?.find((sc) => sc.date === dt);
    let cColor = colors.candleDefault;
    let bColor = colors.candleDefault;
    if (sig) {
      if (sig.type === "PURPLE") {
        cColor = "#C724F5";
        bColor = "#C724F5";
      } else if (sig.type === "RED" && showPBNyangkut) {
        cColor = "#D62839";
        bColor = "#D62839";
      } else if (sig.type === "YELLOW") {
        cColor = "#FFFF00";
        bColor = "#FFFF00";
      } else if (sig.type === "ORANGE") {
        cColor = "#FF6820";
        bColor = "#FF6820";
      }
    }
    return {
      value: [opensSlice[idx], closesSlice[idx], lowsSlice[idx], highsSlice[idx]],
      itemStyle: { color: cColor, color0: cColor, borderColor: bColor, borderColor0: bColor },
    };
  });

  const histTopVals = (
    histTopMode === "ff" ? an.ff_daily_val || [] : an.fco_daily_val || []
  ).slice(effectiveWindowStart, effectiveWindowEnd);

  const rawHistBottom =
    histBottomMode === "korban"
      ? an.kf_daily_val
      : histBottomMode === "xl"
        ? an.xl_daily_val
        : histBottomMode === "semiritel"
          ? an.semiritel_daily_val
          : an.values_miliar || [];
  const histBottomVals = (rawHistBottom || []).slice(effectiveWindowStart, effectiveWindowEnd);

  const { maLineSeries, overlaySeries } = buildStockOverlaySeries(input, colors.localNonRitel);

  const flowLines = [
    showFF && { name: "Kumulatif Asing", data: cumFFSlice, color: "#FF6820", width: 2 },
    showFCo && { name: "Foreign & Co", data: cumFCoSlice, color: "#D62839", width: 2 },
    showKF && { name: "Korban Flow", data: cumKFSlice, color: "#2F7676", width: 1.8 },
    showXL && { name: "Broker XL", data: cumXLSlice, color: "#3A7BD5", width: 1.5 },
    showSemiRitel && { name: "Semi Ritel", data: cumSRSlice, color: "#1DB87A", width: 1.5 },
  ]
    .filter((item): item is { name: string; data: number[]; color: string; width: number } => Boolean(item))
    .map((line) => ({
      name: line.name,
      type: "line",
      yAxisIndex: 1,
      data: line.data,
      gridIndex: 0,
      smooth: true,
      lineStyle: { color: line.color, width: line.width },
      showSymbol: false,
    }));

  return {
    backgroundColor: "transparent",
    animation: false,
    tooltip: {
      trigger: "axis",
      axisPointer: { type: "cross", link: [{ xAxisIndex: "all" }] },
      backgroundColor: colors.tooltipBg,
      borderColor: colors.tooltipBorder,
      textStyle: { color: colors.tooltipText, fontSize: 10 },
      confine: true,
    },
    axisPointer: {
      link: [{ xAxisIndex: "all" }],
      label: {
        backgroundColor: colors.pointerBg,
        borderColor: colors.pointerBorder,
        borderWidth: 1,
        color: colors.tooltipText,
        fontSize: 10,
      },
    },
    legend: {
      type: "scroll",
      data: [
        "Candlestick",
        ...flowLines.map((l) => l.name),
        ...(maLineSeries ? [String(maLineSeries.name)] : []),
        ...(overlaySeries ? [String(overlaySeries.name)] : []),
      ],
      textStyle: { color: colors.legendText, fontSize: 10 },
      top: 0,
    },
    grid: [
      { left: 12, right: 12, top: "7%", height: "52%", containLabel: true },
      { left: 12, right: 12, top: "63%", height: "16%", containLabel: true },
      { left: 12, right: 12, top: "82%", height: "14%", containLabel: true },
    ],
    xAxis: [
      { type: "category", data: datesSlice, gridIndex: 0, scale: true, boundaryGap: true, axisLine: { lineStyle: { color: colors.axisLine } }, axisLabel: { show: false }, axisTick: { show: false } },
      { type: "category", data: datesSlice, gridIndex: 1, scale: true, boundaryGap: true, axisLine: { lineStyle: { color: colors.axisLine } }, axisLabel: { show: false }, axisTick: { show: false } },
      {
        type: "category",
        data: datesSlice,
        gridIndex: 2,
        scale: true,
        boundaryGap: true,
        axisLine: { lineStyle: { color: colors.axisLine } },
        axisLabel: {
          show: true,
          color: colors.xAxisLabel,
          fontSize: 10,
          interval: Math.max(Math.floor(datesSlice.length / 10), 1),
          formatter: (dt: string) => {
            if (!dt) return "";
            const parts = dt.split("-");
            return parts.length === 3 ? `${parts[2]}/${parts[1]}` : dt;
          },
        },
        axisTick: { show: true, lineStyle: { color: colors.axisTick } },
      },
    ],
    yAxis: [
      { scale: true, gridIndex: 0, splitLine: { lineStyle: { color: colors.splitLinePrimary } }, axisLabel: { color: colors.yAxisLabel, fontSize: 10, formatter: (v: number) => v.toLocaleString("id-ID") } },
      { scale: true, gridIndex: 0, position: "right", splitLine: { show: false }, axisLabel: { color: colors.rightFlowAxis, fontSize: 10, formatter: (v: number) => `${v >= 0 ? "+" : ""}${Math.round(v)} M` } },
      { scale: true, gridIndex: 1, splitLine: { lineStyle: { color: colors.splitLineSub } }, axisLabel: { color: colors.yAxisSubLabel, fontSize: 9, formatter: (v: number) => `${v >= 0 ? "+" : ""}${Math.round(v)}` } },
      { scale: histBottomMode !== "value", min: histBottomMode === "value" ? 0 : undefined, gridIndex: 2, splitLine: { lineStyle: { color: colors.splitLineSub } }, axisLabel: { color: colors.yAxisSubLabel, fontSize: 9, formatter: (v: number) => Math.round(v).toString() } },
      { scale: true, gridIndex: 2, position: "right", splitLine: { show: false }, axisLabel: { color: colors.rightPctAxis, fontSize: 9, formatter: (v: number) => additionalOverlay === "korelasiKF" ? v.toFixed(2) : `${Math.round(v)}%` } },
    ],
    series: [
      { name: "Candlestick", type: "candlestick", data: candleData, gridIndex: 0 },
      ...flowLines,
      ...(maLineSeries ? [maLineSeries] : []),
      ...(overlaySeries ? [overlaySeries] : []),
      {
        name: "Histogram Top",
        type: "bar",
        xAxisIndex: 1,
        yAxisIndex: 2,
        barMaxWidth: 8,
        data: histTopVals.map((v) => ({ value: v, itemStyle: { color: v >= 0 ? "#16A34A" : "#DC2626" } })),
      },
      {
        name: "Histogram Bottom",
        type: "bar",
        xAxisIndex: 2,
        yAxisIndex: 3,
        barMaxWidth: 8,
        data: histBottomVals.map((v) => ({
          value: v,
          itemStyle: { color: histBottomMode === "value" ? colors.histValueBar : v >= 0 ? "#16A34A" : "#DC2626" },
        })),
      },
    ],
  };
}
