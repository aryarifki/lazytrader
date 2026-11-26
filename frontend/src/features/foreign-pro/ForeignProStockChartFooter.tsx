import type { HistTopMode } from "./foreignProTypes";

interface ForeignProStockChartFooterProps {
  latestPartisipasi: number;
  latestKat: number;
  latestKfh: number;
  latestKret: number;
  latestCorr: number;
  latestCorrKF: number;
  latestValM: number;
  additionalOverlay: string;
  histTopMode: HistTopMode;
  setHistTopMode: (v: HistTopMode) => void;
  histBottomMode: string;
  setHistBottomMode: (v: string) => void;
  showPBNyangkut: boolean;
  setShowPBNyangkut: (v: boolean) => void;
  zoomDays: number;
  setZoomDays: (v: number) => void;
  windowOffset: number;
  totalStockPoints: number;
  notifOpen: boolean;
  setNotifOpen: (v: boolean) => void;
  highlights: string[];
}

export default function ForeignProStockChartFooter({
  latestPartisipasi,
  latestKat,
  latestKfh,
  latestKret,
  latestCorr,
  latestCorrKF,
  latestValM,
  additionalOverlay,
  histTopMode,
  setHistTopMode,
  histBottomMode,
  setHistBottomMode,
  showPBNyangkut,
  setShowPBNyangkut,
  zoomDays,
  setZoomDays,
  windowOffset,
  totalStockPoints,
  notifOpen,
  setNotifOpen,
  highlights,
}: ForeignProStockChartFooterProps) {
  return (
    <>
      <div
        id="histInfoBar"
        className="my-2 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/60 rounded px-3 py-1.5 flex items-center justify-between text-xs font-mono text-[var(--md-sys-color-on-surface-variant)] overflow-x-auto gap-4 whitespace-nowrap"
      >
        <span>
          Partisipasi Asing <b className="text-[var(--md-sys-color-on-surface)]">{latestPartisipasi.toFixed(1)}%</b>
        </span>
        <span className="text-[var(--md-sys-color-outline-variant)]">|</span>
        <span>
          Kepemilikan Foreign Total <b className="text-[var(--md-sys-color-on-surface)]">{latestKat.toFixed(1)}%</b>
        </span>
        <span className="text-[var(--md-sys-color-outline-variant)]">|</span>
        <span>
          Kepemilikan Foreign Hidden <b className="text-[var(--md-sys-color-on-surface)]">{latestKfh.toFixed(1)}%</b>
        </span>
        <span className="text-[var(--md-sys-color-outline-variant)]">|</span>
        <span>
          Kepemilikan Ritel <b className="text-[var(--color-negative)]">{latestKret.toFixed(1)}%</b>
        </span>
        <span className="text-[var(--md-sys-color-outline-variant)]">|</span>
        <span>
          Korelasi FF dengan Harga{" "}
          <b className={latestCorr >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"}>
            ({latestCorr >= 0 ? `+${latestCorr}` : latestCorr})
          </b>
        </span>
        {additionalOverlay === "korelasiKF" && (
          <>
            <span className="text-[var(--md-sys-color-outline-variant)]">|</span>
            <span>
              Korelasi Korban Flow{" "}
              <b className="text-[#2F7676]">({latestCorrKF >= 0 ? `+${latestCorrKF}` : latestCorrKF})</b>
            </span>
          </>
        )}
        <span className="text-[var(--md-sys-color-outline-variant)]">|</span>
        <span id="histInfoValueStandalone">
          Value <b className="text-amber-600 dark:text-[#FBBF24]">({latestValM.toFixed(1)} M)</b>
        </span>
      </div>

      <div className="flex items-center justify-between pt-1 border-t border-[var(--md-sys-color-outline-variant)]/40 text-xs font-mono flex-wrap gap-2">
        <div className="flex items-center gap-2 flex-wrap">
          <span id="ffHistLabel" className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
            Histogram Atas:
          </span>
          <button
            id="btnSwitchTop"
            type="button"
            onClick={() => setHistTopMode(histTopMode === "ff" ? "fco" : "ff")}
            className="px-2 py-0.5 rounded bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/50 text-[var(--md-sys-color-primary)] font-bold cursor-pointer"
          >
            {histTopMode === "ff" ? "Daily FF" : "Daily F&Co"}
          </button>

          <span id="valueHistLabel" className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] ml-2">
            Bawah:
          </span>
          <select
            id="selectBottom"
            value={histBottomMode}
            onChange={(e) => setHistBottomMode(e.target.value)}
            className="bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/50 rounded px-1.5 py-0.5 text-[var(--md-sys-color-on-surface)] text-[11px]"
          >
            <option value="value">Value (Miliar)</option>
            <option value="korban">Korban Flow</option>
            <option value="xl">Broker XL</option>
            <option value="semiritel">Semi Ritel</option>
          </select>
        </div>

        <div className="flex items-center gap-3">
          <label className="flex items-center gap-1.5 text-[11px] cursor-pointer">
            <input
              type="checkbox"
              id="chkPBNyangkut"
              checked={showPBNyangkut}
              onChange={(e) => setShowPBNyangkut(e.target.checked)}
              className="accent-[#D62839]"
            />
            <span
              id="lblPBNyangkut"
              className={showPBNyangkut ? "text-[#D62839] font-bold" : "text-[var(--md-sys-color-on-surface-variant)]"}
            >
              PB Nyangkut
            </span>
          </label>

          <div
            id="scrollTrack"
            className="w-32 h-2.5 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/50 rounded relative cursor-pointer"
          >
            <div
              id="scrollThumb"
              className="h-full bg-[var(--md-sys-color-primary)] rounded absolute"
              style={{
                width: `${Math.max(15, (zoomDays / Math.max(totalStockPoints, 1)) * 100)}%`,
                right: `${Math.min(85, (windowOffset / Math.max(totalStockPoints, 1)) * 100)}%`,
              }}
            />
          </div>

          <div className="flex items-center gap-1">
            <button
              id="btnOut"
              type="button"
              onClick={() => setZoomDays(Math.min(250, zoomDays + 20))}
              title="Zoom Out"
              className="px-2 py-0.5 bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)]/50 rounded text-[var(--md-sys-color-on-surface)] cursor-pointer"
            >
              −
            </button>
            <span id="zoomLabel" className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] px-1">
              {zoomDays}
            </span>
            <button
              id="btnIn"
              type="button"
              onClick={() => setZoomDays(Math.max(40, zoomDays - 20))}
              title="Zoom In"
              className="px-2 py-0.5 bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)]/50 rounded text-[var(--md-sys-color-on-surface)] cursor-pointer"
            >
              +
            </button>
          </div>
        </div>
      </div>

      <div
        id="notifBox"
        className="mt-3 font-mono text-xs border border-[var(--md-sys-color-outline-variant)]/60 rounded-lg overflow-hidden bg-[var(--md-sys-color-surface-container-low)]"
      >
        <div
          id="notifHd"
          onClick={() => setNotifOpen(!notifOpen)}
          className="flex items-center gap-2 bg-[var(--md-sys-color-surface-container-high)] px-3 py-2 cursor-pointer hover:bg-[var(--md-sys-color-surface-container-highest)] transition-colors select-none"
        >
          <span className="w-2 h-2 rounded-full bg-[var(--md-sys-color-primary)] animate-pulse" />
          <span className="font-bold text-[var(--md-sys-color-on-surface)] text-xs">
            Highlight Hari Ini
          </span>
          <span
            id="notifCount"
            className="px-2 py-0.5 bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] rounded-full text-[10px] font-bold"
          >
            {highlights.length}
          </span>
          <span className="text-[var(--md-sys-color-on-surface-variant)] text-xs ml-auto">
            {notifOpen ? "Sembunyikan ▲" : "Tampilkan ▼"}
          </span>
        </div>

        {notifOpen && (
          <div
            id="notifBody"
            className="p-3 bg-[var(--md-sys-color-surface-container-lowest)]/90 border-t border-[var(--md-sys-color-outline-variant)]/40 grid grid-cols-1 md:grid-cols-2 gap-2"
          >
            {highlights.map((h, idx) => (
              <div
                key={idx}
                className="flex items-start gap-2 text-[11px] text-[var(--md-sys-color-on-surface)] bg-[var(--md-sys-color-surface-container)]/70 p-2 rounded border border-[var(--md-sys-color-outline-variant)]/30"
              >
                <span className="text-[var(--md-sys-color-primary)] font-bold mt-0.5">•</span>
                <span className="leading-relaxed">{h}</span>
              </div>
            ))}
            {highlights.length === 0 && (
              <div className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] italic col-span-2 py-1">
                Tidak ada highlight transaksi anomali untuk hari ini.
              </div>
            )}
          </div>
        )}
      </div>
    </>
  );
}
