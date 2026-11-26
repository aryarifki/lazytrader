import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

export interface CacheLoadingState {
  isOpen: boolean;
  title?: string;
  step?: string;
  detail?: string;
  progress?: number;
  isSuccess?: boolean;
}

interface AppState {
  theme: "dark" | "light";
  activeTicker: string;
  sidebarOpen: boolean;
  summaryDays: number;
  universe: string;
  analysisDate: string;
  isDateManual: boolean;
  windowDays: number;
  horizon: number;
  minEvents: number;
  minNetBuy: number;
  localWatchlist: string[];
  cacheLoading: CacheLoadingState;
  lastSyncTimestamp: number;

  setTheme: (theme: "dark" | "light") => void;
  toggleTheme: () => void;
  setActiveTicker: (ticker: string) => void;
  setSidebarOpen: (isOpen: boolean) => void;
  setSummaryDays: (days: number) => void;
  setUniverse: (universe: string) => void;
  setAnalysisDate: (date: string) => void;
  setIsDateManual: (manual: boolean) => void;
  setWindowDays: (days: number) => void;
  setHorizon: (horizon: number) => void;
  setMinEvents: (events: number) => void;
  setMinNetBuy: (buy: number) => void;
  setLocalWatchlist: (list: string[]) => void;
  setCacheLoading: (state: Partial<CacheLoadingState> | boolean) => void;
  triggerSync: () => void;
}

export const useAppStore = create<AppState>()(
  persist(
    (set) => ({
      theme: "dark",
      activeTicker: "BBCA",
      sidebarOpen: false,
      summaryDays: 30,
      universe: "watchlist",
      analysisDate: "",
      isDateManual: false,
      windowDays: 20,
      horizon: 10,
      minEvents: 5,
      minNetBuy: 0,
      localWatchlist: [],
      cacheLoading: {
        isOpen: false,
        title: "Sinkronisasi Cache & Data",
        step: "Memuat cache...",
        detail: "Redis & In-Memory Layer",
        progress: 0,
        isSuccess: false,
      },
      lastSyncTimestamp: 0,

      setTheme: (theme) => set({ theme }),
      toggleTheme: () => set((state) => ({ theme: state.theme === "dark" ? "light" : "dark" })),
      setActiveTicker: (ticker) => set({ activeTicker: ticker.toUpperCase() }),
      setSidebarOpen: (isOpen) => set({ sidebarOpen: isOpen }),
      setSummaryDays: (days) => set({ summaryDays: days }),
      setUniverse: (universe) => set({ universe }),
      setAnalysisDate: (date) => set({ analysisDate: date }),
      setIsDateManual: (manual) => set({ isDateManual: manual }),
      setWindowDays: (days) => set({ windowDays: days }),
      setHorizon: (horizon) => set({ horizon }),
      setMinEvents: (events) => set({ minEvents: events }),
      setMinNetBuy: (buy) => set({ minNetBuy: buy }),
      setLocalWatchlist: (list) => set({ localWatchlist: list }),
      triggerSync: () => set({ lastSyncTimestamp: Date.now() }),
      setCacheLoading: (val) =>
        set((state) => {
          if (typeof val === "boolean") {
            return {
              cacheLoading: {
                ...state.cacheLoading,
                isOpen: val,
                progress: val ? 25 : 100,
                isSuccess: false,
              },
            };
          }
          return {
            cacheLoading: {
              ...state.cacheLoading,
              ...val,
            },
          };
        }),
    }),
    {
      name: "investowl-global-state",
      version: 3,
      storage: createJSONStorage(() => localStorage),
      migrate: (persistedState: unknown, version: number) => {
        const st = (persistedState || {}) as Record<string, unknown>;
        if (!persistedState || version < 3) {
          return { ...st, windowDays: 20 };
        }
        return st;
      },
      partialize: (state) => {
        const {
          sidebarOpen: _sidebarOpen,
          analysisDate: _analysisDate,
          isDateManual: _isDateManual,
          cacheLoading: _cacheLoading,
          lastSyncTimestamp: _lastSyncTimestamp,
          ...rest
        } = state;
        return rest;
      },
    }
  )
);
