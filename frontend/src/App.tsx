import React, { Suspense, lazy } from "react";
import { BrowserRouter, Navigate, Route, Routes, useLocation } from "react-router-dom";
import useSWR from "swr";
import { AuthProvider, useSession } from "@/lib/auth";
import ClientLayout from "@/components/layout/ClientLayout";
import FeaturePlaceholder from "@/components/ui/FeaturePlaceholder";

const LoginPage = lazy(() => import("@/features/auth/LoginPage"));
const PendingPage = lazy(() => import("@/features/auth/PendingPage"));
const AdminPage = lazy(() => import("@/features/admin/AdminPage"));

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
                <Route path="/" element={<FeaturePlaceholder title="Home — InvestOwl Dashboard" phase="Fase 4" />} />
                <Route path="/login" element={<LoginPage />} />
                <Route path="/pending" element={<PendingPage />} />
                <Route path="/admin" element={<AdminPage />} />
                <Route path="/admin/market-report" element={<FeaturePlaceholder title="Market Report Generator Admin" phase="Fase 4" />} />
                <Route path="/signal" element={<FeaturePlaceholder title="AI Signal & Backtest" phase="Fase 4" />} />
                <Route path="/backtest" element={<Navigate to="/signal?tab=backtest" replace />} />
                <Route path="/foreign" element={<FeaturePlaceholder title="Foreign Flow & Sector Rotation" phase="Fase 4" />} />
                <Route path="/ownership" element={<FeaturePlaceholder title="Ownership Intelligence" phase="Fase 4" />} />
                <Route path="/group" element={<FeaturePlaceholder title="Konglomerasi & Jejaring" phase="Fase 4" />} />
                <Route path="/explorer" element={<FeaturePlaceholder title="Ownership Explorer" phase="Fase 4" />} />
                <Route path="/mmd" element={<FeaturePlaceholder title="MM Detector Pro" phase="Fase 4" />} />
                <Route path="/ipo-hunter" element={<FeaturePlaceholder title="IPO Hunter Pro" phase="Fase 4" />} />
                <Route path="/foreign-pro" element={<FeaturePlaceholder title="Foreign Flow Pro" phase="Fase 4" />} />
                <Route path="/konglo" element={<FeaturePlaceholder title="Konglo Screener" phase="Fase 4" />} />
                <Route path="/:ticker" element={<FeaturePlaceholder title="Detail Saham & Bandarmologi" phase="Fase 4" />} />
              </Routes>
            </RouteGuard>
          </Suspense>
        </ClientLayout>
      </AuthProvider>
    </BrowserRouter>
  );
}
