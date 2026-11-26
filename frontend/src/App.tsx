import React, { Suspense, lazy } from "react";
import { BrowserRouter, Navigate, Route, Routes, useLocation } from "react-router-dom";
import useSWR from "swr";
import { AuthProvider, useSession } from "@/lib/auth";
import ClientLayout from "@/components/layout/ClientLayout";

const LoginPage = lazy(() => import("@/features/auth/LoginPage"));
const PendingPage = lazy(() => import("@/features/auth/PendingPage"));
const AdminPage = lazy(() => import("@/features/admin/AdminPage"));
const MarketReportPage = lazy(() => import("@/features/market-report/MarketReportPage"));
const HomePage = lazy(() => import("@/features/home/HomePage"));
const TickerPage = lazy(() => import("@/features/ticker/TickerPage"));
const SignalPage = lazy(() => import("@/features/signal/SignalPage"));
const ForeignPage = lazy(() => import("@/features/foreign/ForeignPage"));
const OwnershipPage = lazy(() => import("@/features/ownership/OwnershipPage"));
const GroupPage = lazy(() => import("@/features/group/GroupPage"));
const ExplorerPage = lazy(() => import("@/features/explorer/ExplorerPage"));
const MmdPage = lazy(() => import("@/features/mmd/MmdPage"));
const IpoHunterPage = lazy(() => import("@/features/ipo-hunter/IpoHunterPage"));
const ForeignProPage = lazy(() => import("@/features/foreign-pro/ForeignProPage"));
const KongloPage = lazy(() => import("@/features/konglo/KongloPage"));

const PROTECTED_PREFIXES = [
  "/signal",
  "/foreign",
  "/backtest",
  "/ownership",
  "/group",
  "/explorer",
  "/mmd",
  "/ipo-hunter",
  "/foreign-pro",
  "/konglo",
];

const gateFetcher = (url: string) =>
  fetch(url, { cache: "no-store" }).then((r) => (r.ok ? r.json() : { require_login_full_feature: true }));

function RouteGuard({ children }: { children: React.ReactNode }) {
  const { pathname } = useLocation();
  const { data: session, status } = useSession();
  const { data: gateData } = useSWR<{ require_login_full_feature: boolean }>("/api/users/settings/auth-gate", gateFetcher, {
    dedupingInterval: 5000,
    revalidateOnFocus: false,
  });

  const isAdminRoute = pathname.startsWith("/admin");
  const isProtected = PROTECTED_PREFIXES.some((p) => pathname === p || pathname.startsWith(`${p}/`));
  const isGateActive = gateData?.require_login_full_feature ?? true;
  const isApproved = session?.user?.is_approved === true || session?.user?.role === "admin";
  const isAdmin = session?.user?.role === "admin";

  if (status === "loading" && (isAdminRoute || isProtected || pathname === "/pending")) {
    return (
      <div className="max-w-7xl mx-auto p-8 text-center text-xs font-mono text-[var(--md-sys-color-on-surface-variant)] animate-pulse">
        Memeriksa sesi...
      </div>
    );
  }

  if (isAdminRoute) {
    if (!session) return <Navigate to="/login" replace />;
    if (!isAdmin) return <Navigate to="/" replace />;
  }

  if (isGateActive && isProtected) {
    if (!session) return <Navigate to="/login" replace />;
    if (!isApproved) return <Navigate to="/pending" replace />;
  }

  if (pathname === "/pending" && (!isGateActive || isApproved)) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}

const PageLoader = () => (
  <div className="max-w-7xl mx-auto p-8 text-center text-xs font-mono text-[var(--md-sys-color-on-surface-variant)] animate-pulse">
    Memuat halaman...
  </div>
);

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <ClientLayout>
          <Suspense fallback={<PageLoader />}>
            <RouteGuard>
              <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/login" element={<LoginPage />} />
                <Route path="/pending" element={<PendingPage />} />
                <Route path="/admin" element={<AdminPage />} />
                <Route path="/admin/market-report" element={<MarketReportPage />} />
                <Route path="/signal" element={<SignalPage />} />
                <Route path="/backtest" element={<Navigate to="/signal?tab=backtest" replace />} />
                <Route path="/foreign" element={<ForeignPage />} />
                <Route path="/ownership" element={<OwnershipPage />} />
                <Route path="/group" element={<GroupPage />} />
                <Route path="/explorer" element={<ExplorerPage />} />
                <Route path="/mmd" element={<MmdPage />} />
                <Route path="/ipo-hunter" element={<IpoHunterPage />} />
                <Route path="/foreign-pro" element={<ForeignProPage />} />
                <Route path="/konglo" element={<KongloPage />} />
                <Route path="/:ticker" element={<TickerPage />} />
              </Routes>
            </RouteGuard>
          </Suspense>
        </ClientLayout>
      </AuthProvider>
    </BrowserRouter>
  );
}
