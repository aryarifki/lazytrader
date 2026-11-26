import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import useSWR from "swr";
import { Icon } from "@iconify/react";
import { fmtRp, signedColor } from "@/lib/format";
import ConfluenceFilters from "./ConfluenceFilters";
import { getBandarSignalBadge, getHmmBadge } from "./signalHelpers";
import type { ConfluenceRow, ConfluenceSortKey } from "./types";

const fetcher = async (url: string) => {
  const res = await fetch(url, { credentials: "include" });
  if (!res.ok) {
    throw new Error(`HTTP ${res.status}: Gagal memuat data`);
  }
  return res.json();
};

export default function ConfluenceScreenerTab() {
  const [universe, setUniverse] = useState("lq45");
  const [windowDays, setWindowDays] = useState(20);
  const [selectedTag, setSelectedTag] = useState<string>("ALL");
  const [selectedSignal, setSelectedSignal] = useState<string>("ALL");
  const [minConviction, setMinConviction] = useState<number>(0);
  const [minLiq, setMinLiq] = useState<number>(0);
  const [search, setSearch] = useState("");
  const [sortBy, setSortBy] = useState<ConfluenceSortKey>("conviction");
  const [sortAsc, setSortAsc] = useState(false);

  useEffect(() => {
    try {
      const savedUniv = localStorage.getItem("investowl_confluence_universe");
      if (savedUniv) setUniverse(savedUniv);
      const savedDays = localStorage.getItem("investowl_confluence_window");
      if (savedDays) setWindowDays(Number(savedDays));
      const savedTag = localStorage.getItem("investowl_confluence_tag");
      if (savedTag) setSelectedTag(savedTag);
      const savedSig = localStorage.getItem("investowl_confluence_signal");
      if (savedSig) setSelectedSignal(savedSig);
      const savedLiq = localStorage.getItem("investowl_confluence_min_liq");
      if (savedLiq) setMinLiq(Number(savedLiq));
      const savedSort = localStorage.getItem("investowl_confluence_sort") as ConfluenceSortKey | null;
      if (savedSort && ["conviction", "distance", "foreign", "gorengan"].includes(savedSort)) {
        setSortBy(savedSort);
      }
    } catch {}
  }, []);

  const handleUniverseChange = (val: string) => {
    setUniverse(val);
    try {
      localStorage.setItem("investowl_confluence_universe", val);
    } catch {}
  };
  const handleWindowDaysChange = (val: number) => {
    setWindowDays(val);
    try {
      localStorage.setItem("investowl_confluence_window", String(val));
    } catch {}
  };
  const handleTagChange = (val: string) => {
    setSelectedTag(val);
    try {
      localStorage.setItem("investowl_confluence_tag", val);
    } catch {}
  };
  const handleSignalChange = (val: string) => {
    setSelectedSignal(val);
    try {
      localStorage.setItem("investowl_confluence_signal", val);
    } catch {}
  };
  const handleMinLiqChange = (val: number) => {
    setMinLiq(val);
    try {
      localStorage.setItem("investowl_confluence_min_liq", String(val));
    } catch {}
  };

  const endpoint = `/api/signal/confluence-screener?universe_mode=${universe}&window_days=${windowDays}`;
  const { data, error, isLoading, isValidating } = useSWR<{ data?: ConfluenceRow[] }>(endpoint, fetcher, {
    revalidateOnFocus: false,
    dedupingInterval: 60000,
  });

  const rawRows: ConfluenceRow[] = useMemo(() => data?.data || [], [data?.data]);

  const processedRows = useMemo(() => {
    const result = rawRows.filter((row) => {
      const matchSearch =
        !search ||
        (row.ticker && row.ticker.toUpperCase().includes(search.toUpperCase())) ||
        (row.company_name && row.company_name.toLowerCase().includes(search.toLowerCase())) ||
        (row.group_name && row.group_name.toLowerCase().includes(search.toLowerCase()));

      const matchTag = selectedTag === "ALL" || (row.quick_tags && row.quick_tags.includes(selectedTag));
      const matchSignal = selectedSignal === "ALL" || (row.signal && row.signal.toLowerCase().includes(selectedSignal.toLowerCase()));
      const matchConviction = (row.conviction_score ?? 0) >= minConviction;
      const matchLiq = !minLiq || (row.total_value ?? 0) >= minLiq;

      return Boolean(matchSearch && matchTag && matchSignal && matchConviction && matchLiq);
    });

    result.sort((a, b) => {
      let valA = 0;
      let valB = 0;
      if (sortBy === "conviction") {
        valA = a.conviction_score;
        valB = b.conviction_score;
      } else if (sortBy === "distance") {
        valA = a.bandar_price_distance_pct;
        valB = b.bandar_price_distance_pct;
      } else if (sortBy === "foreign") {
        valA = a.smart_foreign_net;
        valB = b.smart_foreign_net;
      } else if (sortBy === "gorengan") {
        valA = a.gorengan_net;
        valB = b.gorengan_net;
      }
      return sortAsc ? valA - valB : valB - valA;
    });

    return result;
  }, [rawRows, search, selectedTag, selectedSignal, minConviction, minLiq, sortBy, sortAsc]);

  const handleSort = (col: ConfluenceSortKey) => {
    if (sortBy === col) {
      setSortAsc(!sortAsc);
    } else {
      setSortBy(col);
      setSortAsc(false);
      try {
        localStorage.setItem("investowl_confluence_sort", col);
      } catch {}
    }
  };

  return (
    <div className="space-y-5">
      <ConfluenceFilters
        selectedTag={selectedTag}
        onTagChange={handleTagChange}
        universe={universe}
        onUniverseChange={handleUniverseChange}
        windowDays={windowDays}
        onWindowDaysChange={handleWindowDaysChange}
        selectedSignal={selectedSignal}
        onSignalChange={handleSignalChange}
        minLiq={minLiq}
        onMinLiqChange={handleMinLiqChange}
        minConviction={minConviction}
        onMinConvictionChange={setMinConviction}
        search={search}
        onSearchChange={setSearch}
        filteredCount={processedRows.length}
        totalCount={rawRows.length}
      />

      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-5 shadow-xs space-y-3">
        <div className="flex items-center justify-between">
          <div className="text-xs text-[var(--md-sys-color-on-surface-variant)]">
            Menampilkan <strong className="text-[var(--md-sys-color-on-surface)] font-bold">{processedRows.length}</strong> emiten
            dengan metrik bandarmology dan foreign confluence.
          </div>
          {isValidating && (
            <div className="inline-flex items-center gap-1.5 text-[11px] font-mono text-[var(--md-sys-color-primary)]">
              <Icon icon="ph:spinner-bold" className="animate-spin" width="14" />
              <span>Memperbarui data...</span>
            </div>
          )}
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse text-xs">
            <thead>
              <tr className="border-b border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface-variant)] font-bold uppercase text-[10px] tracking-wider">
                <th className="py-2.5 px-3">Ticker</th>
                <th className="py-2.5 px-3 text-right">Harga</th>
                <th className="py-2.5 px-3 text-right">Modal Bandar</th>
                <th
                  className="py-2.5 px-3 text-right cursor-pointer hover:text-[var(--md-sys-color-primary)] transition-colors select-none"
                  onClick={() => handleSort("distance")}
                >
                  <span className="inline-flex items-center gap-1">
                    Jarak ke Modal
                    <Icon
                      icon={sortBy === "distance" ? (sortAsc ? "ph:sort-ascending-bold" : "ph:sort-descending-bold") : "ph:arrows-down-up-bold"}
                      width="12"
                    />
                  </span>
                </th>
                <th className="py-2.5 px-3 text-center">Sinyal Bandar</th>
                <th
                  className="py-2.5 px-3 text-right cursor-pointer hover:text-[var(--md-sys-color-primary)] transition-colors select-none"
                  onClick={() => handleSort("conviction")}
                >
                  <span className="inline-flex items-center gap-1">
                    Conviction
                    <Icon
                      icon={sortBy === "conviction" ? (sortAsc ? "ph:sort-ascending-bold" : "ph:sort-descending-bold") : "ph:arrows-down-up-bold"}
                      width="12"
                    />
                  </span>
                </th>
                <th className="py-2.5 px-3 text-center">Top Brokers</th>
                <th
                  className="py-2.5 px-3 text-right cursor-pointer hover:text-[var(--md-sys-color-primary)] transition-colors select-none"
                  onClick={() => handleSort("foreign")}
                >
                  <span className="inline-flex items-center gap-1" title="Net Flow Broker Asing Pintar (AK, BK, KZ, CS)">
                    Smart Foreign
                    <Icon
                      icon={sortBy === "foreign" ? (sortAsc ? "ph:sort-ascending-bold" : "ph:sort-descending-bold") : "ph:arrows-down-up-bold"}
                      width="12"
                    />
                  </span>
                </th>
                <th
                  className="py-2.5 px-3 text-right cursor-pointer hover:text-[var(--md-sys-color-primary)] transition-colors select-none"
                  onClick={() => handleSort("gorengan")}
                >
                  <span className="inline-flex items-center gap-1" title="Net Flow Broker Gorengan / Retail Cepat (MG, SQ, dll)">
                    Gorengan Flow
                    <Icon
                      icon={sortBy === "gorengan" ? (sortAsc ? "ph:sort-ascending-bold" : "ph:sort-descending-bold") : "ph:arrows-down-up-bold"}
                      width="12"
                    />
                  </span>
                </th>
                <th className="py-2.5 px-3 text-center">Foreign HMM</th>
                <th className="py-2.5 px-3 text-right">Free Float</th>
                <th className="py-2.5 px-3">Tags</th>
                <th className="py-2.5 px-3 text-center">Navigasi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/60">
              {isLoading ? (
                <tr>
                  <td colSpan={13} className="py-12 text-center text-[var(--md-sys-color-on-surface-variant)] animate-pulse">
                    Menghitung parameter konfluensi bandarmology &amp; foreign flow...
                  </td>
                </tr>
              ) : error ? (
                <tr>
                  <td colSpan={13} className="py-12 text-center text-rose-700 dark:text-rose-400">
                    Gagal memuat data screener: {error.message || "Koneksi backend terputus"}. Silakan coba lagi.
                  </td>
                </tr>
              ) : processedRows.length === 0 ? (
                <tr>
                  <td colSpan={13} className="py-12 text-center text-[var(--md-sys-color-on-surface-variant)]">
                    Tidak ada emiten yang memenuhi kriteria filter saat ini.
                  </td>
                </tr>
              ) : (
                processedRows.map((row) => {
                  const dist = row.bandar_price_distance_pct;
                  const isStealthDiscount = dist < 0;
                  const isSafeNearBandar = dist >= 0 && dist <= 5;
                  const isOverextended = dist > 10;

                  return (
                    <tr
                      key={row.ticker}
                      className="hover:bg-[var(--md-sys-color-surface-container-high)]/50 transition-colors group"
                    >
                      <td className="py-2.5 px-3">
                        <Link
                          to={`/${row.ticker}`}
                          className="font-mono font-bold text-sm text-[var(--md-sys-color-on-surface)] group-hover:text-[var(--md-sys-color-primary)] hover:underline block"
                        >
                          {row.ticker}
                        </Link>
                        <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] truncate max-w-[130px]">
                          {row.company_name}
                        </div>
                      </td>

                      <td className="py-2.5 px-3 text-right font-mono font-semibold text-[var(--md-sys-color-on-surface)]">
                        <div>{row.current_price ? `Rp ${row.current_price.toLocaleString("id-ID")}` : "-"}</div>
                        {row.total_value && row.total_value > 0 ? (
                          <div
                            className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-normal tabular-nums"
                            title="Nilai Transaksi Harian"
                          >
                            Vol: {fmtRp(row.total_value)}
                          </div>
                        ) : null}
                      </td>

                      <td className="py-2.5 px-3 text-right font-mono font-medium text-[var(--md-sys-color-on-surface-variant)]">
                        {row.bandar_avg_price ? `Rp ${row.bandar_avg_price.toLocaleString("id-ID")}` : "-"}
                      </td>

                      <td className="py-2.5 px-3 text-right font-mono font-bold">
                        <span
                          className={`px-1.5 py-0.5 rounded text-[11px] ${
                            isStealthDiscount
                              ? "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400"
                              : isSafeNearBandar
                              ? "bg-cyan-500/15 text-cyan-700 dark:text-cyan-400"
                              : isOverextended
                              ? "bg-rose-500/15 text-rose-700 dark:text-rose-400"
                              : "text-[var(--md-sys-color-on-surface-variant)]"
                          }`}
                        >
                          {dist > 0 ? `+${dist}%` : `${dist}%`}
                        </span>
                      </td>

                      <td className="py-2.5 px-3 text-center">
                        <span className={`text-[10px] font-mono px-2 py-0.5 rounded-full border font-bold ${getBandarSignalBadge(row.signal)}`}>
                          {row.signal}
                        </span>
                      </td>

                      <td className="py-2.5 px-3 text-right font-mono font-bold">
                        <span
                          className={`px-2 py-0.5 rounded border text-xs ${
                            row.conviction_score >= 70
                              ? "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/30"
                              : row.conviction_score >= 50
                              ? "bg-cyan-500/15 text-cyan-700 dark:text-cyan-400 border-cyan-500/30"
                              : row.conviction_score >= 40
                              ? "bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/30"
                              : "bg-rose-500/15 text-rose-700 dark:text-rose-400 border-rose-500/30"
                          }`}
                        >
                          {row.conviction_score}
                        </span>
                      </td>

                      <td className="py-2.5 px-3 text-center font-mono text-[11px]">
                        <span className="text-emerald-700 dark:text-emerald-400 font-bold" title="Top Buyer">
                          {row.top_buyer}
                        </span>
                        <span className="text-[var(--md-sys-color-on-surface-variant)] mx-1">/</span>
                        <span className="text-rose-700 dark:text-rose-400 font-bold" title="Top Seller">
                          {row.top_seller}
                        </span>
                      </td>

                      <td className="py-2.5 px-3 text-right font-mono text-[11px] font-semibold" style={{ color: signedColor(row.smart_foreign_net) }}>
                        {fmtRp(row.smart_foreign_net)}
                      </td>

                      <td className="py-2.5 px-3 text-right font-mono text-[11px] font-semibold" style={{ color: signedColor(row.gorengan_net) }}>
                        {fmtRp(row.gorengan_net)}
                      </td>

                      <td className="py-2.5 px-3 text-center">
                        <span className={`text-[9px] font-mono px-1.5 py-0.5 rounded border font-semibold ${getHmmBadge(row.foreign_hmm_state)}`}>
                          {row.foreign_hmm_state}
                        </span>
                      </td>

                      <td className="py-2.5 px-3 text-right font-mono text-[11px]">
                        <span
                          className={
                            row.free_float_est > 0 && row.free_float_est <= 25
                              ? "text-amber-700 dark:text-amber-400 font-bold"
                              : "text-[var(--md-sys-color-on-surface-variant)]"
                          }
                        >
                          {row.free_float_est > 0 ? `${row.free_float_est}%` : "-"}
                        </span>
                      </td>

                      <td className="py-2.5 px-3">
                        <div className="flex flex-wrap gap-1">
                          {row.quick_tags && row.quick_tags.length > 0 ? (
                            row.quick_tags.map((tag) => {
                              const isGolden = tag === "GOLDEN_CONFLUENCE";
                              const isStealth = tag === "STEALTH_ACCUMULATION";
                              const isGorengan = tag === "GORENGAN_MOMENTUM";
                              const isTrap = tag === "DISTRIBUTION_TRAP";
                              return (
                                <span
                                  key={tag}
                                  className={`text-[9px] font-mono px-1.5 py-0.5 rounded border font-bold ${
                                    isGolden
                                      ? "bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/30"
                                      : isStealth
                                      ? "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/30"
                                      : isGorengan
                                      ? "bg-purple-500/15 text-purple-700 dark:text-purple-400 border-purple-500/30"
                                      : isTrap
                                      ? "bg-rose-500/15 text-rose-700 dark:text-rose-400 border-rose-500/30"
                                      : "bg-sky-500/15 text-sky-700 dark:text-sky-400 border-sky-500/30"
                                  }`}
                                >
                                  {tag.replace("_", " ")}
                                </span>
                              );
                            })
                          ) : (
                            <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] opacity-40">-</span>
                          )}
                        </div>
                      </td>

                      <td className="py-2.5 px-3 text-center">
                        <div className="inline-flex items-center gap-1.5">
                          <Link
                            to={`/${row.ticker}`}
                            title="Buka Detail Bandarmology"
                            className="p-1 rounded-md bg-[var(--md-sys-color-surface-container-highest)] hover:bg-[var(--md-sys-color-primary)] hover:text-[var(--md-sys-color-on-primary)] transition-all text-[var(--md-sys-color-on-surface-variant)]"
                          >
                            <Icon icon="ph:chart-line-up-bold" width="13" />
                          </Link>
                          <Link
                            to={`/foreign?ticker=${row.ticker}`}
                            title="Buka Foreign Flow"
                            className="p-1 rounded-md bg-[var(--md-sys-color-surface-container-highest)] hover:bg-[var(--md-sys-color-primary)] hover:text-[var(--md-sys-color-on-primary)] transition-all text-[var(--md-sys-color-on-surface-variant)]"
                          >
                            <Icon icon="ph:globe-stand-bold" width="13" />
                          </Link>
                        </div>
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
