import { signOut, useSession } from "@/lib/auth";
import { Icon } from "@iconify/react";

export default function PendingPage() {
  const { data: session } = useSession();

  const waText = session?.user?.email
    ? `Halo Admin InvestOwl, saya telah melakukan login dengan email ${session.user.email}. Mohon untuk menyetujui (approve) akses saya ke dalam Terminal.`
    : "Halo Admin InvestOwl, mohon persetujuan akses saya ke dalam Terminal.";

  const waLink = `https://wa.me/6287882725845?text=${encodeURIComponent(waText)}`;

  return (
    <div className="min-h-[100dvh] bg-[var(--md-sys-color-surface)] text-[var(--md-sys-color-on-surface)] flex items-center justify-center p-4 transition-colors duration-300">
      <div className="w-full max-w-md bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/60 rounded-3xl p-8 shadow-xl flex flex-col items-center text-center relative overflow-hidden">
        <div className="w-16 h-16 rounded-full bg-amber-500/10 flex items-center justify-center border border-amber-500/20 mb-5">
          <Icon icon="ph:hourglass-medium-duotone" className="text-amber-500" width="32" />
        </div>

        <h1 className="text-xl font-bold tracking-tight text-[var(--md-sys-color-on-surface)] mb-2">
          Akses Menunggu Persetujuan
        </h1>
        <p className="text-sm text-[var(--md-sys-color-on-surface-variant)] mb-6 leading-relaxed">
          Akun Anda (<strong className="text-[var(--md-sys-color-on-surface)] font-mono">{session?.user?.email}</strong>) telah berhasil disinkronisasi, namun belum diizinkan oleh Administrator untuk masuk ke dalam terminal analitik.
        </p>

        <div className="w-full space-y-3">
          <div className="grid grid-cols-2 gap-2.5 w-full">
            <a
              href={waLink}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center justify-center gap-2 bg-[#25D366]/10 hover:bg-[#25D366]/20 text-[#0f5132] dark:text-[#4ade80] border border-[#25D366]/40 dark:border-[#25D366]/35 px-3 py-3 rounded-xl font-bold text-xs sm:text-sm transition-all active:scale-[0.98] shadow-xs cursor-pointer group"
            >
              <Icon icon="logos:whatsapp-icon" width="18" className="flex-shrink-0 transition-transform group-hover:scale-110" />
              <span className="truncate">WhatsApp</span>
            </a>

            <a
              href="https://t.me/theinvestowl_id"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center justify-center gap-2 bg-[#0088cc]/10 hover:bg-[#0088cc]/20 text-[#0c5460] dark:text-[#38bdf8] border border-[#0088cc]/40 dark:border-[#0088cc]/35 px-3 py-3 rounded-xl font-bold text-xs sm:text-sm transition-all active:scale-[0.98] shadow-xs cursor-pointer group"
            >
              <Icon icon="logos:telegram" width="18" className="flex-shrink-0 transition-transform group-hover:scale-110" />
              <span className="truncate">Grup Telegram</span>
            </a>
          </div>

          <button
            onClick={() => signOut({ callbackUrl: "/login" })}
            className="w-full flex items-center justify-center gap-2 bg-transparent text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--color-negative)] px-4 py-2.5 rounded-full font-semibold text-xs sm:text-sm transition-all cursor-pointer"
          >
            <Icon icon="ph:sign-out-bold" />
            <span>Ganti Akun (Sign Out)</span>
          </button>
        </div>
      </div>
    </div>
  );
}
