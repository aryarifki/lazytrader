import { formatIndoDate } from "./foreignProReportMath";
import type { AdminReportMode, IhsgResponse } from "./foreignProTypes";

interface ForeignProScreenerAdminProps {
  adminMode: AdminReportMode;
  setAdminMode: (m: AdminReportMode) => void;
  adminReportText: string;
  adminCopied: boolean;
  onCopyReport: () => void;
  ihsgData: IhsgResponse | null;
}

export default function ForeignProScreenerAdmin({
  adminMode,
  setAdminMode,
  adminReportText,
  adminCopied,
  onCopyReport,
  ihsgData,
}: ForeignProScreenerAdminProps) {
  return (
    <div
      id="adminPanel"
      className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-3 rounded-xl space-y-3"
    >
      <div className="flex items-center justify-between pb-2 border-b border-[var(--md-sys-color-outline-variant)]/40 flex-wrap gap-2">
        <div className="flex items-center gap-2">
          <span className="w-2.5 h-2.5 rounded-full bg-[var(--md-sys-color-primary)]" />
          <span className="font-bold text-[var(--md-sys-color-on-surface)] text-xs">
            Generator Market Report Resmi (InvestOwl Format)
          </span>
        </div>

        <div className="flex items-center gap-2">
          <button
            id="adminModeSesi1Btn"
            type="button"
            onClick={() => setAdminMode("sesi1")}
            className={`px-3 py-1 rounded-lg text-xs font-bold transition cursor-pointer ${
              adminMode === "sesi1"
                ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-2xs"
                : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)]"
            }`}
          >
            Sesi 1
          </button>
          <button
            id="adminModeEndBtn"
            type="button"
            onClick={() => setAdminMode("end")}
            className={`px-3 py-1 rounded-lg text-xs font-bold transition cursor-pointer ${
              adminMode === "end"
                ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-2xs"
                : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)]"
            }`}
          >
            Akhir Sesi
          </button>
        </div>
      </div>

      <div className="space-y-2">
        <textarea
          id="adminReportText"
          readOnly
          rows={16}
          value={adminReportText}
          className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/60 rounded p-3 text-xs text-[var(--md-sys-color-on-surface)] font-mono leading-relaxed focus:outline-none"
        />
        <div className="flex items-center justify-between flex-wrap gap-2">
          <button
            id="adminCopyBtn"
            type="button"
            onClick={onCopyReport}
            className="bg-[var(--md-sys-color-primary)] hover:opacity-90 text-[var(--md-sys-color-on-primary)] px-4 py-1.5 rounded-lg text-xs font-bold transition shadow-2xs cursor-pointer"
          >
            {adminCopied ? "✓ Berhasil Disalin!" : "Copy ke Clipboard"}
          </button>
          <span id="adminIhsgStatus" className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
            Status IHSG:{" "}
            <span className="font-semibold text-[var(--color-positive)]">
              {ihsgData?.summary?.status || "EOD"}
            </span>{" "}
            (Close:{" "}
            <span className="font-mono text-[var(--md-sys-color-on-surface)]">
              {ihsgData?.summary?.close || 5096.83}
            </span>
            ) •{" "}
            {ihsgData?.summary?.date ? formatIndoDate(ihsgData.summary.date) : "Senin, 02 Nov 2026"}
          </span>
        </div>
      </div>
    </div>
  );
}
