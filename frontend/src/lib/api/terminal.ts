import type { StockOwnership } from "@/features/ownership/stock/stockTypes";

const API_BASE = "/api/v1/terminal";

async function terminalGet<T = unknown>(
  path: string,
  extraHeaders?: Record<string, string>
): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    headers: extraHeaders,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ detail: res.statusText }));
    throw new Error(err.detail || err.message || `API error ${res.status}: ${res.statusText}`);
  }
  return res.json();
}

async function terminalPost<T = unknown>(
  path: string,
  body: unknown,
  extraHeaders?: Record<string, string>
): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", ...extraHeaders },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ detail: res.statusText }));
    throw new Error(err.detail || err.message || `API error ${res.status}: ${res.statusText}`);
  }
  return res.json();
}

export function getAdminKey(): string {
  if (typeof window !== "undefined") {
    return localStorage.getItem("terminal_admin_key") || "";
  }
  return "";
}

export function setAdminKey(key: string): void {
  if (typeof window !== "undefined") {
    if (key && key.trim()) {
      localStorage.setItem("terminal_admin_key", key.trim());
    } else {
      localStorage.removeItem("terminal_admin_key");
    }
  }
}

function getAdminHeaders(): Record<string, string> {
  const key = getAdminKey();
  const headers: Record<string, string> = {};
  if (key) {
    headers["X-Admin-Key"] = key;
  }
  return headers;
}

export async function getMmdTickers<T = any>(): Promise<T> {
  return terminalGet<T>("/mmd/tickers");
}

export async function getMmdDetail<T = any>(ticker: string, limit = 60): Promise<T> {
  return terminalGet<T>(`/mmd/${encodeURIComponent(ticker)}?limit=${limit}`);
}

export async function getMmdAnalysis<T = any>(
  ticker: string,
  dateFrom?: string,
  dateTo?: string,
  autoRetail = true,
  retailOverride?: string
): Promise<T> {
  const params = new URLSearchParams();
  if (dateFrom) params.append("date_from", dateFrom);
  if (dateTo) params.append("date_to", dateTo);
  if (autoRetail !== undefined) params.append("auto_retail", String(autoRetail));
  if (retailOverride) params.append("retail_override", retailOverride);
  const q = params.toString() ? `?${params.toString()}` : "";
  return terminalGet<T>(`/mmd/${encodeURIComponent(ticker)}/analysis${q}`);
}

export async function getMmdT5<T = any>(ticker: string): Promise<T> {
  return terminalGet<T>(`/mmd/${encodeURIComponent(ticker)}/t5`);
}

export async function getMmdBrokerDrilldown<T = any>(
  ticker: string,
  broker: string,
  dateFrom?: string,
  dateTo?: string
): Promise<T> {
  const params = new URLSearchParams();
  params.append("broker", broker);
  if (dateFrom) params.append("date_from", dateFrom);
  if (dateTo) params.append("date_to", dateTo);
  return terminalGet<T>(`/mmd/${encodeURIComponent(ticker)}/broker-drilldown?${params.toString()}`);
}

export async function scanMmdExplorer<T = any>(
  category = "normal",
  criteria?: string,
  minVolumeVal?: number,
  limit = 50
): Promise<T> {
  const params = new URLSearchParams();
  params.append("category", category);
  if (criteria && criteria !== "ALL") params.append("criteria", criteria);
  if (minVolumeVal && minVolumeVal > 0) params.append("min_volume_val", String(minVolumeVal));
  params.append("limit", String(limit));
  return terminalGet<T>(`/mmd/explorer/scan?${params.toString()}`);
}

export async function getIpoHunterMeta<T = any>(): Promise<T> {
  return terminalGet<T>("/ipo-hunter/meta");
}

export async function getIpoHunterList<T = any>(params?: {
  uw?: string;
  lead_only?: boolean;
  val_bucket?: string;
  year?: number | string;
  q?: string;
  sort?: string;
  dir?: number;
}): Promise<T> {
  const sp = new URLSearchParams();
  if (params?.uw) sp.append("uw", params.uw);
  if (params?.lead_only) sp.append("lead_only", "true");
  if (params?.val_bucket) sp.append("val_bucket", params.val_bucket);
  if (params?.year) sp.append("year", String(params.year));
  if (params?.q) sp.append("q", params.q);
  if (params?.sort) sp.append("sort", params.sort);
  if (params?.dir !== undefined) sp.append("dir", String(params.dir));
  const qs = sp.toString() ? `?${sp.toString()}` : "";
  return terminalGet<T>(`/ipo-hunter/list${qs}`);
}

