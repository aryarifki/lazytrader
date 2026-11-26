import React from "react";
import { Icon } from "@iconify/react";
import { formatDuration, type ReportData } from "./marketReportTypes";

interface MarketReportViewerCardProps {
  activeReport: ReportData | null;
  onCopy: (text: string, label: string) => void;
}

export default function MarketReportViewerCard({
  activeReport,
  onCopy,
}: MarketReportViewerCardProps) {
  const breadth = activeReport?.market_breadth;
  const totalBreadth = breadth ? breadth.up + breadth.down + breadth.flat : 0;

  return (
    <div className="space-y-4">
      {/* Report Metadata & Market Breadth Summary */}
      {activeReport ? (
        <div className="p-4 rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] flex flex-col lg:flex-row lg:items-center justify-between gap-3">
          <div className="flex items-center gap-2 flex-wrap text-xs font-mono">
            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] font-bold">
              <Icon icon="ph:file-text-bold" width="14" />
              Sesi {activeReport.session}
            </span>
            {activeReport.date ? (
              <span className="px-2.5 py-1 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)] font-semibold">
                {activeReport.date}
              </span>
            ) : null}
            {activeReport.total_scanned ? (
              <span className="px-2.5 py-1 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)]">
                {activeReport.total_scanned} emiten
              </span>
            ) : null}
            {activeReport.duration_seconds ? (
              <span className="px-2.5 py-1 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)]">
                Durasi: {formatDuration(activeReport.duration_seconds)}
              </span>
            ) : null}
            {activeReport.generated_at ? (
              <span className="px-2.5 py-1 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)]">
                Dibuat: {activeReport.generated_at}
              </span>
            ) : null}
          </div>

          {breadth ? (
            <div className="flex items-center gap-2 flex-wrap text-xs font-mono">
              <span className="text-[11px] font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
                Breadth:
              </span>
              <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-xl bg-emerald-500/12 text-emerald-700 dark:text-emerald-400 border border-emerald-500/25 font-bold">
                <Icon icon="ph:trend-up-bold" width="13" />
                Naik {breadth.up}
              </span>
              <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-xl bg-rose-500/12 text-rose-700 dark:text-rose-400 border border-rose-500/25 font-bold">
                <Icon icon="ph:trend-down-bold" width="13" />
                Turun {breadth.down}
              </span>
              <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)] font-bold">
                <Icon icon="ph:minus-bold" width="13" />
                Tetap {breadth.flat}
              </span>
              {totalBreadth > 0 ? (
                <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
                  ({totalBreadth} total)
                </span>
              ) : null}
            </div>
          ) : null}
        </div>
      ) : null}

      {/* Dual Side-by-Side Output Viewers */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-5">
        {/* Left: Masked (Public) */}
        <div className="flex flex-col rounded-3xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] overflow-hidden shadow-xs">
          <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)] flex items-center justify-between gap-2 bg-[var(--md-sys-color-surface-container-high)]/50">
            <div className="flex items-center gap-2.5 min-w-0">
              <div className="w-8 h-8 rounded-xl bg-emerald-500/12 text-emerald-700 dark:text-emerald-400 flex items-center justify-center font-bold shrink-0">
                <Icon icon="ph:paper-plane-tilt-bold" width="16" />
              </div>
              <div className="min-w-0">
                <h3 className="text-sm font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5 flex-wrap">
                  Versi Masked (Publik)
                  <span className="text-[10px] font-mono font-bold px-1.5 py-0.5 rounded bg-emerald-500/12 text-emerald-700 dark:text-emerald-400 border border-emerald-500/25">
                    Siap Share
                  </span>
                </h3>
                <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] truncate">
                  Kode saham disamarkan (T**A, B**S) untuk channel Telegram &amp; WhatsApp
                </p>
              </div>
            </div>

            <button
              type="button"
              onClick={() => onCopy(activeReport?.masked_text || "", "Laporan Masked")}
              disabled={!activeReport?.masked_text}
              className="px-3 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs flex items-center gap-1.5 shadow-xs transition disabled:opacity-40 cursor-pointer shrink-0"
            >
              <Icon icon="ph:copy-bold" width="14" />
              <span>Salin Teks</span>
            </button>
          </div>

          <div className="p-4 sm:p-5 flex-1 flex flex-col">
            <textarea
              readOnly
              value={activeReport?.masked_text || "Belum ada laporan yang dipilih."}
              rows={18}
              className="w-full flex-1 p-3.5 rounded-2xl bg-[var(--md-sys-color-surface-container-lowest)] text-[var(--md-sys-color-on-surface)] font-mono text-xs leading-relaxed border border-[var(--md-sys-color-outline-variant)]/70 focus:outline-none resize-none selection:bg-emerald-500/30"
            />
            <div className="flex items-center justify-between text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-2.5 font-mono gap-2 flex-wrap">
              <span>Panjang Teks: {activeReport?.masked_text?.length || 0} karakter</span>
              <span>File: {activeReport?.filename || "-"}</span>
            </div>
          </div>
        </div>

        {/* Right: Unmasked (Internal) */}
        <div className="flex flex-col rounded-3xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] overflow-hidden shadow-xs">
          <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)] flex items-center justify-between gap-2 bg-[var(--md-sys-color-surface-container-high)]/50">
            <div className="flex items-center gap-2.5 min-w-0">
              <div className="w-8 h-8 rounded-xl bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] flex items-center justify-center font-bold shrink-0">
                <Icon icon="ph:eye-bold" width="16" />
              </div>
              <div className="min-w-0">
                <h3 className="text-sm font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5 flex-wrap">
                  Versi Unmasked (Internal)
                  <span className="text-[10px] font-mono font-bold px-1.5 py-0.5 rounded bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border border-[var(--md-sys-color-primary)]/20">
                    Full Code
                  </span>
                </h3>
                <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] truncate">
                  Kode saham lengkap (TPIA, BBRI, dsb.) untuk analisa internal tim
                </p>
              </div>
            </div>

            <button
              type="button"
              onClick={() => onCopy(activeReport?.unmasked_text || "", "Laporan Unmasked")}
              disabled={!activeReport?.unmasked_text}
              className="px-3 py-1.5 rounded-xl bg-[var(--md-sys-color-primary)] hover:opacity-90 text-[var(--md-sys-color-on-primary)] font-bold text-xs flex items-center gap-1.5 shadow-xs transition disabled:opacity-40 cursor-pointer shrink-0"
            >
              <Icon icon="ph:copy-bold" width="14" />
              <span>Salin Data</span>
            </button>
          </div>

          <div className="p-4 sm:p-5 flex-1 flex flex-col">
            <textarea
              readOnly
              value={activeReport?.unmasked_text || "Belum ada laporan yang dipilih."}
              rows={18}
              className="w-full flex-1 p-3.5 rounded-2xl bg-[var(--md-sys-color-surface-container-lowest)] text-[var(--md-sys-color-on-surface)] font-mono text-xs leading-relaxed border border-[var(--md-sys-color-outline-variant)]/70 focus:outline-none resize-none selection:bg-[var(--md-sys-color-primary)]/25"
            />
            <div className="flex items-center justify-between text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-2.5 font-mono gap-2 flex-wrap">
              <span>Panjang Teks: {activeReport?.unmasked_text?.length || 0} karakter</span>
              <span>Branding: by Investowl ID</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
