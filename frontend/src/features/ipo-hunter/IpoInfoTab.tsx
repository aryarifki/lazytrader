import React, { useEffect, useMemo, useState } from "react";
import { getStockOwnershipProfile } from "@/lib/api/terminal";
import type { StockOwnership } from "@/features/ownership/stock/stockTypes";
import type { IpoInfoCard } from "./ipoTypes";
import { fmtInt, fmtSignedPct0 } from "./ipoFormat";

interface NormalizedHolder {
  name: string;
  pct: number | null;
}

export default function IpoInfoTab({ infoCard, ticker }: { infoCard?: IpoInfoCard | null; ticker: string }) {
  const [filterGroup, setFilterGroup] = useState<"all" | "off" | "time" | "perf">("all");
  const [ownership, setOwnership] = useState<StockOwnership | null>(null);
  const [profileLoading, setProfileLoading] = useState(false);

  useEffect(() => {
    if (!ticker) return;
    let cancelled = false;
    setProfileLoading(true);
    getStockOwnershipProfile(ticker)
      .then((res) => {
        if (cancelled) return;
        setOwnership(res);
        setProfileLoading(false);
      })
      .catch(() => {
        if (cancelled) return;
        setOwnership(null);
        setProfileLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [ticker]);

  const holders = useMemo<NormalizedHolder[]>(() => {
    if (!ownership) return [];
    if (ownership.staging_holders?.holders?.length) {
      return ownership.staging_holders.holders.map((h) => ({
        name: h.name,
        pct: h.pct ?? null,
      }));
    }
    if (ownership.ct_holders?.length) {
      return ownership.ct_holders.map((h) => ({
        name: h.name,
        pct: h.pct_jun ?? h.pct_mei ?? null,
      }));
    }
    if (ownership.holders?.length) {
      return ownership.holders.map((h) => ({
        name: h.name,
        pct: h.pct ?? null,
      }));
    }
    return [];
  }, [ownership]);

  if (!infoCard) return null;

  // 3-pin Price Scale (Close, Average, IPO Price)
  const rawPts = [
    { label: "Close", val: infoCard.last_close, kind: "close", color: "var(--md-sys-color-primary)" },
    { label: "Average", val: infoCard.vwap, kind: "avg", color: "#3B82F6" },
    { label: "IPO Price", val: infoCard.ipo_price, kind: "ipo", color: "#8B5CF6" },
  ].filter((x): x is { label: string; val: number; kind: string; color: string } => x.val != null && x.val > 0);

  let scalePts: Array<{ label: string; val: number; kind: string; color: string; pos: number; lab: number }> = [];
  if (rawPts.length >= 2) {
    const lo = Math.min(...rawPts.map((x) => x.val));
    const hi = Math.max(...rawPts.map((x) => x.val));
    const rg = hi - lo || 1;
    scalePts = rawPts.map((x) => {
      const p = 5 + ((x.val - lo) / rg) * 90;
      return { ...x, pos: p, lab: p };
    });
    const sorted = [...scalePts].sort((a, b) => a.pos - b.pos);
    for (let k = 1; k < sorted.length; k++) {
      if (sorted[k].lab - sorted[k - 1].lab < 18) sorted[k].lab = sorted[k - 1].lab + 18;
    }
    const over = sorted[sorted.length - 1].lab - 95;
    if (over > 0) sorted.forEach((x) => (x.lab -= over));
  }

  const dimCls = (g: "off" | "time" | "perf") =>
    filterGroup !== "all" && filterGroup !== g ? "opacity-35 transition-opacity" : "transition-opacity";

  const totalPct = ownership?.summary?.total_pct;
  const freeFloatPct = ownership?.profile?.free_float_pct ?? ownership?.summary?.free_float_est;

  return (
    <div className="space-y-4">
      {/* INFO IPO CARD */}
      <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-4 shadow-sm space-y-3">
        <div className="flex items-center justify-between">
          <h3 className="text-xs font-extrabold uppercase tracking-wider">INFO IPO</h3>
          <span className="font-mono-nums font-extrabold text-xs text-[var(--md-sys-color-primary)]">
            {infoCard.ticker}
          </span>
        </div>

        <div>
          <div className="text-sm font-bold text-[var(--md-sys-color-on-surface)]">{infoCard.company_name}</div>
          <div className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
            IPO {infoCard.listing_date || "—"}
            {infoCard.trading_days_since_ipo != null && (
              <>
                {" "}
                · <b className="text-[var(--md-sys-color-on-surface)]">{fmtInt(infoCard.trading_days_since_ipo)}</b> hari
                bursa
              </>
            )}
          </div>
        </div>

        {/* 3-pin Price Scale */}
        {scalePts.length >= 2 && (
          <div className={`pt-2 pb-1 ${dimCls("perf")}`}>
            <div className="relative h-14 mx-1">
              <div className="absolute left-0 right-0 top-3 h-1 rounded-full bg-[var(--md-sys-color-outline-variant)]/60" />
              {scalePts.map((pt) => (
                <React.Fragment key={pt.kind}>
                  <div
                    className="absolute top-1.5 -translate-x-1/2 w-3.5 h-3.5 rounded-full border-2 border-white shadow"
                    style={{ left: `${pt.pos}%`, backgroundColor: pt.color }}
                    title={`${pt.label}: ${fmtInt(pt.val)}`}
                  />
                  <div
                    className="absolute top-6 -translate-x-1/2 text-center leading-tight whitespace-nowrap"
                    style={{ left: `${pt.lab}%` }}
                  >
                    <div className="text-[9.5px] font-semibold text-[var(--md-sys-color-on-surface-variant)]">
                      {pt.label}
                    </div>
                    <div className="text-[11px] font-mono-nums font-extrabold">{fmtInt(pt.val)}</div>
                  </div>
                </React.Fragment>
              ))}
            </div>
          </div>
        )}

        {/* Performance KPIs */}
        <div className={`grid grid-cols-2 gap-2.5 ${dimCls("perf")}`}>
          <div className="rounded-xl bg-[var(--md-sys-color-surface-container-high)]/60 border border-[var(--md-sys-color-outline-variant)]/50 p-2.5">
            <div className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">
              IPO vs Close Price
            </div>
            <div
              className={`text-lg font-mono-nums font-black ${
                (infoCard.ipo_vs_close_pct ?? 0) > 0
                  ? "text-emerald-600 dark:text-emerald-400"
                  : (infoCard.ipo_vs_close_pct ?? 0) < 0
                  ? "text-rose-600 dark:text-rose-400"
                  : ""
              }`}
            >
              {fmtSignedPct0(infoCard.ipo_vs_close_pct)}
            </div>
            <div className="text-[10.5px] font-mono-nums text-[var(--md-sys-color-on-surface-variant)]">
              {fmtInt(infoCard.ipo_price)} → {fmtInt(infoCard.last_close)}
            </div>
          </div>

          <div className="rounded-xl bg-[var(--md-sys-color-surface-container-high)]/60 border border-[var(--md-sys-color-outline-variant)]/50 p-2.5">
            <div className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">
              Average vs Close Price
            </div>
            <div
              className={`text-lg font-mono-nums font-black ${
                (infoCard.avg_vs_close_pct ?? 0) > 0
                  ? "text-emerald-600 dark:text-emerald-400"
                  : (infoCard.avg_vs_close_pct ?? 0) < 0
                  ? "text-rose-600 dark:text-rose-400"
                  : ""
              }`}
            >
              {fmtSignedPct0(infoCard.avg_vs_close_pct)}
            </div>
            <div className="text-[10.5px] font-mono-nums text-[var(--md-sys-color-on-surface-variant)]">
              {fmtInt(infoCard.vwap)} → {fmtInt(infoCard.last_close)}
            </div>
          </div>
        </div>

        {/* Filter chips */}
        <div className="flex flex-wrap gap-1.5 pt-1">
          {(
            [
              ["all", "Semua"],
              ["off", "Penawaran"],
              ["time", "Waktu"],
              ["perf", "Performa"],
            ] as const
          ).map(([k, lbl]) => (
            <button
              key={k}
              type="button"
              onClick={() => setFilterGroup(k)}
              className={`px-2.5 py-0.5 rounded-full text-[10.5px] font-bold border transition cursor-pointer ${
                filterGroup === k
                  ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] border-[var(--md-sys-color-primary)]"
                  : "bg-[var(--md-sys-color-surface-container-high)] border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface-variant)]"
              }`}
            >
              {lbl}
            </button>
          ))}
        </div>

        {/* 8 Metadata Tiles */}
        <div className="grid grid-cols-2 gap-2 text-xs">
          {[
            { g: "off" as const, k: "Jumlah Lot IPO", v: fmtInt(infoCard.lot_offer) },
            { g: "off" as const, k: "% Dijual saat IPO", v: infoCard.pct_ipo || "—" },
            { g: "off" as const, k: "Nilai IPO (Miliar)", v: fmtInt(infoCard.nilai_miliar) },
            {
              g: "off" as const,
              k: "Est. Penjatahan Retail",
              v: infoCard.penjatahan != null ? `${infoCard.penjatahan}%` : "NA",
              accent: "text-rose-600 dark:text-rose-400",
            },
            { g: "time" as const, k: "Tanggal IPO", v: infoCard.listing_date || "—" },
            {
              g: "time" as const,
              k: "Hari Bursa Pasca IPO",
              v: infoCard.trading_days_since_ipo != null ? `${fmtInt(infoCard.trading_days_since_ipo)} hari` : "—",
            },
            { g: "off" as const, k: "Total Saham Scriptless", v: fmtInt(infoCard.lot_total) },
            { g: "off" as const, k: "Papan", v: infoCard.papan || "—" },
          ].map((tile) => (
            <div
              key={tile.k}
              className={`rounded-lg bg-[var(--md-sys-color-surface-container-high)]/40 border border-[var(--md-sys-color-outline-variant)]/40 p-2 ${dimCls(
                tile.g
              )}`}
            >
              <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">{tile.k}</div>
              <div className={`font-mono-nums font-bold text-xs mt-0.5 ${tile.accent || ""}`}>{tile.v}</div>
            </div>
          ))}
        </div>

        {/* Warnings */}
        {infoCard.gap && (
          <div className="p-2 rounded-lg bg-amber-500/10 border border-amber-500/30 text-amber-700 dark:text-amber-300 text-[11px]">
            ⚠ Ada jeda data: {infoCard.gap.from} – {infoCard.gap.to} tidak terhitung
          </div>
        )}
        {infoCard.no_d0 && (
          <div className="p-2 rounded-lg bg-amber-500/10 border border-amber-500/30 text-amber-700 dark:text-amber-300 text-[11px]">
            ⚠ Data D0 tidak tersedia — dihitung dari data harian pertama yang ada.
          </div>
        )}
      </div>

      {/* PROFILE SAHAM CARD (KSEI >1%) */}
      <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-4 shadow-sm">
        <div className="flex items-center justify-between mb-2">
          <h3 className="text-xs font-extrabold uppercase tracking-wider">PROFILE SAHAM</h3>
          <span className="text-[10.5px] font-mono-nums text-[var(--md-sys-color-on-surface-variant)]">KSEI &gt;1%</span>
        </div>

        {profileLoading ? (
          <div className="py-4 text-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
            Memuat profil kepemilikan…
          </div>
        ) : holders.length === 0 ? (
          <div className="py-4 text-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
            Data kepemilikan belum tersedia untuk <b>{ticker}</b> (belum masuk rilis KSEI &gt;1%).
          </div>
        ) : (
          <div className="space-y-3 text-xs">
            <div className="grid grid-cols-3 gap-2">
              <div className="rounded-lg bg-[var(--md-sys-color-surface-container-high)]/50 p-2">
                <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">Pemegang &gt;1%</div>
                <div className="font-mono-nums font-bold">{holders.length}</div>
              </div>
              <div className="rounded-lg bg-[var(--md-sys-color-surface-container-high)]/50 p-2">
                <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">Total &gt;1%</div>
                <div className="font-mono-nums font-bold">
                  {totalPct != null ? `${Number(totalPct).toFixed(1)}%` : "—"}
                </div>
              </div>
              <div className="rounded-lg bg-[var(--md-sys-color-surface-container-high)]/50 p-2">
                <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">Free Float</div>
                <div className="font-mono-nums font-bold">
                  {freeFloatPct != null ? `${Number(freeFloatPct).toFixed(1)}%` : "—"}
                </div>
              </div>
            </div>

            <div className="max-h-56 overflow-y-auto border border-[var(--md-sys-color-outline-variant)]/40 rounded-lg">
              <table className="w-full text-[11px]">
                <thead className="bg-[var(--md-sys-color-surface-container-high)]/60 sticky top-0">
                  <tr>
                    <th className="px-2.5 py-1.5 text-left font-bold">Pemegang Saham</th>
                    <th className="px-2.5 py-1.5 text-right font-bold">Porsi</th>
                  </tr>
                </thead>
                <tbody>
                  {holders.slice(0, 12).map((sh, idx) => (
                    <tr key={`${sh.name}-${idx}`} className="border-t border-[var(--md-sys-color-outline-variant)]/30">
                      <td className="px-2.5 py-1.5 truncate max-w-[200px]" title={sh.name}>
                        {sh.name}
                      </td>
                      <td className="px-2.5 py-1.5 text-right font-mono-nums font-semibold">
                        {sh.pct != null ? `${Number(sh.pct).toFixed(2)}%` : "—"}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
