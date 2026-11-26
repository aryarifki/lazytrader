import { Icon } from "@iconify/react";

export default function FeaturePlaceholder({ title, phase }: { title: string; phase: string }) {
  return (
    <main className="min-h-[70vh] flex items-center justify-center p-6">
      <div className="max-w-md w-full rounded-3xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)] p-8 text-center shadow-sm space-y-3">
        <div className="w-12 h-12 rounded-2xl mx-auto flex items-center justify-center bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)]">
          <Icon icon="ph:wrench-duotone" width="24" />
        </div>
        <h1 className="text-lg font-extrabold text-[var(--md-sys-color-on-surface)]">{title}</h1>
        <p className="text-xs text-[var(--md-sys-color-on-surface-variant)]">
          Modul ini dijadwalkan pada <span className="font-mono font-bold">{phase}</span> dengan UI/UX 100% InvestOwl MD3.
        </p>
      </div>
    </main>
  );
}
