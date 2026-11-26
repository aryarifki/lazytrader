import React from "react";
import { Icon } from "@iconify/react";

export function cx(...parts: Array<string | false | null | undefined>) {
  return parts.filter(Boolean).join(" ");
}

export function Card({
  children,
  className,
  as: Tag = "section",
}: {
  children: React.ReactNode;
  className?: string;
  as?: "section" | "div" | "article";
}) {
  return (
    <Tag
      className={cx(
        "bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl shadow-xs",
        className
      )}
    >
      {children}
    </Tag>
  );
}

export function SectionTitle({
  icon,
  title,
  subtitle,
  action,
}: {
  icon: string;
  title: string;
  subtitle?: string;
  action?: React.ReactNode;
}) {
  return (
    <div className="flex items-start justify-between gap-3">
      <div className="flex items-center gap-2.5 min-w-0">
        <div className="w-8 h-8 rounded-lg bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0">
          <Icon icon={icon} width="17" />
        </div>
        <div className="min-w-0">
          <h2 className="text-sm font-bold text-[var(--md-sys-color-on-surface)] tracking-tight truncate">
            {title}
          </h2>
          {subtitle ? (
            <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-0.5">
              {subtitle}
            </p>
          ) : null}
        </div>
      </div>
      {action}
    </div>
  );
}

export function StatCard({
  label,
  value,
  sub,
  icon,
  tone = "neutral",
}: {
  label: string;
  value: React.ReactNode;
  sub?: React.ReactNode;
  icon?: string;
  tone?: "neutral" | "up" | "down" | "primary";
}) {
  const toneClass =
    tone === "up"
      ? "text-emerald-600 dark:text-emerald-400"
      : tone === "down"
        ? "text-rose-600 dark:text-rose-400"
        : tone === "primary"
          ? "text-[var(--md-sys-color-primary)]"
          : "text-[var(--md-sys-color-on-surface)]";

  return (
    <div className="rounded-xl bg-[var(--md-sys-color-surface-container-high)]/60 border border-[var(--md-sys-color-outline-variant)]/50 p-3.5 min-w-0">
      <div className="flex items-center justify-between gap-2">
        <span className="text-[10px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
          {label}
        </span>
        {icon ? (
          <Icon
            icon={icon}
            width="15"
            className="text-[var(--md-sys-color-on-surface-variant)] shrink-0"
          />
        ) : null}
      </div>
      <div className={cx("text-lg sm:text-xl font-extrabold tracking-tight mt-1.5", toneClass)}>
        {value}
      </div>
      {sub ? (
        <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] mt-1 leading-relaxed">
          {sub}
        </div>
      ) : null}
    </div>
  );
}

export function Badge({
  children,
  tone = "neutral",
}: {
  children: React.ReactNode;
  tone?: "neutral" | "primary" | "up" | "down" | "warning";
}) {
  const cls =
    tone === "primary"
      ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border-[var(--md-sys-color-primary)]/20"
      : tone === "up"
        ? "bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border-emerald-500/20"
        : tone === "down"
          ? "bg-rose-500/10 text-rose-700 dark:text-rose-400 border-rose-500/20"
          : tone === "warning"
            ? "bg-amber-500/10 text-amber-700 dark:text-amber-400 border-amber-500/20"
            : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]";
  return (
    <span className={cx("inline-flex items-center gap-1 rounded-md border px-1.5 py-0.5 text-[9px] font-mono font-bold uppercase tracking-wide", cls)}>
      {children}
    </span>
  );
}

export function DataNote({
  title = "Catatan data",
  children,
  tone = "info",
}: {
  title?: string;
  children: React.ReactNode;
  tone?: "info" | "warning";
}) {
  const cls =
    tone === "warning"
      ? "border-amber-500/30 bg-amber-500/5 text-amber-900 dark:text-amber-200"
      : "border-sky-500/25 bg-sky-500/5 text-sky-900 dark:text-sky-200";
  return (
    <div className={cx("rounded-xl border p-3.5 text-[11px] leading-relaxed", cls)}>
      <div className="flex items-center gap-1.5 font-bold mb-1">
        <Icon icon={tone === "warning" ? "ph:warning-duotone" : "ph:info-duotone"} width="15" />
        <span>{title}</span>
      </div>
      {children}
    </div>
  );
}

export function EmptyState({
  icon = "ph:tray-duotone",
  title,
  description,
}: {
  icon?: string;
  title: string;
  description?: string;
}) {
  return (
    <div className="rounded-2xl border border-dashed border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-low)]/40 p-8 text-center">
      <div className="w-11 h-11 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] flex items-center justify-center mx-auto mb-3">
        <Icon icon={icon} width="22" />
      </div>
      <h3 className="text-sm font-bold text-[var(--md-sys-color-on-surface)]">{title}</h3>
      {description ? (
        <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-1.5 max-w-md mx-auto leading-relaxed">
          {description}
        </p>
      ) : null}
    </div>
  );
}

