import { useMemo } from "react";
import { Icon } from "@iconify/react";
import type { StockOwnership } from "@/features/ownership/stock/stockTypes";
import MmdChartCard from "./MmdChartCard";
import type { MmdAnalysis, MmdDetailResponse, MmdT5Item } from "./mmdTypes";

interface MmdBigPlayerViewProps {
  selectedTicker: string;
  candleData: MmdDetailResponse | null;
  t5Data: MmdT5Item[];
  analysisData: MmdAnalysis | null;
  ownership: StockOwnership | null;
  showT5Markers: boolean;
  showTooltip: boolean;
}

interface NormalizedHolderRow {
  name: string;
  typeLabel: string;
  isTreasury: boolean;
  prevLot: number | null;
  prevPct: number | null;
  currLot: number | null;
  currPct: number | null;
  valueM: number | null;
  diffLot: number;
  diffPct: number;
  status: string;
}

function mapHolders(ownership: StockOwnership | null): NormalizedHolderRow[] {
  const closePrice = ownership?.price?.close ?? 0;

  if (ownership?.staging_holders?.holders && ownership.staging_holders.holders.length > 0) {
    return ownership.staging_holders.holders.map((h) => {
      const prevLot =
        h.previous_total_shares != null ? Math.round(h.previous_total_shares / 100) : null;
      const currLot = h.total_shares != null ? Math.round(h.total_shares / 100) : null;
      const diffLot =
        h.change_shares != null
          ? Math.round(h.change_shares / 100)
          : currLot != null && prevLot != null
            ? currLot - prevLot
            : 0;
      const valueM =
        h.total_shares != null && closePrice > 0 ? (h.total_shares * closePrice) / 1e9 : null;
      const clsLower = (h.classification || h.name || "").toLowerCase();
      return {
        name: h.name,
        typeLabel: h.investor_type || h.location || "ID",
        isTreasury: clsLower.includes("treasury"),
        prevLot,
        prevPct: h.previous_pct ?? null,
        currLot,
        currPct: h.pct ?? null,
        valueM,
        diffLot,
        diffPct: h.change_pct ?? 0,
        status: (h.status || "").toLowerCase(),
      };
    });
  }

  if (ownership?.ct_holders && ownership.ct_holders.length > 0) {
    return ownership.ct_holders.map((h) => {
      const prevLot = h.sl_mei ?? null;
      const currLot = h.sl_jun ?? null;
      const diffLot = currLot != null && prevLot != null ? currLot - prevLot : 0;
      const prevPct = h.pct_mei ?? null;
      const currPct = h.pct_jun ?? null;
      const diffPct = currPct != null && prevPct != null ? currPct - prevPct : 0;
      const valueM = currLot != null && closePrice > 0 ? (currLot * 100 * closePrice) / 1e9 : null;
      return {
        name: h.name,
        typeLabel: h.lf || "ID",
        isTreasury: h.name.toLowerCase().includes("treasury"),
        prevLot,
        prevPct,
        currLot,
        currPct,
        valueM,
        diffLot,
        diffPct,
        status: (h.status || "").toLowerCase(),
      };
    });
  }

  if (ownership?.holders && ownership.holders.length > 0) {
    return ownership.holders.map((h) => ({
      name: h.name,
      typeLabel: h.is_controlling ? "Pengendali" : h.category || "ID",
      isTreasury: (h.category || "").toLowerCase().includes("treasury"),
      prevLot: null,
      prevPct: null,
      currLot: null,
      currPct: h.pct ?? null,
      valueM: null,
      diffLot: 0,
      diffPct: 0,
      status: "",
    }));
  }

  return [];
}

