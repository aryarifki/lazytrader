import type { IpoCategoryCard as IpoCategoryCardData } from "./ipoTypes";
import { fmtLot, fmtPct, fmtSignedLot } from "./ipoFormat";

export default function IpoCategoryCard({ catCard }: { catCard?: IpoCategoryCardData | null }) {
  if (!catCard) return null;

  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-4 shadow-sm space-y-2">
      <h3 className="text-xs font-extrabold uppercase tracking-wider">TOTAL TRANSAKSI PER KATEGORI</h3>
      <div className="grid grid-cols-[1fr_70px_54px_70px_54px] text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] pb-1 border-b border-[var(--md-sys-color-outline-variant)]/50">
        <span>Kategori</span>
        <span className="text-right">Net</span>
        <span className="text-right">%Net</span>
        <span className="text-right">Trans</span>
        <span className="text-right">%Trans</span>
      </div>
      <div className="divide-y divide-[var(--md-sys-color-outline-variant)]/30 text-xs">
        {catCard.rows.map((r) => (
          <div key={r.key} className="grid grid-cols-[1fr_70px_54px_70px_54px] items-center py-1.5">
            <span className="font-semibold truncate">
              {r.label}
              {r.bandar_brokers?.length > 0 && (
                <span className="ml-1 text-[10px] opacity-65">({r.bandar_brokers.join(", ")})</span>
              )}
            </span>
            <span
              className={`text-right font-mono-nums ${
                r.net >= 0 ? "text-emerald-600 dark:text-emerald-400" : "text-rose-600 dark:text-rose-400"
              }`}
            >
              {fmtSignedLot(r.net)}
            </span>
            <span
              className={`text-right font-mono-nums ${
                (r.pct_net ?? 0) >= 0 ? "text-emerald-600 dark:text-emerald-400" : "text-rose-600 dark:text-rose-400"
              }`}
            >
              {fmtPct(r.pct_net)}
            </span>
            <span className="text-right font-mono-nums">{fmtLot(r.trans)}</span>
            <span className="text-right font-mono-nums">{r.pct_trans != null ? `${r.pct_trans.toFixed(1)}%` : "—"}</span>
          </div>
        ))}
        <div className="grid grid-cols-[1fr_70px_54px_70px_54px] items-center py-1.5 font-extrabold border-t border-[var(--md-sys-color-outline-variant)]">
          <span>Total</span>
          <span className="text-right font-mono-nums">{fmtSignedLot(catCard.total.net)}</span>
          <span className="text-right font-mono-nums">{fmtPct(catCard.total.pct_net)}</span>
          <span className="text-right font-mono-nums">{fmtLot(catCard.total.trans)}</span>
          <span className="text-right font-mono-nums">
            {catCard.total.pct_trans != null ? `${catCard.total.pct_trans.toFixed(1)}%` : "—"}
          </span>
        </div>
      </div>
    </div>
  );
}
