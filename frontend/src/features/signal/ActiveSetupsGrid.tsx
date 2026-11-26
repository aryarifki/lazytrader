import { Link } from "react-router-dom";
import { Icon } from "@iconify/react";
import { getWyckoffBadge, getScoreBadge } from "./signalHelpers";
import type { ActiveSignal } from "./types";

interface ActiveSetupsGridProps {
  activeSignals: ActiveSignal[];
  activeLoading: boolean;
}

export default function ActiveSetupsGrid({ activeSignals, activeLoading }: ActiveSetupsGridProps) {
  return (
    <div className="space-y-3">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <div className="flex items-center gap-2">
          <div className="p-1.5 rounded-lg bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)]">
            <Icon icon="ph:crosshair-bold" width="18" />
          </div>
          <div>
            <h2 className="text-base font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-2">
              Actionable Setups (Open Position)
              <span className="text-[10px] font-mono px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20 font-semibold">
                {activeSignals.length} Setups
              </span>
            </h2>
            <p className="text-xs text-[var(--md-sys-color-on-surface-variant)]">
              Kalkulasi target profit (TP), batasan risiko (SL), dan rasio risk-to-reward algoritma SMtracker.
            </p>
          </div>
        </div>
      </div>

      {activeLoading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="h-44 rounded-2xl bg-[var(--md-sys-color-surface-container)] animate-pulse border border-[var(--md-sys-color-outline-variant)]"
            />
          ))}
        </div>
      ) : activeSignals.length === 0 ? (
        <div className="p-8 text-center bg-[var(--md-sys-color-surface-container)] rounded-2xl border border-[var(--md-sys-color-outline-variant)] text-xs text-[var(--md-sys-color-on-surface-variant)]">
          Tidak ada setup open aktif saat ini.
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3.5">
          {activeSignals.map((item) => {
            const phase = item.wyckoff_phase || "-";
            return (
              <div
                key={item.id}
                className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] hover:border-[var(--md-sys-color-primary)]/50 rounded-2xl p-4 transition-all duration-200 shadow-2xs relative flex flex-col justify-between group"
              >
                <div>
                  <div className="flex items-start justify-between gap-2 mb-2.5">
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center gap-1.5 flex-wrap">
                        <Link
                          to={`/${item.ticker}`}
                          className="font-mono font-black text-lg text-[var(--md-sys-color-on-surface)] group-hover:text-[var(--md-sys-color-primary)] transition-colors hover:underline"
                        >
                          {item.ticker}
                        </Link>
                        <span
                          className={`text-[10px] font-mono px-2 py-0.5 rounded-full border font-bold ${getWyckoffBadge(phase)}`}
                          title={item.wyckoff_name || undefined}
                        >
                          Phase {phase}
                          {item.wyckoff_name ? ` · ${item.wyckoff_name}` : ""}
                        </span>
                        {item.vcp_grade && item.vcp_grade !== "NONE" && (
                          <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-purple-500/15 text-purple-700 dark:text-purple-400 border border-purple-500/30 font-bold">
                            VCP {item.vcp_grade}
                          </span>
                        )}
                        {(item.confirmation_count ?? 1) > 1 && (
                          <span
                            className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-amber-500/15 text-amber-700 dark:text-amber-400 border border-amber-500/30 font-bold"
                            title={`Sinyal pertama ${item.signal_date}; konfirmasi terakhir ${item.last_confirmed_date ?? item.signal_date}`}
                          >
                            {item.confirmation_count}× konfirmasi
                          </span>
                        )}
                      </div>
                      <div className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] line-clamp-1 mt-0.5">
                        {item.company_name}
                      </div>
                      <div className="flex items-center gap-1.5 mt-1 flex-wrap">
                        {item.persistence_5d !== undefined && (
                          <span
                            className={`text-[9px] font-mono px-1.5 py-0.2 rounded border font-semibold ${
                              item.persistence_5d >= 0.6
                                ? "bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border-emerald-500/20"
                                : item.persistence_5d <= 0.4
                                ? "bg-rose-500/10 text-rose-700 dark:text-rose-400 border-rose-500/20"
                                : "bg-neutral-500/10 text-neutral-600 dark:text-neutral-400 border-neutral-500/20"
                            }`}
                          >
                            {item.persistence_5d >= 0.6 ? "🟢" : item.persistence_5d <= 0.4 ? "🔴" : "⚪"}{" "}
                            {(item.persistence_5d * 100).toFixed(0)}% 5D Flow
                          </span>
                        )}
                        {item.sector_notes && (
                          <span
                            className="text-[9px] font-mono px-1.5 py-0.2 rounded bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)] truncate max-w-[140px]"
                            title={item.sector_notes}
                          >
                            {item.sector_notes}
                          </span>
                        )}
                      </div>
                    </div>

                    <div className="text-right shrink-0">
                      <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-medium">Score</div>
                      <div className={`text-xs font-mono font-black px-2 py-0.5 rounded-md border ${getScoreBadge(item.composite_score)}`}>
                        {item.composite_score}
                      </div>
                      {item.ml_win_prob !== undefined && (
                        <div className="text-[9px] font-mono mt-0.5 text-[var(--md-sys-color-on-surface-variant)]">
                          {(() => {
                            if (item.ml_win_prob === null) return "N/A";
                            const norm = item.ml_win_prob > 1.0 ? item.ml_win_prob / 100.0 : item.ml_win_prob;
                            return `${(norm * 100).toFixed(0)}% Win`;
                          })()}
                        </div>
                      )}
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-2 p-2.5 rounded-xl bg-[var(--md-sys-color-surface-container-high)]/70 border border-[var(--md-sys-color-outline-variant)]/60 mb-3 text-center">
                    <div>
                      <div className="text-[9px] uppercase font-bold text-[var(--md-sys-color-on-surface-variant)]">Entry</div>
                      <div className="font-mono font-bold text-xs text-[var(--md-sys-color-on-surface)]">
                        {item.entry_price ? `Rp ${item.entry_price.toLocaleString("id-ID")}` : "-"}
                      </div>
                    </div>
                    <div>
                      <div className="text-[9px] uppercase font-bold text-emerald-700 dark:text-emerald-400">
                        TP (+{item.target_gain_pct}%)
                      </div>
                      <div className="font-mono font-bold text-xs text-emerald-700 dark:text-emerald-400">
                        {item.tp_price ? `Rp ${item.tp_price.toLocaleString("id-ID")}` : "-"}
                      </div>
                    </div>
                    <div>
                      <div className="text-[9px] uppercase font-bold text-rose-700 dark:text-rose-400">
                        SL (-{item.stop_risk_pct}%)
                      </div>
                      <div className="font-mono font-bold text-xs text-rose-700 dark:text-rose-400">
                        {item.sl_price ? `Rp ${item.sl_price.toLocaleString("id-ID")}` : "-"}
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center justify-between text-[11px] text-[var(--md-sys-color-on-surface-variant)] mb-1.5 font-mono">
                    <span>
                      R:R: <strong className="text-[var(--md-sys-color-primary)] font-bold">{item.risk_reward_ratio}:1</strong>
                    </span>
                    {item.rsi !== null && <span>RSI: {Number(item.rsi).toFixed(1)}</span>}
                    {item.cmf !== null && <span>CMF: {Number(item.cmf).toFixed(2)}</span>}
                  </div>

                  {item.smart_money_notes && (
                    <div
                      className="text-[10px] text-emerald-700 dark:text-emerald-400/90 font-mono mb-1.5 truncate"
                      title={item.smart_money_notes}
                    >
                      ⚡ <span className="opacity-80">SM:</span> {item.smart_money_notes}
                    </div>
                  )}

                  {item.gate_notes && (
                    <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] bg-[var(--md-sys-color-surface)]/60 px-2 py-1 rounded border border-[var(--md-sys-color-outline-variant)]/40 line-clamp-1">
                      <Icon icon="ph:info-duotone" className="inline mr-1 text-[var(--md-sys-color-primary)]" width="12" />
                      {item.gate_notes}
                    </div>
                  )}
                </div>

                <div className="mt-3 pt-2.5 border-t border-[var(--md-sys-color-outline-variant)]/60 flex items-center justify-between">
                  <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono">{item.signal_date}</span>
                  <Link
                    to={`/${item.ticker}`}
                    className="text-xs font-semibold text-[var(--md-sys-color-primary)] hover:underline inline-flex items-center gap-1"
                  >
                    Buka Analitik
                    <Icon icon="ph:arrow-right-bold" width="12" />
                  </Link>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
