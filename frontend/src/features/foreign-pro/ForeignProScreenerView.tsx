import { Icon } from "@iconify/react";
import ForeignProScreenerBasic from "./ForeignProScreenerBasic";
import ForeignProScreenerCustom from "./ForeignProScreenerCustom";
import ForeignProScreenerManual from "./ForeignProScreenerManual";
import ForeignProScreenerAdmin from "./ForeignProScreenerAdmin";
import ForeignProScreenerResults from "./ForeignProScreenerResults";
import type {
  AdminReportMode,
  IhsgResponse,
  MarketReportData,
  PeriodReportData,
  ScreenerRow,
  ScreenerSubTab,
} from "./foreignProTypes";

interface ForeignProScreenerViewProps {
  universe: string;
  setUniverse: (u: string) => void;
  windowLookback: number;
  setWindowLookback: (w: number) => void;
  screenerSubTab: ScreenerSubTab;
  setScreenerSubTab: (t: ScreenerSubTab) => void;
  screenerSearch: string;
  setScreenerSearch: (s: string) => void;
  filteredScreenerList: ScreenerRow[];
  activeStrategyId: string;
  setActiveStrategyId: (id: string) => void;
  strategyDesc: string;
  setStrategyDesc: (d: string) => void;
  activeReportType: string | null;
  setActiveReportType: (t: string | null) => void;
  periodReportWindow: number | null;
  setPeriodReportWindow: (w: number | null) => void;
  customCategory: string;
  setCustomCategory: (c: string) => void;
  transBucket: string;
  setTransBucket: (b: string) => void;
  priceClass: string;
  setPriceClass: (p: string) => void;
  stateOnly: boolean;
  setStateOnly: (v: boolean) => void;
  hfpOnly: boolean;
  setHfpOnly: (v: boolean) => void;
  mfpOnly: boolean;
  setMfpOnly: (v: boolean) => void;
  adminMode: AdminReportMode;
  setAdminMode: (m: AdminReportMode) => void;
  adminReportText: string;
  adminCopied: boolean;
  onCopyReport: () => void;
  ihsgData: IhsgResponse | null;
  marketReportData: MarketReportData;
  periodReportData: PeriodReportData;
  accumHeaderLabel: string;
  isLightMode: boolean;
  onSelectStock: (ticker: string) => void;
}

const UNIVERSE_BUTTONS = [
  { key: "ff", label: "Foreign Flow" },
  { key: "fco", label: "F&Co" },
  { key: "kf", label: "Retail" },
  { key: "xl", label: "XL" },
  { key: "semiritel", label: "SR" },
];

const SUB_TABS: Array<{ id: ScreenerSubTab; domId: string; label: string }> = [
  { id: "basic", domId: "tabBasicBtn", label: "Basic Search" },
  { id: "custom", domId: "tabCustomBtn", label: "Advanced" },
  { id: "manual", domId: "tabManualBtn", label: "Pro Search" },
  { id: "admin", domId: "tabAdminBtn", label: "Admin" },
];

