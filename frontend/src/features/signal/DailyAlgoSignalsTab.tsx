import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import useSWR from "swr";
import { Icon } from "@iconify/react";
import ActiveSetupsGrid from "./ActiveSetupsGrid";
import { getWyckoffBadge, getScoreBadge } from "./signalHelpers";
import type { ActiveSignal, DailySignalRow } from "./types";

const fetcher = async (url: string) => {
  const res = await fetch(url, { credentials: "include" });
  if (!res.ok) {
    throw new Error(`HTTP ${res.status}: Gagal memuat data`);
  }
  return res.json();
};

export default function DailyAlgoSignalsTab() {
  const [search, setSearch] = useState("");
  const [selectedPhase, setSelectedPhase] = useState("ALL");
  const [selectedMlLabel, setSelectedMlLabel] = useState("ALL");
  const minComposite = 0;

  useEffect(() => {
    try {
      const savedPhase = localStorage.getItem("investowl_algo_phase");
      if (savedPhase) setSelectedPhase(savedPhase);
      const savedLabel = localStorage.getItem("investowl_algo_label");
      if (savedLabel) setSelectedMlLabel(savedLabel);
    } catch {}
  }, []);

  const handlePhaseChange = (phase: string) => {
    setSelectedPhase(phase);
    try {
      localStorage.setItem("investowl_algo_phase", phase);
    } catch {}
  };

  const handleLabelChange = (label: string) => {
    setSelectedMlLabel(label);
    try {
      localStorage.setItem("investowl_algo_label", label);
    } catch {}
  };

  const { data: activeData, isLoading: activeLoading } = useSWR<{ data?: ActiveSignal[] }>(
    "/api/signal/active?limit=30",
    fetcher,
    { refreshInterval: 60000, dedupingInterval: 30000 }
  );

  const { data: dailyData, isLoading: dailyLoading, error: dailyError } = useSWR<{
    latest_date?: string;
    data?: DailySignalRow[];
  }>("/api/signal/daily?limit=100", fetcher, { refreshInterval: 60000, dedupingInterval: 30000 });

  const activeSignals: ActiveSignal[] = activeData?.data || [];
  const dailySignals: DailySignalRow[] = useMemo(() => dailyData?.data || [], [dailyData?.data]);
  const latestDate = dailyData?.latest_date || "-";

  const filteredDaily = useMemo(() => {
    return dailySignals.filter((item) => {
      const matchSearch = !search || (item.ticker && item.ticker.toUpperCase().includes(search.toUpperCase()));
      const phase = (item.wyckoff_phase || item.features?.wyckoff_phase || item.features?.wp || "-").toUpperCase();
      const matchPhase = selectedPhase === "ALL" || phase === selectedPhase;
      const matchLabel = selectedMlLabel === "ALL" || item.ml_label === selectedMlLabel;
      const matchComposite = (item.composite_score ?? 0) >= minComposite;
      return Boolean(matchSearch && matchPhase && matchLabel && matchComposite);
    });
  }, [dailySignals, search, selectedPhase, selectedMlLabel, minComposite]);

  return (
    <div className="space-y-6">
      <ActiveSetupsGrid activeSignals={activeSignals} activeLoading={activeLoading} />

      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-5 shadow-xs space-y-4">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-3 pb-3 border-b border-[var(--md-sys-color-outline-variant)]">
          <div>
            <h2 className="text-base font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-2">
              <Icon icon="ph:ranking-duotone" className="text-[var(--md-sys-color-primary)]" width="20" />
              Daily Algorithmic Market Ranking
              <span className="text-[11px] font-mono font-medium px-2 py-0.5 rounded-full bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface-variant)]">
                Scan Date: {latestDate}
              </span>
            </h2>
            <p className="text-xs text-[var(--md-sys-color-on-surface-variant)]">
              Output model kuantitatif multi-faktor (Technical, Smart Money, Sector, Wyckoff Phase) dari{" "}
              <code className="text-[var(--md-sys-color-primary)]">scanner.py</code>.
            </p>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <div className="relative">
              <Icon
                icon="ph:magnifying-glass-duotone"
                className="absolute left-2.5 top-1/2 -translate-y-1/2 text-[var(--md-sys-color-on-surface-variant)]"
                width="14"
              />
              <input
                type="text"
                placeholder="Cari ticker..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-8 pr-3 py-1.5 text-xs rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] focus:outline-none focus:border-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-surface)] w-36 sm:w-44"
              />
            </div>

            <div className="flex items-center gap-1 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] p-1 rounded-xl text-xs font-semibold overflow-x-auto">
              <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] px-1.5 uppercase font-mono">Phase:</span>
              {["ALL", "A", "B", "C", "D", "E", "UNKNOWN"].map((phase) => (
                <button
                  key={phase}
                  type="button"
                  onClick={() => handlePhaseChange(phase)}
                  className={`px-2 py-0.5 rounded-lg text-xs transition-all whitespace-nowrap cursor-pointer ${
                    selectedPhase === phase
                      ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-2xs font-bold"
                      : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                  }`}
                >
                  {phase}
                </button>
              ))}
            </div>

            <select
              value={selectedMlLabel}
              onChange={(e) => handleLabelChange(e.target.value)}
              className="py-1.5 px-2.5 text-xs rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface)] focus:outline-none cursor-pointer"
            >
              <option value="ALL">Semua Label</option>
              <option value="WIN">ML Label: WIN</option>
              <option value="HOLD">ML Label: HOLD</option>
              <option value="LOSS">ML Label: LOSS</option>
            </select>
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse text-xs">
            <thead>
              <tr className="border-b border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface-variant)] font-bold uppercase text-[10px] tracking-wider">
                <th className="py-2.5 px-3">#</th>
                <th className="py-2.5 px-3">Ticker</th>
                <th className="py-2.5 px-3 text-center">Wyckoff</th>
                <th className="py-2.5 px-3 text-right">Composite</th>
                <th className="py-2.5 px-3 text-right">ML Win Prob</th>
                <th className="py-2.5 px-3 text-right">Smart Money</th>
                <th className="py-2.5 px-3 text-right">Technical</th>
                <th className="py-2.5 px-3 text-right">Sector</th>
                <th className="py-2.5 px-3 text-center">VCP</th>
                <th className="py-2.5 px-3">Gate Note / Analysis</th>
                <th className="py-2.5 px-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/60">
              {dailyLoading ? (
                <tr>
                  <td colSpan={11} className="py-8 text-center text-[var(--md-sys-color-on-surface-variant)] animate-pulse">
                    Memuat data sinyal harian algoritma...
                  </td>
                </tr>
              ) : dailyError ? (
                <tr>
                  <td colSpan={11} className="py-8 text-center text-rose-700 dark:text-rose-400">
                    Gagal memuat data sinyal: {dailyError.message || "Koneksi backend terputus"}. Silakan coba lagi.
                  </td>
                </tr>
              ) : filteredDaily.length === 0 ? (
                <tr>
                  <td colSpan={11} className="py-8 text-center text-[var(--md-sys-color-on-surface-variant)]">
                    Tidak ada sinyal yang sesuai dengan filter pencarian.
                  </td>
                </tr>
              ) : (
                filteredDaily.map((row, idx) => {
                  const phase = (row.wyckoff_phase || row.features?.wyckoff_phase || row.features?.wp || "-").toUpperCase();
                  const wyckoffName = row.wyckoff_name || row.features?.wyckoff_name || "";
                  const vcp = row.features?.vcp_grade || "NONE";
                  const p5d = row.persistence_5d ?? row.features?.persistence_5d;
                  const secNotes = row.sector_notes || row.features?.sector_notes;

                  return (
                    <tr key={row.ticker} className="hover:bg-[var(--md-sys-color-surface-container-high)]/50 transition-colors group">
                      <td className="py-2.5 px-3 font-mono text-[var(--md-sys-color-on-surface-variant)]">{idx + 1}</td>
                      <td className="py-2.5 px-3 font-mono font-bold text-[var(--md-sys-color-on-surface)] group-hover:text-[var(--md-sys-color-primary)]">
                        <Link to={`/${row.ticker}`} className="hover:underline">
                          {row.ticker}
                        </Link>
                      </td>
                      <td className="py-2.5 px-3 text-center">
                        <span
                          className={`text-[10px] font-mono px-2 py-0.5 rounded-full border font-bold inline-block ${getWyckoffBadge(phase)}`}
                          title={wyckoffName || undefined}
                        >
                          Phase {phase}
                          {wyckoffName && (
                            <span className="block text-[8px] opacity-85 font-normal tracking-tight">{wyckoffName}</span>
                          )}
                        </span>
                      </td>
                      <td className="py-2.5 px-3 text-right font-mono font-bold">
                        <span className={`px-2 py-0.5 rounded border ${getScoreBadge(row.composite_score)}`}>
                          {row.composite_score}
                        </span>
                      </td>
                      <td className="py-2.5 px-3 text-right font-mono font-semibold">
                        {(() => {
                          const raw = row.ml_win_prob;
                          if (raw === null || raw === undefined) {
                            return <span className="text-[var(--md-sys-color-on-surface-variant)] opacity-50">N/A</span>;
                          }
                          const norm = raw > 1.0 ? raw / 100.0 : raw;
                          const color =
                            norm >= 0.55
                              ? "text-emerald-700 dark:text-emerald-400"
                              : norm >= 0.45
                              ? "text-amber-700 dark:text-amber-400"
                              : "text-[var(--md-sys-color-on-surface-variant)]";
                          return <span className={color}>{(norm * 100).toFixed(1)}%</span>;
                        })()}
                      </td>
                      <td className="py-2.5 px-3 text-right font-mono text-[var(--md-sys-color-on-surface)]">
                        <div className="font-bold">{row.scores.smart_money}</div>
                        {p5d !== undefined && p5d !== null && (
                          <div
                            className={`text-[9px] ${
                              p5d >= 0.6
                                ? "text-emerald-700 dark:text-emerald-400"
                                : p5d <= 0.4
                                ? "text-rose-700 dark:text-rose-400"
                                : "text-neutral-600 dark:text-neutral-400"
                            }`}
                          >
                            {(p5d * 100).toFixed(0)}% flow
                          </div>
                        )}
                      </td>
                      <td className="py-2.5 px-3 text-right font-mono text-[var(--md-sys-color-on-surface)]">{row.scores.technical}</td>
                      <td className="py-2.5 px-3 text-right font-mono text-[var(--md-sys-color-on-surface)]">
                        <div className="font-bold">{row.scores.sector}</div>
                        {secNotes && (
                          <div className="text-[8px] text-[var(--md-sys-color-on-surface-variant)] truncate max-w-[85px] ml-auto" title={secNotes}>
                            {secNotes}
                          </div>
                        )}
                      </td>
                      <td className="py-2.5 px-3 text-center font-mono">
                        {vcp !== "NONE" ? (
                          <span className="px-1.5 py-0.5 rounded bg-purple-500/15 text-purple-700 dark:text-purple-400 border border-purple-500/30 text-[10px] font-bold">
                            {vcp}
                          </span>
                        ) : (
                          <span className="text-[var(--md-sys-color-on-surface-variant)] opacity-40">-</span>
                        )}
                      </td>
                      <td className="py-2.5 px-3 text-[11px] text-[var(--md-sys-color-on-surface-variant)] max-w-xs truncate" title={row.gate_notes}>
                        {row.gate_notes || "Clear"}
                      </td>
                      <td className="py-2.5 px-3 text-center">
                        <Link
                          to={`/${row.ticker}`}
                          className="px-2.5 py-1 rounded-lg bg-[var(--md-sys-color-surface-container-highest)] hover:bg-[var(--md-sys-color-primary)] hover:text-[var(--md-sys-color-on-primary)] transition-all font-semibold text-[11px] inline-flex items-center gap-1"
                        >
                          Detail
                          <Icon icon="ph:arrow-square-out-bold" width="12" />
                        </Link>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
