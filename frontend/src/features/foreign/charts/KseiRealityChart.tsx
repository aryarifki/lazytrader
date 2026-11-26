import { useMemo } from "react";
import ReactECharts from "echarts-for-react";
import { useAppStore } from "@/store/useAppStore";

export interface KseiHistoryItem {
  period_date: string;
  pct_asing: number;
  pct_asing_gt1: number;
  pct_asing_lt1: number;
  pct_ritel_pure: number;
  pct_local_institution: number;
}

export default function KseiRealityChart({
  history,
}: {
  history?: KseiHistoryItem[] | null;
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

  if (!history || history.length === 0) {
    return (
      <div className="flex h-full items-center justify-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
        Data kepemilikan bulanan KSEI tidak tersedia
      </div>
    );
  }

  const periods = history.map((h) => h.period_date.slice(0, 7));

  const option = {
    tooltip: {
      trigger: "axis",
      backgroundColor: styles.surfaceHigh,
      borderColor: styles.outlineVariant,
      textStyle: { color: styles.onSurface, fontSize: 11, fontFamily: "Inter, sans-serif", fontWeight: 600 },
      padding: [8, 12],
      borderRadius: 12,
      valueFormatter: (val: number) => `${Number(val).toFixed(2)}%`,
    },
    legend: {
      top: 0,
      right: 0,
      textStyle: { color: styles.onSurfaceVariant, fontSize: 10, fontWeight: 600 },
    },
    grid: { left: "2%", right: "2%", bottom: "4%", top: "16%", containLabel: true },
    xAxis: {
      type: "category",
      data: periods,
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
        formatter: (v: number) => `${v}%`,
      },
    },
    series: [
      {
        name: "Asing >1% (Named)",
        type: "bar",
        stack: "foreign",
        data: history.map((h) => h.pct_asing_gt1),
        itemStyle: { color: "#059669" },
      },
      {
        name: "Asing <1% (Hidden)",
        type: "bar",
        stack: "foreign",
        data: history.map((h) => h.pct_asing_lt1),
        itemStyle: { color: "#34d399", borderRadius: [4, 4, 0, 0] },
      },
      {
        name: "Ritel Murni (%)",
        type: "line",
        data: history.map((h) => h.pct_ritel_pure),
        smooth: true,
        lineStyle: { width: 2.5, color: "#f59e0b" },
        itemStyle: { color: "#f59e0b" },
      },
    ],
  };

  return <ReactECharts option={option} style={{ height: "100%", width: "100%", zIndex: 0 }} />;
}
