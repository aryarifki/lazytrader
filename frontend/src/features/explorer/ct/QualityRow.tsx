import { Icon } from "@iconify/react";

export function QualityRow({ label, ok }: { label: string; ok: boolean }) {
  return (
    <div className="flex items-center gap-1.5 rounded-md border border-[var(--md-sys-color-outline-variant)]/50 px-2 py-1.5">
      <Icon
        icon={ok ? "ph:check-circle-duotone" : "ph:minus-circle-duotone"}
        width="13"
        className={ok ? "text-emerald-500" : "text-[var(--md-sys-color-on-surface-variant)]"}
      />
      <span className="text-[var(--md-sys-color-on-surface-variant)]">{label}</span>
    </div>
  );
}
