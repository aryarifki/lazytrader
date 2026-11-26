import { Badge, Card, EmptyState } from "@/components/research/ui";
import { fmtDate, fmtNumber, fmtSignedPct, toneFor } from "@/components/research/format";
import { pctClass } from "./ihsgHelpers";
import type { IndexPoint, MarketOverview } from "./ihsgTypes";

export function IhsgHeaderCard({ data }: { data: MarketOverview }) {
  return (
    <Card className="p-4 sm:p-5 overflow-hidden relative">
      <div className="absolute -right-16 -top-20 h-56 w-56 rounded-full bg-[var(--md-sys-color-primary)]/5 blur-3xl pointer-events-none" />
      <div className="relative flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 flex-wrap">
            <Badge tone="primary">IHSG</Badge>
            <Badge tone="neutral">COMPOSITE</Badge>
            <span className="text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
              Per {fmtDate(data.index.date)}
            </span>
          </div>
          <div className="mt-2 flex items-end gap-3 flex-wrap">
            <span className="text-3xl sm:text-4xl font-extrabold tracking-tight">
              {fmtNumber(data.index.close, 3)}
            </span>
            <span
              className={`text-sm font-mono font-bold mb-1 ${
                toneFor(data.index.change_1d_pct) === "up"
                  ? "text-emerald-600 dark:text-emerald-400"
                  : toneFor(data.index.change_1d_pct) === "down"
                    ? "text-rose-600 dark:text-rose-400"
                    : "text-[var(--md-sys-color-on-surface-variant)]"
              }`}
            >
              {fmtSignedPct(data.index.change_1d_pct)}
            </span>
          </div>
          <p className="mt-1 text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
            {fmtDate(data.history[0]?.date)} — {fmtDate(data.history[data.history.length - 1]?.date)}
            {" · "}
            {data.history.length} titik data
          </p>
        </div>
        <div className="grid grid-cols-3 gap-2 w-full sm:w-auto">
          <MiniStat label="1W" value={data.index.return_1w_pct} />
          <MiniStat label="1M" value={data.index.return_1m_pct} />
          <MiniStat label="YTD" value={data.index.return_ytd_pct} />
        </div>
      </div>

      <div className="relative mt-5">
        <IndexChart points={data.history} />
      </div>
    </Card>
  );
}

function MiniStat({ label, value }: { label: string; value: number | null }) {
  return (
    <div className="rounded-lg border border-[var(--md-sys-color-outline-variant)]/60 bg-[var(--md-sys-color-surface-container-high)]/50 px-2.5 py-2 text-center">
      <div className="text-[8px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
        {label}
      </div>
      <div className={`text-xs font-mono font-bold mt-0.5 ${pctClass(value)}`}>
        {fmtSignedPct(value)}
      </div>
    </div>
  );
}

function IndexChart({ points }: { points: IndexPoint[] }) {
  if (points.length < 2) return <EmptyState title="Riwayat IHSG belum cukup" />;

  const width = 900;
  const height = 220;
  const closes = points.map((point) => point.close);
  const min = Math.min(...closes);
  const max = Math.max(...closes);
  const range = max - min || 1;
  const coords = points.map((point, index) => {
    const x = (index / (points.length - 1)) * width;
    const y = height - 12 - ((point.close - min) / range) * (height - 24);
    return { x, y };
  });
  const line = coords.map((point) => `${point.x.toFixed(2)},${point.y.toFixed(2)}`).join(" ");
  const area = `0,${height} ${line} ${width},${height}`;

  return (
    <svg viewBox={`0 0 ${width} ${height}`} className="w-full h-[190px] sm:h-[230px]" preserveAspectRatio="none" role="img" aria-label="Grafik IHSG">
      <defs>
        <linearGradient id="ihsg-area" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor="var(--md-sys-color-primary)" stopOpacity="0.22" />
          <stop offset="100%" stopColor="var(--md-sys-color-primary)" stopOpacity="0" />
        </linearGradient>
      </defs>
      {[0.25, 0.5, 0.75].map((ratio) => (
        <line
          key={ratio}
          x1="0"
          x2={width}
          y1={height * ratio}
          y2={height * ratio}
          stroke="var(--md-sys-color-outline-variant)"
          strokeWidth="0.5"
          strokeDasharray="3 5"
          opacity="0.55"
        />
      ))}
      <polygon points={area} fill="url(#ihsg-area)" />
      <polyline
        points={line}
        fill="none"
        stroke="var(--md-sys-color-primary)"
        strokeWidth="2"
        strokeLinejoin="round"
        strokeLinecap="round"
        vectorEffect="non-scaling-stroke"
      />
    </svg>
  );
}
