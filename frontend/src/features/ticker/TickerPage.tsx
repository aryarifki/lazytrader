import { Suspense, lazy, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import useSWR from "swr";
import { Icon } from "@iconify/react";
import TickerSearch from "@/components/layout/TickerSearch";
import { useSession } from "@/lib/auth";
import { useAppStore } from "@/store/useAppStore";
import TickerHeader from "./TickerHeader";
import OverviewTab from "./OverviewTab";
import type { TickerDetailData } from "./types";

function TabSkeleton() {
  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-6 animate-pulse space-y-4">
      <div className="flex items-center justify-between">
        <div className="h-5 w-40 bg-[var(--md-sys-color-surface-container-highest)] rounded-md" />
        <div className="h-4 w-24 bg-[var(--md-sys-color-surface-container-highest)] rounded-full" />
      </div>
      <div className="h-64 bg-[var(--md-sys-color-surface-container-high)] rounded-xl" />
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div className="h-16 bg-[var(--md-sys-color-surface-container-high)] rounded-lg" />
        <div className="h-16 bg-[var(--md-sys-color-surface-container-high)] rounded-lg" />
        <div className="h-16 bg-[var(--md-sys-color-surface-container-high)] rounded-lg" />
      </div>
    </div>
  );
}

const BrokerFlowTab = lazy(() => import("./BrokerFlowTab"));
const CausalityTab = lazy(() => import("./CausalityTab"));
const ValidationTab = lazy(() => import("./ValidationTab"));
const RawTablesTab = lazy(() => import("./RawTablesTab"));

const fetcher = (url: string) => fetch(url, { credentials: "include" }).then((r) => r.json());

const TABS = [
  { id: "Overview", label: "Overview", icon: "ph:chart-line-up-duotone" },
  { id: "Broker Flow", label: "Broker Flow", icon: "ph:arrows-left-right-duotone" },
  { id: "Causality", label: "Causality", icon: "ph:git-commit-duotone" },
  { id: "Validation", label: "Validation", icon: "ph:check-circle-duotone" },
  { id: "Raw Tables", label: "Raw Records", icon: "ph:table-duotone" },
] as const;

type TabId = (typeof TABS)[number]["id"];

