import React from "react";
import { Icon } from "@iconify/react";
import { getLogLineClass, type StatusResponse } from "./marketReportTypes";

interface MarketReportLogcatModalProps {
  open: boolean;
  status: StatusResponse | null;
  logs: string[];
  autoScrollLogs: boolean;
  logLoading: boolean;
  logContainerRef: React.RefObject<HTMLDivElement | null>;
  onClose: () => void;
  onToggleAutoScroll: () => void;
  onRefreshLogs: () => void;
  onCopyLogs: () => void;
}

export default function MarketReportLogcatModal({
  open,
  status,
  logs,
  autoScrollLogs,
  logLoading,
  logContainerRef,
  onClose,
  onToggleAutoScroll,
  onRefreshLogs,
  onCopyLogs,
}: MarketReportLogcatModalProps) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-6 bg-black/70 backdrop-blur-sm">
      <div className="w-full max-w-4xl max-h-[88vh] rounded-3xl bg-zinc-950 border border-zinc-800 shadow-2xl flex flex-col overflow-hidden text-zinc-200">
        <div className="p-4 sm:px-6 py-3.5 border-b border-zinc-800 bg-zinc-900/80 flex items-center justify-between gap-2 flex-wrap">
          <div className="flex items-center gap-2.5">
            <div className="flex items-center gap-1.5">
              <div className="w-3 h-3 rounded-full bg-rose-500/80" />
              <div className="w-3 h-3 rounded-full bg-amber-500/80" />
              <div className="w-3 h-3 rounded-full bg-emerald-500/80" />
            </div>
            <div className="h-4 w-px bg-zinc-700 mx-1" />
            <span className="font-mono text-xs font-bold text-zinc-100 flex items-center gap-1.5">
              <Icon icon="ph:terminal-bold" width="14" className="text-emerald-400" />
              Terminal Logcat Eksekusi (Investowl ID Engine)
            </span>
            {status?.is_running ? (
              <span className="px-2 py-0.5 rounded-full bg-emerald-500/20 text-emerald-400 text-[10px] font-mono animate-pulse border border-emerald-500/30">
                LIVE STREAMING
              </span>
            ) : null}
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={onToggleAutoScroll}
              className={`px-2.5 py-1 rounded-lg text-[10px] font-mono font-bold transition flex items-center gap-1 cursor-pointer border ${
                autoScrollLogs
                  ? "bg-emerald-500/20 text-emerald-300 border-emerald-500/40"
                  : "bg-zinc-800 text-zinc-400 border-zinc-700"
              }`}
            >
              <Icon icon="ph:arrows-down-up-bold" width="12" />
              <span>Auto-Scroll: {autoScrollLogs ? "ON" : "OFF"}</span>
            </button>

            <button
              type="button"
              onClick={onRefreshLogs}
              disabled={logLoading}
              className="p-1.5 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-300 transition cursor-pointer"
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
              className="px-2.5 py-1 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-[10px] font-mono font-bold transition flex items-center gap-1 cursor-pointer border border-zinc-700"
            >
              <Icon icon="ph:copy-bold" width="12" />
              <span>Salin Log</span>
            </button>

            <button
              type="button"
              onClick={onClose}
              className="p-1.5 rounded-lg hover:bg-zinc-800 text-zinc-400 hover:text-white transition cursor-pointer"
              title="Tutup"
            >
              <Icon icon="ph:x-bold" width="16" />
            </button>
          </div>
        </div>

        <div
          ref={logContainerRef}
          className="p-4 sm:p-5 flex-1 overflow-y-auto font-mono text-xs leading-relaxed space-y-1 bg-zinc-950 selection:bg-emerald-500/30 min-h-[350px] max-h-[65vh]"
        >
          {logs.length === 0 ? (
            <div className="text-zinc-500 italic py-12 text-center flex flex-col items-center gap-2">
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

        <div className="p-3 px-5 border-t border-zinc-800 bg-zinc-900/80 flex items-center justify-between text-[11px] text-zinc-400 font-mono">
          <span>Total: {logs.length} baris log tercatat</span>
          <span className="flex items-center gap-1.5">
            <span
              className={`w-2 h-2 rounded-full ${
                status?.is_running ? "bg-emerald-400 animate-ping" : "bg-zinc-500"
              }`}
            />
            Status Engine: {status?.status?.toUpperCase() || (status?.is_running ? "RUNNING" : "IDLE")}
          </span>
        </div>
      </div>
    </div>
  );
}
