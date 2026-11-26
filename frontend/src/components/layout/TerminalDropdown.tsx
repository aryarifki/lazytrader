import { useState, useRef, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { Icon } from "@iconify/react";
import { useSession } from "@/lib/auth";

export const TERMINAL_ITEMS = [
  {
    path: "/mmd",
    label: "MM Detector Pro",
    desc: "Deteksi bandar & klaster broker",
    icon: "ph:radar-duotone",
    adminOnly: false,
  },
  {
    path: "/ipo-hunter",
    label: "IPO Hunter Pro",
    desc: "Underwriter & performa pasca-IPO",
    icon: "ph:rocket-launch-duotone",
    adminOnly: false,
  },
  {
    path: "/foreign-pro",
    label: "Foreign Flow Pro",
    desc: "Screener & akumulasi asing",
    icon: "ph:currency-dollar-duotone",
    adminOnly: false,
  },
  {
    path: "/konglo",
    label: "Konglo Screener",
    desc: "Kinerja grup konglomerasi",
    icon: "ph:users-three-duotone",
    adminOnly: false,
  },
  {
    path: "/admin/market-report",
    label: "Market Report",
    desc: "Generator laporan Sesi 1 & 2",
    icon: "ph:file-text-duotone",
    adminOnly: true,
  },
];

export default function TerminalDropdown() {
  const { pathname } = useLocation();
  const { data: session } = useSession();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  const isAdmin = session?.user?.role === "admin";
  const visibleItems = TERMINAL_ITEMS.filter((item) => !item.adminOnly || isAdmin);
  const isTerminalActive = TERMINAL_ITEMS.some(
    (item) => pathname === item.path || pathname.startsWith(`${item.path}/`)
  );

  useEffect(() => {
    const onDocClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", onDocClick);
    return () => document.removeEventListener("mousedown", onDocClick);
  }, []);

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold tracking-tight transition-all duration-150 cursor-pointer ${
          isTerminalActive
            ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] shadow-2xs font-bold"
            : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] hover:bg-[var(--md-sys-color-surface-container-high)]"
        }`}
        aria-expanded={open}
        aria-haspopup="menu"
      >
        <Icon icon={isTerminalActive ? "ph:terminal-window-fill" : "ph:terminal-window-duotone"} width="15" />
        <span>Terminal</span>
        <Icon icon={open ? "ph:caret-up-bold" : "ph:caret-down-bold"} width="11" className="opacity-75" />
      </button>

      {open && (
        <div
          role="menu"
          className="absolute right-0 mt-2 w-64 rounded-2xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)] shadow-xl p-1.5 z-50"
        >
          {visibleItems.map((item) => {
            const active = pathname === item.path || pathname.startsWith(`${item.path}/`);
            return (
              <Link
                key={item.path}
                to={item.path}
                onClick={() => setOpen(false)}
                className={`flex items-start gap-2.5 px-3 py-2 rounded-xl transition-colors ${
                  active
                    ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)]"
                    : "text-[var(--md-sys-color-on-surface)] hover:bg-[var(--md-sys-color-surface-container-high)]"
                }`}
              >
                <Icon icon={item.icon} width="17" className="mt-0.5 shrink-0 text-[var(--md-sys-color-primary)]" />
                <div className="min-w-0">
                  <div className="text-xs font-bold leading-tight">{item.label}</div>
                  <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] truncate">{item.desc}</div>
                </div>
              </Link>
            );
          })}
        </div>
      )}
    </div>
  );
}
