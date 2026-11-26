import { Icon } from "@iconify/react";

import { Card, EmptyState, SectionTitle } from "@/components/research/ui";
import { fmtNumber } from "@/components/research/format";
import { StockMonthlyMomTable } from "./StockMonthlyMomTable";
import { TrendMini } from "./TrendMini";
import { useStockMonthlyRows } from "./useStockMonthlyRows";
import type { MonthlyPoint, StockOwnership } from "./stockTypes";

interface StockMonthlyTrendCardProps {
  data: StockOwnership;
  history: MonthlyPoint[];
  showMonthlyTable: boolean;
  onToggleMonthlyTable: () => void;
}

export function StockMonthlyTrendCard({
  data,
  history,
  showMonthlyTable,
  onToggleMonthlyTable,
}: StockMonthlyTrendCardProps) {
  const monthlyComparisonRows = useStockMonthlyRows(history);

  const latestMonth = history.length > 0 ? history[history.length - 1] : null;
  const prevMonth = history.length > 1 ? history[history.length - 2] : null;
  const comparisonLabel = prevMonth?.period_label ?? null;

  const latestForeignTotal = latestMonth
    ? (latestMonth.foreign_total_pct ?? Number((latestMonth.foreign_named_pct + latestMonth.foreign_hidden_pct).toFixed(4)))
    : 0;
  const prevForeignTotal = prevMonth
    ? (prevMonth.foreign_total_pct ?? Number((prevMonth.foreign_named_pct + prevMonth.foreign_hidden_pct).toFixed(4)))
    : null;
  const foreignMomDelta = prevForeignTotal != null ? Number((latestForeignTotal - prevForeignTotal).toFixed(4)) : null;

  const latestForeignHidden = latestMonth?.foreign_hidden_pct ?? 0;
  const prevForeignHidden = prevMonth?.foreign_hidden_pct ?? null;
  const foreignHiddenMomDelta = prevForeignHidden != null ? Number((latestForeignHidden - prevForeignHidden).toFixed(4)) : null;

  const latestRetail = latestMonth?.retail_pct ?? 0;
  const prevRetail = prevMonth?.retail_pct ?? null;
  const retailMomDelta = prevRetail != null ? Number((latestRetail - prevRetail).toFixed(4)) : null;

  const latestLocalNonRetail = latestMonth?.local_non_retail_pct ?? 0;
  const prevLocalNonRetail = prevMonth?.local_non_retail_pct ?? null;
  const localNonRetailMomDelta = prevLocalNonRetail != null ? Number((latestLocalNonRetail - prevLocalNonRetail).toFixed(4)) : null;

  const latestSid = latestMonth?.total_sid ?? 0;
  const prevSid = prevMonth?.total_sid ?? null;
  const sidMomDelta = latestMonth?.sid_change_mom ?? (prevSid != null && latestSid ? latestSid - prevSid : null);
  const sidMomDeltaPct = latestMonth?.sid_change_pct ?? (prevSid && prevSid > 0 && sidMomDelta != null ? Number(((sidMomDelta / prevSid) * 100).toFixed(2)) : null);

  return (
    <Card className="p-4 sm:p-5 flex flex-col justify-between">
      <div>
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
          <SectionTitle
            icon="ph:chart-line-duotone"
            title="Tren kepemilikan bulanan"
            subtitle={
              history.length > 0
                ? `${history.length} snapshot ${data.ct_history?.length ? "baseline CT" : "KSEI"} tersedia · Perbandingan MoM`
                : "Snapshot bulanan belum tersedia"
            }
          />
          {history.length > 1 ? (
            <button
              type="button"
              onClick={onToggleMonthlyTable}
              className="inline-flex items-center gap-1.5 px-2.5 py-1 text-[11px] font-medium rounded-lg border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] transition-colors self-start sm:self-auto cursor-pointer"
            >
              <Icon icon={showMonthlyTable ? "ph:table-duotone" : "ph:table-bold"} width="13" />
              <span>{showMonthlyTable ? "Tutup Tabel" : "Lihat Tabel MoM"}</span>
            </button>
          ) : null}
        </div>

        {history.length > 0 ? (
          <>
            <div className="mt-4 grid grid-cols-2 gap-2.5">
              <TrendMini
                label="Asing Total"
                value={latestForeignTotal}
                points={history.map((item) => item.foreign_total_pct ?? Number((item.foreign_named_pct + item.foreign_hidden_pct).toFixed(4)))}
                color="#3A7BD5"
                delta={foreignMomDelta}
                comparisonLabel={comparisonLabel}
              />
              <TrendMini
                label="Foreign Hidden Player"
                value={latestForeignHidden}
                points={history.map((item) => item.foreign_hidden_pct)}
                color="#2A5BA0"
                delta={foreignHiddenMomDelta}
                comparisonLabel={comparisonLabel}
              />
              <TrendMini
                label="Retail"
                value={latestRetail}
                points={history.map((item) => item.retail_pct)}
                color="#1DB87A"
                delta={retailMomDelta}
                reverseColor
                comparisonLabel={comparisonLabel}
              />
              <TrendMini
                label="Lokal Non-Ritel"
                value={latestLocalNonRetail}
                points={history.map((item) => item.local_non_retail_pct)}
                color="#8C6A00"
                delta={localNonRetailMomDelta}
                comparisonLabel={comparisonLabel}
              />
              {history.some((item) => item.total_sid != null) ? (
                <TrendMini
                  label="Investor (SID)"
                  value={latestSid}
                  points={history
                    .map((item) => item.total_sid)
                    .filter((value): value is number => value != null)}
                  color="#E07A5F"
                  format={(value) => fmtNumber(value)}
                  delta={sidMomDelta}
                  deltaPct={sidMomDeltaPct}
                  deltaUnit=" SID"
                  isSid
                  comparisonLabel={comparisonLabel}
                  className="col-span-2"
                />
              ) : null}
            </div>

            {showMonthlyTable && monthlyComparisonRows.length > 0 ? (
              <StockMonthlyMomTable rows={monthlyComparisonRows} />
            ) : null}
          </>
        ) : (
          <div className="mt-4">
            <EmptyState
              icon="ph:chart-line-duotone"
              title="Tren bulanan belum tersedia"
              description="Emiten ini belum memiliki snapshot KSEI bulanan di database. Komposisi di samping memakai snapshot pemegang saham saat ini."
            />
          </div>
        )}
      </div>
    </Card>
  );
}
