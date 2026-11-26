import React from "react";
import { Icon } from "@iconify/react";
import MarketReportLogcatModal from "./MarketReportLogcatModal";
import type { StatusResponse } from "./marketReportTypes";

interface MarketReportModalsProps {
  status: StatusResponse | null;
  actionLoading: boolean;
  showLogcatModal: boolean;
  logs: string[];
  autoScrollLogs: boolean;
  logLoading: boolean;
  logContainerRef: React.RefObject<HTMLDivElement | null>;
  onCloseLogcat: () => void;
  onToggleAutoScroll: () => void;
  onRefreshLogs: () => void;
  onCopyLogs: () => void;
  showTokenModal: boolean;
  tokenInput: string;
  tokenSubmitting: boolean;
  tokenError: string | null;
  tokenSuccess: string | null;
  onTokenInputChange: (value: string) => void;
  onCloseTokenModal: () => void;
  onSubmitToken: (e: React.FormEvent) => void;
  onAutoRenewFromModal: () => void;
  showAdminKeyModal: boolean;
  adminKeyState: string;
  adminKeyInput: string;
  onAdminKeyInputChange: (value: string) => void;
  onCloseAdminKeyModal: () => void;
  onSubmitAdminKey: (e: React.FormEvent) => void;
  onClearAdminKey: () => void;
}

