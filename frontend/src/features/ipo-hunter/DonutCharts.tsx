import React, { useState } from "react";
import type { IpoAfterIpoDistCard, IpoBuyerSellerPieCard, IpoDonutSlice } from "./ipoTypes";
import { fmtInt } from "./ipoFormat";
import { DonutRing } from "./DonutRing";

interface TooltipState {
  x: number;
  y: number;
  title: string;
  line1: string;
  line2?: string;
}

function DonutTooltip({ tip }: { tip: TooltipState | null }) {
  if (!tip) return null;
  return (
    <div
      className="fixed z-50 pointer-events-none px-3 py-2 rounded-lg bg-[var(--md-sys-color-surface-container-highest)] text-[var(--md-sys-color-on-surface)] text-[11px] shadow-xl border border-[var(--md-sys-color-outline-variant)]"
      style={{ left: tip.x, top: tip.y }}
    >
      <div className="font-bold">{tip.title}</div>
      <div className="font-mono-nums">{tip.line1}</div>
      {tip.line2 && <div className="text-[var(--md-sys-color-on-surface-variant)] text-[10px] mt-0.5">{tip.line2}</div>}
    </div>
  );
}

export function BuyerSellerDonutCard({ data }: { data?: IpoBuyerSellerPieCard | null }) {
  const [pieMode, setPieMode] = useState<"net" | "gross">("net");
  const [tip, setTip] = useState<TooltipState | null>(null);

  if (!data) return null;
  const active = pieMode === "net" ? data.net : data.gross;
  const buyer = active?.buyer || { total: 0, slices: [] };
  const seller = active?.seller || { total: 0, slices: [] };

  const legendMap = new Map<string, { label: string; color: string }>();
  [...buyer.slices, ...seller.slices].forEach((s: IpoDonutSlice) => {
    if (!legendMap.has(s.key)) legendMap.set(s.key, { label: s.label, color: s.color });
  });

  const handleHover = (e: React.MouseEvent, s: IpoDonutSlice, side: string) => {
    setTip({
      x: e.clientX + 12,
      y: e.clientY + 12,
      title: `${side} · ${s.label}`,
      line1: `${(s.pct || 0).toFixed(1).replace(".", ",")}% · ${fmtInt(s.value)} lot`,
      line2: s.top_brokers || undefined,
    });
  };

  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-4 shadow-sm relative">
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-xs font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface)]">
          {pieMode === "net" ? "NET BUYER vs NET SELLER" : "BUYER vs SELLER"}
        </h3>
        <div className="inline-flex rounded-lg bg-[var(--md-sys-color-surface-container-high)] p-0.5 border border-[var(--md-sys-color-outline-variant)]/50">
          <button
            type="button"
            onClick={() => setPieMode("gross")}
            className={`px-2.5 py-1 text-[10px] font-bold rounded-md transition cursor-pointer ${
              pieMode === "gross"
                ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
            }`}
          >
            Transaction
          </button>
          <button
            type="button"
            onClick={() => setPieMode("net")}
            className={`px-2.5 py-1 text-[10px] font-bold rounded-md transition cursor-pointer ${
              pieMode === "net"
                ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
            }`}
          >
            Net
          </button>
        </div>
      </div>

      <svg viewBox="-40 22 620 212" className="w-full h-auto">
        <DonutRing
          slices={buyer.slices}
          total={buyer.total}
          cx={140}
          cy={130}
          R={82}
          r={40}
          sideLabel={pieMode === "net" ? "Net Buyer" : "Buyer"}
          centerTitle={pieMode === "net" ? "NET BUYER" : "BUYER"}
          gapRight={true}
          onHover={handleHover}
          onLeave={() => setTip(null)}
        />
        <DonutRing
          slices={seller.slices}
          total={seller.total}
          cx={400}
          cy={130}
          R={82}
          r={40}
          sideLabel={pieMode === "net" ? "Net Seller" : "Seller"}
          centerTitle={pieMode === "net" ? "NET SELLER" : "SELLER"}
          gapRight={false}
          onHover={handleHover}
          onLeave={() => setTip(null)}
        />
        {/* 3-chevron Seller -> Buyer arrow */}
        <g>
          {(
            [
              [292, 0.3],
              [270, 0.65],
              [248, 1.0],
            ] as const
          ).map(([x, op]) => (
            <path
              key={x}
              d={`M${x + 9} 116 L${x - 3} 130 L${x + 9} 144`}
              fill="none"
              stroke="var(--md-sys-color-primary)"
              strokeWidth={4}
              strokeLinecap="round"
              strokeLinejoin="round"
              opacity={op}
            />
          ))}
        </g>
      </svg>

      <div className="flex flex-wrap items-center justify-center gap-x-4 gap-y-1.5 mt-2 pt-2 border-t border-[var(--md-sys-color-outline-variant)]/40">
        {Array.from(legendMap.entries()).map(([k, v]) => (
          <span key={k} className="inline-flex items-center gap-1.5 text-[11px] font-medium text-[var(--md-sys-color-on-surface-variant)]">
            <span className="w-2.5 h-2.5 rounded-sm inline-block" style={{ backgroundColor: v.color }} />
            {v.label}
          </span>
        ))}
      </div>

      <DonutTooltip tip={tip} />
    </div>
  );
}

