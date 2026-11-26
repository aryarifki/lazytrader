import { Icon } from "@iconify/react";
import {
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Bar,
  ComposedChart,
  ReferenceLine,
} from "recharts";
import { fmtRp, signedColor, getProfileTextColor } from "@/lib/format";
import {
  TopBrokerSummaryCard,
  PricePerformanceCard,
  ProfileBrokerDetailCard,
} from "./OverviewTables";
import type { SmartDailyRow, TickerDetailData } from "./types";

interface SmartBarProps {
  x?: number;
  y?: number;
  width?: number;
  height?: number;
  payload?: SmartDailyRow;
}

export default function OverviewTab({ data, isLoading }: { data?: TickerDetailData; isLoading: boolean }) {
  if (isLoading) {
    return (
      <div className="text-[var(--md-sys-color-on-surface-variant)] text-xs font-mono font-bold animate-pulse">
        Memuat ringkasan overview...
      </div>
    );
  }
  if (!data || data.error) {
    return <div className="text-rose-500 text-xs font-mono font-bold">{data?.error || "Tidak ada data"}</div>;
  }

  const chartData = (data.price_chart || []).map((p) => {
    const sig = (data.signal_overlay || []).find((s) => s.date === p.date);
    return {
      ...p,
      signal: sig?.signal || null,
      signalScore: sig?.score ?? null,
      bandar_vwap: p.bandar_vwap ?? null,
    };
  });

  const retail = data.retail_summary || null;
  const profileFlow = data.profile_flow || [];
  const maxProfileAbs = Math.max(...profileFlow.map((r) => Math.abs(r.net)), 1);

  return (
    <div className="space-y-5">
      <div className="grid grid-cols-1 lg:grid-cols-[1.6fr_1fr] gap-5">
        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-5 shadow-xs">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-4">
            <div>
              <div className="flex items-center gap-2">
                <h3 className="text-sm sm:text-base font-bold text-[var(--md-sys-color-on-surface)] tracking-tight">
                  Price, Volume &amp; Bandar VWAP
                </h3>
                {data.bandar_pnl_pct !== null && data.bandar_pnl_pct !== undefined && (
                  <span
                    className={`text-[10px] font-mono font-bold px-2 py-0.5 rounded-full border ${
                      data.bandar_pnl_pct >= 0
                        ? "bg-emerald-500/10 text-emerald-500 border-emerald-500/30"
                        : "bg-rose-500/10 text-rose-500 border-rose-500/30"
                    }`}
                  >
                    Bandar {data.bandar_pnl_pct >= 0 ? `+${data.bandar_pnl_pct}% Floating Profit` : `${data.bandar_pnl_pct}% Underwater`}
                  </span>
                )}
              </div>
              <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
                Korelasi pergerakan harga close terhadap volume &amp; harga modal akumulasi bandar
              </p>
            </div>
            <div className="flex items-center gap-3 text-[10px] font-mono font-bold flex-wrap">
              <span className="flex items-center gap-1.5 text-amber-500">
                <span className="w-2.5 h-0.5 bg-amber-500"></span> Close
              </span>
              {data.bandar_vwap && (
                <span
                  className="flex items-center gap-1.5 text-pink-500"
                  title="Volume-Weighted Average Price dari Top 5 Broker Akumulasi"
                >
                  <span className="w-2.5 h-0.5 bg-pink-500 border-b border-dashed"></span> Bandar VWAP (Rp{" "}
                  {data.bandar_vwap.toLocaleString("id-ID")})
                </span>
              )}
              <span className="flex items-center gap-1.5 text-[var(--md-sys-color-on-surface-variant)]">
                <span className="w-2 h-2 bg-[var(--md-sys-color-outline)] rounded-xs"></span> Volume
              </span>
            </div>
          </div>

          <div className="h-72">
            <ResponsiveContainer width="100%" height="100%">
              <ComposedChart data={chartData} margin={{ top: 5, right: 5, left: 0, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="var(--md-sys-color-outline-variant)" opacity={0.6} />
                <XAxis
                  dataKey="date"
                  tick={{ fontSize: 10, fill: "var(--md-sys-color-on-surface-variant)", fontWeight: 600 }}
                  stroke="var(--md-sys-color-outline-variant)"
                />
                <YAxis
                  yAxisId="left"
                  tick={{ fontSize: 10, fill: "var(--md-sys-color-on-surface-variant)", fontWeight: 600 }}
                  stroke="var(--md-sys-color-outline-variant)"
                  domain={["auto", "auto"]}
                />
                <YAxis
                  yAxisId="right"
                  orientation="right"
                  tick={{ fontSize: 10, fill: "var(--md-sys-color-on-surface-variant)", fontWeight: 600 }}
                  stroke="var(--md-sys-color-outline-variant)"
                />
                <Tooltip
                  contentStyle={{
                    background: "var(--md-sys-color-surface-container-highest)",
                    border: "1px solid var(--md-sys-color-outline-variant)",
                    borderRadius: "12px",
                    fontSize: "11px",
                    color: "var(--md-sys-color-on-surface)",
                    fontWeight: 700,
                  }}
                  labelStyle={{ color: "var(--md-sys-color-on-surface-variant)", fontWeight: 800, marginBottom: "4px" }}
                  formatter={(value, name) => [
                    String(name) === "Bandar VWAP" ? `Rp ${Number(value ?? 0).toLocaleString("id-ID")}` : fmtRp(Number(value ?? 0)),
                    String(name ?? ""),
                  ]}
                />
                <Bar yAxisId="right" dataKey="volume" fill="var(--md-sys-color-outline)" opacity={0.3} radius={[3, 3, 0, 0]} />
                <Line yAxisId="left" type="monotone" dataKey="close" stroke="#f59e0b" strokeWidth={2.5} dot={false} />
                <Line
                  yAxisId="left"
                  type="monotone"
                  dataKey="bandar_vwap"
                  stroke="#ec4899"
                  strokeWidth={2}
                  strokeDasharray="4 4"
                  dot={false}
                  name="Bandar VWAP"
                />
              </ComposedChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="space-y-5">
          <TopBrokerSummaryCard rows={data.broker_summary || []} />
          <PricePerformanceCard rows={data.price_performance || []} />
        </div>
      </div>

      {retail && (
        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-5 shadow-xs">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 pb-3 border-b border-[var(--md-sys-color-outline-variant)]/60">
            <div className="flex items-center gap-2.5">
              <div className="w-8 h-8 rounded-xl bg-indigo-500/10 text-indigo-500 flex items-center justify-center shrink-0">
                <Icon icon="ph:users-three-bold" width="18" />
              </div>
              <div>
                <h3 className="text-sm sm:text-base font-bold text-[var(--md-sys-color-on-surface)] tracking-tight">
                  Retail vs. Smart Money Divergence Radar
                </h3>
                <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
                  Monitoring aksi broker ritel (YP, PD, XC, KK, NI, XL, SQ) vs institusi smart money
                </p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <span
                className={`text-xs font-bold px-3 py-1 rounded-full border shadow-2xs ${
                  retail.sentiment_tone === "positive"
                    ? "bg-emerald-500/10 text-emerald-500 border-emerald-500/30"
                    : retail.sentiment_tone === "negative"
                    ? "bg-rose-500/10 text-rose-500 border-rose-500/30"
                    : retail.sentiment_tone === "primary"
                    ? "bg-blue-500/10 text-blue-500 border-blue-500/30"
                    : "bg-amber-500/10 text-amber-500 border-amber-500/30"
                }`}
              >
                {retail.sentiment}
              </span>
            </div>
          </div>

          <p className="text-xs font-medium text-[var(--md-sys-color-on-surface-variant)] leading-relaxed mb-4">
            💡 {retail.sentiment_desc}
          </p>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-4">
            <div className="bg-[var(--md-sys-color-surface-container-high)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40">
              <div className="text-[10px] font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider mb-1">
                Net Smart Money (Akumulasi)
              </div>
              <div className="text-sm sm:text-base font-mono font-extrabold" style={{ color: signedColor(retail.smart_cum_net) }}>
                {fmtRp(retail.smart_cum_net)}
              </div>
            </div>

            <div className="bg-[var(--md-sys-color-surface-container-high)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40">
              <div className="text-[10px] font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider mb-1">
                Net Retail Flow (Partisipasi)
              </div>
              <div className="text-sm sm:text-base font-mono font-extrabold" style={{ color: signedColor(retail.retail_cum_net) }}>
                {fmtRp(retail.retail_cum_net)}
              </div>
            </div>

            <div className="bg-[var(--md-sys-color-surface-container-high)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40">
              <div className="text-[10px] font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider mb-1">
                Top Broker Ritel Teraktif
              </div>
              <div className="flex flex-wrap gap-1.5 mt-1">
                {(retail.top_retail_brokers || []).map((b, idx) => (
                  <span
                    key={idx}
                    className="text-[10px] font-mono font-bold px-2 py-0.5 rounded-md bg-[var(--md-sys-color-surface)] border border-[var(--md-sys-color-outline-variant)]"
                    style={{ color: signedColor(b.net) }}
                    title={`Buy: ${fmtRp(b.buy)} | Sell: ${fmtRp(b.sell)}`}
                  >
                    {b.broker}: {fmtRp(b.net)}
                  </span>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-[1.1fr_0.9fr] gap-5">
        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-5 shadow-xs">
          <h3 className="text-sm sm:text-base font-bold text-[var(--md-sys-color-on-surface)] mb-4 tracking-tight">
            Smart-Money Daily Flow
          </h3>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <ComposedChart data={data.smart_daily || []} margin={{ top: 5, right: 5, left: 0, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="var(--md-sys-color-outline-variant)" opacity={0.6} />
                <XAxis
                  dataKey="date"
                  tick={{ fontSize: 10, fill: "var(--md-sys-color-on-surface-variant)", fontWeight: 600 }}
                  stroke="var(--md-sys-color-outline-variant)"
                />
                <YAxis
                  yAxisId="left"
                  tick={{ fontSize: 10, fill: "var(--md-sys-color-on-surface-variant)", fontWeight: 600 }}
                  stroke="var(--md-sys-color-outline-variant)"
                />
                <YAxis
                  yAxisId="right"
                  orientation="right"
                  tick={{ fontSize: 10, fill: "var(--md-sys-color-on-surface-variant)", fontWeight: 600 }}
                  stroke="var(--md-sys-color-outline-variant)"
                />
                <Tooltip
                  contentStyle={{
                    background: "var(--md-sys-color-surface-container-highest)",
                    border: "1px solid var(--md-sys-color-outline-variant)",
                    borderRadius: "12px",
                    fontSize: "11px",
                    color: "var(--md-sys-color-on-surface)",
                    fontWeight: 700,
                  }}
                  labelStyle={{ color: "var(--md-sys-color-on-surface-variant)", fontWeight: 800, marginBottom: "4px" }}
                  formatter={(value, name) => [fmtRp(Number(value ?? 0)), String(name ?? "")]}
                />
                <Bar
                  yAxisId="left"
                  dataKey="smart_net"
                  shape={(rawProps: unknown) => {
                    const { x = 0, y = 0, width = 0, height = 0, payload } = rawProps as SmartBarProps;
                    const color = (payload?.smart_net ?? 0) >= 0 ? "var(--color-positive)" : "var(--color-negative)";
                    return <rect x={x} y={y} width={width} height={height} fill={color} opacity={0.9} rx={3} />;
                  }}
                />
                <Line yAxisId="right" type="monotone" dataKey="cumulative_net" stroke="#38bdf8" strokeWidth={2.5} dot={false} />
                <ReferenceLine yAxisId="left" y={0} stroke="var(--md-sys-color-outline)" strokeWidth={1} />
              </ComposedChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-5 shadow-xs">
          <h3 className="text-sm sm:text-base font-bold text-[var(--md-sys-color-on-surface)] mb-4 tracking-tight">
            Participant Profile Net Flow
          </h3>
          {profileFlow.length === 0 ? (
            <p className="text-xs font-mono font-medium text-[var(--md-sys-color-on-surface-variant)]">
              Tidak ada arus profil transaksi untuk periode ini.
            </p>
          ) : (
            <div className="space-y-4">
              {profileFlow.map((row, i) => {
                const width = Math.max(3, (Math.abs(row.net) / maxProfileAbs) * 100);
                return (
                  <div key={i}>
                    <div className="flex justify-between items-center text-xs mb-1.5">
                      <span className="font-bold uppercase tracking-wider" style={{ color: getProfileTextColor(row.label) }}>
                        {row.label}
                      </span>
                      <span className="font-mono font-bold" style={{ color: signedColor(row.net) }}>
                        {fmtRp(row.net)}
                      </span>
                    </div>
                    <div className="h-2 bg-[var(--md-sys-color-surface-container-high)] rounded-full overflow-hidden border border-[var(--md-sys-color-outline-variant)]">
                      <div
                        className="h-full rounded-full transition-all duration-700"
                        style={{ width: `${width}%`, backgroundColor: signedColor(row.net) }}
                      />
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>

      <ProfileBrokerDetailCard rows={data.profile_broker_detail || []} />
    </div>
  );
}
