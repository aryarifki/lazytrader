import { useEffect, useState } from "react";
import { runIpoHunterScreener } from "@/lib/api/terminal";
import type { IpoApiResponse, IpoMeta, IpoScreenerCondition, IpoScreenerData } from "./ipoTypes";
import { SCR_METRICS, SCR_PRESETS, type ScreenerMetricKey } from "./ipoScreenerConstants";
import IpoScreenerOptions from "./IpoScreenerOptions";
import IpoScreenerTable from "./IpoScreenerTable";

export default function IpoScreenerView({
  meta,
  onOpenTicker,
}: {
  meta: IpoMeta | null;
  onOpenTicker: (ticker: string, opts: { after: boolean; lo?: string; hi?: string }) => void;
}) {
  const dates = meta?.available_dates || [];
  const [quick, setQuick] = useState<number | null>(5);
  const [after, setAfter] = useState(false);
  const [dateTo, setDateTo] = useState<string>("");
  const [dateFrom, setDateFrom] = useState<string>("");
  const [year, setYear] = useState<string>("");
  const [chg, setChg] = useState<string>("all");
  const [metric, setMetric] = useState<string>("topPerf");
  const [conds, setConds] = useState<IpoScreenerCondition[]>([]);
  const [sortKey, setSortKey] = useState<string | null>(null);
  const [sortDir, setSortDir] = useState<number>(-1);
  const [res, setRes] = useState<IpoScreenerData | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (dates.length > 0 && !dateTo) {
      setDateTo(dates[dates.length - 1]);
      setDateFrom(dates[Math.max(0, dates.length - 5)]);
    }
  }, [dates, dateTo]);

  useEffect(() => {
    if (!dateTo) return;
    let cancelled = false;
    const t = setTimeout(() => {
      setLoading(true);
      runIpoHunterScreener<IpoApiResponse<IpoScreenerData>>({
        after,
        date_from: after ? undefined : dateFrom || undefined,
        date_to: dateTo,
        quick,
        year: year ? Number(year) : null,
        chg,
        metric,
        conds,
        sort_key: sortKey,
        dir: sortDir,
      })
        .then((r) => {
          if (cancelled) return;
          setRes(r.data);
          setLoading(false);
        })
        .catch(() => {
          if (!cancelled) setLoading(false);
        });
    }, 200);
    return () => {
      cancelled = true;
      clearTimeout(t);
    };
  }, [after, dateFrom, dateTo, quick, year, chg, metric, conds, sortKey, sortDir]);

  const handleQuick = (n: number) => {
    setAfter(false);
    setQuick(n);
    if (dates.length && dateTo) {
      const idx = dates.indexOf(dateTo);
      const endIdx = idx >= 0 ? idx : dates.length - 1;
      setDateFrom(dates[Math.max(0, endIdx - n + 1)]);
    }
  };

  const resetAll = () => {
    setQuick(5);
    setAfter(false);
    setYear("");
    setChg("all");
    setMetric("topPerf");
    setConds([]);
    setSortKey(null);
    setSortDir(-1);
    if (dates.length) {
      setDateTo(dates[dates.length - 1]);
      setDateFrom(dates[Math.max(0, dates.length - 5)]);
    }
  };

  const activePreset = SCR_PRESETS.find((p) => p.id === metric);
  const activeCol: ScreenerMetricKey = activePreset ? activePreset.col : (metric as ScreenerMetricKey);
  const activeDesc = activePreset ? activePreset.desc : SCR_METRICS.find((m) => m.k === metric)?.desc || "";

  const cols: ScreenerMetricKey[] = ["chg", "ratrans"];
  if (!cols.includes(activeCol)) cols.push(activeCol);
  conds.forEach((c) => {
    const k = c.k as ScreenerMetricKey;
    if (!cols.includes(k)) cols.push(k);
  });

  const currentSk = sortKey && cols.includes(sortKey as ScreenerMetricKey) ? sortKey : activeCol;

  return (
    <div className="space-y-4">
      <IpoScreenerOptions
        dates={dates}
        years={meta?.years || []}
        quick={quick}
        after={after}
        dateFrom={dateFrom}
        dateTo={dateTo}
        year={year}
        chg={chg}
        metric={metric}
        conds={conds}
        nDays={res?.n_days}
        hiDate={res?.hi}
        activeDesc={activeDesc}
        onQuick={handleQuick}
        onToggleAfter={() => setAfter((a) => !a)}
        onDateFromChange={(v) => {
          setDateFrom(v);
          setQuick(null);
          setAfter(false);
        }}
        onDateToChange={(v) => {
          setDateTo(v);
          if (!after && quick) {
            const idx = dates.indexOf(v);
            if (idx >= 0) setDateFrom(dates[Math.max(0, idx - quick + 1)]);
          }
        }}
        onYearChange={setYear}
        onChgChange={setChg}
        onMetricChange={(val) => {
          setMetric(val);
          setSortKey(null);
          const pr = SCR_PRESETS.find((p) => p.id === val);
          setSortDir(pr ? pr.dir : -1);
        }}
        onCondsChange={setConds}
        onReset={resetAll}
      />

      <IpoScreenerTable
        res={res}
        loading={loading}
        cols={cols}
        after={after}
        activeCol={activeCol}
        currentSk={currentSk}
        sortDir={sortDir}
        onSortCol={(k, isSorted) => {
          if (isSorted) setSortDir((d) => -d);
          else {
            setSortKey(k);
            setSortDir(-1);
          }
        }}
        onOpenTicker={onOpenTicker}
      />
    </div>
  );
}