export default function MarketReportModals({
  status,
  actionLoading,
  showLogcatModal,
  logs,
  autoScrollLogs,
  logLoading,
  logContainerRef,
  onCloseLogcat,
  onToggleAutoScroll,
  onRefreshLogs,
  onCopyLogs,
  showTokenModal,
  tokenInput,
  tokenSubmitting,
  tokenError,
  tokenSuccess,
  onTokenInputChange,
  onCloseTokenModal,
  onSubmitToken,
  onAutoRenewFromModal,
  showAdminKeyModal,
  adminKeyState,
  adminKeyInput,
  onAdminKeyInputChange,
  onCloseAdminKeyModal,
  onSubmitAdminKey,
  onClearAdminKey,
}: MarketReportModalsProps) {
  return (
    <>
      <MarketReportLogcatModal
        open={showLogcatModal}
        status={status}
        logs={logs}
        autoScrollLogs={autoScrollLogs}
        logLoading={logLoading}
        logContainerRef={logContainerRef}
        onClose={onCloseLogcat}
        onToggleAutoScroll={onToggleAutoScroll}
        onRefreshLogs={onRefreshLogs}
        onCopyLogs={onCopyLogs}
      />

      {/* ── UPDATE TOKEN STOCKBIT MODAL ── */}
      {showTokenModal ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
          <div className="w-full max-w-lg rounded-3xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] shadow-2xl p-6 relative">
            <div className="flex items-center justify-between pb-4 border-b border-[var(--md-sys-color-outline-variant)] mb-4">
              <div className="flex items-center gap-2.5">
                <div className="w-8 h-8 rounded-xl bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 flex items-center justify-center font-bold">
                  <Icon icon="ph:key-bold" width="18" />
                </div>
                <div>
                  <h3 className="text-sm font-bold text-[var(--md-sys-color-on-surface)]">
                    Perbarui Token Stockbit
                  </h3>
                  <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
                    Validasi dan perbarui kredensial Bearer API Stockbit
                  </p>
                </div>
              </div>
              <button
                type="button"
                onClick={onCloseTokenModal}
                className="p-1.5 rounded-xl hover:bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] transition cursor-pointer"
              >
                <Icon icon="ph:x-bold" width="16" />
              </button>
            </div>

            <div className="p-3.5 rounded-2xl bg-emerald-500/10 border border-emerald-500/30 mb-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
              <div>
                <span className="text-[10px] font-mono font-bold uppercase tracking-wider text-emerald-700 dark:text-emerald-400 flex items-center gap-1">
                  <Icon icon="ph:sparkle-bold" width="12" />
                  METODE OTOMATIS SMTRACKER (BACKFILL DAILY)
                </span>
                <h4 className="text-xs font-bold text-[var(--md-sys-color-on-surface)] mt-0.5">
                  Two-Tier Auto-Renew (Tanpa Copy-Paste)
                </h4>
                <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] leading-snug">
                  Menjalankan login headless SMtracker &amp; mengirim konfirmasi ke aplikasi HP Stockbit Anda.
                </p>
              </div>
              <button
                type="button"
                onClick={onAutoRenewFromModal}
                disabled={status?.is_running || actionLoading}
                className="px-3.5 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs flex items-center gap-1.5 shadow-xs transition shrink-0 cursor-pointer disabled:opacity-50"
              >
                <Icon icon="ph:device-mobile-bold" width="15" />
                <span>Mulai Auto-Renew</span>
              </button>
            </div>

            <div className="p-3.5 rounded-2xl bg-[var(--md-sys-color-surface-container-high)]/60 border border-[var(--md-sys-color-outline-variant)] text-xs text-[var(--md-sys-color-on-surface-variant)] mb-4 leading-relaxed space-y-1.5">
              <div className="font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5 text-[11px]">
                <Icon icon="ph:lightbulb-filament-bold" width="14" className="text-amber-600 dark:text-amber-400" />
                Panduan Cepat Menyalin Token:
              </div>
              <ol className="list-decimal list-inside space-y-1 text-[11px] pl-1">
                <li>Buka peramban di komputer Anda dan login ke <b>stockbit.com</b>.</li>
                <li>Tekan tombol <b>F12</b> (Inspect) lalu buka tab <b>Network</b>.</li>
                <li>
                  Ketik{" "}
                  <code className="px-1 py-0.5 rounded bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)]">
                    marketdetectors
                  </code>{" "}
                  atau{" "}
                  <code className="px-1 py-0.5 rounded bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)]">
                    price-feed
                  </code>{" "}
                  di kolom filter.
                </li>
                <li>Klik salah satu request, buka tab <b>Headers</b>, dan cari <b>Authorization</b>.</li>
                <li>
                  Salin seluruh nilai token (format{" "}
                  <code className="text-emerald-700 dark:text-emerald-400">Bearer eyJ...</code> atau{" "}
                  <code className="text-emerald-700 dark:text-emerald-400">eyJ...</code>).
                </li>
              </ol>
            </div>

            <form onSubmit={onSubmitToken} className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-[var(--md-sys-color-on-surface)] mb-1.5">
                  Tempel Authorization Token / JWT:
                </label>
                <textarea
                  value={tokenInput}
                  onChange={(e) => onTokenInputChange(e.target.value)}
                  placeholder="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                  rows={4}
                  required
                  className="w-full p-3.5 rounded-2xl bg-[var(--md-sys-color-surface-container-lowest)] text-[var(--md-sys-color-on-surface)] font-mono text-xs border border-[var(--md-sys-color-outline-variant)] focus:outline-none focus:border-[var(--md-sys-color-primary)] transition resize-none"
                />
              </div>

              {tokenError ? (
                <div className="p-3 rounded-xl bg-rose-500/12 border border-rose-500/30 text-rose-700 dark:text-rose-300 text-xs font-mono flex items-start gap-2">
                  <Icon icon="ph:warning-circle-fill" width="16" className="text-rose-600 dark:text-rose-400 shrink-0 mt-0.5" />
                  <span>{tokenError}</span>
                </div>
              ) : null}

              {tokenSuccess ? (
                <div className="p-3 rounded-xl bg-emerald-500/12 border border-emerald-500/30 text-emerald-700 dark:text-emerald-300 text-xs font-mono flex items-start gap-2">
                  <Icon icon="ph:check-circle-fill" width="16" className="text-emerald-600 dark:text-emerald-400 shrink-0 mt-0.5" />
                  <span>{tokenSuccess}</span>
                </div>
              ) : null}

              <div className="flex items-center justify-end gap-2.5 pt-2">
                <button
                  type="button"
                  onClick={onCloseTokenModal}
                  disabled={tokenSubmitting}
                  className="px-4 py-2.5 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] text-xs font-bold hover:bg-[var(--md-sys-color-surface-container-highest)] transition cursor-pointer"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  disabled={tokenSubmitting || !tokenInput.trim()}
                  className="px-5 py-2.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs flex items-center gap-2 shadow-xs transition disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                >
                  <Icon
                    icon={tokenSubmitting ? "ph:spinner-gap-bold" : "ph:check-bold"}
                    width="16"
                    className={tokenSubmitting ? "animate-spin" : ""}
                  />
                  <span>{tokenSubmitting ? "Memverifikasi Token..." : "Validasi & Simpan Token"}</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      ) : null}

      {/* ── OPTIONAL X-ADMIN-KEY MODAL ── */}
      {showAdminKeyModal ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
          <div className="w-full max-w-md rounded-3xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] shadow-2xl p-6 relative">
            <div className="flex items-center justify-between pb-4 border-b border-[var(--md-sys-color-outline-variant)] mb-4">
              <div className="flex items-center gap-2.5">
                <div className="w-8 h-8 rounded-xl bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] flex items-center justify-center font-bold">
                  <Icon icon="ph:shield-check-bold" width="18" />
                </div>
                <div>
                  <h3 className="text-sm font-black text-[var(--md-sys-color-on-surface)]">
                    Kunci Admin Terminal (Opsional)
                  </h3>
                  <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
                    Override header X-Admin-Key di luar sesi admin
                  </p>
                </div>
              </div>
              <button
                type="button"
                onClick={onCloseAdminKeyModal}
                className="p-1.5 rounded-xl hover:bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] transition cursor-pointer"
              >
                <Icon icon="ph:x-bold" width="16" />
              </button>
            </div>

            <form onSubmit={onSubmitAdminKey} className="space-y-4">
              <p className="text-xs text-[var(--md-sys-color-on-surface-variant)] leading-relaxed">
                Sesi admin InvestOwl Anda sudah otomatis terotorisasi melalui cookie terenkripsi. Kolom ini bersifat opsional jika Anda perlu memasang header <code className="font-mono">X-Admin-Key</code> tambahan.
              </p>

              <div>
                <label className="block text-xs font-bold text-[var(--md-sys-color-on-surface)] mb-1.5">
                  Terminal Admin Key (Override Opsional)
                </label>
                <input
                  type="password"
                  value={adminKeyInput}
                  onChange={(e) => onAdminKeyInputChange(e.target.value)}
                  placeholder="Masukkan TERMINAL_ADMIN_KEY..."
                  className="w-full px-3.5 py-2.5 rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface)] text-xs font-mono placeholder:text-[var(--md-sys-color-on-surface-variant)]/60 focus:outline-none focus:border-[var(--md-sys-color-primary)] transition"
                />
              </div>

              <div className="flex items-center justify-between pt-2">
                {adminKeyState ? (
                  <button
                    type="button"
                    onClick={onClearAdminKey}
                    className="px-3 py-2 rounded-xl text-rose-700 dark:text-rose-400 hover:bg-rose-500/12 text-xs font-bold transition cursor-pointer"
                  >
                    Hapus Override
                  </button>
                ) : (
                  <div />
                )}

                <div className="flex items-center gap-2">
                  <button
                    type="button"
                    onClick={onCloseAdminKeyModal}
                    className="px-4 py-2.5 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] text-xs font-bold hover:bg-[var(--md-sys-color-surface-container-highest)] transition cursor-pointer"
                  >
                    Batal
                  </button>
                  <button
                    type="submit"
                    className="px-5 py-2.5 rounded-xl bg-[var(--md-sys-color-primary)] hover:opacity-90 text-[var(--md-sys-color-on-primary)] font-bold text-xs flex items-center gap-1.5 shadow-xs transition cursor-pointer"
                  >
                    <Icon icon="ph:check-bold" width="14" />
                    <span>Simpan Kunci</span>
                  </button>
                </div>
              </div>
            </form>
          </div>
        </div>
      ) : null}
    </>
  );
}
