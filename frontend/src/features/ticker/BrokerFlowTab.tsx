import { useEffect, useMemo, useState } from "react";
import useSWR from "swr";
import { Icon } from "@iconify/react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  ReferenceLine,
} from "recharts";
import { fmtRp, signedColor, getProfileTextColor } from "@/lib/format";
import BrokerFlowDetails from "./BrokerFlowDetails";
import type { BrokerFlowData } from "./types";

const fetcher = (url: string) => fetch(url, { credentials: "include" }).then((r) => r.json());

const COLORS = [
  "#3b82f6",
  "var(--color-negative)",
  "var(--color-positive)",
  "#f59e0b",
  "#8b5cf6",
  "#06b6d4",
  "#ec4899",
  "#84cc16",
  "#f97316",
  "#6366f1",
];

export default function BrokerFlowTab({
  ticker,
  analysisDate,
  windowDays,
}: {
  ticker: string;
  analysisDate: string;
  windowDays: number;
}) {
  const [selectedCodes, setSelectedCodes] = useState<string[]>([]);
  const [flowMode, setFlowMode] = useState("Cumulative");
  const [distMode, setDistMode] = useState("Date range");
  const [customDistDate, setCustomDistDate] = useState<string | null>(null);
  const [customDistStart, setCustomDistStart] = useState<string | null>(null);
  const [customDistEnd, setCustomDistEnd] = useState<string | null>(null);
  const [profileFilter, setProfileFilter] = useState("all");

  useEffect(() => {
    try {
      const savedFlow = localStorage.getItem("investowl_brokerflow_flow_mode");
      if (savedFlow) setFlowMode(savedFlow);
      const savedDist = localStorage.getItem("investowl_brokerflow_dist_mode");
      if (savedDist) setDistMode(savedDist);
      const savedProfile = localStorage.getItem("investowl_brokerflow_profile_filter");
      if (savedProfile) setProfileFilter(savedProfile);
    } catch {}
  }, []);

  const handleSetFlowMode = (m: string) => {
    setFlowMode(m);
    try {
      localStorage.setItem("investowl_brokerflow_flow_mode", m);
    } catch {}
  };

  const handleSetDistMode = (m: string) => {
    setDistMode(m);
    try {
      localStorage.setItem("investowl_brokerflow_dist_mode", m);
    } catch {}
  };

  const handleSetProfileFilter = (p: string) => {
    setProfileFilter(p);
    try {
      localStorage.setItem("investowl_brokerflow_profile_filter", p);
    } catch {}
  };

  const defaultDistEnd = analysisDate;
  const defaultDistStart = useMemo(() => {
    if (!analysisDate) return "";
    const dateObj = new Date(analysisDate);
    if (isNaN(dateObj.getTime())) return analysisDate;
    dateObj.setDate(dateObj.getDate() - windowDays);
    const yyyy = dateObj.getFullYear();
    const mm = String(dateObj.getMonth() + 1).padStart(2, "0");
    const dd = String(dateObj.getDate()).padStart(2, "0");
    return `${yyyy}-${mm}-${dd}`;
  }, [analysisDate, windowDays]);

  const distDate = customDistDate ?? analysisDate;
  const distStart = customDistStart ?? defaultDistStart;
  const distEnd = customDistEnd ?? defaultDistEnd;

  let qs = `?window_days=${windowDays}`;
  if (analysisDate) qs += `&analysis_date=${analysisDate}`;
  if (flowMode) qs += `&flow_mode=${flowMode}`;
  if (distMode) qs += `&dist_mode=${encodeURIComponent(distMode)}`;
  if (distDate) qs += `&dist_date=${distDate}`;
  if (distStart) qs += `&dist_start=${distStart}`;
  if (distEnd) qs += `&dist_end=${distEnd}`;
  if (selectedCodes.length > 0) qs += `&broker_codes=${selectedCodes.join(",")}`;

  const { data, error, isLoading, isValidating } = useSWR<BrokerFlowData>(
    ticker ? `/api/bandar/broker-flow/${ticker}${qs}` : null,
    fetcher,
    { refreshInterval: 60000, revalidateOnFocus: false, keepPreviousData: true }
  );

  const activeCodes: string[] = selectedCodes.length > 0 ? selectedCodes : data?.default_codes?.slice(0, 3) || [];

  const toggleCode = (code: string) => {
    const current = selectedCodes.length > 0 ? selectedCodes : data?.default_codes?.slice(0, 3) || [];
    if (current.includes(code)) {
      setSelectedCodes(current.filter((c) => c !== code));
    } else {
      if (current.length >= 10) return;
      setSelectedCodes([...current, code]);
    }
  };

  if (isLoading && !data) {
    return (
      <div className="p-8 border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-low)] rounded-[24px] flex items-center justify-center gap-3 text-[var(--md-sys-color-on-surface-variant)] shadow-sm animate-pulse">
        <Icon icon="ph:spinner-gap-bold" className="animate-spin text-[var(--md-sys-color-primary)]" width="24" />
        <span className="text-sm font-bold tracking-wide">Extracting broker flow data...</span>
      </div>
    );
  }
  if (error) {
    return (
      <div className="p-5 bg-rose-500/10 border border-rose-500/30 text-rose-500 font-bold text-sm rounded-[24px] shadow-sm">
        Error processing data.
      </div>
    );
  }

  const chartData = data?.compare_chart || [];
  const dist = data?.distribution || { edges: [] };
  const summary = data?.summary || [];
  const detailRows = data?.detail_rows || [];
  const profileFlow = data?.profile_flow || [];
  const profileDetail = data?.profile_broker_detail || [];
  const maxProfileAbs = Math.max(...profileFlow.map((r) => Math.abs(r.net)), 1);

  const filteredProfileDetail =
    profileFilter === "all" ? profileDetail : profileDetail.filter((r) => r.profile_key === profileFilter);

  return (
    <div
      className={`space-y-6 sm:space-y-8 transition-opacity duration-300 ${isValidating ? "opacity-70" : "opacity-100"} animate-fade-in`}
    >
      <div className="space-y-4">
        <div className="flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={() => {
              const retailPreset = ["YP", "PD", "XC", "XL", "KK"].filter((c) => (data?.all_codes || []).includes(c));
              setSelectedCodes(retailPreset.length > 0 ? retailPreset : ["YP", "PD", "XC"]);
            }}
            className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-bold bg-amber-500/10 text-amber-500 border border-amber-500/30 hover:bg-amber-500/20 transition-all cursor-pointer shadow-2xs"
            title="Tampilkan deret akumulasi/distribusi broker ritel populer"
          >
            <Icon icon="ph:users-three-bold" width="14" />
            <span>Preset Retail (YP, PD, XC, XL, KK)</span>
          </button>
          <button
            type="button"
            onClick={() => {
              const smartPreset = ["AK", "BK", "KZ", "CS", "RX", "ZP"].filter((c) => (data?.all_codes || []).includes(c));
              setSelectedCodes(smartPreset.length > 0 ? smartPreset : ["AK", "BK"]);
            }}
            className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-bold bg-emerald-500/10 text-emerald-500 border border-emerald-500/30 hover:bg-emerald-500/20 transition-all cursor-pointer shadow-2xs"
            title="Tampilkan deret akumulasi/distribusi broker institusi asing tier 1"
          >
            <Icon icon="ph:crown-bold" width="14" />
            <span>Preset Smart Foreign (AK, BK, KZ, CS)</span>
          </button>
          <button
            type="button"
            onClick={() => {
              const top3 = data?.default_codes?.slice(0, 3) || [];
              setSelectedCodes(top3);
            }}
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)] hover:text-[var(--md-sys-color-on-surface)] transition-all cursor-pointer shadow-2xs"
          >
            <Icon icon="ph:arrows-clockwise-bold" width="14" />
            <span>Reset Top 3</span>
          </button>
        </div>

        <div className="flex flex-wrap items-center gap-3">
          <select
            className="bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded-full px-4 py-2.5 text-sm font-bold text-[var(--md-sys-color-on-surface)] outline-none focus:border-[var(--md-sys-color-primary)] shadow-sm cursor-pointer"
            onChange={(e) => {
              if (e.target.value) toggleCode(e.target.value);
              e.target.value = "";
            }}
          >
            <option value="">+ Add Broker</option>
            {(data?.all_codes || []).map((code) => (
              <option key={code} value={code} disabled={activeCodes.includes(code)}>
                {code}
              </option>
            ))}
          </select>

          <div className="flex flex-wrap gap-2">
            {activeCodes.map((code) => (
              <div
                key={code}
                className="flex items-center gap-2 pl-3.5 pr-1.5 py-1.5 bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border border-[var(--md-sys-color-primary)] rounded-full text-sm font-extrabold tracking-wide shadow-sm"
              >
                {code}
                <button
                  type="button"
                  onClick={() => toggleCode(code)}
                  className="bg-[var(--md-sys-color-surface)]/20 hover:bg-[var(--color-negative)] hover:text-white transition-colors p-1 rounded-full"
                >
                  <Icon icon="ph:x-bold" width="14" />
                </button>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)] p-4 sm:p-5 rounded-[24px] shadow-sm">
          <label className="block text-[10px] font-extrabold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-widest mb-2.5">
            Flow mode
          </label>
          <select
            className="w-full md:w-72 bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] rounded-full px-4 py-2.5 text-sm font-bold text-[var(--md-sys-color-on-surface)] outline-none focus:border-[var(--md-sys-color-primary)] shadow-sm cursor-pointer"
            value={flowMode}
            onChange={(e) => handleSetFlowMode(e.target.value)}
          >
            <option value="Cumulative">Cumulative</option>
            <option value="Daily">Daily</option>
          </select>
          <p className="text-xs font-medium text-[var(--md-sys-color-on-surface-variant)] mt-3">
            Cumulative mode sums broker net flow across the selected broker window. Daily mode shows each date separately.
          </p>
        </div>
      </div>

      {activeCodes.length > 0 && (
        <div className="bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)] rounded-[28px] p-5 sm:p-6 shadow-sm transition-colors duration-300">
          <h3 className="text-sm sm:text-base font-extrabold text-[var(--md-sys-color-on-surface)] mb-5 tracking-tight">
            Broker Flow Comparison{" "}
            <span className="opacity-60 text-[var(--md-sys-color-on-surface-variant)] text-sm font-semibold ml-1">
              ({flowMode})
            </span>
          </h3>
          <div className="h-80 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={chartData} margin={{ top: 5, right: 5, left: -10, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="var(--md-sys-color-outline-variant)" strokeOpacity={0.5} vertical={false} />
                <XAxis
                  dataKey="date"
                  tick={{ fontSize: 10, fill: "var(--md-sys-color-on-surface-variant)", fontWeight: 600 }}
                  stroke="var(--md-sys-color-outline-variant)"
                  axisLine={false}
                  tickLine={false}
                  dy={10}
                />
                <YAxis
                  tick={{ fontSize: 10, fill: "var(--md-sys-color-on-surface-variant)", fontWeight: 600 }}
                  stroke="var(--md-sys-color-outline-variant)"
                  axisLine={false}
                  tickLine={false}
                  label={{
                    value: "Cumulative Net (Rp B)",
                    angle: -90,
                    position: "insideLeft",
                    fill: "var(--md-sys-color-on-surface-variant)",
                    fontSize: 10,
                    fontWeight: 700,
                  }}
                />
                <Tooltip
                  contentStyle={{
                    background: "var(--md-sys-color-surface-container-highest)",
                    border: "1px solid var(--md-sys-color-outline-variant)",
                    borderRadius: "16px",
                    fontSize: "12px",
                    color: "var(--md-sys-color-on-surface)",
                    fontWeight: 700,
                    boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
                  }}
                  labelStyle={{ color: "var(--md-sys-color-on-surface-variant)", fontWeight: 800, marginBottom: "4px" }}
                  formatter={(value, name) => [`Rp ${Number(value ?? 0).toFixed(2)} B`, String(name ?? "")]}
                />
                <ReferenceLine y={0} stroke="var(--md-sys-color-outline)" strokeWidth={1.5} />
                {activeCodes.map((code, i) => (
                  <Line
                    key={code}
                    type="monotone"
                    dataKey={code}
                    stroke={COLORS[i % COLORS.length]}
                    strokeWidth={3}
                    dot={{ r: 3, strokeWidth: 1 }}
                    activeDot={{ r: 6, stroke: "var(--md-sys-color-surface)", strokeWidth: 2 }}
                    connectNulls
                  />
                ))}
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      )}

      <div className="bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)] rounded-[28px] p-5 sm:p-6 shadow-sm transition-colors duration-300">
        <div className="flex items-center gap-2 mb-5 border-b border-[var(--md-sys-color-outline-variant)] pb-3">
          <Icon icon="ph:link-bold" className="text-[var(--md-sys-color-primary)]" width="20" />
          <h3 className="text-sm sm:text-base font-extrabold text-[var(--md-sys-color-on-surface)] tracking-tight">
            Broker Profile Flow
          </h3>
        </div>

        {profileFlow.length === 0 ? (
          <p className="text-xs font-bold text-[var(--md-sys-color-on-surface-variant)] text-center py-6">
            No profile flow detected.
          </p>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-5">
            {profileFlow.map((row, i) => {
              const width = Math.max(3, (Math.abs(row.net) / maxProfileAbs) * 100);
              return (
                <div
                  key={i}
                  className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-[20px] p-5 flex flex-col justify-between transition-colors shadow-sm"
                >
                  <div className="flex justify-between items-start mb-4">
                    <div>
                      <div
                        className="text-sm sm:text-base font-extrabold tracking-tight uppercase"
                        style={{ color: getProfileTextColor(row.label) }}
                      >
                        {row.label}
                      </div>
                      <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] mt-1">
                        {row.description}
                      </div>
                    </div>
                    <span className="font-mono font-extrabold text-sm sm:text-base" style={{ color: signedColor(row.net) }}>
                      {fmtRp(row.net)}
                    </span>
                  </div>
                  <div>
                    <div className="h-2 bg-[var(--md-sys-color-surface-container-highest)] rounded-full overflow-hidden border border-[var(--md-sys-color-outline-variant)] mb-4">
                      <div
                        className="h-full rounded-full transition-all duration-1000"
                        style={{ width: `${width}%`, backgroundColor: signedColor(row.net) }}
                      />
                    </div>
                    <div className="flex flex-wrap gap-2">
                      {(row.top_brokers || []).map((b, j) => (
                        <span
                          key={j}
                          className="inline-flex items-center gap-1.5 text-[10px] bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)] rounded-full px-3 py-1.5 shadow-sm"
                        >
                          <span className="font-mono font-extrabold text-[var(--md-sys-color-on-surface)]">{b.broker_code}</span>
                          <span className="text-[var(--md-sys-color-on-surface-variant)] font-bold uppercase">
                            {b.participant_type === "Asing" ? "FOREIGN" : b.participant_type === "Lokal" ? "LOCAL" : b.participant_type}
                          </span>
                          <span className="font-mono font-extrabold" style={{ color: signedColor(b.net) }}>
                            {fmtRp(b.net)}
                          </span>
                        </span>
                      ))}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      <BrokerFlowDetails
        profileFilter={profileFilter}
        onProfileFilterChange={handleSetProfileFilter}
        filteredProfileDetail={filteredProfileDetail}
        isValidating={isValidating}
        distMode={distMode}
        onDistModeChange={handleSetDistMode}
        distDate={distDate}
        distStart={distStart}
        distEnd={distEnd}
        onDistDateChange={setCustomDistDate}
        onDistStartChange={setCustomDistStart}
        onDistEndChange={setCustomDistEnd}
        edges={dist.edges || []}
        summary={summary}
        detailRows={detailRows}
      />
    </div>
  );
}
