import React, { useDeferredValue, useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import useSWR from "swr";

import { ErrorBlock, LoadingBlock } from "@/components/research/ui";
import { fmtDate, fmtNumber } from "@/components/research/format";
import { ownershipFetcher } from "../ownershipFetcher";
import { StockCompositionCard } from "./StockCompositionCard";
import { StockHoldersCard } from "./StockHoldersCard";
import { StockMonthlyTrendCard } from "./StockMonthlyTrendCard";
import { StockMovementsCard } from "./StockMovementsCard";
import { StockQualityCard } from "./StockQualityCard";
import { StockSearchHeader } from "./StockSearchHeader";
import { StockStatCards } from "./StockStatCards";
import type { Composition, SearchResult, StockOwnership } from "./stockTypes";

export default function StockOwnershipTab() {
  const [searchParams] = useSearchParams();
  const tickerParam = searchParams.get("ticker")?.trim().toUpperCase() ?? "";
  const [ticker, setTicker] = useState(tickerParam || "BBCA");
  const [query, setQuery] = useState("");
  const deferredQuery = useDeferredValue(query.trim());

  useEffect(() => {
    if (tickerParam && tickerParam !== ticker) setTicker(tickerParam);
  }, [tickerParam]);

  const { data: searchResults } = useSWR<SearchResult[]>(
    deferredQuery.length > 0
      ? `/api/v1/stock/search?q=${encodeURIComponent(deferredQuery)}&limit=6`
      : null,
    ownershipFetcher
  );

  const { data, error, isLoading } = useSWR<StockOwnership>(
    `/api/v1/ownership/stock/${encodeURIComponent(ticker)}`,
    ownershipFetcher
  );

  const { data: parityData } = useSWR<any>(
    ticker ? `/api/v1/ownership/staging/parity/${encodeURIComponent(ticker)}` : null,
    ownershipFetcher
  );

  const stagingCompAsComposition: Composition | null = useMemo(() => {
    if (!data?.staging_composition?.items?.length) return null;
    const items = data.staging_composition.items;
    const reportDate = data.staging_composition.report_date;
    const total = items.reduce((acc, it) => acc + (it.pct ?? 0), 0);
    return {
      basis: "staging_normalized",
      period_label: reportDate ? fmtDate(reportDate) : "KSEI / IDX",
      period_date: reportDate,
      total_pct: total,
      items: items.map((it) => ({
        key: it.kind || it.label,
        label: it.label,
        pct: it.pct ?? 0,
        color: it.color || "#3A7BD5",
        description: it.shares != null ? `${fmtNumber(it.shares)} lembar` : "",
      })),
      note: "Data komposisi resmi KSEI / IDX (ownership_staging.v_composition_normalized) yang telah dinormalisasi terhadap total listed shares.",
    };
  }, [data?.staging_composition]);

  const displayComposition =
    stagingCompAsComposition ?? data?.ct_composition ?? data?.composition ?? null;
  const sortedComposition = useMemo(() => {
    if (!displayComposition) return [];
    return [...displayComposition.items].sort((a, b) => b.pct - a.pct);
  }, [displayComposition]);

  const history = data?.ct_history?.length ? data.ct_history : data?.history ?? [];
  const displayHolders = data?.ct_holders?.length ? data.ct_holders : null;
  const transactions = data?.transactions ?? [];
  const freeFloat = data?.profile.ct_free_float_pct ?? data?.profile.free_float_pct;
  const isCtBaseline = Boolean(data?.data_quality.ct_baseline);

  const hasStagingHolders = Boolean(data?.staging_holders?.holders?.length);
  const hasCtHolders = Boolean(data?.ct_holders?.length);
  const [holderTab, setHolderTab] = useState<"staging" | "ct">("staging");

  const hasStagingMovements = Boolean(data?.staging_movements?.length);
  const hasCtTransactions = Boolean(data?.transactions?.length);
  const [movementTab, setMovementTab] = useState<"staging" | "ct">("staging");

  const activeHolderTab = (holderTab === "staging" && hasStagingHolders) || !hasCtHolders ? (hasStagingHolders ? "staging" : "ct") : "ct";
  const activeMovementTab = (movementTab === "staging" && hasStagingMovements) || !hasCtTransactions ? (hasStagingMovements ? "staging" : "ct") : "ct";

  const [showMonthlyTable, setShowMonthlyTable] = useState(true);

  const selectTicker = (value: string) => {
    setTicker(value.toUpperCase());
    setQuery("");
  };

  const handleSearchKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter" && searchResults && searchResults.length > 0) {
      event.preventDefault();
      selectTicker(searchResults[0].ticker);
    }
  };

  return (
    <div className="space-y-4">
      <StockSearchHeader
        ticker={ticker}
        data={data}
        isCtBaseline={isCtBaseline}
        displayHolders={displayHolders}
        query={query}
        onQueryChange={setQuery}
        searchResults={searchResults}
        onSelectTicker={selectTicker}
        onSearchKeyDown={handleSearchKeyDown}
      />

      {isLoading ? <LoadingBlock label="Memuat profil kepemilikan..." /> : null}
      {error ? <ErrorBlock message={`Gagal memuat ${ticker}: ${error.message}`} /> : null}

      {data ? (
        <>
          <StockStatCards
            data={data}
            isCtBaseline={isCtBaseline}
            freeFloat={freeFloat}
          />

          <div className="grid grid-cols-1 xl:grid-cols-[1.05fr_0.95fr] gap-4">
            <StockCompositionCard
              displayComposition={displayComposition}
              sortedComposition={sortedComposition}
            />

            <StockMonthlyTrendCard
              data={data}
              history={history}
              showMonthlyTable={showMonthlyTable}
              onToggleMonthlyTable={() => setShowMonthlyTable(!showMonthlyTable)}
            />
          </div>

          <StockHoldersCard
            data={data}
            activeHolderTab={activeHolderTab}
            onSelectHolderTab={setHolderTab}
            hasStagingHolders={hasStagingHolders}
            hasCtHolders={hasCtHolders}
            displayHolders={displayHolders}
          />

          <div className="space-y-4">
            <StockMovementsCard
              data={data}
              transactions={transactions}
              activeMovementTab={activeMovementTab}
              onSelectMovementTab={setMovementTab}
              hasStagingMovements={hasStagingMovements}
              hasCtTransactions={hasCtTransactions}
            />

            <StockQualityCard
              data={data}
              isCtBaseline={isCtBaseline}
              displayHolders={displayHolders}
              parityData={parityData}
            />
          </div>
        </>
      ) : null}
    </div>
  );
}
