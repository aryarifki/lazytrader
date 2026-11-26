export function PillGroup({
  label,
  options,
  selected,
  onToggle,
}: {
  label: string;
  options: Array<{ id: string; label: string; title?: string }>;
  selected: string[];
  onToggle: (value: string) => void;
}) {
  return (
    <div>
      <div className="text-[9px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)] mb-1.5">
        {label}
      </div>
      <div className="flex flex-wrap gap-1.5">
        {options.map((option) => {
          const active = selected.includes(option.id);
          return (
            <button
              key={option.id}
              type="button"
              title={option.title}
              onClick={() => onToggle(option.id)}
              className={`rounded-lg border px-2.5 py-1.5 text-[10px] font-bold transition-all ${
                active
                  ? "border-[var(--md-sys-color-primary)] bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)]"
                  : "border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)]"
              }`}
            >
              {option.label}
            </button>
          );
        })}
      </div>
    </div>
  );
}
