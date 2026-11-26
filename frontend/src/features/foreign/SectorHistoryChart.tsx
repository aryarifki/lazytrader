import { Icon } from "@iconify/react";
import {
  ResponsiveContainer,
  ComposedChart,
  Line,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  ReferenceLine,
} from "recharts";
import { fmtRadarFlow, signedColor } from "@/lib/format";
import { formatShortDate } from "./sectorHelpers";
import type { SectorHistoryPoint } from "./sectorTypes";

interface SectorHistoryChartProps {
  sectorCode: string;
  historyTimeline: SectorHistoryPoint[];
  isHistoryLoading: boolean;
  isNetOutflow: boolean;
  gradientOffset: number;
}

export default function SectorHistoryChart({
  sectorCode,
  historyTimeline,
  isHistoryLoading,
  isNetOutflow,
  gradientOffset,
}: SectorHistoryChartProps) {
  return (
    <div className="space-y-2">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between text-xs gap-2">
        <span className="font-semibold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
          <Icon icon="ph:chart-line-up-bold" width="16" className="text-[var(--md-sys-color-primary)]" />
          Dinamika Indeks Sektor vs Aliran Modal Asing Kumulatif (30 Hari)
        </span>
        <div className="flex items-center gap-3 text-[11px] font-mono">
          <span className="inline-flex items-center gap-1.5">
            <span className="w-2.5 h-2.5 rounded-full bg-[#38bdf8]"></span>
            <span>Close Indeks (Kiri)</span>
          </span>
          <span className="inline-flex items-center gap-1.5">
            <span
              className={`w-2.5 h-2.5 rounded-full ${isNetOutflow ? "bg-rose-500" : "bg-emerald-500"}`}
            ></span>
            <span>Kumulatif Net Foreign (Kanan)</span>
          </span>
        </div>
      </div>

      <div className="h-[260px] sm:h-[300px] w-full bg-[var(--md-sys-color-surface-container-low)] rounded-2xl p-3 border border-[var(--md-sys-color-outline-variant)]/60">
        {isHistoryLoading ? (
          <div className="h-full flex items-center justify-center text-xs font-mono text-[var(--md-sys-color-on-surface-variant)] animate-pulse gap-2">
            <Icon icon="ph:spinner-gap-bold" className="animate-spin text-[var(--md-sys-color-primary)]" width="16" />
            <span>Memuat dinamika arus asing sektor...</span>
          </div>
        ) : historyTimeline.length === 0 ? (
          <div className="h-full flex items-center justify-center text-xs font-mono text-[var(--md-sys-color-on-surface-variant)]">
            Data historis sektor tidak tersedia.
          </div>
        ) : (
          <ResponsiveContainer width="100%" height="100%">
            <ComposedChart data={historyTimeline} margin={{ top: 10, right: 10, left: -10, bottom: 0 }}>
              <defs>
                <linearGradient id={`flowGrad-${sectorCode}`} x1="0" y1="0" x2="0" y2="1">
                  <stop offset={gradientOffset} stopColor="#10b981" stopOpacity={0.35} />
                  <stop offset={gradientOffset} stopColor="#f43f5e" stopOpacity={0.35} />
                  <stop offset="100%" stopColor={isNetOutflow ? "#f43f5e" : "#10b981"} stopOpacity={0.05} />
                </linearGradient>
              </defs>

              <CartesianGrid strokeDasharray="3 3" stroke="var(--md-sys-color-outline-variant)" opacity={0.3} />

              <XAxis
                dataKey="date"
                tickFormatter={formatShortDate}
                stroke="var(--md-sys-color-on-surface-variant)"
                fontSize={10}
                tickLine={false}
                axisLine={{ stroke: "var(--md-sys-color-outline-variant)", opacity: 0.5 }}
              />

              {/* Left Y-Axis: Sector Close Price */}
              <YAxis
                yAxisId="left"
                domain={["auto", "auto"]}
                stroke="var(--md-sys-color-on-surface-variant)"
                fontSize={10}
                tickFormatter={(v) => Math.round(Number(v)).toLocaleString("id-ID")}
                tickLine={false}
                axisLine={{ stroke: "var(--md-sys-color-outline-variant)", opacity: 0.5 }}
              />

              {/* Right Y-Axis: Cumulative Sector Foreign Net Flow */}
              <YAxis
                yAxisId="right"
                orientation="right"
                domain={["auto", "auto"]}
                stroke="var(--md-sys-color-on-surface-variant)"
                fontSize={10}
                tickFormatter={(v) => fmtRadarFlow(Number(v))}
                tickLine={false}
                axisLine={{ stroke: "var(--md-sys-color-outline-variant)", opacity: 0.5 }}
              />

              <Tooltip
                content={({
                  active,
                  payload,
                }: {
                  active?: boolean;
                  payload?: ReadonlyArray<{ payload?: SectorHistoryPoint }>;
                }) => {
                  if (active && payload && payload.length && payload[0].payload) {
                    const p = payload[0].payload;
                    return (
                      <div className="bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)] rounded-xl p-3 shadow-lg text-xs font-mono space-y-1.5 z-50">
                        <div className="text-[var(--md-sys-color-on-surface-variant)] font-bold border-b border-[var(--md-sys-color-outline-variant)] pb-1">
                          Tanggal: <span className="text-[var(--md-sys-color-on-surface)]">{p.date}</span>
                        </div>
                        <div className="flex justify-between gap-4">
                          <span className="text-[#38bdf8]">Close Indeks:</span>
                          <strong className="text-[var(--md-sys-color-on-surface)]">
                            {p.close.toLocaleString("id-ID")}
                          </strong>
                        </div>
                        <div className="flex justify-between gap-4">
                          <span>Foreign Net (Hari Ini):</span>
                          <strong style={{ color: signedColor(p.foreign_net) }}>
                            {p.foreign_net > 0 ? "+" : ""}
                            {fmtRadarFlow(p.foreign_net)}
                          </strong>
                        </div>
                        <div className="flex justify-between gap-4">
                          <span>Kumulatif Flow:</span>
                          <strong style={{ color: signedColor(p.cumulative_foreign_flow) }}>
                            {p.cumulative_foreign_flow > 0 ? "+" : ""}
                            {fmtRadarFlow(p.cumulative_foreign_flow)}
                          </strong>
                        </div>
                      </div>
                    );
                  }
                  return null;
                }}
              />

              <ReferenceLine
                yAxisId="right"
                y={0}
                stroke="var(--md-sys-color-outline)"
                strokeDasharray="3 3"
                opacity={0.5}
              />

              {/* Cumulative Area */}
              <Area
                yAxisId="right"
                type="monotone"
                dataKey="cumulative_foreign_flow"
                baseValue={0}
                fill={`url(#flowGrad-${sectorCode})`}
                stroke={isNetOutflow ? "#f43f5e" : "#10b981"}
                strokeWidth={2}
                isAnimationActive={false}
              />

              {/* Sector Index Close Line */}
              <Line
                yAxisId="left"
                type="monotone"
                dataKey="close"
                stroke="#38bdf8"
                strokeWidth={2.5}
                dot={false}
                activeDot={{ r: 5, fill: "#38bdf8" }}
                isAnimationActive={false}
              />
            </ComposedChart>
          </ResponsiveContainer>
        )}
      </div>
    </div>
  );
}
