import { Link } from "react-router-dom";
import { fmtNumber, toneFor } from "@/components/research/format";
import { StockBadges } from "@/components/research/StockBadges";
import {
  fmtExplorerPct,
  fmtExplorerValM,
  fmtMtd,
  fmtPlainPct,
  fmtSidDelta,
  mtdMonthLabel,
  sidMonthLabel,
  toneTextClass,
} from "./ctHelpers";
import type { Category, Column, ExplorerPeriod } from "./ctTypes";

export function buildStockColumns(
  numberColumn: Column,
  category: Category,
  condition: string,
  period: ExplorerPeriod | undefined
): Column[] {
  const tickerColumn: Column = {
    key: "ticker",
    label: "Saham",
    sort: "ticker",
    render: (row) =>
      row.ticker ? (
        <div className="flex items-center gap-1.5 flex-wrap">
          <Link
            to={`/ownership?tab=saham&ticker=${encodeURIComponent(row.ticker)}`}
            className="font-mono font-bold text-[var(--md-sys-color-primary)] hover:underline"
          >
            {row.ticker}
          </Link>
          <StockBadges
            badges={row.badges}
            badge={row.badge}
            conglomerate={row.conglomerate || row.konglo}
            size="sm"
          />
        </div>
      ) : (
        "—"
      ),
  };
  const mtdColumn: Column = {
    key: "mtd",
    label: `Δ Harga (${mtdMonthLabel(period)})`,
    align: "right",
    sort: "mtd",
    render: (row) => (
      <span className={toneTextClass(toneFor(row.mtd_pct))}>
        {fmtMtd(row.mtd_pct)}
      </span>
    ),
  };

  if (category === "kepemilikan") {
    if (condition === "sidNaik" || condition === "sidTurun") {
      const month = sidMonthLabel(period);
      return [
        numberColumn,
        tickerColumn,
        {
          key: "sid_delta",
          label: `Δ SID (${month})`,
          align: "right",
          sort: "sid_amount",
          render: (row) => (
            <span className={toneTextClass(toneFor(row.net_pct))}>
              {fmtSidDelta(row.sid_change, row.net_pct)}
            </span>
          ),
        },
        {
          key: "sid",
          label: `SID (${month})`,
          align: "right",
          sort: "sid_value",
          render: (row) => fmtNumber(row.sid, 0),
        },
        mtdColumn,
      ];
    }
    return [
      numberColumn,
      tickerColumn,
      {
        key: "net",
        label: "Δ Kepemilikan >1%",
        align: "right",
        sort: "net",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.net_pct))}>
            {fmtExplorerPct(row.net_pct)}
          </span>
        ),
      },
      mtdColumn,
    ];
  }

  if (category === "transaksi") {
    return [
      numberColumn,
      tickerColumn,
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
      mtdColumn,
    ];
  }

  const advancedColumns: Column[] = [numberColumn, tickerColumn];
  if (condition === "aktifBeli") {
    advancedColumns.push(
      {
        key: "count",
        label: "Frekuensi Transaksi",
        align: "right",
        sort: "count",
        render: (row) => fmtNumber(row.count, 0),
      },
      {
        key: "net",
        label: "Total Perubahan (Rp)",
        align: "right",
        sort: "net",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.net_value_m))}>
            {fmtExplorerValM(row.net_value_m)}
          </span>
        ),
      }
    );
  } else if (condition === "buyback" || condition === "mayoritas") {
    advancedColumns.push(
      {
        key: "pct",
        label: "Perubahan Kepemilikan",
        align: "right",
        sort: "pct",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.pct_change))}>
            {fmtExplorerPct(row.pct_change)}
          </span>
        ),
      },
      {
        key: "net",
        label: "Perubahan (Rp)",
        align: "right",
        sort: "net",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.net_value_m))}>
            {fmtExplorerValM(row.net_value_m)}
          </span>
        ),
      }
    );
  } else if (condition === "baru") {
    advancedColumns.push({
      key: "pct",
      label: "Perubahan Kepemilikan",
      align: "right",
      sort: "pct",
      render: (row) => (
        <span className={toneTextClass(toneFor(row.pct_change))}>
          {fmtExplorerPct(row.pct_change)}
        </span>
      ),
    });
  } else if (condition === "masukFlat") {
    advancedColumns.push(
      {
        key: "pct",
        label: "Perubahan Kepemilikan",
        align: "right",
        sort: "pct",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.pct_change))}>
            {fmtExplorerPct(row.pct_change)}
          </span>
        ),
      },
      {
        key: "net",
        label: "Transaksi Bulan Ini",
        align: "right",
        sort: "net",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.net_value_m))}>
            {fmtExplorerValM(row.net_value_m)}
          </span>
        ),
      }
    );
  } else if (condition === "sid3Turun" || condition === "sid3Naik") {
    advancedColumns.push(
      {
        key: "sid_diff",
        label: "Perubahan SID (3 Bulan)",
        align: "right",
        sort: "sid_change",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.pct_change))}>
            {fmtSidDelta(row.sid_change, row.pct_change)}
          </span>
        ),
      },
      {
        key: "latest",
        label: "SID Terkini",
        align: "right",
        sort: "latest",
        render: (row) => fmtNumber(row.latest, 0),
      }
    );
  } else if (condition === "retailRendah") {
    advancedColumns.push({
      key: "latest",
      label: "% Retail",
      align: "right",
      sort: "latest",
      render: (row) => fmtPlainPct(row.latest),
    });
  } else if (
    condition === "retailTurun" ||
    condition === "hiddenNaik" ||
    condition === "foreignNaik"
  ) {
    advancedColumns.push(
      {
        key: "pct",
        label: "Perubahan",
        align: "right",
        sort: "pct",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.pct_change))}>
            {fmtExplorerPct(row.pct_change)}
          </span>
        ),
      },
      {
        key: "latest",
        label: "% Terkini",
        align: "right",
        sort: "latest",
        render: (row) => fmtPlainPct(row.latest),
      }
    );
  } else if (condition === "foreign3Bulan" || condition === "retail3Bulan") {
    advancedColumns.push(
      {
        key: "pct",
        label: "Perubahan (3 Bulan)",
        align: "right",
        sort: "pct",
        render: (row) => (
          <span className={toneTextClass(toneFor(row.pct_change))}>
            {fmtExplorerPct(row.pct_change)}
          </span>
        ),
      },
      {
        key: "latest",
        label: "% Terkini",
        align: "right",
        sort: "latest",
        render: (row) => fmtPlainPct(row.latest),
      }
    );
  } else {
    advancedColumns.push({
      key: "pct",
      label: "Perubahan Kepemilikan",
      align: "right",
      sort: "pct",
      render: (row) => (
        <span className={toneTextClass(toneFor(row.pct_change))}>
          {fmtExplorerPct(row.pct_change)}
        </span>
      ),
    });
  }
  advancedColumns.push(mtdColumn);
  return advancedColumns;
}
