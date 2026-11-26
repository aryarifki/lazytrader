import type { IpoScreenerData } from "./ipoTypes";
import { SCR_METRICS, fmtScreenerVal, type ScreenerMetricKey } from "./ipoScreenerConstants";

interface IpoScreenerTableProps {
  res: IpoScreenerData | null;
  loading: boolean;
  cols: ScreenerMetricKey[];
  after: boolean;
  activeCol: ScreenerMetricKey;
  currentSk: string;
  sortDir: number;
  onSortCol: (col: ScreenerMetricKey, isSorted: boolean) => void;
  onOpenTicker: (ticker: string, opts: { after: boolean; lo?: string; hi?: string }) => void;
}

export default function IpoScreenerTable({
  res,
  loading,
  cols,
  after,
  activeCol,
  currentSk,
  sortDir,
  onSortCol,
  onOpenTicker,
}: IpoScreenerTableProps) {
  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] shadow-sm overflow-hidden">
      <div className="flex items-center justify-between px-4 py-3 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <h3 className="text-xs font-extrabold uppercase tracking-wider">IPO SCREENER RESULT</h3>
        <span className="text-xs font-mono-nums text-[var(--md-sys-color-on-surface-variant)]">
          {loading ? "Menghitung…" : <><b>{res?.total ?? 0}</b> saham</>}
        </span>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-xs border-collapse">
          <thead>
            <tr className="bg-[var(--md-sys-color-surface-container-high)]/40 border-b border-[var(--md-sys-color-outline-variant)]">
              <th className="px-3 py-2.5 text-left text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] w-10">
                #
              </th>
              <th className="px-3 py-2.5 text-left text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">
                Saham
              </th>
              {cols.map((k) => {
                const m = SCR_METRICS.find((x) => x.k === k)!;
                const lbl = k === "chg" && after ? "vs IPO" : m.label;
                const isSorted = currentSk === k;
                return (
                  <th
                    key={k}
                    title={m.desc}
                    onClick={() => onSortCol(k, isSorted)}
                    className={`px-3 py-2.5 text-right text-[10px] font-bold uppercase cursor-pointer select-none whitespace-nowrap ${
                      isSorted ? "text-[var(--md-sys-color-primary)]" : "text-[var(--md-sys-color-on-surface-variant)]"
                    }`}
                  >
                    {lbl}
                    {isSorted ? (sortDir === -1 ? " ▼" : " ▲") : ""}
                  </th>
                );
              })}
            </tr>
          </thead>
          <tbody>
            {!res?.rows?.length ? (
              <tr>
                <td colSpan={cols.length + 2} className="p-8 text-center text-[var(--md-sys-color-on-surface-variant)]">
                  {loading ? "Menghitung hasil screener…" : "Tidak ada saham yang memenuhi kriteria."}
                </td>
              </tr>
            ) : (
              res.rows.map((r, idx) => (
                <tr
                  key={r.kode}
                  onClick={() => onOpenTicker(r.kode, { after, lo: res.lo ?? undefined, hi: res.hi ?? undefined })}
                  className="border-b border-[var(--md-sys-color-outline-variant)]/40 hover:bg-[var(--md-sys-color-surface-container-high)]/60 cursor-pointer transition-colors"
                >
                  <td className="px-3 py-2 font-mono-nums text-[var(--md-sys-color-on-surface-variant)]">{idx + 1}</td>
                  <td className="px-3 py-2">
                    <span className="font-mono-nums font-extrabold text-[var(--md-sys-color-primary)]">{r.kode}</span>
                    {r.company_name && (
                      <span className="ml-2 text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
                        {r.company_name}
                      </span>
                    )}
                  </td>
                  {cols.map((k) => {
                    const v = r[k];
                    const m = SCR_METRICS.find((x) => x.k === k)!;
                    const isPct = m.kind === "pct";
                    const colorCls =
                      isPct && v != null
                        ? v > 0
                          ? "text-emerald-600 dark:text-emerald-400"
                          : v < 0
                          ? "text-rose-600 dark:text-rose-400"
                          : ""
                        : "";
                    return (
                      <td
                        key={k}
                        className={`px-3 py-2 text-right font-mono-nums font-semibold ${colorCls} ${
                          k === activeCol ? "bg-[var(--md-sys-color-primary)]/5" : ""
                        }`}
                      >
                        {fmtScreenerVal(k, v)}
                      </td>
                    );
                  })}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
