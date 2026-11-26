import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { Icon } from "@iconify/react";
import { getIpoHunterAnalysis, getIpoHunterMeta } from "@/lib/api/terminal";
import type {
  IpoAnalysisData,
  IpoAnalysisMode,
  IpoApiResponse,
  IpoMainView,
  IpoMeta,
  IpoRightTab,
} from "./ipoTypes";
import IpoAnalysisControls from "./IpoAnalysisControls";
import IpoCandlestickChart from "./IpoCandlestickChart";
import { AfterIpoDistCard, BuyerSellerDonutCard } from "./DonutCharts";
import IpoTopBuyerSellerCard from "./IpoTopBuyerSellerCard";
import IpoInfoTab from "./IpoInfoTab";
import IpoUnderwriterCard from "./IpoUnderwriterCard";
import IpoCategoryCard from "./IpoCategoryCard";
import IpoRetailActionCard from "./IpoRetailActionCard";
import IpoListView from "./IpoListView";
import IpoScreenerView from "./IpoScreenerView";

export default function IpoHunterPage() {
  const [searchParams] = useSearchParams();
  const qTicker = searchParams.get("ticker");

  const [mainView, setMainView] = useState<IpoMainView>("analisis");
  const [rightTab, setRightTab] = useState<IpoRightTab>("tx");
  const [meta, setMeta] = useState<IpoMeta | null>(null);

  const [selectedTicker, setSelectedTicker] = useState<string>(qTicker?.toUpperCase() || "");
  const [tickerInput, setTickerInput] = useState<string>(qTicker?.toUpperCase() || "");
  const [tickerError, setTickerError] = useState<string>("");

  const [mode, setMode] = useState<IpoAnalysisMode>("after");
  const [dateFrom, setDateFrom] = useState<string>("");
  const [dateTo, setDateTo] = useState<string>("");
  const [autoRetail, setAutoRetail] = useState<boolean>(true);
  const [retailOverride, setRetailOverride] = useState<string>("");

  const [analysis, setAnalysis] = useState<IpoAnalysisData | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  // Load meta once on mount
  useEffect(() => {
    let cancelled = false;
    getIpoHunterMeta<IpoApiResponse<IpoMeta>>()
      .then((res) => {
        if (cancelled) return;
        const d = res.data;
        setMeta(d);
        if (d?.quick_picks?.length) {
          const def = d.quick_picks[0].ticker;
          setSelectedTicker((prev) => prev || def);
          setTickerInput((prev) => prev || def);
        }
        if (d?.latest_date) {
          setDateTo((prev) => prev || d.latest_date || "");
        }
        if (d?.available_dates?.length) {
          const ds = d.available_dates;
          const defFrom = ds[Math.max(0, ds.length - 5)];
          setDateFrom((prev) => prev || defFrom);
        }
      })
      .catch((err) => {
        if (!cancelled) console.error("Error loading IPO Hunter meta:", err);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  // Fetch analysis when parameters change
  useEffect(() => {
    if (!selectedTicker || mainView !== "analisis") return;
    let cancelled = false;
    setLoading(true);
    setTickerError("");
    getIpoHunterAnalysis<IpoApiResponse<IpoAnalysisData>>(selectedTicker, {
      mode,
      date_from: mode === "periodic" ? dateFrom || undefined : undefined,
      date_to: dateTo || undefined,
      auto_retail: autoRetail,
      retail_override: retailOverride || undefined,
    })
      .then((res) => {
        if (cancelled) return;
        setAnalysis(res.data);
        setLoading(false);
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        const msg = err instanceof Error ? err.message : "Gagal memuat analisis IPO";
        setTickerError(msg);
        setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [selectedTicker, mode, dateFrom, dateTo, autoRetail, retailOverride, mainView]);

  const ipoTickerSet = useMemo(() => {
    const s = new Set<string>();
    (meta?.tickers || []).forEach((t) => s.add(t.ticker));
    return s;
  }, [meta]);

  const selectTicker = (tk: string) => {
    const up = tk.trim().toUpperCase();
    if (!up) return;
    if (ipoTickerSet.size > 0 && !ipoTickerSet.has(up)) {
      setTickerError(`Kode ${up} tidak ada di daftar IPO. Pilih kode dari IPO List.`);
      return;
    }
    setTickerError("");
    setSelectedTicker(up);
    setTickerInput(up);
    setAutoRetail(true);
    setRetailOverride("");
  };

  const cards = analysis?.cards;
  const infoCard = cards?.info_card;

  return (
    <div className="space-y-4">
      {/* Top Header & View Switcher */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2.5">
          <div className="w-9 h-9 rounded-xl bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] flex items-center justify-center border border-[var(--md-sys-color-primary)]/30">
            <Icon icon="ph:rocket-launch-duotone" width="20" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-lg sm:text-xl font-black tracking-tight">IPO Hunter Pro</h1>
              <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] border border-[var(--md-sys-color-primary)]/30 uppercase">
                {meta?.total_ipos ?? 317} Emiten
              </span>
            </div>
            <p className="text-[11px] text-[var(--md-sys-color-on-surface-variant)]">
              Analisis broker After-IPO (D0 + harian), track record underwriter, distribusi pasca-IPO, &amp; screener.
            </p>
          </div>
        </div>

        <div className="inline-flex rounded-xl bg-[var(--md-sys-color-surface-container)] p-1 border border-[var(--md-sys-color-outline-variant)]">
          {(
            [
              ["analisis", "Analisis", "ph:chart-line-up-duotone"],
              ["list", "IPO List", "ph:table-duotone"],
              ["screener", "IPO Screener", "ph:funnel-duotone"],
            ] as const
          ).map(([v, lbl, ico]) => (
            <button
              key={v}
              type="button"
              onClick={() => setMainView(v)}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-bold transition cursor-pointer ${
                mainView === v
                  ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-xs"
                  : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
              }`}
            >
              <Icon icon={ico} width="15" />
              <span>{lbl}</span>
            </button>
          ))}
        </div>
      </div>

      {/* VIEW 2: IPO LIST */}
      {mainView === "list" && (
        <IpoListView
          meta={meta}
          onSelectTicker={(tk) => {
            setSelectedTicker(tk);
            setTickerInput(tk);
            setMainView("analisis");
          }}
        />
      )}

      {/* VIEW 3: IPO SCREENER */}
      {mainView === "screener" && (
        <IpoScreenerView
          meta={meta}
          onOpenTicker={(tk, opts) => {
            setSelectedTicker(tk);
            setTickerInput(tk);
            setMode(opts.after ? "after" : "periodic");
            if (!opts.after && opts.lo) setDateFrom(opts.lo);
            if (opts.hi) setDateTo(opts.hi);
            setMainView("analisis");
          }}
        />
      )}

      {/* VIEW 1: ANALISIS */}
      {mainView === "analisis" && (
        <div className="space-y-4">
          <IpoAnalysisControls
            meta={meta}
            selectedTicker={selectedTicker}
            tickerInput={tickerInput}
            onTickerInputChange={setTickerInput}
            onSelectTicker={selectTicker}
            mode={mode}
            onModeChange={setMode}
            dateFrom={dateFrom}
            onDateFromChange={setDateFrom}
            dateTo={dateTo}
            onDateToChange={setDateTo}
            infoCard={infoCard}
            tickerError={tickerError}
            onOpenList={() => setMainView("list")}
          />

          {loading && !analysis ? (
            <div className="p-12 text-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
              Memuat analisis {selectedTicker}…
            </div>
          ) : analysis && !analysis.has_data ? (
            <div className="p-12 text-center rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] text-xs text-[var(--md-sys-color-on-surface-variant)]">
              Kode saham <b>{selectedTicker}</b> tidak memiliki data transaksi pada periode ini.
            </div>
          ) : (
            analysis && (
              <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 items-start">
                {/* LEFT COLUMN (7/12): Candlestick + Dual Donut + Top Buyer/Seller */}
                <div className="lg:col-span-7 space-y-4">
                  <IpoCandlestickChart
                    ticker={selectedTicker}
                    candles={analysis.candles || []}
                    ipoPrice={infoCard?.ipo_price}
                    mode={mode}
                    onAfterIpoClick={() => setMode("after")}
                  />

                  <BuyerSellerDonutCard data={cards?.buyer_seller_pie} />

                  <IpoTopBuyerSellerCard tbs={cards?.top_buyer_seller} />
                </div>

                {/* RIGHT COLUMN (5/12): Transaction | Info Tabs */}
                <div className="lg:col-span-5 space-y-4">
                  <div className="flex border-b border-[var(--md-sys-color-outline-variant)]">
                    {(
                      [
                        ["tx", "Transaction"],
                        ["info", "Info"],
                      ] as const
                    ).map(([t, lbl]) => (
                      <button
                        key={t}
                        type="button"
                        onClick={() => setRightTab(t)}
                        className={`px-4 py-2 text-xs font-extrabold uppercase tracking-wider border-b-2 -mb-px transition cursor-pointer ${
                          rightTab === t
                            ? "border-[var(--md-sys-color-primary)] text-[var(--md-sys-color-primary)]"
                            : "border-transparent text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                        }`}
                      >
                        {lbl}
                      </button>
                    ))}
                  </div>

                  {rightTab === "info" ? (
                    <IpoInfoTab infoCard={infoCard} ticker={selectedTicker} />
                  ) : (
                    <div className="space-y-4">
                      <IpoUnderwriterCard uwCard={cards?.underwriter_card} ticker={selectedTicker} />
                      <AfterIpoDistCard data={cards?.after_ipo_distribution} />
                      <IpoCategoryCard catCard={cards?.category_card} />
                      <IpoRetailActionCard
                        raCard={cards?.retail_action_card}
                        autoRetail={autoRetail}
                        onToggleAutoRetail={() => setAutoRetail((a) => !a)}
                        retailOverride={retailOverride}
                        onRetailOverrideChange={setRetailOverride}
                      />
                    </div>
                  )}
                </div>
              </div>
            )
          )}
        </div>
      )}
    </div>
  );
}