export default function ForeignProScreenerView(props: ForeignProScreenerViewProps) {
  const {
    universe,
    setUniverse,
    windowLookback,
    setWindowLookback,
    screenerSubTab,
    setScreenerSubTab,
    screenerSearch,
    setScreenerSearch,
    filteredScreenerList,
    activeStrategyId,
    setActiveStrategyId,
    strategyDesc,
    setStrategyDesc,
    activeReportType,
    setActiveReportType,
    periodReportWindow,
    setPeriodReportWindow,
    customCategory,
    setCustomCategory,
    transBucket,
    setTransBucket,
    priceClass,
    setPriceClass,
    stateOnly,
    setStateOnly,
    hfpOnly,
    setHfpOnly,
    mfpOnly,
    setMfpOnly,
    adminMode,
    setAdminMode,
    adminReportText,
    adminCopied,
    onCopyReport,
    ihsgData,
    marketReportData,
    periodReportData,
    accumHeaderLabel,
    isLightMode,
    onSelectStock,
  } = props;

  return (
    <div className="space-y-3 font-mono">
      {/* Screener Universe Selector Bar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 px-3 py-2 rounded-xl gap-2 text-xs">
        <div className="flex items-center gap-2 flex-wrap">
          <span className="text-[11px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider flex items-center gap-1">
            <Icon icon="ph:globe-stand-duotone" width={14} />
            Universe:
          </span>
          <div id="screenerUniverseBtnRow" className="flex items-center gap-1.5 flex-wrap">
            {UNIVERSE_BUTTONS.map((u) => (
              <button
                key={u.key}
                id={`universeBtn_${u.key}`}
                type="button"
                onClick={() => setUniverse(u.key)}
                className={`px-2.5 py-1 rounded text-xs font-bold transition cursor-pointer ${
                  universe === u.key
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-xs"
                    : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                }`}
              >
                {u.label}
              </button>
            ))}
          </div>
        </div>
        <div className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] font-semibold hidden sm:block">
          Total Screener:{" "}
          <span className="text-[var(--md-sys-color-on-surface)] font-bold">
            {filteredScreenerList.length}
          </span>{" "}
          Emiten
        </div>
      </div>

      {/* Screener Top Sub-Tabs */}
      <div className="flex items-center justify-between bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-2.5 rounded-xl flex-wrap gap-2">
        <div className="flex items-center gap-2 flex-wrap">
          {SUB_TABS.map((t) => (
            <button
              key={t.id}
              id={t.domId}
              type="button"
              onClick={() => {
                setScreenerSubTab(t.id);
                setActiveReportType(null);
                setPeriodReportWindow(null);
              }}
              className={`px-3 py-1 rounded text-xs font-bold transition cursor-pointer ${
                screenerSubTab === t.id
                  ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                  : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>

        <div className="flex items-center gap-2">
          <input
            type="text"
            id="searchInput"
            value={screenerSearch}
            onChange={(e) => setScreenerSearch(e.target.value)}
            placeholder="Cari kode saham..."
            className="bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/60 rounded px-2.5 py-1 text-xs text-[var(--md-sys-color-on-surface)] uppercase focus:outline-none focus:border-[var(--md-sys-color-primary)]"
          />
        </div>
      </div>

      {screenerSubTab === "basic" && (
        <ForeignProScreenerBasic
          activeStrategyId={activeStrategyId}
          setActiveStrategyId={setActiveStrategyId}
          strategyDesc={strategyDesc}
          setStrategyDesc={setStrategyDesc}
          activeReportType={activeReportType}
          setActiveReportType={setActiveReportType}
          periodReportWindow={periodReportWindow}
          setPeriodReportWindow={setPeriodReportWindow}
        />
      )}

      {screenerSubTab === "custom" && (
        <ForeignProScreenerCustom
          universe={universe}
          setUniverse={setUniverse}
          customCategory={customCategory}
          setCustomCategory={setCustomCategory}
          activeStrategyId={activeStrategyId}
          setActiveStrategyId={setActiveStrategyId}
          strategyDesc={strategyDesc}
          setStrategyDesc={setStrategyDesc}
        />
      )}

      {screenerSubTab === "manual" && (
        <ForeignProScreenerManual
          universe={universe}
          setUniverse={setUniverse}
          windowLookback={windowLookback}
          setWindowLookback={setWindowLookback}
          transBucket={transBucket}
          setTransBucket={setTransBucket}
          priceClass={priceClass}
          setPriceClass={setPriceClass}
          stateOnly={stateOnly}
          setStateOnly={setStateOnly}
          hfpOnly={hfpOnly}
          setHfpOnly={setHfpOnly}
          mfpOnly={mfpOnly}
          setMfpOnly={setMfpOnly}
        />
      )}

      {screenerSubTab === "admin" && (
        <ForeignProScreenerAdmin
          adminMode={adminMode}
          setAdminMode={setAdminMode}
          adminReportText={adminReportText}
          adminCopied={adminCopied}
          onCopyReport={onCopyReport}
          ihsgData={ihsgData}
        />
      )}

      <ForeignProScreenerResults
        activeReportType={activeReportType}
        marketReportData={marketReportData}
        periodReportWindow={periodReportWindow}
        periodReportData={periodReportData}
        filteredScreenerList={filteredScreenerList}
        universe={universe}
        windowLookback={windowLookback}
        accumHeaderLabel={accumHeaderLabel}
        isLightMode={isLightMode}
        onSelectStock={onSelectStock}
      />
    </div>
  );
}
