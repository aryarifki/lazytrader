import { useEffect } from "react";
import useSWR from "swr";
import { useAppStore } from "@/store/useAppStore";

const fetcher = async (url: string) => {
  const res = await fetch(url, { credentials: "include" });

  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(errorData.detail || "Gagal mengambil data analitik");
  }

  return res.json();
};

export function useForeignFlowAnalytics(ticker: string, lookbackDays: number = 20) {
  const { data, error, isLoading, isValidating, mutate } = useSWR(
    ticker ? `/api/foreign-flow/${ticker}?lookback_days=${lookbackDays}` : null,
    fetcher,
    {
      revalidateOnFocus: false,
      dedupingInterval: 60000,
      shouldRetryOnError: false,
      keepPreviousData: true,
    }
  );

  const lastSyncTimestamp = useAppStore((state) => state.lastSyncTimestamp);

  useEffect(() => {
    if (lastSyncTimestamp > 0) {
      mutate();
    }
  }, [lastSyncTimestamp, mutate]);

  useEffect(() => {
    const handleSync = () => {
      mutate();
    };
    window.addEventListener("investowl:sync", handleSync);
    return () => window.removeEventListener("investowl:sync", handleSync);
  }, [mutate]);

  return {
    data,
    isLoading,
    isFetching: isValidating,
    isError: error,
    mutate,
  };
}