export default function TickerPage() {
  const params = useParams<{ ticker: string }>();
  const tickerFromUrl = String(params.ticker || "").toUpperCase();
  const { data: session } = useSession();

  const {
    activeTicker,
    setActiveTicker,
    analysisDate,
    windowDays,
    universe,
    horizon,
    minEvents,
    localWatchlist,
    setLocalWatchlist,
  } = useAppStore();

  const effectiveTicker = (tickerFromUrl || activeTicker || "").toUpperCase();
  const [activeTab, setActiveTab] = useState<TabId>("Overview");
  const [visitedTabs, setVisitedTabs] = useState<Record<TabId, boolean>>({
    Overview: true,
    "Broker Flow": false,
    Causality: false,
    Validation: false,
    "Raw Tables": false,
  });

  useEffect(() => {
    try {
      const searchParams = new URLSearchParams(window.location.search);
      const urlTab = searchParams.get("tab");
      const savedTab = localStorage.getItem("investowl_ticker_active_tab");

      const validTabIds = TABS.map((t) => t.id);
      let target: TabId | undefined;
      if (urlTab) {
        target = validTabIds.find(
          (id) =>
            id.toLowerCase() === urlTab.toLowerCase() ||
            id.toLowerCase().replace(/\s+/g, "-") === urlTab.toLowerCase()
        );
      }
      if (!target && savedTab) {
        target = validTabIds.find((id) => id === savedTab);
      }
      if (target) {
        setActiveTab(target);
        setVisitedTabs((prev) => ({ ...prev, [target]: true }));
      }
    } catch {}
  }, []);

  const handleTabChange = (tabId: TabId) => {
    setActiveTab(tabId);
    setVisitedTabs((prev) => (prev[tabId] ? prev : { ...prev, [tabId]: true }));
    try {
      localStorage.setItem("investowl_ticker_active_tab", tabId);
      const url = new URL(window.location.href);
      url.searchParams.set("tab", tabId.toLowerCase().replace(/\s+/g, "-"));
      window.history.replaceState(null, "", url.toString());
    } catch {}
  };

  useEffect(() => {
    if (tickerFromUrl && tickerFromUrl !== activeTicker) {
      setActiveTicker(tickerFromUrl);
    }
  }, [tickerFromUrl, activeTicker, setActiveTicker]);

  useEffect(() => {
    try {
      const ls = localStorage.getItem("investowl_watchlist") || localStorage.getItem("tradepulse_watchlist");
      if (ls) {
        const parsed = JSON.parse(ls);
        if (Array.isArray(parsed) && parsed.length > 0) {
          setLocalWatchlist(parsed);
        }
      }
    } catch {}

    const userEmail = session?.user?.email;
    if (userEmail) {
      fetch(`/api/users/${encodeURIComponent(userEmail)}/watchlist`, { credentials: "include" })
        .then((res) => res.json())
        .then((d) => {
          if (Array.isArray(d?.watchlist) && d.watchlist.length > 0) {
            setLocalWatchlist(d.watchlist);
            try {
              localStorage.setItem("investowl_watchlist", JSON.stringify(d.watchlist));
              localStorage.setItem("tradepulse_watchlist", JSON.stringify(d.watchlist));
            } catch {}
          }
        })
        .catch(() => {});
    }
  }, [setLocalWatchlist, session?.user?.email]);

  const isBookmarked = localWatchlist.includes(effectiveTicker);

  const toggleBookmark = () => {
    const nextList = isBookmarked
      ? localWatchlist.filter((t) => t !== effectiveTicker)
      : [...localWatchlist, effectiveTicker];
    setLocalWatchlist(nextList);
    try {
      localStorage.setItem("investowl_watchlist", JSON.stringify(nextList));
      localStorage.setItem("tradepulse_watchlist", JSON.stringify(nextList));
    } catch {}

    const userEmail = session?.user?.email;
    if (userEmail) {
      fetch(`/api/users/${encodeURIComponent(userEmail)}/watchlist`, {
        method: "PUT",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ watchlist: nextList }),
      }).catch((e) => console.error("Gagal sinkronisasi bookmark ke server:", e));
    }
  };

  const qs = `?window_days=${windowDays}${analysisDate ? `&analysis_date=${analysisDate}` : ""}`;
  const { data, error, isLoading } = useSWR<TickerDetailData>(
    effectiveTicker ? `/api/bandar/detail/${effectiveTicker}${qs}` : null,
    fetcher,
    { refreshInterval: 60000, dedupingInterval: 30000 }
  );

  useEffect(() => {
    if (!effectiveTicker || isLoading || !data) return;

    const timer = setTimeout(() => {
      const qsBase = `?window_days=${windowDays}${analysisDate ? `&analysis_date=${analysisDate}` : ""}`;
      const urlsToPreload = [
        `/api/bandar/broker-flow/${effectiveTicker}${qsBase}&flow_mode=Cumulative&dist_mode=Date%20range`,
        `/api/bandar/causality/${effectiveTicker}${qsBase}`,
        `/api/bandar/validation-v2/${effectiveTicker}${qsBase}&horizon=${horizon}&min_events=${minEvents}&universe_mode=${universe}`,
        `/api/bandar/stocks/${effectiveTicker}/raw-tables${qsBase}`,
      ];

      urlsToPreload.forEach((url) => {
        const opts: RequestInit & { priority?: string } = { credentials: "include", priority: "low" };
        fetch(url, opts).catch(() => {});
      });
    }, 400);

    return () => clearTimeout(timer);
  }, [effectiveTicker, isLoading, data, windowDays, analysisDate, horizon, minEvents, universe]);

  if (!effectiveTicker) {
    return (
      <div className="min-h-screen bg-[var(--md-sys-color-surface)] text-[var(--md-sys-color-on-surface)] flex items-center justify-center">
        Pilih emiten untuk melihat analitik kuantitatif
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 py-6 pb-24 transition-colors duration-200">
      <TickerSearch />

      <TickerHeader
        activeTicker={effectiveTicker}
        windowDays={windowDays}
        isBookmarked={isBookmarked}
        onToggleBookmark={toggleBookmark}
        data={data}
        isLoading={isLoading}
        error={error}
      />

      <div className="mb-6 border-b border-[var(--md-sys-color-outline-variant)]">
        <div className="flex gap-1.5 overflow-x-auto scrollbar-none pb-2">
          {TABS.map((tab) => {
            const isActive = activeTab === tab.id;
            return (
              <button
                key={tab.id}
                type="button"
                onClick={() => handleTabChange(tab.id)}
                className={`flex items-center gap-2 px-4 py-2 text-xs font-bold rounded-lg transition-all cursor-pointer shadow-2xs whitespace-nowrap ${
                  isActive
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-sm"
                    : "bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] hover:bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]"
                }`}
              >
                <Icon icon={tab.icon} width="16" />
                <span>{tab.label}</span>
              </button>
            );
          })}
        </div>
      </div>

      <div className="transition-colors duration-200">
        <div className={activeTab === "Overview" ? "block" : "hidden"}>
          <OverviewTab data={data} isLoading={isLoading} />
        </div>
        {visitedTabs["Broker Flow"] && (
          <div className={activeTab === "Broker Flow" ? "block" : "hidden"}>
            <Suspense fallback={<TabSkeleton />}>
              <BrokerFlowTab ticker={effectiveTicker} analysisDate={analysisDate} windowDays={windowDays} />
            </Suspense>
          </div>
        )}
        {visitedTabs.Causality && (
          <div className={activeTab === "Causality" ? "block" : "hidden"}>
            <Suspense fallback={<TabSkeleton />}>
              <CausalityTab
                ticker={effectiveTicker}
                analysisDate={analysisDate}
                windowDays={windowDays}
                initialDetail={data}
              />
            </Suspense>
          </div>
        )}
        {visitedTabs.Validation && (
          <div className={activeTab === "Validation" ? "block" : "hidden"}>
            <Suspense fallback={<TabSkeleton />}>
              <ValidationTab
                ticker={effectiveTicker}
                analysisDate={analysisDate}
                windowDays={windowDays}
                universeMode={universe}
                horizon={horizon}
                minEvents={minEvents}
              />
            </Suspense>
          </div>
        )}
        {visitedTabs["Raw Tables"] && (
          <div className={activeTab === "Raw Tables" ? "block" : "hidden"}>
            <Suspense fallback={<TabSkeleton />}>
              <RawTablesTab ticker={effectiveTicker} analysisDate={analysisDate} windowDays={windowDays} />
            </Suspense>
          </div>
        )}
      </div>
    </div>
  );
}
