import React from "react";
import { Link, useLocation } from "react-router-dom";
import { Icon } from "@iconify/react";
import BottomNav from "./BottomNav";
import RefreshButton from "./RefreshButton";
import Sidebar from "./Sidebar";
import TerminalDropdown from "./TerminalDropdown";
import ThemeToggle from "./ThemeToggle";
import LogoutButton from "./LogoutButton";
import CacheLoadingOverlay from "@/components/ui/CacheLoadingOverlay";
import { useAppStore } from "@/store/useAppStore";

export default function ClientLayout({ children }: { children: React.ReactNode }) {
  const { pathname } = useLocation();
  const { activeTicker, sidebarOpen, setSidebarOpen } = useAppStore();
  const isAuthPage = pathname === "/login" || pathname === "/pending";

  const isHome = pathname === "/";
  const isSignal = pathname.includes("/signal") || pathname.includes("/backtest");
  const isForeign = pathname === "/foreign";
  const isOwnership = pathname.startsWith("/ownership");
  const isGroup = pathname.startsWith("/group");
  const isExplorer = pathname.startsWith("/explorer");
  const isTerminal =
    pathname.startsWith("/mmd") ||
    pathname.startsWith("/ipo-hunter") ||
    pathname.startsWith("/foreign-pro") ||
    pathname.startsWith("/konglo") ||
    pathname.startsWith("/admin/market-report");
  const isResearch = isOwnership || isGroup || isExplorer;
  const isDetail = !isHome && !isSignal && !isForeign && !isResearch && !isTerminal && !pathname.includes("/admin");

  const detailPath = `/${activeTicker || "BBCA"}`;

  const desktopNavLinkClass = (active: boolean) =>
    `inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold tracking-tight transition-all duration-150 ${
      active
        ? "bg-[var(--md-sys-color-primary-container)] text-[var(--md-sys-color-on-primary-container)] shadow-2xs font-bold"
        : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] hover:bg-[var(--md-sys-color-surface-container-high)]"
    }`;

  return (
    <>
      <CacheLoadingOverlay />

      <header className="sticky top-0 z-40 w-full border-b border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)]/90 backdrop-blur-md transition-colors duration-200">
        <div className="flex items-center justify-between px-3 sm:px-6 h-14 max-w-7xl mx-auto">
          <div className="flex items-center gap-3">
            {!isAuthPage && (
              <button
                onClick={() => setSidebarOpen(!sidebarOpen)}
                className="p-1.5 rounded-lg bg-[var(--md-sys-color-surface-container-high)] hover:bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)] text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] transition-all cursor-pointer shadow-2xs"
                title="Buka Panel Parameter & Model"
                aria-label="Toggle Controls Panel"
              >
                <Icon icon="ph:sliders-horizontal-duotone" width="18" />
              </button>
            )}

            <Link to="/" className="flex items-center gap-2.5 group">
              <div className="w-8.5 h-8.5 relative flex-shrink-0 flex items-center justify-center rounded-xl bg-[var(--md-sys-color-surface-container-high)]/70 dark:bg-white/[0.08] border border-[var(--md-sys-color-outline-variant)]/60 dark:border-white/15 p-1 shadow-2xs group-hover:border-[var(--md-sys-color-primary)]/60 transition-all">
                <img
                  src="/logo-v2.png"
                  alt="Investowl Logo"
                  className="w-full h-full object-contain p-0.5 transition-all dark:brightness-110 dark:drop-shadow-[0_0_6px_rgba(255,184,121,0.35)]"
                />
              </div>
              <div className="flex items-baseline gap-1.5">
                <span className="font-extrabold text-base tracking-tight text-[var(--md-sys-color-on-surface)] group-hover:text-[var(--md-sys-color-primary)] transition-colors">
                  Investowl
                </span>
                <span className="hidden sm:inline-block text-[9px] font-mono font-bold px-1.5 py-0.5 rounded bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-primary)] border border-[var(--md-sys-color-outline-variant)] uppercase tracking-wider">
                  Quant
                </span>
              </div>
            </Link>
          </div>

          {!isAuthPage && (
            <nav className="hidden md:flex items-center gap-1">
              <Link to="/" className={desktopNavLinkClass(isHome)}>
                <Icon icon={isHome ? "ph:house-fill" : "ph:house-duotone"} width="15" />
                <span>Home</span>
              </Link>
              <Link to={detailPath} className={desktopNavLinkClass(isDetail)}>
                <Icon icon={isDetail ? "ph:chart-line-up-fill" : "ph:chart-line-up-duotone"} width="15" />
                <span>Emiten</span>
              </Link>
              <Link to="/signal" className={desktopNavLinkClass(isSignal)}>
                <Icon icon={isSignal ? "ph:broadcast-fill" : "ph:broadcast-duotone"} width="15" />
                <span>Signal</span>
              </Link>
              <Link to="/foreign" className={desktopNavLinkClass(isForeign)}>
                <Icon icon={isForeign ? "ph:globe-stand-fill" : "ph:globe-stand-duotone"} width="15" />
                <span>Foreign</span>
              </Link>
              <Link to="/ownership" className={desktopNavLinkClass(isOwnership)} title="Ownership">
                <Icon icon={isOwnership ? "ph:chart-pie-slice-fill" : "ph:chart-pie-slice-duotone"} width="15" />
                <span className="hidden lg:inline">Ownership</span>
              </Link>
              <Link to="/group" className={desktopNavLinkClass(isGroup)} title="Konglomerasi & Jejaring">
                <Icon icon={isGroup ? "ph:share-network-fill" : "ph:share-network-duotone"} width="15" />
                <span className="hidden lg:inline">Group</span>
              </Link>
              <Link to="/explorer" className={desktopNavLinkClass(isExplorer)} title="Explorer Kepemilikan">
                <Icon icon={isExplorer ? "ph:funnel-fill" : "ph:funnel-duotone"} width="15" />
                <span className="hidden lg:inline">Explorer</span>
              </Link>
              <TerminalDropdown />
            </nav>
          )}

          <div className="flex items-center gap-2">
            {!isAuthPage && <RefreshButton />}
            <ThemeToggle />
            {!isAuthPage && <LogoutButton />}
          </div>
        </div>
      </header>

      <div className="flex min-h-[calc(100vh-3.5rem)] pb-20 md:pb-8">
        {!isAuthPage && <Sidebar />}
        <div className="flex-1 min-w-0">{children}</div>
      </div>

      {!isAuthPage && <BottomNav />}
    </>
  );
}
