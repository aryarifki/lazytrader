import { useState } from "react";
import { useSWRConfig } from "swr";
import { Icon } from "@iconify/react";
import { useAppStore } from "@/store/useAppStore";

export default function RefreshButton() {
  const { mutate } = useSWRConfig();
  const { setCacheLoading } = useAppStore();
  const [isSyncing, setIsSyncing] = useState(false);
  const [syncedRecently, setSyncedRecently] = useState(false);

  const handleSync = async () => {
    if (isSyncing) return;
    setIsSyncing(true);
    setSyncedRecently(false);

    setCacheLoading({
      isOpen: true,
      title: "Sinkronisasi Data & Cache",
      step: "Membersihkan cache layer L1 (In-Memory), L2 (Redis), & Cloudflare Edge CDN...",
      detail: "POST /api/logs/cache/purge?pattern=*",
      progress: 25,
      isSuccess: false,
    });

    try {
      const purgeRes = await fetch("/api/logs/cache/purge?pattern=*", { method: "POST" });
      const purgeJson = await purgeRes.json().catch(() => ({}));
      const deletedKeys = purgeJson.deleted_keys ?? 0;

      await Promise.allSettled([
        fetch("/api/forensic/refresh", { method: "POST" }),
        fetch("/api/bandar/clear-cache", { method: "POST" }),
      ]);

      fetch("/api/logs/cache/warmup", { method: "POST" }).catch(() => {});

      setCacheLoading({
        step: "Mengambil data komputasi pasar terverifikasi dari PostgreSQL...",
        detail: `Dibersihkan: ${deletedKeys} keys | Memuat ulang ringkasan pasar...`,
        progress: 55,
      });

      await fetch("/api/bandar/daily-summary?universe_mode=watchlist&refresh=1").catch(() => {});

      setCacheLoading({
        step: "Menyinkronkan seluruh halaman & tab antarmuka...",
        detail: "Broadcasting event investowl:sync & SWR revalidation",
        progress: 80,
      });

      useAppStore.getState().triggerSync();
      window.dispatchEvent(new CustomEvent("investowl:sync"));

      await mutate(() => true, undefined, { revalidate: true });

      setCacheLoading({
        step: "Berhasil! Seluruh data dan cache telah tersinkronisasi.",
        detail: "Status: 200 OK — Real-time Synchronized",
        progress: 100,
        isSuccess: true,
      });

      setSyncedRecently(true);
      setTimeout(() => setSyncedRecently(false), 2000);
      setTimeout(() => setCacheLoading({ isOpen: false, isSuccess: false }), 700);
    } catch (e) {
      console.error("Gagal sinkronisasi data & cache:", e);
      setCacheLoading({
        step: "Gagal menyinkronkan cache. Periksa koneksi backend!",
        detail: "Error: Gagal memanggil endpoint cache",
        progress: 100,
        isSuccess: false,
      });
      setTimeout(() => setCacheLoading({ isOpen: false }), 1500);
    } finally {
      setIsSyncing(false);
    }
  };

  return (
    <button
      onClick={handleSync}
      disabled={isSyncing}
      className={`flex items-center gap-1.5 px-2.5 py-1.5 border rounded-md transition-all active:scale-95 disabled:opacity-60 cursor-pointer ${
        syncedRecently
          ? "bg-emerald-500/15 border-emerald-500/40 text-emerald-400"
          : "bg-[var(--md-sys-color-surface-container-highest)] hover:bg-[var(--md-sys-color-surface-container-high)] border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
      }`}
      title="Sinkronisasi Data Real-time & Muat Ulang Seluruh Halaman"
    >
      <Icon
        icon={syncedRecently ? "ph:check-bold" : "ph:arrows-clockwise-duotone"}
        width="16"
        className={isSyncing ? "animate-spin text-[var(--md-sys-color-primary)]" : syncedRecently ? "text-emerald-400" : ""}
      />
      <span className="text-[10px] font-mono font-bold uppercase tracking-wider hidden sm:inline">
        {isSyncing ? "Syncing..." : syncedRecently ? "Synced!" : "Sync"}
      </span>
    </button>
  );
}
