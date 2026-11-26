import { Icon } from "@iconify/react";
import { fmtRp, fmtPct, signedColor, signalColor } from "@/lib/format";
import type { TickerDetailData } from "./types";

interface TickerHeaderProps {
  activeTicker: string;
  windowDays: number;
  isBookmarked: boolean;
  onToggleBookmark: () => void;
  data?: TickerDetailData;
  isLoading: boolean;
  error: unknown;
}

function MetricCard({
  label,
  value,
  note,
  tone,
  accent,
}: {
  label: string;
  value: string;
  note: string;
  tone: number | null | undefined;
  accent?: string;
}) {
  let color = "var(--md-sys-color-on-surface)";
  if (accent) color = accent;
  else if (tone !== null && tone !== undefined) color = signedColor(Number(tone));

  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] hover:border-[var(--md-sys-color-primary)] rounded-xl p-3.5 transition-all duration-200 shadow-2xs relative overflow-hidden group">
      <div
        className="absolute top-0 left-0 bottom-0 w-1 transition-all opacity-80 group-hover:opacity-100"
        style={{ backgroundColor: color }}
      ></div>
      <div className="pl-2">
        <div className="text-[10px] font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] uppercase tracking-wider mb-1 truncate">
          {label}
        </div>
        <div className="text-base sm:text-lg font-mono font-extrabold tracking-tight truncate" style={{ color }}>
          {value}
        </div>
        <div className="text-[10px] font-mono font-medium text-[var(--md-sys-color-on-surface-variant)] truncate mt-1">
          {note}
        </div>
      </div>
    </div>
  );
}

