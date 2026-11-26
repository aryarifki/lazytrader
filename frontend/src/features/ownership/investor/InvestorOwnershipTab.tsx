import { useDeferredValue, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import useSWR from "swr";

import { EmptyState, ErrorBlock, LoadingBlock } from "@/components/research/ui";
import { ownershipFetcher } from "../ownershipFetcher";
import { InvestorHoldingsTable } from "./InvestorHoldingsTable";
import { InvestorQuality } from "./InvestorQuality";
import { InvestorSearchPanel } from "./InvestorSearchPanel";
import { InvestorSummary } from "./InvestorSummary";
import { InvestorT5Table } from "./InvestorT5Table";
import type { InvestorDetail, InvestorSearchResponse } from "./investorTypes";

export default function InvestorOwnershipTab() {
  const [searchParams] = useSearchParams();
  const nameParam = searchParams.get("name")?.trim() ?? "";
  const [query, setQuery] = useState("");
  const deferredQuery = useDeferredValue(query.trim());

  const { data: search, isLoading: searchLoading } = useSWR<InvestorSearchResponse>(
    `/api/v1/ownership/investors/ct?q=${encodeURIComponent(deferredQuery)}&limit=15`,
    ownershipFetcher
  );

  const effectiveName = nameParam;
  const { data: detail, error, isLoading } = useSWR<InvestorDetail>(
    effectiveName
      ? `/api/v1/ownership/investors/ct/detail?name=${encodeURIComponent(effectiveName)}`
      : null,
    ownershipFetcher
  );

  const suggestions = useMemo(() => search?.results ?? [], [search]);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[320px_minmax(0,1fr)] gap-4">
      <InvestorSearchPanel
        query={query}
        onQueryChange={setQuery}
        deferredQuery={deferredQuery}
        search={search}
        searchLoading={searchLoading}
        suggestions={suggestions}
        effectiveName={effectiveName}
      />

      <div className="space-y-4 min-w-0">
        {isLoading ? <LoadingBlock label="Memuat portofolio investor..." /> : null}
        {error ? <ErrorBlock message={`Gagal memuat investor: ${error.message}`} /> : null}

        {detail ? (
          <>
            <InvestorSummary detail={detail} />
            <InvestorHoldingsTable detail={detail} />
            {detail.transactions.length > 0 ? (
              <InvestorT5Table transactions={detail.transactions} />
            ) : null}
            <InvestorQuality detail={detail} />
          </>
        ) : !isLoading && !error ? (
          <EmptyState
            icon="ph:user-circle-duotone"
            title="Pilih investor"
            description="Cari nama investor di panel kiri. Data portofolio, T5 lintas saham, dan G/L peak-capital dihitung dari baseline CT."
          />
        ) : null}
      </div>
    </div>
  );
}
