import { Badge, Card, EmptyState, LoadingBlock, SectionTitle } from "@/components/research/ui";
import { fmtPct } from "@/components/research/format";
import type { GroupDetail } from "./groupTypes";

interface GroupMatrixViewProps {
  detail: GroupDetail | undefined;
  detailLoading: boolean;
  matrixTickers: string[];
}

export function GroupMatrixView({
  detail,
  detailLoading,
  matrixTickers,
}: GroupMatrixViewProps) {
  return (
    <Card className="overflow-hidden">
      <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <SectionTitle
          icon="ph:table-duotone"
          title="Matriks kepemilikan silang"
          subtitle="Pemegang yang muncul di minimal 2 emiten dalam grup"
          action={<Badge>{detail?.top_holdings.length ?? 0} entitas</Badge>}
        />
      </div>
      {detailLoading ? (
        <div className="p-5">
          <LoadingBlock label="Memuat matriks..." />
        </div>
      ) : detail && detail.top_holdings.length > 0 ? (
        <div className="overflow-x-auto">
          <table className="w-full text-left text-[10px]">
            <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
              <tr>
                <th className="sticky left-0 z-10 bg-[var(--md-sys-color-surface-container-high)] px-3 py-2.5 font-bold min-w-[210px]">
                  Pemegang
                </th>
                <th className="px-2 py-2.5 font-bold text-right">Total</th>
                {matrixTickers.map((ticker) => (
                  <th key={ticker} className="px-2 py-2.5 font-bold text-center min-w-[62px]">
                    {ticker}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
              {detail.top_holdings.map((holder) => {
                const byTicker = new Map(holder.holdings.map((item) => [item.ticker, item]));
                return (
                  <tr
                    key={holder.name}
                    className="hover:bg-[var(--md-sys-color-surface-container-high)]/40"
                  >
                    <td className="sticky left-0 z-10 bg-[var(--md-sys-color-surface-container)] px-3 py-2 font-medium max-w-[260px]">
                      <span className="line-clamp-1">{holder.name}</span>
                      <span className="text-[8px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                        {holder.ticker_count} emiten
                      </span>
                    </td>
                    <td className="px-2 py-2 text-right font-mono font-bold text-[var(--md-sys-color-primary)]">
                      {fmtPct(holder.total_pct)}
                    </td>
                    {matrixTickers.map((ticker) => {
                      const item = byTicker.get(ticker);
                      if (!item) {
                        return (
                          <td key={ticker} className="px-2 py-2 text-center text-[var(--md-sys-color-on-surface-variant)]/40">
                            ·
                          </td>
                        );
                      }
                      return (
                        <td
                          key={ticker}
                          className={`px-2 py-2 text-center font-mono ${
                            item.is_controlling
                              ? "text-amber-700 dark:text-amber-400 font-bold"
                              : "text-[var(--md-sys-color-on-surface)]"
                          }`}
                          title={`${holder.name} di ${ticker}: ${fmtPct(item.pct, 3)}`}
                        >
                          {item.pct > 0 ? item.pct.toFixed(2) : "0"}
                        </td>
                      );
                    })}
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      ) : !detailLoading ? (
        <div className="p-5">
          <EmptyState
            icon="ph:table-duotone"
            title="Matriks belum tersedia"
            description="Tidak ada pemegang yang muncul di lebih dari satu emiten pada grup ini."
          />
        </div>
      ) : null}
      <div className="p-4 border-t border-[var(--md-sys-color-outline-variant)]/60">
        <div className="text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
          Nilai adalah persentase saham emiten. Sel kosong berarti pemegang tidak tercatat di emiten tersebut.
        </div>
      </div>
    </Card>
  );
}
