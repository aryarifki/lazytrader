import React, { useEffect, useMemo, useRef, useState } from "react";
import KongloGroupDetail from "./KongloGroupDetail";
import { KongloGroupBadge, KongloTickerLink } from "./KongloShared";
import { fmtM, fmtMc, getGroupColor, type KongloReportData } from "./kongloTypes";

interface KongloGroupTabProps {
  reportData: KongloReportData | null;
  period: string;
  selectedGroup: string;
  onSelectGroup: (group: string) => void;
}

export default function KongloGroupTab({
  reportData,
  period,
  selectedGroup,
  onSelectGroup,
}: KongloGroupTabProps) {
  const [groupSearchQuery, setGroupSearchQuery] = useState<string>(
    selectedGroup || "Barito Group"
  );
  const [groupDropdownOpen, setGroupDropdownOpen] = useState<boolean>(false);
  const dropdownContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setGroupSearchQuery(selectedGroup);
  }, [selectedGroup]);

  useEffect(() => {
    const handleOutsideClick = (event: MouseEvent) => {
      if (
        dropdownContainerRef.current &&
        !dropdownContainerRef.current.contains(event.target as Node)
      ) {
        setGroupDropdownOpen(false);
      }
    };
    document.addEventListener("mousedown", handleOutsideClick);
    return () => document.removeEventListener("mousedown", handleOutsideClick);
  }, []);

  const filteredGroups = useMemo(() => {
    if (!reportData?.all_group_stats) return [];
    const q = groupSearchQuery.toLowerCase();
    return reportData.all_group_stats.filter((g) =>
      g.group.toLowerCase().includes(q)
    );
  }, [reportData, groupSearchQuery]);

  const curGroupObj = useMemo(() => {
    if (!reportData?.all_group_stats) return null;
    return (
      reportData.all_group_stats.find((g) => g.group === selectedGroup) ||
      reportData.all_group_stats[0] ||
      null
    );
  }, [reportData, selectedGroup]);

  const selectGroup = (group: string) => {
    onSelectGroup(group);
    setGroupSearchQuery(group);
    setGroupDropdownOpen(false);
  };

  return (
    <div id="tab-konglo" className="space-y-4 animate-fadeIn">
      <div className="flex items-center gap-2 text-xs font-bold text-[var(--md-sys-color-on-surface)]">
        <span className="w-2.5 h-2.5 rounded-full bg-[var(--md-sys-color-primary)]" />
        <span className="text-sm">Konglo Report</span>
        <span className="text-[var(--md-sys-color-on-surface-variant)] text-xs font-normal font-mono">
          {period}
        </span>
      </div>

      {/* 5 Konglomerasi Terbaik Table */}
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl overflow-hidden shadow-2xs">
        <div className="px-4 py-3 border-b border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-high)]/60 flex items-center justify-between gap-2">
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-emerald-500" />
            <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)]">
              5 Konglomerasi Terbaik
            </span>
          </div>
          <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
            diurutkan avg price change
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead>
              <tr className="border-b border-[var(--md-sys-color-outline-variant)]/40 bg-[var(--md-sys-color-surface-container-high)]/30 text-[var(--md-sys-color-on-surface-variant)] font-mono uppercase text-[10px]">
                <th className="py-2.5 px-3 w-8 font-bold">#</th>
                <th className="py-2.5 px-3 font-bold">Konglomerasi</th>
                <th className="py-2.5 px-3 font-bold">Saham Terbaik</th>
                <th className="py-2.5 px-3 text-right font-bold">Value Total</th>
                <th className="py-2.5 px-3 text-right font-bold">Kenaikan MC</th>
                <th className="py-2.5 px-3 text-right font-bold">Avg Change</th>
              </tr>
            </thead>
            <tbody>
              {reportData?.top5_konglo?.map((g, idx) => {
                const col = getGroupColor(g.group);
                return (
                  <tr
                    key={g.group}
                    className="border-b border-[var(--md-sys-color-outline-variant)]/20 hover:bg-[var(--md-sys-color-surface-container-high)]/40 transition-colors"
                  >
                    <td className="py-2.5 px-3">
                      <span className="inline-flex rounded-full bg-[#1E1A17]/90 dark:bg-transparent">
                        <div
                          className="w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold font-mono"
                          style={{
                            backgroundColor: col.bg,
                            color: col.color,
                            border: `1px solid ${col.border}`,
                          }}
                        >
                          {idx + 1}
                        </div>
                      </span>
                    </td>
                    <td className="py-2.5 px-3">
                      <KongloGroupBadge
                        group={g.group}
                        onClick={() => selectGroup(g.group)}
                        title={`Klik untuk lihat detail ${g.group}`}
                      />
                    </td>
                    <td className="py-2.5 px-3">
                      <div className="flex items-center gap-2 flex-wrap">
                        {g.top2?.map((t2, tIdx) => (
                          <React.Fragment key={t2.ticker}>
                            {tIdx > 0 && <span className="opacity-40">·</span>}
                            <span className="inline-flex items-center gap-1 font-bold text-xs font-mono">
                              <KongloTickerLink ticker={t2.ticker} />
                              <span
                                className={
                                  t2.pct >= 0
                                    ? "text-emerald-600 dark:text-emerald-400"
                                    : "text-rose-600 dark:text-rose-400"
                                }
                              >
                                {t2.pct >= 0
                                  ? `+${t2.pct.toFixed(1)}%`
                                  : `${t2.pct.toFixed(1)}%`}
                              </span>
                            </span>
                          </React.Fragment>
                        ))}
                      </div>
                    </td>
                    <td className="py-2.5 px-3 text-right text-[var(--md-sys-color-on-surface-variant)] font-mono">
                      {fmtM(g.total_val)}
                    </td>
                    <td
                      className={`py-2.5 px-3 text-right font-semibold font-mono ${
                        g.mc_chg >= 0
                          ? "text-emerald-600 dark:text-emerald-400"
                          : "text-rose-600 dark:text-rose-400"
                      }`}
                    >
                      {fmtMc(g.mc_chg)}
                    </td>
                    <td
                      className={`py-2.5 px-3 text-right font-bold font-mono ${
                        g.avg_pct >= 0
                          ? "text-emerald-600 dark:text-emerald-400"
                          : "text-rose-600 dark:text-rose-400"
                      }`}
                    >
                      {g.avg_pct >= 0
                        ? `+${g.avg_pct.toFixed(1)}%`
                        : `${g.avg_pct.toFixed(1)}%`}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>

      {/* Search Box & Dropdown */}
      <div className="relative max-w-sm" ref={dropdownContainerRef}>
        <input
          id="kongloSearch"
          type="text"
          value={groupSearchQuery}
          onChange={(e) => {
            setGroupSearchQuery(e.target.value);
            setGroupDropdownOpen(true);
          }}
          onFocus={() => setGroupDropdownOpen(true)}
          onKeyDown={(e) => {
            if (e.key === "Escape") setGroupDropdownOpen(false);
          }}
          placeholder="🔍 Ketik atau pilih konglomerasi..."
          className="w-full rounded-xl px-3.5 py-2 text-xs border border-[var(--md-sys-color-outline-variant)]/60 bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface)] focus:outline-none focus:border-[var(--md-sys-color-primary)] shadow-2xs"
        />
        {groupDropdownOpen && (
          <div
            id="kongloDropdown"
            className="absolute top-full left-0 mt-1 w-full border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface)] rounded-xl shadow-2xl z-50 max-h-60 overflow-y-auto"
          >
            {filteredGroups.map((g) => (
              <div
                key={g.group}
                className="px-3.5 py-2 text-xs hover:bg-[var(--md-sys-color-primary-container)] hover:text-[var(--md-sys-color-on-primary-container)] cursor-pointer flex justify-between transition-colors"
                onClick={() => selectGroup(g.group)}
              >
                <span className="font-medium">{g.group}</span>
                <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
                  {g.total_members} saham
                </span>
              </div>
            ))}
          </div>
        )}
      </div>

      <KongloGroupDetail group={curGroupObj} />
    </div>
  );
}
