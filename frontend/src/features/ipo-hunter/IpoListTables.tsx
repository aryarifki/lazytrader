import React from "react";
import type { IpoKinerjaSummary, IpoListRow } from "./ipoTypes";
import { uwChipClass } from "./ipoFormat";

const PERF_COLS = ["d1", "d2", "d3", "d4", "d5", "d6", "d7", "d8", "d9", "d10", "perf2w", "perf1m", "perf3m"] as const;

export function fmtPerf(v: number | null | undefined): React.ReactNode {
  if (v == null) return <span className="text-[var(--md-sys-color-on-surface-variant)]/50">—</span>;
  const rv = Math.round(v) || 0;
  const cls =
    rv > 0
      ? "text-emerald-600 dark:text-emerald-400 font-bold"
      : rv < 0
      ? "text-rose-600 dark:text-rose-400 font-bold"
      : "";
  return (
    <span className={cls}>
      {rv > 0 ? "+" : ""}
      {rv.toLocaleString("en-US")}%
    </span>
  );
}

interface SortThProps {
  col: string;
  label: string;
  align?: "left" | "right";
  sortBy: string;
  sortDir: number;
  onSort: (col: string) => void;
}

function SortTh({ col, label, align = "left", sortBy, sortDir, onSort }: SortThProps) {
  return (
    <th
      onClick={() => onSort(col)}
      className={`px-3 py-2.5 text-[10px] font-bold uppercase tracking-wider cursor-pointer select-none whitespace-nowrap border-b border-[var(--md-sys-color-outline-variant)] ${
        align === "right" ? "text-right" : "text-left"
      } ${sortBy === col ? "text-[var(--md-sys-color-primary)]" : "text-[var(--md-sys-color-on-surface-variant)]"}`}
    >
      {label}
      {sortBy === col ? (sortDir > 0 ? " ▲" : " ▼") : ""}
    </th>
  );
}

interface TableProps {
  rows: IpoListRow[];
  sortBy: string;
  sortDir: number;
  onSort: (col: string) => void;
  onSelectTicker: (ticker: string) => void;
  summary?: IpoKinerjaSummary;
}

export function IpoListInfoTable({ rows, sortBy, sortDir, onSort, onSelectTicker }: TableProps) {
  return (
    <table className="w-full text-xs border-collapse">
      <thead>
        <tr className="bg-[var(--md-sys-color-surface-container-high)]/40">
          <SortTh col="kode" label="Kode" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="nama" label="Nama Emiten" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="price" label="IPO Price" align="right" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="nilai" label="Nilai IPO (M)" align="right" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="papan" label="Papan" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="uw" label="Underwriter" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="sid" label="Jumlah SID" align="right" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="perf2w" label="2W" align="right" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="perf1m" label="1M" align="right" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="perf3m" label="3M" align="right" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
        </tr>
      </thead>
      <tbody>
        {rows.length === 0 ? (
          <tr>
            <td colSpan={10} className="p-8 text-center text-[var(--md-sys-color-on-surface-variant)]">
              Tidak ada saham yang cocok dengan filter.
            </td>
          </tr>
        ) : (
          rows.map((r) => (
            <tr
              key={r.ticker}
              onClick={() => onSelectTicker(r.ticker)}
              className="border-b border-[var(--md-sys-color-outline-variant)]/40 hover:bg-[var(--md-sys-color-surface-container-high)]/60 cursor-pointer transition-colors"
            >
              <td className="px-3 py-2 font-mono-nums font-extrabold text-[var(--md-sys-color-primary)]">{r.ticker}</td>
              <td className="px-3 py-2 max-w-[240px] truncate">{r.company_name || "—"}</td>
              <td className="px-3 py-2 text-right font-mono-nums">{r.ipo_price?.toLocaleString("en-US") ?? "—"}</td>
              <td className="px-3 py-2 text-right font-mono-nums">
                {r.nilai_miliar != null ? Math.round(r.nilai_miliar).toLocaleString("en-US") : "—"}
              </td>
              <td className="px-3 py-2 whitespace-nowrap">{r.papan || "—"}</td>
              <td className="px-3 py-2">
                <div className="flex flex-wrap gap-1">
                  {(r.underwriters || []).length === 0
                    ? "—"
                    : r.underwriters.map((u) => (
                        <span
                          key={u.broker}
                          title={u.role}
                          className={`px-1.5 py-0.5 rounded border text-[10px] font-mono-nums font-bold ${uwChipClass(
                            u.category
                          )} ${u.role === "LEAD" ? "underline decoration-2" : ""}`}
                        >
                          {u.broker}
                        </span>
                      ))}
                </div>
              </td>
              <td className="px-3 py-2 text-right font-mono-nums">
                {r.jumlah_sid != null ? r.jumlah_sid.toLocaleString("en-US") : "—"}
              </td>
              <td className="px-3 py-2 text-right font-mono-nums">{fmtPerf(r.perf_2w)}</td>
              <td className="px-3 py-2 text-right font-mono-nums">{fmtPerf(r.perf_1m)}</td>
              <td className="px-3 py-2 text-right font-mono-nums">{fmtPerf(r.perf_3m)}</td>
            </tr>
          ))
        )}
      </tbody>
    </table>
  );
}

