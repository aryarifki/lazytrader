import { csvCell, fmtExplorerPct } from "./ctHelpers";
import type { Category, ExplorerRow, Mode } from "./ctTypes";

export function exportCtCsv(mode: Mode, category: Category, rows: ExplorerRow[]): void {
  const header =
    mode === "stocks"
      ? [
          "Saham",
          "Perubahan (%)",
          "Delta SID",
          "SID",
          "Net Value Transaksi (M)",
          "Frekuensi Transaksi",
          "Perubahan (%)",
          "Terkini (%)",
          "Harga MTD (%)",
          "Market Cap (M)",
        ]
      : [
          "Investor",
          "Saham yang Cocok",
          "Net Value Transaksi (M)",
          "Frekuensi Transaksi",
          "Jumlah Saham",
          "Nilai Transaksi (M)",
          "G/L (Rp M)",
          "G/L (%)",
          "Tipe",
        ];
  const body = rows.map((row) =>
    mode === "stocks"
      ? [
          row.ticker,
          row.net_pct,
          row.sid_change,
          row.sid,
          row.net_value_m,
          row.count,
          row.pct_change,
          row.latest,
          row.mtd_pct,
          row.market_cap_b,
        ]
      : [
          row.name,
          (row.matches ?? [])
            .map((match) => `${match.ticker} ${fmtExplorerPct(match.delta)}`)
            .join("; "),
          row.net_value_m,
          row.count,
          row.ticker_count,
          row.abs_value_m,
          row.gl_rp_m,
          row.pct,
          row.lf_label ?? row.lf,
        ]
  );
  const csv = [header, ...body]
    .map((line) => line.map(csvCell).join(","))
    .join("\r\n");
  const blob = new Blob([`\uFEFF${csv}`], { type: "text/csv;charset=utf-8;" });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = `investowl-ct-explorer-${mode}-${category}-${new Date()
    .toISOString()
    .slice(0, 10)}.csv`;
  document.body.appendChild(anchor);
  anchor.click();
  anchor.remove();
  URL.revokeObjectURL(url);
}
