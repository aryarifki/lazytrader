import React from "react";
import { Link } from "react-router-dom";
import { Icon } from "@iconify/react";

interface StockBadgesProps {
  badges?: string[] | null;
  badge?: string | null;
  conglomerate?: string | null;
  konglo?: string | null;
  size?: "sm" | "md";
  showConglomerate?: boolean;
  showBadges?: boolean;
  className?: string;
}

export function StockBadges({
  badges,
  badge,
  conglomerate,
  konglo,
  size = "md",
  showConglomerate = true,
  showBadges = true,
  className = "",
}: StockBadgesProps) {
  // Collect all unique badges
  const badgeList: string[] = [];
  if (Array.isArray(badges)) {
    for (const b of badges) {
      if (b && typeof b === "string" && !badgeList.includes(b.trim())) {
        badgeList.push(b.trim());
      }
    }
  }
  if (badge && typeof badge === "string" && !badgeList.includes(badge.trim())) {
    badgeList.push(badge.trim());
  }

  const groupName = (conglomerate || konglo || "").trim();
  const hasGroup = groupName && groupName !== "-" && groupName.toLowerCase() !== "none";

  if (!badgeList.length && (!showConglomerate || !hasGroup)) {
    return null;
  }

  const isSm = size === "sm";
  const badgeClass = isSm
    ? "px-1 py-0.5 text-[9px] font-bold rounded tracking-wider uppercase border leading-none shrink-0"
    : "px-1.5 py-0.5 text-[10px] font-bold rounded tracking-wider uppercase border leading-none shrink-0";

  const groupClass = isSm
    ? "inline-flex items-center gap-1 px-1.5 py-0.5 text-[9px] font-medium rounded-full border border-emerald-500/25 bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20 transition-colors leading-none shrink-0"
    : "inline-flex items-center gap-1 px-2 py-0.5 text-[10px] font-medium rounded-full border border-emerald-500/25 bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20 transition-colors leading-none shrink-0";

  return (
    <div className={`inline-flex items-center gap-1 flex-wrap ${className}`}>
      {showBadges &&
        badgeList.map((item) => {
          const upper = item.toUpperCase();
          let style = "bg-slate-500/15 text-slate-300 border-slate-500/30";
          if (upper.includes("MSCI")) {
            style = "bg-blue-500/15 text-blue-400 border-blue-500/30";
          } else if (upper.includes("HSC")) {
            style = "bg-purple-500/15 text-purple-400 border-purple-500/30";
          } else if (upper.includes("FCA")) {
            style = "bg-amber-500/15 text-amber-400 border-amber-500/30";
          }

          return (
            <span key={item} className={`${badgeClass} ${style}`}>
              {item}
            </span>
          );
        })}

      {showConglomerate && hasGroup && (
        <Link
          to="/group"
          title={`Konglomerasi: ${groupName}`}
          className={groupClass}
          onClick={(e) => e.stopPropagation()}
        >
          <Icon icon="ph:buildings-duotone" className={isSm ? "w-2.5 h-2.5" : "w-3 h-3"} />
          <span className="truncate max-w-[130px]">{groupName}</span>
        </Link>
      )}
    </div>
  );
}

export default StockBadges;
