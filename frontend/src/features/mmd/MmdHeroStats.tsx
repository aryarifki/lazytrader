import { useMemo } from "react";
import { bandarBadgeClass } from "./mmdBadges";
import type { MmdAnalysis, MmdTicker } from "./mmdTypes";

interface MmdHeroStatsProps {
  tickerInput: string;
  setTickerInput: (val: string) => void;
  suggestOpen: boolean;
  setSuggestOpen: (open: boolean) => void;
  tickers: MmdTicker[];
  onSelectTicker: (ticker: string) => void;
  analysisData: MmdAnalysis | null;
  dateFrom: string;
  dateTo: string;
}

export default function MmdHeroStats({
  tickerInput,
  setTickerInput,
  suggestOpen,
  setSuggestOpen,
  tickers,
  onSelectTicker,
  analysisData,
  dateFrom,
  dateTo,
}: MmdHeroStatsProps) {
  const filteredSuggestions = useMemo(() => {
    if (!tickerInput) return [];
    const term = tickerInput.toUpperCase();
    return tickers.filter((t) => t.ticker.toUpperCase().includes(term)).slice(0, 8);
  }, [tickers, tickerInput]);

  const badgeClass = analysisData?.bandar_action?.cls || "undetected";
  const priceChange = analysisData?.price_change || 0;
  const korelasi = analysisData?.korelasi;
  const aiScore = analysisData?.ai_score;

  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-3.5 shadow-2xs space-y-3">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div className="flex items-center gap-3 flex-wrap">
          <div className="relative">
            <div className="flex items-center gap-1.5 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded-xl px-2.5 py-1.5 shadow-2xs">
              <input
                type="text"
                value={tickerInput}
                onChange={(e) => {
                  setTickerInput(e.target.value.toUpperCase());
                  setSuggestOpen(true);
                }}
                onKeyDown={(e) => {
                  if (e.key === "Enter" && tickerInput.trim()) {
                    onSelectTicker(tickerInput.trim().toUpperCase());
                  } else if (e.key === "Escape") {
                    setSuggestOpen(false);
                  }
                }}
                onFocus={() => setSuggestOpen(true)}
                onBlur={() => {
                  setTimeout(() => setSuggestOpen(false), 150);
                }}
                placeholder="CARI SAHAM..."
                className="w-24 bg-transparent border-none text-xs font-bold font-mono text-[var(--md-sys-color-on-surface)] tracking-wider focus:outline-none uppercase"
              />
              <button
                type="button"
                onClick={() => {
                  if (tickerInput.trim()) onSelectTicker(tickerInput.trim().toUpperCase());
                }}
                className="px-2 py-0.5 rounded-lg bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] text-[10px] font-mono font-bold hover:opacity-90 transition cursor-pointer"
              >
                GO
              </button>
            </div>
            {suggestOpen && filteredSuggestions.length > 0 && (
              <div className="absolute top-full left-0 mt-1 w-48 bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)] rounded-xl shadow-2xl z-50 max-h-48 overflow-y-auto p-1 backdrop-blur-md">
                {filteredSuggestions.map((t) => (
                  <button
                    key={t.ticker}
                    type="button"
                    onMouseDown={(e) => {
                      e.preventDefault();
                      onSelectTicker(t.ticker);
                    }}
                    className="w-full text-left px-3 py-1.5 hover:bg-[var(--md-sys-color-primary)] hover:text-[var(--md-sys-color-on-primary)] rounded-lg cursor-pointer text-xs font-mono font-bold flex items-center justify-between transition"
                  >
                    <span>{t.ticker}</span>
                    <span className="text-[10px] opacity-70 truncate max-w-[84px]">
                      {t.company_name || t.name || ""}
                    </span>
                  </button>
                ))}
              </div>
            )}
          </div>

          <div>
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-semibold uppercase tracking-wider">
              Last Price
            </div>
            <div className="text-sm sm:text-base font-mono-nums font-extrabold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
              <span>
                Rp {analysisData?.last_price ? Number(analysisData.last_price).toLocaleString("id-ID") : "—"}
              </span>
              <span
                className={`text-xs font-bold ${
                  priceChange >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                }`}
              >
                {priceChange > 0 ? `+${priceChange}%` : `${priceChange}%`}
              </span>
            </div>
          </div>

          <div className="ml-auto sm:ml-2">
            <span className={`${bandarBadgeClass(badgeClass)} text-[11px] py-1 px-2.5 min-w-[96px]`}>
              {analysisData?.bandar_action?.label || "—"}
            </span>
          </div>
        </div>

        <div className="flex items-center gap-2 text-xs font-mono">
          <div className="bg-[var(--md-sys-color-surface-container-high)] px-3 py-1.5 rounded-xl border border-[var(--md-sys-color-outline-variant)]">
            <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] block leading-tight">
              Periode:
            </span>
            <span className="font-bold text-[var(--md-sys-color-on-surface)] text-[11px]">
              {analysisData?.date_from || dateFrom || "—"} ~ {analysisData?.date_to || dateTo || "—"}
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-2 pt-2 border-t border-[var(--md-sys-color-outline-variant)]/60 text-xs font-mono-nums">
        <div className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl text-center border border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[9px] text-[var(--color-positive)] font-bold">TOTAL BUY</div>
          <div className="text-[var(--md-sys-color-on-surface)] font-bold text-xs">
            {analysisData?.total_buy?.toLocaleString("id-ID") || "0"}
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl text-center border border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[9px] text-[var(--color-negative)] font-bold">TOTAL SELL</div>
          <div className="text-[var(--md-sys-color-on-surface)] font-bold text-xs">
            {analysisData?.total_sell?.toLocaleString("id-ID") || "0"}
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl text-center border border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[9px] text-[var(--md-sys-color-primary)] font-bold">NET VOLUME</div>
          <div className="text-[var(--md-sys-color-on-surface)] font-bold text-xs">
            {analysisData?.tot_net_volume?.toLocaleString("id-ID") || "0"}
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl text-center border border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-bold">HARI BURSA</div>
          <div className="text-[var(--md-sys-color-on-surface)] font-bold text-xs">
            {analysisData?.total_days || 0}
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl text-center border border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-bold">KORELASI</div>
          <div
            className={`font-bold text-xs ${
              (korelasi || 0) >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
            }`}
          >
            {korelasi !== undefined ? (korelasi > 0 ? `+${korelasi}` : korelasi) : "—"}
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl text-center border border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-bold">AI SCORE</div>
          <div
            className={`font-bold text-xs ${
              (aiScore || 0) >= 70
                ? "text-[var(--color-positive)]"
                : (aiScore || 0) <= 40
                  ? "text-[var(--color-negative)]"
                  : "text-[var(--md-sys-color-primary)]"
            }`}
          >
            {aiScore !== undefined ? `${aiScore}/100` : "—"}
          </div>
        </div>
      </div>
    </div>
  );
}
