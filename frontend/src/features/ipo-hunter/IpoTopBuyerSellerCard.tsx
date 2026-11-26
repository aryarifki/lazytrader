import type { IpoTopBuyerSellerCard } from "./ipoTypes";
import { fmtLot, tbsColorClass } from "./ipoFormat";

export default function IpoTopBuyerSellerCard({ tbs }: { tbs?: IpoTopBuyerSellerCard | null }) {
  if (!tbs) return null;

  const buyers = tbs.buyers || [];
  const sellers = tbs.sellers || [];
  const rankCount = Math.min(10, Math.max(buyers.length, sellers.length));

  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-4 shadow-sm overflow-x-auto">
      <div className="grid grid-cols-1 md:grid-cols-[1fr_28px_1fr] gap-2 min-w-[540px]">
        {/* Top Buyer */}
        <div>
          <div className="text-xs font-extrabold uppercase text-emerald-600 dark:text-emerald-400 mb-1.5">
            Top Buyer
          </div>
          <div className="grid grid-cols-[52px_1fr_68px_68px] text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] pb-1 border-b border-[var(--md-sys-color-outline-variant)]/50">
            <span>Broker</span>
            <span>Kategori</span>
            <span className="text-right">Volume</span>
            <span className="text-right">Value (jt)</span>
          </div>
          {buyers.map((r) => (
            <div
              key={r.broker}
              className={`grid grid-cols-[52px_1fr_68px_68px] items-center py-1 text-xs border-b border-[var(--md-sys-color-outline-variant)]/25 ${
                r.is_uw ? "bg-rose-500/5" : ""
              }`}
            >
              <span className={`font-mono-nums ${tbsColorClass(r.style_key)}`}>{r.broker}</span>
              <span className="text-[11px] truncate">{r.category_label}</span>
              <span className="text-right font-mono-nums">{fmtLot(r.net)}</span>
              <span className="text-right font-mono-nums text-[var(--md-sys-color-on-surface-variant)]">
                {r.value_jt != null ? fmtLot(r.value_jt) : "—"}
              </span>
            </div>
          ))}
        </div>

        {/* Rank # */}
        <div className="hidden md:flex flex-col items-center">
          <div className="text-xs font-extrabold mb-1.5">&nbsp;</div>
          <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] pb-1 border-b border-[var(--md-sys-color-outline-variant)]/50 w-full text-center">
            #
          </div>
          {Array.from({ length: rankCount }).map((_, i) => (
            <div
              key={i}
              className="py-1 text-[11px] font-mono-nums text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/25 w-full text-center"
            >
              {i + 1}
            </div>
          ))}
        </div>

        {/* Top Seller */}
        <div>
          <div className="text-xs font-extrabold uppercase text-rose-600 dark:text-rose-400 mb-1.5">
            Top Seller
          </div>
          <div className="grid grid-cols-[52px_1fr_68px_68px] text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] pb-1 border-b border-[var(--md-sys-color-outline-variant)]/50">
            <span>Broker</span>
            <span>Kategori</span>
            <span className="text-right">Volume</span>
            <span className="text-right">Value (jt)</span>
          </div>
          {sellers.map((r) => (
            <div
              key={r.broker}
              className={`grid grid-cols-[52px_1fr_68px_68px] items-center py-1 text-xs border-b border-[var(--md-sys-color-outline-variant)]/25 ${
                r.is_uw ? "bg-rose-500/5" : ""
              }`}
            >
              <span className={`font-mono-nums ${tbsColorClass(r.style_key)}`}>{r.broker}</span>
              <span className="text-[11px] truncate">{r.category_label}</span>
              <span className="text-right font-mono-nums">{fmtLot(r.net)}</span>
              <span className="text-right font-mono-nums text-[var(--md-sys-color-on-surface-variant)]">
                {r.value_jt != null ? fmtLot(r.value_jt) : "—"}
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
