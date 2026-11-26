import { useNavigate, useSearchParams } from "react-router-dom";
import { Icon } from "@iconify/react";

import ResearchTabs from "@/components/research/ResearchTabs";
import StockOwnershipTab from "./stock/StockOwnershipTab";
import InvestorOwnershipTab from "./investor/InvestorOwnershipTab";
import IhsgTab from "./ihsg/IhsgTab";

const TABS = [
  {
    id: "saham",
    label: "By Saham",
    description: "Profil, komposisi, tren",
    icon: "ph:chart-pie-slice-duotone",
  },
  {
    id: "investor",
    label: "By Investor",
    description: "Portofolio tercatat",
    icon: "ph:users-three-duotone",
  },
  {
    id: "ihsg",
    label: "IHSG",
    description: "Pasar & rotasi sektor",
    icon: "ph:chart-line-up-duotone",
  },
] as const;

type TabId = (typeof TABS)[number]["id"];

function isTabId(value: string | null): value is TabId {
  return value === "saham" || value === "investor" || value === "ihsg";
}

export default function OwnershipDashboard() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const tabParam = searchParams.get("tab");
  const activeTab: TabId = isTabId(tabParam) ? tabParam : "saham";

  const handleTabChange = (id: TabId) => {
    const params = new URLSearchParams(window.location.search);
    params.set("tab", id);
    navigate(`/ownership?${params.toString()}`, {
      replace: true,
      preventScrollReset: true,
    });
  };

  return (
    <div className="max-w-7xl mx-auto space-y-4">
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-3">
        <div className="flex items-center gap-3.5">
          <div className="w-11 h-11 sm:w-12 sm:h-12 rounded-2xl bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0 border border-[var(--md-sys-color-outline-variant)] shadow-2xs">
            <Icon icon="ph:chart-pie-slice-duotone" width="26" />
          </div>
          <div>
            <div className="inline-flex items-center gap-2 text-[9px] font-mono font-bold text-[var(--md-sys-color-primary)] uppercase tracking-wider mb-0.5">
              <span>OWNERSHIP INTELLIGENCE</span>
              <span className="opacity-40">•</span>
              <span>KSEI + IDX</span>
            </div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight">
              Ownership & Market Structure
            </h1>
          </div>
        </div>
        <ResearchTabs />
      </div>

      <div className="border-b border-[var(--md-sys-color-outline-variant)]">
        <div className="flex gap-1.5 overflow-x-auto scrollbar-none pb-2">
          {TABS.map((tab) => {
            const active = activeTab === tab.id;
            return (
              <button
                key={tab.id}
                type="button"
                onClick={() => handleTabChange(tab.id)}
                className={`flex items-center gap-2 rounded-lg border px-3.5 py-2 text-xs font-bold transition-all whitespace-nowrap ${
                  active
                    ? "border-[var(--md-sys-color-primary)] bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-sm"
                    : "border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface-variant)] hover:bg-[var(--md-sys-color-surface-container-high)] hover:text-[var(--md-sys-color-on-surface)]"
                }`}
                title={tab.description}
              >
                <Icon icon={tab.icon} width="16" />
                <span>{tab.label}</span>
              </button>
            );
          })}
        </div>
      </div>

      <div className="transition-colors duration-200">
        {activeTab === "saham" ? <StockOwnershipTab /> : null}
        {activeTab === "investor" ? <InvestorOwnershipTab /> : null}
        {activeTab === "ihsg" ? <IhsgTab /> : null}
      </div>
    </div>
  );
}
