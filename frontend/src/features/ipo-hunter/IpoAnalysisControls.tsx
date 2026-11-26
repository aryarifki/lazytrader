import { useMemo, useState } from "react";
import type { IpoAnalysisMode, IpoInfoCard, IpoMasterTicker, IpoMeta } from "./ipoTypes";
import { fmtLot } from "./ipoFormat";

interface IpoAnalysisControlsProps {
  meta: IpoMeta | null;
  selectedTicker: string;
  tickerInput: string;
  onTickerInputChange: (val: string) => void;
  onSelectTicker: (ticker: string) => void;
  mode: IpoAnalysisMode;
  onModeChange: (mode: IpoAnalysisMode) => void;
  dateFrom: string;
  onDateFromChange: (val: string) => void;
  dateTo: string;
  onDateToChange: (val: string) => void;
  infoCard?: IpoInfoCard | null;
  tickerError: string;
  onOpenList: () => void;
}

export default function IpoAnalysisControls({
  meta,
  selectedTicker,
  tickerInput,
  onTickerInputChange,
  onSelectTicker,
  mode,
  onModeChange,
  dateFrom,
  onDateFromChange,
  dateTo,
  onDateToChange,
  infoCard,
  tickerError,
  onOpenList,
}: IpoAnalysisControlsProps) {
  const [suggestOpen, setSuggestOpen] = useState(false);
  const dates = meta?.available_dates || [];

  const suggestions = useMemo<IpoMasterTicker[]>(() => {
    const q = tickerInput.trim().toUpperCase();
    if (!q || !meta?.tickers) return [];
    return meta.tickers
      .filter((t) => t.ticker.includes(q) || (t.company_name || "").toUpperCase().includes(q))
      .slice(0, 8);
  }, [tickerInput, meta]);

  const handleSelect = (tk: string) => {
    setSuggestOpen(false);
    onSelectTicker(tk);
  };

  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-3.5 shadow-sm space-y-3">
      {/* IPO Terbaru Chips */}
      {meta?.quick_picks && meta.quick_picks.length > 0 && (
        <div className="flex items-center gap-1.5 overflow-x-auto pb-1 no-scrollbar">
          <span className="text-[10px] font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)] whitespace-nowrap mr-1">
            IPO Terbaru
          </span>
          {meta.quick_picks.slice(0, 18).map((qp) => (
            <button
              key={qp.ticker}
              type="button"
              title={`${qp.company_name} · IPO ${qp.listing_date}`}
              onClick={() => handleSelect(qp.ticker)}
              className={`px-2.5 py-1 rounded-md text-[11px] font-mono-nums font-bold border whitespace-nowrap transition cursor-pointer ${
                selectedTicker === qp.ticker
                  ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] border-[var(--md-sys-color-primary)]"
                  : "bg-[var(--md-sys-color-surface-container-high)] border-[var(--md-sys-color-outline-variant)] hover:border-[var(--md-sys-color-primary)]"
              }`}
            >
              {qp.ticker}
            </button>
          ))}
        </div>
      )}

      {/* Controls Row */}
      <div className="flex flex-wrap items-end gap-3">
        {/* Ticker Input */}
        <div className="relative">
          <label className="block text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] mb-1">
            Input Kode IPO
          </label>
          <div className="flex items-center gap-1">
            <input
              type="text"
              maxLength={6}
              value={tickerInput}
              onFocus={() => setSuggestOpen(true)}
              onBlur={() => setTimeout(() => setSuggestOpen(false), 180)}
              onChange={(e) => {
                onTickerInputChange(e.target.value.toUpperCase());
                setSuggestOpen(true);
              }}
              onKeyDown={(e) => {
                if (e.key === "Enter") handleSelect(tickerInput);
              }}
              placeholder="KODE IPO"
              className="h-8 w-28 px-2.5 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] font-mono-nums font-bold text-xs uppercase"
            />
            <button
              type="button"
              onClick={() => handleSelect(tickerInput)}
              className="h-8 px-2.5 rounded-lg bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] text-xs font-bold cursor-pointer"
            >
              Cari
            </button>
          </div>
          {suggestOpen && suggestions.length > 0 && (
            <div className="absolute left-0 top-full mt-1 w-64 z-40 rounded-xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] shadow-xl max-h-56 overflow-y-auto">
              {suggestions.map((s) => (
                <button
                  key={s.ticker}
                  type="button"
                  onMouseDown={() => handleSelect(s.ticker)}
                  className="w-full text-left px-3 py-1.5 hover:bg-[var(--md-sys-color-surface-container-high)] flex items-center justify-between text-xs border-b border-[var(--md-sys-color-outline-variant)]/30 last:border-0 cursor-pointer"
                >
                  <span className="truncate">
                    <b className="font-mono-nums text-[var(--md-sys-color-primary)]">{s.ticker}</b>{" "}
                    <span className="text-[10.5px] text-[var(--md-sys-color-on-surface-variant)] truncate">
                      {s.company_name}
                    </span>
                  </span>
                  <span className="text-[10px] font-mono-nums text-[var(--md-sys-color-on-surface-variant)] ml-2 shrink-0">
                    {s.listing_date}
                  </span>
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Start Mode */}
        <div>
          <label className="block text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] mb-1">
            Start
          </label>
          <div className="flex items-center gap-1.5">
            <button
              type="button"
              onClick={() => onModeChange("after")}
              className={`h-8 px-2.5 rounded-lg text-xs font-bold border transition cursor-pointer ${
                mode === "after"
                  ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border-[var(--md-sys-color-primary)]/40"
                  : "bg-[var(--md-sys-color-surface-container-high)] border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface-variant)]"
              }`}
            >
              IPO Date
            </button>
            <select
              value={mode === "periodic" ? dateFrom : ""}
              onChange={(e) => {
                if (e.target.value) {
                  onModeChange("periodic");
                  onDateFromChange(e.target.value);
                }
              }}
              className="h-8 px-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs font-mono-nums"
            >
              <option value="">Manual Start…</option>
              {dates.map((d) => (
                <option key={d} value={d}>
                  {d}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* End Date */}
        <div>
          <label className="block text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] mb-1">
            End
          </label>
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

        {/* Last Price & Change */}
        {infoCard && (
          <div className="ml-auto flex items-center gap-4">
            <div className="text-right">
              <div className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">
                Last Price
              </div>
              <div className="text-base font-mono-nums font-black">
                {fmtLot(infoCard.last_close)}
                {infoCard.ipo_vs_close_pct != null && (
                  <span
                    className={`ml-1.5 text-xs ${
                      infoCard.ipo_vs_close_pct >= 0
                        ? "text-emerald-600 dark:text-emerald-400"
                        : "text-rose-600 dark:text-rose-400"
                    }`}
                  >
                    ({infoCard.ipo_vs_close_pct >= 0 ? "+" : ""}
                    {infoCard.ipo_vs_close_pct}% vs IPO)
                  </span>
                )}
              </div>
            </div>
          </div>
        )}
      </div>

      {tickerError && (
        <div className="text-xs text-rose-500 font-semibold flex items-center gap-2">
          <span>{tickerError}</span>
          <button type="button" onClick={onOpenList} className="underline cursor-pointer">
            Lihat IPO List
          </button>
        </div>
      )}
    </div>
  );
}
