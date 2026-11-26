import React from "react";
import { Icon } from "@iconify/react";

import { Badge, Card } from "@/components/research/ui";
import { fmtDate } from "@/components/research/format";
import { StockBadges } from "@/components/research/StockBadges";
import { QUICK_TICKERS } from "./stockHelpers";
import type { CtHolder, SearchResult, StockOwnership } from "./stockTypes";

interface StockSearchHeaderProps {
  ticker: string;
  data: StockOwnership | undefined;
  isCtBaseline: boolean;
  displayHolders: CtHolder[] | null;
  query: string;
  onQueryChange: (value: string) => void;
  searchResults: SearchResult[] | undefined;
  onSelectTicker: (ticker: string) => void;
  onSearchKeyDown: (event: React.KeyboardEvent<HTMLInputElement>) => void;
}

export function StockSearchHeader({
  ticker,
  data,
  isCtBaseline,
  displayHolders,
  query,
  onQueryChange,
  searchResults,
  onSelectTicker,
  onSearchKeyDown,
}: StockSearchHeaderProps) {
  return (
    <Card className="p-4 sm:p-5">
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
        <div className="min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <Badge tone="primary">{ticker}</Badge>
            <StockBadges
              badges={data?.badges ?? data?.profile.badges}
              badge={data?.badge ?? data?.profile.badge}
              conglomerate={data?.conglomerate ?? data?.profile.conglomerate}
              size="md"
            />
            {data?.staging_holders?.report_date ? (
              <Badge tone="up">KSEI Staging ({fmtDate(data.staging_holders.report_date)})</Badge>
            ) : (
              <Badge tone={isCtBaseline ? "primary" : data?.data_quality.monthly_composition ? "up" : "warning"}>
                {isCtBaseline
                  ? "Baseline CT (statis)"
                  : data?.data_quality.monthly_composition
                    ? "KSEI bulanan"
                    : "Snapshot pemegang saham"}
              </Badge>
            )}
            <span className="text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
              {data?.staging_holders?.report_date
                ? `KSEI ${fmtDate(data.staging_holders.report_date)} · Data Mutakhir`
                : isCtBaseline
                  ? `CT ${data?.as_of.ct_baseline ?? "baseline"} · arsip`
                  : data?.as_of.monthly
                    ? `KSEI ${fmtDate(data.as_of.monthly)}`
                    : "Snapshot saat ini"}
            </span>
          </div>
          <h2 className="text-lg sm:text-2xl font-extrabold tracking-tight mt-2 truncate">
            {data?.name ?? "Memuat emiten..."}
          </h2>
          <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-0.5">
            {data?.sector ?? "—"} ·{" "}
            {data?.staging_holders?.holders?.length
              ? `${data.staging_holders.holders.length} pemegang saham (KSEI ${fmtDate(data.staging_holders.report_date)})`
              : displayHolders
                ? `${displayHolders.length} baris pemegang >1% (CT)`
                : `${data?.summary.holder_count ?? 0} baris pemegang saham tercatat`}
          </p>
        </div>

        <div className="w-full lg:w-[360px] space-y-2">
          <div className="relative">
            <Icon
              icon="ph:magnifying-glass-duotone"
              width="16"
              className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--md-sys-color-on-surface-variant)]"
            />
            <input
              value={query}
              onChange={(event) => onQueryChange(event.target.value)}
              onKeyDown={onSearchKeyDown}
              placeholder="Cari kode atau nama emiten..."
              aria-label="Cari emiten"
              className="w-full rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] pl-9 pr-3 py-2.5 text-xs font-mono uppercase text-[var(--md-sys-color-on-surface)] placeholder:normal-case placeholder:font-sans placeholder:text-[var(--md-sys-color-on-surface-variant)] focus:outline-none focus:border-[var(--md-sys-color-primary)]"
            />
            {query && searchResults ? (
              <div className="absolute z-30 mt-1.5 w-full overflow-hidden rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)] shadow-lg">
                {searchResults.length === 0 ? (
                  <div className="px-3 py-3 text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
                    Kode “{query}” tidak ditemukan.
                  </div>
                ) : (
                  searchResults.map((item) => (
                    <button
                      key={item.ticker}
                      type="button"
                      onClick={() => onSelectTicker(item.ticker)}
                      className="flex w-full items-center justify-between gap-3 px-3 py-2.5 text-left hover:bg-[var(--md-sys-color-surface-container-high)] transition-colors"
                    >
                      <span className="min-w-0">
                        <span className="block text-xs font-mono font-bold text-[var(--md-sys-color-primary)]">
                          {item.ticker}
                        </span>
                        <span className="block text-[10px] text-[var(--md-sys-color-on-surface-variant)] truncate">
                          {item.name}
                        </span>
                      </span>
                      <span className="text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)] shrink-0">
                        {item.sector}
                      </span>
                    </button>
                  ))
                )}
              </div>
            ) : null}
          </div>

          <div className="flex items-center gap-1.5 overflow-x-auto scrollbar-none">
            {QUICK_TICKERS.map((item) => (
              <button
                key={item}
                type="button"
                onClick={() => onSelectTicker(item)}
                className={`rounded-md border px-2 py-1 text-[9px] font-mono font-bold transition-colors ${
                  ticker === item
                    ? "border-[var(--md-sys-color-primary)] bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)]"
                    : "border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                }`}
              >
                {item}
              </button>
            ))}
          </div>
        </div>
      </div>
    </Card>
  );
}
