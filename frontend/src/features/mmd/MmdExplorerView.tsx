import { Icon } from "@iconify/react";
import { bandarBadgeClass } from "./mmdBadges";
import type { ExplorerCategory, MmdExplorerRow } from "./mmdTypes";

interface MmdExplorerViewProps {
  explorerCat: ExplorerCategory;
  setExplorerCat: (cat: ExplorerCategory) => void;
  explorerCriteria: string;
  setExplorerCriteria: (crit: string) => void;
  explorerVolumeBucket: number;
  setExplorerVolumeBucket: (vol: number) => void;
  explorerResults: MmdExplorerRow[];
  explorerLoading: boolean;
  onScan: (cat: ExplorerCategory, crit: string, vol: number) => void;
  onSelectTicker: (ticker: string) => void;
}

const CATEGORIES: Array<{ id: ExplorerCategory; label: string }> = [
  { id: "normal", label: "NORMAL" },
  { id: "advanced", label: "ADVANCED" },
  { id: "bigPlayer", label: "BIG PLAYER" },
];

export default function MmdExplorerView({
  explorerCat,
  setExplorerCat,
  explorerCriteria,
  setExplorerCriteria,
  explorerVolumeBucket,
  setExplorerVolumeBucket,
  explorerResults,
  explorerLoading,
  onScan,
  onSelectTicker,
}: MmdExplorerViewProps) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
      <div className="lg:col-span-4 space-y-3">
        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 shadow-2xs space-y-3.5">
          <div className="text-xs font-bold text-[var(--md-sys-color-on-surface)] pb-2 border-b border-[var(--md-sys-color-outline-variant)] flex items-center justify-between">
            <span className="flex items-center gap-1.5 uppercase tracking-wider">
              <Icon icon="ph:magnifying-glass-duotone" width={16} className="text-[var(--md-sys-color-primary)]" />
              PARAMETER SCANNER
            </span>
            <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
              Real-Time
            </span>
          </div>

          <div className="space-y-1.5">
            <label className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider block">
              Kategori:
            </label>
            <div className="grid grid-cols-3 gap-1 bg-[var(--md-sys-color-surface-container-low)] p-1 rounded-xl border border-[var(--md-sys-color-outline-variant)]/60">
              {CATEGORIES.map((c) => (
                <button
                  key={c.id}
                  type="button"
                  onClick={() => {
                    setExplorerCat(c.id);
                    setExplorerCriteria("ALL");
                  }}
                  className={`py-1.5 text-xs font-bold rounded-lg cursor-pointer transition ${
                    explorerCat === c.id
                      ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-2xs"
                      : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                  }`}
                >
                  {c.label}
                </button>
              ))}
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider block">
              Kondisi / Sinyal:
            </label>
            <select
              value={explorerCriteria}
              onChange={(e) => setExplorerCriteria(e.target.value)}
              className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/70 rounded-xl px-3 py-2 text-xs text-[var(--md-sys-color-on-surface)] font-medium outline-none focus:border-[var(--md-sys-color-primary)] cursor-pointer"
            >
              {explorerCat === "normal" && (
                <>
                  <option value="ALL">Semua Hasil Bandar</option>
                  <option value="Big Acc">Big Accumulation (Akumulasi Besar)</option>
                  <option value="Normal Acc">Normal Accumulation</option>
                  <option value="Small Acc">Small Accumulation</option>
                  <option value="Netral">Netral</option>
                  <option value="Distribusi">Distribusi (Semua Tingkat)</option>
                </>
              )}
              {explorerCat === "advanced" && (
                <>
                  <option value="ALL">Semua Kondisi Advanced (18 Sinyal)</option>
                  <option value="Mark Up">1. Mark Up Bandar (Harga Naik Didukung Bandar)</option>
                  <option value="Mark Down">2. Mark Down Bandar (Harga Turun Tanpa Distribusi)</option>
                  <option value="Sideways Acc">3. Sideways Akumulasi Tersembunyi</option>
                  <option value="Sideways Dist">4. Sideways Distribusi Halus</option>
                  <option value="High Volatility">5. Volatilitas Tinggi dengan Partisipasi Bandar</option>
                  <option value="Retail Panic">6. Retail Panic Selling (Bandar Tampung)</option>
                  <option value="Strong Corr">7. Korelasi Kuat Harga vs Aliran Bandar</option>
                  <option value="Acc">8. Akumulasi Masif 5 Hari Terakhir</option>
                  <option value="Dist">9. Distribusi Masif 5 Hari Terakhir</option>
                </>
              )}
              {explorerCat === "bigPlayer" && (
                <>
                  <option value="ALL">Semua Aktivitas Big Player</option>
                  <option value="Active">Aktivitas Transaksi T5 Terkini</option>
                  <option value="Mayoritas">Akumulasi Pemegang Saham Mayoritas</option>
                  <option value="Masif">Transaksi Bernilai Masif (&gt; 10 Miliar)</option>
                </>
              )}
            </select>
          </div>

          <div className="space-y-1.5">
            <label className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider block">
              Nilai Transaksi Harian (Bucket):
            </label>
            <select
              value={explorerVolumeBucket}
              onChange={(e) => setExplorerVolumeBucket(Number(e.target.value))}
              className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/70 rounded-xl px-3 py-2 text-xs text-[var(--md-sys-color-on-surface)] font-medium outline-none focus:border-[var(--md-sys-color-primary)] cursor-pointer"
            >
              <option value={0}>Semua Nilai Transaksi</option>
              <option value={1000000000}>&gt; Rp 1 Miliar</option>
              <option value={5000000000}>&gt; Rp 5 Miliar</option>
              <option value={10000000000}>&gt; Rp 10 Miliar</option>
              <option value={50000000000}>&gt; Rp 50 Miliar</option>
              <option value={100000000000}>&gt; Rp 100 Miliar</option>
            </select>
          </div>

          <button
            type="button"
            onClick={() => onScan(explorerCat, explorerCriteria, explorerVolumeBucket)}
            disabled={explorerLoading}
            className="w-full py-2.5 px-4 rounded-xl bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] font-bold text-xs flex items-center justify-center gap-2 hover:opacity-90 active:scale-[0.98] transition shadow-xs cursor-pointer disabled:opacity-50"
          >
            {explorerLoading ? (
              <>
                <div className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
                <span>Memindai Pasar...</span>
              </>
            ) : (
              <>
                <Icon icon="ph:play-bold" width={14} />
                <span>Scan Pasar Sekarang</span>
              </>
            )}
          </button>
        </div>
      </div>

      <div className="lg:col-span-8 space-y-3">
        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl overflow-hidden shadow-2xs">
          <div className="p-3 bg-[var(--md-sys-color-surface-container-high)] border-b border-[var(--md-sys-color-outline-variant)] flex items-center justify-between">
            <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)]">
              HASIL SCAN EXPLORER
            </span>
            <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-mono">
              {explorerResults.length} Saham
            </span>
          </div>

          {explorerLoading ? (
            <div className="p-12 text-center text-xs text-[var(--md-sys-color-on-surface-variant)] flex flex-col items-center gap-2">
              <div className="w-6 h-6 border-2 border-[var(--md-sys-color-primary)] border-t-transparent rounded-full animate-spin" />
              <span>Memindai data saham...</span>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs font-mono">
                <thead>
                  <tr className="bg-[var(--md-sys-color-surface-container-high)] text-[10px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]">
                    <th className="py-2 px-3">KODE SAHAM</th>
                    <th className="py-2 px-3">HARGA</th>
                    <th className="py-2 px-3">DETAIL TRANSAKSI</th>
                    <th className="py-2 px-3 text-center">HASIL BANDAR</th>
                    <th className="py-2 px-3 text-center">AKSI</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40 font-mono-nums">
                  {explorerResults.length > 0 ? (
                    explorerResults.map((r, idx) => (
                      <tr key={idx} className="hover:bg-[var(--md-sys-color-surface-container-high)] transition">
                        <td className="py-2 px-3 font-bold text-[var(--md-sys-color-on-surface)]">{r.ticker}</td>
                        <td className="py-2 px-3 font-bold text-teal-700 dark:text-[#2FD4C6]">
                          Rp {Number(r.last_close || r.last_price || 0).toLocaleString("id-ID")}
                        </td>
                        <td className="py-2 px-3 text-[var(--md-sys-color-on-surface-variant)] text-[11px]">
                          {r.detail || `${r.change_pct || r.price_change || 0}% change`}
                        </td>
                        <td className="py-2 px-3 text-center">
                          <span className={`${bandarBadgeClass(r.badge?.cls)} text-[10px] py-0.5 px-2`}>
                            {r.badge?.label || "—"}
                          </span>
                        </td>
                        <td className="py-2 px-3 text-center">
                          <button
                            type="button"
                            onClick={() => onSelectTicker(r.ticker)}
                            className="px-2.5 py-1 rounded-lg bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] text-[10px] font-bold hover:opacity-90 transition cursor-pointer"
                          >
                            Lihat MMD
                          </button>
                        </td>
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <td colSpan={5} className="py-8 text-center text-[var(--md-sys-color-on-surface-variant)]">
                        Tidak ada hasil scan.
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
