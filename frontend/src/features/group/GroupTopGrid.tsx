import { Icon } from "@iconify/react";
import { fmtSignedPct } from "@/components/research/format";
import { pctClass } from "./groupHelpers";
import type { GroupSummary } from "./groupTypes";

interface GroupTopGridProps {
  topGroups: GroupSummary[];
  effectiveGroup: string;
  onSelectGroup: (groupId: string) => void;
}

export function GroupTopGrid({
  topGroups,
  effectiveGroup,
  onSelectGroup,
}: GroupTopGridProps) {
  return (
    <div>
      <div className="flex items-center justify-between gap-2 mb-2">
        <h2 className="text-xs font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)] flex items-center gap-1.5">
          <Icon icon="ph:trophy-duotone" width="15" className="text-amber-500" />
          Top 10 grup · rata-rata MTD
        </h2>
        <span className="text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
          Klik untuk memilih
        </span>
      </div>
      <div className="grid grid-cols-2 md:grid-cols-5 gap-2">
        {topGroups.map((group, index) => {
          const active = group.id === effectiveGroup;
          return (
            <button
              key={group.id}
              type="button"
              onClick={() => onSelectGroup(group.id)}
              className={`rounded-xl border p-3 text-left transition-all ${
                active
                  ? "border-[var(--md-sys-color-primary)]/50 bg-[var(--md-sys-color-primary-container)]/55 shadow-xs"
                  : "border-[var(--md-sys-color-outline-variant)]/60 bg-[var(--md-sys-color-surface-container)] hover:bg-[var(--md-sys-color-surface-container-high)]/60"
              }`}
            >
              <div className="flex items-center justify-between gap-2">
                <div className="flex items-center gap-1.5 min-w-0">
                  <span className="text-[9px] font-mono font-bold text-[var(--md-sys-color-on-surface-variant)]">
                    #{index + 1}
                  </span>
                  {group.logo_url ? (
                    <div className="w-4 h-4 rounded-full bg-white/90 p-0.5 border border-[var(--md-sys-color-outline-variant)]/40 flex items-center justify-center shrink-0 overflow-hidden">
                      <img
                        src={group.logo_url}
                        alt=""
                        className="max-h-full max-w-full object-contain"
                        onError={(e) => {
                          (e.target as HTMLElement).style.display = "none";
                        }}
                      />
                    </div>
                  ) : null}
                </div>
                <span className={`text-[11px] font-mono font-bold ${pctClass(group.avg_mtd_pct)}`}>
                  {fmtSignedPct(group.avg_mtd_pct)}
                </span>
              </div>
              <div className="mt-1 text-[11px] font-bold leading-snug line-clamp-2 min-h-[28px]">
                {group.name}
              </div>
              <div className="mt-1.5 text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                {group.ticker_count} saham · {group.covered_tickers} tercakup
              </div>
              {group.top_stocks[0] ? (
                <div className="mt-1 text-[9px] text-[var(--md-sys-color-on-surface-variant)] truncate">
                  Top: <span className="font-mono font-bold text-[var(--md-sys-color-primary)]">{group.top_stocks[0].ticker}</span>{" "}
                  {fmtSignedPct(group.top_stocks[0].mtd_pct)}
                </div>
              ) : null}
            </button>
          );
        })}
      </div>
    </div>
  );
}
