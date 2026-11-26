import { Icon } from "@iconify/react";
import { Badge, Card } from "@/components/research/ui";
import { fmtDate, fmtNumber } from "@/components/research/format";
import type { GroupDetail, GroupsResponse } from "./groupTypes";

interface GroupHeaderCardProps {
  groupsData: GroupsResponse;
  detail: GroupDetail | undefined;
  effectiveGroup: string;
  minPct: number;
  isSyncing: boolean;
  syncStatus: string | null;
  onSelectGroup: (groupId: string) => void;
  onChangeMinPct: (minPct: number) => void;
  onSyncCatalog: () => void;
  onDismissSyncStatus: () => void;
}

export function GroupHeaderCard({
  groupsData,
  detail,
  effectiveGroup,
  minPct,
  isSyncing,
  syncStatus,
  onSelectGroup,
  onChangeMinPct,
  onSyncCatalog,
  onDismissSyncStatus,
}: GroupHeaderCardProps) {
  return (
    <Card className="p-4 sm:p-5">
      <div className="flex flex-col lg:flex-row lg:items-end justify-between gap-4">
        <div className="min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <Badge tone="primary">Grup terpilih</Badge>
            <Badge>{groupsData.coverage.mapped_groups} grup terpetakan</Badge>
            {detail?.group.ct_baseline || groupsData.coverage.ct_baseline ? (
              <Badge tone="up">Baseline CT (statis)</Badge>
            ) : null}
            {detail?.group.catalog_code ? (
              <Badge tone="neutral">Katalog ({detail.group.catalog_code})</Badge>
            ) : null}
            <span className="text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
              Harga per {fmtDate(groupsData.as_of)}
            </span>
          </div>

          <div className="mt-2.5 flex items-center gap-3">
            {detail?.group.logo_url ? (
              <div className="w-12 h-12 rounded-xl bg-white/95 p-1.5 flex items-center justify-center border border-[var(--md-sys-color-outline-variant)] shadow-2xs shrink-0 overflow-hidden">
                <img
                  src={detail.group.logo_url}
                  alt={detail.group.name}
                  className="max-h-full max-w-full object-contain"
                  onError={(e) => {
                    (e.target as HTMLElement).style.display = "none";
                  }}
                />
              </div>
            ) : null}
            <div className="min-w-0">
              <h2 className="text-lg sm:text-2xl font-extrabold tracking-tight truncate">
                {detail?.group.name ?? groupsData.top_groups[0]?.name ?? "Konglomerasi"}
              </h2>
              {detail?.group.catalog_code ? (
                <span className="text-[10px] font-mono text-[var(--md-sys-color-primary)] font-semibold">
                  Katalog: {detail.group.catalog_code}
                </span>
              ) : null}
            </div>
          </div>

          <p className="mt-1 text-[11px] text-[var(--md-sys-color-on-surface-variant)] max-w-2xl leading-relaxed">
            {detail?.group.description ?? groupsData.coverage.note}
          </p>
          {detail ? (
            <div className="mt-2 flex flex-wrap items-center gap-2 text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
              <span>Hub: {detail.group.hub_entity || "—"}</span>
              <span className="opacity-40">•</span>
              <span>{detail.group.covered_tickers}/{detail.group.ticker_count} emiten tercakup</span>
              <span className="opacity-40">•</span>
              <span>
                {detail.group.ct_baseline
                  ? `${fmtNumber(detail.group.sid_total)} SID tercakup`
                  : `${detail.group.monthly_covered} emiten punya snapshot KSEI`}
              </span>
            </div>
          ) : null}
        </div>

        <div className="flex flex-col sm:flex-row gap-2 w-full lg:w-auto items-stretch sm:items-end">
          <label className="block">
            <span className="block text-[9px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)] mb-1">
              Grup
            </span>
            <select
              value={effectiveGroup}
              onChange={(event) => onSelectGroup(event.target.value)}
              className="w-full sm:w-[280px] rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] px-3 py-2.5 text-xs text-[var(--md-sys-color-on-surface)] focus:outline-none focus:border-[var(--md-sys-color-primary)]"
            >
              {groupsData.groups.map((group) => (
                <option key={group.id} value={group.id}>
                  {group.name} · {group.ticker_count} saham
                </option>
              ))}
            </select>
          </label>
          <label className="block">
            <span className="block text-[9px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)] mb-1">
              Min. porsi
            </span>
            <select
              value={minPct}
              onChange={(event) => onChangeMinPct(Number(event.target.value))}
              className="w-full sm:w-[120px] rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] px-3 py-2.5 text-xs font-mono text-[var(--md-sys-color-on-surface)] focus:outline-none focus:border-[var(--md-sys-color-primary)]"
            >
              {[0.5, 1, 2, 5].map((value) => (
                <option key={value} value={value}>
                  ≥ {value}%
                </option>
              ))}
            </select>
          </label>
          <button
            type="button"
            onClick={onSyncCatalog}
            disabled={isSyncing}
            className="px-3 py-2.5 rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-primary-container)] text-xs font-medium text-[var(--md-sys-color-on-surface)] flex items-center justify-center gap-1.5 transition-colors disabled:opacity-50 h-[38px] cursor-pointer"
            title="Sinkronisasi Katalog Konglomerasi Otomatis"
          >
            <Icon
              icon={isSyncing ? "ph:spinner-gap-bold" : "ph:arrows-clockwise-bold"}
              width="15"
              className={isSyncing ? "animate-spin text-[var(--md-sys-color-primary)]" : ""}
            />
            <span>{isSyncing ? "Syncing..." : "Sync Katalog"}</span>
          </button>
        </div>
      </div>
      {syncStatus ? (
        <div className="mt-3 text-[11px] font-mono px-3 py-1.5 rounded-lg bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] flex items-center justify-between border border-[var(--md-sys-color-outline-variant)]/60">
          <span>{syncStatus}</span>
          <button
            type="button"
            onClick={onDismissSyncStatus}
            className="text-[10px] text-[var(--md-sys-color-outline)] hover:text-[var(--md-sys-color-on-surface)] ml-2"
          >
            ✕
          </button>
        </div>
      ) : null}
    </Card>
  );
}
