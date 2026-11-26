import { useMemo } from "react";

import type { MonthlyComparisonRow, MonthlyPoint } from "./stockTypes";

export function useStockMonthlyRows(history: MonthlyPoint[]): MonthlyComparisonRow[] {
  return useMemo(() => {
    return history.map((curr, idx) => {
      const prev = idx > 0 ? history[idx - 1] : null;

      const currForeign = curr.foreign_total_pct ?? Number((curr.foreign_named_pct + curr.foreign_hidden_pct).toFixed(4));
      const prevForeign = prev ? (prev.foreign_total_pct ?? Number((prev.foreign_named_pct + prev.foreign_hidden_pct).toFixed(4))) : null;
      const foreignDelta = prevForeign != null ? Number((currForeign - prevForeign).toFixed(4)) : null;

      const currForeignHidden = curr.foreign_hidden_pct;
      const prevForeignHidden = prev ? prev.foreign_hidden_pct : null;
      const foreignHiddenDelta = prevForeignHidden != null ? Number((currForeignHidden - prevForeignHidden).toFixed(4)) : null;

      const currRetail = curr.retail_pct;
      const prevRetail = prev ? prev.retail_pct : null;
      const retailDelta = prevRetail != null ? Number((currRetail - prevRetail).toFixed(4)) : null;

      const currLocalNonRetail = curr.local_non_retail_pct;
      const prevLocalNonRetail = prev ? prev.local_non_retail_pct : null;
      const localNonRetailDelta = prevLocalNonRetail != null ? Number((currLocalNonRetail - prevLocalNonRetail).toFixed(4)) : null;

      const currSid = curr.total_sid;
      const prevSid = prev ? prev.total_sid : null;
      const sidDelta = curr.sid_change_mom ?? (prevSid != null && currSid != null ? currSid - prevSid : null);
      const sidDeltaPct = curr.sid_change_pct ?? (prevSid && prevSid > 0 && sidDelta != null ? Number(((sidDelta / prevSid) * 100).toFixed(2)) : null);

      let flowSignal = "Baseline Awal";
      let flowTone: "up" | "down" | "primary" | "warning" | "neutral" = "neutral";
      if (idx > 0) {
        const fUp = (foreignDelta ?? 0) > 0.01;
        const fDown = (foreignDelta ?? 0) < -0.01;
        const rUp = (retailDelta ?? 0) > 0.01;
        const rDown = (retailDelta ?? 0) < -0.01;

        if (fUp && rDown) {
          flowSignal = "Akumulasi Asing";
          flowTone = "up";
        } else if (fDown && rUp) {
          flowSignal = "Distribusi Asing ke Ritel";
          flowTone = "down";
        } else if (fUp && rUp) {
          flowSignal = "Inflow Bersama";
          flowTone = "primary";
        } else if (rDown && (localNonRetailDelta ?? 0) > 0.01) {
          flowSignal = "Akumulasi Institusi Lokal";
          flowTone = "primary";
        } else if (fDown && rDown) {
          flowSignal = "Outflow Asing & Ritel";
          flowTone = "warning";
        } else {
          flowSignal = "Konsolidasi / Netral";
          flowTone = "neutral";
        }
      }

      return {
        period_label: curr.period_label,
        isLatest: idx === history.length - 1,
        currForeign,
        currForeignNamed: curr.foreign_named_pct,
        currForeignHidden,
        foreignDelta,
        foreignHiddenDelta,
        currRetail,
        retailDelta,
        currLocalNonRetail,
        localNonRetailDelta,
        currSid,
        sidDelta,
        sidDeltaPct,
        flowSignal,
        flowTone,
      };
    }).reverse();
  }, [history]);
}
