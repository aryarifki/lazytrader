import { useEffect, useMemo, useState } from "react";
import useSWR from "swr";
import { Icon } from "@iconify/react";
import { useAppStore } from "@/store/useAppStore";
import ExodusChart from "./ExodusChart";
import ExodusDetailsGrid from "./ExodusDetailsGrid";
import type { MarketState, RadarData } from "./types";

const fetcher = (url: string) =>
  fetch(url, { credentials: "include" }).then((res) => {
    if (!res.ok) throw new Error("Gagal mengambil data radar");
    return res.json();
  });

const STATE_CONFIG: Record<
  MarketState,
  {
    label: string;
    sublabel: string;
    icon: string;
    badgeClass: string;
    dotClass: string;
    glowClass: string;
  }
> = {
  CRASH: {
    label: "CRASH",
    sublabel: "Forced Liquidation Detected",
    icon: "ph:warning-octagon-fill",
    badgeClass: "bg-rose-500/15 text-rose-600 dark:text-rose-400 border-rose-500/30",
    dotClass: "bg-rose-500",
    glowClass: "shadow-[0_0_12px_rgba(244,63,94,0.35)]",
  },
  WARNING: {
    label: "WARNING",
    sublabel: "Hidden Foreign Distribution",
    icon: "ph:warning-circle-fill",
    badgeClass: "bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/30",
    dotClass: "bg-amber-500",
    glowClass: "shadow-[0_0_12px_rgba(245,158,11,0.35)]",
  },
  BOTTOMING: {
    label: "BOTTOMING",
    sublabel: "Capitulation Buyback (Initial Serap)",
    icon: "ph:arrow-circle-up-fill",
    badgeClass: "bg-sky-500/15 text-sky-600 dark:text-sky-400 border-sky-500/30",
    dotClass: "bg-sky-500",
    glowClass: "shadow-[0_0_12px_rgba(14,165,233,0.35)]",
  },
  RECOVERY_TESTING: {
    label: "TESTING RECOVERY",
    sublabel: "Pengujian Bottoming (5 Hari Bursa)",
    icon: "ph:hourglass-medium-fill",
    badgeClass: "bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/30",
    dotClass: "bg-amber-500 animate-pulse",
    glowClass: "shadow-[0_0_12px_rgba(245,158,11,0.35)]",
  },
  CONFIRMED_BOTTOM: {
    label: "CONFIRMED BOTTOM",
    sublabel: "Bottom Terkonfirmasi - Safe to DCA",
    icon: "ph:check-circle-fill",
    badgeClass: "bg-emerald-600/15 text-emerald-700 dark:text-emerald-400 border-emerald-600/30",
    dotClass: "bg-emerald-600",
    glowClass: "shadow-[0_0_12px_rgba(5,150,105,0.35)]",
  },
  NORMAL: {
    label: "NORMAL",
    sublabel: "Stable Foreign Flow",
    icon: "ph:shield-check-fill",
    badgeClass: "bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/30",
    dotClass: "bg-emerald-500",
    glowClass: "shadow-[0_0_12px_rgba(16,185,129,0.35)]",
  },
};

