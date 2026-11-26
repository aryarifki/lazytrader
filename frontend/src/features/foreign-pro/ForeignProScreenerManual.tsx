interface ForeignProScreenerManualProps {
  universe: string;
  setUniverse: (u: string) => void;
  windowLookback: number;
  setWindowLookback: (w: number) => void;
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
}

const UNIVERSES = [
  { key: "ff", label: "FF" },
  { key: "fco", label: "F&Co" },
  { key: "kf", label: "Retail" },
  { key: "xl", label: "XL" },
  { key: "semiritel", label: "SR" },
];

const WINDOWS = [5, 10, 20, 40, 60, 120];

export default function ForeignProScreenerManual({
  universe,
  setUniverse,
  windowLookback,
  setWindowLookback,
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
}: ForeignProScreenerManualProps) {
  const activeChip =
    "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]";
  const idleChip =
    "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]";

  return (
    <div
      id="manualControls"
      className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 p-3 rounded-xl space-y-3"
    >
      <div className="grid grid-cols-1 md:grid-cols-12 gap-3 text-xs">
        <div className="md:col-span-3 space-y-1">
          <label className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
            Universe:
          </label>
          <div id="universeChipRow" className="flex items-center gap-1.5 flex-wrap">
            {UNIVERSES.map((u) => (
              <button
                key={u.key}
                type="button"
                onClick={() => setUniverse(u.key)}
                className={`px-2 py-0.5 rounded text-xs font-bold transition cursor-pointer ${
                  universe === u.key ? activeChip : idleChip
                }`}
              >
                {u.label}
              </button>
            ))}
          </div>
        </div>

        <div className="md:col-span-3 space-y-1">
          <label className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
            Lookback Window
          </label>
          <div id="windowPresetRow" className="flex items-center gap-1 flex-wrap">
            {WINDOWS.map((w) => (
              <button
                key={w}
                type="button"
                onClick={() => setWindowLookback(w)}
                className={`px-2 py-0.5 rounded text-xs font-bold transition cursor-pointer ${
                  windowLookback === w ? activeChip : idleChip
                }`}
              >
                {w}d
              </button>
            ))}
          </div>
        </div>

        <div className="md:col-span-3 space-y-1">
          <label className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
            Turnover 10D
          </label>
          <select
            id="transBucket"
            value={transBucket}
            onChange={(e) => setTransBucket(e.target.value)}
            className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/60 rounded px-2 py-1 text-xs text-[var(--md-sys-color-on-surface)]"
          >
            <option value="all">All</option>
            <option value="gt35">&gt; 35 M</option>
            <option value="10to35">10 M - 35 M</option>
            <option value="lt10">&lt; 10 M</option>
          </select>
        </div>

        <div className="md:col-span-3 space-y-1">
          <label className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
            Price Change
          </label>
          <select
            id="priceClassSelect"
            value={priceClass}
            onChange={(e) => setPriceClass(e.target.value)}
            className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)]/60 rounded px-2 py-1 text-xs text-[var(--md-sys-color-on-surface)]"
          >
            <option value="All">All</option>
            <option value="Up">Up</option>
            <option value="Down">Down</option>
            <option value="Sideways">Sideways</option>
          </select>
        </div>
      </div>

      <div className="flex items-center gap-4 text-xs pt-1 border-t border-[var(--md-sys-color-outline-variant)]/40 flex-wrap">
        <label className="flex items-center gap-1.5 cursor-pointer">
          <input
            type="checkbox"
            id="ffStateOnly"
            checked={stateOnly}
            onChange={(e) => setStateOnly(e.target.checked)}
            className="accent-[var(--md-sys-color-primary)]"
          />
          <span id="ffStateLabel" className="text-[var(--md-sys-color-on-surface)] font-bold">
            {universe === "ff" || universe === "fco" ? "Accumulation Only" : "Retail Sell Only"}
          </span>
        </label>

        <label className="flex items-center gap-1.5 cursor-pointer">
          <input
            type="checkbox"
            id="hfpOnly"
            checked={hfpOnly}
            onChange={(e) => setHfpOnly(e.target.checked)}
            className="accent-[var(--md-sys-color-primary)]"
          />
          <span className="text-[var(--md-sys-color-on-surface)]">HFP Only</span>
        </label>

        <label className="flex items-center gap-1.5 cursor-pointer">
          <input
            type="checkbox"
            id="mfpOnly"
            checked={mfpOnly}
            onChange={(e) => setMfpOnly(e.target.checked)}
            className="accent-[var(--md-sys-color-primary)]"
          />
          <span className="text-[var(--md-sys-color-on-surface)]">MFP Only</span>
        </label>
      </div>
    </div>
  );
}
