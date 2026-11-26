import { Icon } from "@iconify/react";
import type { MmdBrokerDrilldown } from "./mmdTypes";

interface MmdBrokerModalProps {
  open: boolean;
  loading: boolean;
  data: MmdBrokerDrilldown | null;
  onClose: () => void;
}

export default function MmdBrokerModal({ open, loading, data, onClose }: MmdBrokerModalProps) {
  if (!open) return null;

  const dailyRows = data?.daily_breakdown ?? data?.daily ?? data?.transactions ?? [];
  const summary = data?.summary;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-4 bg-black/60 backdrop-blur-xs">
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl w-full max-w-2xl overflow-hidden shadow-2xl flex flex-col max-h-[90vh]">
        <div className="bg-[var(--md-sys-color-surface-container-high)] px-4 py-3 border-b border-[var(--md-sys-color-outline-variant)] flex items-center justify-between">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="font-extrabold text-sm text-[var(--md-sys-color-on-surface)]">
              DETAIL TRANSAKSI BROKER
            </span>
            <span className="px-2 py-0.5 rounded-lg bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] font-mono font-extrabold text-xs">
              {data?.broker_code || "—"}
            </span>
            {data?.broker_name && (
              <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] truncate max-w-[220px]">
                {data.broker_name}
              </span>
            )}
          </div>
          <button
            type="button"
            onClick={onClose}
            className="p-1 rounded-lg hover:bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] cursor-pointer"
          >
            <Icon icon="ph:x-bold" width={18} />
          </button>
        </div>

        <div className="p-4 overflow-y-auto space-y-4 font-mono-nums text-xs">
          {loading ? (
            <div className="py-12 text-center text-[var(--md-sys-color-on-surface-variant)]">
              Memuat detail transaksi broker...
            </div>
          ) : summary ? (
            <div className="space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
                <div className="bg-[var(--md-sys-color-surface-container-low)] p-3 rounded-xl border border-[var(--color-positive)]/30 space-y-1">
                  <div className="text-[10px] text-[var(--color-positive)] font-bold uppercase">TOTAL BUY</div>
                  <div className="text-sm font-extrabold text-[var(--md-sys-color-on-surface)]">
                    {Math.round(summary.total_buy_lot).toLocaleString("id-ID")} lot
                  </div>
                  <div className="text-[11px] text-[var(--color-positive)]">
                    Rp {summary.total_buy_value_b.toLocaleString("id-ID")} M
                  </div>
                </div>

                <div className="bg-[var(--md-sys-color-surface-container-low)] p-3 rounded-xl border border-[var(--color-negative)]/30 space-y-1">
                  <div className="text-[10px] text-[var(--color-negative)] font-bold uppercase">TOTAL SELL</div>
                  <div className="text-sm font-extrabold text-[var(--md-sys-color-on-surface)]">
                    {Math.round(summary.total_sell_lot).toLocaleString("id-ID")} lot
                  </div>
                  <div className="text-[11px] text-[var(--color-negative)]">
                    Rp {summary.total_sell_value_b.toLocaleString("id-ID")} M
                  </div>
                </div>

                <div className="bg-[var(--md-sys-color-surface-container-low)] p-3 rounded-xl border border-[var(--md-sys-color-primary)]/30 space-y-1">
                  <div className="text-[10px] text-[var(--md-sys-color-primary)] font-bold uppercase">
                    NET POSITION
                  </div>
                  <div
                    className={`text-sm font-extrabold ${
                      summary.net_lot >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                    }`}
                  >
                    {summary.net_lot > 0 ? "+" : ""}
                    {Math.round(summary.net_lot).toLocaleString("id-ID")} lot
                  </div>
                  <div className="text-[11px] text-[var(--md-sys-color-on-surface)]">
                    Rp {summary.net_value_b.toLocaleString("id-ID")} M
                  </div>
                </div>
              </div>

              {dailyRows.length > 0 && (
                <div className="bg-[var(--md-sys-color-surface-container-low)] rounded-xl border border-[var(--md-sys-color-outline-variant)]/40 overflow-x-auto">
                  <table className="w-full text-left text-[11px] font-mono-nums">
                    <thead>
                      <tr className="bg-[var(--md-sys-color-surface-container-high)] text-[10px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]">
                        <th className="py-2 px-2.5">Tanggal</th>
                        <th className="py-2 px-2 text-right">Buy (Lot)</th>
                        <th className="py-2 px-2 text-right">Avg Buy</th>
                        <th className="py-2 px-2 text-right">Sell (Lot)</th>
                        <th className="py-2 px-2 text-right">Avg Sell</th>
                        <th className="py-2 px-2.5 text-right">Net (Lot)</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/30">
                      {dailyRows.map((d, idx) => (
                        <tr key={idx} className="hover:bg-[var(--md-sys-color-surface-container-high)]/60">
                          <td className="py-1.5 px-2.5 font-semibold text-[var(--md-sys-color-on-surface)]">
                            {d.date}
                          </td>
                          <td className="py-1.5 px-2 text-right text-[var(--color-positive)]">
                            {Math.round(d.buy_lot).toLocaleString("id-ID")}
                          </td>
                          <td className="py-1.5 px-2 text-right text-[var(--md-sys-color-on-surface-variant)]">
                            {d.buy_avg > 0 ? d.buy_avg.toLocaleString("id-ID") : "—"}
                          </td>
                          <td className="py-1.5 px-2 text-right text-[var(--color-negative)]">
                            {Math.round(d.sell_lot).toLocaleString("id-ID")}
                          </td>
                          <td className="py-1.5 px-2 text-right text-[var(--md-sys-color-on-surface-variant)]">
                            {d.sell_avg > 0 ? d.sell_avg.toLocaleString("id-ID") : "—"}
                          </td>
                          <td
                            className={`py-1.5 px-2.5 text-right font-bold ${
                              d.net_lot >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                            }`}
                          >
                            {d.net_lot > 0 ? "+" : ""}
                            {Math.round(d.net_lot).toLocaleString("id-ID")}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          ) : (
            <div className="py-8 text-center text-[var(--md-sys-color-on-surface-variant)]">
              Tidak ada data transaksi.
            </div>
          )}
        </div>

        <div className="bg-[var(--md-sys-color-surface-container-high)] px-4 py-2.5 border-t border-[var(--md-sys-color-outline-variant)] flex justify-end">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-1.5 rounded-xl bg-[var(--md-sys-color-surface-container-highest)] hover:bg-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface)] text-xs font-bold cursor-pointer"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  );
}
