import { calculateSMA } from "./foreignProMath";
import type { StockChartInput } from "./foreignProTypes";

export function buildStockOverlaySeries(
  input: StockChartInput,
  localNonRitelColor: string
): {
  maLineSeries: Record<string, unknown> | null;
  overlaySeries: Record<string, unknown> | null;
} {
  const {
    stockData,
    effectiveWindowStart,
    effectiveWindowEnd,
    cumFFSlice,
    cumFCoSlice,
    cumKFSlice,
    cumXLSlice,
    cumSRSlice,
    cumLNRSlice,
    showMA,
    maPeriod,
    maTarget,
    additionalOverlay,
  } = input;

  const an = stockData?.analytics;

  let maLineSeries: Record<string, unknown> | null = null;
  if (showMA) {
    const targetArr =
      maTarget === "fco"
        ? cumFCoSlice
        : maTarget === "kf"
          ? cumKFSlice
          : maTarget === "xl"
            ? cumXLSlice
            : maTarget === "semiritel"
              ? cumSRSlice
              : cumFFSlice;
    maLineSeries = {
      name: `MA-${maPeriod} ${maTarget.toUpperCase()}`,
      type: "line",
      yAxisIndex: 1,
      data: calculateSMA(targetArr, maPeriod),
      gridIndex: 0,
      smooth: true,
      lineStyle: { color: "#FFB380", width: 1.5, type: "solid" },
      showSymbol: false,
    };
  }

  let overlaySeries: Record<string, unknown> | null = null;
  if (additionalOverlay === "localNonRitel") {
    overlaySeries = {
      name: "Local Non Ritel",
      type: "line",
      yAxisIndex: 1,
      data: cumLNRSlice,
      gridIndex: 0,
      smooth: true,
      lineStyle: { color: localNonRitelColor, width: 1.2, type: "dashed" },
      showSymbol: false,
    };
  } else if (additionalOverlay === "fcoV2") {
    overlaySeries = {
      name: "F&Co v2",
      type: "line",
      yAxisIndex: 1,
      data: cumFCoSlice,
      gridIndex: 0,
      smooth: true,
      lineStyle: { color: "#D62839", width: 1.8 },
      showSymbol: false,
    };
  } else if (additionalOverlay === "korelasiKF") {
    overlaySeries = {
      name: "Korelasi Korban Flow",
      type: "line",
      xAxisIndex: 2,
      yAxisIndex: 4,
      data: an?.correlation_kf_20?.slice(effectiveWindowStart, effectiveWindowEnd) || [],
      lineStyle: { color: "#2F7676", width: 1.5 },
      showSymbol: false,
    };
  } else if (additionalOverlay === "kepRitel") {
    overlaySeries = {
      name: "% Kep Ritel",
      type: "line",
      xAxisIndex: 2,
      yAxisIndex: 4,
      data: an?.kepemilikan_ritel?.slice(effectiveWindowStart, effectiveWindowEnd) || [],
      lineStyle: { color: "#F87171", width: 1.5 },
      showSymbol: false,
    };
  } else if (additionalOverlay === "kepForeign") {
    overlaySeries = {
      name: "% Kep Foreign",
      type: "line",
      xAxisIndex: 2,
      yAxisIndex: 4,
      data: an?.kepemilikan_asing_total?.slice(effectiveWindowStart, effectiveWindowEnd) || [],
      lineStyle: { color: "#FF6820", width: 1.5 },
      showSymbol: false,
    };
  } else if (additionalOverlay === "kepHiddenF") {
    overlaySeries = {
      name: "% Kep Hidden F",
      type: "line",
      xAxisIndex: 2,
      yAxisIndex: 4,
      data: an?.kepemilikan_foreign_hidden?.slice(effectiveWindowStart, effectiveWindowEnd) || [],
      lineStyle: { color: "#2FD4C6", width: 1.5 },
      showSymbol: false,
    };
  }

  return { maLineSeries, overlaySeries };
}
