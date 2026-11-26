import { useState } from "react";
import useSWR from "swr";
import { Icon } from "@iconify/react";
import { useSession } from "@/lib/auth";

interface AdminUser {
  id: number;
  name?: string;
  email: string;
  role: string;
  is_approved: boolean;
}

const fetcher = (url: string) =>
  fetch(url, {
    credentials: "include",
    cache: "no-store",
    headers: { "Cache-Control": "no-cache, no-store, must-revalidate", Pragma: "no-cache" },
  }).then(async (res) => {
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.detail || `HTTP error! status: ${res.status}`);
    }
    return res.json();
  });

const ADMIN_WHITELIST = ["adryanrifky716@gmail.com", "amadryanrifky@gmail.com", "adryan@gmail.com", "arya@gmail.com"];

export default function AdminPage() {
  const { data: session } = useSession();

  const userEmail = (session?.user?.email || "").trim().toLowerCase();
  const isOwner = ADMIN_WHITELIST.includes(userEmail) || userEmail.includes("adryan") || userEmail.includes("arya");
  const isAdmin = session?.user?.role === "admin" || isOwner;

  const { data: users = [], isLoading, mutate: refreshUsers } = useSWR<AdminUser[]>(isAdmin ? "/api/users" : null, fetcher, {
    revalidateOnFocus: true,
    dedupingInterval: 0,
  });

  const { data: gateSetting, mutate: refreshGateSetting } = useSWR<{ require_login_full_feature: boolean }>(
    isAdmin ? "/api/users/settings/auth-gate" : null,
    fetcher,
    { revalidateOnFocus: false, dedupingInterval: 0 }
  );

  const [isUpdatingGate, setIsUpdatingGate] = useState(false);
  const [actionMessage, setActionMessage] = useState<string | null>(null);

  const isGateActive = gateSetting?.require_login_full_feature ?? true;

  const toggleAuthGate = async () => {
    if (isUpdatingGate) return;
    const newStatus = !isGateActive;
    setIsUpdatingGate(true);
    refreshGateSetting({ require_login_full_feature: newStatus }, false);

    try {
      const res = await fetch("/api/users/settings/auth-gate", {
        method: "PATCH",
        credentials: "include",
        headers: { "Content-Type": "application/json", "Cache-Control": "no-cache, no-store, must-revalidate", Pragma: "no-cache" },
        body: JSON.stringify({ require_login_full_feature: newStatus, admin_email: session?.user?.email || "" }),
      });

      if (res.ok) {
        const data = await res.json();
        await refreshGateSetting(data, false);
        setActionMessage(
          newStatus
            ? "Sukses: Opsi wajib login & persetujuan admin untuk full fitur AKTIF."
            : "Sukses: Opsi dinonaktifkan. Pengunjung publik dapat mengakses full fitur tanpa login."
        );
        setTimeout(() => setActionMessage(null), 4000);
      } else {
        const err = await res.json().catch(() => ({}));
        alert(err.detail || "Gagal mengubah pengaturan");
        await refreshGateSetting();
      }
    } catch {
      alert("Terjadi kesalahan koneksi saat mengubah pengaturan.");
      await refreshGateSetting();
    } finally {
      setIsUpdatingGate(false);
    }
  };

  const toggleApproval = async (id: number, currentStatus: boolean) => {
    const nextStatus = !currentStatus;
    refreshUsers(users.map((u) => (u.id === id ? { ...u, is_approved: nextStatus } : u)), false);
    try {
      const res = await fetch(`/api/users/${id}/approve`, {
        method: "PATCH",
        credentials: "include",
        headers: { "Content-Type": "application/json", "Cache-Control": "no-cache, no-store, must-revalidate", Pragma: "no-cache" },
        body: JSON.stringify({ is_approved: nextStatus }),
      });
      if (res.ok) {
        const updated = await res.json();
        refreshUsers(users.map((u) => (u.id === id ? updated : u)), false);
      } else {
        const err = await res.json().catch(() => ({}));
        alert(err.detail || "Gagal mengubah status persetujuan");
        refreshUsers();
      }
    } catch {
      alert("Terjadi kesalahan koneksi.");
      refreshUsers();
    }
  };

  const deleteUser = async (id: number, email: string) => {
    if (!window.confirm(`Apakah Anda yakin ingin menghapus akun ${email} secara permanen?`)) return;
    const previousUsers = users;
    refreshUsers(users.filter((u) => u.id !== id), false);
    try {
      const res = await fetch(`/api/users/${id}`, {
        method: "DELETE",
        credentials: "include",
        headers: { "Cache-Control": "no-cache, no-store, must-revalidate", Pragma: "no-cache" },
      });
      if (res.ok) {
        await refreshUsers();
      } else {
        const err = await res.json().catch(() => ({}));
        alert(err.detail || "Gagal menghapus user");
        refreshUsers(previousUsers, false);
      }
    } catch {
      alert("Terjadi kesalahan koneksi saat menghapus user.");
      refreshUsers(previousUsers, false);
    }
  };

  if (!session || !isAdmin) {
    return (
      <div className="flex flex-col items-center justify-center h-[70vh] text-[var(--color-negative)]">
        <Icon icon="ph:shield-warning-duotone" width="64" className="mb-4" />
        <h1 className="text-xl font-bold">Akses Ditolak</h1>
        <p className="text-sm">Halaman ini khusus untuk Administrator.</p>
      </div>
    );
  }

  return (
    <div className="p-4 sm:p-6 md:p-8 max-w-5xl mx-auto space-y-8">
      {actionMessage && (
        <div className="p-4 rounded-2xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-500 text-sm font-semibold flex items-center gap-2">
          <Icon icon="ph:check-circle-duotone" width="20" />
          <span>{actionMessage}</span>
        </div>
      )}

      <div className="bg-[var(--md-sys-color-surface-container)] rounded-[24px] border border-[var(--md-sys-color-outline-variant)] p-5 sm:p-6 shadow-sm">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="flex items-start sm:items-center gap-3.5">
            <div className={`w-12 h-12 rounded-2xl flex items-center justify-center flex-shrink-0 transition-colors shadow-xs ${isGateActive ? "bg-amber-500/15 text-amber-500" : "bg-emerald-500/15 text-emerald-500"}`}>
              <Icon icon={isGateActive ? "ph:lock-key-duotone" : "ph:lock-key-open-duotone"} width="26" />
            </div>
            <div>
              <div className="flex flex-wrap items-center gap-2">
                <h2 className="text-lg font-bold text-[var(--md-sys-color-on-surface)]">Login Akses Full Fitur</h2>
                <span className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold border transition-colors ${isGateActive ? "bg-amber-500/10 text-amber-500 border-amber-500/20" : "bg-emerald-500/10 text-emerald-500 border-emerald-500/20"}`}>
                  <span className={`w-2 h-2 rounded-full ${isGateActive ? "bg-amber-500" : "bg-emerald-500"}`} />
                  {isGateActive ? "Wajib Login & Disetujui" : "Akses Terbuka (Bebas Login)"}
                </span>
              </div>
              <p className="text-xs sm:text-sm text-[var(--md-sys-color-on-surface-variant)] mt-1 max-w-2xl leading-relaxed">
                {isGateActive
                  ? "User harus login dan akun disetujui lewat panel admin agar bisa akses full fitur (Signal, Foreign Flow, Backtest, Radar Nyangkut). Pengunjung tanpa akun yang disetujui dibatasi hanya pada Home dan Detail Saham."
                  : "Pengunjung dapat mengakses full fitur aplikasi InvestOwl secara bebas tanpa harus login."}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-3 sm:self-center">
            <button
              type="button"
              onClick={toggleAuthGate}
              disabled={isUpdatingGate}
              className={`relative inline-flex h-8 w-14 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none disabled:opacity-50 ${isGateActive ? "bg-[var(--md-sys-color-primary)]" : "bg-neutral-600 dark:bg-neutral-700"}`}
              aria-label="Toggle login akses full fitur"
            >
              <span className={`pointer-events-none inline-flex items-center justify-center h-7 w-7 transform rounded-full bg-white text-neutral-800 shadow-md ring-0 transition duration-200 ease-in-out ${isGateActive ? "translate-x-6" : "translate-x-0"}`}>
                <Icon icon={isUpdatingGate ? "eos-icons:loading" : isGateActive ? "ph:lock-fill" : "ph:lock-open-fill"} width="14" className={isGateActive ? "text-[var(--md-sys-color-primary)]" : "text-neutral-600"} />
              </span>
            </button>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <div className="w-12 h-12 rounded-2xl bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] flex items-center justify-center shadow-sm">
          <Icon icon="ph:users-three-duotone" width="28" />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-[var(--md-sys-color-on-surface)] tracking-tight">User Management</h1>
          <p className="text-sm text-[var(--md-sys-color-on-surface-variant)] mt-1">Kelola akun dan persetujuan akses pengguna InvestOwl.</p>
        </div>
      </div>

      <div className="bg-[var(--md-sys-color-surface-container)] rounded-[24px] border border-[var(--md-sys-color-outline-variant)] overflow-hidden shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] text-sm">
                <th className="p-4 font-semibold whitespace-nowrap">Pengguna</th>
                <th className="p-4 font-semibold text-center whitespace-nowrap">Role</th>
                <th className="p-4 font-semibold text-center whitespace-nowrap">Status Akses</th>
                <th className="p-4 font-semibold text-right whitespace-nowrap">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]">
              {isLoading ? (
                <tr><td colSpan={4} className="p-8 text-center text-[var(--md-sys-color-on-surface-variant)]"><Icon icon="eos-icons:loading" width="24" className="mx-auto" /></td></tr>
              ) : users.length === 0 ? (
                <tr><td colSpan={4} className="p-8 text-center text-[var(--md-sys-color-on-surface-variant)]">Belum ada pengguna terdaftar.</td></tr>
              ) : (
                users.map((user) => {
                  const isCurrentUser = user.email.toLowerCase() === userEmail;
                  return (
                    <tr key={user.id} className="hover:bg-[var(--md-sys-color-surface-container-highest)] transition-colors">
                      <td className="p-4">
                        <div className="flex flex-col">
                          <div className="font-semibold text-[var(--md-sys-color-on-surface)] whitespace-nowrap flex items-center gap-2">
                            <span>{user.name || "Tanpa Nama"}</span>
                            {isCurrentUser && <span className="text-[10px] px-1.5 py-0.2 rounded bg-neutral-500/20 text-[var(--md-sys-color-primary)] font-mono">Anda</span>}
                          </div>
                          <div className="text-xs text-[var(--md-sys-color-on-surface-variant)] whitespace-nowrap">{user.email}</div>
                        </div>
                      </td>
                      <td className="p-4 text-center">
                        <span className={`px-3 py-1 rounded-full text-xs font-bold ${user.role === "admin" ? "bg-purple-500/10 text-purple-500" : "bg-blue-500/10 text-blue-500"}`}>{user.role.toUpperCase()}</span>
                      </td>
                      <td className="p-4 text-center">
                        <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold ${user.is_approved ? "bg-emerald-500/10 text-emerald-500" : "bg-amber-500/10 text-amber-500"}`}>
                          <span className={`w-2 h-2 rounded-full ${user.is_approved ? "bg-emerald-500" : "bg-amber-500"}`} />
                          {user.is_approved ? "DISETUJUI" : "PENDING"}
                        </span>
                      </td>
                      <td className="p-4 text-right">
                        {!isCurrentUser && user.role !== "admin" && (
                          <div className="flex items-center justify-end gap-2">
                            <button
                              onClick={() => toggleApproval(user.id, user.is_approved)}
                              className={`px-3.5 py-1.5 rounded-full text-xs font-semibold transition-all active:scale-95 whitespace-nowrap cursor-pointer ${user.is_approved ? "bg-red-500/15 text-red-400 hover:bg-red-500 hover:text-white" : "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] hover:opacity-90 shadow-sm"}`}
                            >
                              {user.is_approved ? "Cabut Akses" : "Izinkan"}
                            </button>
                            <button onClick={() => deleteUser(user.id, user.email)} className="p-1.5 rounded-full text-[var(--md-sys-color-on-surface-variant)] hover:text-red-500 hover:bg-red-500/10 transition-colors cursor-pointer" title={`Hapus pengguna ${user.email}`}>
                              <Icon icon="ph:trash-duotone" width="18" />
                            </button>
                          </div>
                        )}
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
