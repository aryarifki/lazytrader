import { fmtDate, fmtNumber } from "@/components/research/format";
import type { Mode } from "./snapshotTypes";

export function SnapshotPagination({
  mode,
  shown,
  total,
  offset,
  limit,
  priceDate,
  source,
  onPrevPage,
  onNextPage,
}: {
  mode: Mode;
  shown: number;
  total: number;
  offset: number;
  limit: number;
  priceDate: string | null | undefined;
  source: string;
  onPrevPage: () => void;
  onNextPage: () => void;
}) {
  return (
    <>
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-t border-[var(--md-sys-color-outline-variant)]/60 p-4">
        <div className="text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
          Menampilkan {fmtNumber(shown)} dari {fmtNumber(total)} baris
          {mode === "stocks" ? ` • Harga per ${fmtDate(priceDate)}` : ""}
        </div>
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={onPrevPage}
            disabled={offset === 0}
            className="rounded-lg border border-[var(--md-sys-color-outline-variant)] px-2.5 py-1.5 text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] disabled:cursor-not-allowed disabled:opacity-40 hover:border-[var(--md-sys-color-primary)]/60"
          >
            Sebelumnya
          </button>
          <span className="text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
            {total === 0 ? "0" : `${fmtNumber(offset + 1)}–${fmtNumber(offset + shown)}`}
          </span>
          <button
            type="button"
            onClick={onNextPage}
            disabled={offset + limit >= total}
            className="rounded-lg border border-[var(--md-sys-color-outline-variant)] px-2.5 py-1.5 text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] disabled:cursor-not-allowed disabled:opacity-40 hover:border-[var(--md-sys-color-primary)]/60"
          >
            Berikutnya
          </button>
        </div>
      </div>

      <div className="px-4 pb-4">
        <div className="flex flex-wrap items-center gap-2 text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
          <span>{source}</span>
          <span className="opacity-40">•</span>
          <span>Snapshot, bukan delta bulanan</span>
          <span className="opacity-40">•</span>
          <span>Market cap = listed shares IDX</span>
          <span className="opacity-40">•</span>
          <span>Sorting & filter di server</span>
        </div>
      </div>
    </>
  );
}
