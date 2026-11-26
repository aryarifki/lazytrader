import React from "react";
import { buildInvestorColumns } from "./ctColumnsInvestors";
import { buildStockColumns } from "./ctColumnsStocks";
import type { Category, Column, ExplorerPeriod, Mode } from "./ctTypes";

export function buildColumns(
  mode: Mode,
  category: Category,
  condition: string,
  period: ExplorerPeriod | undefined
): Column[] {
  const numberColumn: Column = {
    key: "no",
    label: "#",
    align: "right",
    render: (_row, index) =>
      React.createElement(
        "span",
        {
          className:
            "text-[10px] font-mono text-[var(--md-sys-color-on-surface-variant)]",
        },
        index + 1
      ),
  };

  if (mode === "stocks") {
    return buildStockColumns(numberColumn, category, condition, period);
  }
  return buildInvestorColumns(numberColumn, category, condition);
}
