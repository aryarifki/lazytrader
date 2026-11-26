export interface MarketBreadth {
  up: number;
  down: number;
  flat: number;
}

export interface ReportData {
  session: number;
  filename: string;
  date: string;
  generated_at: string;
  masked_text: string;
  unmasked_text: string;
  duration_seconds?: number;
  total_scanned?: number;
  market_breadth?: MarketBreadth;
}

export interface StatusResponse {
  status?: string; // "idle" | "running" | "completed" | "error"
  is_running: boolean;
  session: number;
  progress: number;
  total: number;
  pct: number;
  message: string;
  started_at: string | null;
  completed_at: string | null;
  error: string | null;
  latest_report: ReportData | null;
  token_valid: boolean;
  current_time_wib: string;
  recent_logs?: string[];
}

export interface HistoryItem {
  filename: string;
  session: number;
  date: string;
  modified_at: string;
  size_bytes: number;
}

export interface ToastState {
  message: string;
  type: "info" | "success" | "error";
  visible: boolean;
}

export interface MarketSessionInfo {
  label: string;
  badgeClass: string;
}

/** Gentle Web Audio notification chime when a market scan completes. */
export function playCompletionChime(): void {
  try {
    const win = window as Window & { webkitAudioContext?: typeof AudioContext };
    const AudioCtx = window.AudioContext || win.webkitAudioContext;
    if (!AudioCtx) return;

    const ctx = new AudioCtx();
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.connect(gain);
    gain.connect(ctx.destination);
    osc.type = "sine";

    const now = ctx.currentTime;
    osc.frequency.setValueAtTime(587.33, now); // D5
    osc.frequency.setValueAtTime(880, now + 0.15); // A5
    gain.gain.setValueAtTime(0.12, now);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.6);

    osc.start(now);
    osc.stop(now + 0.6);
  } catch {
    // Ignore audio autoplay restrictions
  }
}

export function formatTimer(sec: number): string {
  const safe = Math.max(0, Math.floor(sec));
  const m = Math.floor(safe / 60);
  const s = safe % 60;
  return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
}

export function formatDuration(sec?: number): string {
  if (typeof sec !== "number" || !Number.isFinite(sec) || sec <= 0) return "-";
  if (sec < 60) return `${sec.toFixed(1)}d`;
  const m = Math.floor(sec / 60);
  const s = Math.round(sec % 60);
  return `${m}m ${s}d`;
}

export function formatFileSize(bytes?: number): string {
  if (typeof bytes !== "number" || !Number.isFinite(bytes) || bytes <= 0) return "-";
  if (bytes < 1024) return `${bytes} B`;
  return `${(bytes / 1024).toFixed(1)} KB`;
}

export function getMarketSessionStatus(): MarketSessionInfo {
  const now = new Date();
  const utc = now.getTime() + now.getTimezoneOffset() * 60000;
  const wib = new Date(utc + 3600000 * 7);
  const day = wib.getDay();
  const timeVal = wib.getHours() * 60 + wib.getMinutes();

  if (day === 0 || day === 6) {
    return {
      label: "Bursa Tutup (Akhir Pekan)",
      badgeClass:
        "text-[var(--md-sys-color-on-surface-variant)] bg-[var(--md-sys-color-surface-container-high)] border-[var(--md-sys-color-outline-variant)]",
    };
  }
  if (timeVal >= 540 && timeVal < 720) {
    return {
      label: "Pasar Berjalan: Sesi 1",
      badgeClass: "text-emerald-700 dark:text-emerald-400 bg-emerald-500/12 border-emerald-500/25",
    };
  }
  if (timeVal >= 720 && timeVal < 810) {
    return {
      label: "Istirahat Siang (Sesi 1 Selesai)",
      badgeClass: "text-amber-700 dark:text-amber-400 bg-amber-500/12 border-amber-500/25",
    };
  }
  if (timeVal >= 810 && timeVal < 960) {
    return {
      label: "Pasar Berjalan: Sesi 2",
      badgeClass: "text-sky-700 dark:text-sky-400 bg-sky-500/12 border-sky-500/25",
    };
  }
  return {
    label: "Pasar Tutup (Pasca Sesi 2)",
    badgeClass: "text-purple-700 dark:text-purple-400 bg-purple-500/12 border-purple-500/25",
  };
}

export function getErrorMessage(err: unknown): string {
  return err instanceof Error ? err.message : String(err);
}

export function isAuthError(msg: string): boolean {
  return msg.includes("401") || msg.includes("Akses ditolak") || msg.includes("Kunci Admin");
}

export function notifyDesktopCompletion(session: number): void {
  if (typeof window === "undefined" || !("Notification" in window)) return;
  if (Notification.permission !== "granted") return;
  try {
    new Notification(`Laporan Pasar Sesi ${session} Selesai!`, {
      body: "Pemindaian emiten BEI berhasil diselesaikan. Teks Masked & Unmasked siap disalin.",
    });
  } catch {
    // Ignore Notification errors
  }
}

export function getLogLineClass(line: string): string {
  const isErr = line.includes("[ERROR]") || line.includes("❌") || line.includes("Exception");
  const isProgress = line.includes("[PROGRESS]") || line.includes("Memindai");
  const isWarn = line.includes("[WARN]") || line.includes("⚠️");
  const isCheck =
    line.includes("[Pre-Check]") ||
    line.includes("🔐") ||
    line.includes("🚀") ||
    line.includes("[SMtracker]");
  const isComplete = line.includes("✅") || line.includes("🎉") || line.includes("[COMPLETE]");

  if (isErr) return "bg-rose-950/45 text-rose-300 font-semibold border-l-2 border-rose-500";
  if (isProgress) return "text-emerald-400 font-medium";
  if (isWarn) return "text-amber-300 bg-amber-950/25";
  if (isCheck) return "text-cyan-300 font-medium";
  if (isComplete) return "text-emerald-300 font-bold bg-emerald-950/25";
  return "text-zinc-300";
}
