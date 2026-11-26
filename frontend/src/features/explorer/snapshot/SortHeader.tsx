import { Icon } from "@iconify/react";

import type { SortOrder } from "./snapshotTypes";

export function SortHeader({
  field,
  label,
  align = "left",
  current,
  order,
  onSort,
}: {
  field: string;
  label: string;
  align?: "left" | "right";
  current: string;
  order: SortOrder;
  onSort: (field: string) => void;
}) {
  const active = current === field;
  return (
    <th
      className={`px-4 py-2.5 font-bold ${align === "right" ? "text-right" : ""}`}
      aria-sort={active ? (order === "asc" ? "ascending" : "descending") : "none"}
    >
      <button
        type="button"
        onClick={() => onSort(field)}
        className={`inline-flex items-center gap-1 hover:text-[var(--md-sys-color-primary)] transition-colors ${
          align === "right" ? "justify-end w-full" : ""
        } ${active ? "text-[var(--md-sys-color-primary)]" : ""}`}
      >
        {label}
        <Icon
          icon={
            active
              ? order === "asc"
                ? "ph:sort-ascending-duotone"
                : "ph:sort-descending-duotone"
              : "ph:arrows-down-up-duotone"
          }
          width="12"
          className={active ? "opacity-100" : "opacity-45"}
        />
      </button>
    </th>
  );
}
