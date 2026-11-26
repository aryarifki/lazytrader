import React from "react";
import type { IpoDonutSlice } from "./ipoTypes";
import { fmtInt } from "./ipoFormat";

function polar(cx: number, cy: number, a: number, r: number): [number, number] {
  return [cx + r * Math.cos(a), cy + r * Math.sin(a)];
}

function buildRingPath(cx: number, cy: number, R: number, r: number, a1: number, a2: number): string {
  const [x1, y1] = polar(cx, cy, a1, R);
  const [x2, y2] = polar(cx, cy, a2, R);
  const [x3, y3] = polar(cx, cy, a2, r);
  const [x4, y4] = polar(cx, cy, a1, r);
  const large = a2 - a1 > Math.PI ? 1 : 0;
  return `M${x1.toFixed(2)} ${y1.toFixed(2)} A${R} ${R} 0 ${large} 1 ${x2.toFixed(2)} ${y2.toFixed(2)} L${x3.toFixed(2)} ${y3.toFixed(2)} A${r} ${r} 0 ${large} 0 ${x4.toFixed(2)} ${y4.toFixed(2)} Z`;
}

export interface DonutRingProps {
  slices: IpoDonutSlice[];
  total: number;
  cx: number;
  cy: number;
  R: number;
  r: number;
  sideLabel: string;
  centerTitle: string;
  centerSub?: string | false;
  gapRight?: boolean;
  isDistribution?: boolean;
  onHover: (e: React.MouseEvent, s: IpoDonutSlice, side: string) => void;
  onLeave: () => void;
}

export function DonutRing({
  slices,
  total,
  cx,
  cy,
  R,
  r,
  sideLabel,
  centerTitle,
  centerSub,
  gapRight,
  isDistribution,
  onHover,
  onLeave,
}: DonutRingProps) {
  let ang = -Math.PI / 2;
  const outer: Array<{ mid: number; label: string; txt: string }> = [];
  const elements: React.ReactNode[] = [];

  slices.forEach((s, idx) => {
    const frac = total > 0 ? s.value / total : 0;
    if (frac <= 0) return;
    const a2 = ang + frac * 2 * Math.PI;
    const fillCol = s.key === "ipo" ? "#94A3B8" : s.color;
    const labelTxt = isDistribution
      ? s.key === "ipo"
        ? `${fmtInt(s.value / 1000)} k`
        : `${(s.pct_ipo || 0).toFixed(1).replace(".", ",")}% IPO`
      : `${(frac * 100).toFixed(1)}%`;

    if (frac >= 0.9999) {
      elements.push(
        <circle
          key={`${s.key}-${idx}`}
          cx={cx}
          cy={cy}
          r={(R + r) / 2}
          fill="none"
          stroke={fillCol}
          strokeWidth={R - r}
          className="cursor-pointer transition-opacity hover:opacity-85"
          onMouseMove={(e) => onHover(e, s, sideLabel)}
          onMouseLeave={onLeave}
        />
      );
    } else {
      elements.push(
        <path
          key={`${s.key}-${idx}`}
          d={buildRingPath(cx, cy, R, r, ang, a2)}
          fill={fillCol}
          stroke="var(--md-sys-color-surface-container)"
          strokeWidth={2}
          strokeLinejoin="round"
          className="cursor-pointer transition-opacity hover:opacity-85"
          onMouseMove={(e) => onHover(e, s, sideLabel)}
          onMouseLeave={onLeave}
        />
      );
    }

    const mid = (ang + a2) / 2;
    if (frac >= 0.15) {
      const [tx, ty] = polar(cx, cy, mid, (R + r) / 2);
      const shortName =
        s.label === "Ritel dipakai Bandar"
          ? ["Ritel dipakai", "Bandar"]
          : s.label === "Underwriter"
          ? ["UW"]
          : s.label === "Jumlah IPO"
          ? ["IPO"]
          : [s.label];
      const lines = [...shortName, labelTxt];
      const y0 = ty - 1 - (lines.length - 2) * 4.5;
      lines.forEach((ln, li) => {
        elements.push(
          <text
            key={`${s.key}-in-${li}`}
            x={tx.toFixed(1)}
            y={(y0 + li * 10).toFixed(1)}
            textAnchor="middle"
            fontSize="8.5"
            fontWeight={700}
            fill="#ffffff"
            className="pointer-events-none select-none"
          >
            {ln}
          </text>
        );
      });
    } else if (frac >= 0.045) {
      outer.push({ mid, label: s.key === "ipo" ? "IPO" : s.label, txt: labelTxt });
    }
    ang = a2;
  });

  // Outer labels with collision avoidance
  const items = outer.map((o) => {
    const [ex, ey] = polar(cx, cy, o.mid, R);
    const [lx, ly] = polar(cx, cy, o.mid, R + 10);
    return { ...o, ex, ey, lx, ly, right: Math.cos(o.mid) >= 0, ty: ly };
  });
  items.forEach((o) => {
    if (gapRight !== undefined && o.right === gapRight && Math.abs(o.ty - cy) < 20) {
      o.ty = o.ty < cy ? cy - 22 : cy + 22;
    }
  });
  [true, false].forEach((sideFlag) => {
    const grp = items.filter((i) => i.right === sideFlag).sort((a, b) => a.ty - b.ty);
    for (let i = 1; i < grp.length; i++) {
      if (grp[i].ty - grp[i - 1].ty < 11) grp[i].ty = grp[i - 1].ty + 11;
    }
  });
  items.forEach((o, idx) => {
    const hx = o.lx + (o.right ? 8 : -8);
    elements.push(
      <g key={`out-${idx}`} className="pointer-events-none select-none">
        <polyline
          points={`${o.ex.toFixed(1)},${o.ey.toFixed(1)} ${o.lx.toFixed(1)},${o.ly.toFixed(1)} ${hx.toFixed(1)},${o.ty.toFixed(1)}`}
          fill="none"
          stroke="var(--md-sys-color-outline)"
          strokeWidth={1}
        />
        <text
          x={(hx + (o.right ? 3 : -3)).toFixed(1)}
          y={(o.ty + 3).toFixed(1)}
          textAnchor={o.right ? "start" : "end"}
          fontSize="8.5"
          fontWeight={600}
          fill="var(--md-sys-color-on-surface-variant)"
        >
          {o.label} {o.txt}
        </text>
      </g>
    );
  });

  // Center text
  elements.push(
    <text
      key="ct"
      x={cx}
      y={centerSub === false ? cy + 4 : cy - 1}
      textAnchor="middle"
      fontSize="10"
      fontWeight={800}
      fill="var(--md-sys-color-on-surface)"
      className="pointer-events-none select-none"
    >
      {centerTitle}
    </text>
  );
  if (centerSub !== false) {
    const subStr =
      centerSub ||
      (total >= 1000 ? `${fmtInt(Math.round(total / 1000))} k LOT` : `${fmtInt(Math.round(total))} LOT`);
    elements.push(
      <text
        key="cs"
        x={cx}
        y={cy + 11}
        textAnchor="middle"
        fontSize="8.5"
        fontWeight={700}
        fill="var(--md-sys-color-on-surface-variant)"
        className="pointer-events-none select-none"
      >
        {subStr}
      </text>
    );
  }

  return <>{elements}</>;
}
