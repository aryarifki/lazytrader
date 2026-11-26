import type { IpoScreenerCondition } from "./ipoTypes";
import { SCR_METRICS, SCR_PRESETS, SCR_PRIMARY } from "./ipoScreenerConstants";

interface IpoScreenerOptionsProps {
  dates: string[];
  years: number[];
  quick: number | null;
  after: boolean;
  dateFrom: string;
  dateTo: string;
  year: string;
  chg: string;
  metric: string;
  conds: IpoScreenerCondition[];
  nDays?: number;
  hiDate?: string | null;
  activeDesc: string;
  onQuick: (n: number) => void;
  onToggleAfter: () => void;
  onDateFromChange: (val: string) => void;
  onDateToChange: (val: string) => void;
  onYearChange: (val: string) => void;
  onChgChange: (val: string) => void;
  onMetricChange: (val: string) => void;
  onCondsChange: (conds: IpoScreenerCondition[]) => void;
  onReset: () => void;
}

export default function IpoScreenerOptions({
  dates,
  years,
  quick,
  after,
  dateFrom,
  dateTo,
  year,
  chg,
  metric,
  conds,
  nDays,
  hiDate,
  activeDesc,
  onQuick,
  onToggleAfter,
  onDateFromChange,
  onDateToChange,
  onYearChange,
  onChgChange,
  onMetricChange,
  onCondsChange,
  onReset,
}: IpoScreenerOptionsProps) {
  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-4 shadow-sm max-w-4xl">
      <div className="flex items-center justify-between pb-3 mb-3 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <h2 className="text-xs font-extrabold uppercase tracking-wider">SCREENER OPTIONS</h2>
        <span className="text-xs font-mono-nums text-[var(--md-sys-color-on-surface-variant)]">
          {after ? `sejak IPO s/d ${hiDate || dateTo}` : `${nDays ?? quick ?? 5} hari`}
        </span>
      </div>

      <div className="space-y-3 text-xs">
        {/* Periode */}
        <div className="grid grid-cols-1 sm:grid-cols-[140px_1fr] gap-2 items-center">
          <span className="font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase text-[11px]">Periode</span>
          <div className="flex flex-wrap items-center gap-1.5">
            {[1, 5, 10, 20].map((n) => (
              <button
                key={n}
                type="button"
                onClick={() => onQuick(n)}
                className={`px-2.5 py-1 rounded-md text-xs font-bold border transition cursor-pointer ${
                  !after && quick === n
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] border-[var(--md-sys-color-primary)]"
                    : "bg-[var(--md-sys-color-surface-container-high)] border-[var(--md-sys-color-outline-variant)]"
                }`}
              >
                {n} hari
              </button>
            ))}
            <button
              type="button"
              onClick={onToggleAfter}
              className={`px-2.5 py-1 rounded-md text-xs font-bold border transition cursor-pointer ${
                after
                  ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border-[var(--md-sys-color-primary)]/40"
                  : "bg-[var(--md-sys-color-surface-container-high)] border-[var(--md-sys-color-outline-variant)]"
              }`}
            >
              After IPO
            </button>
            <select
              value={dateFrom}
              disabled={after}
              onChange={(e) => onDateFromChange(e.target.value)}
              className="h-8 px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs font-mono-nums disabled:opacity-40"
            >
              {dates.map((d) => (
                <option key={d} value={d}>
                  {d}
                </option>
              ))}
            </select>
            <span>–</span>
            <select
              value={dateTo}
              onChange={(e) => onDateToChange(e.target.value)}
              className="h-8 px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs font-mono-nums"
            >
              {dates.map((d) => (
                <option key={d} value={d}>
                  {d}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Tahun IPO */}
        <div className="grid grid-cols-1 sm:grid-cols-[140px_1fr] gap-2 items-center">
          <span className="font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase text-[11px]">Tahun IPO</span>
          <select
            value={year}
            onChange={(e) => onYearChange(e.target.value)}
            className="h-8 w-fit px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs"
          >
            <option value="">Semua</option>
            {years.map((y) => (
              <option key={y} value={y}>
                IPO {y}
              </option>
            ))}
          </select>
        </div>

        {/* Perubahan Harga */}
        <div className="grid grid-cols-1 sm:grid-cols-[140px_1fr] gap-2 items-center">
          <span className="font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase text-[11px]">
            Perubahan Harga
          </span>
          <select
            value={chg}
            onChange={(e) => onChgChange(e.target.value)}
            className="h-8 w-fit px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs"
          >
            <option value="all">Semua</option>
            <option value="up">Naik</option>
            <option value="down">Turun</option>
          </select>
        </div>

        {/* Jenis Filter */}
        <div className="grid grid-cols-1 sm:grid-cols-[140px_1fr] gap-2 items-center">
          <span className="font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase text-[11px]">
            Jenis Filter
          </span>
          <div className="flex flex-wrap items-center gap-2">
            <select
              value={metric}
              onChange={(e) => onMetricChange(e.target.value)}
              className="h-8 px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs font-semibold"
            >
              <optgroup label="Peringkat">
                {SCR_PRESETS.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.label}
                  </option>
                ))}
              </optgroup>
              <optgroup label="Aliran broker & harga">
                {SCR_PRIMARY.map((k) => {
                  const m = SCR_METRICS.find((x) => x.k === k)!;
                  return (
                    <option key={k} value={k}>
                      {m.label}
                    </option>
                  );
                })}
              </optgroup>
            </select>
            <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">{activeDesc}</span>
          </div>
        </div>

        {/* Kondisi Tambahan */}
        <div className="grid grid-cols-1 sm:grid-cols-[140px_1fr] gap-2 items-start">
          <span className="font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase text-[11px] pt-1.5">
            Kondisi Tambahan
          </span>
          <div className="space-y-2">
            {conds.map((c, idx) => {
              const m = SCR_METRICS.find((x) => x.k === c.k);
              const unit = m?.kind === "num" ? "miliar" : m?.kind === "x" ? "x" : "%";
              return (
                <div key={idx} className="flex flex-wrap items-center gap-1.5">
                  <select
                    value={c.k}
                    onChange={(e) => {
                      const next = [...conds];
                      next[idx] = { ...c, k: e.target.value };
                      onCondsChange(next);
                    }}
                    className="h-8 px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs"
                  >
                    {SCR_METRICS.map((mm) => (
                      <option key={mm.k} value={mm.k}>
                        {mm.label}
                      </option>
                    ))}
                  </select>
                  <select
                    value={c.op}
                    onChange={(e) => {
                      const next = [...conds];
                      next[idx] = { ...c, op: e.target.value };
                      onCondsChange(next);
                    }}
                    className="h-8 px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs font-mono-nums"
                  >
                    {[">", ">=", "<", "<="].map((op) => (
                      <option key={op} value={op}>
                        {op}
                      </option>
                    ))}
                  </select>
                  <input
                    type="text"
                    value={c.v}
                    placeholder="nilai"
                    onChange={(e) => {
                      const next = [...conds];
                      next[idx] = { ...c, v: e.target.value.replace(",", ".") };
                      onCondsChange(next);
                    }}
                    className="h-8 w-24 px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs font-mono-nums"
                  />
                  <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">{unit}</span>
                  <button
                    type="button"
                    onClick={() => onCondsChange(conds.filter((_, i) => i !== idx))}
                    className="w-7 h-7 rounded-lg text-rose-500 hover:bg-rose-500/10 font-bold cursor-pointer"
                  >
                    ×
                  </button>
                </div>
              );
            })}
            <div className="flex items-center gap-2">
              <button
                type="button"
                onClick={() => onCondsChange([...conds, { k: "ratrans", op: ">", v: "" }])}
                className="px-2.5 py-1 rounded-md text-xs font-semibold bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] hover:border-[var(--md-sys-color-primary)] cursor-pointer"
              >
                + Kondisi
              </button>
              <button
                type="button"
                onClick={onReset}
                className="px-2.5 py-1 rounded-md text-xs font-semibold bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] cursor-pointer"
              >
                Reset
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
