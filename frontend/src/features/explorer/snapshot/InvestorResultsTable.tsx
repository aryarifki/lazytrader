import { Link } from "react-router-dom";

import { fmtNumber, fmtPct } from "@/components/research/format";
import type { InvestorRow, SortOrder } from "./snapshotTypes";
import { SortHeader } from "./SortHeader";

export function InvestorResultsTable({
  rows,
  sort,
  order,
  onSort,
}: {
  rows: InvestorRow[];
  sort: string;
  order: SortOrder;
  onSort: (field: string) => void;
}) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left text-[11px]">
        <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
          <tr>
            <SortHeader
              field="name"
              label="Investor"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <SortHeader
              field="tickers"
              label="Emiten"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <SortHeader
              field="total_pct"
              label="Total porsi"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <SortHeader
              field="max_pct"
              label="Porsi max"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <SortHeader
              field="controlling"
              label="Pengendali"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <SortHeader
              field="substantial"
              label="Substantial"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <SortHeader
              field="sectors"
              label="Sektor"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <th className="px-4 py-2.5 font-bold">Contoh emiten</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
          {rows.map((row) => (
            <tr
              key={row.name}
              className="hover:bg-[var(--md-sys-color-surface-container-high)]/40"
            >
              <td className="px-4 py-2.5 font-medium max-w-[360px]">
                <Link
                  to={`/ownership?tab=investor&name=${encodeURIComponent(row.name)}`}
                  className="line-clamp-2 hover:text-[var(--md-sys-color-primary)] hover:underline"
                >
                  {row.name}
                </Link>
              </td>
              <td className="px-4 py-2.5 text-right font-mono font-bold">
                {fmtNumber(row.ticker_count)}
              </td>
              <td className="px-4 py-2.5 text-right font-mono">
                {fmtPct(row.total_pct)}
              </td>
              <td className="px-4 py-2.5 text-right font-mono font-bold">
                {fmtPct(row.max_pct)}
              </td>
              <td className="px-4 py-2.5 text-right font-mono">
                {fmtNumber(row.controlling_count)}
              </td>
              <td className="px-4 py-2.5 text-right font-mono">
                {fmtNumber(row.substantial_count)}
              </td>
              <td className="px-4 py-2.5 text-right font-mono">
                {fmtNumber(row.sector_count)}
              </td>
              <td className="px-4 py-2.5">
                <div className="flex flex-wrap gap-1">
                  {row.sample_tickers.map((ticker) => (
                    <Link
                      key={ticker}
                      to={`/ownership?tab=saham&ticker=${encodeURIComponent(ticker)}`}
                      className="rounded border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] px-1.5 py-0.5 font-mono text-[9px] text-[var(--md-sys-color-primary)] hover:border-[var(--md-sys-color-primary)]/60"
                    >
                      {ticker}
                    </Link>
                  ))}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
