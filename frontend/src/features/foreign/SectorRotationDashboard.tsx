import React, { useState, useMemo, useRef, useEffect } from "react";
import useSWR from "swr";
import { Icon } from "@iconify/react";
import { useAppStore } from "@/store/useAppStore";
import {
  SECTOR_STORAGE_KEY,
  TIMEFRAME_STORAGE_KEY,
  QUADRANT_STORAGE_KEY,
} from "./sectorHelpers";
import SectorMacroOverview from "./SectorMacroOverview";
import SectorFilterBar from "./SectorFilterBar";
import SectorCard from "./SectorCard";
import SectorRadarDrawer from "./SectorRadarDrawer";
import type { SectorHistoryResponse, SectorSummaryResponse } from "./sectorTypes";

const fetcher = (url: string) =>
  fetch(url, { credentials: "include" }).then((res) => {
    if (!res.ok) throw new Error("Gagal mengambil data rotasi sektor");
    return res.json();
  });

export default function SectorRotationDashboard() {
  const [timeframe, setTimeframe] = useState<"1d" | "5d" | "20d">("1d");
  const [selectedQuadrant, setSelectedQuadrant] = useState<string>("all");
  const [selectedSector, setSelectedSector] = useState<string | null>("IDXTECHNO");
  const [isRefreshing, setIsRefreshing] = useState(false);

  const radarDetailRef = useRef<HTMLDivElement>(null);

  // Restore state from URL search params or localStorage on mount
  useEffect(() => {
    try {
      const params = new URLSearchParams(window.location.search);
      const urlSector = params.get("sector");
      const savedSector = localStorage.getItem(SECTOR_STORAGE_KEY);
      if (urlSector) {
        setSelectedSector(urlSector);
      } else if (savedSector !== null && savedSector !== "") {
        setSelectedSector(savedSector);
      }

      const savedTf = localStorage.getItem(TIMEFRAME_STORAGE_KEY);
      if (savedTf === "1d" || savedTf === "5d" || savedTf === "20d") {
        setTimeframe(savedTf);
      }

      const savedQ = localStorage.getItem(QUADRANT_STORAGE_KEY);
      if (savedQ) {
        setSelectedQuadrant(savedQ);
      }
    } catch {}
  }, []);

  const handleSelectSector = (sec: string | null) => {
    setSelectedSector(sec);
    try {
      if (sec) {
        localStorage.setItem(SECTOR_STORAGE_KEY, sec);
        const url = new URL(window.location.href);
        url.searchParams.set("sector", sec);
        window.history.replaceState(null, "", url.toString());
      } else {
        localStorage.removeItem(SECTOR_STORAGE_KEY);
        const url = new URL(window.location.href);
        url.searchParams.delete("sector");
        window.history.replaceState(null, "", url.toString());
      }
    } catch {}
  };

  const handleTimeframeChange = (tf: "1d" | "5d" | "20d") => {
    setTimeframe(tf);
    try {
      localStorage.setItem(TIMEFRAME_STORAGE_KEY, tf);
    } catch {}
  };

  const handleQuadrantChange = (q: string) => {
    setSelectedQuadrant(q);
    try {
      localStorage.setItem(QUADRANT_STORAGE_KEY, q);
    } catch {}
  };

  const { data, isLoading, error, mutate } = useSWR<SectorSummaryResponse>(
    "/api/foreign-flow/sectors/summary",
    fetcher,
    { revalidateOnMount: true, revalidateOnFocus: true, dedupingInterval: 5000 }
  );

  const { data: historyData, isLoading: isHistoryLoading, mutate: mutateHistory } = useSWR<SectorHistoryResponse>(
    selectedSector ? `/api/foreign-flow/sectors/history?sector=${selectedSector}&days=30` : null,
    fetcher,
    { revalidateOnMount: true, revalidateOnFocus: true, dedupingInterval: 5000 }
  );

  const handleRefresh = async () => {
    setIsRefreshing(true);
    try {
      await Promise.all([mutate(), mutateHistory()]);
    } finally {
      setTimeout(() => setIsRefreshing(false), 600);
    }
  };

  const lastSyncTimestamp = useAppStore((state) => state.lastSyncTimestamp);

  // Sinkronisasi instan saat tombol Master Sync diklik
  useEffect(() => {
    if (lastSyncTimestamp > 0) {
      handleRefresh();
    }
  }, [lastSyncTimestamp]);

  useEffect(() => {
    const handleSyncEvent = () => {
      handleRefresh();
    };
    window.addEventListener("investowl:sync", handleSyncEvent);
    return () => window.removeEventListener("investowl:sync", handleSyncEvent);
  }, []);

  // Smooth scroll into radar detail when a card is selected
  useEffect(() => {
    if (selectedSector && radarDetailRef.current) {
      const timer = setTimeout(() => {
        radarDetailRef.current?.scrollIntoView({ behavior: "smooth", block: "nearest" });
      }, 80);
      return () => clearTimeout(timer);
    }
  }, [selectedSector]);

  const filteredSectors = useMemo(() => {
    if (!data?.sectors) return [];
    let list = [...data.sectors];

    if (selectedQuadrant !== "all") {
      list = list.filter((s) => s.quadrant.toLowerCase() === selectedQuadrant.toLowerCase());
    }

    list.sort((a, b) => {
      if (timeframe === "5d") return b.ret_5d - a.ret_5d;
      if (timeframe === "20d") return b.ret_20d - a.ret_20d;
      return b.ret_1d - a.ret_1d;
    });

    return list;
  }, [data?.sectors, selectedQuadrant, timeframe]);

  // Quadrant counts
  const quadrantCounts = useMemo(() => {
    if (!data?.sectors) return { Leading: 0, Improving: 0, Weakening: 0, Lagging: 0 };
    return data.sectors.reduce(
      (acc, s) => {
        acc[s.quadrant] = (acc[s.quadrant] || 0) + 1;
        return acc;
      },
      { Leading: 0, Improving: 0, Weakening: 0, Lagging: 0 } as Record<string, number>
    );
  }, [data?.sectors]);

  // Foreign Exodus state counts across all 11 sectors
  const foreignStateCounts = useMemo(() => {
    if (!data?.sectors) return { surge: 0, warning: 0, exodus: 0, bottoming: 0 };
    return data.sectors.reduce(
      (acc, s) => {
        if (s.foreign_state === "INFLOW SURGE") acc.surge += 1;
        else if (s.foreign_state === "HIDDEN DISTRIBUTION") acc.warning += 1;
        else if (s.foreign_state === "SECTOR EXODUS") acc.exodus += 1;
        else if (s.foreign_state === "BOTTOMING") acc.bottoming += 1;
        return acc;
      },
      { surge: 0, warning: 0, exodus: 0, bottoming: 0 }
    );
  }, [data?.sectors]);

  // Total Foreign Net 1D with comprehensive fallback
  const computedTotalForeign1d = useMemo(() => {
    if (data?.total_foreign_net_1d != null && !Number.isNaN(data.total_foreign_net_1d)) {
      return data.total_foreign_net_1d;
    }
    if (data?.sectors && data.sectors.length > 0) {
      return data.sectors.reduce((acc, s) => acc + (s.foreign_net_1d || 0), 0);
    }
    return 0;
  }, [data]);

  const activeSectorDetail = useMemo(() => {
    if (!selectedSector || !data?.sectors) return null;
    return data.sectors.find((s) => s.code === selectedSector) || null;
  }, [selectedSector, data?.sectors]);

  // Dual-axis history data processing
  const historyTimeline = useMemo(() => {
    return historyData?.data || [];
  }, [historyData]);

  // 1D and 5D Net Foreign flow with multi-source fallback
  const fn1d = useMemo(() => {
    if (activeSectorDetail?.foreign_net_1d != null && !Number.isNaN(activeSectorDetail.foreign_net_1d)) {
      return activeSectorDetail.foreign_net_1d;
    }
    if (historyTimeline.length > 0) {
      return historyTimeline[historyTimeline.length - 1].foreign_net ?? null;
    }
    return null;
  }, [activeSectorDetail?.foreign_net_1d, historyTimeline]);

  const fn5d = useMemo(() => {
    if (activeSectorDetail?.foreign_net_5d != null && !Number.isNaN(activeSectorDetail.foreign_net_5d)) {
      return activeSectorDetail.foreign_net_5d;
    }
    if (historyTimeline.length >= 5) {
      return historyTimeline.slice(-5).reduce((acc, curr) => acc + (curr.foreign_net || 0), 0);
    }
    if (historyTimeline.length > 0) {
      return historyTimeline.reduce((acc, curr) => acc + (curr.foreign_net || 0), 0);
    }
    return null;
  }, [activeSectorDetail?.foreign_net_5d, historyTimeline]);

  const isNetOutflow = useMemo(() => {
    if (!historyTimeline.length) return false;
    return historyTimeline[historyTimeline.length - 1].cumulative_foreign_flow < 0;
  }, [historyTimeline]);

  const gradientOffset = useMemo(() => {
    if (!historyTimeline.length) return 0;
    const flows = historyTimeline.map((d) => d.cumulative_foreign_flow);
    const max = Math.max(...flows);
    const min = Math.min(...flows);
    if (max <= 0) return 0;
    if (min >= 0) return 1;
    return max / (max - min);
  }, [historyTimeline]);

  if (isLoading) {
    return (
      <div className="py-24 text-center bg-[var(--md-sys-color-surface-container)] rounded-2xl border border-[var(--md-sys-color-outline-variant)] shadow-xs flex flex-col items-center justify-center gap-3">
        <Icon icon="ph:spinner-gap-bold" className="animate-spin text-[var(--md-sys-color-primary)]" width="28" />
        <span className="text-xs font-semibold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
          Menganalisis rotasi sektor &amp; radar modal asing...
        </span>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="p-6 text-center bg-[var(--md-sys-color-surface-container-low)] rounded-2xl border border-[var(--md-sys-color-outline-variant)] space-y-2">
        <Icon icon="ph:warning-circle-duotone" width="36" className="text-amber-500 mx-auto" />
        <h3 className="text-sm font-bold text-[var(--md-sys-color-on-surface)]">Gagal Memuat Data Rotasi Sektor</h3>
        <p className="text-xs text-[var(--md-sys-color-on-surface-variant)]">
          Pastikan data `idx_index_summary` dan `broker_flow` telah tersedia di database.
        </p>
      </div>
    );
  }

  const { benchmark, date } = data;

  return (
    <div className="space-y-5">
      {/* MARKET MACRO & ROTATION OVERVIEW CARD */}
      <SectorMacroOverview
        date={date}
        benchmark={benchmark}
        macroMarketState={data?.macro_market_state}
        computedTotalForeign1d={computedTotalForeign1d}
        selectedQuadrant={selectedQuadrant}
        onSelectQuadrant={setSelectedQuadrant}
        quadrantCounts={quadrantCounts}
        foreignStateCounts={foreignStateCounts}
        isRefreshing={isRefreshing}
        onRefresh={handleRefresh}
      />

      {/* FILTER & TIMEFRAME BAR */}
      <SectorFilterBar
        selectedQuadrant={selectedQuadrant}
        onQuadrantChange={handleQuadrantChange}
        quadrantCounts={quadrantCounts}
        timeframe={timeframe}
        onTimeframeChange={handleTimeframeChange}
      />

      {/* SECTOR CARDS GRID WITH INLINE DYNAMIC RADAR DRAWER */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3.5">
        {filteredSectors.map((sector) => {
          const isSelected = selectedSector === sector.code;
          return (
            <React.Fragment key={sector.code}>
              <SectorCard
                sector={sector}
                timeframe={timeframe}
                isSelected={isSelected}
                onSelect={() => handleSelectSector(isSelected ? null : sector.code)}
              />

              {/* INLINE DYNAMIC RADAR DRAWER LANGSUNG MUNCUL DIBAWAH KARTU SEKTOR */}
              {isSelected && activeSectorDetail && (
                <SectorRadarDrawer
                  drawerRef={radarDetailRef}
                  activeSectorDetail={activeSectorDetail}
                  historyData={historyData}
                  historyTimeline={historyTimeline}
                  isHistoryLoading={isHistoryLoading}
                  fn1d={fn1d}
                  fn5d={fn5d}
                  isNetOutflow={isNetOutflow}
                  gradientOffset={gradientOffset}
                  onClose={() => handleSelectSector(null)}
                />
              )}
            </React.Fragment>
          );
        })}
      </div>
    </div>
  );
}
