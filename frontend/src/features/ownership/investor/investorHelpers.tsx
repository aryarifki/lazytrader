import React from "react";

import { Badge } from "@/components/research/ui";
import type { InvestorHolding } from "./investorTypes";

export function fmtLot(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return "—";
  const v = Math.round(value / 100 / 1000);
  return `${v.toLocaleString("id-ID")}k`;
}

export function fmtLotT5(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return "—";
  return `${Math.round(value / 1000).toLocaleString("id-ID")}k`;
}

export function fmtRpM(valueM: number | null | undefined, withSign = false): string {
  if (valueM === null || valueM === undefined || Number.isNaN(valueM)) return "—";
  const sign = valueM < 0 ? "-" : withSign ? "+" : "";
  const abs = Math.abs(valueM);
  let value = abs;
  let unit = "M";
  if (abs >= 1000) {
    value = abs / 1000;
    unit = "T";
  }
  const decimals = value >= 100 ? 0 : value >= 10 ? 1 : 2;
  return `${sign}Rp ${value.toFixed(decimals)} ${unit}`;
}

export function fmtValueM(value: number | null | undefined): string {
  if (value === null || value === undefined || Number.isNaN(value)) return "—";
  return value.toLocaleString("id-ID", { maximumFractionDigits: 2 });
}

export function fmtDelta(row: InvestorHolding): React.ReactNode {
  if (row.status === "new") return <Badge tone="up">NEW</Badge>;
  if (row.status === "out") return <Badge tone="down">OUT</Badge>;
  if (row.sl_mei === null || row.sl_jun === null) return "—";
  const diff = row.sl_jun - row.sl_mei;
  if (diff === 0) return <span className="text-[var(--md-sys-color-on-surface-variant)]">—</span>;
  const value = Math.round(diff / 100 / 1000);
  return (
    <span
      className={
        diff > 0
          ? "text-emerald-600 dark:text-emerald-400"
          : "text-rose-600 dark:text-rose-400"
      }
    >
      {diff > 0 ? "+" : ""}
      {value.toLocaleString("id-ID")}k
    </span>
  );
}
