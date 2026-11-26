import { Badge, Card, EmptyState, SectionTitle } from "@/components/research/ui";
import { fmtNumber, fmtPct } from "@/components/research/format";
import { fmtLot } from "./groupHelpers";
import type { GroupDetail } from "./groupTypes";

interface GroupT5TableProps {
  detail: GroupDetail | undefined;
}

export function GroupT5Table({ detail }: GroupT5TableProps) {
  return (
    <Card className="p-4 sm:p-5">
      <SectionTitle
        icon="ph:receipt-duotone"
        title="Transaksi T5 grup"
        subtitle={
          detail?.t5_summary
            ? `${detail.t5_summary.count} transaksi · ${detail.t5_summary.buy_count} beli / ${detail.t5_summary.sell_count} jual`
            : "Transaksi besar emiten dalam grup"
        }
        action={
          detail?.t5_summary && detail.t5_summary.count > 0 ? (
            <Badge tone={detail.t5_summary.net_value_m >= 0 ? "up" : "down"}>
              Net {detail.t5_summary.net_value_m >= 0 ? "+" : ""}
              {fmtNumber(detail.t5_summary.net_value_m, 2)} M
            </Badge>
          ) : null
        }
      />
      {detail && detail.transactions.length > 0 ? (
        <div className="mt-3 max-h-[520px] overflow-auto rounded-xl border border-[var(--md-sys-color-outline-variant)]/50">
          <table className="w-full min-w-[1050px] text-left text-[11px]">
            <thead className="sticky top-0 z-10 bg-[var(--md-sys-color-surface-container-high)] text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
              <tr>
                <th className="px-3 py-2.5 font-bold">Tanggal</th>
                <th className="px-3 py-2.5 font-bold">Emiten</th>
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
              {detail.transactions.map((tx, index) => {
                const action = (tx.action ?? "").toUpperCase();
                const lot = tx.change_lot ?? 0;
                return (
                  <tr
                    key={`${tx.ticker}-${tx.date}-${tx.broker}-${tx.name}-${index}`}
                    className="hover:bg-[var(--md-sys-color-surface-container-high)]/40"
                  >
                    <td className="px-3 py-2.5 font-mono whitespace-nowrap">{tx.date ?? "—"}</td>
                    <td className="px-3 py-2.5 font-mono font-bold text-[var(--md-sys-color-primary)]">
                      {tx.ticker}
                    </td>
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
        <div className="mt-3">
          <EmptyState
            icon="ph:receipt-duotone"
            title="Belum ada transaksi T5 untuk grup ini"
            description="Baseline CT tidak memuat transaksi T5 untuk emiten di grup ini. Tidak ada data sintetis yang ditampilkan."
          />
        </div>
      )}
    </Card>
  );
}
