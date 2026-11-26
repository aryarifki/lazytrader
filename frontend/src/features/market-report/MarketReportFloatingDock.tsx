import React from "react";
import { Icon } from "@iconify/react";
import { formatTimer, type StatusResponse, type ToastState } from "./marketReportTypes";

interface MarketReportFloatingDockProps {
  toast: ToastState;
  status: StatusResponse | null;
  floatingMinimized: boolean;
  elapsedSeconds: number;
  onSetFloatingMinimized: (minimized: boolean) => void;
  onOpenLogcat: () => void;
  onOpenTokenModal: () => void;
  onDismissError: () => void;
}

export default function MarketReportFloatingDock({
  toast,
  status,
  floatingMinimized,
  elapsedSeconds,
  onSetFloatingMinimized,
  onOpenLogcat,
  onOpenTokenModal,
  onDismissError,
}: MarketReportFloatingDockProps) {
  return (
    <>
      {/* ── TOAST NOTIFICATION ── */}
      {toast.visible ? (
        <div
          className={`fixed top-18 right-4 sm:right-6 z-50 flex items-center gap-2.5 px-4 py-3 rounded-2xl shadow-2xl border backdrop-blur-md font-medium text-xs max-w-md ${
            toast.type === "success"
              ? "bg-emerald-600/95 border-emerald-400/30 text-white"
              : toast.type === "error"
                ? "bg-rose-600/95 border-rose-400/30 text-white"
                : "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] border-[var(--md-sys-color-outline-variant)]/40"
          }`}
        >
          <Icon
            icon={
              toast.type === "success"
                ? "ph:check-circle-fill"
                : toast.type === "error"
                  ? "ph:warning-circle-fill"
                  : "ph:info-fill"
            }
            width="18"
            className="shrink-0"
          />
          <span className="flex-1 leading-snug">{toast.message}</span>
        </div>
      ) : null}

      {/* ── FLOATING LIVE DOCK (PROGRESS OR ERROR) ── */}
      {status?.is_running || status?.error ? (
        <div className="fixed bottom-6 right-4 sm:right-6 z-40">
          {status.is_running ? (
            floatingMinimized ? (
              <button
                type="button"
                onClick={() => onSetFloatingMinimized(false)}
                className="flex items-center gap-2.5 px-4 py-2.5 rounded-full bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface)] border-2 border-[var(--md-sys-color-primary)] shadow-2xl hover:scale-105 transition-all cursor-pointer font-bold text-xs"
              >
                <span className="w-2.5 h-2.5 rounded-full bg-[var(--md-sys-color-primary)] animate-ping" />
                <span className="font-mono text-[var(--md-sys-color-primary)]">
                  {status.pct.toFixed(0)}%
                </span>
                <span>
                  Memindai {status.progress}/{status.total}
                </span>
                <Icon
                  icon="ph:arrows-out-simple-bold"
                  width="14"
                  className="text-[var(--md-sys-color-on-surface-variant)]"
                />
              </button>
            ) : (
              <div className="w-[340px] sm:w-[380px] p-4.5 rounded-3xl bg-[var(--md-sys-color-surface-container)]/95 backdrop-blur-lg border border-[var(--md-sys-color-primary)]/40 shadow-2xl text-[var(--md-sys-color-on-surface)] relative overflow-hidden">
                <div className="absolute top-0 left-0 right-0 h-1.5 bg-gradient-to-r from-[var(--md-sys-color-primary)] via-emerald-500 to-[var(--md-sys-color-primary)] animate-pulse" />

                <div className="flex items-center justify-between mb-2.5">
                  <div className="flex items-center gap-2">
                    <div className="w-7 h-7 rounded-xl bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] flex items-center justify-center font-bold">
                      <Icon icon="ph:spinner-gap-bold" width="16" className="animate-spin" />
                    </div>
                    <div>
                      <span className="text-[10px] font-bold uppercase tracking-wider text-[var(--md-sys-color-primary)] font-mono">
                        LIVE SCANNING (SESI {status.session})
                      </span>
                      <h4 className="text-xs font-black">Generator Market Report</h4>
                    </div>
                  </div>

                  <div className="flex items-center gap-1">
                    <button
                      type="button"
                      onClick={onOpenLogcat}
                      className="px-2 py-1 rounded-lg bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-[10px] font-bold font-mono text-[var(--md-sys-color-primary)] flex items-center gap-1 transition cursor-pointer"
                      title="Buka Terminal Logcat"
                    >
                      <Icon icon="ph:terminal-window-bold" width="12" />
                      <span>Logcat</span>
                    </button>
                    <button
                      type="button"
                      onClick={() => onSetFloatingMinimized(true)}
                      className="p-1 rounded-lg hover:bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] transition cursor-pointer"
                      title="Kecilkan Notifikasi"
                    >
                      <Icon icon="ph:minus-bold" width="14" />
                    </button>
                  </div>
                </div>

                <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mb-2.5 line-clamp-2 leading-relaxed">
                  {status.message || "Memindai harga 962 emiten di Bursa Efek Indonesia..."}
                </p>

                <div className="w-full h-2.5 rounded-full bg-[var(--md-sys-color-surface-container-highest)] overflow-hidden p-0.5 border border-[var(--md-sys-color-outline-variant)] mb-2">
                  <div
                    className="h-full rounded-full bg-gradient-to-r from-[var(--md-sys-color-primary)] to-emerald-500 transition-all duration-300 ease-out shadow-xs"
                    style={{ width: `${Math.max(3, status.pct)}%` }}
                  />
                </div>

                <div className="flex items-center justify-between text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)] pt-1 border-t border-[var(--md-sys-color-outline-variant)]/60">
                  <span className="font-bold text-[var(--md-sys-color-on-surface)]">
                    {status.progress} / {status.total} emiten ({status.pct.toFixed(1)}%)
                  </span>
                  <span className="flex items-center gap-1">
                    <Icon icon="ph:clock-bold" width="12" />
                    {formatTimer(elapsedSeconds)}
                  </span>
                </div>
              </div>
            )
          ) : (
            <div className="w-[340px] sm:w-[400px] p-4.5 rounded-3xl bg-[var(--md-sys-color-surface-container)]/95 backdrop-blur-lg border border-rose-500/45 shadow-2xl text-[var(--md-sys-color-on-surface)] relative overflow-hidden">
              <div className="absolute top-0 left-0 right-0 h-1.5 bg-rose-500" />
              <div className="flex items-start justify-between gap-2 mb-2">
                <div className="flex items-center gap-2">
                  <div className="w-7 h-7 rounded-xl bg-rose-500/15 text-rose-600 dark:text-rose-400 flex items-center justify-center font-bold shrink-0">
                    <Icon icon="ph:warning-octagon-bold" width="18" />
                  </div>
                  <div>
                    <span className="text-[10px] font-bold uppercase tracking-wider text-rose-600 dark:text-rose-400 font-mono">
                      PEMINDAIAN TERHENTI
                    </span>
                    <h4 className="text-xs font-black text-rose-700 dark:text-rose-300">
                      Eksekusi Mengalami Error
                    </h4>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={onDismissError}
                  className="p-1 rounded-lg hover:bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] transition cursor-pointer"
                  title="Tutup Pesan"
                >
                  <Icon icon="ph:x-bold" width="14" />
                </button>
              </div>

              <p className="text-[11px] mb-3 bg-rose-500/10 p-2.5 rounded-xl border border-rose-500/25 font-mono text-rose-700 dark:text-rose-300 leading-relaxed">
                {status.error}
              </p>

              <div className="flex items-center gap-2">
                <button
                  type="button"
                  onClick={onOpenLogcat}
                  className="flex-1 py-2 px-3 rounded-xl bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-[11px] font-bold text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)] flex items-center justify-center gap-1.5 transition cursor-pointer"
                >
                  <Icon icon="ph:terminal-window-bold" width="14" />
                  <span>Buka Logcat</span>
                </button>
                <button
                  type="button"
                  onClick={onOpenTokenModal}
                  className="flex-1 py-2 px-3 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-[11px] font-bold flex items-center justify-center gap-1.5 transition cursor-pointer shadow-xs"
                >
                  <Icon icon="ph:key-bold" width="14" />
                  <span>Perbarui Token</span>
                </button>
              </div>
            </div>
          )}
        </div>
      ) : null}
    </>
  );
}
