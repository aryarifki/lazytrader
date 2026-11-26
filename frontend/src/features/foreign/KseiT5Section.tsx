import React, { Suspense } from "react";
import { Icon } from "@iconify/react";
import type { KseiRealityInfo, T5Entry } from "./deepdiveTypes";

const KseiRealityChart = React.lazy(() => import("./charts/KseiRealityChart"));

interface KseiT5SectionProps {
  activeTicker: string;
  ksei?: KseiRealityInfo;
}

export default function KseiT5Section({ activeTicker, ksei }: KseiT5SectionProps) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 sm:gap-5">
      <div className="lg:col-span-7 bg-[var(--md-sys-color-surface-container)] p-4 sm:p-5 rounded-2xl border border-[var(--md-sys-color-outline-variant)] shadow-xs flex flex-col justify-between">
        <div>
          <div className="flex flex-wrap items-center justify-between gap-2 pb-3 mb-3 border-b border-[var(--md-sys-color-outline-variant)]/60">
            <div className="flex items-center gap-2">
              <Icon icon="ph:bank-duotone" className="text-[var(--md-sys-color-primary)]" width="20" />
              <h3 className="text-xs sm:text-sm font-extrabold text-[var(--md-sys-color-on-surface)] uppercase tracking-tight">
                KSEI Reality Check (6 Bulan Kepemilikan Resmi)
              </h3>
            </div>
            <span className="text-[10px] font-bold px-2.5 py-0.5 rounded-full bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)]">
              {ksei?.alignment_status}
            </span>
          </div>
          <div className="w-full h-[250px] relative z-0">
            <Suspense fallback={<div className="h-full w-full" />}>
              <KseiRealityChart history={ksei?.history} />
            </Suspense>
          </div>
        </div>
        <div className="mt-3 pt-2.5 border-t border-[var(--md-sys-color-outline-variant)]/50 flex flex-wrap items-center justify-between gap-2 text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
          <span>
            Periode KSEI terakhir: <strong className="text-[var(--md-sys-color-on-surface)]">{ksei?.latest_period || "-"}</strong>
          </span>
          <span className="font-mono">
            Asing &gt;1%: <strong>{ksei?.pct_asing_gt1 ?? 0}%</strong> · Asing &lt;1%: <strong>{ksei?.pct_asing_lt1 ?? 0}%</strong> · Ritel Murni: <strong>{ksei?.latest_pct_ritel_pure ?? 0}%</strong>
          </span>
        </div>
      </div>

      <div className="lg:col-span-5 bg-[var(--md-sys-color-surface-container)] p-4 sm:p-5 rounded-2xl border border-[var(--md-sys-color-outline-variant)] shadow-xs flex flex-col">
        <div className="flex items-center justify-between gap-2 pb-3 mb-3 border-b border-[var(--md-sys-color-outline-variant)]/60">
          <div className="flex items-center gap-2">
            <Icon icon="ph:list-magnifying-glass-duotone" className="text-[var(--md-sys-color-primary)]" width="20" />
            <h3 className="text-xs sm:text-sm font-extrabold text-[var(--md-sys-color-on-surface)] uppercase tracking-tight">
              Transaksi Pemegang Saham &gt;1% (T5)
            </h3>
          </div>
        </div>
        <div className="overflow-x-auto flex-1">
          {!ksei?.t5_recent || ksei.t5_recent.length === 0 ? (
            <div className="flex h-48 items-center justify-center text-xs text-[var(--md-sys-color-on-surface-variant)] text-center px-4">
              Belum ada catatan mutasi pemegang saham &gt;1% (T5) terbaru untuk {activeTicker}.
            </div>
          ) : (
            <table className="w-full text-left text-xs">
              <thead>
                <tr className="text-[10px] uppercase text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/60">
                  <th className="py-2 pr-2">Tanggal / Investor</th>
                  <th className="py-2 px-2">Aksi</th>
                  <th className="py-2 pl-2 text-right">Nilai (M)</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
                {(ksei.t5_recent as T5Entry[]).map((t, idx) => (
                  <tr key={`${t.date}-${t.investor_name}-${idx}`} className="hover:bg-[var(--md-sys-color-surface-container-low)]">
                    <td className="py-2 pr-2">
                      <div className="font-semibold text-[var(--md-sys-color-on-surface)] truncate max-w-[180px]">
                        {t.investor_name}
                      </div>
                      <div className="text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                        {t.date} · Broker {t.broker}
                      </div>
                    </td>
                    <td className="py-2 px-2">
                      <span
                        className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${
                          t.action.toUpperCase().includes("BUY")
                            ? "bg-emerald-500/15 text-emerald-500"
                            : "bg-rose-500/15 text-rose-500"
                        }`}
                      >
                        {t.action}
                      </span>
                    </td>
                    <td className="py-2 pl-2 text-right font-mono font-bold text-[var(--md-sys-color-on-surface)]">
                      {t.value_m > 0 ? `Rp ${t.value_m.toLocaleString("id-ID")} M` : `${t.change_lot.toLocaleString("id-ID")} lot`}
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
