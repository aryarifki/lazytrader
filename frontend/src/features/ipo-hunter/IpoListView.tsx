import { useEffect, useState } from "react";
import { getIpoHunterList } from "@/lib/api/terminal";
import type { IpoApiResponse, IpoListData, IpoMeta } from "./ipoTypes";
import { IpoListInfoTable, IpoListPerfTable } from "./IpoListTables";

const VAL_BUCKETS = [
  ["", "Semua nilai"],
  ["lt100", "< 100 M"],
  ["100_500", "100 – 500 M"],
  ["500_1000", "500 M – 1 T"],
  ["1_5t", "1 – 5 T"],
  ["gt5t", "> 5 T"],
] as const;

export default function IpoListView({
  meta,
  onSelectTicker,
}: {
  meta: IpoMeta | null;
  onSelectTicker: (ticker: string) => void;
}) {
  const [viewMode, setViewMode] = useState<"info" | "perf">("info");
  const [uw, setUw] = useState("");
  const [leadOnly, setLeadOnly] = useState(false);
  const [valBucket, setValBucket] = useState("");
  const [year, setYear] = useState("");
  const [q, setQ] = useState("");
  const [sortBy, setSortBy] = useState("date");
  const [sortDir, setSortDir] = useState(-1);
  const [data, setData] = useState<IpoListData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    const t = setTimeout(
      () => {
        setLoading(true);
        getIpoHunterList<IpoApiResponse<IpoListData>>({
          uw: uw || undefined,
          lead_only: leadOnly,
          val_bucket: valBucket || undefined,
          year: year || undefined,
          q: q || undefined,
          sort: sortBy,
          dir: sortDir,
        })
          .then((res) => {
            if (cancelled) return;
            setData(res.data);
            setLoading(false);
          })
          .catch(() => {
            if (!cancelled) setLoading(false);
          });
      },
      q ? 200 : 0
    );
    return () => {
      cancelled = true;
      clearTimeout(t);
    };
  }, [uw, leadOnly, valBucket, year, q, sortBy, sortDir]);

  const handleSort = (col: string) => {
    if (sortBy === col) {
      setSortDir((d) => -d);
    } else {
      setSortBy(col);
      setSortDir(["kode", "nama", "papan", "uw"].includes(col) ? 1 : -1);
    }
  };

  const resetFilters = () => {
    setUw("");
    setLeadOnly(false);
    setValBucket("");
    setYear("");
    setQ("");
    setSortBy("date");
    setSortDir(-1);
  };

  const rows = data?.rows || [];

  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] shadow-sm overflow-hidden">
      <div className="flex flex-wrap items-center justify-between gap-3 px-4 py-3 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <div className="flex items-center gap-3">
          <h2 className="text-xs font-extrabold uppercase tracking-wider text-[var(--md-sys-color-on-surface)]">
            IPO LIST
          </h2>
          <div className="inline-flex rounded-lg bg-[var(--md-sys-color-surface-container-high)] p-0.5 border border-[var(--md-sys-color-outline-variant)]/50">
            {(
              [
                ["info", "Info"],
                ["perf", "Kinerja"],
              ] as const
            ).map(([m, lbl]) => (
              <button
                key={m}
                type="button"
                onClick={() => setViewMode(m)}
                className={`px-3 py-1 text-[11px] font-bold rounded-md transition cursor-pointer ${
                  viewMode === m
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                    : "text-[var(--md-sys-color-on-surface-variant)]"
                }`}
              >
                {lbl}
              </button>
            ))}
          </div>
        </div>
        <span className="text-xs font-mono-nums text-[var(--md-sys-color-on-surface-variant)]">
          <b className="text-[var(--md-sys-color-on-surface)]">{data?.total ?? 0}</b> dari{" "}
          {data?.total_all ?? meta?.total_ipos ?? 0} saham
        </span>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap items-end gap-3 px-4 py-3 bg-[var(--md-sys-color-surface-container-low)]/60 border-b border-[var(--md-sys-color-outline-variant)]/50 text-xs">
        <label className="flex flex-col gap-1">
          <span className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">Underwriter</span>
          <select
            value={uw}
            onChange={(e) => setUw(e.target.value)}
            className="h-8 px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs"
          >
            <option value="">Semua</option>
            {(meta?.underwriters || []).map((u) => (
              <option key={u.broker} value={u.broker}>
                {u.broker} ({u.count})
              </option>
            ))}
          </select>
        </label>

        <label className="flex items-center gap-1.5 h-8 px-2 cursor-pointer select-none">
          <input type="checkbox" checked={leadOnly} onChange={(e) => setLeadOnly(e.target.checked)} className="rounded" />
          <span className="text-xs font-semibold">Lead saja</span>
        </label>

        <label className="flex flex-col gap-1">
          <span className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">Nilai IPO</span>
          <select
            value={valBucket}
            onChange={(e) => setValBucket(e.target.value)}
            className="h-8 px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs"
          >
            {VAL_BUCKETS.map(([k, lbl]) => (
              <option key={k} value={k}>
                {lbl}
              </option>
            ))}
          </select>
        </label>

        <label className="flex flex-col gap-1">
          <span className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">Tahun IPO</span>
          <select
            value={year}
            onChange={(e) => setYear(e.target.value)}
            className="h-8 px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs"
          >
            <option value="">Semua</option>
            {(meta?.years || []).map((y) => (
              <option key={y} value={y}>
                {y}
              </option>
            ))}
          </select>
        </label>

        <label className="flex flex-col gap-1 flex-1 min-w-[160px] max-w-xs">
          <span className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">Kode / Nama</span>
          <input
            type="text"
            value={q}
            onChange={(e) => setQ(e.target.value)}
            placeholder="cari kode / nama…"
            className="h-8 px-2.5 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs"
          />
        </label>

        <button
          type="button"
          onClick={resetFilters}
          className="h-8 px-3 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] hover:border-[var(--md-sys-color-primary)] text-xs font-semibold cursor-pointer"
        >
          Reset
        </button>
      </div>

      <div className="px-4 py-1.5 text-[11px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/40">
        Kode Underwriter: warna = kategori broker (sama dengan Top Buyer/Seller) · <u className="font-semibold">garis bawah</u> = Lead
      </div>

      <div className="overflow-x-auto">
        {loading && !data ? (
          <div className="p-8 text-center text-xs text-[var(--md-sys-color-on-surface-variant)]">Memuat daftar IPO…</div>
        ) : viewMode === "info" ? (
          <IpoListInfoTable
            rows={rows}
            sortBy={sortBy}
            sortDir={sortDir}
            onSort={handleSort}
            onSelectTicker={onSelectTicker}
          />
        ) : (
          <IpoListPerfTable
            rows={rows}
            summary={data?.kinerja_summary}
            sortBy={sortBy}
            sortDir={sortDir}
            onSort={handleSort}
            onSelectTicker={onSelectTicker}
          />
        )}
      </div>
    </div>
  );
}