export function IpoListPerfTable({ rows, summary, sortBy, sortDir, onSort, onSelectTicker }: TableProps) {
  return (
    <table className="w-full text-xs border-collapse">
      <thead>
        <tr className="bg-[var(--md-sys-color-surface-container-high)]/40">
          <SortTh col="kode" label="Kode" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="date" label="Tgl IPO" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="penjatahan" label="Penjatahan" align="right" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((n) => (
            <SortTh
              key={n}
              col={`d${n}`}
              label={n === 1 ? "Day One" : `Day ${n}`}
              align="right"
              sortBy={sortBy}
              sortDir={sortDir}
              onSort={onSort}
            />
          ))}
          <SortTh col="perf2w" label="2W" align="right" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="perf1m" label="1M" align="right" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
          <SortTh col="perf3m" label="3M" align="right" sortBy={sortBy} sortDir={sortDir} onSort={onSort} />
        </tr>
      </thead>
      <tbody>
        {rows.length > 0 && (
          <>
            <tr className="bg-[var(--md-sys-color-surface-container-high)]/70 font-bold border-b border-[var(--md-sys-color-outline-variant)]">
              <td colSpan={3} className="px-3 py-2 text-left">
                Rata-rata
              </td>
              {PERF_COLS.map((k) => (
                <td key={`avg-${k}`} className="px-3 py-2 text-right font-mono-nums">
                  {fmtPerf(summary?.rata_rata?.[k])}
                </td>
              ))}
            </tr>
            <tr className="bg-[var(--md-sys-color-surface-container-high)]/50 font-bold border-b-2 border-[var(--md-sys-color-outline-variant)]">
              <td colSpan={3} className="px-3 py-2 text-left">
                Kemungkinan Harga Naik
              </td>
              {PERF_COLS.map((k) => (
                <td key={`up-${k}`} className="px-3 py-2 text-right font-mono-nums">
                  {summary?.kemungkinan_naik?.[k] != null ? `${summary.kemungkinan_naik[k]}%` : "—"}
                </td>
              ))}
            </tr>
          </>
        )}
        {rows.map((r) => (
          <tr
            key={r.ticker}
            onClick={() => onSelectTicker(r.ticker)}
            className="border-b border-[var(--md-sys-color-outline-variant)]/40 hover:bg-[var(--md-sys-color-surface-container-high)]/60 cursor-pointer transition-colors"
          >
            <td className="px-3 py-2 font-mono-nums font-extrabold text-[var(--md-sys-color-primary)]">{r.ticker}</td>
            <td className="px-3 py-2 whitespace-nowrap font-mono-nums">{r.listing_date || "—"}</td>
            <td className="px-3 py-2 text-right font-mono-nums">
              {r.penjatahan != null ? `${Number(r.penjatahan.toFixed(2))}%` : "—"}
            </td>
            {[0, 1, 2, 3, 4, 5, 6, 7, 8, 9].map((idx) => (
              <td key={idx} className="px-3 py-2 text-right font-mono-nums">
                {fmtPerf(r.day_returns?.[idx])}
              </td>
            ))}
            <td className="px-3 py-2 text-right font-mono-nums">{fmtPerf(r.perf_2w)}</td>
            <td className="px-3 py-2 text-right font-mono-nums">{fmtPerf(r.perf_1m)}</td>
            <td className="px-3 py-2 text-right font-mono-nums">{fmtPerf(r.perf_3m)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
