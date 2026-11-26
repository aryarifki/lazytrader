import { StatCard } from "@/components/research/ui";
import {
  fmtDate,
  fmtNumber,
  fmtPct,
  fmtRupiah,
  fmtSignedPct,
  toneFor,
} from "@/components/research/format";
import type { StockOwnership } from "./stockTypes";

export function StockStatCards({
  data,
  isCtBaseline,
  freeFloat,
}: {
  data: StockOwnership;
  isCtBaseline: boolean;
  freeFloat: number | null | undefined;
}) {
  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-2.5">
      <StatCard
        label="Harga terakhir"
        value={fmtRupiah(data.price?.close)}
        sub={data.price?.date ? `Per ${fmtDate(data.price.date)}` : "Harga belum tersedia"}
        icon="ph:currency-circle-dollar-duotone"
      />
      <StatCard
        label="Perubahan 1D"
        value={fmtSignedPct(data.price?.change_1d_pct)}
        sub="Dibanding penutupan sebelumnya"
        icon="ph:trend-up-duotone"
        tone={toneFor(data.price?.change_1d_pct)}
      />
      <StatCard
        label="Total SID"
        value={data.profile.total_sid ? fmtNumber(data.profile.total_sid) : "—"}
        sub={
          !data.profile.total_sid
            ? "Belum dilaporkan oleh emiten"
            : isCtBaseline
              ? `Baseline CT (statis) · ${data.as_of.ct_baseline ?? "Sep 2026"}`
              : data.data_quality.monthly_composition
                ? "Snapshot KSEI bulanan"
                : "Belum tersedia untuk emiten ini"
        }
        icon="ph:users-three-duotone"
      />
      <StatCard
        label="Free float est."
        value={fmtPct(freeFloat)}
        sub={
          data.profile.ct_free_float_pct != null
            ? `CT ${fmtPct(data.profile.ct_free_float_pct)} · estimasi pemegang saham`
            : "Estimasi dari pemegang saham pengendali & treasury"
        }
        icon="ph:arrows-out-line-horizontal-duotone"
        tone="primary"
      />
    </div>
  );
}
