import { useMemo, useState } from "react";
import useSWR from "swr";
import { Icon } from "@iconify/react";

import ResearchTabs from "@/components/research/ResearchTabs";
import { EmptyState, ErrorBlock, LoadingBlock } from "@/components/research/ui";
import { fetcher, VIEWS } from "./groupHelpers";
import { useGroupGraphOption } from "./useGroupGraphOption";
import { GroupHeaderCard } from "./GroupHeaderCard";
import { GroupTopGrid } from "./GroupTopGrid";
import { GroupStatCards } from "./GroupStatCards";
import { GroupNetworkView } from "./GroupNetworkView";
import { GroupMatrixView } from "./GroupMatrixView";
import { GroupT5Table } from "./GroupT5Table";
import { GroupConstituentsTable } from "./GroupConstituentsTable";
import type {
  GroupDetail,
  GroupsResponse,
  GroupView,
  NetworkData,
} from "./groupTypes";

export default function GroupDashboard() {
  const [selectedGroup, setSelectedGroup] = useState<string | null>(null);
  const [minPct, setMinPct] = useState(1);
  const [view, setView] = useState<GroupView>("jejaring");
  const [isSyncing, setIsSyncing] = useState(false);
  const [syncStatus, setSyncStatus] = useState<string | null>(null);

  const { data: groupsData, error: groupsError, isLoading: groupsLoading, mutate: mutateGroups } = useSWR<GroupsResponse>(
    "/api/v1/ownership/groups",
    fetcher
  );

  const effectiveGroup = selectedGroup || groupsData?.top_groups[0]?.id || "barito";

  const { data: detail, error: detailError, isLoading: detailLoading, mutate: mutateDetail } = useSWR<GroupDetail>(
    `/api/v1/ownership/groups/${encodeURIComponent(effectiveGroup)}`,
    fetcher
  );

  const handleSyncCatalog = async () => {
    setIsSyncing(true);
    setSyncStatus(null);
    try {
      const resp = await fetch("/api/v1/ownership/groups/sync-catalog", {
        method: "POST",
        credentials: "include",
      });
      const res = await resp.json();
      if (resp.ok && res.status === "success") {
        setSyncStatus(`Sukses: ${res.catalogs_synced} katalog disinkronkan (${res.unique_tickers} ticker).`);
        await mutateGroups();
        await mutateDetail();
      } else {
        setSyncStatus(res.message || "Gagal sinkronisasi katalog");
      }
    } catch {
      setSyncStatus("Gagal menghubungi server");
    } finally {
      setIsSyncing(false);
    }
  };

  const { data: network, error: networkError, isLoading: networkLoading } = useSWR<NetworkData>(
    `/api/bandar/shareholders/network-graph?group=${encodeURIComponent(effectiveGroup)}&min_pct=${minPct}`,
    fetcher
  );

  const graphOption = useGroupGraphOption(network);

  const matrixTickers = useMemo(
    () => (detail?.constituents ?? []).map((item) => item.ticker),
    [detail]
  );

  if (groupsLoading) return <LoadingBlock label="Memuat peta konglomerasi..." />;
  if (groupsError) return <ErrorBlock message={`Gagal memuat grup: ${groupsError.message}`} />;
  if (!groupsData) return <EmptyState title="Data grup tidak tersedia" />;

  return (
    <div className="max-w-7xl mx-auto space-y-4">
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-3">
        <div className="flex items-center gap-3.5">
          <div className="w-11 h-11 sm:w-12 sm:h-12 rounded-2xl bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0 border border-[var(--md-sys-color-outline-variant)] shadow-2xs">
            <Icon icon="ph:share-network-duotone" width="26" />
          </div>
          <div>
            <div className="inline-flex items-center gap-2 text-[9px] font-mono font-bold text-[var(--md-sys-color-primary)] uppercase tracking-wider mb-0.5">
              <span>CONGLOMERATE MAP</span>
              <span className="opacity-40">•</span>
              <span>OWNERSHIP GRAPH</span>
            </div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight">
              Konglomerasi & Jejaring
            </h1>
          </div>
        </div>
        <ResearchTabs />
      </div>

      <GroupHeaderCard
        groupsData={groupsData}
        detail={detail}
        effectiveGroup={effectiveGroup}
        minPct={minPct}
        isSyncing={isSyncing}
        syncStatus={syncStatus}
        onSelectGroup={setSelectedGroup}
        onChangeMinPct={setMinPct}
        onSyncCatalog={handleSyncCatalog}
        onDismissSyncStatus={() => setSyncStatus(null)}
      />

      <GroupTopGrid
        topGroups={groupsData.top_groups}
        effectiveGroup={effectiveGroup}
        onSelectGroup={setSelectedGroup}
      />

      <GroupStatCards group={detail?.group} />

      <div className="border-b border-[var(--md-sys-color-outline-variant)]">
        <div className="flex gap-1.5 overflow-x-auto scrollbar-none pb-2">
          {VIEWS.map((item) => {
            const active = view === item.id;
            return (
              <button
                key={item.id}
                type="button"
                onClick={() => setView(item.id)}
                className={`flex items-center gap-2 rounded-lg border px-3.5 py-2 text-xs font-bold transition-all whitespace-nowrap ${
                  active
                    ? "border-[var(--md-sys-color-primary)] bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)]"
                    : "border-[var(--md-sys-color-outline-variant)] bg-[var(--md-sys-color-surface-container)] text-[var(--md-sys-color-on-surface-variant)] hover:bg-[var(--md-sys-color-surface-container-high)]"
                }`}
              >
                <Icon icon={item.icon} width="15" />
                {item.label}
              </button>
            );
          })}
        </div>
      </div>

      {view === "jejaring" ? (
        <GroupNetworkView
          network={network}
          networkError={networkError}
          networkLoading={networkLoading}
          minPct={minPct}
          graphOption={graphOption}
        />
      ) : null}

      {view === "matriks" ? (
        <div className="space-y-4">
          <GroupMatrixView
            detail={detail}
            detailLoading={detailLoading}
            matrixTickers={matrixTickers}
          />
          <GroupT5Table detail={detail} />
        </div>
      ) : null}

      {view === "emiten" ? (
        <GroupConstituentsTable
          detail={detail}
          detailLoading={detailLoading}
          detailError={detailError}
          groupsData={groupsData}
        />
      ) : null}
    </div>
  );
}