export default function ForeignExodusRadar() {
  const [isExpanded, setIsExpanded] = useState<boolean>(false);

  useEffect(() => {
    try {
      const saved = localStorage.getItem("investowl_radar_expanded");
      if (saved === "true") setIsExpanded(true);
    } catch {}
  }, []);

  const toggleExpanded = () => {
    setIsExpanded((prev) => {
      const next = !prev;
      try {
        localStorage.setItem("investowl_radar_expanded", String(next));
      } catch {}
      return next;
    });
  };

  const { data, isLoading, mutate: mutateRadar } = useSWR<RadarData>("/api/forensic/foreign-exodus", fetcher, {
    revalidateOnFocus: false,
    dedupingInterval: 30000,
  });

  const lastSyncTimestamp = useAppStore((state) => state.lastSyncTimestamp);

  useEffect(() => {
    if (lastSyncTimestamp > 0) {
      mutateRadar();
    }
  }, [lastSyncTimestamp, mutateRadar]);

  useEffect(() => {
    const handleSyncEvent = () => {
      mutateRadar();
    };
    window.addEventListener("investowl:sync", handleSyncEvent);
    return () => window.removeEventListener("investowl:sync", handleSyncEvent);
  }, [mutateRadar]);

  const stateConfig = useMemo(() => {
    const s = data?.market_state || "NORMAL";
    return STATE_CONFIG[s] || STATE_CONFIG.NORMAL;
  }, [data?.market_state]);

  const timeline = useMemo(() => data?.timeline || [], [data?.timeline]);
  const culpritBrokers = useMemo(() => data?.culprit_brokers || [], [data?.culprit_brokers]);

  const isNetOutflow = useMemo(() => {
    if (!timeline.length) return true;
    return timeline[timeline.length - 1].cumulative_flow < 0;
  }, [timeline]);

  const gradientOffset = useMemo(() => {
    if (!timeline.length) return 0;
    const flows = timeline.map((d) => d.cumulative_flow);
    const max = Math.max(...flows);
    const min = Math.min(...flows);
    if (max <= 0) return 0;
    if (min >= 0) return 1;
    return max / (max - min);
  }, [timeline]);

  const syncDateFormatted = useMemo(() => {
    if (!timeline.length) return "";
    const lastDate = timeline[timeline.length - 1]?.date;
    if (!lastDate) return "";
    try {
      const parts = lastDate.split("-");
      if (parts.length === 3) {
        const months = [
          "Januari",
          "Februari",
          "Maret",
          "April",
          "Mei",
          "Juni",
          "Juli",
          "Agustus",
          "September",
          "Oktober",
          "November",
          "Desember",
        ];
        return `${parseInt(parts[2], 10)} ${months[parseInt(parts[1], 10) - 1]} ${parts[0]}`;
      }
      return lastDate;
    } catch {
      return lastDate;
    }
  }, [timeline]);

  return (
    <section className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/70 rounded-3xl p-5 sm:p-6 shadow-xs transition-all relative overflow-hidden">
      <div className="flex flex-col gap-3.5">
        <div className="flex items-start sm:items-center gap-3.5">
          <div className="w-12 h-12 rounded-2xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0 border border-[var(--md-sys-color-outline-variant)] shadow-2xs mt-0.5 sm:mt-0">
            <svg xmlns="http://www.w3.org/2000/svg" width="26" height="26" viewBox="0 0 24 24" className="animate-spin-slow text-[var(--md-sys-color-primary)]">
              <g fill="currentColor">
                <path d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z" opacity=".35" />
                <path d="M11.1256 6.82338C11.846 6.70116 12.6041 6.72817 13.3586 6.93033C16.1585 7.68056 17.8201 10.5585 17.0699 13.3584C16.3196 16.1583 13.4417 17.8199 10.6418 17.0697C7.84187 16.3194 6.18028 13.4415 6.93051 10.6416C7.05879 10.1629 7.24878 9.7184 7.48934 9.31389L11.5283 12.583C11.8503 12.8436 12.3226 12.7938 12.5832 12.4718C12.8438 12.1499 12.794 11.6776 12.472 11.417L8.27423 8.01941C7.7127 7.56491 6.80868 7.60723 6.3558 8.29795C5.97114 8.8846 5.67251 9.54098 5.48162 10.2534C4.51698 13.8535 6.65344 17.5539 10.2535 18.5186C13.8536 19.4832 17.5541 17.3467 18.5187 13.7466C19.4834 10.1465 17.3469 6.44608 13.7468 5.48144C12.7798 5.22232 11.8035 5.18693 10.8747 5.34452C10.4664 5.41381 10.1915 5.80103 10.2608 6.20941C10.33 6.61779 10.7173 6.89267 11.1256 6.82338Z" />
              </g>
            </svg>
          </div>
          <div>
            <div className="inline-flex items-center gap-2 text-[10px] font-mono font-bold text-[var(--md-sys-color-primary)] uppercase tracking-wider mb-0.5">
              <span>FOREIGN EXODUS RADAR</span>
              <span className="opacity-40">•</span>
              <span>IHSG COMPOSITE</span>
            </div>

            <div className="flex items-baseline gap-2.5">
              <h2 className="text-2xl sm:text-3xl font-extrabold font-mono text-[var(--md-sys-color-on-surface)] tracking-tight">
                {isLoading ? (
                  <span className="inline-block w-28 h-7 bg-[var(--md-sys-color-surface-container-highest)] animate-pulse rounded-md" />
                ) : (
                  (data?.ihsg_current || 0).toLocaleString("id-ID", {
                    minimumFractionDigits: 2,
                    maximumFractionDigits: 2,
                  })
                )}
              </h2>

              {!isLoading && data && (
                <span
                  className={`inline-flex items-center gap-0.5 text-xs sm:text-sm font-mono font-bold ${
                    data.ihsg_change_pct >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                  }`}
                >
                  <Icon icon={data.ihsg_change_pct >= 0 ? "ph:trend-up-bold" : "ph:trend-down-bold"} width="15" />
                  <span>
                    {data.ihsg_change_pct > 0 ? "+" : ""}
                    {data.ihsg_change_pct.toFixed(2)}%
                  </span>
                </span>
              )}
            </div>

            {isLoading ? (
              <div className="w-28 h-3.5 bg-[var(--md-sys-color-surface-container-highest)] animate-pulse rounded mt-1.5" />
            ) : syncDateFormatted ? (
              <div className="text-[10px] sm:text-[11px] font-mono font-medium text-[var(--md-sys-color-on-surface-variant)] mt-1 flex items-center gap-1.5">
                <span className="inline-block w-1.5 h-1.5 rounded-full bg-emerald-500/80 animate-pulse shrink-0" />
                <span>sync : {syncDateFormatted}</span>
              </div>
            ) : null}
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-2.5 pt-0.5">
          {isLoading ? (
            <div className="w-32 h-8 bg-[var(--md-sys-color-surface-container-highest)] animate-pulse rounded-xl" />
          ) : (
            <div
              className={`inline-flex items-center gap-2 px-3 py-1.5 rounded-xl border text-xs font-mono font-extrabold tracking-wider transition-all ${stateConfig.badgeClass} ${stateConfig.glowClass}`}
            >
              <span className="relative flex h-2 w-2">
                <span
                  className={`animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 ${stateConfig.dotClass}`}
                ></span>
                <span className={`relative inline-flex rounded-full h-2 w-2 ${stateConfig.dotClass}`}></span>
              </span>
              <Icon icon={stateConfig.icon} width="16" />
              <span>{stateConfig.label}</span>
            </div>
          )}

          {!isLoading && (
            <span className="text-[11px] font-medium text-[var(--md-sys-color-on-surface-variant)]">
              {stateConfig.sublabel}
            </span>
          )}
        </div>
      </div>

      {isExpanded && (
        <div className="mt-6 pt-6 border-t border-[var(--md-sys-color-outline-variant)]/60 space-y-6 animate-fade-in">
          <ExodusChart
            timeline={timeline}
            isLoading={isLoading}
            isNetOutflow={isNetOutflow}
            gradientOffset={gradientOffset}
          />
          <ExodusDetailsGrid culpritBrokers={culpritBrokers} timeline={timeline} />
        </div>
      )}

      <div className="mt-4 pt-3 border-t border-[var(--md-sys-color-outline-variant)]/40 flex justify-center">
        <button
          type="button"
          onClick={toggleExpanded}
          className="inline-flex items-center justify-center gap-2 text-xs font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-primary)] transition-all cursor-pointer select-none py-1.5 px-4 rounded-full hover:bg-[var(--md-sys-color-surface-container-high)]/60 active:scale-95 group"
          aria-expanded={isExpanded}
          title={isExpanded ? "Tutup detail radar" : "Buka grafik dan detail arus asing"}
        >
          <span>{isExpanded ? "Tutup Detail Radar" : "Detail Radar"}</span>
          <Icon
            icon={isExpanded ? "ph:caret-up-bold" : "ph:caret-down-bold"}
            width="15"
            className="transition-transform duration-200 text-[var(--md-sys-color-primary)] group-hover:translate-y-0.5"
          />
        </button>
      </div>
    </section>
  );
}
