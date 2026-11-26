import { useMemo, useState } from "react";
import ReactECharts from "echarts-for-react";
import { useAppStore } from "@/store/useAppStore";
import type { IpoAnalysisMode, IpoCandle } from "./ipoTypes";

export default function IpoCandlestickChart({
  ticker,
  candles,
  ipoPrice,
  mode,
  onAfterIpoClick,
}: {
  ticker: string;
  candles: IpoCandle[];
  ipoPrice?: number | null;
  mode: IpoAnalysisMode;
  onAfterIpoClick?: () => void;
}) {
  const [visibleCount, setVisibleCount] = useState<number>(0); // 0 = show all
  const theme = useAppStore((s) => s.theme);
  const isDark = theme === "dark";

  const displayedCandles = useMemo(() => {
    if (!candles || candles.length === 0) return [];
    if (!visibleCount || visibleCount >= candles.length) return candles;
    return candles.slice(candles.length - visibleCount);
  }, [candles, visibleCount]);

  const option = useMemo(() => {
    if (!displayedCandles.length) return {};

    const palette = isDark
      ? {
          tooltipBg: "rgba(33, 29, 26, 0.96)",
          tooltipBorder: "#52443C",
          tooltipText: "#EBE0D9",
          axisLine: "rgba(215, 194, 180, 0.25)",
          axisLabel: "#D7C2B4",
          splitLine: "rgba(255, 255, 255, 0.08)",
          ipoLine: "#FFB879",
        }
      : {
          tooltipBg: "rgba(255, 255, 255, 0.96)",
          tooltipBorder: "#D7C2B4",
          tooltipText: "#1E1A17",
          axisLine: "rgba(78, 69, 63, 0.28)",
          axisLabel: "#4E453F",
          splitLine: "rgba(30, 26, 23, 0.08)",
          ipoLine: "#8C5000",
        };

    const dates = displayedCandles.map((c) => c.date);
    const ohlc = displayedCandles.map((c) => [c.open ?? c.close, c.close, c.low ?? c.close, c.high ?? c.close]);
    const volumes = displayedCandles.map((c) => ({
      value: c.volume ?? 0,
      itemStyle: {
        color:
          (c.close ?? 0) >= (c.open ?? c.close ?? 0)
            ? "rgba(22, 163, 74, 0.55)"
            : "rgba(220, 38, 38, 0.55)",
      },
    }));

    const prices = displayedCandles.flatMap((c) =>
      [c.low, c.high, c.close, c.open].filter((v): v is number => v != null && v > 0)
    );
    const pMin = prices.length ? Math.min(...prices) : 0;
    const pMax = prices.length ? Math.max(...prices) : 0;
    const showIpoLine = ipoPrice && ipoPrice > 0 && ipoPrice >= pMin * 0.85 && ipoPrice <= pMax * 1.15;

    return {
      animation: false,
      tooltip: {
        trigger: "axis",
        axisPointer: { type: "cross" },
        backgroundColor: palette.tooltipBg,
        borderColor: palette.tooltipBorder,
        textStyle: { color: palette.tooltipText, fontSize: 11 },
        formatter: (params: Array<{ dataIndex: number }>) => {
          if (!params || !params.length) return "";
          const idx = params[0].dataIndex;
          const c = displayedCandles[idx];
          if (!c) return "";
          const chg = c.prev && c.close ? ((c.close - c.prev) / c.prev) * 100 : null;
          const chgStr = chg != null ? `${chg >= 0 ? "+" : ""}${chg.toFixed(2)}%` : "—";
          return `
            <div style="font-weight:700;margin-bottom:4px;">${ticker} · ${c.date}</div>
            <div>O: <b>${c.open?.toLocaleString() ?? "—"}</b> &nbsp; H: <b>${c.high?.toLocaleString() ?? "—"}</b></div>
            <div>L: <b>${c.low?.toLocaleString() ?? "—"}</b> &nbsp; C: <b>${c.close?.toLocaleString() ?? "—"}</b> (${chgStr})</div>
            <div>Avg: <b>${c.tavg ? Math.round(c.tavg).toLocaleString() : "—"}</b></div>
            <div>Vol: <b>${c.volume?.toLocaleString() ?? "0"}</b> lot</div>
          `;
        },
      },
      axisPointer: { link: [{ xAxisIndex: "all" }] },
      grid: [
        { left: 56, right: 20, top: 20, height: "62%" },
        { left: 56, right: 20, top: "76%", height: "16%" },
      ],
      xAxis: [
        {
          type: "category",
          data: dates,
          boundaryGap: true,
          axisLine: { lineStyle: { color: palette.axisLine } },
          axisLabel: { color: palette.axisLabel, fontSize: 10 },
          splitLine: { show: false },
        },
        {
          type: "category",
          gridIndex: 1,
          data: dates,
          boundaryGap: true,
          axisLabel: { show: false },
          axisTick: { show: false },
          axisLine: { lineStyle: { color: palette.axisLine } },
        },
      ],
      yAxis: [
        {
          scale: true,
          splitArea: { show: false },
          splitLine: { lineStyle: { color: palette.splitLine } },
          axisLabel: { color: palette.axisLabel, fontSize: 10 },
        },
        {
          scale: true,
          gridIndex: 1,
          splitNumber: 2,
          axisLabel: { show: false },
          axisLine: { show: false },
          axisTick: { show: false },
          splitLine: { show: false },
        },
      ],
      series: [
        {
          name: ticker,
          type: "candlestick",
          data: ohlc,
          itemStyle: {
            color: "#16A34A",
            color0: "#DC2626",
            borderColor: "#16A34A",
            borderColor0: "#DC2626",
          },
          markLine: showIpoLine
            ? {
                symbol: "none",
                silent: true,
                lineStyle: { color: palette.ipoLine, width: 1.4, type: "dashed", opacity: 0.75 },
                label: {
                  formatter: `IPO ${ipoPrice?.toLocaleString()}`,
                  position: "insideEndTop",
                  color: palette.ipoLine,
                  fontSize: 10,
                  fontWeight: 700,
                },
                data: [{ yAxis: ipoPrice }],
              }
            : undefined,
        },
        {
          name: "Volume",
          type: "bar",
          xAxisIndex: 1,
          yAxisIndex: 1,
          data: volumes,
        },
      ],
    };
  }, [displayedCandles, ticker, ipoPrice, isDark]);

  const handleZoomIn = () => {
    const cur = visibleCount || candles.length;
    setVisibleCount(Math.max(10, Math.round(cur * 0.7)));
  };

  const handleZoomOut = () => {
    const cur = visibleCount || candles.length;
    const next = Math.round(cur * 1.4);
    setVisibleCount(next >= candles.length ? 0 : next);
  };

  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-4 shadow-sm">
      <div className="flex flex-wrap items-center justify-between gap-2 mb-2">
        <div className="flex items-center gap-2">
          <h3 className="text-xs font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface)]">
            {ticker} Candlestick
          </h3>
          <span
            className={`px-2 py-0.5 rounded-md text-[10px] font-extrabold uppercase border ${
              mode === "after"
                ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border-[var(--md-sys-color-primary)]/40"
                : "bg-[var(--md-sys-color-primary)]/10 text-[var(--md-sys-color-primary)] border-[var(--md-sys-color-primary)]/40"
            }`}
          >
            {mode === "after" ? "After IPO" : "Periodic"}
          </span>
        </div>

        <div className="flex items-center gap-1.5">
          {mode !== "after" && onAfterIpoClick && (
            <button
              type="button"
              onClick={onAfterIpoClick}
              className="px-2 py-1 rounded-md text-[10px] font-bold border border-[var(--md-sys-color-primary)]/50 text-[var(--md-sys-color-primary)] hover:bg-[var(--md-sys-color-primary)]/10 transition cursor-pointer"
              title="Aktifkan mode After IPO (sejak D0)"
            >
              After IPO
            </button>
          )}
          <button
            type="button"
            onClick={handleZoomOut}
            className="w-7 h-7 rounded-md text-xs font-bold bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] hover:border-[var(--md-sys-color-primary)] transition cursor-pointer"
            title="Zoom out"
          >
            −
          </button>
          <button
            type="button"
            onClick={handleZoomIn}
            className="w-7 h-7 rounded-md text-xs font-bold bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] hover:border-[var(--md-sys-color-primary)] transition cursor-pointer"
            title="Zoom in"
          >
            +
          </button>
        </div>
      </div>

      {displayedCandles.length > 0 ? (
        <ReactECharts option={option} style={{ height: 340, width: "100%" }} notMerge={true} />
      ) : (
        <div className="h-[260px] flex items-center justify-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
          Belum ada data candlestick harian untuk {ticker} pada rentang ini.
        </div>
      )}
    </div>
  );
}
