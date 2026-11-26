import { Icon } from "@iconify/react";

import { MiniLine } from "@/components/research/ui";
import { fmtNumber, fmtPct } from "@/components/research/format";

export function TrendMini({
  label,
  value,
  points,
  color,
  format,
  delta,
  deltaPct: _deltaPct,
  deltaUnit = "%",
  isSid = false,
  reverseColor = false,
  comparisonLabel,
  className = "",
}: {
  label: string;
  value: number;
  points: number[];
  color: string;
  format?: (value: number) => string;
  delta?: number | null;
  deltaPct?: number | null;
  deltaUnit?: string;
  isSid?: boolean;
  reverseColor?: boolean;
  comparisonLabel?: string | null;
  className?: string;
}) {
  const hasDelta = delta != null && !Number.isNaN(delta);
  const isUp = hasDelta && (isSid ? delta! > 0 : delta! > 0.005);
  const isDown = hasDelta && (isSid ? delta! < 0 : delta! < -0.005);

  return (
    <div className={`rounded-xl border border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-high)]/45 p-3 flex flex-col justify-between ${className}`}>
      <div>
        <div className="flex items-center justify-between gap-1.5">
          <span className="text-[10px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)] truncate" title={label}>
            {label}
          </span>
          <span className="text-[12px] font-mono font-bold tracking-tight" style={{ color }}>
            {(format ?? fmtPct)(value)}
          </span>
        </div>

        {hasDelta ? (
          <div className="mt-1 flex items-center justify-between gap-1 text-[10px] flex-wrap">
            <span className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-medium">
              vs {comparisonLabel ?? "Bulan Lalu"}:
            </span>
            <div className="flex items-center gap-1 font-mono font-semibold">
              {isUp ? (
                <span
                  className={`inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded text-[10px] font-mono ${
                    reverseColor
                      ? "text-rose-500 bg-rose-500/10"
                      : "text-emerald-500 bg-emerald-500/10"
                  }`}
                >
                  <Icon icon="ph:trend-up-bold" width="10" />
                  <span>
                    +{isSid ? fmtNumber(delta!) : delta!.toFixed(2)}
                    {deltaUnit}
                  </span>
                  <span className="text-[8px] uppercase tracking-wider font-sans font-bold ml-0.5">
                    Naik
                  </span>
                </span>
              ) : isDown ? (
                <span
                  className={`inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded text-[10px] font-mono ${
                    reverseColor
                      ? "text-emerald-500 bg-emerald-500/10"
                      : "text-rose-500 bg-rose-500/10"
                  }`}
                >
                  <Icon icon="ph:trend-down-bold" width="10" />
                  <span>
                    {isSid ? fmtNumber(delta!) : delta!.toFixed(2)}
                    {deltaUnit}
                  </span>
                  <span className="text-[8px] uppercase tracking-wider font-sans font-bold ml-0.5">
                    Turun
                  </span>
                </span>
              ) : (
                <span className="inline-flex items-center gap-0.5 text-[var(--md-sys-color-on-surface-variant)] bg-[var(--md-sys-color-surface-container-highest)] px-1.5 py-0.5 rounded text-[10px] font-mono">
                  <Icon icon="ph:minus-bold" width="10" />
                  <span>0.00{deltaUnit}</span>
                  <span className="text-[8px] uppercase tracking-wider font-sans font-bold ml-0.5">
                    Tetap
                  </span>
                </span>
              )}
            </div>
          </div>
        ) : null}
      </div>

      <div className="mt-2.5">
        <MiniLine points={points} color={color} width={130} height={32} />
      </div>
    </div>
  );
}
