import { Link, useLocation, useSearchParams } from "react-router-dom";
import { Icon } from "@iconify/react";

const ITEMS = [
  {
    href: "/ownership",
    label: "Ownership",
    description: "Saham, investor, IHSG",
    icon: "ph:chart-pie-slice-duotone",
  },
  {
    href: "/group",
    label: "Group",
    description: "Konglomerasi & jejaring",
    icon: "ph:buildings-duotone",
  },
  {
    href: "/explorer",
    label: "Explorer",
    description: "Screener kepemilikan",
    icon: "ph:funnel-duotone",
  },
] as const;

export default function ResearchTabs() {
  const { pathname } = useLocation();
  const [searchParams] = useSearchParams();

  return (
    <div className="flex items-center gap-1.5 overflow-x-auto scrollbar-none">
      <span className="hidden sm:inline-flex items-center gap-1 text-[9px] font-mono font-bold uppercase tracking-widest text-[var(--md-sys-color-on-surface-variant)] mr-1 shrink-0">
        <Icon icon="ph:compass-duotone" width="13" />
        Riset
      </span>
      {ITEMS.map((item) => {
        const active = pathname === item.href || pathname.startsWith(`${item.href}/`);
        let targetHref = item.href as string;
        if (item.href === "/explorer" && pathname.startsWith("/explorer")) {
          const currentSource = searchParams.get("source");
          if (currentSource) {
            targetHref = `/explorer?source=${currentSource}`;
          }
        }
        return (
          <Link
            key={item.href}
            to={targetHref}
            title={item.description}
            className={`inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1.5 text-[10px] sm:text-[11px] font-bold whitespace-nowrap transition-all ${
              active
                ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] border-[var(--md-sys-color-primary)] shadow-2xs"
                : "bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface-variant)] border-[var(--md-sys-color-outline-variant)] hover:text-[var(--md-sys-color-on-surface)] hover:bg-[var(--md-sys-color-surface-container-high)]"
            }`}
          >
            <Icon icon={item.icon} width="14" />
            {item.label}
          </Link>
        );
      })}
    </div>
  );
}
