import {
  CONSTANT_CONDITION_LABELS,
  INVESTOR_ADVANCED_IDS,
  STOCK_ADVANCED_IDS,
} from "./ctConstants";
import type {
  Category,
  ConditionOption,
  ExplorerPeriod,
  ExplorerRow,
  ExplorerTickerValue,
  Mode,
} from "./ctTypes";

export const fetcher = async (url: string) => {
  const response = await fetch(url, { credentials: "include" });
  if (!response.ok) {
    let detail = `HTTP ${response.status}`;
    try {
      const body = await response.json();
      detail = body?.detail ?? detail;
    } catch {
      // keep status
    }
    throw new Error(detail);
  }
  return response.json();
};

export function defaultCondition(mode: Mode, category: Category): string {
  if (category === "advanced") {
    return mode === "stocks" ? "naikPorsi" : "portofolioNaik";
  }
  return "bertambah";
}

export function conditionsFor(mode: Mode, category: Category): ConditionOption[] {
  if (category === "kepemilikan") {
    const ids = mode === "stocks"
      ? ["bertambah", "berkurang", "sidNaik", "sidTurun"]
      : ["bertambah", "berkurang"];
    const short: Record<string, string> = {
      bertambah: "Net Buy",
      berkurang: "Net Sell",
      sidNaik: "SID Bertambah",
      sidTurun: "SID Berkurang",
    };
    return ids.map((id) => ({ id, label: short[id] ?? id }));
  }
  if (category === "transaksi") {
    return ["bertambah", "berkurang"].map((id) => ({
      id,
      label: id === "bertambah" ? "Net Buy" : "Net Sell",
    }));
  }
  const ids = mode === "stocks" ? STOCK_ADVANCED_IDS : INVESTOR_ADVANCED_IDS;
  return ids.map((id) => ({ id, label: CONSTANT_CONDITION_LABELS[id] ?? id }));
}

export function fmtExplorerPct(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return "—";
  const abs = Math.abs(value);
  let text: string;
  if (abs >= 9.95) text = String(Math.round(abs));
  else if (abs >= 0.995) text = abs.toFixed(1);
  else text = abs.toFixed(2);
  return `${value >= 0 ? "+" : "-"}${text}%`;
}

export function fmtExplorerValM(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return "—";
  const sign = value < 0 ? "-" : "+";
  const abs = Math.abs(value);
  let display = abs;
  let unit = "M";
  if (abs >= 1000) {
    display = abs / 1000;
    unit = "T";
  }
  const digits = display >= 100 ? 0 : display >= 10 ? 1 : 2;
  return `${sign}${display.toFixed(digits)} ${unit}`;
}

export function fmtMtd(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return "—";
  return `${value >= 0 ? "+" : ""}${value.toFixed(2)}%`;
}

export function fmtPlainPct(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return "—";
  return `${value.toFixed(2)}%`;
}

export function fmtPortfolioPct(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return "—";
  return `${value >= 0 ? "+" : ""}${value.toFixed(2)}%`;
}

export function fmtSidDelta(diff: number | null | undefined, pct: number | null | undefined): string {
  if (diff === null || diff === undefined) return "—";
  const diffText = `${diff >= 0 ? "+" : "-"}${Math.abs(diff).toLocaleString("id-ID")} SID`;
  return `${diffText} | ${fmtExplorerPct(pct)}`;
}

export function csvCell(value: unknown): string {
  const text = value === null || value === undefined ? "" : String(value);
  if (/[",\r\n]/.test(text)) return `"${text.replace(/"/g, '""')}"`;
  return text;
}

export function toggleValue(values: string[], value: string): string[] {
  if (!values.includes(value)) return [...values, value];
  if (values.length === 1) return values;
  return values.filter((item) => item !== value);
}

export function tickerValueItems(tickers: ExplorerRow["tickers"]): ExplorerTickerValue[] {
  if (!Array.isArray(tickers) || tickers.length === 0) return [];
  if (typeof tickers[0] === "string") return [];
  return tickers as ExplorerTickerValue[];
}

export function mtdMonthLabel(period: ExplorerPeriod | undefined): string {
  const month = (period?.last_update_date ?? "").split(" ")[1];
  return month ? month.slice(0, 3) : "Okt";
}

export function sidMonthLabel(period: ExplorerPeriod | undefined): string {
  return (period?.month_curr ?? "September").slice(0, 3);
}

export function toneTextClass(tone: "up" | "down" | "neutral"): string {
  if (tone === "up") return "text-emerald-600 dark:text-emerald-400";
  if (tone === "down") return "text-rose-600 dark:text-rose-400";
  return "";
}
