import { Icon } from "@iconify/react";
import type { MmdDetailResponse } from "./mmdTypes";

interface MmdParameterDrawerProps {
  open: boolean;
  onClose: () => void;
  candleData: MmdDetailResponse | null;
  dateFrom: string;
  setDateFrom: (val: string) => void;
  dateTo: string;
  setDateTo: (val: string) => void;
  showT5Markers: boolean;
  setShowT5Markers: (val: boolean) => void;
  showTooltip: boolean;
  setShowTooltip: (val: boolean) => void;
  quickModeActive: boolean;
  setQuickModeActive: (val: boolean) => void;
  retailBandarAutoEnabled: boolean;
  setRetailBandarAutoEnabled: (val: boolean) => void;
}

export default function MmdParameterDrawer({
  open,
  onClose,
  candleData,
  dateFrom,
  setDateFrom,
  dateTo,
  setDateTo,
  showT5Markers,
  setShowT5Markers,
  showTooltip,
  setShowTooltip,
  quickModeActive,
  setQuickModeActive,
  retailBandarAutoEnabled,
  setRetailBandarAutoEnabled,
}: MmdParameterDrawerProps) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex justify-end bg-black/60 backdrop-blur-xs">
      <div className="w-full max-w-md bg-[var(--md-sys-color-surface-container)] border-l border-[var(--md-sys-color-outline-variant)] h-full overflow-y-auto p-5 flex flex-col justify-between shadow-2xl">
        <div className="space-y-5">
          <div className="flex items-center justify-between pb-3 border-b border-[var(--md-sys-color-outline-variant)]">
            <div className="flex items-center gap-2">
              <Icon icon="ph:sliders-horizontal-duotone" className="text-[var(--md-sys-color-primary)]" width={20} />
              <span className="font-extrabold text-base text-[var(--md-sys-color-on-surface)]">
                Parameter &amp; Filter
              </span>
            </div>
            <button
              type="button"
              onClick={onClose}
              className="p-1.5 rounded-lg hover:bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] cursor-pointer"
            >
              <Icon icon="ph:x-bold" width={18} />
            </button>
          </div>

          <div className="space-y-3 font-mono-nums">
            <span className="text-xs font-bold uppercase text-[var(--md-sys-color-on-surface)] tracking-wider">
              Rentang Waktu Bursa
            </span>
            <div className="grid grid-cols-2 gap-3 text-xs">
              <div>
                <label className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] block mb-1">
                  Tanggal Mulai:
                </label>
                <select
                  value={dateFrom}
                  onChange={(e) => {
                    const newFrom = e.target.value;
                    setDateFrom(newFrom);
                    if (dateTo && newFrom > dateTo) setDateTo(newFrom);
                  }}
                  className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded-xl px-2.5 py-2 text-xs font-bold text-[var(--md-sys-color-on-surface)] outline-none"
                >
                  {(candleData?.candles || []).map((c) => (
                    <option key={c.date_iso} value={c.date_iso}>
                      {c.date_iso}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] block mb-1">
                  Tanggal Akhir:
                </label>
                <select
                  value={dateTo}
                  onChange={(e) => {
                    const newTo = e.target.value;
                    setDateTo(newTo);
                    if (dateFrom && newTo < dateFrom) setDateFrom(newTo);
                  }}
                  className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded-xl px-2.5 py-2 text-xs font-bold text-[var(--md-sys-color-on-surface)] outline-none"
                >
                  {(candleData?.candles || []).map((c) => (
                    <option key={c.date_iso} value={c.date_iso}>
                      {c.date_iso}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          </div>

          <div className="space-y-2.5 pt-2 border-t border-[var(--md-sys-color-outline-variant)]">
            <span className="text-xs font-bold uppercase text-[var(--md-sys-color-on-surface)] tracking-wider">
              Opsi Tampilan Grafik
            </span>

            <label className="flex items-center justify-between p-2.5 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]/60 cursor-pointer">
              <span className="text-xs font-semibold text-[var(--md-sys-color-on-surface)]">
                Tampilkan Marker Big Player (T5)
              </span>
              <input
                type="checkbox"
                checked={showT5Markers}
                onChange={(e) => setShowT5Markers(e.target.checked)}
                className="accent-[var(--md-sys-color-primary)] w-4 h-4 cursor-pointer"
              />
            </label>

            <label className="flex items-center justify-between p-2.5 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]/60 cursor-pointer">
              <span className="text-xs font-semibold text-[var(--md-sys-color-on-surface)]">
                Aktifkan Tooltip Interaktif
              </span>
              <input
                type="checkbox"
                checked={showTooltip}
                onChange={(e) => setShowTooltip(e.target.checked)}
                className="accent-[var(--md-sys-color-primary)] w-4 h-4 cursor-pointer"
              />
            </label>

            <label className="flex items-center justify-between p-2.5 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]/60 cursor-pointer">
              <span className="text-xs font-semibold text-[var(--md-sys-color-on-surface)]">
                Quick Mode (Multi-Periode)
              </span>
              <input
                type="checkbox"
                checked={quickModeActive}
                onChange={(e) => setQuickModeActive(e.target.checked)}
                className="accent-[var(--md-sys-color-primary)] w-4 h-4 cursor-pointer"
              />
            </label>
          </div>

          <div className="space-y-2 pt-2 border-t border-[var(--md-sys-color-outline-variant)]">
            <span className="text-xs font-bold uppercase text-[var(--md-sys-color-on-surface)] tracking-wider">
              Deteksi Broker Retail
            </span>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => setRetailBandarAutoEnabled(true)}
                className={`flex-1 py-2 rounded-xl text-xs font-bold border transition cursor-pointer ${
                  retailBandarAutoEnabled
                    ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border-[var(--md-sys-color-primary)]"
                    : "bg-[var(--md-sys-color-surface-container-low)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]"
                }`}
              >
                Auto Dynamic (ML)
              </button>
              <button
                type="button"
                onClick={() => setRetailBandarAutoEnabled(false)}
                className={`flex-1 py-2 rounded-xl text-xs font-bold border transition cursor-pointer ${
                  !retailBandarAutoEnabled
                    ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border-[var(--md-sys-color-primary)]"
                    : "bg-[var(--md-sys-color-surface-container-low)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]"
                }`}
              >
                Manual Override
              </button>
            </div>
          </div>
        </div>

        <div className="pt-4 border-t border-[var(--md-sys-color-outline-variant)]">
          <button
            type="button"
            onClick={onClose}
            className="w-full py-2.5 rounded-xl bg-[var(--md-sys-color-primary)] hover:opacity-90 text-[var(--md-sys-color-on-primary)] font-bold text-xs shadow-xs cursor-pointer transition"
          >
            Tutup &amp; Terapkan
          </button>
        </div>
      </div>
    </div>
  );
}
