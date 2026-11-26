import { Link } from "react-router-dom";
import { Icon } from "@iconify/react";

import ResearchTabs from "@/components/research/ResearchTabs";
import { Badge, Card } from "@/components/research/ui";
import { fmtNumber } from "@/components/research/format";
import {
  CONSTANT_CONDITION_LABELS,
  FOREIGN_TYPES,
  LOCAL_TYPES,
  MKTCAP_TIERS,
  OWNERSHIP_RANGE_OPTIONS,
  PORTFOLIO_RANGE_OPTIONS,
  RETAIL_RANGE_OPTIONS,
  SID_RANGE_OPTIONS,
  VALUE_RANGE_OPTIONS,
} from "./ctConstants";
import type { Category, ConditionOption, ExplorerResponse, Mode } from "./ctTypes";
import { PillGroup } from "./PillGroup";

interface CTExplorerControlsProps {
  embedded: boolean;
  mode: Mode;
  category: Category;
  condition: string;
  conditions: ConditionOption[];
  data: ExplorerResponse | undefined;
  total: number;
  isLoading: boolean;
  isValidating: boolean;
  isSidCondition: boolean;
  isPortfolioCondition: boolean;
  sidRanges: string[];
  ownershipRanges: string[];
  transaksiRange: string;
  retailRanges: string[];
  portfolioRanges: string[];
  mktcapTier: string;
  origin: string;
  localTypes: string[];
  foreignTypes: string[];
  limit: number;
  hasRows: boolean;
  onChangeMode: (nextMode: Mode) => void;
  onChangeCategory: (nextCategory: Category) => void;
  onChangeCondition: (nextCondition: string) => void;
  onToggleSidRange: (value: string) => void;
  onToggleOwnershipRange: (value: string) => void;
  onSelectTransaksiRange: (value: string) => void;
  onToggleRetailRange: (value: string) => void;
  onTogglePortfolioRange: (value: string) => void;
  onSelectMktcapTier: (value: string) => void;
  onSelectOrigin: (value: string) => void;
  onToggleOriginType: (value: string) => void;
  onResetAll: () => void;
  onExportCsv: () => void;
  onChangeLimit: (nextLimit: number) => void;
}

