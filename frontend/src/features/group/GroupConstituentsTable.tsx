import { Link } from "react-router-dom";
import { Badge, Card, DataNote, ErrorBlock, LoadingBlock, SectionTitle } from "@/components/research/ui";
import { fmtNumber, fmtPct, fmtRupiah, fmtSignedPct } from "@/components/research/format";
import { StockBadges } from "@/components/research/StockBadges";
import { fmtMarketCapB, pctClass } from "./groupHelpers";
import type { GroupDetail, GroupsResponse } from "./groupTypes";

interface GroupConstituentsTableProps {
  detail: GroupDetail | undefined;
  detailLoading: boolean;
  detailError: Error | undefined;
  groupsData: GroupsResponse;
}

export function GroupConstituentsTable({
  detail,
  detailLoading,
  detailError,
  groupsData,
}: GroupConstituentsTableProps) {
  return (
    <Card className="overflow-hidden">
      <div className="p-4 sm:p-5 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <SectionTitle
          icon="ph:buildings-duotone"
          title="Emiten dalam grup"
          subtitle="Harga, kepemilikan, dan snapshot KSEI"
          action={<Badge>{detail?.constituents.length ?? 0} emiten</Badge>}
        />
      </div>
      {detailLoading ? (
        <div className="p-5">
          <LoadingBlock label="Memuat emiten grup..." />
        </div>
      ) : detailError ? (
        <div className="p-5">
          <ErrorBlock message={`Gagal memuat grup: ${detailError.message}`} />
        </div>
      ) : detail ? (
        <div className="overflow-x-auto">
          <table className="w-full min-w-[1420px] text-left text-[11px]">
            <thead className="bg-[var(--md-sys-color-surface-container-high)]/70 text-[9px] font-mono uppercase tracking-wider text-[var(--md-sys-color-on-surface-variant)]">
              <tr>
                <th className="px-4 py-2.5 font-bold">#</th>
                <th className="px-4 py-2.5 font-bold">Emiten</th>
                <th className="px-4 py-2.5 font-bold">Sektor</th>
                <th className="px-4 py-2.5 font-bold text-right">Harga</th>
                <th className="px-4 py-2.5 font-bold text-right">1D</th>
                <th className="px-4 py-2.5 font-bold text-right">MTD</th>
                <th className="px-4 py-2.5 font-bold text-right">Big Player</th>
                <th className="px-4 py-2.5 font-bold text-right">Foreign Hidden</th>
                <th className="px-4 py-2.5 font-bold text-right">Retail</th>
                <th className="px-4 py-2.5 font-bold text-right">SID</th>
                <th className="px-4 py-2.5 font-bold text-right">Market Cap</th>
                <th className="px-4 py-2.5 font-bold text-right">Holders</th>
                <th className="px-4 py-2.5 font-bold text-right">Pengendali</th>
                <th className="px-4 py-2.5 font-bold text-right">Publik</th>
                <th className="px-4 py-2.5 font-bold text-right">Free float</th>
                <th className="px-4 py-2.5 font-bold">KSEI</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--md-sys-color-outline-variant)]/40">
              {detail.constituents.map((item, index) => (
                <tr key={item.ticker} className="hover:bg-[var(--md-sys-color-surface-container-high)]/40">
                  <td className="px-4 py-2.5 font-mono text-[var(--md-sys-color-on-surface-variant)]">
                    {index + 1}
                  </td>
                  <td className="px-4 py-2.5">
                    <div className="flex items-center gap-1.5 flex-wrap">
                      <Link
                        to={`/ownership?tab=saham&ticker=${encodeURIComponent(item.ticker)}`}
                        className="font-mono font-bold text-[var(--md-sys-color-primary)] hover:underline"
                      >
                        {item.ticker}
                      </Link>
                      <StockBadges
                        badges={item.badges}
                        badge={item.badge}
                        showConglomerate={false}
                        size="sm"
                      />
                    </div>
                    <div className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] max-w-[220px] line-clamp-1">
                      {item.name}
                    </div>
                  </td>
                  <td className="px-4 py-2.5 text-[10px] text-[var(--md-sys-color-on-surface-variant)]">
                    {item.sector}
                  </td>
                  <td className="px-4 py-2.5 text-right font-mono">{fmtRupiah(item.price)}</td>
                  <td className={`px-4 py-2.5 text-right font-mono font-bold ${pctClass(item.change_1d_pct)}`}>
                    {fmtSignedPct(item.change_1d_pct)}
                  </td>
                  <td className={`px-4 py-2.5 text-right font-mono font-bold ${pctClass(item.mtd_pct)}`}>
                    {fmtSignedPct(item.mtd_pct)}
                  </td>
                  <td className="px-4 py-2.5 text-right font-mono whitespace-nowrap">{fmtPct(item.big_player_pct)}</td>
                  <td className="px-4 py-2.5 text-right font-mono whitespace-nowrap">{fmtPct(item.foreign_hidden_pct)}</td>
                  <td className="px-4 py-2.5 text-right font-mono whitespace-nowrap">{fmtPct(item.retail_pct)}</td>
                  <td className="px-4 py-2.5 text-right font-mono whitespace-nowrap">
                    {item.sid != null ? fmtNumber(item.sid) : "—"}
                  </td>
                  <td className="px-4 py-2.5 text-right font-mono whitespace-nowrap">
                    {fmtMarketCapB(item.market_cap_b)}
                  </td>
                  <td className="px-4 py-2.5 text-right font-mono">{fmtNumber(item.holder_count)}</td>
                  <td className="px-4 py-2.5 text-right font-mono">{fmtPct(item.controlling_pct)}</td>
                  <td className="px-4 py-2.5 text-right font-mono">{fmtPct(item.public_pct)}</td>
                  <td className="px-4 py-2.5 text-right font-mono">{fmtPct(item.free_float_est)}</td>
                  <td className="px-4 py-2.5 text-[9px] font-mono text-[var(--md-sys-color-on-surface-variant)]">
                    {item.monthly?.period_label ?? "—"}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : null}
      <div className="p-4 border-t border-[var(--md-sys-color-outline-variant)]/60">
        <DataNote tone={detail?.data_quality.ct_baseline ? "info" : "warning"}>
          <ul className="list-disc pl-4 space-y-1">
            <li>
              {detail?.data_quality.market_cap
                ? "Market cap dihitung dari baseline CT atau harga terakhir x listed shares IDX."
                : "Market cap per emiten belum tersedia, sehingga urutan memakai MTD."}
            </li>
            <li>
              {detail?.data_quality.t5_transactions
                ? "Transaksi T5 berasal dari baseline CT (quarantined), bukan data sinyal."
                : "Transaksi T5 belum tersedia untuk emiten di grup ini."}
            </li>
            <li>
              {groupsData.coverage.ct_baseline
                ? `Grup yang dipetakan: ${groupsData.coverage.mapped_groups} grup baseline CT.`
                : `Grup yang dipetakan: ${groupsData.coverage.mapped_groups} preset.`}
            </li>
          </ul>
        </DataNote>
      </div>
    </Card>
  );
}
