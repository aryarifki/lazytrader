import { Icon } from "@iconify/react";
import ForeignExodusRadar from "./ForeignExodusRadar";
import WatchlistSection from "./WatchlistSection";

export default function HomePage() {
  return (
    <div className="min-h-screen bg-[var(--md-sys-color-surface)] text-[var(--md-sys-color-on-surface)] transition-colors duration-200">
      <main className="max-w-7xl mx-auto p-4 sm:p-6 lg:p-8 space-y-6">
        <ForeignExodusRadar />

        <div className="flex items-center gap-2.5 px-1.5 -mt-3 sm:-mt-2">
          <span className="text-xs font-semibold text-[var(--md-sys-color-on-surface-variant)] select-none">
            Connect with me
          </span>

          <div className="flex items-center gap-2">
            <a
              href="mailto:adryanrifky716@gmail.com"
              title="Email: adryanrifky716@gmail.com"
              className="w-8 h-8 rounded-xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/50 flex items-center justify-center text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-primary)] hover:border-[var(--md-sys-color-primary)]/50 hover:bg-[var(--md-sys-color-surface-container-high)] transition-all duration-200 shadow-2xs"
            >
              <Icon icon="ph:envelope-simple-duotone" width="16" />
            </a>
            <a
              href="https://t.me/theinvestowl_id"
              target="_blank"
              rel="noopener noreferrer"
              title="Komunitas Telegram: @theinvestowl_id"
              className="w-8 h-8 rounded-xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/50 flex items-center justify-center text-[var(--md-sys-color-on-surface-variant)] hover:text-sky-500 hover:border-sky-500/50 hover:bg-[var(--md-sys-color-surface-container-high)] transition-all duration-200 shadow-2xs"
            >
              <Icon icon="ph:telegram-logo-duotone" width="16" />
            </a>
            <a
              href="https://www.instagram.com/rifqiaarya?dlrf=bzJzbzZhZW1qanFr"
              target="_blank"
              rel="noopener noreferrer"
              title="Instagram: @rifqiaarya"
              className="w-8 h-8 rounded-xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/50 flex items-center justify-center text-[var(--md-sys-color-on-surface-variant)] hover:text-pink-500 hover:border-pink-500/50 hover:bg-[var(--md-sys-color-surface-container-high)] transition-all duration-200 shadow-2xs"
            >
              <Icon icon="ph:instagram-logo-duotone" width="16" />
            </a>
            <a
              href="https://github.com/aryarifki"
              target="_blank"
              rel="noopener noreferrer"
              title="GitHub: aryarifki"
              className="w-8 h-8 rounded-xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)]/50 flex items-center justify-center text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] hover:border-[var(--md-sys-color-outline)] hover:bg-[var(--md-sys-color-surface-container-high)] transition-all duration-200 shadow-2xs"
            >
              <Icon icon="ph:github-logo-duotone" width="16" />
            </a>
          </div>
        </div>

        <WatchlistSection />
      </main>
    </div>
  );
}
