import { Icon } from "@iconify/react";
import { fmtRp, signedColor, getProfileTextColor } from "@/lib/format";
import type {
  BrokerComparisonSummaryRow,
  DetailedBrokerRow,
  DistributionEdge,
  ProfileBrokerDetailRow,
} from "./types";

const PROFILE_OPTIONS = [
  { key: "all", label: "All Profiles" },
  { key: "smart_foreign", label: "Foreign Smart Money" },
  { key: "local_institutional", label: "Local Institutions" },
  { key: "market_maker", label: "Market Makers" },
  { key: "bandar_gorengan", label: "Speculative Operators" },
  { key: "retail", label: "Retail-Dominant" },
  { key: "lainnya", label: "Other Brokers" },
];

interface BrokerFlowDetailsProps {
  profileFilter: string;
  onProfileFilterChange: (value: string) => void;
  filteredProfileDetail: ProfileBrokerDetailRow[];
  isValidating: boolean;
  distMode: string;
  onDistModeChange: (value: string) => void;
  distDate: string;
  distStart: string;
  distEnd: string;
  onDistDateChange: (value: string) => void;
  onDistStartChange: (value: string) => void;
  onDistEndChange: (value: string) => void;
  edges: DistributionEdge[];
  summary: BrokerComparisonSummaryRow[];
  detailRows: DetailedBrokerRow[];
}

