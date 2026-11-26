import { Link } from "react-router-dom";

import { Badge, Card, SectionTitle } from "@/components/research/ui";
import { fmtPct } from "@/components/research/format";
import { fmtLotT5, fmtValueM } from "./investorHelpers";
import type { InvestorTransaction } from "./investorTypes";

export function InvestorT5Table({ transactions }: { transactions: InvestorTransaction[] }) {
  return (
    <Card className="overflow-hidden">
      <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <SectionTitle
          icon="ph:swap-duotone"
          title="Transaksi T5"
          subtitle={`Semua saham · ${transactions.length} transaksi`}
          action={<Badge tone="primary">T5</Badge>}
        />
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left text-[11px]">
          <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
            <tr>
              <th className="px-3 py-2.5 font-bold">Tanggal</th>
              <th className="px-3 py-2.5 font-bold">Emiten</th>
              <th className="px-3 py-2.5 font-bold">Broker</th>
              <th className="px-3 py-2.5 font-bold">Action</th>
              <th className="px-3 py-2.5 font-bold text-right">Lot</th>
              <th className="px-3 py-2.5 font-bold text-right">Value (M)</th>
              <th className="px-3 py-2.5 font-bold text-right">(Lot) Setelah</th>
              <th className="px-3 py-2.5 font-bold text-right">%</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
            {transactions.map((tx, index) => {
              const action = tx.action?.toUpperCase();
              const actionTone =
                action === "BUY"
                  ? "text-emerald-600 dark:text-emerald-400"
                  : action === "SELL"
                    ? "text-rose-600 dark:text-rose-400"
                    : "text-[var(--md-sys-color-on-surface-variant)]";
              return (
                <tr
                  key={`${tx.ticker}-${tx.date}-${tx.broker}-${index}`}
                  className="hover:bg-[var(--md-sys-color-surface-container-high)]/40"
                >
                  <td className="px-3 py-2.5 font-mono text-[10px]">{tx.date}</td>
                  <td className="px-3 py-2.5">
                    <Link
                      to={`/ownership?tab=saham&ticker=${encodeURIComponent(tx.ticker)}`}
                      className="font-mono font-bold text-[var(--md-sys-color-primary)] hover:underline"
                    >
                      {tx.ticker}
                    </Link>
                  </td>
                  <td className="px-3 py-2.5 font-mono text-[10px]">{tx.broker}</td>
                  <td className={`px-3 py-2.5 font-mono font-bold ${actionTone}`}>
                    {action || "—"}
                  </td>
                  <td className="px-3 py-2.5 text-right font-mono">
                    {fmtLotT5(Math.abs(tx.change_lot ?? 0))}
                  </td>
                  <td className="px-3 py-2.5 text-right font-mono">
                    {fmtValueM(tx.value_m)}
                  </td>
                  <td className="px-3 py-2.5 text-right font-mono">
                    {fmtLotT5(tx.lot_after)}
                  </td>
                  <td className="px-3 py-2.5 text-right font-mono">
                    {fmtPct(tx.pct_juli)}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </Card>
  );
}
