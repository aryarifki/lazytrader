import { useMemo } from "react";
import ReactECharts from "echarts-for-react";
import { buildIhsgChartOption } from "./foreignProIhsgChartOption";
import type { IhsgMacroMetrics, IhsgResponse } from "./foreignProTypes";

interface ForeignProIhsgViewProps {
  ihsgData: IhsgResponse | null;
  ihsgMacroMetrics: IhsgMacroMetrics;
  isLightMode: boolean;
  catIhsg: boolean;
  setCatIhsg: (v: boolean) => void;
  catKonglo: boolean;
  setCatKonglo: (v: boolean) => void;
  catFund: boolean;
  setCatFund: (v: boolean) => void;
  metFF: boolean;
  setMetFF: (v: boolean) => void;
  metFCo: boolean;
  setMetFCo: (v: boolean) => void;
  metKF: boolean;
  setMetKF: (v: boolean) => void;
  metXL: boolean;
  setMetXL: (v: boolean) => void;
  metSR: boolean;
  setMetSR: (v: boolean) => void;
  flipKF: boolean;
  setFlipKF: (v: boolean) => void;
  flipXL: boolean;
  setFlipXL: (v: boolean) => void;
  flipSR: boolean;
  setFlipSR: (v: boolean) => void;
}

export default function ForeignProIhsgView(props: ForeignProIhsgViewProps) {
  const {
    ihsgData,
    ihsgMacroMetrics,
    isLightMode,
    catIhsg,
    setCatIhsg,
    catKonglo,
    setCatKonglo,
    catFund,
    setCatFund,
    metFF,
    setMetFF,
    metFCo,
    setMetFCo,
    metKF,
    setMetKF,
    metXL,
    setMetXL,
    metSR,
    setMetSR,
    flipKF,
    setFlipKF,
    flipXL,
    setFlipXL,
    flipSR,
    setFlipSR,
  } = props;

  const chartOption = useMemo(
    () =>
      buildIhsgChartOption({
        ihsgData,
        isLightMode,
        catIhsg,
        catKonglo,
        catFund,
        metFF,
        metFCo,
        metKF,
        metXL,
        metSR,
        flipKF,
        flipXL,
        flipSR,
      }),
    [
      ihsgData,
      isLightMode,
      catIhsg,
      catKonglo,
      catFund,
      metFF,
      metFCo,
      metKF,
      metXL,
      metSR,
      flipKF,
      flipXL,
      flipSR,
    ]
  );

  const changePct = ihsgData?.summary?.changePct || 0;

  return (
    <div className="space-y-3">
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-xl p-3 shadow-md space-y-3 font-mono">
        <div className="flex items-center justify-between pb-2 border-b border-[var(--md-sys-color-outline-variant)]/40 flex-wrap gap-2">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="w-3 h-3 rounded-full bg-[var(--md-sys-color-primary)]" />
            <span className="text-base font-black text-[var(--md-sys-color-on-surface)]">
              IHSG COMPOSITE MACRO VIEW — GRAFIK COMPOSITE FLOW IHSG
            </span>
            <span
              id="dateRangeInfo"
              className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] ml-2"
            >
              {ihsgData?.points && ihsgData.points.length > 0
                ? `${ihsgData.points[Math.max(0, ihsgData.points.length - 120)].date} – ${ihsgData.points[ihsgData.points.length - 1].date} (120 hari)`
                : "12 Mei 2026 – 30 Okt 2026 (120 hari)"}
            </span>
          </div>

          <div className="text-right">
            <div className="text-sm font-bold text-[var(--md-sys-color-on-surface)]">
              Level: {ihsgData?.summary?.close || 5057.05}
            </div>
            <div
              className={`text-[10px] font-bold ${
                changePct >= 0 ? "text-[var(--color-positive)]" : "text-[var(--color-negative)]"
              }`}
            >
              {changePct > 0 ? `+${changePct}%` : `${changePct}%`}
            </div>
          </div>
        </div>

        {/* Controls Strip: Categories & Metrics */}
        <div className="grid grid-cols-1 md:grid-cols-12 gap-3 text-xs">
          <div className="md:col-span-4 bg-[var(--md-sys-color-surface-container-high)] p-2.5 rounded border border-[var(--md-sys-color-outline-variant)]/40 space-y-1.5">
            <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
              Stock Universe Category
            </div>
            <div className="flex items-center gap-4 flex-wrap">
              {[
                { id: "chkCatIhsg", label: "IHSG (Solid)", checked: catIhsg, set: setCatIhsg },
                { id: "chkCatKonglo", label: "Konglo (Dashed)", checked: catKonglo, set: setCatKonglo },
                { id: "chkCatFund", label: "Fundamental (Dotted)", checked: catFund, set: setCatFund },
              ].map((item) => (
                <label key={item.id} className="flex items-center gap-1.5 cursor-pointer">
                  <input
                    type="checkbox"
                    id={item.id}
                    checked={item.checked}
                    onChange={(e) => item.set(e.target.checked)}
                    className="accent-[var(--md-sys-color-primary)]"
                  />
                  <span className="text-[var(--md-sys-color-on-surface)]">{item.label}</span>
                </label>
              ))}
            </div>
          </div>

          <div className="md:col-span-8 bg-[var(--md-sys-color-surface-container-high)] p-2.5 rounded border border-[var(--md-sys-color-outline-variant)]/40 space-y-1.5">
            <div className="text-[10px] font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase">
              Metrics &amp; Cumulative Flow (Triliun Rp)
            </div>
            <div className="flex items-center justify-between flex-wrap gap-2 text-xs">
              <label className="flex items-center gap-1.5 cursor-pointer">
                <input
                  type="checkbox"
                  id="chkMetFF"
                  checked={metFF}
                  onChange={(e) => setMetFF(e.target.checked)}
                  className="accent-[#FF6820]"
                />
                <span className="text-[#FF6820] font-bold">FF</span>
                <span id="subFF" className="text-[10px] text-[var(--md-sys-color-on-surface)]">
                  {ihsgMacroMetrics.ff}
                </span>
              </label>

              <label className="flex items-center gap-1.5 cursor-pointer">
                <input
                  type="checkbox"
                  id="chkMetFCoV2"
                  checked={metFCo}
                  onChange={(e) => setMetFCo(e.target.checked)}
                  className="accent-[#D62839]"
                />
                <span className="text-[#D62839] font-bold">F&amp;Co</span>
                <span id="subFCoV2" className="text-[10px] text-[var(--md-sys-color-on-surface)]">
                  {ihsgMacroMetrics.fco}
                </span>
              </label>

              {[
                {
                  chkId: "chkMetKF",
                  flipId: "flipKF",
                  subId: "subKF",
                  label: "KF",
                  color: "#2F7676",
                  checked: metKF,
                  setChecked: setMetKF,
                  val: ihsgMacroMetrics.kf,
                  flipped: flipKF,
                  setFlipped: setFlipKF,
                },
                {
                  chkId: "chkMetXL",
                  flipId: "flipXL",
                  subId: "subXL",
                  label: "XL",
                  color: "#3A7BD5",
                  checked: metXL,
                  setChecked: setMetXL,
                  val: ihsgMacroMetrics.xl,
                  flipped: flipXL,
                  setFlipped: setFlipXL,
                },
                {
                  chkId: "chkMetSemiRitel",
                  flipId: "flipSemiRitel",
                  subId: "subSemiRitel",
                  label: "SR",
                  color: "#1DB87A",
                  checked: metSR,
                  setChecked: setMetSR,
                  val: ihsgMacroMetrics.sr,
                  flipped: flipSR,
                  setFlipped: setFlipSR,
                },
              ].map((m) => (
                <div key={m.chkId} className="flex items-center gap-1">
                  <label className="flex items-center gap-1.5 cursor-pointer">
                    <input
                      type="checkbox"
                      id={m.chkId}
                      checked={m.checked}
                      onChange={(e) => m.setChecked(e.target.checked)}
                      style={{ accentColor: m.color }}
                    />
                    <span className="font-bold" style={{ color: m.color }}>
                      {m.label}
                    </span>
                    <span
                      id={m.subId}
                      className="text-[10px] text-[var(--md-sys-color-on-surface)]"
                    >
                      {m.val}
                    </span>
                  </label>
                  <button
                    id={m.flipId}
                    type="button"
                    onClick={() => m.setFlipped(!m.flipped)}
                    className={`px-1.5 py-0.5 rounded text-[10px] font-bold transition cursor-pointer ${
                      m.flipped
                        ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-xs"
                        : "bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface)] hover:bg-[var(--md-sys-color-primary)] hover:text-[var(--md-sys-color-on-primary)]"
                    }`}
                    title="Invert sign"
                  >
                    ⇅
                  </button>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Macro Chart */}
        <div className="h-[520px] w-full">
          <ReactECharts
            option={chartOption}
            style={{ height: "100%", width: "100%" }}
            notMerge={true}
          />
        </div>

        {/* Highlights Box for IHSG */}
        <div
          id="hlBox"
          className="bg-[var(--md-sys-color-surface-container-high)] p-3 rounded border border-[var(--md-sys-color-outline-variant)]/40 text-xs space-y-1"
        >
          <div className="flex items-center gap-2 border-b border-[var(--md-sys-color-outline-variant)]/40 pb-1">
            <span className="w-2 h-2 rounded-full bg-[var(--md-sys-color-primary)]" />
            <span id="hlTitle" className="font-bold text-[var(--md-sys-color-on-surface)]">
              Ringkasan Foreign Flow Makro Bursa
            </span>
            <span
              id="hlBadge"
              className="ml-auto text-[10px] text-[var(--md-sys-color-on-primary-container)] bg-[var(--md-sys-color-primary-container)] px-1.5 py-0.5 rounded font-bold"
            >
              EOD UPDATE
            </span>
          </div>
          <div
            id="hlList"
            className="space-y-1 text-[var(--md-sys-color-on-surface-variant)] text-[11px] pt-1"
          >
            <div>
              • Penutupan IHSG berada di level {ihsgData?.summary?.close || 5057.05} (
              {changePct > 0 ? `+${changePct}%` : `${changePct || "-0.80"}%`}) dengan transaksi Rp 11.64 Triliun.
            </div>
            <div>
              • Net Foreign {(ihsgData?.summary?.ffValue ?? -176.08) >= 0 ? "Buy" : "Sell"} bursa hari ini sebesar Rp{" "}
              {Math.abs(ihsgData?.summary?.ffValue ?? 176.08).toFixed(2)} Miliar.
            </div>
            <div>
              • Akumulasi asing bursa 120 hari terakhir tercatat {ihsgMacroMetrics.ff} ({ihsgData?.summary?.status || "EOD"}).
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
