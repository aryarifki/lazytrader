import { Badge } from "@/components/research/ui";
import { fmtNumber } from "@/components/research/format";

export const QUICK_TICKERS = ["BBCA", "BBRI", "BMRI", "GOTO", "ASII", "BRPT"];

export function fmtLot(value: number | null | undefined): string {
  if (value == null || Number.isNaN(value)) return "—";
  const abs = Math.abs(value);
  const sign = value < 0 ? "-" : "";
  if (abs >= 1_000_000) return `${sign}${(abs / 1_000_000).toFixed(1)} jt`;
  if (abs >= 1_000) return `${sign}${(abs / 1_000).toFixed(0)}k`;
  return `${sign}${fmtNumber(abs)}`;
}

export function HolderStatus({ status }: { status: string | null }) {
  const normalized = (status ?? "").toLowerCase();
  if (normalized === "new" || normalized === "baru") return <Badge tone="up">Baru</Badge>;
  if (normalized === "naik") return <Badge tone="up">Naik</Badge>;
  if (normalized === "out" || normalized === "keluar") return <Badge tone="down">Keluar</Badge>;
  if (normalized === "turun") return <Badge tone="down">Turun</Badge>;
  if (normalized === "both" || normalized === "tetap") return <Badge tone="neutral">Tetap</Badge>;
  return <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)]">—</span>;
}

export function actionTone(action: string): "up" | "down" | "warning" | "primary" | "neutral" {
  const a = action.toUpperCase();
  if (a === "BUY" || a === "BUYBACK") return "up";
  if (a === "SELL") return "down";
  if (a === "TRANSFER" || a === "TENDER_OFFER") return "warning";
  if (a === "RIGHT_ISSUE" || a === "PRIVATE_PLACEMENT" || a === "MESOP_OPTION" || a === "WARRANT_EXERCISE") return "primary";
  return "neutral";
}
