import { Icon } from "@iconify/react";
import type { FFMainTab, MasterTickerItem } from "./foreignProTypes";

interface ForeignProHeaderProps {
  activeMainTab: FFMainTab;
  setActiveMainTab: (tab: FFMainTab) => void;
  tickerSearch: string;
  setTickerSearch: (val: string) => void;
  suggestOpen: boolean;
  setSuggestOpen: (open: boolean) => void;
  suggestions: MasterTickerItem[];
  onSelectTicker: (tk: string) => void;
  onOpenParameters: () => void;
  onRefresh: () => void;
}

const MAIN_TABS: Array<{ id: FFMainTab; icon: string; label: string; domId?: string }> = [
  { id: "stock", icon: "ph:magnifying-glass-duotone", label: "Deep Dive Saham" },
  { id: "ihsg", icon: "ph:globe-stand-duotone", label: "Foreign Flow IHSG (IHSG Macro)", domId: "btnIhsgMacro" },
  { id: "screener", icon: "ph:funnel-duotone", label: "Stock Screener", domId: "btnStockScreener" },
];

export default function ForeignProHeader({
  activeMainTab,
  setActiveMainTab,
  tickerSearch,
  setTickerSearch,
  suggestOpen,
  setSuggestOpen,
  suggestions,
  onSelectTicker,
  onOpenParameters,
  onRefresh,
}: ForeignProHeaderProps) {
  return (
    <>
      {/* TOP RESEARCH NAV & PAGE BANNER */}
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-3">
        <div className="flex items-center gap-3.5">
          <div className="w-11 h-11 sm:w-12 sm:h-12 rounded-2xl bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0 border border-[var(--md-sys-color-outline-variant)] shadow-2xs">
            <Icon icon="ph:chart-polar-duotone" width={26} />
          </div>
          <div>
            <div className="inline-flex items-center gap-2 text-[9px] font-mono font-bold text-[var(--md-sys-color-primary)] uppercase tracking-wider mb-0.5">
              <span>FOREIGN FLOW INTELLIGENCE</span>
              <span className="opacity-40">•</span>
              <span>IDX + ASING VS RITEL</span>
            </div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight text-[var(--md-sys-color-on-surface)]">
              Foreign Flow Pro
            </h1>
          </div>
        </div>

        {/* Global Ticker Search + Autocomplete (accessible on mobile & desktop) */}
        <div className="relative w-full lg:w-64">
          <div className="flex items-center gap-1.5 bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-xl px-2.5 py-1.5 shadow-2xs">
            <Icon
              icon="ph:magnifying-glass-bold"
              width={14}
              className="text-[var(--md-sys-color-on-surface-variant)] shrink-0"
            />
            <input
              type="text"
              value={tickerSearch}
              onChange={(e) => {
                setTickerSearch(e.target.value.toUpperCase());
                setSuggestOpen(true);
              }}
              onKeyDown={(e) => {
                if (e.key === "Enter" && tickerSearch.trim()) {
                  onSelectTicker(tickerSearch.trim().toUpperCase());
                  setActiveMainTab("stock");
                } else if (e.key === "Escape") {
                  setSuggestOpen(false);
                }
              }}
              onFocus={() => setSuggestOpen(true)}
              onBlur={() => setTimeout(() => setSuggestOpen(false), 150)}
              placeholder="Cari Ticker (e.g. BBCA)..."
              className="w-full bg-transparent text-xs font-mono font-bold text-[var(--md-sys-color-on-surface)] uppercase focus:outline-none"
            />
            <button
              type="button"
              onClick={() => {
                if (tickerSearch.trim()) {
                  onSelectTicker(tickerSearch.trim().toUpperCase());
                  setActiveMainTab("stock");
                }
              }}
              className="px-2 py-0.5 rounded-lg bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] text-[10px] font-mono font-bold hover:opacity-90 transition cursor-pointer shrink-0"
            >
              GO
            </button>
          </div>
          {suggestOpen && suggestions.length > 0 && (
            <div className="absolute top-full left-0 mt-1 w-full bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)] rounded-xl shadow-2xl z-50 max-h-56 overflow-y-auto p-1">
              {suggestions.map((t) => (
                <button
                  key={t.ticker}
                  type="button"
                  onMouseDown={(e) => {
                    e.preventDefault();
                    onSelectTicker(t.ticker);
                    setActiveMainTab("stock");
                  }}
                  className="w-full px-2.5 py-1.5 text-xs font-mono text-[var(--md-sys-color-on-surface)] hover:bg-[var(--md-sys-color-primary)] hover:text-[var(--md-sys-color-on-primary)] rounded-lg cursor-pointer flex justify-between items-center transition"
                >
                  <span className="font-bold">{t.ticker}</span>
                  <span className="opacity-70 text-[10px]">
                    {t.last_close !== undefined ? `Rp ${t.last_close}` : t.company_name || t.name || ""}
                  </span>
                </button>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* SUB-NAV VIEW SWITCHER & PARAMETER TRIGGER */}
      <div className="flex items-center justify-between gap-2 border-b border-[var(--md-sys-color-outline-variant)] pb-3 flex-wrap sm:flex-nowrap">
        <div className="flex items-center bg-[var(--md-sys-color-surface-container-high)] p-1 rounded-xl border border-[var(--md-sys-color-outline-variant)] shadow-2xs flex-wrap gap-1">
          {MAIN_TABS.map((tab) => (
            <button
              key={tab.id}
              id={tab.domId}
              type="button"
              onClick={() => setActiveMainTab(tab.id)}
              className={`px-3.5 py-1.5 rounded-lg text-xs font-bold transition flex items-center gap-1.5 cursor-pointer ${
                activeMainTab === tab.id
                  ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-xs font-extrabold"
                  : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
              }`}
            >
              <Icon icon={tab.icon} width={16} />
              <span>{tab.label}</span>
            </button>
          ))}
        </div>

        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={onOpenParameters}
            title="Atur Indikator & Parameter"
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)] cursor-pointer shadow-2xs transition"
          >
            <Icon icon="ph:sliders-horizontal-duotone" width={16} />
            <span>Parameter &amp; Indikator</span>
          </button>
          <button
            type="button"
            onClick={onRefresh}
            title="Refresh data"
            className="p-2 rounded-lg bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)] cursor-pointer transition shadow-2xs"
          >
            <Icon icon="ph:arrows-clockwise-bold" width={16} />
          </button>
        </div>
      </div>
    </>
  );
}