export default function TickerHeader({
  activeTicker,
  windowDays,
  isBookmarked,
  onToggleBookmark,
  data,
  isLoading,
  error,
}: TickerHeaderProps) {
  return (
    <>
      <div className="mt-2 mb-6 bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-3xl p-5 sm:p-6 shadow-xs relative overflow-hidden transition-all">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-3">
          <div>
            <div className="inline-flex items-center gap-2 text-[10px] font-mono font-bold text-[var(--md-sys-color-primary)] uppercase tracking-wider mb-1">
              <span>BANDARMOLOGY &amp; SMART MONEY DEEP-DIVE</span>
              <span className="opacity-40">•</span>
              <span className="text-[var(--color-positive)] font-bold">LOOKBACK: {windowDays} HARI</span>
            </div>
            <div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-3">
              <div className="flex items-center gap-2.5">
                <h1 className="text-3xl sm:text-4xl font-mono font-extrabold text-[var(--md-sys-color-on-surface)] tracking-tight">
                  {activeTicker}
                </h1>
                <button
                  type="button"
                  onClick={onToggleBookmark}
                  className={`p-1.5 sm:p-2 rounded-full border transition-all cursor-pointer shadow-2xs active:scale-95 flex items-center justify-center ${
                    isBookmarked
                      ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-primary)] border-[var(--md-sys-color-primary)]"
                      : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)]/60 hover:text-[var(--md-sys-color-primary)] hover:border-[var(--md-sys-color-primary)]"
                  }`}
                  title={isBookmarked ? `Hapus ${activeTicker} dari Watchlist` : `Simpan ${activeTicker} ke Watchlist`}
                  aria-label={isBookmarked ? `Hapus ${activeTicker} dari Watchlist` : `Simpan ${activeTicker} ke Watchlist`}
                >
                  <Icon
                    icon={isBookmarked ? "ph:star-fill" : "ph:star-bold"}
                    width="20"
                    className={isBookmarked ? "text-amber-500" : ""}
                  />
                </button>
              </div>
              <span className="text-xs sm:text-sm font-semibold text-[var(--md-sys-color-on-surface-variant)]">
                {data?.company?.name || "Bursa Efek Indonesia"}
              </span>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <span className="text-xs font-mono font-bold bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-xl px-3 py-1.5 shadow-2xs">
              Periode: {data?.window_start || "..."} <span className="opacity-50 mx-1">s/d</span>{" "}
              {data?.analysis_date || "..."}
            </span>
            <span className="text-xs font-semibold px-3 py-1.5 rounded-xl bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border border-[var(--md-sys-color-outline-variant)]/60">
              {data?.company?.group || "Equities IDX"}
            </span>
          </div>
        </div>
      </div>

      {isLoading && (
        <div className="flex items-center gap-2 text-[var(--md-sys-color-on-surface-variant)] text-xs font-mono font-medium mb-4 animate-pulse">
          <Icon icon="ph:spinner-gap-duotone" className="animate-spin" width="16" />
          <span>Memuat data analitik presisi {activeTicker}...</span>
        </div>
      )}

      {Boolean(error) && (
        <div className="p-4 bg-rose-500/10 border border-rose-500/20 text-rose-600 dark:text-rose-400 text-xs font-mono font-bold rounded-xl mb-4">
          Gagal mengambil data sistem bandarmologi. Pastikan backend aktif.
        </div>
      )}

      {data && !data.error && (
        <div className="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-3 mb-6">
          <MetricCard
            label="Conviction Score"
            value={data.conviction_score != null ? `${data.conviction_score.toFixed(1)}/100` : "-"}
            note="weighted composite"
            tone={data.conviction_score}
          />
          <MetricCard
            label="Bandar Signal"
            value={data.signal || "NEUTRAL"}
            note="selected date"
            tone={null}
            accent={signalColor(data.signal_score)}
          />
          <MetricCard label="5D Return" value={fmtPct(data.ret_5d)} note="price context" tone={data.ret_5d} />
          <MetricCard label="Foreign Net 5D" value={fmtRp(data.foreign_5d)} note="broker summary" tone={data.foreign_5d} />
          <MetricCard
            label="Top Buyer"
            value={data.top_buyer?.broker || "-"}
            note={fmtRp(data.top_buyer?.net)}
            tone={1}
          />
          <MetricCard
            label="Smart Cumulative"
            value={fmtRp(data.smart_cumulative)}
            note={`${data.smart_daily?.length || 0} broker days`}
            tone={data.smart_cumulative}
          />
        </div>
      )}

      {data?.alerts && data.alerts.length > 0 && (
        <div className="mb-6 bg-rose-500/10 border border-rose-500/20 rounded-xl px-4 py-3 shadow-2xs">
          <div className="flex items-center gap-2 mb-1.5 text-rose-500 text-xs font-bold uppercase tracking-wider">
            <Icon icon="ph:warning-circle-bold" width="16" />
            <span>Deteksi Anomali / Kontradiksi Arus Modal</span>
          </div>
          {data.alerts.map((a, i) => (
            <div key={i} className="text-xs font-medium text-rose-400 flex items-start gap-2 mb-1 last:mb-0">
              <span className="w-1.5 h-1.5 rounded-full bg-rose-500 mt-1.5 flex-shrink-0"></span>
              <span>{a}</span>
            </div>
          ))}
        </div>
      )}

      {data?.verdict && (
        <div className="mb-6 bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-xl p-4 sm:p-5 shadow-2xs relative overflow-hidden">
          <div className="absolute left-0 top-0 bottom-0 w-1.5 bg-[var(--md-sys-color-primary)]"></div>
          <div className="pl-2">
            <div className="text-[10px] font-mono font-bold text-[var(--md-sys-color-primary)] uppercase tracking-wider mb-1">
              Current Quantitative Read
            </div>
            <div className="text-xs sm:text-sm text-[var(--md-sys-color-on-surface)] leading-relaxed font-medium">
              {data.verdict}
            </div>
          </div>
        </div>
      )}

      {data?.affiliated_brokers_detected && data.affiliated_brokers_detected.length > 0 && (
        <div className="mb-6 bg-amber-500/10 border border-amber-500/25 rounded-2xl p-4 sm:p-5 shadow-xs relative overflow-hidden transition-all">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-3">
            <div className="flex items-center gap-2 text-amber-500 dark:text-amber-400 text-xs font-bold uppercase tracking-wider">
              <Icon icon="ph:buildings-bold" width="18" />
              <span>Deteksi Relasi Afiliasi Broker &amp; Pengendali (Governance Breadcrumb)</span>
            </div>
            <span className="text-[10px] font-mono font-bold bg-amber-500/20 text-amber-600 dark:text-amber-300 px-2.5 py-0.5 rounded-full border border-amber-500/30 w-fit">
              {data.affiliated_brokers_detected.length} Broker Terafiliasi
            </span>
          </div>

          <div className="space-y-3">
            {data.affiliated_brokers_detected.map((aff, i) => (
              <div
                key={i}
                className="bg-[var(--md-sys-color-surface-container)] border border-amber-500/20 rounded-xl p-3.5 shadow-2xs"
              >
                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-2">
                  <div className="flex items-center gap-2">
                    <span className="px-2 py-0.5 bg-amber-500 text-black font-mono font-extrabold text-xs rounded-md shadow-2xs">
                      {aff.broker_code}
                    </span>
                    <span className="text-xs sm:text-sm font-extrabold text-[var(--md-sys-color-on-surface)]">
                      {aff.broker_name}
                    </span>
                  </div>
                  <span className="text-[10px] font-mono font-bold px-2.5 py-0.5 rounded-full bg-amber-500/15 text-amber-500 border border-amber-500/30 w-fit">
                    {aff.badge}
                  </span>
                </div>

                <p className="text-xs text-[var(--md-sys-color-on-surface-variant)] leading-relaxed mb-2.5 font-medium">
                  {aff.relation}
                </p>

                <div className="flex flex-wrap items-center gap-2 sm:gap-3 text-[11px] font-mono text-[var(--md-sys-color-on-surface-variant)] pt-2 border-t border-[var(--md-sys-color-outline-variant)]/40">
                  <span>
                    Aktif di Jendela: <strong className="text-[var(--md-sys-color-on-surface)]">{aff.days_active} Hari</strong>
                  </span>
                  <span className="opacity-40">•</span>
                  <span>
                    Akumulasi Kumulatif:{" "}
                    <strong style={{ color: signedColor(aff.cumulative_net) }}>{fmtRp(aff.cumulative_net)}</strong>
                  </span>
                  <span className="opacity-40">•</span>
                  <span>
                    Transaksi Terakhir:{" "}
                    <strong className="text-[var(--md-sys-color-on-surface)]">
                      {aff.latest_date} ({fmtRp(aff.latest_net)})
                    </strong>
                  </span>
                </div>
              </div>
            ))}
          </div>

          <div className="mt-3 text-[10px] font-medium text-[var(--md-sys-color-on-surface-variant)]/80 italic flex items-start sm:items-center gap-1.5">
            <Icon icon="ph:info-bold" width="13" className="flex-shrink-0 mt-0.5 sm:mt-0" />
            <span>
              * Observational research hypothesis berdasarkan data keterbukaan publik (komposisi dewan &amp; pemegang saham). Kode broker mencerminkan anggota bursa pelaksana, bukan konfirmasi identitas nasabah akhir atau tuduhan insider trading.
            </span>
          </div>
        </div>
      )}
    </>
  );
}
