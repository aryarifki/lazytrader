import { useEffect, useMemo, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import useSWR from "swr";

import { buildColumns } from "./ctColumns";
import {
  FOREIGN_TYPES,
  LOCAL_TYPES,
  PORTFOLIO_RANGES,
  RETAIL_RANGES,
} from "./ctConstants";
import { CTExplorerControls } from "./CTExplorerControls";
import { CTExplorerQuality } from "./CTExplorerQuality";
import { CTExplorerResults } from "./CTExplorerResults";
import { exportCtCsv } from "./ctExport";
import {
  conditionsFor,
  defaultCondition,
  fetcher,
  toggleValue,
} from "./ctHelpers";
import type { Category, Column, ExplorerResponse, Mode, SortOrder } from "./ctTypes";

export default function CTExplorer({ embedded = false }: { embedded?: boolean } = {}) {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const modeParam = searchParams.get("mode");
  const initialMode: Mode = modeParam === "investors" ? "investors" : "stocks";
  const [mode, setMode] = useState<Mode>(initialMode);
  const [category, setCategory] = useState<Category>("kepemilikan");
  const [condition, setCondition] = useState("bertambah");
  const [ownershipRanges, setOwnershipRanges] = useState<string[]>([">10"]);
  const [sidRanges, setSidRanges] = useState<string[]>([">10"]);
  const [mktcapTier, setMktcapTier] = useState("big");
  const [transaksiRange, setTransaksiRange] = useState("<5M");
  const [origin, setOrigin] = useState("Local");
  const [localTypes, setLocalTypes] = useState<string[]>(LOCAL_TYPES.map((item) => item.id));
  const [foreignTypes, setForeignTypes] = useState<string[]>(FOREIGN_TYPES.map((item) => item.id));
  const [retailRanges, setRetailRanges] = useState<string[]>(RETAIL_RANGES);
  const [portfolioRanges, setPortfolioRanges] = useState<string[]>(PORTFOLIO_RANGES);
  const [sort, setSort] = useState("default");
  const [order, setOrder] = useState<SortOrder>("desc");
  const [limit, setLimit] = useState(50);
  const [offset, setOffset] = useState(0);

  useEffect(() => {
    if (modeParam === "investors" && mode !== "investors") {
      setMode("investors");
    } else if (modeParam === "stocks" && mode !== "stocks") {
      setMode("stocks");
    }
  }, [modeParam]);

  const isSidCondition = condition === "sidNaik" || condition === "sidTurun";
  const isPortfolioCondition =
    condition === "portofolioNaik" || condition === "portofolioTurun";

  const queryString = useMemo(() => {
    const params = new URLSearchParams();
    params.set("mode", mode);
    params.set("category", category);
    params.set("condition", condition);
    params.set("mktcap_tier", mktcapTier);
    params.set("transaksi_range", transaksiRange);
    params.set("sort", sort);
    params.set("order", order);
    params.set("limit", String(limit));
    params.set("offset", String(offset));

    if (category === "kepemilikan") {
      params.set("ranges", (isSidCondition ? sidRanges : ownershipRanges).join(","));
    }
    if (category === "advanced" && condition === "retailRendah") {
      params.set("retail_ranges", retailRanges.join(","));
    }
    if (category === "advanced" && isPortfolioCondition) {
      params.set("portfolio_ranges", portfolioRanges.join(","));
    }
    if (mode === "investors") {
      params.set("investor_origin", origin);
      params.set(
        "investor_types",
        (origin === "Local" ? localTypes : foreignTypes).join(",")
      );
    }
    return params.toString();
  }, [
    mode,
    category,
    condition,
    mktcapTier,
    transaksiRange,
    sort,
    order,
    limit,
    offset,
    isSidCondition,
    isPortfolioCondition,
    sidRanges,
    ownershipRanges,
    retailRanges,
    portfolioRanges,
    origin,
    localTypes,
    foreignTypes,
  ]);

  const { data, error, isLoading, isValidating } = useSWR<ExplorerResponse>(
    `/api/v1/ownership/explorer/ct?${queryString}`,
    fetcher
  );

  const conditions = useMemo(() => {
    const fallback = conditionsFor(mode, category);
    if (category !== "advanced") return fallback;
    const apiConditions = data?.conditions ?? [];
    return fallback.map((item) => {
      const fromApi = apiConditions.find((candidate) => candidate.id === item.id);
      return fromApi ?? item;
    });
  }, [mode, category, data?.conditions]);

  const rows = data?.results ?? [];
  const total = data?.total ?? 0;
  const shown = data?.count ?? rows.length;

  const changeMode = (nextMode: Mode) => {
    setMode(nextMode);
    setCondition(defaultCondition(nextMode, category));
    setSort("default");
    setOffset(0);
    const params = new URLSearchParams(window.location.search);
    params.set("mode", nextMode);
    navigate(`/explorer?${params.toString()}`, {
      replace: true,
      preventScrollReset: true,
    });
  };

  const changeCategory = (nextCategory: Category) => {
    setCategory(nextCategory);
    setCondition(defaultCondition(mode, nextCategory));
    setSort("default");
    setOffset(0);
  };

  const changeCondition = (nextCondition: string) => {
    setCondition(nextCondition);
    setSort("default");
    setOffset(0);
  };

  const resetAll = () => {
    if (category !== "advanced") setCondition("bertambah");
    setOwnershipRanges([">10"]);
    setSidRanges([">10"]);
    setMktcapTier("big");
    setTransaksiRange("<5M");
    setOrigin("Local");
    setLocalTypes(LOCAL_TYPES.map((item) => item.id));
    setForeignTypes(FOREIGN_TYPES.map((item) => item.id));
    setRetailRanges(RETAIL_RANGES);
    setPortfolioRanges(PORTFOLIO_RANGES);
    setSort("default");
    setOrder("desc");
    setOffset(0);
  };

  const toggleSort = (column: Column) => {
    if (!column.sort) return;
    setOffset(0);
    if (sort === column.sort) {
      setOrder((previous) => (previous === "asc" ? "desc" : "asc"));
      return;
    }
    setSort(column.sort);
    setOrder(column.sort === "ticker" || column.sort === "name" ? "asc" : "desc");
  };

  const columns = useMemo(
    () => buildColumns(mode, category, condition, data?.period),
    [mode, category, condition, data?.period]
  );

  return (
    <div className={embedded ? "space-y-4" : "max-w-7xl mx-auto space-y-4"}>
      <CTExplorerControls
        embedded={embedded}
        mode={mode}
        category={category}
        condition={condition}
        conditions={conditions}
        data={data}
        total={total}
        isLoading={isLoading}
        isValidating={isValidating}
        isSidCondition={isSidCondition}
        isPortfolioCondition={isPortfolioCondition}
        sidRanges={sidRanges}
        ownershipRanges={ownershipRanges}
        transaksiRange={transaksiRange}
        retailRanges={retailRanges}
        portfolioRanges={portfolioRanges}
        mktcapTier={mktcapTier}
        origin={origin}
        localTypes={localTypes}
        foreignTypes={foreignTypes}
        limit={limit}
        hasRows={rows.length > 0}
        onChangeMode={changeMode}
        onChangeCategory={changeCategory}
        onChangeCondition={changeCondition}
        onToggleSidRange={(value) => {
          setOffset(0);
          setSidRanges((previous) => toggleValue(previous, value));
        }}
        onToggleOwnershipRange={(value) => {
          setOffset(0);
          setOwnershipRanges((previous) => toggleValue(previous, value));
        }}
        onSelectTransaksiRange={(value) => {
          setOffset(0);
          setTransaksiRange(value);
        }}
        onToggleRetailRange={(value) => {
          setOffset(0);
          setRetailRanges((previous) => toggleValue(previous, value));
        }}
        onTogglePortfolioRange={(value) => {
          setOffset(0);
          setPortfolioRanges((previous) => toggleValue(previous, value));
        }}
        onSelectMktcapTier={(value) => {
          setOffset(0);
          setMktcapTier(value);
        }}
        onSelectOrigin={(value) => {
          setOffset(0);
          setOrigin(value);
        }}
        onToggleOriginType={(value) => {
          setOffset(0);
          if (origin === "Local") {
            setLocalTypes((previous) => toggleValue(previous, value));
          } else {
            setForeignTypes((previous) => toggleValue(previous, value));
          }
        }}
        onResetAll={resetAll}
        onExportCsv={() => exportCtCsv(mode, category, rows)}
        onChangeLimit={(nextLimit) => {
          setOffset(0);
          setLimit(nextLimit);
        }}
      />

      <CTExplorerResults
        data={data}
        error={error}
        isLoading={isLoading}
        rows={rows}
        columns={columns}
        mode={mode}
        category={category}
        shown={shown}
        total={total}
        sort={sort}
        order={order}
        offset={offset}
        limit={limit}
        onToggleSort={toggleSort}
        onPrevPage={() => setOffset((previous) => Math.max(0, previous - limit))}
        onNextPage={() => setOffset((previous) => previous + limit)}
      />

      <CTExplorerQuality dataQuality={data?.data_quality} />
    </div>
  );
}