export default function BrokerFlowDetails({
  profileFilter,
  onProfileFilterChange,
  filteredProfileDetail,
  isValidating,
  distMode,
  onDistModeChange,
  distDate,
  distStart,
  distEnd,
  onDistDateChange,
  onDistStartChange,
  onDistEndChange,
  edges,
  summary,
  detailRows,
}: BrokerFlowDetailsProps) {
  return (
    <>
      <div className="bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)] rounded-[28px] p-5 sm:p-6 shadow-sm transition-colors duration-300">
        <label className="block text-[10px] font-extrabold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-widest mb-2.5">
          Profile detail filter
        </label>
        <select
          className="w-full md:w-72 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded-full px-4 py-2.5 text-sm font-bold text-[var(--md-sys-color-on-surface)] outline-none focus:border-[var(--md-sys-color-primary)] shadow-sm cursor-pointer mb-5"
          value={profileFilter}
          onChange={(e) => onProfileFilterChange(e.target.value)}
        >
          {PROFILE_OPTIONS.map((p) => (
            <option key={p.key} value={p.key}>
              {p.label}
            </option>
          ))}
        </select>

        <div className="bg-[var(--md-sys-color-surface)] border border-[var(--md-sys-color-outline-variant)] rounded-[20px] shadow-sm overflow-hidden transition-colors duration-300">
          <div className="max-h-80 overflow-y-auto scrollbar-thin scrollbar-thumb-[var(--md-sys-color-outline-variant)] pb-2">
            <table className="w-full text-xs text-left whitespace-nowrap">
              <thead className="sticky top-0 bg-[var(--md-sys-color-surface-container-highest)] z-10 shadow-sm">
                <tr className="text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)] text-[9px] uppercase tracking-widest font-extrabold">
                  <th className="py-3 px-4">Profile</th>
                  <th className="py-3 px-4">Broker</th>
                  <th className="py-3 px-4">Type</th>
                  <th className="py-3 px-4 text-right">Buy</th>
                  <th className="py-3 px-4 text-right">Sell</th>
                  <th className="py-3 px-4 text-right">Net</th>
                  <th className="py-3 px-4 text-right">Freq</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
                {filteredProfileDetail.map((row, i) => (
                  <tr key={i} className="hover:bg-[var(--md-sys-color-surface-container-highest)] transition-colors">
                    <td
                      className="py-3 px-4 font-extrabold uppercase tracking-wide"
                      style={{ color: getProfileTextColor(row.profile) }}
                    >
                      {row.profile}
                    </td>
                    <td className="py-3 px-4 text-[var(--md-sys-color-on-surface)] font-extrabold font-mono">{row.broker}</td>
                    <td className="py-3 px-4 font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">{row.type}</td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface)]">
                      {fmtRp(row.buy)}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface)]">
                      {fmtRp(row.sell)}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-extrabold" style={{ color: signedColor(row.net) }}>
                      {fmtRp(row.net)}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface-variant)]">
                      {row.freq.toLocaleString("id-ID")}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {filteredProfileDetail.length === 0 && (
              <div className="text-center text-xs font-bold text-[var(--md-sys-color-on-surface-variant)] py-8">
                No data for selected profile.
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)] rounded-[28px] p-5 sm:p-6 shadow-sm transition-colors duration-300">
        <div className="flex items-center gap-2 mb-5 border-b border-[var(--md-sys-color-outline-variant)] pb-3">
          <Icon icon="ph:link-bold" className="text-[var(--md-sys-color-primary)]" width="20" />
          <h3 className="text-sm sm:text-base font-extrabold text-[var(--md-sys-color-on-surface)] flex items-center gap-2 tracking-tight">
            Broker Distribution
            {isValidating && <Icon icon="ph:spinner-gap-bold" className="animate-spin text-[var(--md-sys-color-primary)]" width="16" />}
          </h3>
        </div>

        <div className="flex flex-col md:flex-row gap-5 mb-5">
          <div className="flex-1">
            <label className="block text-[10px] font-extrabold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-widest mb-2.5">
              Distribution mode
            </label>
            <select
              className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded-full px-4 py-2.5 text-sm font-bold text-[var(--md-sys-color-on-surface)] outline-none focus:border-[var(--md-sys-color-primary)] shadow-sm cursor-pointer"
              value={distMode}
              onChange={(e) => onDistModeChange(e.target.value)}
            >
              <option value="Single day">Single day</option>
              <option value="Date range">Date range</option>
            </select>
          </div>

          <div className="flex-1">
            <label className="block text-[10px] font-extrabold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-widest mb-2.5">
              Distribution date
            </label>
            {distMode === "Single day" ? (
              <input
                type="date"
                className="w-full bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded-full px-4 py-2.5 text-sm font-bold text-[var(--md-sys-color-on-surface)] outline-none focus:border-[var(--md-sys-color-primary)] shadow-sm"
                value={distDate}
                onChange={(e) => onDistDateChange(e.target.value)}
              />
            ) : (
              <div className="flex flex-col sm:flex-row items-center gap-2">
                <input
                  type="date"
                  className="w-full sm:w-1/2 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded-full px-4 py-2.5 text-sm font-bold text-[var(--md-sys-color-on-surface)] outline-none focus:border-[var(--md-sys-color-primary)] shadow-sm"
                  value={distStart}
                  onChange={(e) => onDistStartChange(e.target.value)}
                />
                <span className="text-[var(--md-sys-color-on-surface-variant)] text-xs font-bold uppercase">to</span>
                <input
                  type="date"
                  className="w-full sm:w-1/2 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded-full px-4 py-2.5 text-sm font-bold text-[var(--md-sys-color-on-surface)] outline-none focus:border-[var(--md-sys-color-primary)] shadow-sm"
                  value={distEnd}
                  onChange={(e) => onDistEndChange(e.target.value)}
                />
              </div>
            )}
          </div>
        </div>
        <p className="text-xs font-medium text-[var(--md-sys-color-on-surface-variant)] mb-5">
          The flow chart below uses broker-to-broker distribution edges returned by the live API.
        </p>

        {edges.length > 0 ? (
          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-[20px] p-5 shadow-sm transition-colors duration-300">
            <h4 className="text-sm font-extrabold text-[var(--md-sys-color-on-surface)] mb-4 tracking-tight">
              Estimated Counterparties on {distMode === "Single day" ? distDate : `${distStart} to ${distEnd}`}
            </h4>
            <div className="max-h-64 overflow-y-auto scrollbar-thin scrollbar-thumb-[var(--md-sys-color-outline-variant)] space-y-2 pr-2">
              {edges.map((e, i) => (
                <div
                  key={i}
                  className="flex items-center justify-between text-xs bg-[var(--md-sys-color-surface)] border border-[var(--md-sys-color-outline-variant)] rounded-xl px-4 py-3 shadow-sm hover:border-[var(--md-sys-color-primary)] transition-colors"
                >
                  <div className="flex items-center gap-2.5 w-1/3">
                    <span className="text-[var(--color-positive)] font-mono font-extrabold">{e.buyer_code}</span>
                    <span className="text-[9px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
                      {e.buyer_type}
                    </span>
                  </div>
                  <div className="w-1/3 text-center">
                    <Icon icon="ph:arrow-right-bold" className="text-[var(--md-sys-color-outline)] inline-block" width="16" />
                  </div>
                  <div className="flex items-center justify-end gap-2.5 w-1/3">
                    <span className="text-[9px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider">
                      {e.seller_type}
                    </span>
                    <span className="text-[var(--color-negative)] font-mono font-extrabold">{e.seller_code}</span>
                  </div>
                  <span className="text-[var(--md-sys-color-on-surface)] font-mono font-extrabold tabular-nums ml-4">
                    {fmtRp(e.matched_value)}
                  </span>
                </div>
              ))}
            </div>
          </div>
        ) : (
          <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-[20px] p-5 shadow-sm mt-4 text-center text-xs font-bold text-[var(--md-sys-color-on-surface-variant)] py-8 transition-colors duration-300">
            No distribution data found for the selected dates.
          </div>
        )}
      </div>

      <div className="bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)] rounded-[28px] p-5 sm:p-6 shadow-sm transition-colors duration-300">
        <label className="block text-[10px] font-extrabold text-[var(--md-sys-color-on-surface)] uppercase tracking-widest mb-4">
          Broker Summary
        </label>
        <div className="bg-[var(--md-sys-color-surface)] border border-[var(--md-sys-color-outline-variant)] rounded-[20px] shadow-sm overflow-hidden transition-colors duration-300">
          <div className="max-h-80 overflow-y-auto scrollbar-thin scrollbar-thumb-[var(--md-sys-color-outline-variant)] pb-2">
            <table className="w-full text-xs text-left whitespace-nowrap">
              <thead className="sticky top-0 bg-[var(--md-sys-color-surface-container-highest)] z-10 shadow-sm">
                <tr className="text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)] text-[9px] uppercase tracking-widest font-extrabold">
                  <th className="py-3 px-4">Buy Broker</th>
                  <th className="py-3 px-4">Buy Type</th>
                  <th className="py-3 px-4 text-right">Buy Value</th>
                  <th className="py-3 px-4 text-right">Buy Lot</th>
                  <th className="py-3 px-4 text-right">Buy Avg</th>
                  <th className="py-3 px-4">Sell Broker</th>
                  <th className="py-3 px-4">Sell Type</th>
                  <th className="py-3 px-4 text-right">Sell Value</th>
                  <th className="py-3 px-4 text-right">Sell Lot</th>
                  <th className="py-3 px-4 text-right">Sell Avg</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
                {summary.map((row, i) => (
                  <tr key={i} className="hover:bg-[var(--md-sys-color-surface-container-highest)] transition-colors">
                    <td className="py-3 px-4 text-[var(--md-sys-color-on-surface)] font-extrabold font-mono">
                      {row.buy_broker || ""}
                    </td>
                    <td className="py-3 px-4 font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">
                      {row.buy_type || ""}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface)]">
                      {fmtRp(row.buy_value)}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface-variant)]">
                      {row.buy_lot ? row.buy_lot.toLocaleString("id-ID") : "None"}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface-variant)]">
                      {row.buy_avg ? row.buy_avg.toFixed(0) : "None"}
                    </td>
                    <td className="py-3 px-4 text-[var(--md-sys-color-on-surface)] font-extrabold font-mono">
                      {row.sell_broker || ""}
                    </td>
                    <td className="py-3 px-4 font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">
                      {row.sell_type || ""}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface)]">
                      {fmtRp(row.sell_value)}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface-variant)]">
                      {row.sell_lot ? row.sell_lot.toLocaleString("id-ID") : "None"}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface-variant)]">
                      {row.sell_avg ? row.sell_avg.toFixed(0) : "None"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)] rounded-[28px] p-5 sm:p-6 shadow-sm transition-colors duration-300">
        <label className="block text-[10px] font-extrabold text-[var(--md-sys-color-on-surface)] uppercase tracking-widest mb-4">
          Detailed broker rows
        </label>
        <div className="bg-[var(--md-sys-color-surface)] border border-[var(--md-sys-color-outline-variant)] rounded-[20px] shadow-sm overflow-hidden transition-colors duration-300">
          <div className="max-h-80 overflow-y-auto scrollbar-thin scrollbar-thumb-[var(--md-sys-color-outline-variant)] pb-2">
            <table className="w-full text-xs text-left whitespace-nowrap">
              <thead className="sticky top-0 bg-[var(--md-sys-color-surface-container-highest)] z-10 shadow-sm">
                <tr className="text-[var(--md-sys-color-on-surface-variant)] border-b border-[var(--md-sys-color-outline-variant)] text-[9px] uppercase tracking-widest font-extrabold">
                  <th className="py-3 px-4">Broker</th>
                  <th className="py-3 px-4">Type</th>
                  <th className="py-3 px-4">Sub-type</th>
                  <th className="py-3 px-4 text-right">Buy</th>
                  <th className="py-3 px-4 text-right">Sell</th>
                  <th className="py-3 px-4 text-right">Net</th>
                  <th className="py-3 px-4 text-right">Freq</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
                {detailRows.map((row, i) => (
                  <tr key={i} className="hover:bg-[var(--md-sys-color-surface-container-highest)] transition-colors">
                    <td className="py-3 px-4 text-[var(--md-sys-color-on-surface)] font-extrabold font-mono">{row.broker}</td>
                    <td className="py-3 px-4 font-bold uppercase text-[var(--md-sys-color-on-surface-variant)]">{row.type}</td>
                    <td className="py-3 px-4 font-bold text-xs">
                      {row.subtype && row.subtype !== "-" ? (
                        <span
                          className={`inline-block px-2 py-0.5 rounded-full text-[10px] font-extrabold tracking-wide ${
                            row.subtype === "Institutional"
                              ? "bg-emerald-500/10 text-emerald-500 border border-emerald-500/20"
                              : row.subtype === "Speculative"
                              ? "bg-amber-500/10 text-amber-500 border border-amber-500/20"
                              : "bg-blue-500/10 text-blue-500 border border-blue-500/20"
                          }`}
                        >
                          {row.subtype}
                        </span>
                      ) : (
                        <span className="text-[var(--md-sys-color-on-surface-variant)] opacity-50">-</span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface)]">
                      {row.buy ? fmtRp(row.buy) : "Rp 0"}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface)]">
                      {row.sell ? fmtRp(row.sell) : "Rp 0"}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-extrabold" style={{ color: signedColor(row.net) }}>
                      {fmtRp(row.net)}
                    </td>
                    <td className="py-3 px-4 text-right font-mono font-bold text-[var(--md-sys-color-on-surface-variant)]">
                      {row.freq ? row.freq.toLocaleString("id-ID") : "0"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </>
  );
}
