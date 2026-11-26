import { useMemo } from "react";
import { Icon } from "@iconify/react";
import { brokerCategoryClass } from "./mmdBadges";
import type { MmdAnalysis, MmdBrokerItem, MobileBrokerTab } from "./mmdTypes";

interface MmdBrokerTableProps {
  analysisData: MmdAnalysis | null;
  mobileBrokerTab: MobileBrokerTab;
  setMobileBrokerTab: (tab: MobileBrokerTab) => void;
  onOpenBroker: (code: string) => void;
}

function getBrokerVol(item?: MmdBrokerItem): number {
  if (!item) return 0;
  return Math.abs(item.volume ?? item.net ?? 0);
}

export default function MmdBrokerTable({
  analysisData,
  mobileBrokerTab,
  setMobileBrokerTab,
  onOpenBroker,
}: MmdBrokerTableProps) {
  const maxBuyVol = useMemo(() => {
    if (!analysisData?.buyers?.length) return 1;
    return Math.max(1, ...analysisData.buyers.map((b) => getBrokerVol(b)));
  }, [analysisData?.buyers]);

  const maxSellVol = useMemo(() => {
    if (!analysisData?.sellers?.length) return 1;
    return Math.max(1, ...analysisData.sellers.map((s) => getBrokerVol(s)));
  }, [analysisData?.sellers]);

  const renderMobileSide = (side: "buyer" | "seller") => {
    const list = side === "buyer" ? analysisData?.buyers : analysisData?.sellers;
    const maxVol = side === "buyer" ? maxBuyVol : maxSellVol;
    const colorClass = side === "buyer" ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]";
    const barClass = side === "buyer" ? "bg-[var(--color-positive)]" : "bg-[var(--color-negative)]";

    return (
      <div className="space-y-1.5">
        {Array.from({ length: 10 }).map((_, i) => {
          const row = list?.[i];
          if (!row) return null;
          const volPct = Math.min(100, Math.round((getBrokerVol(row) / maxVol) * 100));
          return (
            <div
              key={i}
              onClick={() => onOpenBroker(row.broker)}
              className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl border border-[var(--md-sys-color-outline-variant)]/60 flex items-center justify-between text-xs font-mono-nums active:bg-[var(--md-sys-color-surface-container-highest)] cursor-pointer transition"
            >
              <div className="flex items-center gap-2">
                <span className="w-5 text-center text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)]">
                  #{i + 1}
                </span>
                <div>
                  <div className="flex items-center gap-1.5">
                    <span className={`font-extrabold text-sm ${colorClass}`}>{row.broker}</span>
                    <span className={brokerCategoryClass(row.kategori)}>{row.kategori}</span>
                  </div>
                  <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] mt-0.5">
                    Rp {row.value_jt ? row.value_jt.toLocaleString("id-ID") : 0} Jt
                  </div>
                </div>
              </div>
              <div className="text-right">
                <div className={`font-bold ${colorClass}`}>{row.volume_fmt || "0"}</div>
                <div className="w-16 bg-[var(--md-sys-color-surface-container-highest)] h-1.5 rounded-full mt-1 overflow-hidden ml-auto">
                  <div className={`${barClass} h-full`} style={{ width: `${volPct}%` }} />
                </div>
              </div>
            </div>
          );
        })}
      </div>
    );
  };

  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl overflow-hidden shadow-2xs">
      <div className="bg-[var(--md-sys-color-surface-container-high)] px-3.5 py-2.5 border-b border-[var(--md-sys-color-outline-variant)] flex items-center justify-between">
        <div>
          <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] tracking-wide">
            TOP BUYER VS TOP SELLER
          </span>
          <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
            Net Volume &amp; Value (Jt IDR)
          </div>
        </div>

        <div className="md:hidden flex items-center bg-[var(--md-sys-color-surface-container-highest)] p-0.5 rounded-lg border border-[var(--md-sys-color-outline-variant)]">
          {(["buyer", "seller", "accdist"] as const).map((tab) => {
            const label = tab === "buyer" ? "Buyer" : tab === "seller" ? "Seller" : "Net";
            const activeColor =
              tab === "buyer"
                ? "text-[var(--color-positive)]"
                : tab === "seller"
                  ? "text-[var(--color-negative)]"
                  : "text-[var(--md-sys-color-primary)]";
            return (
              <button
                key={tab}
                type="button"
                onClick={() => setMobileBrokerTab(tab)}
                className={`px-2 py-1 rounded text-[10px] font-bold transition cursor-pointer ${
                  mobileBrokerTab === tab
                    ? `bg-[var(--md-sys-color-surface)] ${activeColor} shadow-xs`
                    : "text-[var(--md-sys-color-on-surface-variant)]"
                }`}
              >
                {label}
              </button>
            );
          })}
        </div>
      </div>

      <div className="md:hidden p-2.5 space-y-2">
        {mobileBrokerTab === "buyer" && renderMobileSide("buyer")}
        {mobileBrokerTab === "seller" && renderMobileSide("seller")}
        {mobileBrokerTab === "accdist" && (
          <div className="space-y-1.5">
            {Array.from({ length: 10 }).map((_, i) => {
              const b = analysisData?.buyers?.[i];
              const s = analysisData?.sellers?.[i];
              const diff = analysisData?.acc_dist_rows?.[i];
              if (!b && !s) return null;
              return (
                <div
                  key={i}
                  className="bg-[var(--md-sys-color-surface-container-low)] p-2.5 rounded-xl border border-[var(--md-sys-color-outline-variant)]/60 text-xs font-mono-nums space-y-1.5"
                >
                  <div className="flex items-center justify-between text-[11px]">
                    <span className="font-bold text-[var(--md-sys-color-on-surface-variant)]">Row #{i + 1}</span>
                    <span
                      className={`font-extrabold text-xs ${
                        (diff || 0) >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                      }`}
                    >
                      Net: {diff !== undefined ? diff.toLocaleString("id-ID") : "—"}
                    </span>
                  </div>
                  <div className="grid grid-cols-2 gap-2 text-[11px] pt-1 border-t border-[var(--md-sys-color-outline-variant)]/40">
                    <div className="text-[var(--color-positive)]">
                      Buy: <b>{b?.broker || "—"}</b> ({b?.volume_fmt || "0"})
                    </div>
                    <div className="text-[var(--color-negative)] text-right">
                      Sell: <b>{s?.broker || "—"}</b> ({s?.volume_fmt || "0"})
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      <div className="hidden md:block overflow-x-auto">
        <table className="w-full text-left text-xs font-mono-nums">
          <thead>
            <tr className="bg-[var(--md-sys-color-surface-container-highest)]/50 text-[10px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]">
              <th className="py-2.5 px-3 text-[var(--color-positive)]">BUYER</th>
              <th className="py-2.5 px-2 text-right">VOL</th>
              <th className="py-2.5 px-2">KAT</th>
              <th className="py-2.5 px-2 text-right">VAL(JT)</th>
              <th className="py-2.5 px-2 text-center text-[var(--md-sys-color-outline)]">#</th>
              <th className="py-2.5 px-3 text-[var(--color-negative)]">SELLER</th>
              <th className="py-2.5 px-2 text-right">VOL</th>
              <th className="py-2.5 px-2">KAT</th>
              <th className="py-2.5 px-2 text-right">VAL(JT)</th>
              <th className="py-2.5 px-3 text-right text-[var(--md-sys-color-primary)]">ACC/DIST</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
            {Array.from({ length: 10 }).map((_, i) => {
              const b = analysisData?.buyers?.[i];
              const s = analysisData?.sellers?.[i];
              const diff = analysisData?.acc_dist_rows?.[i];
              return (
                <tr key={i} className="hover:bg-[var(--md-sys-color-surface-container-high)] transition">
                  <td className="py-1.5 px-3 font-bold text-[var(--md-sys-color-on-surface)]">
                    {b?.broker ? (
                      <button
                        type="button"
                        onClick={() => onOpenBroker(b.broker)}
                        className="hover:text-[var(--md-sys-color-primary)] hover:underline cursor-pointer focus:outline-none flex items-center gap-1 group text-left"
                        title={`Buka detail drilldown transaksi broker ${b.broker}`}
                      >
                        <span>{b.broker}</span>
                        <Icon
                          icon="ph:magnifying-glass"
                          width={11}
                          className="opacity-0 group-hover:opacity-100 text-[var(--md-sys-color-primary)] transition"
                        />
                      </button>
                    ) : (
                      "—"
                    )}
                  </td>
                  <td className="py-1.5 px-2 text-right text-[var(--color-positive)]">{b?.volume_fmt || "—"}</td>
                  <td className="py-1.5 px-2">
                    {b ? <span className={brokerCategoryClass(b.kategori)}>{b.kategori}</span> : "—"}
                  </td>
                  <td className="py-1.5 px-2 text-right text-[var(--md-sys-color-on-surface)]">
                    {b?.value_jt ? b.value_jt.toLocaleString("id-ID") : "—"}
                  </td>
                  <td className="py-1.5 px-2 text-center text-[10px] text-[var(--md-sys-color-outline)] font-bold">
                    {i + 1}
                  </td>
                  <td className="py-1.5 px-3 font-bold text-[var(--md-sys-color-on-surface)]">
                    {s?.broker ? (
                      <button
                        type="button"
                        onClick={() => onOpenBroker(s.broker)}
                        className="hover:text-[var(--md-sys-color-primary)] hover:underline cursor-pointer focus:outline-none flex items-center gap-1 group text-left"
                        title={`Buka detail drilldown transaksi broker ${s.broker}`}
                      >
                        <span>{s.broker}</span>
                        <Icon
                          icon="ph:magnifying-glass"
                          width={11}
                          className="opacity-0 group-hover:opacity-100 text-[var(--md-sys-color-primary)] transition"
                        />
                      </button>
                    ) : (
                      "—"
                    )}
                  </td>
                  <td className="py-1.5 px-2 text-right text-[var(--color-negative)]">{s?.volume_fmt || "—"}</td>
                  <td className="py-1.5 px-2">
                    {s ? <span className={brokerCategoryClass(s.kategori)}>{s.kategori}</span> : "—"}
                  </td>
                  <td className="py-1.5 px-2 text-right text-[var(--md-sys-color-on-surface)]">
                    {s?.value_jt ? s.value_jt.toLocaleString("id-ID") : "—"}
                  </td>
                  <td
                    className={`py-1.5 px-3 text-right font-bold ${
                      (diff || 0) >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                    }`}
                  >
                    {diff !== undefined ? diff.toLocaleString("id-ID") : "—"}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
