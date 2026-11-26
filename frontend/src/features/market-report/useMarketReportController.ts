import { useCallback, useEffect, useRef, useState } from "react";
import {
  getAdminKey,
  getMarketReportDetail,
  getMarketReportHistory,
  getMarketReportLogs,
  getMarketReportStatus,
  resetMarketReportStatus,
  setAdminKey,
  triggerGenerateMarketReport,
  triggerSmtrackerAutoRenew,
  updateStockbitToken,
} from "@/lib/api/terminal";
import {
  getErrorMessage,
  isAuthError,
  notifyDesktopCompletion,
  playCompletionChime,
  type HistoryItem,
  type ReportData,
  type StatusResponse,
  type ToastState,
} from "./marketReportTypes";

export function useMarketReportController() {
  const [status, setStatus] = useState<StatusResponse | null>(null);
  const [activeReport, setActiveReport] = useState<ReportData | null>(null);
  const [history, setHistory] = useState<HistoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);
  const [selectedHistoryFile, setSelectedHistoryFile] = useState("");
  const [floatingMinimized, setFloatingMinimized] = useState(false);
  const [elapsedSeconds, setElapsedSeconds] = useState(0);
  const [rateLimit, setRateLimit] = useState(180);
  const [workers, setWorkers] = useState(8);

  const [showLogcatModal, setShowLogcatModal] = useState(false);
  const [terminalCollapsed, setTerminalCollapsed] = useState(false);
  const [logs, setLogs] = useState<string[]>([]);
  const [autoScrollLogs, setAutoScrollLogs] = useState(true);
  const [logLoading, setLogLoading] = useState(false);
  const logContainerRef = useRef<HTMLDivElement | null>(null);
  const embeddedLogContainerRef = useRef<HTMLDivElement | null>(null);

  const [showTokenModal, setShowTokenModal] = useState(false);
  const [tokenInput, setTokenInput] = useState("");
  const [tokenSubmitting, setTokenSubmitting] = useState(false);
  const [tokenError, setTokenError] = useState<string | null>(null);
  const [tokenSuccess, setTokenSuccess] = useState<string | null>(null);

  const [adminKeyState, setAdminKeyState] = useState("");
  const [showAdminKeyModal, setShowAdminKeyModal] = useState(false);
  const [adminKeyInput, setAdminKeyInput] = useState("");

  const [toast, setToast] = useState<ToastState>({ message: "", type: "info", visible: false });
  const toastTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const prevRunningRef = useRef(false);

  const showToast = useCallback(
    (message: string, type: "info" | "success" | "error" = "info", duration = 3500) => {
      if (toastTimerRef.current) clearTimeout(toastTimerRef.current);
      setToast({ message, type, visible: true });
      toastTimerRef.current = setTimeout(() => setToast((p) => ({ ...p, visible: false })), duration);
    },
    []
  );

  const handleCopy = useCallback(
    (text: string, label: string) => {
      if (!text) return;
      navigator.clipboard
        .writeText(text)
        .then(() => showToast(`Berhasil menyalin ${label} ke clipboard!`, "success", 2500))
        .catch(() => showToast(`Gagal menyalin ${label}.`, "error", 2500));
    },
    [showToast]
  );

  const handleActionAuthError = useCallback(
    (err: unknown, authMsg: string, fallbackPrefix: string) => {
      const msg = getErrorMessage(err);
      if (isAuthError(msg)) {
        showToast(authMsg, "error", 5000);
        setShowAdminKeyModal(true);
      } else {
        showToast(`${fallbackPrefix}: ${msg}`, "error", 4000);
      }
    },
    [showToast]
  );

  const fetchState = useCallback(async () => {
    try {
      const [res, histRes] = await Promise.all([
        getMarketReportStatus<{ data?: StatusResponse }>(),
        getMarketReportHistory<{ history?: HistoryItem[] }>(),
      ]);
      if (res?.data) {
        setStatus(res.data);
        if (res.data.recent_logs?.length) setLogs(res.data.recent_logs);
        if (res.data.latest_report) {
          const latest = res.data.latest_report;
          setActiveReport((prev) => prev ?? latest);
          setSelectedHistoryFile((prev) => prev || latest.filename);
        }
      }
      if (histRes?.history) setHistory(histRes.history);
    } catch (err) {
      console.error("Gagal mengambil status market report:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchFullLogs = useCallback(async () => {
    setLogLoading(true);
    try {
      const res = await getMarketReportLogs<{ logs?: string[] }>();
      if (res?.logs) setLogs(res.logs);
    } catch (err) {
      console.error("Gagal mengambil logs:", err);
    } finally {
      setLogLoading(false);
    }
  }, []);

  useEffect(() => {
    if (typeof window !== "undefined") {
      const urlParams = new URLSearchParams(window.location.search);
      const urlKey = urlParams.get("key") || urlParams.get("admin_key");
      if (urlKey) {
        setAdminKey(urlKey);
        setAdminKeyState(urlKey);
        window.history.replaceState({}, document.title, window.location.pathname);
        showToast("Kunci Admin Terminal berhasil dikonfigurasi dari URL!", "success", 3500);
      } else {
        setAdminKeyState(getAdminKey());
      }
      if ("Notification" in window && Notification.permission === "default") {
        Notification.requestPermission().catch(() => {});
      }
    }
    fetchState();
    fetchFullLogs();
    return () => { if (toastTimerRef.current) clearTimeout(toastTimerRef.current); };
  }, [fetchFullLogs, fetchState, showToast]);

  useEffect(() => {
    if (!autoScrollLogs) return;
    if (logContainerRef.current) logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight;
    if (embeddedLogContainerRef.current) embeddedLogContainerRef.current.scrollTop = embeddedLogContainerRef.current.scrollHeight;
  }, [logs, showLogcatModal, autoScrollLogs, terminalCollapsed]);

  useEffect(() => {
    if (!status?.is_running) return;
    setElapsedSeconds(0);
    const timer = setInterval(() => setElapsedSeconds((prev) => prev + 1), 1000);
    return () => clearInterval(timer);
  }, [status?.is_running]);

  useEffect(() => {
    if (!status?.is_running && !showLogcatModal) return;
    const poll = setInterval(async () => {
      try {
        const res = await getMarketReportStatus<{ data?: StatusResponse }>();
        const newData = res?.data;
        if (!newData) return;
        setStatus(newData);
        if (newData.recent_logs?.length) setLogs(newData.recent_logs);

        if (prevRunningRef.current && !newData.is_running) {
          if (newData.error) {
            showToast(`Pemindaian terhenti: ${newData.error}`, "error", 6000);
          } else if (newData.latest_report) {
            setActiveReport(newData.latest_report);
            setSelectedHistoryFile(newData.latest_report.filename);
            playCompletionChime();
            showToast(`Laporan Sesi ${newData.session} Berhasil Dibuat!`, "success", 5000);
            notifyDesktopCompletion(newData.session);
            getMarketReportHistory<{ history?: HistoryItem[] }>()
              .then((h) => h?.history && setHistory(h.history))
              .catch(() => {});
          }
        }
        prevRunningRef.current = newData.is_running;
      } catch (err) {
        console.error("Poll error:", err);
      }
    }, 1500);
    return () => clearInterval(poll);
  }, [showLogcatModal, showToast, status?.is_running]);

  const applyAdminKey = (key: string, message: string, tone: "success" | "info") => {
    setAdminKey(key);
    setAdminKeyState(key);
    setAdminKeyInput("");
    setShowAdminKeyModal(false);
    showToast(message, tone, 3000);
    fetchState();
    fetchFullLogs();
  };

  const handleSaveAdminKey = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = adminKeyInput.trim();
    applyAdminKey(
      trimmed,
      trimmed ? "Override X-Admin-Key berhasil disimpan!" : "Override X-Admin-Key dihapus.",
      trimmed ? "success" : "info"
    );
  };

  const handleClearAdminKey = () =>
    applyAdminKey("", "Override X-Admin-Key dihapus dari browser.", "info");

  const handleGenerate = async (session: number) => {
    const sesLabel = session === 1 ? "Sesi 1 (12:00 WIB)" : "End-of-Day Sesi 2 (16:00 WIB)";
    if (!window.confirm(`Mulai pemindaian 962 emiten untuk Laporan Pasar ${sesLabel}?`)) return;
    setActionLoading(true);
    setFloatingMinimized(false);
    try {
      await triggerGenerateMarketReport(session, rateLimit, workers);
      showToast(`Memulai pemindaian 962 emiten untuk Laporan Sesi ${session}...`, "info", 4000);
      await fetchState();
    } catch (err) {
      handleActionAuthError(err, "Memerlukan otorisasi admin untuk memulai pemindaian.", "Gagal memulai");
    } finally {
      setActionLoading(false);
    }
  };

  const handleUpdateToken = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!tokenInput.trim()) return setTokenError("Harap masukkan token Stockbit.");
    setTokenSubmitting(true);
    setTokenError(null);
    setTokenSuccess(null);
    try {
      const res = await updateStockbitToken<{ success?: boolean; message?: string }>(tokenInput.trim());
      if (res?.success) {
        setTokenSuccess(res.message || "Token berhasil divalidasi dan disimpan!");
        showToast("Token Stockbit berhasil diperbarui!", "success", 4000);
        await fetchState();
        setTimeout(() => {
          setShowTokenModal(false);
          setTokenInput("");
          setTokenSuccess(null);
        }, 1800);
      } else {
        setTokenError(res?.message || "Token tidak valid.");
      }
    } catch (err) {
      const msg = getErrorMessage(err);
      if (isAuthError(msg)) { setTokenError("Akses ditolak: Silakan periksa sesi admin atau X-Admin-Key."); setShowAdminKeyModal(true); }
      else setTokenError(msg || "Gagal menghubungi server untuk memvalidasi token.");
    } finally { setTokenSubmitting(false); }
  };

  const handleAutoRenew = async () => {
    setActionLoading(true);
    try {
      await triggerSmtrackerAutoRenew();
      setShowLogcatModal(true);
      fetchFullLogs();
      showToast("Memulai Auto-Renew SMtracker. Cek HP Anda bila ada konfirmasi Stockbit!", "info", 5000);
      await fetchState();
    } catch (err) {
      handleActionAuthError(err, "Memerlukan otorisasi admin untuk memicu Auto-Renew.", "Gagal auto-renew");
    } finally { setActionLoading(false); }
  };

  const handleResetStatus = async () => {
    try {
      await resetMarketReportStatus();
      showToast("Status generator berhasil direset.", "info", 2500);
      await fetchState();
    } catch (err) {
      if (isAuthError(getErrorMessage(err))) setShowAdminKeyModal(true);
      console.error(err);
    }
  };

  const handleSelectHistory = async (filename: string) => {
    setSelectedHistoryFile(filename);
    try {
      const res = await getMarketReportDetail<{ report?: ReportData }>(filename);
      if (res?.report) {
        setActiveReport({ ...res.report, date: res.report.date || filename.split("_").pop()?.replace(".txt", "") || "", generated_at: res.report.generated_at || "" });
        showToast(`Memuat arsip: ${filename}`, "info", 2000);
      }
    } catch (err) {
      console.error("Gagal memuat detail file arsip:", err);
      showToast("Gagal memuat arsip laporan.", "error", 3000);
    }
  };

  return {
    status, activeReport, history, loading, actionLoading, selectedHistoryFile,
    floatingMinimized, setFloatingMinimized, elapsedSeconds, rateLimit, setRateLimit,
    workers, setWorkers, showLogcatModal, setShowLogcatModal, terminalCollapsed,
    setTerminalCollapsed, logs, autoScrollLogs, setAutoScrollLogs, logLoading,
    logContainerRef, embeddedLogContainerRef, showTokenModal, setShowTokenModal,
    tokenInput, setTokenInput, tokenSubmitting, tokenError, setTokenError,
    tokenSuccess, setTokenSuccess, adminKeyState, showAdminKeyModal, setShowAdminKeyModal,
    adminKeyInput, setAdminKeyInput, toast, handleCopy, fetchState, fetchFullLogs,
    handleSaveAdminKey, handleClearAdminKey, handleGenerate, handleUpdateToken,
    handleAutoRenew, handleResetStatus, handleSelectHistory,
  };
}
