import { useMemo } from "react";
import type React from "react";

import type { ActiveFilterItem, Filters, Mode } from "./snapshotTypes";

export function useActiveFilters({
  condition,
  conditionLabel,
  deferredQuery,
  filters,
  mode,
  setCondition,
  setQuery,
  setFilters,
}: {
  condition: string;
  conditionLabel: string | undefined;
  deferredQuery: string;
  filters: Filters;
  mode: Mode;
  setCondition: (value: string) => void;
  setQuery: (value: string) => void;
  setFilters: React.Dispatch<React.SetStateAction<Filters>>;
}): ActiveFilterItem[] {
  return useMemo(() => {
    const items: ActiveFilterItem[] = [];
    if (condition !== "all") {
      items.push({
        key: "condition",
        label: conditionLabel ?? condition,
        clear: () => setCondition("all"),
      });
    }
    if (deferredQuery) {
      items.push({
        key: "q",
        label: `Cari: ${deferredQuery}`,
        clear: () => setQuery(""),
      });
    }
    if (mode === "stocks") {
      if (filters.sector) {
        items.push({
          key: "sector",
          label: `Sektor: ${filters.sector}`,
          clear: () => setFilters((prev) => ({ ...prev, sector: "" })),
        });
      }
      if (filters.minHolders) {
        items.push({
          key: "minHolders",
          label: `Holders ≥ ${filters.minHolders}`,
          clear: () => setFilters((prev) => ({ ...prev, minHolders: "" })),
        });
      }
      if (filters.maxHolders) {
        items.push({
          key: "maxHolders",
          label: `Holders ≤ ${filters.maxHolders}`,
          clear: () => setFilters((prev) => ({ ...prev, maxHolders: "" })),
        });
      }
      if (filters.minControlling) {
        items.push({
          key: "minControlling",
          label: `Pengendali ≥ ${filters.minControlling}%`,
          clear: () => setFilters((prev) => ({ ...prev, minControlling: "" })),
        });
      }
      if (filters.maxControlling) {
        items.push({
          key: "maxControlling",
          label: `Pengendali ≤ ${filters.maxControlling}%`,
          clear: () => setFilters((prev) => ({ ...prev, maxControlling: "" })),
        });
      }
      if (filters.minPublic) {
        items.push({
          key: "minPublic",
          label: `Publik ≥ ${filters.minPublic}%`,
          clear: () => setFilters((prev) => ({ ...prev, minPublic: "" })),
        });
      }
      if (filters.maxPublic) {
        items.push({
          key: "maxPublic",
          label: `Publik ≤ ${filters.maxPublic}%`,
          clear: () => setFilters((prev) => ({ ...prev, maxPublic: "" })),
        });
      }
      if (filters.minTopHolder) {
        items.push({
          key: "minTopHolder",
          label: `Top holder ≥ ${filters.minTopHolder}%`,
          clear: () => setFilters((prev) => ({ ...prev, minTopHolder: "" })),
        });
      }
      if (filters.maxTopHolder) {
        items.push({
          key: "maxTopHolder",
          label: `Top holder ≤ ${filters.maxTopHolder}%`,
          clear: () => setFilters((prev) => ({ ...prev, maxTopHolder: "" })),
        });
      }
      if (filters.minMarketCap) {
        items.push({
          key: "minMarketCap",
          label: `Mkt cap ≥ ${filters.minMarketCap} M`,
          clear: () => setFilters((prev) => ({ ...prev, minMarketCap: "" })),
        });
      }
      if (filters.maxMarketCap) {
        items.push({
          key: "maxMarketCap",
          label: `Mkt cap ≤ ${filters.maxMarketCap} M`,
          clear: () => setFilters((prev) => ({ ...prev, maxMarketCap: "" })),
        });
      }
      if (filters.mtd !== "all") {
        items.push({
          key: "mtd",
          label: filters.mtd === "up" ? "MTD naik" : "MTD turun",
          clear: () => setFilters((prev) => ({ ...prev, mtd: "all" })),
        });
      }
      if (filters.foreignFlow !== "all") {
        items.push({
          key: "foreignFlow",
          label: filters.foreignFlow === "inflow" ? "Foreign inflow" : "Foreign outflow",
          clear: () => setFilters((prev) => ({ ...prev, foreignFlow: "all" })),
        });
      }
      if (filters.hasTreasury) {
        items.push({
          key: "hasTreasury",
          label: "Ada treasury",
          clear: () => setFilters((prev) => ({ ...prev, hasTreasury: false })),
        });
      }
    } else {
      if (filters.minTickers) {
        items.push({
          key: "minTickers",
          label: `Emiten ≥ ${filters.minTickers}`,
          clear: () => setFilters((prev) => ({ ...prev, minTickers: "" })),
        });
      }
      if (filters.maxTickers) {
        items.push({
          key: "maxTickers",
          label: `Emiten ≤ ${filters.maxTickers}`,
          clear: () => setFilters((prev) => ({ ...prev, maxTickers: "" })),
        });
      }
      if (filters.minMaxPct) {
        items.push({
          key: "minMaxPct",
          label: `Porsi max ≥ ${filters.minMaxPct}%`,
          clear: () => setFilters((prev) => ({ ...prev, minMaxPct: "" })),
        });
      }
      if (filters.minSectorCount) {
        items.push({
          key: "minSectorCount",
          label: `Sektor ≥ ${filters.minSectorCount}`,
          clear: () => setFilters((prev) => ({ ...prev, minSectorCount: "" })),
        });
      }
    }
    return items;
  }, [condition, conditionLabel, deferredQuery, filters, mode, setCondition, setQuery, setFilters]);
}
