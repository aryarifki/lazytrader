import { Icon } from "@iconify/react";

interface ForeignProParameterDrawerProps {
  open: boolean;
  onClose: () => void;
  showFF: boolean;
  setShowFF: (v: boolean) => void;
  showFCo: boolean;
  setShowFCo: (v: boolean) => void;
  showKF: boolean;
  setShowKF: (v: boolean) => void;
  showXL: boolean;
  setShowXL: (v: boolean) => void;
  showSemiRitel: boolean;
  setShowSemiRitel: (v: boolean) => void;
  zoomDays: number;
  setZoomDays: (v: number) => void;
  akumDays: number;
  setAkumDays: (v: number) => void;
  showPBNyangkut: boolean;
  setShowPBNyangkut: (v: boolean) => void;
}

export default function ForeignProParameterDrawer({
  open,
  onClose,
  showFF,
  setShowFF,
  showFCo,
  setShowFCo,
  showKF,
  setShowKF,
  showXL,
  setShowXL,
  showSemiRitel,
  setShowSemiRitel,
  zoomDays,
  setZoomDays,
  akumDays,
  setAkumDays,
  showPBNyangkut,
  setShowPBNyangkut,
}: ForeignProParameterDrawerProps) {
  if (!open) return null;

  const indicatorRows = [
    {
      label: "Foreign Flow (Asing Murni)",
      checked: showFF,
      onChange: setShowFF,
      dotColor: "#FF6820",
    },
    {
      label: "Foreign & Co",
      checked: showFCo,
      onChange: setShowFCo,
      dotColor: "#D62839",
    },
    {
      label: "Korban Flow (Retail Inflow)",
      checked: showKF,
      onChange: setShowKF,
      dotColor: "#2F7676",
    },
    {
      label: "Broker XL",
      checked: showXL,
      onChange: setShowXL,
      dotColor: "#3A7BD5",
    },
    {
      label: "Semi Ritel",
      checked: showSemiRitel,
      onChange: setShowSemiRitel,
      dotColor: "#1DB87A",
    },
  ];

  const presets = [
    {
      label: "Default Asing",
      apply: () => {
        setShowFF(true);
        setShowFCo(false);
        setShowKF(false);
        setShowXL(false);
        setShowSemiRitel(false);
      },
    },
    {
      label: "Smart Money",
      apply: () => {
        setShowFF(true);
        setShowFCo(true);
        setShowKF(false);
        setShowXL(true);
        setShowSemiRitel(false);
      },
    },
    {
      label: "Retail vs Asing",
      apply: () => {
        setShowFF(true);
        setShowFCo(false);
        setShowKF(true);
        setShowXL(false);
        setShowSemiRitel(true);
      },
    },
  ];

  return (
    <div className="fixed inset-0 z-50 flex justify-end bg-black/60 backdrop-blur-xs">
      <div className="w-full max-w-md bg-[var(--md-sys-color-surface-container)] border-l border-[var(--md-sys-color-outline-variant)] h-full overflow-y-auto p-5 flex flex-col justify-between shadow-2xl">
        <div className="space-y-4">
          <div className="flex items-center justify-between pb-3 border-b border-[var(--md-sys-color-outline-variant)]">
            <div className="flex items-center gap-2">
              <Icon
                icon="ph:sliders-horizontal-duotone"
                className="text-[var(--md-sys-color-primary)]"
                width={20}
              />
              <span className="font-extrabold text-base text-[var(--md-sys-color-on-surface)]">
                Parameter &amp; Indikator
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

          {/* Quick Presets */}
          <div className="space-y-1.5">
            <span className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] tracking-wider">
              Preset Cepat
            </span>
            <div className="flex flex-wrap gap-1.5">
              {presets.map((p) => (
                <button
                  key={p.label}
                  type="button"
                  onClick={p.apply}
                  className="px-2.5 py-1 rounded-lg text-xs font-semibold bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-surface)] hover:text-[var(--md-sys-color-on-primary-container)] border border-[var(--md-sys-color-outline-variant)] transition cursor-pointer"
                >
                  {p.label}
                </button>
              ))}
            </div>
          </div>

          {/* Main Indicators Checkboxes */}
          <div className="space-y-2 pt-2 border-t border-[var(--md-sys-color-outline-variant)]">
            <span className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] tracking-wider">
              Indikator Aliran Dana
            </span>

            {indicatorRows.map((row) => (
              <label
                key={row.label}
                className="flex items-center justify-between p-2.5 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]/60 cursor-pointer"
              >
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={row.checked}
                    onChange={(e) => row.onChange(e.target.checked)}
                    className="accent-[var(--md-sys-color-primary)] w-4 h-4 cursor-pointer"
                  />
                  <span className="text-xs font-semibold text-[var(--md-sys-color-on-surface)]">
                    {row.label}
                  </span>
                </div>
                <span
                  className="w-2.5 h-2.5 rounded-full"
                  style={{ backgroundColor: row.dotColor }}
                />
              </label>
            ))}
          </div>

          {/* Moving Average & Zoom */}
          <div className="space-y-2.5 pt-2 border-t border-[var(--md-sys-color-outline-variant)] text-xs">
            <span className="text-[10px] font-bold uppercase text-[var(--md-sys-color-on-surface-variant)] tracking-wider">
              Rentang Waktu &amp; Moving Average
            </span>

            <div className="flex items-center justify-between">
              <span className="text-[var(--md-sys-color-on-surface)] font-medium">
                Timeline Zoom:
              </span>
              <div className="flex items-center gap-1 font-mono-nums">
                {[50, 100, 200, 300].map((d) => (
                  <button
                    key={d}
                    type="button"
                    onClick={() => setZoomDays(d)}
                    className={`px-2 py-0.5 rounded text-xs font-bold border transition cursor-pointer ${
                      zoomDays === d
                        ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] border-[var(--md-sys-color-primary)]"
                        : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]"
                    }`}
                  >
                    {d}D
                  </button>
                ))}
              </div>
            </div>

            <div className="flex items-center justify-between pt-1">
              <span className="text-[var(--md-sys-color-on-surface)] font-medium">
                Rata-rata Akumulasi:
              </span>
              <div className="flex items-center gap-1">
                <input
                  type="number"
                  min={1}
                  max={100}
                  value={akumDays}
                  onChange={(e) => setAkumDays(Math.max(1, parseInt(e.target.value, 10) || 1))}
                  className="w-14 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded px-1.5 py-0.5 text-center text-xs font-bold text-[var(--md-sys-color-on-surface)]"
                />
                <span className="text-[var(--md-sys-color-on-surface-variant)] text-[11px]">
                  hari
                </span>
              </div>
            </div>

            <label className="flex items-center justify-between p-2 rounded-xl bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)]/60 cursor-pointer mt-1">
              <span className="text-xs font-semibold text-[var(--md-sys-color-on-surface)]">
                Candle PB Nyangkut
              </span>
              <input
                type="checkbox"
                checked={showPBNyangkut}
                onChange={(e) => setShowPBNyangkut(e.target.checked)}
                className="accent-[var(--md-sys-color-primary)] w-4 h-4 cursor-pointer"
              />
            </label>
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
