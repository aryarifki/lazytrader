import React from "react";
import { Icon } from "@iconify/react";
import { getLogLineClass, type StatusResponse } from "./marketReportTypes";

interface MarketReportLogsPanelProps {
  logs: string[];
  status: StatusResponse | null;
  autoScrollLogs: boolean;
  logLoading: boolean;
  collapsed: boolean;
  logRef: React.RefObject<HTMLDivElement | null>;
  onToggleAutoScroll: () => void;
  onRefreshLogs: () => void;
  onCopyLogs: () => void;
  onOpenFullscreen: () => void;
  onToggleCollapsed: () => void;
}

export default function MarketReportLogsPanel({
  logs,
  status,
  autoScrollLogs,
  logLoading,
  collapsed,
  logRef,
  onToggleAutoScroll,
  onRefreshLogs,
  onCopyLogs,
  onOpenFullscreen,
  onToggleCollapsed,
}: MarketReportLogsPanelProps) {
  return (
    <section className="rounded-3xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] overflow-hidden shadow-xs">
      {/* Header */}
      <div className="p-4 sm:px-6 py-3.5 border-b border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)]/60 flex items-center justify-between flex-wrap gap-2.5">
        <div className="flex items-center gap-2.5 flex-wrap">
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded-full bg-rose-500/80" />
            <div className="w-3 h-3 rounded-full bg-amber-500/80" />
            <div className="w-3 h-3 rounded-full bg-emerald-500/80" />
          </div>
          <div className="h-4 w-px bg-[var(--md-sys-color-outline-variant)] mx-1" />
          <span className="font-mono text-xs font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
            <Icon
              icon="ph:terminal-bold"
              width="16"
              className="text-emerald-600 dark:text-emerald-400"
            />
            Terminal Logcat Eksekusi (Investowl ID Engine)
          </span>

          {status?.is_running ? (
            <span className="px-2 py-0.5 rounded-full bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 text-[10px] font-mono animate-pulse border border-emerald-500/30 flex items-center gap-1">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
              LIVE STREAMING
            </span>
          ) : status?.error ? (
            <span className="px-2 py-0.5 rounded-full bg-rose-500/15 text-rose-700 dark:text-rose-400 text-[10px] font-mono border border-rose-500/30 flex items-center gap-1">
              <span className="w-1.5 h-1.5 rounded-full bg-rose-500" />
              ERROR
            </span>
          ) : (
            <span className="px-2 py-0.5 rounded-full bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] text-[10px] font-mono border border-[var(--md-sys-color-outline-variant)] flex items-center gap-1">
              <span className="w-1.5 h-1.5 rounded-full bg-[var(--md-sys-color-outline)]" />
              IDLE
            </span>
          )}
        </div>

        <div className="flex items-center gap-2 flex-wrap">
          <button
            type="button"
            onClick={onToggleAutoScroll}
            className={`px-2.5 py-1 rounded-lg text-[10px] font-mono font-bold transition flex items-center gap-1 cursor-pointer border ${
              autoScrollLogs
                ? "bg-emerald-500/15 text-emerald-700 dark:text-emerald-300 border-emerald-500/35"
                : "bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]"
            }`}
            title="Auto-scroll ke baris terbaru"
          >
            <Icon icon="ph:arrows-down-up-bold" width="12" />
            <span>Auto-Scroll: {autoScrollLogs ? "ON" : "OFF"}</span>
          </button>

          <button
            type="button"
            onClick={onRefreshLogs}
            disabled={logLoading}
            className="p-1.5 rounded-lg bg-[var(--md-sys-color-surface-container)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)] transition cursor-pointer"
            title="Refresh Logs"
          >
            <Icon
              icon="ph:arrows-clockwise-bold"
              width="14"
              className={logLoading ? "animate-spin" : ""}
            />
          </button>

          <button
            type="button"
            onClick={onCopyLogs}
            className="px-2.5 py-1 rounded-lg bg-[var(--md-sys-color-surface-container)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface)] text-[10px] font-mono font-bold transition flex items-center gap-1 cursor-pointer border border-[var(--md-sys-color-outline-variant)]"
            title="Salin Semua Log"
          >
            <Icon icon="ph:copy-bold" width="12" />
            <span>Salin Log</span>
          </button>

          <button
            type="button"
            onClick={onOpenFullscreen}
            className="px-2.5 py-1 rounded-lg bg-[var(--md-sys-color-surface-container)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-primary)] text-[10px] font-mono font-bold transition flex items-center gap-1 cursor-pointer border border-[var(--md-sys-color-outline-variant)]"
            title="Buka Layar Penuh (Pop-out)"
          >
            <Icon icon="ph:arrows-out-simple-bold" width="12" />
            <span>Layar Penuh</span>
          </button>

          <button
            type="button"
            onClick={onToggleCollapsed}
            className="p-1.5 rounded-lg hover:bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] transition cursor-pointer"
            title={collapsed ? "Bentangkan Konsol" : "Ciutkan Konsol"}
          >
            <Icon icon={collapsed ? "ph:caret-down-bold" : "ph:caret-up-bold"} width="14" />
          </button>
        </div>
      </div>

      {/* Console Body */}
      {!collapsed ? (
        <div
          ref={logRef}
          className="p-4 sm:p-5 font-mono text-xs leading-relaxed space-y-1 bg-zinc-950 text-zinc-200 selection:bg-emerald-500/30 min-h-[220px] max-h-[380px] overflow-y-auto"
        >
          {logs.length === 0 ? (
            <div className="text-zinc-500 italic py-10 text-center flex flex-col items-center gap-2">
              <Icon icon="ph:file-dashed-bold" width="28" />
              <span>Belum ada log eksekusi terbaru dari server.</span>
            </div>
          ) : (
            logs.map((line, idx) => (
              <div
                key={idx}
                className={`font-mono text-[11px] sm:text-xs break-all px-1.5 py-0.5 rounded ${getLogLineClass(line)}`}
              >
                {line}
              </div>
            ))
          )}
        </div>
      ) : null}

      {/* Footer */}
      <div className="p-2.5 px-5 border-t border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)]/40 flex items-center justify-between text-[11px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
        <span>Total: {logs.length} baris log tersimpan (disk &amp; memory)</span>
        <span>
          Status Engine: {status?.status?.toUpperCase() || (status?.is_running ? "RUNNING" : "IDLE")}
        </span>
      </div>
    </section>
  );
}
