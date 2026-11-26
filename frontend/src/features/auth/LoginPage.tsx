import { signIn } from "@/lib/auth";
import { Icon } from "@iconify/react";

export default function LoginPage() {
  return (
    <div className="min-h-[92vh] flex flex-col items-center justify-center p-4 relative">
      <div className="absolute inset-0 bg-terminal-grid pointer-events-none opacity-40" />

      <div className="w-full max-w-lg flex flex-col gap-6 relative z-10 my-8">
        <div className="bg-[var(--md-sys-color-surface-container)] rounded-3xl p-7 sm:p-9 shadow-md border border-[var(--md-sys-color-outline-variant)]/60 relative overflow-hidden backdrop-blur-xs">
          <div className="flex flex-col items-center text-center mb-8 mt-2">
            <div className="w-18 h-18 sm:w-20 sm:h-20 mb-4 relative flex items-center justify-center rounded-2xl bg-[var(--md-sys-color-surface-container-high)]/60 dark:bg-white/[0.08] border border-[var(--md-sys-color-outline-variant)]/60 dark:border-white/15 p-2 shadow-xs">
              <img
                src="/logo-v2.png"
                alt="The Investowl Logo"
                className="w-full h-full object-contain p-1.5 transition-all dark:brightness-110 dark:drop-shadow-[0_0_8px_rgba(255,184,121,0.35)]"
              />
            </div>

            <div className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] text-[10px] font-mono font-bold tracking-wider uppercase mb-2">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
              Quantitative Engine v2.0
            </div>

            <h1 className="text-2xl sm:text-3xl font-extrabold text-[var(--md-sys-color-on-surface)] tracking-tight">
              Investowl Terminal
            </h1>
            <p className="text-xs sm:text-sm font-medium text-[var(--md-sys-color-on-surface-variant)] mt-1.5 max-w-sm">
              Sistem analitik bandarmologi dan arus modal institusi Bursa Efek Indonesia.
            </p>
          </div>

          <div className="space-y-3.5">
            <button
              onClick={() => signIn("google", { callbackUrl: "/pending" })}
              className="w-full flex items-center justify-center gap-3 bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface)] border border-[var(--md-sys-color-outline-variant)] hover:border-[var(--md-sys-color-primary)] rounded-xl px-5 py-3.5 text-sm font-bold transition-all active:scale-[0.99] shadow-xs cursor-pointer"
            >
              <Icon icon="logos:google-icon" width="20" />
              <span>Masuk dengan Google</span>
            </button>

            <div className="flex items-center py-1">
              <div className="flex-grow border-t border-[var(--md-sys-color-outline-variant)]" />
              <span className="flex-shrink-0 mx-3 text-[10px] font-mono font-bold tracking-[0.2em] text-[var(--md-sys-color-on-surface-variant)] uppercase">
                Request Akses Persetujuan
              </span>
              <div className="flex-grow border-t border-[var(--md-sys-color-outline-variant)]" />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <a
                href="https://wa.me/6287882725845?text=Halo%20Admin%20InvestOwl,%20saya%20ingin%20request%20persetujuan%20akses%20akun%20saya."
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center justify-center gap-2 bg-[#25D366]/10 hover:bg-[#25D366]/20 text-[#0f5132] dark:text-[#4ade80] border border-[#25D366]/40 dark:border-[#25D366]/35 rounded-xl px-3 sm:px-4 py-3 text-xs sm:text-sm font-bold transition-all active:scale-[0.98] shadow-xs cursor-pointer group"
              >
                <Icon icon="logos:whatsapp-icon" width="20" className="flex-shrink-0 transition-transform group-hover:scale-110" />
                <span className="truncate">WhatsApp</span>
              </a>

              <a
                href="https://t.me/theinvestowl_id"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center justify-center gap-2 bg-[#0088cc]/10 hover:bg-[#0088cc]/20 text-[#0c5460] dark:text-[#38bdf8] border border-[#0088cc]/40 dark:border-[#0088cc]/35 rounded-xl px-3 sm:px-4 py-3 text-xs sm:text-sm font-bold transition-all active:scale-[0.98] shadow-xs cursor-pointer group"
              >
                <Icon icon="logos:telegram" width="20" className="flex-shrink-0 transition-transform group-hover:scale-110" />
                <span className="truncate">Grup Telegram</span>
              </a>
            </div>
          </div>

          <div className="mt-8 pt-6 border-t border-[var(--md-sys-color-outline-variant)]">
            <div className="bg-[var(--md-sys-color-surface-container-low)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 sm:p-5 shadow-2xs">
              <div className="flex items-center gap-2 mb-2.5 text-amber-500 dark:text-amber-400 font-semibold text-xs uppercase tracking-wider">
                <Icon icon="ph:shield-warning-duotone" width="18" />
                <span>Disclaimer Edukasi &amp; Riset</span>
              </div>
              <p className="text-[11px] sm:text-xs leading-relaxed text-[var(--md-sys-color-on-surface-variant)] text-justify">
                The Investowl adalah platform analitik komputasi dan riset kuantitatif berbasis data science dan pemodelan statistik pasar modal. Seluruh metrik, grafik bandarmologi, estimasi arus modal institusi (&quot;Smart Money&quot;), dan sinyal algoritma merupakan estimasi probabilitas matematis dari data transaksi historis Bursa Efek Indonesia (IDX) dan <strong>bukan merupakan rekomendasi, anjuran, atau ajakan membeli/menjual instrumen finansial</strong>. The Investowl bukan penasihat investasi berlisensi (<em>financial advisor</em>). Segala keputusan transaksi pasar modal, potensi keuntungan, maupun risiko kerugian berada sepenuhnya di bawah tanggung jawab mandiri pengguna (<em>Do Your Own Research - DYOR</em>).
              </p>
            </div>
          </div>
        </div>

        <footer className="text-center space-y-2.5 py-3">
          <div className="flex flex-wrap items-center justify-center gap-3 sm:gap-4 text-xs font-medium text-[var(--md-sys-color-on-surface-variant)]">
            <a href="mailto:adryanrifky716@gmail.com" className="flex items-center gap-1.5 hover:text-[var(--md-sys-color-primary)] transition-colors">
              <Icon icon="ph:envelope-simple-duotone" width="16" />
              <span>adryanrifky716@gmail.com</span>
            </a>
            <span className="opacity-30">•</span>
            <a href="https://t.me/theinvestowl_id" target="_blank" rel="noopener noreferrer" className="flex items-center gap-1.5 hover:text-[#0088cc] dark:hover:text-[#38bdf8] transition-colors">
              <Icon icon="ph:telegram-logo-duotone" width="16" />
              <span>Grup Telegram</span>
            </a>
            <span className="opacity-30">•</span>
            <a href="https://www.instagram.com/rifqiaarya?stkn=dzl1d2RpczIyd2ts" target="_blank" rel="noopener noreferrer" className="flex items-center gap-1.5 hover:text-pink-500 transition-colors">
              <Icon icon="ph:instagram-logo-duotone" width="16" />
              <span>@rifqiaarya</span>
            </a>
          </div>

          <div className="text-[11px] font-mono text-[var(--md-sys-color-on-surface-variant)] opacity-85">
            Proprietary Academic &amp; Research License &bull; Copyright &copy; 2026 Adryan Rifky. All rights reserved.
          </div>
        </footer>
      </div>
    </div>
  );
}
