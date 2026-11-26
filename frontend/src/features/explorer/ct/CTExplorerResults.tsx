import { Icon } from "@iconify/react";

import {
  Badge,
  Card,
  EmptyState,
  ErrorBlock,
  LoadingBlock,
  SectionTitle,
} from "@/components/research/ui";
import { fmtNumber } from "@/components/research/format";
import type { Category, Column, ExplorerResponse, ExplorerRow, Mode, SortOrder } from "./ctTypes";

interface CTExplorerResultsProps {
  data: ExplorerResponse | undefined;
  error: Error | undefined;
  isLoading: boolean;
  rows: ExplorerRow[];
  columns: Column[];
  mode: Mode;
  category: Category;
  shown: number;
  total: number;
  sort: string;
  order: SortOrder;
  offset: number;
  limit: number;
  onToggleSort: (column: Column) => void;
  onPrevPage: () => void;
  onNextPage: () => void;
}

export function CTExplorerResults({
  data,
  error,
  isLoading,
  rows,
  columns,
  mode,
  category,
  shown,
  total,
  sort,
  order,
  offset,
  onToggleSort,
  onPrevPage,
  onNextPage,
}: CTExplorerResultsProps) {
  return (
    <Card className="overflow-hidden">
      <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <SectionTitle
          icon="ph:table-duotone"
          title={data?.condition_label ?? "Hasil explorer"}
          subtitle={`${fmtNumber(shown)} dari ${fmtNumber(total)} baris · ${mode === "stocks" ? "By Stocks" : "By Investors"}`}
          action={
            <div className="flex items-center gap-2">
              <Badge tone="primary">{category}</Badge>
              <Badge>{mode}</Badge>
            </div>
          }
        />
      </div>

      {isLoading ? <div className="p-4"><LoadingBlock label="Menghitung explorer..." /></div> : null}
      {error ? (
        <div className="p-4">
          <ErrorBlock message={`Gagal memuat explorer: ${error.message}`} />
        </div>
      ) : null}

      {!isLoading && !error && rows.length === 0 ? (
        <div className="p-4">
          <EmptyState
            icon="ph:funnel-duotone"
            title="Tidak ada hasil"
            description="Kombinasi default CT memang bisa menghasilkan 0 baris. Longgarkan rentang kepemilikan, ganti market cap, atau tambah tipe investor."
          />
        </div>
      ) : null}

      {!isLoading && !error && rows.length > 0 ? (
        <>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-[11px]">
              <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
                <tr>
                  {columns.map((column) => {
                    const active = column.sort && sort === column.sort;
                    return (
                      <th
                        key={column.key}
                        className={`px-3 py-2.5 font-bold ${column.align === "right" ? "text-right" : ""}`}
                      >
                        {column.sort ? (
                          <button
                            type="button"
                            onClick={() => onToggleSort(column)}
                            className={`inline-flex items-center gap-1 ${
                              active ? "text-[var(--md-sys-color-primary)]" : ""
                            }`}
                          >
                            {column.label}
                            <Icon
                              icon={
                                active
                                  ? order === "desc"
                                    ? "ph:sort-descending-duotone"
                                    : "ph:sort-ascending-duotone"
                                  : "ph:arrows-down-up-duotone"
                              }
                              width="12"
                              className={active ? "opacity-100" : "opacity-40"}
                            />
                          </button>
                        ) : (
                          column.label
                        )}
                      </th>
                    );
                  })}
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
                {rows.map((row, index) => (
                  <tr
                    key={`${row.ticker ?? row.name ?? "row"}-${index}`}
                    className="hover:bg-[var(--md-sys-color-surface-container-high)]/40"
                  >
                    {columns.map((column) => (
                      <td
                        key={column.key}
                        className={`px-3 py-2.5 ${column.align === "right" ? "text-right font-mono" : ""}`}
                      >
                        {column.render(row, index)}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className="p-4 border-t border-[var(--md-sys-color-outline-variant)]/60 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
            <span className="text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
              Offset {offset} · {shown} baris · total {total}
            </span>
            <div className="flex items-center gap-2">
              <button
                type="button"
                onClick={onPrevPage}
                disabled={offset <= 0}
                className="rounded-lg border border-[var(--md-sys-color-outline-variant)] px-3 py-1.5 text-[10px] font-bold disabled:opacity-40"
              >
                Sebelumnya
              </button>
              <button
                type="button"
                onClick={onNextPage}
                disabled={offset + shown >= total}
                className="rounded-lg border border-[var(--md-sys-color-outline-variant)] px-3 py-1.5 text-[10px] font-bold disabled:opacity-40"
              >
                Berikutnya
              </button>
            </div>
          </div>
        </>
      ) : null}
    </Card>
  );
}
