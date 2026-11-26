import { bandarBadgeClass } from "./mmdBadges";
import type { MmdAnalysis, MmdGroupInfo } from "./mmdTypes";

interface MmdBandarResultTabProps {
  analysisData: MmdAnalysis | null;
  selectedTicker: string;
  onSelectTicker: (ticker: string) => void;
  onOpenBroker: (code: string) => void;
}

function GroupPerformanceCard({
  group,
  tagLabel,
  accentDotClass,
  accentTextClass,
  selectedTicker,
  onSelectTicker,
}: {
  group: MmdGroupInfo;
  tagLabel: string;
  accentDotClass: string;
  accentTextClass: string;
  selectedTicker: string;
  onSelectTicker: (ticker: string) => void;
}) {
  return (
    <div className="bg-[var(--md-sys-color-surface-container-low)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40 space-y-2">
      <div className="flex items-center justify-between pb-1.5 border-b border-[var(--md-sys-color-outline-variant)]/40">
        <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
          <span className={`w-2 h-2 rounded-full ${accentDotClass}`} />
          {group.title}
        </span>
        <span className={`text-[9px] font-mono font-bold ${accentTextClass}`}>{tagLabel}</span>
      </div>

      <table className="w-full text-xs font-mono text-left">
        <thead>
          <tr className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/30">
            <th className="py-1">Saham</th>
            <th className="py-1 text-right">Kinerja Hari Ini</th>
            <th className="py-1 text-right">Kinerja 1 Bulan</th>
            <th className="py-1 text-right">Avg Value</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/20">
          {(group.members || []).map((m) => {
            const isActive = m.ticker === selectedTicker;
            const avgVal = m.avg_val || 0;
            const valFmt =
              avgVal >= 1e12 ? `${(avgVal / 1e12).toFixed(2)} T` : `${(avgVal / 1e9).toFixed(2)} M`;
            const todayVal = m.today_pct !== undefined ? Number(m.today_pct) : 0;
            const m1Val =
              m.m1_pct !== undefined
                ? Number(m.m1_pct)
                : m.month_pct !== undefined
                  ? Number(m.month_pct)
                  : 0;
            return (
              <tr
                key={m.ticker}
                onClick={() => onSelectTicker(m.ticker)}
                className={`cursor-pointer transition hover:bg-[var(--md-sys-color-surface-container-high)] ${
                  isActive ? "bg-[var(--md-sys-color-primary)]/10 font-bold" : ""
                }`}
              >
                <td className="py-1 font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1">
                  {isActive && <span className="text-[var(--md-sys-color-primary)]">▶</span>}
                  <span className={isActive ? "text-[var(--md-sys-color-primary)]" : accentTextClass}>
                    {m.ticker}
                  </span>
                </td>
                <td
                  className={`py-1 text-right font-bold ${
                    todayVal >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                  }`}
                >
                  {todayVal > 0 ? `+${todayVal.toFixed(1)}%` : `${todayVal.toFixed(1)}%`}
                </td>
                <td
                  className={`py-1 text-right font-bold ${
                    m1Val >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                  }`}
                >
                  {m1Val > 0 ? `+${m1Val.toFixed(1)}%` : `${m1Val.toFixed(1)}%`}
                </td>
                <td className="py-1 text-right text-[var(--md-sys-color-on-surface-variant)]">{valFmt}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

export default function MmdBandarResultTab({
  analysisData,
  selectedTicker,
  onSelectTicker,
  onOpenBroker,
}: MmdBandarResultTabProps) {
  const badgeClass = analysisData?.bandar_action?.cls || "undetected";
  const weekly = analysisData?.previous_results?.weekly || [];
  const monthly = analysisData?.previous_results?.monthly || [];

  const periodPills = [
    { title: "1 MINGGU LALU", item: weekly[1], fallbackBadge: undefined },
    { title: "2 MINGGU LALU", item: weekly[2], fallbackBadge: undefined },
    { title: "BULAN INI", item: monthly[0], fallbackBadge: analysisData?.bandar_action },
    { title: "BULAN LALU", item: monthly[1], fallbackBadge: undefined },
  ];

  const primaryGroup =
    analysisData?.group_performance?.conglomerate || analysisData?.group_performance?.sector;
  const primaryIsConglomerate = Boolean(analysisData?.group_performance?.conglomerate);
  const secondarySector =
    analysisData?.group_performance?.conglomerate && analysisData?.group_performance?.sector
      ? analysisData.group_performance.sector
      : null;

  const commentary = analysisData?.ai_commentary || analysisData?.commentary;

  return (
    <div className="space-y-3">
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
        {periodPills.map(({ title, item, fallbackBadge }) => {
          const badge = item?.badge || fallbackBadge;
          const pct = item?.change_pct;
          return (
            <div
              key={title}
              className="bg-[var(--md-sys-color-surface-container-low)] p-2 rounded-xl border border-[var(--md-sys-color-outline-variant)]/40 text-center"
            >
              <div className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] uppercase font-bold">
                {title}
              </div>
              <span className={`${bandarBadgeClass(badge?.cls)} text-[10px] py-0.5 px-2 mt-1 min-w-[76px]`}>
                {badge?.label || "—"}
              </span>
              <div
                className={`text-[10px] font-mono-nums mt-0.5 font-bold ${
                  (pct || 0) >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                }`}
              >
                {pct !== undefined ? `${pct > 0 ? "+" : ""}${pct}%` : "—"}
              </div>
            </div>
          );
        })}
      </div>

      <div className="bg-[var(--md-sys-color-surface-container-low)] rounded-xl p-3.5 border border-[var(--md-sys-color-outline-variant)]/40 space-y-2.5">
        <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono-nums border-b border-[var(--md-sys-color-outline-variant)]/60 pb-1.5 flex justify-between">
          <span>Periode : {analysisData?.date_from || "—"} ~ {analysisData?.date_to || "—"}</span>
          <span className="text-[var(--md-sys-color-primary)] font-bold">Bandarmologi Engine</span>
        </div>

        {[
          {
            label: "Buyer × Seller",
            ratio: analysisData?.bandar_table?.[3]?.ratio || 0,
            cls: analysisData?.bandar_table?.[3]?.badge?.cls,
            badgeLabel: analysisData?.bandar_table?.[3]?.badge?.label || "—",
          },
          {
            label: "Broker Pengendali",
            ratio: analysisData?.pengendali_stats?.ratio || 0,
            cls: analysisData?.pengendali_stats?.action?.cls,
            badgeLabel: analysisData?.pengendali_stats?.action?.label || "—",
          },
          {
            label: "Retail Action",
            ratio: analysisData?.retail_ratio || 0,
            cls: analysisData?.retail_action?.cls,
            badgeLabel: analysisData?.retail_action?.label || "Do Nothing",
          },
        ].map((row) => (
          <div key={row.label} className="flex items-center justify-between text-xs py-1 border-b border-[var(--md-sys-color-outline-variant)]/30 font-mono-nums">
            <span className="font-semibold text-[var(--md-sys-color-on-surface)]">{row.label}</span>
            <span className="text-[var(--md-sys-color-on-surface-variant)]">{row.ratio}%</span>
            <span className={`${bandarBadgeClass(row.cls)} text-[10px] py-0.5 px-2 min-w-[84px]`}>
              {row.badgeLabel}
            </span>
          </div>
        ))}

        <div className="flex items-center justify-between pt-1">
          <span className="text-sm font-extrabold text-[var(--md-sys-color-on-surface)]">Verdict</span>
          <span className={`${bandarBadgeClass(badgeClass)} text-xs py-1.5 px-3.5 font-black shadow-xs min-w-[104px]`}>
            {analysisData?.bandar_action?.label || "—"}
          </span>
        </div>
      </div>

      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
            <span className="w-2 h-2 rounded-full bg-[var(--md-sys-color-primary)]" />
            AI Commentary
          </span>
          <span className="text-[9px] text-[var(--md-sys-color-on-primary-container)] bg-[var(--md-sys-color-primary-container)] font-bold px-2 py-0.5 rounded">
            Versi 3.0
          </span>
        </div>

        <div className="space-y-1.5 bg-[var(--md-sys-color-surface-container-low)] p-3.5 rounded-xl border border-[var(--md-sys-color-outline-variant)]/40 text-xs text-[var(--md-sys-color-on-surface-variant)]">
          {commentary && commentary.length > 0 ? (
            commentary.map((c, idx) => (
              <p key={idx} className="leading-relaxed">
                {c}
              </p>
            ))
          ) : (
            <p>Memuat narasi analisis bandarmologi...</p>
          )}
        </div>
      </div>

      {primaryGroup && (
        <GroupPerformanceCard
          group={primaryGroup}
          tagLabel={primaryIsConglomerate ? "KONGLOMERASI" : "SEKTOR"}
          accentDotClass="bg-teal-600 dark:bg-[#2FD4C6]"
          accentTextClass="text-teal-700 dark:text-[#2FD4C6]"
          selectedTicker={selectedTicker}
          onSelectTicker={onSelectTicker}
        />
      )}

      {secondarySector && (
        <GroupPerformanceCard
          group={secondarySector}
          tagLabel="SEKTOR"
          accentDotClass="bg-indigo-600 dark:bg-[#818CF8]"
          accentTextClass="text-indigo-700 dark:text-[#818CF8]"
          selectedTicker={selectedTicker}
          onSelectTicker={onSelectTicker}
        />
      )}

      <div className="bg-[var(--md-sys-color-surface-container-low)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40 space-y-2">
        <div className="flex items-center justify-between pb-1.5 border-b border-[var(--md-sys-color-outline-variant)]/40">
          <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
            <span className="w-2 h-2 rounded-full bg-amber-500 dark:bg-[#FFC24D]" />
            TOP 10 COMBINED BROKER TABLE
          </span>
          <span className="text-[9px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
            Net Lot &amp; Acc/Dist
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-xs font-mono text-left">
            <thead>
              <tr className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/30">
                <th className="py-1.5 px-2">Broker Buy</th>
                <th className="py-1.5 px-2 text-right">Net Buy</th>
                <th className="py-1.5 px-1 text-center">No</th>
                <th className="py-1.5 px-2 text-right">Net Sell</th>
                <th className="py-1.5 px-2">Broker Sell</th>
                <th className="py-1.5 px-2 text-right">Acc/Dist</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/20 font-mono-nums">
              {(analysisData?.top10_combined || []).map((r) => (
                <tr key={r.rank} className="hover:bg-[var(--md-sys-color-surface-container-high)] transition">
                  <td className="py-1.5 px-2 font-bold">
                    {r.buyer_broker !== "—" ? (
                      <button
                        type="button"
                        onClick={() => onOpenBroker(r.buyer_broker)}
                        className="text-[var(--md-sys-color-on-surface)] hover:text-[var(--md-sys-color-primary)] hover:underline cursor-pointer focus:outline-none flex items-center gap-1"
                      >
                        <span>{r.buyer_broker}</span>
                        <span className="text-[8px] text-[var(--md-sys-color-on-surface-variant)]">
                          ({r.buyer_kategori})
                        </span>
                      </button>
                    ) : (
                      "—"
                    )}
                  </td>
                  <td className="py-1.5 px-2 text-right text-[var(--color-positive)] font-semibold">
                    {r.buyer_net_fmt}
                  </td>
                  <td className="py-1.5 px-1 text-center text-[var(--md-sys-color-outline)] font-bold">
                    {r.rank}
                  </td>
                  <td className="py-1.5 px-2 text-right text-[var(--color-negative)] font-semibold">
                    {r.seller_net_fmt}
                  </td>
                  <td className="py-1.5 px-2 font-bold">
                    {r.seller_broker !== "—" ? (
                      <button
                        type="button"
                        onClick={() => onOpenBroker(r.seller_broker)}
                        className="text-[var(--md-sys-color-on-surface)] hover:text-[var(--md-sys-color-primary)] hover:underline cursor-pointer focus:outline-none flex items-center gap-1"
                      >
                        <span>{r.seller_broker}</span>
                        <span className="text-[8px] text-[var(--md-sys-color-on-surface-variant)]">
                          ({r.seller_kategori})
                        </span>
                      </button>
                    ) : (
                      "—"
                    )}
                  </td>
                  <td
                    className={`py-1.5 px-2 text-right font-bold ${
                      r.acc_dist >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                    }`}
                  >
                    {r.acc_dist_fmt}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
