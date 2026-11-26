import { useDeferredValue, useEffect, useMemo, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import useSWR from "swr";

import {
  Badge,
  Card,
  EmptyState,
  ErrorBlock,
  LoadingBlock,
  SectionTitle,
} from "@/components/research/ui";
import { fmtNumber } from "@/components/research/format";
import { InvestorResultsTable } from "./InvestorResultsTable";
import {
  EMPTY_FILTERS,
  INVESTOR_SORT_FALLBACK,
  STOCK_SORT_FALLBACK,
} from "./snapshotConstants";
import { exportSnapshotCsv } from "./snapshotExport";
import { SnapshotFilters } from "./SnapshotFilters";
import { SnapshotHeaderCard } from "./SnapshotHeaderCard";
import { fetcher, intParam, numberParam, pctParam } from "./snapshotHelpers";
import { SnapshotPagination } from "./SnapshotPagination";
import { SnapshotStatCards } from "./SnapshotStatCards";
import type {
  Condition,
  ExplorerResponse,
  Filters,
  InvestorRow,
  Mode,
  SortOrder,
  StockRow,
} from "./snapshotTypes";
import { SnapshotUnsupported } from "./SnapshotUnsupported";
import { StockResultsTable } from "./StockResultsTable";
import { useActiveFilters } from "./useActiveFilters";

export default function ExplorerDashboard({ embedded = false }: { embedded?: boolean } = {}) {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const modeParam = searchParams.get("mode");
  const initialMode: Mode = modeParam === "investors" ? "investors" : "stocks";
  const [mode, setMode] = useState<Mode>(initialMode);
  const [condition, setCondition] = useState("all");
  const [query, setQuery] = useState("");
  const [filters, setFilters] = useState<Filters>({ ...EMPTY_FILTERS });
  const [sort, setSort] = useState("default");
  const [order, setOrder] = useState<SortOrder>("desc");
  const [limit, setLimit] = useState(50);
  const [offset, setOffset] = useState(0);
  const deferredQuery = useDeferredValue(query.trim());

  useEffect(() => {
    if (modeParam === "investors" && mode !== "investors") {
      setMode("investors");
    } else if (modeParam === "stocks" && mode !== "stocks") {
      setMode("stocks");
    }
  }, [modeParam]);

  useEffect(() => {
    setOffset(0);
  }, [mode, condition, deferredQuery, filters, sort, order, limit]);

  const queryString = useMemo(() => {
    const params = new URLSearchParams();
    params.set("mode", mode);
    params.set("condition", condition);
    params.set("q", deferredQuery);
    params.set("sort", sort);
    params.set("order", order);
    params.set("limit", String(limit));
    params.set("offset", String(offset));

    if (mode === "stocks") {
      if (filters.sector.trim()) params.set("sector", filters.sector.trim());
      const intFields: Array<[string, string]> = [
        ["min_holders", filters.minHolders],
        ["max_holders", filters.maxHolders],
      ];
      for (const [name, value] of intFields) {
        const parsed = intParam(value);
        if (parsed !== null) params.set(name, parsed);
      }
      const pctFields: Array<[string, string]> = [
        ["min_controlling_pct", filters.minControlling],
        ["max_controlling_pct", filters.maxControlling],
        ["min_public_pct", filters.minPublic],
        ["max_public_pct", filters.maxPublic],
        ["min_top_holder_pct", filters.minTopHolder],
        ["max_top_holder_pct", filters.maxTopHolder],
      ];
      for (const [name, value] of pctFields) {
        const parsed = pctParam(value);
        if (parsed !== null) params.set(name, parsed);
      }
      const minCap = numberParam(filters.minMarketCap);
      if (minCap !== null) params.set("min_market_cap_b", minCap);
      const maxCap = numberParam(filters.maxMarketCap);
      if (maxCap !== null) params.set("max_market_cap_b", maxCap);
      if (filters.mtd !== "all") params.set("mtd", filters.mtd);
      if (filters.foreignFlow !== "all") {
        params.set("foreign_flow", filters.foreignFlow);
      }
      if (filters.hasTreasury) params.set("has_treasury", "true");
    } else {
      const intFields: Array<[string, string]> = [
        ["min_tickers", filters.minTickers],
        ["max_tickers", filters.maxTickers],
        ["min_sector_count", filters.minSectorCount],
      ];
      for (const [name, value] of intFields) {
        const parsed = intParam(value);
        if (parsed !== null) params.set(name, parsed);
      }
      const parsedMaxPct = pctParam(filters.minMaxPct);
      if (parsedMaxPct !== null) params.set("min_max_pct", parsedMaxPct);
    }

    return params.toString();
  }, [mode, condition, deferredQuery, filters, sort, order, limit, offset]);

  const { data, error, isLoading, isValidating } = useSWR<ExplorerResponse>(
    `/api/v1/ownership/explorer?${queryString}`,
    fetcher
  );

  const stockRows = useMemo(
    () => (data?.mode === "stocks" ? (data.results as StockRow[]) : []),
    [data]
  );
  const investorRows = useMemo(
    () => (data?.mode === "investors" ? (data.results as InvestorRow[]) : []),
    [data]
  );

  const sortOptions =
    data?.sorts ?? (mode === "stocks" ? STOCK_SORT_FALLBACK : INVESTOR_SORT_FALLBACK);
  const sectors = data?.sectors ?? [];
  const conditions: Condition[] =
    data?.conditions ?? [{ id: "all", label: "Semua", available: true }];
  const conditionGroups = data?.condition_groups ?? [];
  const total = data?.total ?? data?.count ?? 0;
  const shown = data?.count ?? 0;

  const activeFilters = useActiveFilters({
    condition,
    conditionLabel: data?.condition_label,
    deferredQuery,
    filters,
    mode,
    setCondition,
    setQuery,
    setFilters,
  });

  const changeMode = (nextMode: Mode) => {
    setMode(nextMode);
    setCondition("all");
    setFilters({ ...EMPTY_FILTERS });
    setSort("default");
    setOrder("desc");
    setOffset(0);
    const params = new URLSearchParams(window.location.search);
    params.set("mode", nextMode);
    navigate(`/explorer?${params.toString()}`, {
      replace: true,
      preventScrollReset: true,
    });
  };

  const resetAll = () => {
    setCondition("all");
    setQuery("");
    setFilters({ ...EMPTY_FILTERS });
    setSort("default");
    setOrder("desc");
    setOffset(0);
  };

  const toggleSort = (field: string) => {
    if (sort === field) {
      setOrder((prev) => (prev === "asc" ? "desc" : "asc"));
      return;
    }
    setSort(field);
    setOrder(field === "ticker" || field === "name" ? "asc" : "desc");
  };

  return (
    <div className={embedded ? "space-y-4" : "max-w-7xl mx-auto space-y-4"}>
      <SnapshotHeaderCard
        embedded={embedded}
        mode={mode}
        condition={condition}
        query={query}
        data={data}
        total={total}
        isLoading={isLoading}
        isValidating={isValidating}
        onChangeMode={changeMode}
        onChangeQuery={setQuery}
        onSelectPreset={setCondition}
        onResetAll={resetAll}
      />

      <SnapshotStatCards
        total={total}
        shown={shown}
        activeFilterCount={activeFilters.length}
        source={data?.data_quality.source}
      />

      <SnapshotFilters
        mode={mode}
        condition={condition}
        conditions={conditions}
        conditionGroups={conditionGroups}
        sectors={sectors}
        filters={filters}
        sort={sort}
        order={order}
        limit={limit}
        sortOptions={sortOptions}
        shown={shown}
        activeFilters={activeFilters}
        onChangeCondition={setCondition}
        setFilters={setFilters}
        onChangeSort={setSort}
        onToggleOrder={() => setOrder((prev) => (prev === "asc" ? "desc" : "asc"))}
        onChangeLimit={setLimit}
        onExportCsv={() => exportSnapshotCsv(mode, stockRows, investorRows)}
      />

      {isLoading ? <LoadingBlock label="Menjalankan screener..." /> : null}
      {error ? <ErrorBlock message={`Gagal menjalankan explorer: ${error.message}`} /> : null}

      {data && total === 0 && !isLoading ? (
        <EmptyState
          icon="ph:funnel-duotone"
          title="Tidak ada hasil"
          description="Coba ubah kondisi, longgarkan rentang, atau hapus beberapa filter."
        />
      ) : null}

      {data && total > 0 ? (
        <Card className="overflow-hidden">
          <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)]/60">
            <SectionTitle
              icon="ph:list-magnifying-glass-duotone"
              title={mode === "stocks" ? "Hasil screener saham" : "Hasil screener investor"}
              subtitle={data.condition_label}
              action={
                <Badge tone="primary">
                  {fmtNumber(shown)} / {fmtNumber(total)}
                </Badge>
              }
            />
          </div>

          {mode === "stocks" ? (
            <StockResultsTable
              rows={stockRows}
              sort={sort}
              order={order}
              onSort={toggleSort}
            />
          ) : (
            <InvestorResultsTable
              rows={investorRows}
              sort={sort}
              order={order}
              onSort={toggleSort}
            />
          )}

          <SnapshotPagination
            mode={mode}
            shown={shown}
            total={total}
            offset={offset}
            limit={limit}
            priceDate={stockRows[0]?.price_date}
            source={data.data_quality.source}
            onPrevPage={() => setOffset(Math.max(0, offset - limit))}
            onNextPage={() => setOffset(offset + limit)}
          />
        </Card>
      ) : null}

      {data ? (
        <SnapshotUnsupported
          unsupportedConditions={data.unsupported_conditions}
          notes={data.data_quality.notes}
        />
      ) : null}
    </div>
  );
}
