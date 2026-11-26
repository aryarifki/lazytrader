import type { IpoUnderwriterCard as IpoUnderwriterCardData } from "./ipoTypes";
import { fmtLot, fmtPct, fmtSignedLot } from "./ipoFormat";

export default function IpoUnderwriterCard({
  uwCard,
  ticker,
}: {
  uwCard?: IpoUnderwriterCardData | null;
  ticker: string;
}) {
  if (!uwCard) return null;

  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-4 shadow-sm space-y-3">
      <h3 className="text-xs font-extrabold uppercase tracking-wider">UNDERWRITER</h3>
      {!uwCard.has_uw ? (
        <div className="py-4 text-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
          Tidak ada data underwriter untuk {ticker}.
        </div>
      ) : (
        <>
          <div className="flex items-end justify-between gap-2 pb-3 border-b border-[var(--md-sys-color-outline-variant)]/50">
            <div>
              <div className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">
                Total net underwriter
              </div>
              <div className="flex items-baseline gap-2">
                <span
                  className={`text-2xl font-mono-nums font-black ${
                    uwCard.total_net >= 0
                      ? "text-emerald-600 dark:text-emerald-400"
                      : "text-rose-600 dark:text-rose-400"
                  }`}
                >
                  {fmtSignedLot(uwCard.total_net)}
                </span>
                <span
                  className={`text-sm font-mono-nums font-bold ${
                    (uwCard.total_pct_net ?? 0) >= 0
                      ? "text-emerald-600 dark:text-emerald-400"
                      : "text-rose-600 dark:text-rose-400"
                  }`}
                >
                  {fmtPct(uwCard.total_pct_net)}
                </span>
              </div>
            </div>
            <div className="text-right">
              <span
                className={`inline-block px-2.5 py-0.5 rounded-md text-[10px] font-extrabold border ${
                  uwCard.total_net >= 0
                    ? "border-emerald-500/50 text-emerald-600 dark:text-emerald-400"
                    : "border-rose-500/50 text-rose-600 dark:text-rose-400"
                }`}
              >
                {uwCard.status}
              </span>
              <div className="mt-1.5 text-[10px] uppercase text-[var(--md-sys-color-on-surface-variant)]">
                Total IPO Shares
              </div>
              <div className="text-xs font-mono-nums font-bold">
                {uwCard.lot_ipo ? fmtLot(uwCard.lot_ipo) : "—"}{" "}
                <span className="text-[var(--md-sys-color-on-surface-variant)]">{fmtPct(uwCard.pct_of_ipo)}</span>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-[1fr_72px_56px_72px_56px] text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">
            <span>Broker</span>
            <span className="text-right">Net</span>
            <span className="text-right">%Net</span>
            <span className="text-right">Trans</span>
            <span className="text-right">%Trans</span>
          </div>
          <div className="divide-y divide-[var(--md-sys-color-outline-variant)]/30">
            {uwCard.rows.map((u) => (
              <div
                key={u.broker}
                className="grid grid-cols-[1fr_72px_56px_72px_56px] items-center py-1.5 text-xs font-mono-nums"
              >
                <span className="flex items-center gap-1.5 font-extrabold">
                  {u.broker}
                  <span
                    className={`px-1.5 py-0.5 rounded text-[9px] font-bold ${
                      u.role === "LEAD"
                        ? "bg-amber-500/20 text-amber-600 dark:text-amber-400"
                        : "bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)]"
                    }`}
                  >
                    {u.role}
                  </span>
                </span>
                <span
                  className={`text-right font-semibold ${
                    u.net >= 0 ? "text-emerald-600 dark:text-emerald-400" : "text-rose-600 dark:text-rose-400"
                  }`}
                >
                  {fmtSignedLot(u.net)}
                </span>
                <span
                  className={`text-right ${
                    (u.pct_net ?? 0) >= 0
                      ? "text-emerald-600 dark:text-emerald-400"
                      : "text-rose-600 dark:text-rose-400"
                  }`}
                >
                  {fmtPct(u.pct_net)}
                </span>
                <span className="text-right">{fmtLot(u.trans)}</span>
                <span className="text-right">{u.pct_trans != null ? `${u.pct_trans.toFixed(1)}%` : "—"}</span>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
