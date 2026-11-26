import { Link } from "react-router-dom";
import { fmtNumber, toneFor } from "@/components/research/format";
import {
  fmtExplorerPct,
  fmtExplorerValM,
  fmtPortfolioPct,
  tickerValueItems,
  toneTextClass,
} from "./ctHelpers";
import type { Category, Column } from "./ctTypes";

export function buildInvestorColumns(
  numberColumn: Column,
  category: Category,
  condition: string
): Column[] {
  const nameColumn: Column = {
    key: "name",
    label: "Investor",
    sort: "name",
    render: (row) =>
      row.name ? (
        <Link
          to={`/ownership?tab=investor&name=${encodeURIComponent(row.name)}`}
          className="font-semibold text-[var(--md-sys-color-primary)] hover:underline"
        >
          {row.name}
        </Link>
      ) : (
        "—"
      ),
  };

  if (category === "kepemilikan") {
    return [
      numberColumn,
      nameColumn,
      {
        key: "matches",
        label: "Saham yang Cocok",
        sort: "matches",
        render: (row) => {
          const matches = row.matches ?? [];
          if (!matches.length) return "—";
          return (
            <div className="flex flex-wrap items-center gap-1">
              {matches.slice(0, 2).map((match) => (
                <span
                  key={match.ticker}
                  className={`inline-flex items-center gap-1 rounded-md border px-1.5 py-0.5 font-mono text-[9px] ${
                    match.delta >= 0
                      ? "border-emerald-500/30 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400"
                      : "border-rose-500/30 bg-rose-500/10 text-rose-700 dark:text-rose-400"
                  }`}
                >
                  <Link
                    to={`/ownership?tab=saham&ticker=${encodeURIComponent(match.ticker)}`}
                    className="hover:underline"
                  >
                    {match.ticker}
                  </Link>
                  <span>{fmtExplorerPct(match.delta)}</span>
                </span>
              ))}
              {matches.length > 2 ? (
                <span className="rounded-md border border-[var(--md-sys-color-outline-variant)] px-1.5 py-0.5 text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                  +{matches.length - 2} lagi
                </span>
              ) : null}
            </div>
          );
        },
      },
    ];
  }

  if (category === "transaksi") {
    return [
      numberColumn,
      nameColumn,
      {
        key: "net",
        label: "Net Value Transaksi",
        align: "right",
        sort: "net",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.net_value_m))}>
            {fmtExplorerValM(row.net_value_m)}
          </span>
        ),
      },
      {
        key: "count",
        label: "Frekuensi Transaksi",
        align: "right",
        sort: "count",
        render: (row) => fmtNumber(row.count, 0),
      },
      {
        key: "ticker_count",
        label: "Jumlah Saham",
        align: "right",
        sort: "ticker_count",
        render: (row) => fmtNumber(row.ticker_count, 0),
      },
    ];
  }

  if (condition === "investorAktif") {
    return [
      numberColumn,
      nameColumn,
      {
        key: "count",
        label: "Frekuensi Transaksi",
        align: "right",
        sort: "count",
        render: (row) => fmtNumber(row.count, 0),
      },
      {
        key: "net",
        label: "Net Transaksi",
        align: "right",
        sort: "net",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.net_value_m))}>
            {fmtExplorerValM(row.net_value_m)}
          </span>
        ),
      },
      {
        key: "abs_value",
        label: "Nilai Transaksi",
        align: "right",
        sort: "abs_value",
        render: (row) => fmtExplorerValM(row.abs_value_m),
      },
    ];
  }

  if (condition === "investorNambah") {
    return [
      numberColumn,
      nameColumn,
      {
        key: "net",
        label: "Net Transaksi",
        align: "right",
        sort: "net",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.net_value_m))}>
            {fmtExplorerValM(row.net_value_m)}
          </span>
        ),
      },
      {
        key: "tickers",
        label: "Saham",
        render: (row) => {
          const tickers = tickerValueItems(row.tickers);
          if (!tickers.length) return "—";
          return (
            <div className="flex flex-wrap items-center gap-1">
              {tickers.slice(0, 2).map((item) => (
                <span
                  key={item.ticker}
                  className={`inline-flex items-center gap-1 rounded-md border px-1.5 py-0.5 font-mono text-[9px] ${
                    item.net_value_m >= 0
                      ? "border-emerald-500/30 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400"
                      : "border-rose-500/30 bg-rose-500/10 text-rose-700 dark:text-rose-400"
                  }`}
                >
                  <Link
                    to={`/ownership?tab=saham&ticker=${encodeURIComponent(item.ticker)}`}
                    className="hover:underline"
                  >
                    {item.ticker}
                  </Link>
                  <span>{fmtExplorerValM(item.net_value_m)}</span>
                </span>
              ))}
              {tickers.length > 2 ? (
                <span className="rounded-md border border-[var(--md-sys-color-outline-variant)] px-1.5 py-0.5 text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                  +{tickers.length - 2} lagi
                </span>
              ) : null}
            </div>
          );
        },
      },
    ];
  }

  return [
    numberColumn,
    nameColumn,
    {
      key: "gl",
      label: "G/L (Rp)",
      align: "right",
      sort: "gl",
      render: (row) => (
        <span className={toneTextClass(toneFor(row.gl_rp_m))}>
          {fmtExplorerValM(row.gl_rp_m)}
        </span>
      ),
    },
    {
      key: "pct",
      label: "G/L (%)",
      align: "right",
      sort: "pct",
      render: (row) => (
        <span className={toneTextClass(toneFor(row.pct))}>
          {fmtPortfolioPct(row.pct)}
        </span>
      ),
    },
  ];
}
