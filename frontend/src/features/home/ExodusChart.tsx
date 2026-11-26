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
import type { TimelinePoint } from "./types";

interface ExodusChartProps {
  timeline: TimelinePoint[];
  isLoading: boolean;
  isNetOutflow: boolean;
  gradientOffset: number;
}

function formatShortDate(dateStr: string): string {
  if (!dateStr) return "";
  try {
    const parts = dateStr.split("-");
    if (parts.length === 3) {
      const months = ["Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"];
      return `${parseInt(parts[2], 10)} ${months[parseInt(parts[1], 10) - 1]}`;
    }
    return dateStr;
  } catch {
    return dateStr;
  }
}

export default function ExodusChart({
  timeline,
  isLoading,
  isNetOutflow,
  gradientOffset,
}: ExodusChartProps) {
  return (
    <>
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 text-xs text-[var(--md-sys-color-on-surface-variant)]">
        <div className="flex items-center gap-2">
          <Icon icon="ph:chart-line-up-duotone" width="16" className="text-[var(--md-sys-color-primary)]" />
          <span className="font-bold text-[var(--md-sys-color-on-surface)]">
            Dinamika Kapital Asing vs IHSG (45 Hari Terakhir)
          </span>
        </div>
        <div className="flex items-center gap-4 text-[11px] font-mono">
          <span className="inline-flex items-center gap-1.5">
            <span className="w-2.5 h-2.5 rounded-full bg-[#38bdf8]"></span>
            <span>IHSG Close (Kiri)</span>
          </span>
          <span className="inline-flex items-center gap-1.5">
            <span className={`w-2.5 h-2.5 rounded-full ${isNetOutflow ? "bg-rose-500" : "bg-emerald-500"}`}></span>
            <span>Kumulatif Flow (Kanan)</span>
          </span>
        </div>
      </div>

      <div className="h-[280px] sm:h-[320px] w-full bg-[var(--md-sys-color-surface-container-low)] rounded-2xl p-3 border border-[var(--md-sys-color-outline-variant)]/50">
        {isLoading ? (
          <div className="h-full flex items-center justify-center text-xs font-mono text-[var(--md-sys-color-on-surface-variant)] animate-pulse">
            Memuat data radar historis...
          </div>
        ) : timeline.length === 0 ? (
          <div className="h-full flex items-center justify-center text-xs font-mono text-[var(--md-sys-color-on-surface-variant)]">
            Data historis tidak tersedia.
          </div>
        ) : (
          <ResponsiveContainer width="100%" height="100%">
            <ComposedChart data={timeline} margin={{ top: 10, right: 10, left: -15, bottom: 0 }}>
              <defs>
                <linearGradient id="flowGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset={gradientOffset} stopColor="#10b981" stopOpacity={0.35} />
                  <stop offset={gradientOffset} stopColor="#f43f5e" stopOpacity={0.35} />
                  <stop offset="100%" stopColor={isNetOutflow ? "#f43f5e" : "#10b981"} stopOpacity={0.05} />
                </linearGradient>
              </defs>

              <CartesianGrid strokeDasharray="3 3" stroke="var(--md-sys-color-outline-variant)" opacity={0.35} />

              <XAxis
                dataKey="date"
                tickFormatter={formatShortDate}
                stroke="var(--md-sys-color-on-surface-variant)"
                fontSize={10}
                tickLine={false}
                axisLine={{ stroke: "var(--md-sys-color-outline-variant)", opacity: 0.5 }}
              />

              <YAxis
                yAxisId="left"
                domain={["auto", "auto"]}
                stroke="var(--md-sys-color-on-surface-variant)"
                fontSize={10}
                tickFormatter={(v) => Math.round(Number(v)).toLocaleString("id-ID")}
                tickLine={false}
                axisLine={{ stroke: "var(--md-sys-color-outline-variant)", opacity: 0.5 }}
              />

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
                content={({ active, payload }) => {
                  if (active && payload && payload.length) {
                    const p = payload[0].payload as TimelinePoint;
                    return (
                      <div className="bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)] rounded-xl p-3 shadow-lg text-xs font-mono space-y-1.5 z-50">
                        <div className="text-[var(--md-sys-color-on-surface-variant)] font-bold border-b border-[var(--md-sys-color-outline-variant)] pb-1">
                          Tanggal: <span className="text-[var(--md-sys-color-on-surface)]">{p.date}</span>
                        </div>
                        <div className="flex justify-between gap-4">
                          <span className="text-[#38bdf8]">IHSG Close:</span>
                          <strong className="text-[var(--md-sys-color-on-surface)]">
                            {p.ihsg_close.toLocaleString("id-ID")}
                          </strong>
                        </div>
                        <div className="flex justify-between gap-4">
                          <span>Foreign Net (Hari Ini):</span>
                          <strong style={{ color: signedColor(p.daily_flow) }}>
                            {p.daily_flow > 0 ? "+" : ""}
                            {fmtRadarFlow(p.daily_flow)}
                          </strong>
                        </div>
                        <div className="flex justify-between gap-4">
                          <span>Kumulatif Flow:</span>
                          <strong style={{ color: signedColor(p.cumulative_flow) }}>
                            {p.cumulative_flow > 0 ? "+" : ""}
                            {fmtRadarFlow(p.cumulative_flow)}
                          </strong>
                        </div>
                      </div>
                    );
                  }
                  return null;
                }}
              />

              <ReferenceLine yAxisId="right" y={0} stroke="var(--md-sys-color-outline)" strokeDasharray="3 3" opacity={0.6} />

              <Area
                yAxisId="right"
                type="monotone"
                dataKey="cumulative_flow"
                baseValue={0}
                fill="url(#flowGradient)"
                stroke={isNetOutflow ? "#f43f5e" : "#10b981"}
                strokeWidth={2}
                isAnimationActive={false}
              />

              <Line
                yAxisId="left"
                type="monotone"
                dataKey="ihsg_close"
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
    </>
  );
}
