import type React from "react";
import { Icon } from "@iconify/react";

import { Badge, Card, SectionTitle } from "@/components/research/ui";
import { FilterInput, FilterSelect } from "./FilterControls";
import { PAGE_SIZES } from "./snapshotConstants";
import type {
  ActiveFilterItem,
  Condition,
  ConditionGroup,
  Filters,
  FlowFilter,
  Mode,
  MtdFilter,
  SortOption,
  SortOrder,
} from "./snapshotTypes";

interface SnapshotFiltersProps {
  mode: Mode;
  condition: string;
  conditions: Condition[];
  conditionGroups: ConditionGroup[];
  sectors: string[];
  filters: Filters;
  sort: string;
  order: SortOrder;
  limit: number;
  sortOptions: SortOption[];
  shown: number;
  activeFilters: ActiveFilterItem[];
  onChangeCondition: (value: string) => void;
  setFilters: React.Dispatch<React.SetStateAction<Filters>>;
  onChangeSort: (value: string) => void;
  onToggleOrder: () => void;
  onChangeLimit: (value: number) => void;
  onExportCsv: () => void;
}

export function SnapshotFilters({
  mode,
  condition,
  conditions,
  conditionGroups,
  sectors,
  filters,
  sort,
  order,
  limit,
  sortOptions,
  shown,
  activeFilters,
  onChangeCondition,
  setFilters,
  onChangeSort,
  onToggleOrder,
  onChangeLimit,
  onExportCsv,
}: SnapshotFiltersProps) {
  return (
    <Card className="p-4 sm:p-5">
      <div className="flex items-start justify-between gap-3">
        <SectionTitle
          icon="ph:sliders-horizontal-duotone"
          title="Filter lanjutan"
          subtitle={
            mode === "stocks"
              ? "Sektor, rentang kepemilikan, harga, dan arus asing"
              : "Skala portofolio dan sebaran sektor"
          }
        />
        <Badge tone={activeFilters.length ? "primary" : "neutral"}>
          {activeFilters.length} aktif
        </Badge>
      </div>

      {mode === "stocks" ? (
        <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-3">
          <div className="sm:col-span-2 xl:col-span-3">
            <label className="text-[9px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
              Kondisi
            </label>
            <select
              value={condition}
              onChange={(event) => onChangeCondition(event.target.value)}
              className="mt-1 w-full rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] px-3 py-2.5 text-[11px] text-[var(--md-sys-color-on-surface)] focus:outline-none focus:border-[var(--md-sys-color-primary)]"
            >
              {conditionGroups.length > 0
                ? conditionGroups.map((group) => {
                    const items = conditions.filter(
                      (item) => item.category === group.id
                    );
                    if (items.length === 0) return null;
                    return (
                      <optgroup key={group.id} label={group.label}>
                        {items.map((item) => (
                          <option key={item.id} value={item.id}>
                            {item.label}
                          </option>
                        ))}
                      </optgroup>
                    );
                  })
                : conditions.map((item) => (
                    <option key={item.id} value={item.id}>
                      {item.label}
                    </option>
                  ))}
            </select>
          </div>
          <FilterSelect
            label="Sektor"
            value={filters.sector}
            onChange={(value) => setFilters((prev) => ({ ...prev, sector: value }))}
            options={[
              { value: "", label: "Semua sektor" },
              ...sectors.map((item) => ({ value: item, label: item })),
            ]}
          />
          <FilterInput
            label="Min holders"
            value={filters.minHolders}
            onChange={(value) => setFilters((prev) => ({ ...prev, minHolders: value }))}
            placeholder="contoh: 5"
          />
          <FilterInput
            label="Max holders"
            value={filters.maxHolders}
            onChange={(value) => setFilters((prev) => ({ ...prev, maxHolders: value }))}
            placeholder="contoh: 50"
          />
          <FilterInput
            label="Min pengendali (%)"
            value={filters.minControlling}
            onChange={(value) => setFilters((prev) => ({ ...prev, minControlling: value }))}
            placeholder="contoh: 50"
            step="any"
          />
          <FilterInput
            label="Max pengendali (%)"
            value={filters.maxControlling}
            onChange={(value) => setFilters((prev) => ({ ...prev, maxControlling: value }))}
            placeholder="contoh: 80"
            step="any"
          />
          <FilterInput
            label="Min publik (%)"
            value={filters.minPublic}
            onChange={(value) => setFilters((prev) => ({ ...prev, minPublic: value }))}
            placeholder="contoh: 10"
            step="any"
          />
          <FilterInput
            label="Max publik (%)"
            value={filters.maxPublic}
            onChange={(value) => setFilters((prev) => ({ ...prev, maxPublic: value }))}
            placeholder="contoh: 30"
            step="any"
          />
          <FilterInput
            label="Min top holder (%)"
            value={filters.minTopHolder}
            onChange={(value) => setFilters((prev) => ({ ...prev, minTopHolder: value }))}
            placeholder="contoh: 25"
            step="any"
          />
          <FilterInput
            label="Max top holder (%)"
            value={filters.maxTopHolder}
            onChange={(value) => setFilters((prev) => ({ ...prev, maxTopHolder: value }))}
            placeholder="contoh: 75"
            step="any"
          />
          <FilterInput
            label="Min market cap (M)"
            value={filters.minMarketCap}
            onChange={(value) => setFilters((prev) => ({ ...prev, minMarketCap: value }))}
            placeholder="contoh: 10000"
            step="any"
          />
          <FilterInput
            label="Max market cap (M)"
            value={filters.maxMarketCap}
            onChange={(value) => setFilters((prev) => ({ ...prev, maxMarketCap: value }))}
            placeholder="contoh: 100000"
            step="any"
          />
          <FilterSelect
            label="Perubahan MTD"
            value={filters.mtd}
            onChange={(value) =>
              setFilters((prev) => ({ ...prev, mtd: value as MtdFilter }))
            }
            options={[
              { value: "all", label: "Semua arah" },
              { value: "up", label: "Naik MTD" },
              { value: "down", label: "Turun MTD" },
            ]}
          />
          <FilterSelect
            label="Foreign flow 5D"
            value={filters.foreignFlow}
            onChange={(value) =>
              setFilters((prev) => ({ ...prev, foreignFlow: value as FlowFilter }))
            }
            options={[
              { value: "all", label: "Semua arus" },
              { value: "inflow", label: "Net buy" },
              { value: "outflow", label: "Net sell" },
            ]}
          />
          <label className="flex items-center gap-2 rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] px-3 py-2.5 text-[11px] font-bold text-[var(--md-sys-color-on-surface)]">
            <input
              type="checkbox"
              checked={filters.hasTreasury}
              onChange={(event) =>
                setFilters((prev) => ({ ...prev, hasTreasury: event.target.checked }))
              }
              className="h-3.5 w-3.5 accent-[var(--md-sys-color-primary)]"
            />
            Hanya emiten dengan treasury
          </label>
        </div>
      ) : (
        <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-3">
          <div className="sm:col-span-2 xl:col-span-4">
            <FilterSelect
              label="Kondisi"
              value={condition}
              onChange={(value) => onChangeCondition(value)}
              options={conditions.map((item) => ({
                value: item.id,
                label: item.label,
              }))}
            />
          </div>
          <FilterInput
            label="Min emiten"
            value={filters.minTickers}
            onChange={(value) => setFilters((prev) => ({ ...prev, minTickers: value }))}
            placeholder="contoh: 5"
          />
          <FilterInput
            label="Max emiten"
            value={filters.maxTickers}
            onChange={(value) => setFilters((prev) => ({ ...prev, maxTickers: value }))}
            placeholder="contoh: 50"
          />
          <FilterInput
            label="Min porsi max (%)"
            value={filters.minMaxPct}
            onChange={(value) => setFilters((prev) => ({ ...prev, minMaxPct: value }))}
            placeholder="contoh: 5"
            step="any"
          />
          <FilterInput
            label="Min sektor"
            value={filters.minSectorCount}
            onChange={(value) => setFilters((prev) => ({ ...prev, minSectorCount: value }))}
            placeholder="contoh: 2"
          />
        </div>
      )}

      <div className="mt-4 flex flex-col lg:flex-row lg:items-end gap-3">
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-3 flex-1">
          <FilterSelect
            label="Urutkan"
            value={sort}
            onChange={onChangeSort}
            options={sortOptions.map((item) => ({ value: item.id, label: item.label }))}
          />
          <div>
            <label className="text-[9px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
              Arah
            </label>
            <button
              type="button"
              onClick={onToggleOrder}
              className="mt-1 flex w-full items-center justify-between rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] px-3 py-2.5 text-[11px] font-bold text-[var(--md-sys-color-on-surface)] hover:border-[var(--md-sys-color-primary)]/60"
            >
              {order === "asc" ? "Ascending" : "Descending"}
              <Icon
                icon={
                  order === "asc"
                    ? "ph:sort-ascending-duotone"
                    : "ph:sort-descending-duotone"
                }
                width="15"
              />
            </button>
          </div>
          <FilterSelect
            label="Baris per halaman"
            value={String(limit)}
            onChange={(value) => onChangeLimit(Number(value))}
            options={PAGE_SIZES.map((size) => ({
              value: String(size),
              label: `${size} baris`,
            }))}
          />
        </div>
        <button
          type="button"
          onClick={onExportCsv}
          disabled={shown === 0}
          className="inline-flex items-center justify-center gap-1.5 rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] px-3.5 py-2.5 text-[11px] font-bold text-[var(--md-sys-color-on-surface)] hover:border-[var(--md-sys-color-primary)]/60 disabled:cursor-not-allowed disabled:opacity-40"
        >
          <Icon icon="ph:download-simple-duotone" width="15" />
          Export CSV
        </button>
      </div>

      {activeFilters.length > 0 ? (
        <div className="mt-4 flex flex-wrap items-center gap-1.5 border-t border-[var(--md-sys-color-outline-variant)]/50 pt-3">
          {activeFilters.map((item) => (
            <button
              key={item.key}
              type="button"
              onClick={item.clear}
              title="Hapus filter"
              className="inline-flex items-center gap-1 rounded-full border border-[var(--md-sys-color-primary)]/25 bg-[var(--md-sys-color-primary-container)]/60 px-2.5 py-1 text-[9px] font-mono font-bold text-[var(--md-sys-color-on-primary-container)] hover:border-rose-400/60 hover:text-rose-600 dark:hover:text-rose-400"
            >
              {item.label}
              <Icon icon="ph:x-duotone" width="10" />
            </button>
          ))}
        </div>
      ) : null}
    </Card>
  );
}
