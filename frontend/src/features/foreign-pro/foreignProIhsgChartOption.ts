import { roundNum } from "./foreignProMath";
import type { IhsgChartInput } from "./foreignProTypes";

export function buildIhsgChartOption(input: IhsgChartInput): Record<string, unknown> {
  const {
    ihsgData,
    isLightMode,
    catIhsg,
    catKonglo,
    catFund,
    metFF,
    metFCo,
    metKF,
    metXL,
    metSR,
    flipKF,
    flipXL,
    flipSR,
  } = input;

  if (!ihsgData?.points || ihsgData.points.length === 0) return {};

  const pts = ihsgData.points;
  const dates = pts.map((p) => p.date);
  const closes = pts.map((p) => p.close);

  const colors = {
    tooltipBg: isLightMode ? "rgba(255, 255, 255, 0.96)" : "rgba(22, 19, 32, 0.95)",
    tooltipBorder: isLightMode ? "rgba(132, 116, 106, 0.35)" : "rgba(255,255,255,0.15)",
    tooltipText: isLightMode ? "#1E1A17" : "#F2F2F2",
    pointerBg: isLightMode ? "rgba(237, 223, 211, 0.96)" : "rgba(30, 24, 48, 0.95)",
    legendText: isLightMode ? "#4E453F" : "#A39DB5",
    axisLine: isLightMode ? "rgba(30, 26, 23, 0.18)" : "rgba(255,255,255,0.12)",
    axisLabel: isLightMode ? "#4E453F" : "#7A7290",
    splitLine: isLightMode ? "rgba(30, 26, 23, 0.08)" : "rgba(255,255,255,0.05)",
    compositeLine: isLightMode ? "#1E1A17" : "#FFFFFF",
    rightAxisLabel: isLightMode ? "#8C5000" : "#FF6820",
    sliderBg: isLightMode ? "rgba(30, 26, 23, 0.06)" : "rgba(255,255,255,0.05)",
    sliderFiller: isLightMode ? "rgba(140, 80, 0, 0.22)" : "rgba(255, 104, 32, 0.25)",
    sliderHandle: isLightMode ? "#8C5000" : "#FF6820",
    sliderText: isLightMode ? "#4E453F" : "#9E9E9E",
  };

  const seriesList: Array<Record<string, unknown>> = [
    {
      name: "IHSG Composite",
      type: "line",
      data: closes,
      smooth: true,
      lineStyle: { color: colors.compositeLine, width: 2 },
      showSymbol: false,
    },
  ];

  const buildCumLine = (
    catKey: string,
    metKey: string,
    isFlip: boolean,
    color: string,
    dash: number[] | null,
    width: number,
    label: string
  ) => {
    let runSum = 0;
    const dataArr = pts.map((p) => {
      const catObj = p[catKey] as Record<string, number | undefined> | undefined;
      const val = catObj?.[metKey] || 0.0;
      const adjusted = isFlip ? -val : val;
      runSum += adjusted;
      return roundNum(runSum, 2);
    });
    return {
      name: label,
      type: "line",
      yAxisIndex: 1,
      data: dataArr,
      smooth: true,
      lineStyle: { color, width, type: dash ? "dashed" : "solid" },
      showSymbol: false,
    };
  };

  const categories = [
    ...(catIhsg ? [{ key: "ihsg", label: "IHSG", dash: null, width: 2.0 }] : []),
    ...(catKonglo ? [{ key: "konglo", label: "Konglo", dash: [3, 2], width: 1.5 }] : []),
    ...(catFund ? [{ key: "fundamental", label: "Fundamental", dash: [1.5, 3], width: 1.8 }] : []),
  ];

  categories.forEach((cat) => {
    if (metFF) seriesList.push(buildCumLine(cat.key, "ff", false, "#FF6820", cat.dash, cat.width, `${cat.label} - Foreign Flow`));
    if (metFCo) seriesList.push(buildCumLine(cat.key, "fco", false, "#D62839", cat.dash, cat.width, `${cat.label} - Foreign & Co`));
    if (metKF) seriesList.push(buildCumLine(cat.key, "kf", flipKF, "#2F7676", cat.dash, cat.width, `${cat.label} - Korban Flow`));
    if (metXL) seriesList.push(buildCumLine(cat.key, "xl", flipXL, "#3A7BD5", cat.dash, cat.width, `${cat.label} - Broker XL`));
    if (metSR) seriesList.push(buildCumLine(cat.key, "sr", flipSR, "#1DB87A", cat.dash, cat.width, `${cat.label} - Semi Ritel`));
  });

  const startValue = dates.length > 120 ? dates[dates.length - 120] : dates[0];
  const endValue = dates[dates.length - 1];

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
        color: colors.tooltipText,
      },
    },
    legend: {
      type: "scroll",
      data: seriesList.map((s) => String(s.name)),
      textStyle: { color: colors.legendText, fontSize: 10 },
      top: 0,
    },
    grid: [{ left: 12, right: 12, top: "10%", height: "72%", containLabel: true }],
    dataZoom: [
      {
        type: "inside",
        startValue,
        endValue,
      },
      {
        type: "slider",
        show: true,
        bottom: 2,
        height: 20,
        borderColor: "transparent",
        backgroundColor: colors.sliderBg,
        fillerColor: colors.sliderFiller,
        handleStyle: { color: colors.sliderHandle },
        textStyle: { color: colors.sliderText, fontSize: 9 },
        startValue,
        endValue,
      },
    ],
    xAxis: [
      {
        type: "category",
        data: dates,
        axisLine: { lineStyle: { color: colors.axisLine } },
        axisLabel: {
          color: colors.axisLabel,
          fontSize: 10,
          formatter: (dt: string) => {
            if (!dt) return "";
            const parts = dt.split("-");
            return parts.length === 3 ? `${parts[2]}/${parts[1]}` : dt;
          },
        },
      },
    ],
    yAxis: [
      {
        name: "IHSG Level",
        type: "value",
        scale: true,
        splitLine: { lineStyle: { color: colors.splitLine } },
        axisLabel: {
          color: colors.axisLabel,
          fontSize: 10,
          formatter: (v: number) => Math.round(v).toLocaleString("id-ID"),
        },
      },
      {
        name: "Cumulative Flow (Triliun Rp)",
        type: "value",
        scale: true,
        position: "right",
        splitLine: { show: false },
        axisLabel: {
          color: colors.rightAxisLabel,
          fontSize: 10,
          formatter: (v: number) => `${v >= 0 ? "+" : ""}${(v / 1000).toFixed(1)} T`,
        },
      },
    ],
    series: seriesList,
  };
}
