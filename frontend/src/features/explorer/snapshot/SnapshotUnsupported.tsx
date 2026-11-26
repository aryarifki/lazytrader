import { Icon } from "@iconify/react";

import { Card, DataNote, SectionTitle } from "@/components/research/ui";
import type { UnsupportedCondition } from "./snapshotTypes";

export function SnapshotUnsupported({
  unsupportedConditions,
  notes,
}: {
  unsupportedConditions: UnsupportedCondition[];
  notes: string[];
}) {
  return (
    <Card className="p-4 sm:p-5">
      <SectionTitle
        icon="ph:info-duotone"
        title="Kondisi yang belum tersedia"
        subtitle="Transparansi cakupan data"
      />
      <div className="mt-3 grid grid-cols-1 md:grid-cols-2 gap-2.5">
        {unsupportedConditions.map((item) => (
          <div
            key={item.id}
            className="rounded-xl border border-[var(--md-sys-color-outline-variant)]/60 bg-[var(--md-sys-color-surface-container-high)]/40 p-3"
          >
            <div className="flex items-center gap-1.5 text-[11px] font-bold">
              <Icon
                icon="ph:minus-circle-duotone"
                width="14"
                className="text-[var(--md-sys-color-on-surface-variant)]"
              />
              {item.label}
            </div>
            <p className="mt-1 text-[10px] text-[var(--md-sys-color-on-surface-variant)] leading-relaxed">
              {item.reason}
            </p>
          </div>
        ))}
      </div>
      <div className="mt-3">
        <DataNote tone="warning">
          <ul className="list-disc pl-4 space-y-1">
            {notes.map((note) => (
              <li key={note}>{note}</li>
            ))}
          </ul>
        </DataNote>
      </div>
    </Card>
  );
}
