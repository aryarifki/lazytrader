import { Icon } from "@iconify/react";

import { Badge, Card, SectionTitle } from "@/components/research/ui";
import { fmtDate, fmtPct } from "@/components/research/format";
import { CtHoldersTable } from "./CtHoldersTable";
import { SnapshotHoldersTable } from "./SnapshotHoldersTable";
import { StagingHoldersTable } from "./StagingHoldersTable";
import type { CtHolder, StockOwnership } from "./stockTypes";

interface StockHoldersCardProps {
  data: StockOwnership;
  activeHolderTab: "staging" | "ct";
  onSelectHolderTab: (tab: "staging" | "ct") => void;
  hasStagingHolders: boolean;
  hasCtHolders: boolean;
  displayHolders: CtHolder[] | null;
}

export function StockHoldersCard({
  data,
  activeHolderTab,
  onSelectHolderTab,
  hasStagingHolders,
  hasCtHolders,
  displayHolders,
}: StockHoldersCardProps) {
  return (
    <Card className="overflow-hidden">
      <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <SectionTitle
            icon="ph:list-numbers-duotone"
            title="Daftar Pemegang Saham"
            subtitle={
              activeHolderTab === "staging" && data.staging_holders
                ? `Posisi per ${fmtDate(data.staging_holders.report_date)} · Termasuk perbandingan porsi bulan lalu`
                : hasCtHolders
                  ? "Perbandingan posisi periode · status perubahan kepemilikan (CT)"
                  : "Snapshot saat ini · delta kepemilikan per investor belum tersedia"
            }
          />
          <div className="flex items-center gap-1.5 self-start sm:self-auto flex-wrap">
            {hasStagingHolders && hasCtHolders && (
              <div className="inline-flex rounded-lg border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] p-0.5 text-[10px] font-medium">
                <button
                  type="button"
                  onClick={() => onSelectHolderTab("staging")}
                  className={`rounded-md px-2.5 py-1 transition-colors ${
                    activeHolderTab === "staging"
                      ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-sm"
                      : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                  }`}
                >
                  KSEI Staging ({fmtDate(data.staging_holders?.report_date ?? "")})
                </button>
                <button
                  type="button"
                  onClick={() => onSelectHolderTab("ct")}
                  className={`rounded-md px-2.5 py-1 transition-colors ${
                    activeHolderTab === "ct"
                      ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-sm"
                      : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                  }`}
                >
                  Baseline CT
                </button>
              </div>
            )}
            {hasStagingHolders ? (
              <Badge tone="up">KSEI Staging</Badge>
            ) : hasCtHolders ? (
              <Badge tone="neutral">CT Reference</Badge>
            ) : null}
            <Badge>
              {activeHolderTab === "staging" && data.staging_holders
                ? `${data.staging_holders.holders.length} pemegang`
                : `${displayHolders?.length ?? data.holders.length} baris`}
            </Badge>
          </div>
        </div>
      </div>

      {/* UBO Callout if available */}
      {data.staging_ubo && data.staging_ubo.length > 0 && (
        <div className="mx-4 sm:mx-5 mt-4 rounded-xl border border-[var(--md-sys-color-outline-variant)]/60 bg-[var(--md-sys-color-surface-container-low)] p-3 text-xs">
          <div className="flex items-center gap-2 mb-1.5 font-bold text-[var(--md-sys-color-on-surface)]">
            <Icon icon="ph:identification-card-duotone" width="16" className="text-[var(--md-sys-color-primary)]" />
            <span>Ultimate Beneficial Owner (UBO / Pemilik Manfaat Akhir)</span>
          </div>
          <div className="flex flex-wrap gap-2">
            {data.staging_ubo.map((ubo, i) => (
              <div key={i} className="flex items-center gap-1.5 rounded-lg bg-[var(--md-sys-color-surface-container-high)] px-2.5 py-1 text-[11px]">
                <span className="font-semibold text-[var(--md-sys-color-primary)]">{ubo.name}</span>
                <span className="font-mono font-bold text-[var(--md-sys-color-on-surface)]">({fmtPct(ubo.pct)})</span>
                {ubo.subholdings && ubo.subholdings.length > 0 && (
                  <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
                    via {ubo.subholdings.map((s) => s.name).join(", ")}
                  </span>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="overflow-x-auto">
        {activeHolderTab === "staging" && data.staging_holders ? (
          <StagingHoldersTable holders={data.staging_holders.holders} />
        ) : displayHolders ? (
          <CtHoldersTable holders={displayHolders} />
        ) : (
          <SnapshotHoldersTable holders={data.holders} />
        )}
      </div>
    </Card>
  );
}
