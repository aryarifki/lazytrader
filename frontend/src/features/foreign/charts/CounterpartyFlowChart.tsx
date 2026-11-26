import { useMemo } from "react";
import ReactECharts from "echarts-for-react";
import { useAppStore } from "@/store/useAppStore";

export interface CounterpartyData {
  dates: string[];
  cum_foreign: number[];
  cum_local_inst: number[];
  cum_retail: number[];
}

function formatBillion(val: number): string {
  const abs = Math.abs(val);
  const sign = val < 0 ? "-" : "+";
  if (abs >= 1e12) return `${sign}Rp ${(abs / 1e12).toFixed(2)} T`;
  if (abs >= 1e9) return `${sign}Rp ${(abs / 1e9).toFixed(1)} M`;
  if (abs >= 1e6) return `${sign}Rp ${(abs / 1e6).toFixed(0)} Jt`;
  return `${sign}Rp ${abs.toLocaleString("id-ID")}`;
}

export default function CounterpartyFlowChart({
  data,
}: {
  data?: CounterpartyData | null;
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

  if (!data || !data.dates || data.dates.length === 0) {
    return (
      <div className="flex h-full items-center justify-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
        Data counterparty tidak tersedia
      </div>
    );
  }

  const option = {
    tooltip: {
      trigger: "axis",
      backgroundColor: styles.surfaceHigh,
      borderColor: styles.outlineVariant,
      textStyle: { color: styles.onSurface, fontSize: 11, fontFamily: "Inter, sans-serif", fontWeight: 600 },
      padding: [8, 12],
      borderRadius: 12,
      formatter: (params: Array<{ seriesName: string; value: number; marker: string; axisValue: string }>) => {
        if (!params || params.length === 0) return "";
        const rows = params.map(
          (p) => `${p.marker} ${p.seriesName}: <b>${formatBillion(Number(p.value || 0))}</b>`
        );
        return `<div style="font-weight:800;margin-bottom:4px;">${params[0].axisValue}</div>${rows.join("<br/>")}`;
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
      data: data.dates,
      axisLabel: { color: styles.onSurfaceVariant, fontSize: 10, fontWeight: 600 },
    },
    yAxis: {
      type: "value",
      scale: true,
      splitLine: { lineStyle: { color: styles.splitLine } },
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
    series: [
      {
        name: "Asing (Foreign)",
        type: "line",
        data: data.cum_foreign,
        smooth: true,
        symbol: "none",
        lineStyle: { width: 2.5, color: "#10b981" },
        itemStyle: { color: "#10b981" },
      },
      {
        name: "Institusi Lokal",
        type: "line",
        data: data.cum_local_inst,
        smooth: true,
        symbol: "none",
        lineStyle: { width: 2, color: "#6366f1" },
        itemStyle: { color: "#6366f1" },
      },
      {
        name: "Ritel",
        type: "line",
        data: data.cum_retail,
        smooth: true,
        symbol: "none",
        lineStyle: { width: 2, type: "dashed", color: "#f43f5e" },
        itemStyle: { color: "#f43f5e" },
      },
    ],
  };

  return <ReactECharts option={option} style={{ height: "100%", width: "100%", zIndex: 0 }} />;
}
