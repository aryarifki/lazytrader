import { csvCell } from "./snapshotHelpers";
import type { InvestorRow, Mode, StockRow } from "./snapshotTypes";

export function exportSnapshotCsv(
  mode: Mode,
  stockRows: StockRow[],
  investorRows: InvestorRow[]
): void {
  const header =
    mode === "stocks"
      ? [
          "Ticker",
          "Nama",
          "Sektor",
          "Harga",
          "Tanggal harga",
          "MTD %",
          "Holders",
          "Pengendali %",
          "Publik %",
          "Top holder %",
          "Treasury %",
          "Substantial %",
          "Market cap (miliar Rp)",
          "Foreign 5D",
        ]
      : [
          "Investor",
          "Jumlah emiten",
          "Total porsi %",
          "Porsi max %",
          "Pengendali",
          "Substantial",
          "Sektor",
          "Contoh emiten",
        ];
  const body =
    mode === "stocks"
      ? stockRows.map((row) => [
          row.ticker,
          row.name,
          row.sector,
          row.price,
          row.price_date,
          row.mtd_pct,
          row.holder_count,
          row.controlling_pct,
          row.public_pct,
          row.max_holder_pct,
          row.treasury_pct,
          row.substantial_pct,
          row.market_cap_b,
          row.foreign_5d,
        ])
      : investorRows.map((row) => [
          row.name,
          row.ticker_count,
          row.total_pct,
          row.max_pct,
          row.controlling_count,
          row.substantial_count,
          row.sector_count,
          row.sample_tickers.join(" "),
        ]);
  const csv = [header, ...body]
    .map((line) => line.map(csvCell).join(","))
    .join("\r\n");
  const blob = new Blob([`\uFEFF${csv}`], { type: "text/csv;charset=utf-8;" });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = `investowl-explorer-${mode}-${new Date().toISOString().slice(0, 10)}.csv`;
  document.body.appendChild(anchor);
  anchor.click();
  anchor.remove();
  URL.revokeObjectURL(url);
}
