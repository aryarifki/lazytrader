import { useMemo } from "react";
import ReactECharts from "echarts-for-react";
import { useAppStore } from "@/store/useAppStore";
import type { MmdDetailResponse, MmdT5Item } from "./mmdTypes";

interface MmdChartCardProps {
  candleData: MmdDetailResponse | null;
  t5Data: MmdT5Item[];
  showT5Markers: boolean;
  setShowT5Markers?: (v: boolean) => void;
  showTooltip: boolean;
  setShowTooltip?: (v: boolean) => void;
  isBpMode?: boolean;
  ticker?: string;
}

export default function MmdChartCard({
  candleData,
  t5Data,
  showT5Markers,
  setShowT5Markers,
  showTooltip,
  setShowTooltip,
  isBpMode = false,
  ticker,
}: MmdChartCardProps) {
  const theme = useAppStore((s) => s.theme);
  const isDark = theme === "dark";

  const option = useMemo(() => {
    if (!candleData?.candles || candleData.candles.length === 0) return {};

    const colors = {
      tooltipBg: isDark ? "rgba(33, 29, 26, 0.96)" : "rgba(255, 255, 255, 0.96)",
      tooltipBorder: isDark ? "rgba(255, 255, 255, 0.16)" : "rgba(132, 116, 106, 0.35)",
      tooltipText: isDark ? "#EBE0D9" : "#1E1A17",
      legendText: isDark ? "#D7C2B4" : "#4E453F",
      axisLine: isDark ? "rgba(255, 255, 255, 0.14)" : "rgba(30, 26, 23, 0.18)",
      axisLabel: isDark ? "#9F8D83" : "#4E453F",
      splitLine: isDark ? "rgba(255, 255, 255, 0.06)" : "rgba(30, 26, 23, 0.08)",
      upColor: isDark ? "#10B981" : "#166534",
      downColor: isDark ? "#F43F5E" : "#991B1B",
      upVol: isDark ? "rgba(16, 185, 129, 0.45)" : "rgba(22, 101, 52, 0.35)",
      downVol: isDark ? "rgba(244, 63, 94, 0.45)" : "rgba(153, 27, 27, 0.35)",
      tavgColor: isDark ? "#FFC24D" : "#B45309",
      bandarColor: isDark ? "#FD6602" : "#C2410C",
    };

    const dates = candleData.candles.map((c) => c.date_iso);
    const candleValues = candleData.candles.map((c) => [c.open ?? 0, c.close ?? 0, c.low ?? 0, c.high ?? 0]);
    const bandarAcc = (candleData.bandar_accumulation_series || []).map((b) => b.cumulative_bandar);
    const tavgLine = candleData.candles.map((c) => c.tavg || c.close || 0);
    const volumes = candleData.candles.map((c) => ({
      value: c.volume ?? 0,
      itemStyle: {
        color: (c.close ?? 0) >= (c.open ?? 0) ? colors.upVol : colors.downVol,
      },
    }));

    const markPoints: Record<string, unknown>[] = [];
    if (showT5Markers && t5Data.length > 0) {
      t5Data.forEach((t) => {
        const dateIdx = dates.indexOf(t.date);
        if (dateIdx !== -1) {
          const act = (t.action || "").toUpperCase();
          const isBuy = act === "BUY";
          const isSell = act === "SELL";
          markPoints.push({
            name: "T5 Transaksi",
            coord: [t.date, isBuy ? candleValues[dateIdx][2] : candleValues[dateIdx][3]],
            value: isBuy ? "▲" : isSell ? "▼" : "◆",
            symbol: isBuy || isSell ? "triangle" : "diamond",
            symbolRotate: isSell ? 180 : 0,
            symbolSize: 12,
            itemStyle: {
              color: isBuy ? "#16A34A" : isSell ? "#DC2626" : "#8B5CF6",
            },
            label: {
              show: true,
              formatter: "T5",
              fontSize: 8,
              color: "#FFFFFF",
              position: isBuy ? "bottom" : "top",
            },
            tooltip: {
              formatter: () => `
                <div style="font-family:monospace;font-size:10px;padding:2px;color:${colors.tooltipText};">
                  <div style="font-weight:bold;color:#FD6602;margin-bottom:3px;">T5 ${t.action || "-"} (${t.date})</div>
                  <div>Investor: <b>${t.investor || "-"}</b></div>
                  <div>Broker: <b>${t.broker || "-"}</b></div>
                  <div>Volume: <b>${Math.abs(t.lot).toLocaleString("id-ID")} lot</b> (${t.pct > 0 ? "+" : ""}${t.pct}%)</div>
                  <div>Harga: <b>Rp ${t.price}</b> (~${t.value_m.toFixed(1)} M)</div>
                </div>
              `,
            },
          });
        }
      });
    }

    return {
      backgroundColor: "transparent",
      animation: false,
      tooltip: showTooltip
        ? {
            trigger: "axis",
            axisPointer: { type: "cross" },
            backgroundColor: colors.tooltipBg,
            borderColor: colors.tooltipBorder,
            textStyle: { color: colors.tooltipText, fontSize: 10 },
            confine: true,
          }
        : { show: false },
      legend: {
        data: isBpMode ? ["Harga Saham", "Volume"] : ["Harga Saham", "TAVG", "Akumulasi Bandar", "Volume"],
        textStyle: { color: colors.legendText, fontSize: 10 },
        top: 0,
        type: "scroll",
      },
      grid: [
        { left: 10, right: 10, top: "14%", height: "52%", containLabel: true },
        { left: 10, right: 10, top: "72%", height: "20%", containLabel: true },
      ],
      xAxis: [
        {
          type: "category",
          data: dates,
          gridIndex: 0,
          scale: true,
          boundaryGap: true,
          axisLine: { lineStyle: { color: colors.axisLine } },
          axisLabel: { color: colors.axisLabel, fontSize: 9 },
        },
        {
          type: "category",
          data: dates,
          gridIndex: 1,
          scale: true,
          boundaryGap: true,
          axisLine: { lineStyle: { color: colors.axisLine } },
          axisLabel: { show: false },
        },
      ],
      yAxis: [
        {
          scale: true,
          gridIndex: 0,
          splitLine: { lineStyle: { color: colors.splitLine } },
          axisLabel: { color: colors.axisLabel, fontSize: 9 },
        },
        {
          scale: true,
          gridIndex: 0,
          position: "right",
          splitLine: { show: false },
          axisLabel: { color: colors.bandarColor, fontSize: 9 },
        },
        {
          scale: true,
          gridIndex: 1,
          splitLine: { lineStyle: { color: colors.splitLine } },
          axisLabel: { color: colors.axisLabel, fontSize: 8 },
        },
      ],
      series: [
        {
          name: "Harga Saham",
          type: "candlestick",
          data: candleValues,
          gridIndex: 0,
          itemStyle: {
            color: colors.upColor,
            color0: colors.downColor,
            borderColor: colors.upColor,
            borderColor0: colors.downColor,
          },
          markPoint: { data: markPoints },
        },
        ...(isBpMode
          ? []
          : [
              {
                name: "TAVG",
                type: "line",
                data: tavgLine,
                gridIndex: 0,
                smooth: true,
                lineStyle: { color: colors.tavgColor, width: 1.5 },
                itemStyle: { color: colors.tavgColor },
                showSymbol: false,
              },
              {
                name: "Akumulasi Bandar",
                type: "line",
                yAxisIndex: 1,
                data: bandarAcc,
                gridIndex: 0,
                smooth: true,
                lineStyle: { color: colors.bandarColor, width: 2 },
                itemStyle: { color: colors.bandarColor },
                showSymbol: false,
              },
            ]),
        {
          name: "Volume",
          type: "bar",
          xAxisIndex: 1,
          yAxisIndex: 2,
          data: volumes,
        },
      ],
    };
  }, [candleData, t5Data, showT5Markers, showTooltip, isBpMode, isDark]);

  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-3 relative shadow-2xs">
      <div className="flex items-center justify-between pb-2 border-b border-[var(--md-sys-color-outline-variant)]/60">
        {isBpMode ? (
          <>
            <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)]">
              CHART BIG PLAYER — {ticker || "—"}
            </span>
            <span className="text-[10px] text-[var(--md-sys-color-primary)] font-mono font-bold">
              T5 MARKERS
            </span>
          </>
        ) : (
          <>
            <div className="flex items-center gap-2">
              <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] tracking-wide">
                PRICE &amp; ACCUMULATION
              </span>
              <span className="text-[9px] font-mono font-semibold text-[var(--md-sys-color-on-surface-variant)] bg-[var(--md-sys-color-surface-container-high)] px-2 py-0.5 rounded">
                OHLCV + TAVG
              </span>
            </div>
            <div className="flex items-center gap-1.5">
              <button
                type="button"
                onClick={() => setShowT5Markers?.(!showT5Markers)}
                className={`px-2 py-0.5 rounded text-[10px] font-mono font-bold border transition cursor-pointer ${
                  showT5Markers
                    ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border-[var(--md-sys-color-primary)]"
                    : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]"
                }`}
                title="Toggle Marker T5"
              >
                BP
              </button>
              <button
                type="button"
                onClick={() => setShowTooltip?.(!showTooltip)}
                className={`px-2 py-0.5 rounded text-[10px] font-mono font-bold border transition cursor-pointer ${
                  showTooltip
                    ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border-[var(--md-sys-color-primary)]"
                    : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]"
                }`}
                title="Toggle Tooltip"
              >
                i
              </button>
            </div>
          </>
        )}
      </div>

      <div className={`${isBpMode ? "h-[360px] sm:h-[400px]" : "h-[320px] sm:h-[380px]"} w-full pt-2`}>
        <ReactECharts
          option={option}
          notMerge
          style={{ height: "100%", width: "100%" }}
          opts={{ renderer: "canvas" }}
        />
      </div>
    </div>
  );
}
