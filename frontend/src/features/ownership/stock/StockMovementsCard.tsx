import { Badge, Card, EmptyState, SectionTitle } from "@/components/research/ui";
import { fmtNumber, fmtPct, fmtRupiah } from "@/components/research/format";
import { actionTone, fmtLot } from "./stockHelpers";
import type { StockOwnership, T5Transaction } from "./stockTypes";

interface StockMovementsCardProps {
  data: StockOwnership;
  transactions: T5Transaction[];
  activeMovementTab: "staging" | "ct";
  onSelectMovementTab: (tab: "staging" | "ct") => void;
  hasStagingMovements: boolean;
  hasCtTransactions: boolean;
}

export function StockMovementsCard({
  data,
  transactions,
  activeMovementTab,
  onSelectMovementTab,
  hasStagingMovements,
  hasCtTransactions,
}: StockMovementsCardProps) {
  return (
    <Card className="p-4 sm:p-5">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4">
        <SectionTitle
          icon="ph:arrows-left-right-duotone"
          title={activeMovementTab === "staging" ? "Mutasi Pemegang Saham Terbuka (KSEI)" : "Transaksi T5 (CT Baseline)"}
          subtitle={
            activeMovementTab === "staging" && data.staging_movements_summary && data.staging_movements_summary.count > 0
              ? `${data.staging_movements_summary.count} mutasi · ${data.staging_movements_summary.buy_count} beli / ${data.staging_movements_summary.sell_count} jual`
              : data.t5_summary && data.t5_summary.count > 0
                ? `${data.t5_summary.count} transaksi · ${data.t5_summary.buy_count} beli / ${data.t5_summary.sell_count} jual`
                : "Transaksi dan mutasi pemegang saham signifikan"
          }
        />
        <div className="flex items-center gap-1.5 self-start sm:self-auto flex-wrap">
          {hasStagingMovements && hasCtTransactions && (
            <div className="inline-flex rounded-lg border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] p-0.5 text-[10px] font-medium">
              <button
                type="button"
                onClick={() => onSelectMovementTab("staging")}
                className={`rounded-md px-2.5 py-1 transition-colors ${
                  activeMovementTab === "staging"
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-sm"
                    : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                }`}
              >
                KSEI Staging ({data.staging_movements?.length ?? 0})
              </button>
              <button
                type="button"
                onClick={() => onSelectMovementTab("ct")}
                className={`rounded-md px-2.5 py-1 transition-colors ${
                  activeMovementTab === "ct"
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-sm"
                    : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                }`}
              >
                Baseline CT ({data.transactions?.length ?? 0})
              </button>
            </div>
          )}
          {activeMovementTab === "staging" && data.staging_movements_summary && data.staging_movements_summary.count > 0 ? (
            <Badge tone={data.staging_movements_summary.net_value_m >= 0 ? "up" : "down"}>
              Net {data.staging_movements_summary.net_value_m >= 0 ? "+" : ""}
              {fmtNumber(data.staging_movements_summary.net_value_m, 2)} M
            </Badge>
          ) : activeMovementTab === "ct" && data.t5_summary && data.t5_summary.count > 0 ? (
            <Badge tone={data.t5_summary.net_value_m >= 0 ? "up" : "down"}>
              Net {data.t5_summary.net_value_m >= 0 ? "+" : ""}
              {fmtNumber(data.t5_summary.net_value_m, 2)} M
            </Badge>
          ) : null}
        </div>
      </div>

      {activeMovementTab === "staging" ? (
        data.staging_movements && data.staging_movements.length > 0 ? (
          <div className="max-h-[460px] overflow-auto rounded-xl border border-[var(--md-sys-color-outline-variant)]/50">
            <table className="w-full min-w-[900px] text-left text-[11px]">
              <thead className="sticky top-0 z-10 bg-[var(--md-sys-color-surface-container-high)] text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
                <tr>
                  <th className="px-3 py-2.5 font-bold">Tanggal</th>
                  <th className="px-3 py-2.5 font-bold">Pemegang Saham</th>
                  <th className="px-3 py-2.5 font-bold">Broker</th>
                  <th className="px-3 py-2.5 font-bold">Aksi</th>
                  <th className="px-3 py-2.5 font-bold text-right">Perubahan Saham</th>
                  <th className="px-3 py-2.5 font-bold text-right">Harga</th>
                  <th className="px-3 py-2.5 font-bold text-right">Nilai</th>
                  <th className="px-3 py-2.5 font-bold text-right">Porsi Akhir</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
                {data.staging_movements.map((mov, index) => {
                  const action = mov.action.toUpperCase();
                  const chgShares = mov.change_shares ?? 0;
                  return (
                    <tr
                      key={`${mov.date}-${mov.name}-${index}`}
                      className="hover:bg-[var(--md-sys-color-surface-container-high)]/40"
                    >
                      <td className="px-3 py-2.5 font-mono whitespace-nowrap">{mov.date ?? "—"}</td>
                      <td className="px-3 py-2.5 max-w-[260px]">
                        <span className="line-clamp-2 font-medium text-[var(--md-sys-color-on-surface)]">
                          {mov.name}
                        </span>
                      </td>
                      <td className="px-3 py-2.5 font-mono">
                        {mov.broker ? <Badge tone="neutral">{mov.broker}</Badge> : "—"}
                      </td>
                      <td className="px-3 py-2.5">
                        <Badge tone={actionTone(action)}>
                          {action}
                        </Badge>
                      </td>
                      <td
                        className={`px-3 py-2.5 text-right font-mono font-bold whitespace-nowrap ${
                          chgShares > 0 ? "text-emerald-500" : chgShares < 0 ? "text-rose-500" : ""
                        }`}
                      >
                        {chgShares > 0 ? "+" : ""}{fmtLot(Math.round(chgShares / 100))} lot
                        <div className="text-[9px] font-normal text-[var(--md-sys-color-on-surface-variant)]">
                          {chgShares > 0 ? "+" : ""}{fmtNumber(chgShares)} lbr
                        </div>
                      </td>
                      <td className="px-3 py-2.5 text-right font-mono whitespace-nowrap">
                        {mov.price != null ? `Rp ${fmtNumber(mov.price)}` : "—"}
                      </td>
                      <td className="px-3 py-2.5 text-right font-mono font-bold whitespace-nowrap">
                        {mov.value_m != null ? `${fmtNumber(mov.value_m, 2)} M` : mov.value != null ? fmtRupiah(mov.value) : "—"}
                      </td>
                      <td className="px-3 py-2.5 text-right font-mono whitespace-nowrap">
                        {mov.current_pct != null ? fmtPct(mov.current_pct) : "—"}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        ) : (
          <EmptyState
            icon="ph:receipt-duotone"
            title="Tidak ada mutasi KSEI untuk emiten ini"
            description="Belum ada transaksi keterbukaan informasi pemegang saham signifikan yang dilaporkan ke KSEI/OJK untuk emiten ini."
          />
        )
      ) : transactions.length > 0 ? (
        <div className="max-h-[460px] overflow-auto rounded-xl border border-[var(--md-sys-color-outline-variant)]/50">
          <table className="w-full min-w-[900px] text-left text-[11px]">
            <thead className="sticky top-0 z-10 bg-[var(--md-sys-color-surface-container-high)] text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
              <tr>
                <th className="px-3 py-2.5 font-bold">Tanggal</th>
                <th className="px-3 py-2.5 font-bold">Pemegang Saham</th>
                <th className="px-3 py-2.5 font-bold">Broker</th>
                <th className="px-3 py-2.5 font-bold">Aksi</th>
                <th className="px-3 py-2.5 font-bold text-right">Lot</th>
                <th className="px-3 py-2.5 font-bold text-right">Value (M)</th>
                <th className="px-3 py-2.5 font-bold text-right">Lot Setelah</th>
                <th className="px-3 py-2.5 font-bold text-right">% Kepemilikan</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
              {transactions.map((tx, index) => {
                const action = (tx.action ?? "").toUpperCase();
                const lot = tx.change_lot ?? 0;
                return (
                  <tr
                    key={`${tx.date}-${tx.broker}-${tx.name}-${index}`}
                    className="hover:bg-[var(--md-sys-color-surface-container-high)]/40"
                  >
                    <td className="px-3 py-2.5 font-mono whitespace-nowrap">{tx.date ?? "—"}</td>
                    <td className="px-3 py-2.5 max-w-[260px]">
                      <span className="line-clamp-2">{tx.name ?? "—"}</span>
                    </td>
                    <td className="px-3 py-2.5 font-mono">{tx.broker ?? "—"}</td>
                    <td className="px-3 py-2.5">
                      <Badge tone={action === "BUY" ? "up" : action === "SELL" ? "down" : "neutral"}>
                        {action || "—"}
                      </Badge>
                    </td>
                    <td
                      className={`px-3 py-2.5 text-right font-mono font-bold whitespace-nowrap ${
                        lot > 0 ? "text-emerald-500" : lot < 0 ? "text-rose-500" : ""
                      }`}
                    >
                      {lot > 0 ? "+" : lot < 0 ? "-" : ""}
                      {fmtLot(Math.abs(lot))}
                    </td>
                    <td className="px-3 py-2.5 text-right font-mono whitespace-nowrap">
                      {tx.value_m != null ? `${fmtNumber(tx.value_m, 2)} M` : "—"}
                    </td>
                    <td className="px-3 py-2.5 text-right font-mono whitespace-nowrap">
                      {fmtLot(tx.lot_after)}
                    </td>
                    <td className="px-3 py-2.5 text-right font-mono whitespace-nowrap">
                      {fmtPct(tx.pct_juli)}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      ) : (
        <EmptyState
          icon="ph:receipt-duotone"
          title="Tidak ada transaksi T5 untuk emiten ini"
          description="Tidak terdeteksi transaksi pemegang saham signifikan (≥5%) dalam 30 hari terakhir untuk emiten ini. Data disajikan secara transparan tanpa angka sintetis."
        />
      )}
    </Card>
  );
}
