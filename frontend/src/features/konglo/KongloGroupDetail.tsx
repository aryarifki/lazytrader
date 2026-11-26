import { KongloTickerLink } from "./KongloShared";
import { fmtFreq, fmtM, fmtMc, fmtT, type KongloGroupStats } from "./kongloTypes";

export default function KongloGroupDetail({
  group,
}: {
  group: KongloGroupStats | null;
}) {
  if (!group) return null;

  return (
    <div id="kongloResult">
      <div className="space-y-4">
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl p-4 shadow-2xs">
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] uppercase font-bold tracking-wider mb-1 font-mono">
              Total Market Cap
            </div>
            <div className="text-xl font-bold text-[var(--md-sys-color-on-surface)] font-mono">
              Rp {fmtT(group.total_mc)} T
            </div>
            <div className="text-xs text-[var(--md-sys-color-on-surface-variant)] mt-1">
              {group.total_members} saham dalam group
            </div>
          </div>

          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl p-4 shadow-2xs">
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] uppercase font-bold tracking-wider mb-1 font-mono">
              Perubahan MC
            </div>
            <div
              className={`text-xl font-bold font-mono ${
                group.mc_chg >= 0
                  ? "text-emerald-600 dark:text-emerald-400"
                  : "text-rose-600 dark:text-rose-400"
              }`}
            >
              {fmtMc(group.mc_chg)}
            </div>
            <div
              className={`text-xs mt-1 font-mono ${
                group.mc_pct >= 0
                  ? "text-emerald-600 dark:text-emerald-400"
                  : "text-rose-600 dark:text-rose-400"
              }`}
            >
              {group.mc_pct >= 0
                ? `+${group.mc_pct.toFixed(2)}%`
                : `${group.mc_pct.toFixed(2)}%`}{" "}
              dari MC sebelumnya
            </div>
          </div>

          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl p-4 shadow-2xs">
            <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] uppercase font-bold tracking-wider mb-1 font-mono">
              Total Value
            </div>
            <div className="text-xl font-bold text-[var(--md-sys-color-on-surface)] font-mono">
              {fmtM(group.total_val)}
            </div>
            <div className="text-xs text-[var(--md-sys-color-on-surface-variant)] mt-1">
              {fmtFreq(group.total_freq)} transaksi hari ini
            </div>
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-2xl overflow-hidden shadow-2xs">
          <div className="px-4 py-3 border-b border-[var(--md-sys-color-outline-variant)]/50 bg-[var(--md-sys-color-surface-container-high)]/60 flex items-center justify-between">
            <span className="text-xs font-bold tracking-wide text-[var(--md-sys-color-on-surface)] uppercase">
              ANGGOTA SAHAM {group.group}
            </span>
            <span className="text-xs text-[var(--md-sys-color-on-surface-variant)] font-mono">
              {group.members?.length ?? 0} Saham
            </span>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead>
                <tr className="border-b border-[var(--md-sys-color-outline-variant)]/40 bg-[var(--md-sys-color-surface-container-high)]/30 text-[var(--md-sys-color-on-surface-variant)] font-mono uppercase text-[10px]">
                  <th className="py-2.5 px-3 w-8 font-bold">#</th>
                  <th className="py-2.5 px-3 font-bold">Saham</th>
                  <th className="py-2.5 px-3 text-right font-bold">Change %</th>
                  <th className="py-2.5 px-3 text-right font-bold">Value (M)</th>
                  <th className="py-2.5 px-3 text-right font-bold">Freq</th>
                  <th className="py-2.5 px-3 text-right font-bold">Perubahan MC</th>
                </tr>
              </thead>
              <tbody>
                {group.members?.map((m, idx) => (
                  <tr
                    key={m.ticker}
                    className="border-b border-[var(--md-sys-color-outline-variant)]/20 hover:bg-[var(--md-sys-color-surface-container-high)]/40 transition-colors"
                  >
                    <td className="py-2 px-3 text-[var(--md-sys-color-on-surface-variant)] font-mono">
                      {idx + 1}
                    </td>
                    <td className="py-2 px-3">
                      <KongloTickerLink ticker={m.ticker} className="text-sm" />
                    </td>
                    <td
                      className={`py-2 px-3 text-right font-semibold font-mono ${
                        m.pct >= 0
                          ? "text-emerald-600 dark:text-emerald-400"
                          : "text-rose-600 dark:text-rose-400"
                      }`}
                    >
                      {m.pct > 0 ? `+${m.pct.toFixed(1)}%` : `${m.pct.toFixed(1)}%`}
                    </td>
                    <td className="py-2 px-3 text-right text-[var(--md-sys-color-on-surface-variant)] font-mono">
                      {fmtM(m.val)}
                    </td>
                    <td className="py-2 px-3 text-right text-[var(--md-sys-color-on-surface-variant)] font-mono">
                      {fmtFreq(m.freq)}
                    </td>
                    <td
                      className={`py-2 px-3 text-right font-semibold font-mono ${
                        m.mc_chg >= 0
                          ? "text-emerald-600 dark:text-emerald-400"
                          : "text-rose-600 dark:text-rose-400"
                      }`}
                    >
                      {fmtMc(m.mc_chg)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