export function LoadingBlock({ label = "Memuat data..." }: { label?: string }) {
  return (
    <div className="rounded-2xl border border-[var(--md-sys-color-outline-variant)]/60 bg-[var(--md-sys-color-surface-container)] p-8 text-center">
      <div className="w-8 h-8 rounded-full border-2 border-[var(--md-sys-color-primary)]/25 border-t-[var(--md-sys-color-primary)] animate-spin mx-auto mb-3" />
      <p className="text-[11px] font-mono text-[var(--md-sys-color-on-surface-variant)]">{label}</p>
    </div>
  );
}

export function ErrorBlock({
  message = "Gagal memuat data.",
}: {
  message?: string;
}) {
  return (
    <div className="rounded-2xl border border-rose-500/25 bg-rose-500/5 p-5 text-center">
      <Icon icon="ph:warning-circle-duotone" width="22" className="text-rose-500 mx-auto mb-2" />
      <p className="text-[11px] font-mono text-rose-700 dark:text-rose-300">{message}</p>
    </div>
  );
}

export function Donut({
  items,
  centerLabel,
  centerValue,
  size = 168,
}: {
  items: Array<{ label: string; pct: number; color: string }>;
  centerLabel: string;
  centerValue: string;
  size?: number;
}) {
  const total = items.reduce((sum, item) => sum + Math.max(0, item.pct), 0);
  const radius = 42;
  const circumference = 2 * Math.PI * radius;
  let offset = 0;

  return (
    <div className="flex flex-col sm:flex-row items-center gap-5">
      <div className="relative shrink-0" style={{ width: size, height: size }}>
        <svg viewBox="0 0 120 120" className="w-full h-full -rotate-90" role="img" aria-label={centerLabel}>
          <circle
            cx="60"
            cy="60"
            r={radius}
            fill="none"
            stroke="var(--md-sys-color-surface-container-highest)"
            strokeWidth="13"
          />
          {total > 0
            ? items.map((item) => {
                const pct = Math.max(0, item.pct);
                const dash = (pct / total) * circumference;
                const circle = (
                  <circle
                    key={item.label}
                    cx="60"
                    cy="60"
                    r={radius}
                    fill="none"
                    stroke={item.color}
                    strokeWidth="13"
                    strokeDasharray={`${dash} ${circumference - dash}`}
                    strokeDashoffset={-offset}
                    strokeLinecap="butt"
                  />
                );
                offset += dash;
                return circle;
              })
            : null}
        </svg>
        <div className="absolute inset-0 flex flex-col items-center justify-center text-center px-4">
          <span className="text-[9px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
            {centerLabel}
          </span>
          <span className="text-xl font-extrabold tracking-tight text-[var(--md-sys-color-on-surface)] mt-0.5">
            {centerValue}
          </span>
        </div>
      </div>

      <div className="w-full space-y-2">
        {items.map((item) => (
          <div key={item.label} className="flex items-center justify-between gap-3">
            <div className="flex items-center gap-2 min-w-0">
              <span className="w-2.5 h-2.5 rounded-sm shrink-0" style={{ backgroundColor: item.color }} />
              <span className="text-[11px] text-[var(--md-sys-color-on-surface)] truncate">{item.label}</span>
            </div>
            <span className="text-[11px] font-mono font-bold text-[var(--md-sys-color-on-surface)] shrink-0">
              {item.pct.toFixed(2)}%
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}

export function MiniLine({
  points,
  color = "var(--md-sys-color-primary)",
  width = 150,
  height = 42,
}: {
  points: number[];
  color?: string;
  width?: number;
  height?: number;
}) {
  if (points.length < 2) {
    return <div className="h-[42px] rounded-md bg-[var(--md-sys-color-surface-container-high)]/50" />;
  }
  const min = Math.min(...points);
  const max = Math.max(...points);
  const range = max - min || 1;
  const coords = points
    .map((point, index) => {
      const x = (index / (points.length - 1)) * width;
      const y = height - ((point - min) / range) * (height - 4) - 2;
      return `${x.toFixed(2)},${y.toFixed(2)}`;
    })
    .join(" ");

  return (
    <svg viewBox={`0 0 ${width} ${height}`} className="w-full h-[42px]" preserveAspectRatio="none" aria-hidden="true">
      <polyline
        points={coords}
        fill="none"
        stroke={color}
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
        vectorEffect="non-scaling-stroke"
      />
    </svg>
  );
}

export function MiniBars({
  items,
}: {
  items: Array<{ label: string; value: number; color: string }>;
}) {
  const max = Math.max(...items.map((item) => item.value), 1);
  return (
    <div className="space-y-2.5">
      {items.map((item) => (
        <div key={item.label}>
          <div className="flex items-center justify-between text-[10px] mb-1">
            <span className="text-[var(--md-sys-color-on-surface-variant)]">{item.label}</span>
            <span className="font-mono font-bold text-[var(--md-sys-color-on-surface)]">
              {item.value.toLocaleString("id-ID")}
            </span>
          </div>
          <div className="h-1.5 rounded-full bg-[var(--md-sys-color-surface-container-highest)] overflow-hidden">
            <div
              className="h-full rounded-full transition-all"
              style={{ width: `${(item.value / max) * 100}%`, backgroundColor: item.color }}
            />
          </div>
        </div>
      ))}
    </div>
  );
}
