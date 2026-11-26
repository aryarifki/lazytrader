import { useNavigate } from "react-router-dom";
import { Icon } from "@iconify/react";

import { Card, DataNote, SectionTitle } from "@/components/research/ui";
import type { InvestorSearchItem, InvestorSearchResponse } from "./investorTypes";

interface InvestorSearchPanelProps {
  query: string;
  onQueryChange: (value: string) => void;
  deferredQuery: string;
  search: InvestorSearchResponse | undefined;
  searchLoading: boolean;
  suggestions: InvestorSearchItem[];
  effectiveName: string;
}

export function InvestorSearchPanel({
  query,
  onQueryChange,
  deferredQuery,
  search,
  searchLoading,
  suggestions,
  effectiveName,
}: InvestorSearchPanelProps) {
  const navigate = useNavigate();

  return (
    <Card className="p-4 h-fit lg:sticky lg:top-20">
      <SectionTitle
        icon="ph:users-three-duotone"
        title="By Investor"
        subtitle="Snapshot CT · 5.262 investor"
      />
      <div className="relative mt-3">
        <Icon
          icon="ph:magnifying-glass-duotone"
          width="15"
          className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--md-sys-color-on-surface-variant)]"
        />
        <input
          value={query}
          onChange={(event) => onQueryChange(event.target.value)}
          placeholder="Ketik minimal 2 huruf..."
          aria-label="Cari investor"
          className="w-full rounded-xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] pl-9 pr-3 py-2.5 text-xs text-[var(--md-sys-color-on-surface)] placeholder:text-[var(--md-sys-color-on-surface-variant)] focus:outline-none focus:border-[var(--md-sys-color-primary)]"
        />
      </div>

      <div className="mt-2 flex items-center justify-between text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
        <span>{deferredQuery ? "Hasil pencarian" : "Investor teratas"}</span>
        {search ? <span>{search.total_matches} nama</span> : null}
      </div>

      <div className="mt-2 space-y-1.5 max-h-[520px] overflow-y-auto pr-1">
        {searchLoading ? (
          <div className="py-6 text-center text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)] animate-pulse">
            Mencari investor...
          </div>
        ) : suggestions.length > 0 ? (
          suggestions.map((item) => {
            const active = item.name === effectiveName;
            return (
              <button
                key={item.name}
                type="button"
                onClick={() =>
                  navigate(
                    `/ownership?tab=investor&name=${encodeURIComponent(item.name)}`,
                    { replace: true, preventScrollReset: true }
                  )
                }
                className={`w-full rounded-xl border px-3 py-2.5 text-left transition-all ${
                  active
                    ? "border-[var(--md-sys-color-primary)]/40 bg-[var(--md-sys-color-primary-container)]/60"
                    : "border-transparent hover:border-[var(--md-sys-color-outline-variant)] hover:bg-[var(--md-sys-color-surface-container-high)]/60"
                }`}
              >
                <span className="block text-[11px] font-semibold text-[var(--md-sys-color-on-surface)] line-clamp-2">
                  {item.name}
                </span>
                <span className="mt-1 flex items-center gap-2 text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                  <span>{item.ticker_count} emiten</span>
                  <span className="opacity-40">•</span>
                  <span>{item.lf_label || item.lf || "tipe tidak tersedia"}</span>
                </span>
              </button>
            );
          })
        ) : (
          <div className="py-6 text-center text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
            Investor tidak ditemukan.
          </div>
        )}
      </div>

      <div className="mt-3">
        <DataNote>
          Autocomplete menggabungkan indeks CT dan seluruh baris portofolio, sehingga nama
          yang tidak masuk indeks tetap bisa ditemukan.
        </DataNote>
      </div>
    </Card>
  );
}
