import { Icon } from "@iconify/react";
import { getToneColor } from "./deepdiveHelpers";
import type {
  CostBasisInfo,
  CounterpartyInfo,
  ImpactDivergenceInfo,
  KseiRealityInfo,
} from "./deepdiveTypes";

interface DeepDiveKpiCardsProps {
  cb?: CostBasisInfo;
  imp?: ImpactDivergenceInfo;
  cp?: CounterpartyInfo;
  ksei?: KseiRealityInfo;
  pnlPct: number;
  pnlTone: string;
  streakDays: number;
  zScoreValue: number;
}

export default function DeepDiveKpiCards({
  cb,
  imp,
  cp,
  ksei,
  pnlPct,
  pnlTone,
  streakDays,
  zScoreValue,
}: DeepDiveKpiCardsProps) {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3.5 sm:gap-4">
      {/* Card 1: Divergence & Streak */}
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 shadow-xs relative overflow-hidden flex flex-col justify-between">
        <div className="absolute top-0 left-0 w-1.5 h-full" style={{ backgroundColor: getToneColor(imp?.divergence_tone) }} />
        <div className="pl-2">
          <div className="flex items-center justify-between gap-1 text-[10px] font-semibold text-[var(--md-sys-color-on-surface-variant)] mb-1.5">
            <span className="uppercase tracking-wider">Radar Divergensi</span>
            <span
              className="px-2 py-0.5 rounded-full font-bold font-mono"
              style={{ backgroundColor: `${getToneColor(imp?.divergence_tone)}20`, color: getToneColor(imp?.divergence_tone) }}
            >
              {streakDays > 0 ? `+${streakDays}d Beli` : streakDays < 0 ? `${streakDays}d Jual` : "0d"}
            </span>
          </div>
          <div className="text-sm sm:text-base font-extrabold tracking-tight leading-snug" style={{ color: getToneColor(imp?.divergence_tone) }}>
            {imp?.divergence_label || "Netral"}
          </div>
          <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-1.5 leading-relaxed">
            {imp?.summary_text}
          </p>
          {imp?.crossing_warning && (
            <span className="inline-flex items-center gap-1 mt-2 text-[10px] font-bold px-2 py-0.5 rounded-md bg-amber-500/15 text-amber-500 border border-amber-500/30">
              <Icon icon="ph:warning-bold" width="11" />
              Dominasi 1 Hari ({imp.max_single_day_share_pct}% — indikasi crossing)
            </span>
          )}
        </div>
      </div>

      {/* Card 2: Foreign Avg Cost & Floating PnL */}
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 shadow-xs relative overflow-hidden flex flex-col justify-between">
        <div className="absolute top-0 left-0 w-1.5 h-full" style={{ backgroundColor: getToneColor(pnlTone) }} />
        <div className="pl-2">
          <div className="flex items-center justify-between gap-1 text-[10px] font-semibold text-[var(--md-sys-color-on-surface-variant)] mb-1.5">
            <span className="uppercase tracking-wider">Modal Asing (VWAP)</span>
            <span
              className="px-2 py-0.5 rounded-full font-bold font-mono"
              style={{ backgroundColor: `${getToneColor(pnlTone)}20`, color: getToneColor(pnlTone) }}
            >
              {pnlPct >= 0 ? `+${pnlPct.toFixed(2)}%` : `${pnlPct.toFixed(2)}%`}
            </span>
          </div>
          <div className="text-lg font-extrabold font-mono text-[var(--md-sys-color-on-surface)]">
            Rp {Number(cb?.avg_buy_price || 0).toLocaleString("id-ID", { maximumFractionDigits: 0 })}
          </div>
          <div className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-1.5 flex items-center justify-between">
            <span>{cb?.position_status}</span>
            <span className="font-mono font-bold text-[var(--md-sys-color-on-surface)]">
              Retensi {cb?.inventory_retention_pct ?? 0}%
            </span>
          </div>
          <div className="w-full h-1.5 rounded-full bg-[var(--md-sys-color-surface-container-high)] mt-2 overflow-hidden">
            <div
              className="h-full rounded-full bg-emerald-500 transition-all"
              style={{ width: `${Math.min(100, Math.max(0, cb?.inventory_retention_pct || 0))}%` }}
            />
          </div>
        </div>
      </div>

      {/* Card 3: Foreign Power & Impact */}
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 shadow-xs relative overflow-hidden flex flex-col justify-between">
        <div
          className="absolute top-0 left-0 w-1.5 h-full"
          style={{
            backgroundColor:
              imp?.foreign_power_badge === "HFP"
                ? "var(--color-positive)"
                : imp?.foreign_power_badge === "MFP"
                ? "#06b6d4"
                : "var(--md-sys-color-outline)",
          }}
        />
        <div className="pl-2">
          <div className="flex items-center justify-between gap-1 text-[10px] font-semibold text-[var(--md-sys-color-on-surface-variant)] mb-1.5">
            <span className="uppercase tracking-wider">Dampak &amp; Partisipasi</span>
            <span className="px-2 py-0.5 rounded-full font-bold bg-[var(--md-sys-color-primary)]/15 text-[var(--md-sys-color-primary)] font-mono">
              {imp?.foreign_power_badge || "LFP"}
              {imp?.foreign_power_score !== null && imp?.foreign_power_score !== undefined ? ` · ${imp.foreign_power_score}` : ""}
            </span>
          </div>
          <div className="text-lg font-extrabold font-mono text-[var(--md-sys-color-on-surface)]">
            {Number(imp?.avg_participation_pct || 0).toFixed(1)}%{" "}
            <span className="text-xs font-normal text-[var(--md-sys-color-on-surface-variant)]">porsi asing</span>
          </div>
          <div className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-1.5 flex items-center justify-between">
            <span>
              Korelasi Harga:{" "}
              <strong className="font-mono text-[var(--md-sys-color-on-surface)]">
                {Number(imp?.foreign_price_corr || 0).toFixed(2)}
              </strong>
            </span>
            <span>
              Z:{" "}
              <strong className="font-mono text-[var(--md-sys-color-on-surface)]">
                {zScoreValue > 0 ? `+${zScoreValue.toFixed(1)}` : zScoreValue.toFixed(1)}σ
              </strong>
            </span>
          </div>
        </div>
      </div>

      {/* Card 4: Counterparty & KSEI Reality */}
      <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 shadow-xs relative overflow-hidden flex flex-col justify-between">
        <div className="absolute top-0 left-0 w-1.5 h-full" style={{ backgroundColor: getToneColor(cp?.matchup_tone) }} />
        <div className="pl-2">
          <div className="flex items-center justify-between gap-1 text-[10px] font-semibold text-[var(--md-sys-color-on-surface-variant)] mb-1.5">
            <span className="uppercase tracking-wider">Lawan Transaksi &amp; KSEI</span>
            <span className="px-2 py-0.5 rounded-full font-bold font-mono bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface)]">
              KSEI {ksei?.latest_pct_asing !== null && ksei?.latest_pct_asing !== undefined ? `${ksei.latest_pct_asing}%` : "N/A"}
            </span>
          </div>
          <div className="text-xs sm:text-sm font-extrabold text-[var(--md-sys-color-on-surface)] leading-snug">
            {cp?.matchup_verdict}
          </div>
          <div className="text-[11px] text-[var(--md-sys-color-on-surface-variant)] mt-1.5 flex items-center justify-between">
            <span>{ksei?.alignment_status}</span>
            {ksei?.delta_1m_pp !== null && ksei?.delta_1m_pp !== undefined && (
              <span className={`font-mono font-bold ${ksei.delta_1m_pp >= 0 ? "text-emerald-500" : "text-rose-500"}`}>
                MoM {ksei.delta_1m_pp >= 0 ? `+${ksei.delta_1m_pp}` : ksei.delta_1m_pp} pp
              </span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
