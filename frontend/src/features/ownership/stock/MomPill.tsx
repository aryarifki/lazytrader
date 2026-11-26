import { Icon } from "@iconify/react";

export function MomPill({ delta, reverse = false }: { delta: number | null; reverse?: boolean }) {
  if (delta == null || Number.isNaN(delta)) {
    return <span className="text-[9px] text-[var(--md-sys-color-on-surface-variant)]">—</span>;
  }
  const isUp = delta > 0.005;
  const isDown = delta < -0.005;
  if (!isUp && !isDown) {
    return (
      <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] inline-flex items-center gap-0.5 justify-end">
        <Icon icon="ph:minus-bold" width="8" />
        <span>0.00%</span>
        <span className="text-[8px] uppercase">Tetap</span>
      </div>
    );
  }

  const isGreen = reverse ? isDown : isUp;
  const colorClass = isGreen ? "text-emerald-500" : "text-rose-500";
  const icon = isUp ? "ph:trend-up-bold" : "ph:trend-down-bold";
  const label = isUp ? "Naik" : "Turun";
  const sign = isUp ? "+" : "";

  return (
    <div className={`text-[9px] ${colorClass} inline-flex items-center gap-0.5 justify-end`}>
      <Icon icon={icon} width="8" />
      <span>
        {sign}
        {delta.toFixed(2)}%
      </span>
      <span className="text-[8px] uppercase">({label})</span>
    </div>
  );
}
