import { useCallback, useEffect, useState } from "react";
import { Icon } from "@iconify/react";
import { getKongloReport } from "@/lib/api/terminal";
import { useAppStore } from "@/store/useAppStore";
import KongloFundamentalTab from "./KongloFundamentalTab";
import KongloGroupTab from "./KongloGroupTab";
import KongloOverviewCards from "./KongloOverviewCards";
import KongloRankingsTab, { type KongloRankingTabId } from "./KongloRankingsTab";
import KongloSummaryTab from "./KongloSummaryTab";
import {
  KONGLO_TABS,
  type KongloReportApiResponse,
  type KongloReportData,
  type KongloTab,
} from "./kongloTypes";

function isRankingTab(tab: KongloTab): tab is KongloRankingTabId {
  return (
    tab === "pct-up" ||
    tab === "mc-up" ||
    tab === "pct-down" ||
    tab === "mc-down"
  );
}

export default function KongloPage() {
  const theme = useAppStore((s) => s.theme);
  const toggleTheme = useAppStore((s) => s.toggleTheme);

  const [activeTab, setActiveTab] = useState<KongloTab>("summary");
  const [reportData, setReportData] = useState<KongloReportData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [refreshing, setRefreshing] = useState<boolean>(false);
  const [errorMsg, setErrorMsg] = useState<string>("");
  const [selectedGroup, setSelectedGroup] = useState<string>("Barito Group");

  const loadData = useCallback((isManual = false) => {
    if (isManual) setRefreshing(true);
    else setLoading(true);
    setErrorMsg("");

    getKongloReport<KongloReportApiResponse>()
      .then((res) => {
        if (res.report) {
          setReportData(res.report);
        }
        setLoading(false);
        setRefreshing(false);
      })
      .catch((err: unknown) => {
        console.error("Error loading Konglo Report:", err);
        setErrorMsg(
          err instanceof Error ? err.message : "Gagal memuat laporan konglomerasi"
        );
        setLoading(false);
        setRefreshing(false);
      });
  }, []);

  useEffect(() => {
    loadData(false);
  }, [loadData]);

  const meta = reportData?.metadata;
  const dateStr = meta?.date_display || "Senin, 02 Nov 2026";
  const timeStr = meta?.time_display || "16:00";
  const period = `${dateStr} — ${timeStr}`;

  return (
    <main className="w-full max-w-7xl mx-auto p-4 md:p-6 space-y-4 animate-fadeIn pb-16">
      {/* ── HEADER ── */}
      <div className="app-header flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-[var(--md-sys-color-outline-variant)] pb-3">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] flex items-center justify-center shrink-0 shadow-2xs">
            <Icon icon="ph:buildings-duotone" width="22" />
          </div>
          <div>
            <div className="inline-flex items-center gap-2 text-[9px] font-mono font-bold text-[var(--md-sys-color-primary)] uppercase tracking-wider mb-0.5">
              <span>DAILY KONGLO REPORT</span>
              <span className="opacity-40">•</span>
              <span>INVESTOWL TERMINAL</span>
            </div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight text-[var(--md-sys-color-on-surface)]">
              Konglo Screener
            </h1>
          </div>
        </div>

        {/* Header Right: Market Summary Chips + Theme Toggle */}
        <div className="flex items-center gap-2 text-xs font-mono flex-wrap">
          <span className="px-2.5 py-1 rounded-lg bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/25 font-bold">
            ▲ {reportData?.summary?.stocks_up ?? 358}
          </span>
          <span className="px-2.5 py-1 rounded-lg bg-rose-500/10 text-rose-700 dark:text-rose-400 border border-rose-500/25 font-bold">
            ▼ {reportData?.summary?.stocks_down ?? 628}
          </span>
          <span className="px-2.5 py-1 rounded-lg bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)] font-bold">
            ■ {reportData?.summary?.stocks_flat ?? 356}
          </span>
          <button
            id="themeToggle"
            type="button"
            onClick={toggleTheme}
            className="p-1.5 rounded-lg border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] cursor-pointer transition-colors"
            title={`Beralih ke Mode ${theme === "dark" ? "Terang" : "Gelap"}`}
          >
            <Icon
              icon={theme === "dark" ? "ph:sun-duotone" : "ph:moon-stars-duotone"}
              width="16"
            />
          </button>
        </div>
      </div>

      {/* ── SUB-NAV TAB BAR (7 TABS) ── */}
      <div className="flex items-center bg-[var(--md-sys-color-surface-container)] p-1 rounded-xl border border-[var(--md-sys-color-outline-variant)] shadow-2xs overflow-x-auto scrollbar-none gap-1">
        {KONGLO_TABS.map((tab) => {
          const active = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              id={`tab-btn-${tab.id}`}
              type="button"
              onClick={() => {
                setActiveTab(tab.id);
                if (tab.id === "konglo" && !selectedGroup) {
                  setSelectedGroup("Barito Group");
                }
              }}
              className={`px-3 py-1.5 rounded-lg text-xs font-bold transition flex items-center gap-1.5 cursor-pointer whitespace-nowrap ${
                active
                  ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] shadow-2xs font-extrabold"
                  : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
              }`}
            >
              <Icon icon={tab.icon} width="15" />
              <span>{tab.label}</span>
            </button>
          );
        })}
      </div>

      {/* ── STATUS, CONTROLS, AND OVERVIEW CARDS ── */}
      <KongloOverviewCards
        reportData={reportData}
        refreshing={refreshing}
        dateStr={dateStr}
        timeStr={timeStr}
        period={period}
        onRefresh={() => loadData(true)}
      />

      {/* ── RESULTS WRAPPER ── */}
      <div id="results" className="space-y-4 pt-1">
        {loading ? (
          <div className="p-16 text-center text-xs text-[var(--md-sys-color-on-surface-variant)] flex flex-col items-center gap-3">
            <div className="w-7 h-7 border-2 border-[var(--md-sys-color-primary)] border-t-transparent rounded-full animate-spin" />
            <span>Memuat data konglomerasi...</span>
          </div>
        ) : errorMsg && !reportData ? (
          <div className="rounded-2xl border border-rose-500/30 bg-rose-500/10 p-6 text-center space-y-2">
            <Icon
              icon="ph:warning-circle-duotone"
              width="24"
              className="text-rose-600 dark:text-rose-400 mx-auto"
            />
            <p className="text-xs font-mono text-rose-700 dark:text-rose-300">
              {errorMsg}
            </p>
            <button
              type="button"
              onClick={() => loadData(true)}
              className="px-3 py-1.5 rounded-lg text-xs font-bold bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] cursor-pointer"
            >
              Coba Lagi
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            {activeTab === "summary" && (
              <KongloSummaryTab reportData={reportData} period={period} />
            )}

            {isRankingTab(activeTab) && (
              <KongloRankingsTab
                activeTab={activeTab}
                reportData={reportData}
                period={period}
              />
            )}

            {activeTab === "konglo" && (
              <KongloGroupTab
                reportData={reportData}
                period={period}
                selectedGroup={selectedGroup}
                onSelectGroup={setSelectedGroup}
              />
            )}

            {activeTab === "fundamental" && (
              <KongloFundamentalTab reportData={reportData} period={period} />
            )}

            {/* ── Konglo Missing Banner if any ── */}
            {reportData?.missing_konglo &&
              reportData.missing_konglo.length > 0 &&
              activeTab !== "fundamental" && (
                <div
                  id="missingBannerKonglo"
                  className="p-3.5 rounded-xl bg-amber-500/10 border border-amber-500/30 text-xs space-y-1"
                >
                  <div className="font-bold text-amber-700 dark:text-amber-400 flex items-center gap-1.5 font-mono">
                    <span>⚠️</span>
                    <span>
                      {reportData.missing_konglo.length} saham konglomerasi tidak ada
                      datanya
                    </span>
                  </div>
                  <div className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
                    {reportData.missing_konglo
                      .map((m) => `${m.ticker} (${m.group})`)
                      .join(" · ")}
                  </div>
                </div>
              )}
          </div>
        )}
      </div>
    </main>
  );
}
