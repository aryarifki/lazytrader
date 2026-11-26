import { Icon } from "@iconify/react";
import type { BacktestStats } from "./types";

interface BacktestKpiCardsProps {
  stats?: BacktestStats;
  isLoading: boolean;
}

export default function BacktestKpiCards({ stats, isLoading }: BacktestKpiCardsProps) {
  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 sm:p-5 shadow-xs flex flex-col justify-between">
        <div className="flex items-center justify-between mb-2">
          <span className="text-[10px] sm:text-xs font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
            Win Rate Ratio
          </span>
          <div className="w-8 h-8 rounded-lg bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-primary)] flex items-center justify-center border border-[var(--md-sys-color-outline-variant)]">
            <Icon icon="ph:trophy-bold" width="16" />
          </div>
        </div>
        <div>
          <div className="text-2xl sm:text-3xl font-mono font-extrabold text-[var(--md-sys-color-on-surface)] tracking-tight">
            {isLoading ? "..." : `${stats?.win_rate ?? 0}%`}
          </div>
          <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-1 font-medium">
            {stats?.win_count ?? 0} WIN / {stats?.total_closed ?? 0} Selesai
          </p>
        </div>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 sm:p-5 shadow-xs flex flex-col justify-between">
        <div className="flex items-center justify-between mb-2">
          <span className="text-[10px] sm:text-xs font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
            Avg Target Gain
          </span>
          <div className="w-8 h-8 rounded-lg bg-[var(--md-sys-color-surface-container-high)] text-[var(--color-positive)] flex items-center justify-center border border-[var(--md-sys-color-outline-variant)]">
            <Icon icon="ph:trend-up-bold" width="16" />
          </div>
        </div>
        <div>
          <div className="text-2xl sm:text-3xl font-mono font-extrabold text-[var(--color-positive)] tracking-tight">
            {isLoading ? "..." : `+${stats?.avg_win_pct ?? 0}%`}
          </div>
          <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-1 font-medium">
            Profit Factor: <strong className="text-[var(--md-sys-color-on-surface)]">{stats?.profit_factor ?? 0}x</strong>
          </p>
        </div>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 sm:p-5 shadow-xs flex flex-col justify-between">
        <div className="flex items-center justify-between mb-2">
          <span className="text-[10px] sm:text-xs font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
            Avg Loss &amp; Duration
          </span>
          <div className="w-8 h-8 rounded-lg bg-[var(--md-sys-color-surface-container-high)] text-[var(--color-negative)] flex items-center justify-center border border-[var(--md-sys-color-outline-variant)]">
            <Icon icon="ph:shield-warning-bold" width="16" />
          </div>
        </div>
        <div>
          <div className="text-2xl sm:text-3xl font-mono font-extrabold text-[var(--color-negative)] tracking-tight">
            {isLoading ? "..." : `${stats?.avg_loss_pct ?? 0}%`}
          </div>
          <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-1 font-medium">
            Rata-rata Hold: <strong className="text-[var(--md-sys-color-on-surface)]">{stats?.avg_days_held ?? 0} Hari</strong>
          </p>
        </div>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 sm:p-5 shadow-xs flex flex-col justify-between">
        <div className="flex items-center justify-between mb-2">
          <span className="text-[10px] sm:text-xs font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
            Dataset Pipeline
          </span>
          <div className="w-8 h-8 rounded-lg bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-primary)] flex items-center justify-center border border-[var(--md-sys-color-outline-variant)]">
            <Icon icon="ph:database-bold" width="16" />
          </div>
        </div>
        <div>
          <div className="text-2xl sm:text-3xl font-mono font-extrabold text-[var(--md-sys-color-primary)] tracking-tight">
            {isLoading ? "..." : `${stats?.total_signals ?? 0}`}
          </div>
          <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-1 font-medium">
            {stats?.total_open ?? 0} Sinyal Sedang Aktif (OPEN)
          </p>
        </div>
      </div>
    </div>
  );
}
