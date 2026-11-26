import { Icon } from "@iconify/react";

const PRESET_TAGS = [
  { id: "ALL", label: "Semua Strategi", icon: "ph:squares-four-bold" },
  { id: "GOLDEN_CONFLUENCE", label: "Golden Confluence", icon: "ph:star-bold" },
  { id: "STEALTH_ACCUMULATION", label: "Stealth Accumulation", icon: "ph:shield-check-bold" },
  { id: "GORENGAN_MOMENTUM", label: "Gorengan Momentum", icon: "ph:fire-bold" },
  { id: "DISTRIBUTION_TRAP", label: "Distribution Trap", icon: "ph:warning-bold" },
  { id: "LOW_FLOAT", label: "Low Float (<25%)", icon: "ph:feather-bold" },
];

interface ConfluenceFiltersProps {
  selectedTag: string;
  onTagChange: (tag: string) => void;
  universe: string;
  onUniverseChange: (val: string) => void;
  windowDays: number;
  onWindowDaysChange: (val: number) => void;
  selectedSignal: string;
  onSignalChange: (val: string) => void;
  minLiq: number;
  onMinLiqChange: (val: number) => void;
  minConviction: number;
  onMinConvictionChange: (val: number) => void;
  search: string;
  onSearchChange: (val: string) => void;
  filteredCount: number;
  totalCount: number;
}

export default function ConfluenceFilters({
  selectedTag,
  onTagChange,
  universe,
  onUniverseChange,
  windowDays,
  onWindowDaysChange,
  selectedSignal,
  onSignalChange,
  minLiq,
  onMinLiqChange,
  minConviction,
  onMinConvictionChange,
  search,
  onSearchChange,
  filteredCount,
  totalCount,
}: ConfluenceFiltersProps) {
  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 sm:p-5 shadow-xs space-y-4">
      <div>
        <div className="text-[11px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider mb-2 flex items-center justify-between">
          <span className="flex items-center gap-1.5">
            <Icon icon="ph:magic-wand-duotone" width="14" className="text-[var(--md-sys-color-primary)]" />
            Preset Strategi Kuantitatif
          </span>
          <span className="text-[10px] font-mono font-medium text-[var(--md-sys-color-on-surface-variant)]">
            {filteredCount} dari {totalCount} Emiten
          </span>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          {PRESET_TAGS.map((t) => {
            const active = selectedTag === t.id;
            return (
              <button
                key={t.id}
                type="button"
                onClick={() => onTagChange(t.id)}
                className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-semibold border transition-all cursor-pointer ${
                  active
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] border-[var(--md-sys-color-primary)] shadow-2xs font-bold"
                    : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                }`}
              >
                <Icon icon={t.icon} width="14" />
                <span>{t.label}</span>
              </button>
            );
          })}
        </div>
      </div>

      <div className="pt-3 border-t border-[var(--md-sys-color-outline-variant)] flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-1.5">
          <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] font-medium">Universe:</span>
          <select
            value={universe}
            onChange={(e) => onUniverseChange(e.target.value)}
            className="py-1.5 px-2.5 text-xs rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface)] focus:outline-none cursor-pointer font-semibold"
          >
            <option value="lq45">LQ45 (45 Liquid)</option>
            <option value="kompas100">KOMPAS100 (100 Top)</option>
            <option value="idx30">IDX30 (30 Bluechip)</option>
            <option value="watchlist">Watchlist Saya</option>
          </select>
        </div>

        <div className="flex items-center gap-1.5">
          <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] font-medium">Window:</span>
          <div className="flex items-center gap-1 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] p-0.5 rounded-xl text-xs font-semibold">
            {[5, 10, 20, 30, 60].map((w) => (
              <button
                key={w}
                type="button"
                onClick={() => onWindowDaysChange(w)}
                className={`px-2 py-0.5 rounded-lg text-xs transition-all cursor-pointer ${
                  windowDays === w
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-2xs font-bold"
                    : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                }`}
              >
                {w}D
              </button>
            ))}
          </div>
        </div>

        <div className="flex items-center gap-1.5">
          <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] font-medium">Sinyal:</span>
          <select
            value={selectedSignal}
            onChange={(e) => onSignalChange(e.target.value)}
            className="py-1.5 px-2.5 text-xs rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface)] focus:outline-none cursor-pointer"
          >
            <option value="ALL">Semua Sinyal</option>
            <option value="accumulation">Akumulasi (Big / Normal)</option>
            <option value="big accumulation">Big Accumulation Saja</option>
            <option value="neutral">Netral</option>
            <option value="distribution">Distribusi</option>
          </select>
        </div>

        <div className="flex items-center gap-1.5">
          <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] font-medium">Likuiditas:</span>
          <select
            value={minLiq}
            onChange={(e) => onMinLiqChange(Number(e.target.value))}
            className="py-1.5 px-2.5 text-xs rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface)] focus:outline-none cursor-pointer font-semibold"
          >
            <option value={0}>Semua Likuiditas</option>
            <option value={1000000000}>&gt; Rp 1 Miliar / hari</option>
            <option value={5000000000}>&gt; Rp 5 Miliar / hari</option>
            <option value={10000000000}>&gt; Rp 10 Miliar / hari</option>
          </select>
        </div>

        <div className="flex items-center gap-2 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] px-3 py-1 rounded-xl">
          <span className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] font-medium whitespace-nowrap">
            Min Conviction: <strong className="text-[var(--md-sys-color-primary)] font-mono">{minConviction}</strong>
          </span>
          <input
            type="range"
            min={0}
            max={90}
            step={5}
            value={minConviction}
            onChange={(e) => onMinConvictionChange(Number(e.target.value))}
            className="w-20 accent-[var(--md-sys-color-primary)] cursor-pointer"
          />
        </div>

        <div className="relative ml-auto">
          <Icon
            icon="ph:magnifying-glass-duotone"
            className="absolute left-2.5 top-1/2 -translate-y-1/2 text-[var(--md-sys-color-on-surface-variant)]"
            width="14"
          />
          <input
            type="text"
            placeholder="Cari ticker / emiten..."
            value={search}
            onChange={(e) => onSearchChange(e.target.value)}
            className="pl-8 pr-3 py-1.5 text-xs rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] focus:outline-none focus:border-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-surface)] w-40 sm:w-48"
          />
        </div>
      </div>
    </div>
  );
}
