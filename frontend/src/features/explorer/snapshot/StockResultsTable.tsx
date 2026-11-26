import { Link } from "react-router-dom";

import { fmtNumber, fmtPct, fmtSignedPct } from "@/components/research/format";
import { StockBadges } from "@/components/research/StockBadges";
import { pctClass } from "./snapshotHelpers";
import type { SortOrder, StockRow } from "./snapshotTypes";
import { SortHeader } from "./SortHeader";

export function StockResultsTable({
  rows,
  sort,
  order,
  onSort,
}: {
  rows: StockRow[];
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
              field="ticker"
              label="Emiten"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <th className="px-4 py-2.5 font-bold">Sektor</th>
            <SortHeader
              field="mtd"
              label="Harga / MTD"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <SortHeader
              field="holders"
              label="Holders"
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
              field="public"
              label="Publik"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <SortHeader
              field="top_holder"
              label="Top holder"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <SortHeader
              field="market_cap"
              label="Mkt cap (M)"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
            <SortHeader
              field="foreign_5d"
              label="Foreign 5D"
              align="right"
              current={sort}
              order={order}
              onSort={onSort}
            />
          </tr>
        </thead>
        <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
          {rows.map((row) => (
            <tr
              key={row.ticker}
              className="hover:bg-[var(--md-sys-color-surface-container-high)]/40"
            >
              <td className="px-4 py-2.5">
                <div className="flex items-center gap-1.5 flex-wrap">
                  <Link
                    to={`/ownership?tab=saham&ticker=${encodeURIComponent(row.ticker)}`}
                    className="font-mono font-bold text-[var(--md-sys-color-primary)] hover:underline"
                  >
                    {row.ticker}
                  </Link>
                  <StockBadges
                    badges={row.badges}
                    badge={row.badge}
                    conglomerate={row.conglomerate || row.konglo}
                    size="sm"
                  />
                </div>
                <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] max-w-[230px] line-clamp-1">
                  {row.name}
                </div>
              </td>
              <td className="px-4 py-2.5 text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
                {row.sector}
              </td>
              <td className="px-4 py-2.5 text-right font-mono">
                <div>{row.price === null ? "—" : row.price.toLocaleString("id-ID")}</div>
                <div className={`text-[10px] font-bold ${pctClass(row.mtd_pct)}`}>
                  {fmtSignedPct(row.mtd_pct)}
                </div>
              </td>
              <td className="px-4 py-2.5 text-right font-mono">
                {fmtNumber(row.holder_count)}
              </td>
              <td className="px-4 py-2.5 text-right font-mono">
                {fmtPct(row.controlling_pct)}
              </td>
              <td className="px-4 py-2.5 text-right font-mono">
                {fmtPct(row.public_pct)}
              </td>
              <td className="px-4 py-2.5 text-right font-mono font-bold">
                {fmtPct(row.max_holder_pct)}
              </td>
              <td className="px-4 py-2.5 text-right font-mono">
                {row.market_cap_b === null
                  ? "—"
                  : `${fmtNumber(row.market_cap_b, 0)} M`}
              </td>
              <td
                className={`px-4 py-2.5 text-right font-mono ${pctClass(row.foreign_5d)}`}
              >
                {row.foreign_5d === null
                  ? "—"
                  : row.foreign_5d.toLocaleString("id-ID", {
                      maximumFractionDigits: 0,
                    })}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
