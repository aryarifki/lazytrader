import { Link } from "react-router-dom";
import { Icon } from "@iconify/react";

import ResearchTabs from "@/components/research/ResearchTabs";
import { Badge, Card } from "@/components/research/ui";
import { fmtNumber } from "@/components/research/format";
import { INVESTOR_PRESETS, STOCK_PRESETS } from "./snapshotConstants";
import type { ExplorerResponse, Mode } from "./snapshotTypes";

interface SnapshotHeaderCardProps {
  embedded: boolean;
  mode: Mode;
  condition: string;
  query: string;
  data: ExplorerResponse | undefined;
  total: number;
  isLoading: boolean;
  isValidating: boolean;
  onChangeMode: (nextMode: Mode) => void;
  onChangeQuery: (nextQuery: string) => void;
  onSelectPreset: (presetId: string) => void;
  onResetAll: () => void;
}

export function SnapshotHeaderCard({
  embedded,
  mode,
  condition,
  query,
  data,
  total,
  isLoading,
  isValidating,
  onChangeMode,
  onChangeQuery,
  onSelectPreset,
  onResetAll,
}: SnapshotHeaderCardProps) {
  const presets = mode === "stocks" ? STOCK_PRESETS : INVESTOR_PRESETS;

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
                <span>SCREENER</span>
              </div>
              <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight">
                Explorer Kepemilikan
              </h1>
            </div>
          </div>
          <div className="flex items-center gap-2 flex-wrap">
            <ResearchTabs />
            <Link
              to="/explorer?source=parity"
              className="inline-flex items-center gap-1.5 rounded-lg border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)] px-2.5 py-1.5 text-[10px] sm:text-[11px] font-bold text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
            >
              <Icon icon="ph:chart-polar-duotone" width="14" />
              Parity Explorer
            </Link>
          </div>
        </div>
      )}

      <Card className="p-4 sm:p-5">
        <div className="flex flex-col xl:flex-row xl:items-end justify-between gap-4">
          <div>
            <div className="flex items-center gap-2 flex-wrap">
              <Badge tone="primary">Snapshot</Badge>
              <Badge>{mode === "stocks" ? "By Saham" : "By Investor"}</Badge>
              <span className="text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                {data?.condition_label ?? "Memuat kondisi..."}
              </span>
              {isValidating && !isLoading ? (
                <span className="animate-pulse text-[9px] font-mono text-[var(--md-sys-color-primary)]">
                  memperbarui…
                </span>
              ) : null}
            </div>
            <h2 className="mt-2 text-lg sm:text-2xl font-extrabold tracking-tight">
              {data
                ? `${fmtNumber(total)} hasil`
                : "Menyiapkan screener..."}
            </h2>
            <p className="mt-1 text-[11px] text-[var(--md-sys-color-on-surface-variant)] max-w-2xl leading-relaxed">
              Filter server-side terhubung ke SQL view terindeks ownership_staging dan company_shareholders.
              Kondisi Big Player, mutasi majorholder, dan akumulasi asing kini tersedia.
            </p>
          </div>

          <div className="flex flex-col sm:flex-row gap-2 w-full xl:w-auto">
            <div className="flex rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] p-1">
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
                Saham
              </button>
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
                Investor
              </button>
            </div>

            <div className="relative">
              <Icon
                icon="ph:magnifying-glass-duotone"
                width="15"
                className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--md-sys-color-on-surface-variant)]"
              />
              <input
                value={query}
                onChange={(event) => onChangeQuery(event.target.value)}
                placeholder={mode === "stocks" ? "Cari ticker/nama..." : "Cari investor..."}
                aria-label="Cari hasil"
                className="w-full sm:w-[240px] rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] pl-9 pr-3 py-2.5 text-xs text-[var(--md-sys-color-on-surface)] placeholder:text-[var(--md-sys-color-on-surface-variant)] focus:outline-none focus:border-[var(--md-sys-color-primary)]"
              />
            </div>
          </div>
        </div>

        <div className="mt-4 flex flex-wrap items-center gap-2 border-t border-[var(--md-sys-color-outline-variant)]/50 pt-3">
          <span className="text-[9px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
            Preset
          </span>
          {presets.map((preset) => (
            <button
              key={preset.id}
              type="button"
              onClick={() => onSelectPreset(preset.id)}
              className={`inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1.5 text-[10px] font-bold transition-colors ${
                condition === preset.id
                  ? "border-[var(--md-sys-color-primary)] bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)]"
                  : "border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] hover:border-[var(--md-sys-color-primary)]/60"
              }`}
            >
              <Icon icon={preset.icon} width="13" />
              {preset.label}
            </button>
          ))}
          <button
            type="button"
            onClick={onResetAll}
            className="ml-auto inline-flex items-center gap-1.5 rounded-lg border border-[var(--md-sys-color-outline-variant)] px-2.5 py-1.5 text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] hover:border-rose-400/60 hover:text-rose-600 dark:hover:text-rose-400"
          >
            <Icon icon="ph:arrow-counter-clockwise-duotone" width="13" />
            Reset semua
          </button>
        </div>
      </Card>
    </>
  );
}
