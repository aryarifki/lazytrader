import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Icon } from "@iconify/react";

import ResearchTabs from "@/components/research/ResearchTabs";
import CTExplorer from "./ct/CTExplorer";
import ExplorerDashboard from "./snapshot/ExplorerDashboard";

const TABS = [
  {
    id: "snapshot",
    label: "Snapshot Screener",
    description: "Screener SQL snapshot kepemilikan KSEI & UBO",
    icon: "ph:funnel-duotone",
  },
  {
    id: "parity",
    label: "Parity Explorer",
    description: "Baseline CT, delta bulanan & T5",
    icon: "ph:chart-polar-duotone",
  },
] as const;

type TabId = (typeof TABS)[number]["id"];

const STORAGE_KEY = "investowl_explorer_active_tab";

export default function ExplorerPageSwitch() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const sourceParam = searchParams.get("source");
  const initialTab: TabId = sourceParam === "parity" ? "parity" : "snapshot";
  const [activeTab, setActiveTab] = useState<TabId>(initialTab);

  useEffect(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY) as TabId | null;
      if (sourceParam === "snapshot" || sourceParam === "parity") {
        setActiveTab(sourceParam);
        localStorage.setItem(STORAGE_KEY, sourceParam);
      } else if (saved === "snapshot" || saved === "parity") {
        setActiveTab(saved);
        const params = new URLSearchParams(window.location.search);
        params.set("source", saved);
        navigate(`/explorer?${params.toString()}`, {
          replace: true,
          preventScrollReset: true,
        });
      } else {
        setActiveTab("snapshot");
      }
    } catch {
      // Ignore storage errors
    }
  }, [sourceParam, navigate]);

  const handleTabChange = (id: TabId) => {
    setActiveTab(id);
    try {
      localStorage.setItem(STORAGE_KEY, id);
    } catch {
      // Ignore storage errors
    }
    const params = new URLSearchParams(window.location.search);
    params.set("source", id);
    navigate(`/explorer?${params.toString()}`, {
      replace: true,
      preventScrollReset: true,
    });
  };

  return (
    <div className="max-w-7xl mx-auto space-y-4">
      {/* 1. UNIFIED HEADER */}
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-3">
        <div className="flex items-center gap-3.5">
          <div className="w-11 h-11 sm:w-12 sm:h-12 rounded-2xl bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0 border border-[var(--md-sys-color-outline-variant)] shadow-2xs">
            <Icon icon="ph:funnel-duotone" width="26" />
          </div>
          <div>
            <div className="inline-flex items-center gap-2 text-[9px] font-mono font-bold text-[var(--md-sys-color-primary)] uppercase tracking-wider mb-0.5">
              <span>OWNERSHIP EXPLORER</span>
              <span className="opacity-40">•</span>
              <span>{activeTab === "snapshot" ? "SNAPSHOT SCREENER" : "CT PARITY"}</span>
            </div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight">
              Explorer Kepemilikan
            </h1>
          </div>
        </div>
        <ResearchTabs />
      </div>

      {/* 2. PERSISTENT TAB BAR (NEVER DISAPPEARS) */}
      <div className="border-b border-[var(--md-sys-color-outline-variant)]">
        <div className="flex gap-1.5 overflow-x-auto scrollbar-none pb-2">
          {TABS.map((tab) => {
            const active = activeTab === tab.id;
            return (
              <button
                key={tab.id}
                type="button"
                onClick={() => handleTabChange(tab.id)}
                className={`flex items-center gap-2 rounded-lg border px-3.5 py-2 text-xs font-bold transition-all whitespace-nowrap cursor-pointer ${
                  active
                    ? "border-[var(--md-sys-color-primary)] bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-sm"
                    : "border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface-variant)] hover:bg-[var(--md-sys-color-surface-container-high)] hover:text-[var(--md-sys-color-on-surface)]"
                }`}
                title={tab.description}
              >
                <Icon icon={tab.icon} width="16" />
                <span>{tab.label}</span>
              </button>
            );
          })}
        </div>
      </div>

      {/* 3. ACTIVE TAB CONTENT */}
      <div className="transition-colors duration-200">
        {activeTab === "snapshot" ? (
          <ExplorerDashboard embedded />
        ) : (
          <CTExplorer embedded />
        )}
      </div>
    </div>
  );
}
