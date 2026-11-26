import { Suspense } from "react";
import GroupDashboard from "./GroupDashboard";

export default function GroupPage() {
  return (
    <main className="min-h-screen bg-[var(--md-sys-color-surface)] text-[var(--md-sys-color-on-surface)] p-4 md:p-6 transition-colors duration-300 pb-24">
      <Suspense
        fallback={
          <div className="max-w-7xl mx-auto p-8 text-center text-xs font-mono text-[var(--md-sys-color-on-surface-variant)] animate-pulse">
            Memuat konglomerasi & jejaring...
          </div>
        }
      >
        <GroupDashboard />
      </Suspense>
    </main>
  );
}
