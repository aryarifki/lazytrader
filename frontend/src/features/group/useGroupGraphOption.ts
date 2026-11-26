import { useMemo } from "react";
import type { EChartsOption } from "echarts";
import type { NetworkData } from "./groupTypes";

export function useGroupGraphOption(network: NetworkData | undefined): EChartsOption {
  return useMemo<EChartsOption>(() => {
    if (!network || network.nodes.length === 0) return {};
    const nameById = new Map(network.nodes.map((node) => [node.id, node.label]));
    const categoryIndex: Record<string, number> = { hub: 0, emiten: 1, owner: 2 };
    const nodes = network.nodes.map((node) => ({
      name: node.label,
      value: node.val ?? 1,
      category: categoryIndex[node.type] ?? 1,
      symbolSize: node.type === "hub" ? 52 : node.type === "emiten" ? 30 : 18,
      itemStyle: {
        color:
          node.type === "hub"
            ? "#B45309"
            : node.type === "emiten"
              ? "var(--md-sys-color-primary)"
              : "#3A7BD5",
        borderColor: "var(--md-sys-color-surface-container)",
        borderWidth: 2,
      },
      label: {
        show: true,
        position: "right" as const,
        fontSize: 9,
        color: "var(--md-sys-color-on-surface)",
      },
    }));
    const links = network.links.map((link) => ({
      source: nameById.get(link.source) ?? link.source,
      target: nameById.get(link.target) ?? link.target,
      value: link.pct,
      lineStyle: {
        width: Math.min(4, 1 + Math.max(0, link.pct) / 12),
        opacity: link.is_controlling ? 0.8 : 0.45,
        curveness: 0.12,
        color: link.is_controlling ? "#B45309" : "#64748B",
      },
      label: {
        show: link.pct >= 5,
        formatter: `${link.pct.toFixed(1)}%`,
        fontSize: 8,
        color: "var(--md-sys-color-on-surface-variant)",
      },
    }));

    return {
      backgroundColor: "transparent",
      tooltip: { trigger: "item" },
      legend: [
        {
          data: ["Hub", "Emiten", "Pemegang"],
          bottom: 0,
          textStyle: { fontSize: 10, color: "var(--md-sys-color-on-surface-variant)" },
          itemWidth: 10,
          itemHeight: 10,
        },
      ],
      series: [
        {
          type: "graph",
          layout: "force",
          roam: true,
          draggable: true,
          categories: [{ name: "Hub" }, { name: "Emiten" }, { name: "Pemegang" }],
          data: nodes,
          links,
          force: { repulsion: 190, edgeLength: [70, 145], gravity: 0.08 },
          emphasis: { focus: "adjacency", lineStyle: { width: 4 } },
          lineStyle: { color: "source", curveness: 0.12 },
          labelLayout: { hideOverlap: true },
        },
      ],
    };
  }, [network]);
}