export function CTExplorerControls(props: CTExplorerControlsProps) {
  const {
    embedded, mode, category, condition, conditions, data, total,
    isLoading, isValidating, isSidCondition, isPortfolioCondition,
    sidRanges, ownershipRanges, transaksiRange, retailRanges, portfolioRanges,
    mktcapTier, origin, localTypes, foreignTypes, limit, hasRows,
    onChangeMode, onChangeCategory, onChangeCondition,
    onToggleSidRange, onToggleOwnershipRange, onSelectTransaksiRange,
    onToggleRetailRange, onTogglePortfolioRange, onSelectMktcapTier,
    onSelectOrigin, onToggleOriginType, onResetAll, onExportCsv, onChangeLimit,
  } = props;
  const originTypeOptions = origin === "Local" ? LOCAL_TYPES : FOREIGN_TYPES;
  const originTypes = origin === "Local" ? localTypes : foreignTypes;

  return (
    <>
      {!embedded && (
        <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-3">
          <div className="flex items-center gap-3.5">
            <div className="w-11 h-11 sm:w-12 sm:h-12 rounded-2xl bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0 border border-[var(--md-sys-color-outline-variant)] shadow-2xs">
              <Icon icon="ph:funnel-duotone" width="26" />
            </div>
            <div>
              <div className="inline-flex items-center gap-2 text-[9px] font-mono font-bold text-[var(--md-sys-color-primary)] uppercase tracking-wider mb-0.5">
                <span>OWNERSHIP EXPLORER</span>
                <span className="opacity-40">•</span>
                <span>CT PARITY</span>
              </div>
              <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight">
                Explorer Kepemilikan
              </h1>
            </div>
          </div>
          <div className="flex items-center gap-2 flex-wrap">
            <ResearchTabs />
            <Link
              to="/explorer?source=snapshot"
              className="inline-flex items-center gap-1.5 rounded-lg border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)] px-2.5 py-1.5 text-[10px] sm:text-[11px] font-bold text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
            >
              <Icon icon="ph:database-duotone" width="14" />
              Snapshot
            </Link>
          </div>
        </div>
      )}

      <Card className="p-4 sm:p-5 space-y-4">
        <div className="flex flex-col xl:flex-row xl:items-end justify-between gap-4">
          <div>
            <div className="flex items-center gap-2 flex-wrap">
              <Badge tone="primary">CT</Badge>
              <Badge>{mode === "stocks" ? "By Stocks" : "By Investors"}</Badge>
              <span className="text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                {data?.condition_label ?? CONSTANT_CONDITION_LABELS[condition] ?? condition}
              </span>
              {isValidating && !isLoading ? (
                <span className="animate-pulse text-[9px] font-mono text-[var(--md-sys-color-primary)]">
                  memperbarui…
                </span>
              ) : null}
            </div>
            <h2 className="mt-2 text-lg sm:text-2xl font-extrabold tracking-tight">
              {data ? `${fmtNumber(total)} hasil` : "Menyiapkan explorer..."}
            </h2>
            <p className="mt-1 text-[11px] text-[var(--md-sys-color-on-surface-variant)] max-w-2xl leading-relaxed">
              Screener server-side dari baseline CT: kepemilikan, transaksi T5,
              SID, market cap, dan G/L peak-capital.
              {data?.period?.last_update_date
                ? ` Tanggal update ${data.period.last_update_date}.`
                : ""}
            </p>
          </div>

          <div className="flex flex-col sm:flex-row gap-2 w-full xl:w-auto">
            <div className="flex rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] p-1">
              <button
                type="button"
                onClick={() => onChangeMode("investors")}
                className={`flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-[11px] font-bold transition-all ${
                  mode === "investors"
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                    : "text-[var(--md-sys-color-on-surface-variant)]"
                }`}
              >
                <Icon icon="ph:users-three-duotone" width="14" />
                By Investors
              </button>
              <button
                type="button"
                onClick={() => onChangeMode("stocks")}
                className={`flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-[11px] font-bold transition-all ${
                  mode === "stocks"
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                    : "text-[var(--md-sys-color-on-surface-variant)]"
                }`}
              >
                <Icon icon="ph:buildings-duotone" width="14" />
                By Stocks
              </button>
            </div>

            {category === "advanced" ? (
              <select
                value={condition}
                onChange={(event) => onChangeCondition(event.target.value)}
                aria-label="Kondisi explorer"
                className="rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] px-3 py-2 text-[11px] font-semibold text-[var(--md-sys-color-on-surface)] focus:outline-none focus:border-[var(--md-sys-color-primary)]"
              >
                {conditions.map((option) => (
                  <option key={option.id} value={option.id}>
                    {option.label}
                  </option>
                ))}
              </select>
            ) : (
              <div className="flex flex-wrap gap-1.5">
                {conditions.map((option) => {
                  const active = condition === option.id;
                  const positive = option.id === "bertambah" || option.id === "sidNaik";
                  return (
                    <button
                      key={option.id}
                      type="button"
                      onClick={() => onChangeCondition(option.id)}
                      className={`rounded-lg border px-3 py-1.5 text-[11px] font-bold transition-all ${
                        active
                          ? positive
                            ? "border-emerald-500/40 bg-emerald-500/15 text-emerald-700 dark:text-emerald-400"
                            : "border-rose-500/40 bg-rose-500/15 text-rose-700 dark:text-rose-400"
                          : "border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                      }`}
                    >
                      {option.label}
                    </button>
                  );
                })}
              </div>
            )}
          </div>
        </div>

        <div className="flex flex-wrap gap-1.5">
          {(
            [
              { id: "kepemilikan", label: "Kepemilikan", icon: "ph:chart-pie-slice-duotone" },
              { id: "transaksi", label: "Transaksi", icon: "ph:swap-duotone" },
              { id: "advanced", label: "Advanced", icon: "ph:sparkle-duotone" },
            ] as Array<{ id: Category; label: string; icon: string }>
          ).map((item) => {
            const active = category === item.id;
            return (
              <button
                key={item.id}
                type="button"
                onClick={() => onChangeCategory(item.id)}
                className={`flex items-center gap-1.5 rounded-lg border px-3 py-1.5 text-[11px] font-bold transition-all ${
                  active
                    ? "border-[var(--md-sys-color-primary)] bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                    : "border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                }`}
              >
                <Icon icon={item.icon} width="14" />
                {item.label}
              </button>
            );
          })}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 border-t border-[var(--md-sys-color-outline-variant)]/60 pt-4">
          <div className="space-y-3">
            {category === "kepemilikan" ? (
              <PillGroup
                label={isSidCondition ? "Perubahan SID" : "Perubahan Kepemilikan"}
                options={isSidCondition ? SID_RANGE_OPTIONS : OWNERSHIP_RANGE_OPTIONS}
                selected={isSidCondition ? sidRanges : ownershipRanges}
                onToggle={(value) => {
                  if (isSidCondition) {
                    onToggleSidRange(value);
                  } else {
                    onToggleOwnershipRange(value);
                  }
                }}
              />
            ) : null}

            {category === "transaksi" ? (
              <PillGroup
                label="Nilai Transaksi (M = Miliar Rupiah)"
                options={VALUE_RANGE_OPTIONS}
                selected={[transaksiRange]}
                onToggle={onSelectTransaksiRange}
              />
            ) : null}

            {category === "advanced" && condition === "retailRendah" ? (
              <PillGroup
                label="Rentang % Retail"
                options={RETAIL_RANGE_OPTIONS}
                selected={retailRanges}
                onToggle={onToggleRetailRange}
              />
            ) : null}

            {category === "advanced" && isPortfolioCondition ? (
              <PillGroup
                label="Rentang G/L Portofolio"
                options={PORTFOLIO_RANGE_OPTIONS}
                selected={portfolioRanges}
                onToggle={onTogglePortfolioRange}
              />
            ) : null}

            {mode === "stocks" && condition !== "retailRendah" ? (
              <PillGroup
                label="Market Cap"
                options={MKTCAP_TIERS}
                selected={[mktcapTier]}
                onToggle={onSelectMktcapTier}
              />
            ) : null}
          </div>

          <div className="space-y-3">
            {mode === "investors" ? (
              <>
                <PillGroup
                  label="Asal Investor"
                  options={[
                    { id: "Local", label: "Local" },
                    { id: "Foreign", label: "Foreign" },
                  ]}
                  selected={[origin]}
                  onToggle={onSelectOrigin}
                />
                <PillGroup
                  label={`Tipe Investor (${origin})`}
                  options={originTypeOptions}
                  selected={originTypes}
                  onToggle={onToggleOriginType}
                />
              </>
            ) : null}

            <div className="flex flex-wrap items-center gap-2">
              <button
                type="button"
                onClick={onResetAll}
                className="inline-flex items-center gap-1.5 rounded-lg border border-[var(--md-sys-color-outline-variant)] px-3 py-1.5 text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
              >
                <Icon icon="ph:arrow-counter-clockwise-duotone" width="13" />
                Reset filter
              </button>
              <button
                type="button"
                onClick={onExportCsv}
                disabled={!hasRows}
                className="inline-flex items-center gap-1.5 rounded-lg border border-[var(--md-sys-color-outline-variant)] px-3 py-1.5 text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] disabled:opacity-40"
              >
                <Icon icon="ph:download-simple-duotone" width="13" />
                CSV
              </button>
              <label className="inline-flex items-center gap-1.5 text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                Baris
                <select
                  value={limit}
                  onChange={(event) => onChangeLimit(Number(event.target.value))}
                  className="rounded-md border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] px-1.5 py-1 text-[10px]"
                >
                  {[25, 50, 100, 200].map((value) => (
                    <option key={value} value={value}>
                      {value}
                    </option>
                  ))}
                </select>
              </label>
            </div>
          </div>
        </div>
      </Card>
    </>
  );
}
