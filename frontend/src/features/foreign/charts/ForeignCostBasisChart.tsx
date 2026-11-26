import { useMemo } from "react";
import ReactECharts from "echarts-for-react";
import { useAppStore } from "@/store/useAppStore";

export interface CostBasisTimeSeries {
  dates: string[];
  close_prices: number[];
  foreign_vwap: (number | null)[];
  cum_foreign_net: number[];
}

function formatIdrCompact(val: number): string {
  const abs = Math.abs(val);
  const sign = val < 0 ? "-" : "+";
  if (abs >= 1e12) return `${sign}Rp ${(abs / 1e12).toFixed(2)} T`;
  if (abs >= 1e9) return `${sign}Rp ${(abs / 1e9).toFixed(1)} M`;
  if (abs >= 1e6) return `${sign}Rp ${(abs / 1e6).toFixed(0)} Jt`;
  return `${sign}Rp ${abs.toLocaleString("id-ID")}`;
}

export default function ForeignCostBasisChart({
  timeSeries,
}: {
  timeSeries?: CostBasisTimeSeries | null;
}) {
  const theme = useAppStore((state) => state.theme);
  const isDark = theme === "dark";

  const styles = useMemo(
    () => ({
      surfaceHigh: isDark ? "#2B2930" : "#ECE6F0",
      onSurface: isDark ? "#E6E0E9" : "#1D1B20",
      onSurfaceVariant: isDark ? "#CAC4D0" : "#49454F",
      outlineVariant: isDark ? "#49454F" : "#CAC4D0",
      splitLine: isDark ? "rgba(255,255,255,0.08)" : "rgba(0,0,0,0.06)",
    }),
    [isDark]
  );

  if (!timeSeries || !timeSeries.dates || timeSeries.dates.length === 0) {
    return (
      <div className="flex h-full items-center justify-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
        Data harga &amp; modal asing tidak tersedia
      </div>
    );
  }

  const { dates, close_prices, foreign_vwap, cum_foreign_net } = timeSeries;

  const option = {
    tooltip: {
      trigger: "axis",
      backgroundColor: styles.surfaceHigh,
      borderColor: styles.outlineVariant,
      textStyle: {
        color: styles.onSurface,
        fontSize: 11,
        fontFamily: "Inter, sans-serif",
        fontWeight: 600,
      },
      padding: [8, 12],
      borderRadius: 12,
      formatter: (params: Array<{ seriesName: string; value: number | null; marker: string; axisValue: string }>) => {
        if (!params || params.length === 0) return "";
        const dateLabel = params[0].axisValue;
        const lines = params.map((p) => {
          if (p.value === null || p.value === undefined) return "";
          if (p.seriesName === "Cum. Foreign Flow") {
            return `${p.marker} ${p.seriesName}: <b>${formatIdrCompact(Number(p.value))}</b>`;
          }
          return `${p.marker} ${p.seriesName}: <b>Rp ${Number(p.value).toLocaleString("id-ID", { maximumFractionDigits: 0 })}</b>`;
        });
        return `<div style="font-weight:800;margin-bottom:4px;">${dateLabel}</div>${lines.filter(Boolean).join("<br/>")}`;
      },
    },
    legend: {
      top: 0,
      right: 0,
      textStyle: { color: styles.onSurfaceVariant, fontSize: 10, fontWeight: 600 },
    },
    grid: { left: "2%", right: "3%", bottom: "4%", top: "14%", containLabel: true },
    xAxis: {
      type: "category",
      boundaryGap: false,
      data: dates,
      axisLabel: { color: styles.onSurfaceVariant, fontSize: 10, fontWeight: 600 },
    },
    yAxis: [
      {
        type: "value",
        scale: true,
        splitLine: { lineStyle: { color: styles.splitLine } },
        axisLabel: {
          color: styles.onSurfaceVariant,
          fontSize: 10,
          fontWeight: 600,
          formatter: (v: number) => v.toLocaleString("id-ID"),
        },
      },
      {
        type: "value",
        scale: true,
        splitLine: { show: false },
        axisLabel: {
          color: styles.onSurfaceVariant,
          fontSize: 10,
          fontWeight: 600,
          formatter: (v: number) => {
            const abs = Math.abs(v);
            if (abs >= 1e12) return `${(v / 1e12).toFixed(1)}T`;
            if (abs >= 1e9) return `${(v / 1e9).toFixed(0)}M`;
            if (abs >= 1e6) return `${(v / 1e6).toFixed(0)}Jt`;
            return `${v}`;
          },
        },
      },
    ],
    series: [
      {
        name: "Cum. Foreign Flow",
        type: "line",
        yAxisIndex: 1,
        data: cum_foreign_net,
        smooth: true,
        symbol: "none",
        lineStyle: { width: 1.5, color: "#10b981" },
        itemStyle: { color: "#10b981" },
        areaStyle: {
          color: isDark ? "rgba(16, 185, 129, 0.14)" : "rgba(16, 185, 129, 0.12)",
        },
      },
      {
        name: "Close Price",
        type: "line",
        yAxisIndex: 0,
        data: close_prices,
        smooth: true,
        symbol: "none",
        lineStyle: { width: 2.5, color: "#f59e0b" },
        itemStyle: { color: "#f59e0b" },
      },
      {
        name: "Foreign VWAP",
        type: "line",
        yAxisIndex: 0,
        data: foreign_vwap,
        smooth: true,
        symbol: "none",
        connectNulls: true,
        lineStyle: { width: 2, type: "dashed", color: "#06b6d4" },
        itemStyle: { color: "#06b6d4" },
      },
    ],
  };

  return <ReactECharts option={option} style={{ height: "100%", width: "100%", zIndex: 0 }} />;
}