export function AfterIpoDistCard({ data }: { data?: IpoAfterIpoDistCard | null }) {
  const [distMode, setDistMode] = useState<"sell" | "buy">("sell");
  const [tip, setTip] = useState<TooltipState | null>(null);

  if (!data || !data.enabled) return null;
  const side = distMode === "buy" ? data.buy : data.sell;

  const handleHover = (e: React.MouseEvent, s: IpoDonutSlice) => {
    const isIpo = s.key === "ipo";
    const nm = distMode === "buy" ? "Net Buyer" : "Net Seller";
    setTip({
      x: e.clientX + 12,
      y: e.clientY + 12,
      title: isIpo ? "Jumlah IPO" : `${nm} · ${s.label}`,
      line1: `${fmtInt(s.value)} lot${isIpo ? "" : ` · ${(s.pct_ipo || 0).toFixed(1).replace(".", ",")}% IPO`}`,
      line2: s.sub_brokers || undefined,
    });
  };

  return (
    <div className="rounded-2xl bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] p-4 shadow-sm relative">
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-xs font-bold uppercase tracking-wider text-[var(--md-sys-color-on-surface)]">
          AFTER IPO DISTRIBUTION
        </h3>
        <div className="inline-flex rounded-lg bg-[var(--md-sys-color-surface-container-high)] p-0.5 border border-[var(--md-sys-color-outline-variant)]/50">
          <button
            type="button"
            onClick={() => setDistMode("sell")}
            className={`px-2.5 py-1 text-[10px] font-bold rounded-md transition cursor-pointer ${
              distMode === "sell"
                ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
            }`}
          >
            Seller
          </button>
          <button
            type="button"
            onClick={() => setDistMode("buy")}
            className={`px-2.5 py-1 text-[10px] font-bold rounded-md transition cursor-pointer ${
              distMode === "buy"
                ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
            }`}
          >
            Buyer
          </button>
        </div>
      </div>

      {!side?.available ? (
        <div className="py-6 text-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
          Data jumlah IPO (lot) tidak tersedia.
        </div>
      ) : side.empty ? (
        <div className="py-6 text-center text-xs text-[var(--md-sys-color-on-surface-variant)]">
          Tidak ada kategori Net {distMode === "buy" ? "Buyer" : "Seller"} pada periode ini.
        </div>
      ) : (
        <>
          <svg viewBox="0 22 460 212" className="w-full h-auto max-w-md mx-auto">
            <DonutRing
              slices={side.slices}
              total={side.total}
              cx={230}
              cy={130}
              R={82}
              r={40}
              sideLabel="IPO"
              centerTitle="After IPO"
              centerSub={false}
              isDistribution={true}
              onHover={handleHover}
              onLeave={() => setTip(null)}
            />
          </svg>
          <div className="flex flex-wrap items-center justify-center gap-x-4 gap-y-1.5 mt-1 pt-2 border-t border-[var(--md-sys-color-outline-variant)]/40">
            {side.slices.map((s: IpoDonutSlice) => (
              <span
                key={s.key}
                className="inline-flex items-center gap-1.5 text-[11px] font-medium text-[var(--md-sys-color-on-surface-variant)]"
              >
                <span
                  className="w-2.5 h-2.5 rounded-sm inline-block"
                  style={{ backgroundColor: s.key === "ipo" ? "#94A3B8" : s.color }}
                />
                {s.label}
              </span>
            ))}
          </div>
        </>
      )}

      <DonutTooltip tip={tip} />
    </div>
  );
}
