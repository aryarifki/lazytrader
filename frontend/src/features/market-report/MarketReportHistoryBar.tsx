import React from "react";
import { Icon } from "@iconify/react";
import { formatFileSize, type HistoryItem } from "./marketReportTypes";

interface MarketReportHistoryBarProps {
  history: HistoryItem[];
  selectedHistoryFile: string;
  onSelectHistory: (filename: string) => void;
}

export default function MarketReportHistoryBar({
  history,
  selectedHistoryFile,
  onSelectHistory,
}: MarketReportHistoryBarProps) {
  if (history.length === 0) return null;

  return (
    <div className="flex items-center gap-2 overflow-x-auto pb-1 text-xs">
      <span className="font-bold text-[11px] uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)] flex items-center gap-1 shrink-0">
        <Icon icon="ph:archive-duotone" width="14" />
        Riwayat Arsip:
      </span>
      <div className="flex items-center gap-1.5 flex-nowrap">
        {history.map((h) => {
          const isSelected = selectedHistoryFile === h.filename;
          return (
            <button
              key={h.filename}
              type="button"
              onClick={() => onSelectHistory(h.filename)}
              title={`${h.filename} (${formatFileSize(h.size_bytes)})`}
              className={`px-3 py-1.5 rounded-xl font-mono text-xs font-bold transition flex items-center gap-1.5 shrink-0 cursor-pointer border ${
                isSelected
                  ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] border-[var(--md-sys-color-primary)] shadow-xs"
                  : "bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface-variant)] hover:bg-[var(--md-sys-color-surface-container-high)] border-[var(--md-sys-color-outline-variant)]"
              }`}
            >
              <span
                className={`w-1.5 h-1.5 rounded-full ${
                  h.session === 1 ? "bg-emerald-400" : "bg-sky-400"
                }`}
              />
              <span>
                Sesi {h.session} ({h.date})
              </span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
