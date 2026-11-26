import React from "react";
import { Icon } from "@iconify/react";
import MarketReportHistoryBar from "./MarketReportHistoryBar";
import {
  formatTimer,
  getMarketSessionStatus,
  type HistoryItem,
  type StatusResponse,
} from "./marketReportTypes";

interface MarketReportControlCardProps {
  status: StatusResponse | null;
  history: HistoryItem[];
  selectedHistoryFile: string;
  loading: boolean;
  actionLoading: boolean;
  elapsedSeconds: number;
  adminKeyState: string;
  rateLimit: number;
  workers: number;
  onRateLimitChange: (value: number) => void;
  onWorkersChange: (value: number) => void;
  onGenerate: (session: number) => void;
  onRefresh: () => void;
  onAutoRenew: () => void;
  onResetStatus: () => void;
  onOpenTokenModal: () => void;
  onOpenLogcatModal: () => void;
  onOpenAdminKeyModal: () => void;
  onSelectHistory: (filename: string) => void;
}

export default function MarketReportControlCard({
  status,
  history,
  selectedHistoryFile,
  loading,
  actionLoading,
  elapsedSeconds,
  adminKeyState,
  rateLimit,
  workers,
  onRateLimitChange,
  onWorkersChange,
  onGenerate,
  onRefresh,
  onAutoRenew,
  onResetStatus,
  onOpenTokenModal,
  onOpenLogcatModal,
  onOpenAdminKeyModal,
  onSelectHistory,
}: MarketReportControlCardProps) {
  const marketStatus = getMarketSessionStatus();
  const isRunning = Boolean(status?.is_running);
  const tokenValid = Boolean(status?.token_valid);

  return (
    <div className="space-y-4">
      <section className="p-5 sm:p-6 rounded-3xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] shadow-xs">
        <div className="flex flex-col xl:flex-row xl:items-start justify-between gap-5">
          {/* Left: Badges & Title */}
          <div className="space-y-2.5 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] text-xs font-bold border border-[var(--md-sys-color-primary)]/20 font-mono">
                <Icon icon="ph:newspaper-clipping-bold" width="14" />
                ADMIN REPORT ENGINE
              </span>
              <span
                className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-xl text-xs font-semibold border ${marketStatus.badgeClass}`}
              >
                <span className="w-2 h-2 rounded-full bg-current animate-pulse" />
                {marketStatus.label}
              </span>
              {status?.current_time_wib ? (
                <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-xl text-[11px] font-mono font-semibold bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)]">
                  <Icon icon="ph:clock-bold" width="13" />
                  {status.current_time_wib}
                </span>
              ) : null}
              <span
                className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-xl text-xs font-semibold border ${
                  tokenValid
                    ? "text-emerald-700 dark:text-emerald-400 bg-emerald-500/12 border-emerald-500/25"
                    : "text-rose-700 dark:text-rose-400 bg-rose-500/12 border-rose-500/25"
                }`}
              >
                <Icon icon={tokenValid ? "ph:shield-check-fill" : "ph:shield-warning-fill"} width="14" />
                Token Stockbit: {tokenValid ? "Aktif" : "Kedaluwarsa"}
              </span>
            </div>

            <div className="flex items-center gap-2 flex-wrap pt-0.5">
              <button
                type="button"
                onClick={onAutoRenew}
                disabled={isRunning || actionLoading}
                className="px-2.5 py-1 rounded-xl bg-emerald-500/12 hover:bg-emerald-500/20 text-xs font-bold font-mono text-emerald-700 dark:text-emerald-400 border border-emerald-500/30 flex items-center gap-1.5 transition cursor-pointer disabled:opacity-50"
                title="Pembaruan Token Otomatis Two-Tier via SMtracker"
              >
                <Icon icon="ph:device-mobile-bold" width="13" />
                <span>Auto-Renew (HP)</span>
              </button>
              <button
                type="button"
                onClick={onOpenTokenModal}
                className="px-2.5 py-1 rounded-xl bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-xs font-bold font-mono text-[var(--md-sys-color-primary)] border border-[var(--md-sys-color-outline-variant)] flex items-center gap-1.5 transition cursor-pointer"
                title="Perbarui Token Stockbit Secara Manual"
              >
                <Icon icon="ph:key-bold" width="13" />
                <span>Update Manual</span>
              </button>
              <button
                type="button"
                onClick={onOpenLogcatModal}
                className="px-2.5 py-1 rounded-xl bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-xs font-bold font-mono text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)] flex items-center gap-1.5 transition cursor-pointer"
                title="Lihat Logcat Terminal Eksekusi"
              >
                <Icon icon="ph:terminal-window-bold" width="13" />
                <span>Logcat</span>
              </button>
              <button
                type="button"
                onClick={onOpenAdminKeyModal}
                className={`px-2.5 py-1 rounded-xl text-xs font-bold font-mono border flex items-center gap-1.5 transition cursor-pointer ${
                  adminKeyState
                    ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border-[var(--md-sys-color-primary)]/30 hover:opacity-90"
                    : "bg-emerald-500/12 hover:bg-emerald-500/20 text-emerald-700 dark:text-emerald-400 border-emerald-500/30"
                }`}
                title={
                  adminKeyState
                    ? "Override X-Admin-Key aktif (Klik untuk ubah/hapus)"
                    : "Terotorisasi via Sesi Admin InvestOwl (Klik untuk tambah override X-Admin-Key opsional)"
                }
              >
                <Icon icon={adminKeyState ? "ph:key-return-bold" : "ph:shield-check-bold"} width="13" />
                <span>{adminKeyState ? "Admin Key: Aktif (Override)" : "Sesi Admin Aktif"}</span>
              </button>
              {(isRunning || status?.error) && (
                <button
                  type="button"
                  onClick={onResetStatus}
                  disabled={actionLoading}
                  className="px-2.5 py-1 rounded-xl bg-rose-500/12 hover:bg-rose-500/20 text-xs font-bold font-mono text-rose-700 dark:text-rose-400 border border-rose-500/30 flex items-center gap-1 transition cursor-pointer disabled:opacity-50"
                  title="Reset status generator apabila proses macet atau mengalami error"
                >
                  <Icon icon="ph:arrow-counter-clockwise-bold" width="13" />
                  <span>Reset Status</span>
                </button>
              )}
            </div>

            <div className="pt-1">
              <h1 className="text-xl sm:text-2xl font-black tracking-tight text-[var(--md-sys-color-on-surface)]">
                Generator Market Report Resmi
                <span className="ml-2 text-xs font-bold font-mono px-2 py-0.5 rounded-lg bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] align-middle">
                  by Investowl ID
                </span>
              </h1>
              <p className="text-xs text-[var(--md-sys-color-on-surface-variant)] mt-1">
                Pemindaian data harga dan arus dana asing real-time seluruh universe BEI (962 emiten) dengan format standar InvestOwl.
              </p>
            </div>
          </div>

          {/* Right: Rate Limit / Workers & Generate Actions */}
          <div className="flex flex-col items-stretch xl:items-end gap-3 shrink-0">
            <div className="flex items-center gap-2 flex-wrap">
              <label className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-[11px] font-mono">
                <span className="text-[var(--md-sys-color-on-surface-variant)] font-semibold">Rate/m:</span>
                <input
                  type="number"
                  min={30}
                  max={600}
                  step={10}
                  value={rateLimit}
                  disabled={isRunning}
                  onChange={(e) => onRateLimitChange(Math.max(30, Number(e.target.value) || 180))}
                  className="w-14 bg-transparent text-[var(--md-sys-color-on-surface)] font-bold focus:outline-none text-right"
                />
              </label>
              <label className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-[11px] font-mono">
                <span className="text-[var(--md-sys-color-on-surface-variant)] font-semibold">Workers:</span>
                <input
                  type="number"
                  min={1}
                  max={32}
                  step={1}
                  value={workers}
                  disabled={isRunning}
                  onChange={(e) => onWorkersChange(Math.max(1, Math.min(32, Number(e.target.value) || 8)))}
                  className="w-10 bg-transparent text-[var(--md-sys-color-on-surface)] font-bold focus:outline-none text-right"
                />
              </label>
            </div>

            <div className="flex items-center gap-2.5 flex-wrap">
              <button
                type="button"
                onClick={() => onGenerate(1)}
                disabled={isRunning || actionLoading || !tokenValid}
                className="px-4 py-2.5 rounded-2xl bg-[var(--md-sys-color-primary)] hover:opacity-90 text-[var(--md-sys-color-on-primary)] font-bold text-xs flex items-center gap-2 shadow-xs transition disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                <Icon
                  icon={isRunning && status?.session === 1 ? "ph:spinner-gap-bold" : "ph:sun-dim-bold"}
                  width="16"
                  className={isRunning && status?.session === 1 ? "animate-spin" : ""}
                />
                <span>Generate Sesi 1 (12:00)</span>
              </button>
              <button
                type="button"
                onClick={() => onGenerate(2)}
                disabled={isRunning || actionLoading || !tokenValid}
                className="px-4 py-2.5 rounded-2xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs flex items-center gap-2 shadow-xs transition disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                <Icon
                  icon={isRunning && status?.session === 2 ? "ph:spinner-gap-bold" : "ph:moon-stars-bold"}
                  width="16"
                  className={isRunning && status?.session === 2 ? "animate-spin" : ""}
                />
                <span>Generate End-of-Day (16:00)</span>
              </button>
              <button
                type="button"
                onClick={onRefresh}
                disabled={loading}
                className="p-2.5 rounded-2xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)] hover:bg-[var(--md-sys-color-surface-container-highest)] transition cursor-pointer"
                title="Refresh Data & Status"
              >
                <Icon icon="ph:arrows-clockwise-bold" width="16" className={loading ? "animate-spin" : ""} />
              </button>
            </div>
          </div>
        </div>

        {/* Persistent Error Banner */}
        {status?.error ? (
          <div className="mt-5 p-4 rounded-2xl bg-rose-500/10 border border-rose-500/30 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
            <div className="flex items-center gap-2.5">
              <Icon icon="ph:warning-circle-fill" width="20" className="text-rose-600 dark:text-rose-400 shrink-0" />
              <div>
                <p className="text-xs font-bold text-rose-700 dark:text-rose-300">Terjadi Kesalahan Saat Eksekusi:</p>
                <p className="text-xs font-mono text-rose-600 dark:text-rose-400">{status.error}</p>
              </div>
            </div>
            <div className="flex items-center gap-2 self-end sm:self-auto">
              <button
                type="button"
                onClick={onOpenLogcatModal}
                className="px-3 py-1.5 rounded-xl bg-rose-500/15 hover:bg-rose-500/25 text-rose-700 dark:text-rose-200 text-xs font-bold border border-rose-500/30 flex items-center gap-1 transition cursor-pointer font-mono"
              >
                <Icon icon="ph:terminal-window-bold" width="14" />
                <span>Lihat Logcat</span>
              </button>
              <button
                type="button"
                onClick={onOpenTokenModal}
                className="px-3 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-bold flex items-center gap-1 transition cursor-pointer"
              >
                <Icon icon="ph:key-bold" width="14" />
                <span>Perbarui Token</span>
              </button>
              <button
                type="button"
                onClick={onResetStatus}
                className="p-1.5 rounded-xl hover:bg-rose-500/20 text-rose-600 dark:text-rose-400 transition cursor-pointer"
                title="Tutup & reset status error"
              >
                <Icon icon="ph:x-bold" width="14" />
              </button>
            </div>
          </div>
        ) : null}

        {/* Inline Live Progress Bar */}
        {isRunning && status ? (
          <div className="mt-5 p-4 rounded-2xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]">
            <div className="flex items-center justify-between text-xs mb-2 gap-2 flex-wrap">
              <span className="font-bold flex items-center gap-2 text-[var(--md-sys-color-primary)]">
                <span className="w-2 h-2 rounded-full bg-[var(--md-sys-color-primary)] animate-ping" />
                {status.message || "Memproses pemindaian universe..."}
              </span>
              <span className="font-mono font-bold text-[var(--md-sys-color-on-surface)]">
                {status.progress} / {status.total} emiten ({status.pct.toFixed(1)}%)
              </span>
            </div>
            <div className="w-full h-3 rounded-full bg-[var(--md-sys-color-surface-container-highest)] overflow-hidden p-0.5 border border-[var(--md-sys-color-outline-variant)]">
              <div
                className="h-full rounded-full bg-gradient-to-r from-[var(--md-sys-color-primary)] to-emerald-500 transition-all duration-300 ease-out"
                style={{ width: `${Math.max(2, status.pct)}%` }}
              />
            </div>
            <div className="flex items-center justify-between text-[10px] text-[var(--md-sys-color-on-surface-variant)] mt-1.5 font-mono">
              <span>Waktu Berjalan: {formatTimer(elapsedSeconds)}</span>
              <span>
                Rate Limit: {rateLimit} req/menit | {workers} Worker Threads
              </span>
            </div>
          </div>
        ) : null}
      </section>

      <MarketReportHistoryBar
        history={history}
        selectedHistoryFile={selectedHistoryFile}
        onSelectHistory={onSelectHistory}
      />
    </div>
  );
}
