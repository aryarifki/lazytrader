import type { AkumulasiResult, ForeignPowerInfo, StockSeriesResponse } from "./foreignProTypes";

interface ForeignProStockHeroProps {
  ticker: string;
  stockData: StockSeriesResponse | null;
  latestFp: ForeignPowerInfo;
  akumDays: number;
  akumulasiResult: AkumulasiResult;
  valFFText: string;
  latestCorr: number;
}

export default function ForeignProStockHero({
  ticker,
  stockData,
  latestFp,
  akumDays,
  akumulasiResult,
  valFFText,
  latestCorr,
}: ForeignProStockHeroProps) {
  const priceChange = stockData?.analytics?.price_change || 0;

  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-3.5 shadow-2xs space-y-3 font-mono-nums">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <div className="bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] px-3.5 py-1.5 rounded-xl font-mono font-black text-lg tracking-wider shadow-2xs flex items-center gap-1">
            <span>{ticker}</span>
          </div>
          <div>
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-semibold uppercase tracking-wider">
              Last Price
            </div>
            <div className="text-sm sm:text-base font-extrabold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
              <span>Rp {stockData?.analytics?.last_price?.toLocaleString("id-ID") || "—"}</span>
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
            <span
              className="flex items-center gap-1.5 px-3 py-1 rounded-xl border font-bold text-xs bg-[var(--md-sys-color-surface-container-high)]"
              style={{
                borderColor: latestFp.color,
                color: latestFp.color,
              }}
            >
              <span>{latestFp.badge}</span>
              <span className="opacity-40">|</span>
              <span>{latestFp.score} pts</span>
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 pt-2 border-t border-[var(--md-sys-color-outline-variant)]/60 text-xs">
        <div className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl text-center border border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-bold">
            FLOW HARI INI
          </div>
          <div
            className={`font-bold text-xs ${
              valFFText.includes("+") ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
            }`}
          >
            {valFFText || "—"}
          </div>
        </div>
        <div className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl text-center border border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-bold">
            {akumulasiResult.isAcc ? `AVG AKUMULASI (${akumDays}D)` : `AVG DISTRIBUSI (${akumDays}D)`}
          </div>
          <div
            className={`font-bold text-xs ${
              akumulasiResult.isAcc ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
            }`}
          >
            Rp {akumulasiResult.avgPrice?.toLocaleString("id-ID") || "—"}
          </div>
        </div>
        <div className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl text-center border border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-bold">
            TOTAL FLOW ({akumDays}D)
          </div>
          <div
            className={`font-bold text-xs ${
              akumulasiResult.totalFlow >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
            }`}
          >
            {akumulasiResult.totalFlow >= 0 ? "+" : ""}
            {akumulasiResult.totalFlow.toFixed(1)} M
          </div>
        </div>
        <div className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl text-center border border-[var(--md-sys-color-outline-variant)]/40">
          <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-bold">
            KORELASI FF (20D)
          </div>
          <div
            className={`font-bold text-xs ${
              latestCorr >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
            }`}
          >
            {latestCorr >= 0 ? `+${latestCorr}` : latestCorr}
          </div>
        </div>
      </div>
    </div>
  );
}
