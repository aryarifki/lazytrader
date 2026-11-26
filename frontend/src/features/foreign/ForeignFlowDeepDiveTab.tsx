import { useMemo } from "react";
import useSWR from "swr";
import { Icon } from "@iconify/react";
import { useAppStore } from "@/store/useAppStore";
import { useForeignFlowAnalytics } from "@/hooks/useForeignFlow";
import DeepDiveHeader from "./DeepDiveHeader";
import DeepDiveKpiCards from "./DeepDiveKpiCards";
import CostBasisSection from "./CostBasisSection";
import CounterpartySection from "./CounterpartySection";
import KseiT5Section from "./KseiT5Section";
import type { ForeignFlowAnalyticsResponse, UniverseAllResponse } from "./deepdiveTypes";

const fetcher = (url: string) =>
  fetch(url, { credentials: "include" }).then((r) => r.json());

export default function ForeignFlowDeepDiveTab() {
  const { activeTicker, setActiveTicker, windowDays, setWindowDays } = useAppStore();
  const { data: rawData, isLoading, isFetching, isError } = useForeignFlowAnalytics(activeTicker, windowDays);
  const data = rawData as ForeignFlowAnalyticsResponse | undefined;

  const { data: allUniverseData, isLoading: isLoadingUniverse } = useSWR<UniverseAllResponse>(
    "/api/bandar/universe/all",
    fetcher
  );
  const allTickers = useMemo(() => allUniverseData?.tickers || [], [allUniverseData?.tickers]);

  const companyName = data?.company?.name || "Perusahaan Terdaftar IDX";
  const groupName = data?.company?.group || "Independen / Non-Afiliasi";

  const cb = data?.cost_basis;
  const imp = data?.impact_divergence;
  const cp = data?.counterparty;
  const ksei = data?.ksei_reality;

  const hhiValue = data?.features?.foreign_hhi || 0;
  const zScoreValue = data?.features?.foreign_zscore || 0;
  const pnlPct = Number(cb?.floating_pnl_pct || 0);
  const pnlTone = (cb?.total_net_value || 0) < 0 ? "negative" : pnlPct >= 0 ? "positive" : "warning";
  const streakDays = Number(imp?.streak_days || 0);

  return (
    <div className="flex flex-col gap-5 sm:gap-6 pb-20 transition-colors duration-300">
      {/* KONTROL EMITEN & INTELLIGENCE BAR */}
      <DeepDiveHeader
        activeTicker={activeTicker}
        setActiveTicker={setActiveTicker}
        windowDays={windowDays}
        setWindowDays={setWindowDays}
        companyName={companyName}
        groupName={groupName}
        latestDate={data?.latest_date}
        isFetching={isFetching}
        isLoading={isLoading}
        allTickers={allTickers}
        isLoadingUniverse={isLoadingUniverse}
      />

      {/* LOADING STATE */}
      {isLoading && (
        <div className="flex flex-col items-center justify-center py-20 text-[var(--md-sys-color-on-surface-variant)] bg-[var(--md-sys-color-surface-container)] rounded-2xl border border-[var(--md-sys-color-outline-variant)] shadow-xs gap-3">
          <Icon icon="ph:spinner-gap-bold" className="animate-spin text-[var(--md-sys-color-primary)]" width="28" />
          <span className="font-bold tracking-wider uppercase text-xs">
            Memuat Foreign Smart Money Forensic ({activeTicker})...
          </span>
        </div>
      )}

      {/* ERROR STATE */}
      {isError && !isLoading && (
        <div className="p-6 text-center bg-[var(--md-sys-color-surface-container-low)] rounded-2xl border border-[var(--md-sys-color-outline-variant)] space-y-3">
          <Icon icon="ph:warning-circle-duotone" width="36" className="mx-auto text-amber-500" />
          <h3 className="text-sm font-bold text-[var(--md-sys-color-on-surface)]">
            Data Forensik {activeTicker} Belum Mencukupi
          </h3>
          <p className="text-xs text-[var(--md-sys-color-on-surface-variant)] max-w-md mx-auto">
            Histori transaksi asing untuk {activeTicker} pada window {windowDays} hari belum mencukupi untuk analisis forensik.
          </p>
          <div className="pt-2 flex justify-center gap-2">
            {["BBCA", "BBRI", "BMRI"].map((fallback) => (
              <button
                key={fallback}
                onClick={() => setActiveTicker(fallback)}
                className="px-3 py-1 rounded-lg bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] text-xs font-bold shadow-xs hover:opacity-90"
              >
                Coba {fallback}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* DASHBOARD KONTEN: 4 MODUL FOREIGN SMART MONEY FORENSIC */}
      {!isLoading && !isError && data && (
        <>
          <DeepDiveKpiCards
            cb={cb}
            imp={imp}
            cp={cp}
            ksei={ksei}
            pnlPct={pnlPct}
            pnlTone={pnlTone}
            streakDays={streakDays}
            zScoreValue={zScoreValue}
          />

          <CostBasisSection
            windowDays={windowDays}
            cb={cb}
            pnlPct={pnlPct}
            timeSeries={data.timeseries}
          />

          <CounterpartySection
            cp={cp}
            imp={imp}
            hhiValue={hhiValue}
            brokerHeatmap={data.models?.broker_heatmap}
          />

          <KseiT5Section activeTicker={activeTicker} ksei={ksei} />
        </>
      )}
    </div>
  );
}
