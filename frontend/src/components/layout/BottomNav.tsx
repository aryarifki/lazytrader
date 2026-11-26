import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { Icon } from "@iconify/react";
import { useAppStore } from "@/store/useAppStore";
import { useSession } from "@/lib/auth";
import { TERMINAL_ITEMS } from "./TerminalDropdown";

export default function BottomNav() {
  const { pathname } = useLocation();
  const { activeTicker } = useAppStore();
  const { data: session } = useSession();
  const [sheetOpen, setSheetOpen] = useState(false);

  const isAdmin = session?.user?.role === "admin";
  const visibleTerminal = TERMINAL_ITEMS.filter((i) => !i.adminOnly || isAdmin);

  const isHome = pathname === "/";
  const isSignal = pathname.includes("/signal") || pathname.includes("/backtest");
  const isForeign = pathname === "/foreign";
  const isResearch =
    pathname.startsWith("/ownership") ||
    pathname.startsWith("/group") ||
    pathname.startsWith("/explorer");
  const isTerminal = TERMINAL_ITEMS.some((i) => pathname === i.path || pathname.startsWith(`${i.path}/`));
  const isDashboard = !isHome && !isSignal && !isForeign && !isResearch && !isTerminal && !pathname.includes("/admin");

  const targetDashboardPath = `/${activeTicker || "BBCA"}`;

  const navItemClass = (active: boolean) =>
    `group flex flex-col items-center justify-center w-full h-full py-1 transition-all ${
      active
        ? "text-[var(--md-sys-color-primary)] font-bold"
        : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
    }`;

  const iconWrapper = (active: boolean) =>
    `flex items-center justify-center w-10 h-7 rounded-xl transition-all duration-200 ${
      active
        ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)]"
        : "text-[var(--md-sys-color-on-surface-variant)] group-hover:text-[var(--md-sys-color-on-surface)]"
    }`;

  return (
    <>
      {sheetOpen && (
        <div className="md:hidden fixed inset-0 z-50 bg-black/50 backdrop-blur-xs" onClick={() => setSheetOpen(false)}>
          <div
            className="fixed bottom-16 left-3 right-3 rounded-2xl border border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)] p-2 shadow-2xl space-y-1"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="px-3 py-1.5 text-[10px] font-mono font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
              Terminal Modules
            </div>
            {visibleTerminal.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                onClick={() => setSheetOpen(false)}
                className="flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface)]"
              >
                <Icon icon={item.icon} width="18" className="text-[var(--md-sys-color-primary)]" />
                <div>
                  <div className="text-xs font-bold">{item.label}</div>
                  <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">{item.desc}</div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}

      <nav className="md:hidden fixed bottom-0 left-0 right-0 w-full bg-[var(--md-sys-color-surface-container)]/95 backdrop-blur-md border-t border-[var(--md-sys-color-outline-variant)] z-50 transition-colors duration-200">
        <div className="flex justify-around items-center h-16 max-w-lg mx-auto px-1">
          <Link to="/" className={navItemClass(isHome)}>
            <div className={iconWrapper(isHome)}>
              <Icon icon={isHome ? "ph:house-fill" : "ph:house-duotone"} width="20" />
            </div>
            <span className="text-[10px] mt-0.5 font-medium tracking-tight">Home</span>
          </Link>

          <Link to={targetDashboardPath} className={navItemClass(isDashboard)}>
            <div className={iconWrapper(isDashboard)}>
              <Icon icon={isDashboard ? "ph:chart-line-up-fill" : "ph:chart-line-up-duotone"} width="20" />
            </div>
            <span className="text-[10px] mt-0.5 font-medium tracking-tight">Detail</span>
          </Link>

          <Link to="/signal" className={navItemClass(isSignal)}>
            <div className={iconWrapper(isSignal)}>
              <Icon icon={isSignal ? "ph:broadcast-fill" : "ph:broadcast-duotone"} width="20" />
            </div>
            <span className="text-[10px] mt-0.5 font-medium tracking-tight">Signal</span>
          </Link>

          <Link to="/foreign" className={navItemClass(isForeign)}>
            <div className={iconWrapper(isForeign)}>
              <Icon icon={isForeign ? "ph:globe-stand-fill" : "ph:globe-stand-duotone"} width="20" />
            </div>
            <span className="text-[10px] mt-0.5 font-medium tracking-tight">Foreign</span>
          </Link>

          <Link to="/ownership" className={navItemClass(isResearch)}>
            <div className={iconWrapper(isResearch)}>
              <Icon icon={isResearch ? "ph:chart-pie-slice-fill" : "ph:chart-pie-slice-duotone"} width="20" />
            </div>
            <span className="text-[10px] mt-0.5 font-medium tracking-tight">Riset</span>
          </Link>

          <button type="button" onClick={() => setSheetOpen((v) => !v)} className={navItemClass(isTerminal || sheetOpen)}>
            <div className={iconWrapper(isTerminal || sheetOpen)}>
              <Icon icon={isTerminal ? "ph:terminal-window-fill" : "ph:terminal-window-duotone"} width="20" />
            </div>
            <span className="text-[10px] mt-0.5 font-medium tracking-tight">Terminal</span>
          </button>
        </div>
      </nav>
    </>
  );
}
