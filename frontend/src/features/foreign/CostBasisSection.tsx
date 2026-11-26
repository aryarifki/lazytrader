import React, { Suspense, useState } from "react";
import { Icon } from "@iconify/react";
import { formatIdrCompact } from "./deepdiveHelpers";
import type { BrokerEntry, CostBasisInfo } from "./deepdiveTypes";
import type { CostBasisTimeSeries } from "./charts/ForeignCostBasisChart";

const ForeignCostBasisChart = React.lazy(() => import("./charts/ForeignCostBasisChart"));

interface CostBasisSectionProps {
  windowDays: number;
  cb?: CostBasisInfo;
  pnlPct: number;
  timeSeries?: CostBasisTimeSeries | null;
}

export default function CostBasisSection({
  windowDays,
  cb,
  pnlPct,
  timeSeries,
}: CostBasisSectionProps) {
  const [brokerTab, setBrokerTab] = useState<"holders" | "sellers">("holders");

  const brokersList: BrokerEntry[] =
    brokerTab === "holders" ? cb?.top_holders || [] : cb?.top_sellers || [];

  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 sm:gap-5">
      <div className="lg:col-span-7 bg-[var(--md-sys-color-surface-container)] p-4 sm:p-5 rounded-2xl border border-[var(--md-sys-color-outline-variant)] shadow-xs flex flex-col justify-between">
        <div>
          <div className="flex flex-wrap items-center justify-between gap-2 pb-3 mb-3 border-b border-[var(--md-sys-color-outline-variant)]/60">
            <div className="flex items-center gap-2">
              <Icon icon="ph:coins-duotone" className="text-[var(--md-sys-color-primary)]" width="20" />
              <h3 className="text-xs sm:text-sm font-extrabold text-[var(--md-sys-color-on-surface)] uppercase tracking-tight">
                Harga vs Foreign Avg Cost (VWAP) &amp; Akumulasi
              </h3>
            </div>
            <div className="text-[11px] font-mono font-bold text-[var(--md-sys-color-on-surface)]">
              Net {windowDays}D:{" "}
              <span className={(cb?.total_net_value || 0) >= 0 ? "text-emerald-500" : "text-rose-500"}>
                {formatIdrCompact(cb?.total_net_value || 0)}
              </span>
            </div>
          </div>
          <div className="w-full h-[300px] relative z-0">
            <Suspense fallback={<div className="h-full w-full" />}>
              <ForeignCostBasisChart timeSeries={timeSeries} />
            </Suspense>
          </div>
        </div>
        <p className="mt-3 pt-2.5 border-t border-[var(--md-sys-color-outline-variant)]/50 text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
          Garis biru putus-putus adalah rata-rata harga beli asing (VWAP) selama {windowDays} hari. Harga saat ini (Rp{" "}
          {Number(cb?.latest_close || 0).toLocaleString("id-ID")}) berada{" "}
          <strong className="text-[var(--md-sys-color-on-surface)]">
            {pnlPct >= 0 ? `${pnlPct.toFixed(2)}% di atas` : `${Math.abs(pnlPct).toFixed(2)}% di bawah`}
          </strong>{" "}
          modal rata-rata asing.
        </p>
      </div>

      <div className="lg:col-span-5 bg-[var(--md-sys-color-surface-container)] p-4 sm:p-5 rounded-2xl border border-[var(--md-sys-color-outline-variant)] shadow-xs flex flex-col">
        <div className="flex items-center justify-between gap-2 pb-3 mb-3 border-b border-[var(--md-sys-color-outline-variant)]/60">
          <div className="flex items-center gap-2">
            <Icon icon="ph:buildings-duotone" className="text-[var(--md-sys-color-primary)]" width="20" />
            <h3 className="text-xs sm:text-sm font-extrabold text-[var(--md-sys-color-on-surface)] uppercase tracking-tight">
              Top Broker Asing
            </h3>
          </div>
          <div className="flex items-center p-0.5 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-[10px] font-bold">
            <button
              onClick={() => setBrokerTab("holders")}
              className={`px-2.5 py-1 rounded-md transition-all ${
                brokerTab === "holders"
                  ? "bg-emerald-500 text-white"
                  : "text-[var(--md-sys-color-on-surface-variant)]"
              }`}
            >
              Akumulasi
            </button>
            <button
              onClick={() => setBrokerTab("sellers")}
              className={`px-2.5 py-1 rounded-md transition-all ${
                brokerTab === "sellers"
                  ? "bg-rose-500 text-white"
                  : "text-[var(--md-sys-color-on-surface-variant)]"
              }`}
            >
              Distribusi
            </button>
          </div>
        </div>

        <div className="overflow-x-auto flex-1">
          {brokersList.length === 0 ? (
            <div className="flex h-48 items-center justify-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
              Tidak ada broker asing pada kategori ini.
            </div>
          ) : (
            <table className="w-full text-left text-xs">
              <thead>
                <tr className="text-[10px] uppercase text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/60">
                  <th className="py-2 pr-2">Broker</th>
                  <th className="py-2 px-2 text-right">Net Val</th>
                  <th className="py-2 px-2 text-right">Avg {brokerTab === "holders" ? "Beli" : "Jual"}</th>
                  <th className="py-2 pl-2 text-right">PnL / Hari</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
                {brokersList.map((b) => (
                  <tr key={b.broker_code} className="hover:bg-[var(--md-sys-color-surface-container-low)]">
                    <td className="py-2.5 pr-2">
                      <div className="font-mono font-extrabold text-[var(--md-sys-color-on-surface)]">
                        {b.broker_code}
                      </div>
                      <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] truncate max-w-[130px]">
                        {b.broker_name}
                      </div>
                    </td>
                    <td className="py-2.5 px-2 text-right font-mono">
                      <div className={`font-bold ${b.net_value >= 0 ? "text-emerald-500" : "text-rose-500"}`}>
                        {formatIdrCompact(b.net_value)}
                      </div>
                      <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
                        {Math.abs(b.net_lot).toLocaleString("id-ID")} lot
                      </div>
                    </td>
                    <td className="py-2.5 px-2 text-right font-mono font-semibold text-[var(--md-sys-color-on-surface)]">
                      {Number(brokerTab === "holders" ? b.avg_buy_price : b.avg_sell_price).toLocaleString("id-ID", {
                        maximumFractionDigits: 0,
                      })}
                    </td>
                    <td className="py-2.5 pl-2 text-right font-mono">
                      <div className={`font-bold ${b.floating_pnl_pct >= 0 ? "text-emerald-500" : "text-rose-500"}`}>
                        {brokerTab === "holders"
                          ? `${b.floating_pnl_pct >= 0 ? "+" : ""}${b.floating_pnl_pct.toFixed(1)}%`
                          : "-"}
                      </div>
                      <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
                        {b.active_days}h aktif
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}
