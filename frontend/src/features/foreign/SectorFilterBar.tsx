interface SectorFilterBarProps {
  selectedQuadrant: string;
  onQuadrantChange: (q: string) => void;
  quadrantCounts: Record<string, number>;
  timeframe: "1d" | "5d" | "20d";
  onTimeframeChange: (tf: "1d" | "5d" | "20d") => void;
}

export default function SectorFilterBar({
  selectedQuadrant,
  onQuadrantChange,
  quadrantCounts,
  timeframe,
  onTimeframeChange,
}: SectorFilterBarProps) {
  return (
    <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      {/* Quadrant Quick Filter Pills */}
      <div className="flex gap-1.5 overflow-x-auto pb-1 no-scrollbar text-xs font-medium">
        <button
          onClick={() => onQuadrantChange("all")}
          className={`px-3 py-1.5 rounded-full shrink-0 transition-all cursor-pointer ${
            selectedQuadrant === "all"
              ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-xs"
              : "bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface-variant)] hover:bg-[var(--md-sys-color-surface-container-high)]"
          }`}
        >
          Semua Sektor (11)
        </button>
        {["Leading", "Improving", "Weakening", "Lagging"].map((q) => (
          <button
            key={q}
            onClick={() => onQuadrantChange(q.toLowerCase())}
            className={`px-3 py-1.5 rounded-full shrink-0 transition-all cursor-pointer ${
              selectedQuadrant === q.toLowerCase()
                ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-xs"
                : "bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface-variant)] hover:bg-[var(--md-sys-color-surface-container-high)]"
            }`}
          >
            {q} ({quadrantCounts[q] || 0})
          </button>
        ))}
      </div>

      {/* Timeframe Selector */}
      <div className="flex items-center p-1 rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] self-start sm:self-auto text-xs font-semibold">
        <button
          onClick={() => onTimeframeChange("1d")}
          className={`px-3 py-1 rounded-lg transition-all cursor-pointer ${
            timeframe === "1d"
              ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-2xs"
              : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
          }`}
        >
          1 Hari
        </button>
        <button
          onClick={() => onTimeframeChange("5d")}
          className={`px-3 py-1 rounded-lg transition-all cursor-pointer ${
            timeframe === "5d"
              ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-2xs"
              : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
          }`}
        >
          5 Hari
        </button>
        <button
          onClick={() => onTimeframeChange("20d")}
          className={`px-3 py-1 rounded-lg transition-all cursor-pointer ${
            timeframe === "20d"
              ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-2xs"
              : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
          }`}
        >
          20 Hari
        </button>
      </div>
    </div>
  );
}