export async function getIpoHunterAnalysis<T = any>(
  ticker: string,
  params?: {
    mode?: "after" | "periodic";
    date_from?: string;
    date_to?: string;
    auto_retail?: boolean;
    retail_override?: string;
  }
): Promise<T> {
  const sp = new URLSearchParams();
  if (params?.mode) sp.append("mode", params.mode);
  if (params?.date_from) sp.append("date_from", params.date_from);
  if (params?.date_to) sp.append("date_to", params.date_to);
  if (params?.auto_retail !== undefined) sp.append("auto_retail", String(params.auto_retail));
  if (params?.retail_override) sp.append("retail_override", params.retail_override);
  const qs = sp.toString() ? `?${sp.toString()}` : "";
  return terminalGet<T>(`/ipo-hunter/analysis/${encodeURIComponent(ticker)}${qs}`);
}

export interface IpoScreenerPayload {
  after: boolean;
  date_from?: string;
  date_to?: string;
  quick?: number | null;
  year?: number | null;
  chg?: string;
  metric?: string;
  conds?: Array<{ k: string; op: string; v: string }>;
  sort_key?: string | null;
  dir?: number | null;
}

export async function runIpoHunterScreener<T = any>(payload: IpoScreenerPayload): Promise<T> {
  return terminalPost<T>("/ipo-hunter/screener", payload);
}

export async function getStockOwnershipProfile(ticker: string): Promise<StockOwnership | null> {
  try {
    const res = await fetch(`/api/v1/ownership/stock/${encodeURIComponent(ticker)}`, {
      credentials: "include",
    });
    if (!res.ok) return null;
    return (await res.json()) as StockOwnership;
  } catch {
    return null;
  }
}

export async function getForeignFlowScreener<T = any>(
  universe = "ff",
  window = 40,
  fpBadge?: string,
  transBucket?: string,
  badge?: string,
  limit = 200
): Promise<T> {
  const params = new URLSearchParams();
  params.append("universe", universe);
  params.append("window", String(window));
  if (fpBadge) params.append("fp_badge", fpBadge);
  if (transBucket) params.append("trans_bucket", transBucket);
  if (badge) params.append("badge", badge);
  params.append("limit", String(limit));
  return terminalGet<T>(`/foreign-flow/screener?${params.toString()}`);
}

export async function getForeignFlowSeries<T = any>(ticker: string, limit = 250): Promise<T> {
  return terminalGet<T>(`/foreign-flow/series/${encodeURIComponent(ticker)}?limit=${limit}`);
}

export async function getForeignFlowIhsg<T = any>(): Promise<T> {
  return terminalGet<T>("/foreign-flow/ihsg");
}

export async function getKongloReport<T = any>(): Promise<T> {
  return terminalGet<T>("/konglo/report");
}

export async function getKongloGroups<T = any>(): Promise<T> {
  return terminalGet<T>("/konglo/groups");
}

export async function getKongloSummary<T = any>(group?: string, limit = 200): Promise<T> {
  const params = new URLSearchParams();
  if (group) params.append("group", group);
  params.append("limit", String(limit));
  return terminalGet<T>(`/konglo/summary?${params.toString()}`);
}

export async function getKongloGroupDetail<T = any>(groupName: string): Promise<T> {
  return terminalGet<T>(`/konglo/group/${encodeURIComponent(groupName)}`);
}

export async function getMarketReportStatus<T = any>(): Promise<T> {
  return terminalGet<T>("/market-report/status", getAdminHeaders());
}

export async function triggerGenerateMarketReport<T = any>(
  session = 1,
  rateLimit = 180.0,
  workers = 8
): Promise<T> {
  return terminalPost<T>(
    "/market-report/generate",
    { session, rate_limit: rateLimit, workers },
    getAdminHeaders()
  );
}

export async function getMarketReportHistory<T = any>(): Promise<T> {
  return terminalGet<T>("/market-report/history", getAdminHeaders());
}

export async function getMarketReportDetail<T = any>(filename: string): Promise<T> {
  return terminalGet<T>(
    `/market-report/detail?filename=${encodeURIComponent(filename)}`,
    getAdminHeaders()
  );
}

export async function updateStockbitToken<T = any>(token: string): Promise<T> {
  return terminalPost<T>("/market-report/update-token", { token }, getAdminHeaders());
}

export async function triggerSmtrackerAutoRenew<T = any>(): Promise<T> {
  return terminalPost<T>("/market-report/auto-renew", {}, getAdminHeaders());
}

export async function getMarketReportLogs<T = any>(): Promise<T> {
  return terminalGet<T>("/market-report/logs", getAdminHeaders());
}

export async function resetMarketReportStatus<T = any>(): Promise<T> {
  return terminalPost<T>("/market-report/reset-status", {}, getAdminHeaders());
}
