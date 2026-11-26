import { bandarBadgeClass } from "./mmdBadges";
import type { MmdAnalysis } from "./mmdTypes";

interface MmdBandarDetailTabProps {
  analysisData: MmdAnalysis | null;
  retailBandarAutoEnabled: boolean;
  setRetailBandarAutoEnabled: (enabled: boolean) => void;
  retailBandarManualOverride: string;
  setRetailBandarManualOverride: (broker: string) => void;
}

const RETAIL_BROKER_OVERRIDES = ["YP", "PD", "XC", "NI", "KK", "XL", "EP", "CC"];

export default function MmdBandarDetailTab({
  analysisData,
  retailBandarAutoEnabled,
  setRetailBandarAutoEnabled,
  retailBandarManualOverride,
  setRetailBandarManualOverride,
}: MmdBandarDetailTabProps) {
  const pengendali = analysisData?.pengendali_stats;

  return (
    <div className="space-y-3 font-mono-nums text-xs">
      <div className="bg-[var(--md-sys-color-surface-container-low)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40 space-y-2">
        <div className="flex items-center justify-between flex-wrap gap-2 pb-1 border-b border-[var(--md-sys-color-outline-variant)]/40">
          <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)] tracking-wider flex items-center gap-1.5">
            <span className="w-2 h-2 rounded-full bg-blue-500" />
            BROKER RETAIL DIPAKE BANDAR
          </span>
          <div className="flex items-center gap-2">
            {analysisData?.retail_bandar_auto_result && (
              <div className="flex items-center gap-1.5 bg-[var(--md-sys-color-surface-container-high)] px-2 py-0.5 rounded-lg border border-[var(--md-sys-color-outline-variant)]/50">
                <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">Auto:</span>
                <span className="text-xs font-bold text-blue-600 dark:text-blue-400">
                  {analysisData.retail_bandar_auto_result}
                </span>
                <button
                  type="button"
                  onClick={() => setRetailBandarAutoEnabled(!retailBandarAutoEnabled)}
                  className={`px-1.5 py-0.5 rounded text-[9px] font-bold transition cursor-pointer ${
                    retailBandarAutoEnabled
                      ? "bg-emerald-500 text-black"
                      : "bg-rose-500 text-white"
                  }`}
                >
                  {retailBandarAutoEnabled ? "ON" : "OFF"}
                </button>
              </div>
            )}
            <select
              value={retailBandarManualOverride}
              onChange={(e) => setRetailBandarManualOverride(e.target.value)}
              className="bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-[10px] rounded px-1.5 py-0.5 text-[var(--md-sys-color-on-surface)] focus:outline-none focus:border-[var(--md-sys-color-primary)] font-mono cursor-pointer"
            >
              <option value="">— Override Manual —</option>
              {RETAIL_BROKER_OVERRIDES.map((b) => (
                <option key={b} value={b}>
                  {b}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container-low)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40 space-y-2">
        <div className="text-xs font-bold text-[var(--md-sys-color-on-surface)] pb-1 border-b border-[var(--md-sys-color-outline-variant)]/40">
          TOP BANDAR (1, 3, 5) DETAIL
        </div>
        <table className="w-full text-left text-xs">
          <thead>
            <tr className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/40">
              <th className="py-1">Level</th>
              <th className="py-1 text-right">Vol Buy</th>
              <th className="py-1 text-right">Vol Sell</th>
              <th className="py-1 text-right">Ratio</th>
              <th className="py-1 text-center">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/30">
            {(analysisData?.bandar_table || []).map((row, idx) => (
              <tr key={idx}>
                <td className="py-1.5 font-bold text-[var(--md-sys-color-on-surface)]">{row.label}</td>
                <td className="py-1.5 text-right text-[var(--color-positive)]">
                  {row.buy?.toLocaleString("id-ID") || 0}
                </td>
                <td className="py-1.5 text-right text-[var(--color-negative)]">
                  {row.sell?.toLocaleString("id-ID") || 0}
                </td>
                <td className="py-1.5 text-right font-bold">{row.ratio}%</td>
                <td className="py-1.5 text-center">
                  <span className={`${bandarBadgeClass(row.badge?.cls)} text-[9px] py-0.5 px-1.5`}>
                    {row.badge?.label || "—"}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container-low)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40 space-y-2">
        <div className="flex items-center justify-between pb-1 border-b border-[var(--md-sys-color-outline-variant)]/40">
          <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)]">
            RETAIL ACTION (BROKER TRANSAKSI RITEL)
          </span>
          <span className={`${bandarBadgeClass(analysisData?.retail_action?.cls)} text-[9px] py-0.5 px-2`}>
            {analysisData?.retail_action?.label || "Do Nothing"}
          </span>
        </div>
        <table className="w-full text-left text-xs">
          <thead>
            <tr className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)]/40">
              <th className="py-1">Broker</th>
              <th className="py-1 text-right">Trans</th>
              <th className="py-1 text-right">%Trans</th>
              <th className="py-1 text-right">Net</th>
              <th className="py-1 text-right">%Net</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/30">
            {analysisData?.retail_rows && analysisData.retail_rows.length > 0 ? (
              analysisData.retail_rows.map((row, idx) => (
                <tr key={idx} className="hover:bg-[var(--md-sys-color-surface-container-high)] transition">
                  <td className="py-1 font-bold text-[var(--md-sys-color-on-surface)] flex items-center gap-1.5">
                    <span className="w-1.5 h-1.5 rounded-full bg-fuchsia-500 dark:bg-[#E879F9]" />
                    <span>{row.broker}</span>
                  </td>
                  <td className="py-1 text-right">{row.trans?.toLocaleString("id-ID") || 0}</td>
                  <td className="py-1 text-right text-[var(--md-sys-color-on-surface-variant)]">
                    {row.pct_trans || 0}%
                  </td>
                  <td
                    className={`py-1 text-right font-bold ${
                      row.net >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                    }`}
                  >
                    {row.net?.toLocaleString("id-ID") || 0}
                  </td>
                  <td
                    className={`py-1 text-right font-bold ${
                      row.pct_net >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
                    }`}
                  >
                    {row.pct_net > 0 ? `+${row.pct_net}%` : `${row.pct_net}%`}
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={5} className="py-2 text-center text-[var(--md-sys-color-outline)]">
                  Memuat data broker retail...
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container-low)] rounded-xl p-3 border border-[var(--md-sys-color-outline-variant)]/40 space-y-2.5">
        <div className="flex items-center justify-between pb-1 border-b border-[var(--md-sys-color-outline-variant)]/40">
          <span className="text-xs font-bold text-[var(--md-sys-color-on-surface)]">
            BROKER PENGENDALI (KATEGORI PARTISIPASI)
          </span>
          <span className={`${bandarBadgeClass(pengendali?.action?.cls)} text-[9px] py-0.5 px-2`}>
            {pengendali?.action?.label || "—"}
          </span>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-5 gap-2 text-center">
          <div className="bg-[var(--md-sys-color-surface-container)] p-2 rounded-lg border border-[var(--md-sys-color-outline-variant)]/40">
            <div className="text-[9px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
              INSTITUSI
            </div>
            <div className="text-sm font-black text-[var(--md-sys-color-on-surface)] mt-0.5">
              {pengendali?.institusi?.main || 0}%
            </div>
            <div className="text-[10px] text-[var(--color-positive)] font-semibold mt-0.5">
              {pengendali?.institusi?.sub || 0}%
            </div>
          </div>

          <div className="bg-[var(--md-sys-color-surface-container)] p-2 rounded-lg border border-[var(--md-sys-color-outline-variant)]/40">
            <div className="text-[9px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
              SEMI RITEL
            </div>
            <div className="text-sm font-black text-[var(--md-sys-color-on-surface)] mt-0.5">
              {pengendali?.semi_ritel?.main || 0}%
            </div>
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-semibold mt-0.5">
              {pengendali?.semi_ritel?.sub || 0}%
            </div>
          </div>

          <div className="bg-[var(--md-sys-color-surface-container)] p-2 rounded-lg border border-[var(--md-sys-color-outline-variant)]/40">
            <div className="text-[9px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
              RITEL
            </div>
            <div className="text-sm font-black text-[var(--md-sys-color-on-surface)] mt-0.5">
              {pengendali?.ritel?.main || 0}%
            </div>
            <div className="text-[10px] text-[var(--color-negative)] font-semibold mt-0.5">
              {pengendali?.ritel?.sub || 0}%
            </div>
          </div>

          <div className="bg-[var(--md-sys-color-surface-container)] p-2 rounded-lg border border-[var(--md-sys-color-outline-variant)]/40">
            <div className="text-[9px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
              NON RITEL
            </div>
            <div className="text-sm font-black text-[var(--md-sys-color-on-surface)] mt-0.5">
              {pengendali?.non_ritel?.main || 0}%
            </div>
            <div className="text-[10px] text-[var(--color-positive)] font-semibold mt-0.5">
              {pengendali?.non_ritel?.sub || 0}%
            </div>
          </div>

          <div className="bg-[var(--md-sys-color-surface-container)] p-2 rounded-lg border border-[var(--md-sys-color-outline-variant)]/40 col-span-2 sm:col-span-1">
            <div className="text-[9px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
              REPO / LAIN
            </div>
            <div className="text-xs font-bold text-[var(--md-sys-color-on-surface)] mt-1">
              Repo: {pengendali?.repo?.val || 0}%
            </div>
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] mt-0.5">
              Lain: {pengendali?.lainnya?.val || 0}%
            </div>
          </div>
        </div>

        <div className="flex items-center justify-between pt-1 px-1 text-xs border-t border-[var(--md-sys-color-outline-variant)]/30">
          <span className="text-[var(--md-sys-color-on-surface-variant)]">Partisipasi Pengendali:</span>
          <span className="font-extrabold text-[var(--color-positive)]">{pengendali?.ratio || 0}%</span>
        </div>
      </div>
    </div>
  );
}