export default function MmdBigPlayerView({
  selectedTicker,
  candleData,
  t5Data,
  analysisData,
  ownership,
  showT5Markers,
  showTooltip,
}: MmdBigPlayerViewProps) {
  const holders = useMemo(() => mapHolders(ownership), [ownership]);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
      <div className="lg:col-span-7 space-y-4">
        <MmdChartCard
          candleData={candleData}
          t5Data={t5Data}
          showT5Markers={showT5Markers}
          showTooltip={showTooltip}
          isBpMode
          ticker={selectedTicker}
        />

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl overflow-hidden shadow-2xs">
          <div className="bg-[var(--md-sys-color-surface-container-high)] px-3.5 py-2.5 border-b border-[var(--md-sys-color-outline-variant)] flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Icon icon="ph:users-three-duotone" width={16} className="text-[var(--md-sys-color-primary)]" />
              <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] tracking-wide">
                DAFTAR PEMEGANG SAHAM &gt; 1% ({selectedTicker})
              </span>
            </div>
            <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
              {holders.length} Pemegang Saham
            </span>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs font-mono-nums">
              <thead>
                <tr className="bg-[var(--md-sys-color-surface-container-high)] text-[10px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)] uppercase font-mono">
                  <th className="py-2 px-2.5 text-center w-8">#</th>
                  <th className="py-2 px-2.5">Nama Investor</th>
                  <th className="py-2 px-2">Tipe</th>
                  <th className="py-2 px-2 text-right">(Lot) Prev</th>
                  <th className="py-2 px-2 text-right">% Prev</th>
                  <th className="py-2 px-2 text-right">(Lot) Curr</th>
                  <th className="py-2 px-2 text-right">% Curr</th>
                  <th className="py-2 px-2 text-right">Nilai (M)</th>
                  <th className="py-2 px-2.5 text-right">Perubahan</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/30 text-[11px]">
                {holders.length > 0 ? (
                  holders.map((r, idx) => {
                    const isNew = r.status === "new" || r.status === "baru";
                    const isOut = r.status === "out" || r.status === "keluar";
                    return (
                      <tr key={idx} className="hover:bg-[var(--md-sys-color-surface-container-high)] transition">
                        <td className="py-2 px-2.5 text-center text-[var(--md-sys-color-on-surface-variant)]">{idx + 1}</td>
                        <td className="py-2 px-2.5 font-bold text-[var(--md-sys-color-on-surface)] truncate max-w-[160px]">{r.name}</td>
                        <td className="py-2 px-2">
                          {r.isTreasury ? (
                            <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-purple-500/15 text-purple-700 dark:text-purple-400 border border-purple-500/30">Treasury</span>
                          ) : (
                            <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)]">{r.typeLabel}</span>
                          )}
                        </td>
                        <td className="py-2 px-2 text-right text-[var(--md-sys-color-on-surface-variant)]">
                          {isNew || r.prevLot == null ? "—" : r.prevLot.toLocaleString("id-ID")}
                        </td>
                        <td className="py-2 px-2 text-right text-[var(--md-sys-color-on-surface-variant)]">
                          {isNew || r.prevPct == null ? "—" : `${r.prevPct.toFixed(2)}%`}
                        </td>
                        <td className="py-2 px-2 text-right font-bold text-[var(--md-sys-color-on-surface)]">
                          {isOut ? <span className="text-[var(--color-negative)] font-bold">OUT</span> : r.currLot == null ? "—" : r.currLot.toLocaleString("id-ID")}
                        </td>
                        <td className="py-2 px-2 text-right font-bold text-[var(--md-sys-color-on-surface)]">
                          {isOut || r.currPct == null ? "—" : `${r.currPct.toFixed(2)}%`}
                        </td>
                        <td className="py-2 px-2 text-right text-[var(--md-sys-color-on-surface-variant)]">
                          {isOut || r.valueM == null ? "—" : `Rp ${r.valueM.toFixed(1)} M`}
                        </td>
                        <td
                          className={`py-2 px-2.5 text-right font-bold ${
                            isNew || r.diffLot > 0 ? "text-[var(--color-positive)]" : isOut || r.diffLot < 0 ? "text-[var(--color-negative)]" : "text-[var(--md-sys-color-on-surface-variant)]"
                          }`}
                        >
                          {isNew ? (
                            <span className="px-1.5 py-0.5 rounded text-[9px] bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border border-emerald-500/30">NEW</span>
                          ) : isOut ? (
                            <span className="px-1.5 py-0.5 rounded text-[9px] bg-rose-500/15 text-rose-700 dark:text-rose-400 border border-rose-500/30">OUT</span>
                          ) : r.diffLot !== 0 ? (
                            <span>{r.diffLot > 0 ? `+${r.diffLot.toLocaleString("id-ID")}` : r.diffLot.toLocaleString("id-ID")} ({r.diffPct > 0 ? `+${r.diffPct.toFixed(2)}%` : `${r.diffPct.toFixed(2)}%`})</span>
                          ) : (
                            "—"
                          )}
                        </td>
                      </tr>
                    );
                  })
                ) : (
                  <tr>
                    <td colSpan={9} className="py-6 text-center text-[var(--md-sys-color-on-surface-variant)]">
                      Data pemegang saham &gt; 1% tidak tersedia untuk emiten ini.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div className="lg:col-span-5 space-y-4">
        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl overflow-hidden shadow-2xs">
          <div className="bg-[var(--md-sys-color-surface-container-high)] px-3.5 py-2.5 border-b border-[var(--md-sys-color-outline-variant)] flex items-center justify-between">
            <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] tracking-wide">
              TRANSAKSI T5 (&gt; 1%)
            </span>
            <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
              {t5Data.length} Transaksi
            </span>
          </div>

          <div className="max-h-[340px] overflow-y-auto divide-y divide-[var(--md-sys-color-outline-variant)]/40 text-xs font-mono-nums">
            {t5Data.length > 0 ? (
              t5Data.map((t, idx) => (
                <div key={idx} className="p-3 hover:bg-[var(--md-sys-color-surface-container-high)] transition space-y-1">
                  <div className="flex items-center justify-between">
                    <span className="font-bold text-[var(--md-sys-color-on-surface)] truncate max-w-[200px]">
                      {t.investor || "—"}
                    </span>
                    <span
                      className={`px-2 py-0.5 rounded text-[10px] font-bold ${
                        (t.action || "").toUpperCase() === "BUY"
                          ? "bg-[var(--color-positive)]/15 text-[var(--color-positive)]"
                          : "bg-[var(--color-negative)]/15 text-[var(--color-negative)]"
                      }`}
                    >
                      {t.action || "—"}
                    </span>
                  </div>
                  <div className="flex items-center justify-between text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
                    <span>
                      Tgl: {t.date} | Broker: <b>{t.broker || "—"}</b>
                    </span>
                    <span className="font-bold text-[var(--md-sys-color-on-surface)]">
                      {t.lot.toLocaleString("id-ID")} Lot
                    </span>
                  </div>
                </div>
              ))
            ) : (
              <div className="p-8 text-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
                Tidak ada transaksi T5 untuk saham {selectedTicker}.
              </div>
            )}
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 shadow-2xs space-y-2.5">
          <div className="flex items-center gap-2 pb-2 border-b border-[var(--md-sys-color-outline-variant)]/60">
            <Icon icon="ph:user-focus-duotone" width={18} className="text-[var(--md-sys-color-primary)]" />
            <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] uppercase tracking-wider">
              RANGKUMAN PERGERAKAN BIG PLAYER
            </span>
          </div>
          <div className="space-y-2 text-xs text-[var(--md-sys-color-on-surface)] leading-relaxed">
            {analysisData?.big_player_summary && analysisData.big_player_summary.length > 0 ? (
              analysisData.big_player_summary.map((text, idx) => (
                <div
                  key={idx}
                  className="flex items-start gap-2 p-2.5 rounded-xl bg-[var(--md-sys-color-surface-container-high)]/50 border border-[var(--md-sys-color-outline-variant)]/30"
                >
                  <span className="w-1.5 h-1.5 rounded-full bg-[var(--md-sys-color-primary)] mt-1.5 shrink-0" />
                  <span>{text}</span>
                </div>
              ))
            ) : (
              <div className="p-3 text-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
                Tidak ada aktivitas signifikan terdeteksi dari pemegang saham mayoritas pada periode ini.
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
