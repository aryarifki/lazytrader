import type { StockAnalyticsStatistics } from "./foreignProTypes";

interface ForeignProStatisticsBoxProps {
  statType: string;
  setStatType: (v: string) => void;
  signalStat: string;
  setSignalStat: (v: string) => void;
  stats: StockAnalyticsStatistics | undefined;
}

export default function ForeignProStatisticsBox({
  statType,
  setStatType,
  signalStat,
  setSignalStat,
  stats,
}: ForeignProStatisticsBoxProps) {
  const corrRows = [
    { label: "Foreign Flow:", val: stats?.correlations?.ff, color: "#FF6820" },
    { label: "Foreign & Co:", val: stats?.correlations?.fco, color: "#D62839" },
    { label: "Korban Flow:", val: stats?.correlations?.kf, color: "#2F7676" },
    { label: "Broker XL:", val: stats?.correlations?.xl, color: "#3A7BD5" },
    { label: "Semi Ritel:", val: stats?.correlations?.semiritel, color: "#1DB87A" },
  ];

  return (
    <div className="pt-1 border-t border-[var(--md-sys-color-outline-variant)]/40 space-y-1.5">
      <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
        Statistics
      </div>
      <select
        id="statTypeSelect"
        value={statType}
        onChange={(e) => setStatType(e.target.value)}
        className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/50 rounded p-1 text-xs text-[var(--md-sys-color-on-surface)]"
      >
        <option value="indicatorCorrelation">Main Indicator Correlation</option>
        <option value="consecutiveInflow">Consecutive Inflow</option>
        <option value="inflowHargaNaik">Inflow &amp; Harga Naik</option>
        <option value="markDown">Mark Down</option>
        <option value="signalStats">Statistik Sinyal</option>
      </select>

      {statType === "signalStats" && (
        <select
          id="signalStatSelect"
          value={signalStat}
          onChange={(e) => setSignalStat(e.target.value)}
          className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/50 rounded p-1 text-[11px] text-[var(--md-sys-color-on-surface)]"
        >
          <option value="inflow_streak">Inflow beruntun (≥4 hr)</option>
          <option value="outflow_streak">Outflow beruntun (≥5 hr)</option>
          <option value="inflow_3m">Inflow terbesar 3 bln</option>
          <option value="outflow_3m">Outflow terbesar 3 bln</option>
        </select>
      )}

      <div id="statWrap" className="bg-[var(--md-sys-color-surface-container-high)] p-2 rounded text-xs space-y-1.5 font-mono">
        {statType === "indicatorCorrelation" && (
          <div id="statRows" className="space-y-1">
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-bold border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1">
              Indicator Correlation (vs Close)
            </div>
            {corrRows.map((cr) => (
              <div key={cr.label} className="flex justify-between text-[11px]">
                <span className="text-[var(--md-sys-color-on-surface)]">{cr.label}</span>
                <b style={{ color: cr.color }}>
                  {cr.val !== undefined ? (cr.val >= 0 ? `+${cr.val}` : cr.val) : "—"}
                </b>
              </div>
            ))}
          </div>
        )}

        {statType === "consecutiveInflow" && (
          <div id="statRows" className="space-y-1">
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-bold border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1">
              Average Consecutive Inflow
            </div>
            {(["50d", "100d", "200d"] as const).map((k) => (
              <div key={k} className="flex justify-between text-[11px]">
                <span className="text-[var(--md-sys-color-on-surface)]">Periode {k.replace("d", " hr")}:</span>
                <b className="text-[var(--color-positive)]">{stats?.consecutive_inflow?.[k] || "—"} hari</b>
              </div>
            ))}
          </div>
        )}

        {statType === "inflowHargaNaik" && (
          <div id="statRows" className="space-y-1">
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-bold border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1">
              Inflow &amp; Harga Naik
            </div>
            {(["50d", "100d", "200d"] as const).map((k) => (
              <div key={k} className="flex justify-between text-[11px]">
                <span className="text-[var(--md-sys-color-on-surface)]">{k.replace("d", " hr")}:</span>
                <b className="text-[var(--color-positive)]">
                  Inflow {stats?.inflow_harga_naik?.[k]?.count || 0}x ({stats?.inflow_harga_naik?.[k]?.pct ?? 0}%)
                </b>
              </div>
            ))}
          </div>
        )}

        {statType === "markDown" && (
          <div id="statRows" className="space-y-1">
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-bold border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1">
              Mark Down (Harga Turun, FF Inflow)
            </div>
            {(["50d", "100d", "200d"] as const).map((k) => (
              <div key={k} className="flex justify-between text-[11px]">
                <span className="text-[var(--md-sys-color-on-surface)]">{k.replace("d", " hr")}:</span>
                <b className="text-amber-600 dark:text-[#FFC24D]">Mark Down {stats?.mark_down?.[k] || 0}x</b>
              </div>
            ))}
          </div>
        )}

        {statType === "signalStats" && (
          <div id="statRows" className="space-y-1">
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-bold border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1">
              Statistik Sinyal ({signalStat})
            </div>
            <div className="flex justify-between text-[11px]">
              <span className="text-[var(--md-sys-color-on-surface)]">Kejadian:</span>
              <b className="text-[var(--md-sys-color-on-surface)]">{stats?.signal_stats?.massive_acc?.events || 0}x</b>
            </div>
            <div className="flex justify-between text-[11px]">
              <span className="text-[var(--md-sys-color-on-surface)]">Horizon 5 hr:</span>
              <b className="text-[var(--color-positive)]">{stats?.signal_stats?.massive_acc?.horizons?.["5"]?.pct || 0}% Naik</b>
            </div>
            <div className="flex justify-between text-[11px]">
              <span className="text-[var(--md-sys-color-on-surface)]">Horizon 20 hr:</span>
              <b className="text-[var(--color-positive)]">{stats?.signal_stats?.massive_acc?.horizons?.["20"]?.pct || 0}% Naik</b>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
