import type { IpoRetailActionCard as IpoRetailActionCardData } from "./ipoTypes";
import { fmtLot, fmtPct, fmtSignedLot, retailBadgeClass } from "./ipoFormat";

export default function IpoRetailActionCard({
  raCard,
  autoRetail,
  onToggleAutoRetail,
  retailOverride,
  onRetailOverrideChange,
}: {
  raCard?: IpoRetailActionCardData | null;
  autoRetail: boolean;
  onToggleAutoRetail: () => void;
  retailOverride: string;
  onRetailOverrideChange: (val: string) => void;
}) {
  if (!raCard) return null;

  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-4 shadow-sm space-y-3">
      <div className="flex items-center justify-between">
        <h3 className="text-xs font-extrabold uppercase tracking-wider">RETAIL ACTION</h3>
        <div className="flex items-center gap-2">
          <span
            className={`text-xs font-mono-nums font-bold ${
              (raCard.retail_ratio ?? 0) < 0
                ? "text-rose-600 dark:text-rose-400"
                : "text-emerald-600 dark:text-emerald-400"
            }`}
          >
            {fmtPct(raCard.retail_ratio)}
          </span>
          <span
            className={`px-2 py-0.5 rounded-md text-[10.5px] font-bold border ${retailBadgeClass(
              raCard.retail_action?.cls
            )}`}
          >
            {raCard.retail_action?.label || "—"}
          </span>
        </div>
      </div>

      {/* Retail Dipakai Bandar controls */}
      <div className="flex flex-wrap items-center justify-between gap-2 p-2.5 rounded-xl bg-[var(--md-sys-color-surface-container-high)]/60 border border-[var(--md-sys-color-outline-variant)]/50 text-xs">
        <div className="flex items-center gap-2">
          <span className="text-[10.5px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">
            Ritel Dipakai Bandar:
          </span>
          <span className="font-mono-nums font-extrabold text-amber-600 dark:text-amber-400">
            {raCard.retail_bandar_auto_result || "Tidak terdeteksi"}
          </span>
          {raCard.retail_bandar_auto_result && (
            <button
              type="button"
              onClick={onToggleAutoRetail}
              className={`px-2 py-0.5 rounded text-[10px] font-bold border cursor-pointer ${
                autoRetail
                  ? "bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/30"
                  : "bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]"
              }`}
            >
              {autoRetail ? "AUTO ON" : "AUTO OFF"}
            </button>
          )}
        </div>

        <select
          value={retailOverride}
          onChange={(e) => onRetailOverrideChange(e.target.value)}
          className="h-7 px-2 rounded-md bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] text-[11px] font-mono-nums"
        >
          <option value="">— Manual —</option>
          {(raCard.effective_retail_list || [])
            .filter((b) => !(autoRetail && b === raCard.retail_bandar_auto_result))
            .sort()
            .map((b) => (
              <option key={b} value={b}>
                {b}
              </option>
            ))}
        </select>
      </div>

      <div className="grid grid-cols-[54px_70px_54px_70px_1fr] text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] pb-1 border-b border-[var(--md-sys-color-outline-variant)]/50">
        <span>Broker</span>
        <span className="text-right">Net</span>
        <span className="text-right">%Net</span>
        <span className="text-right">Trans</span>
        <span className="text-right">%Trans</span>
      </div>
      <div className="divide-y divide-[var(--md-sys-color-outline-variant)]/30 text-xs font-mono-nums">
        {raCard.rows.map((r) => (
          <div
            key={r.broker}
            className={`grid grid-cols-[54px_70px_54px_70px_1fr] items-center py-1 ${
              r.is_active_bandar ? "opacity-50" : ""
            }`}
          >
            <span className="font-bold">
              {r.broker}
              {r.is_active_bandar && <span className="ml-1 text-[9px] text-amber-500">★</span>}
            </span>
            <span
              className={`text-right ${
                r.net >= 0 ? "text-emerald-600 dark:text-emerald-400" : "text-rose-600 dark:text-rose-400"
              }`}
            >
              {fmtSignedLot(r.net)}
            </span>
            <span
              className={`text-right ${
                (r.pct_net ?? 0) >= 0 ? "text-emerald-600 dark:text-emerald-400" : "text-rose-600 dark:text-rose-400"
              }`}
            >
              {fmtPct(r.pct_net)}
            </span>
            <span className="text-right">{fmtLot(r.trans)}</span>
            <span className="text-right">{r.pct_trans != null ? `${r.pct_trans.toFixed(1)}%` : "—"}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
