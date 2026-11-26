import React from "react";
import MarketReportControlCard from "./MarketReportControlCard";
import MarketReportFloatingDock from "./MarketReportFloatingDock";
import MarketReportLogsPanel from "./MarketReportLogsPanel";
import MarketReportModals from "./MarketReportModals";
import MarketReportViewerCard from "./MarketReportViewerCard";
import { useMarketReportController } from "./useMarketReportController";

export default function MarketReportPage() {
  const ctrl = useMarketReportController();

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 py-6 space-y-6 pb-20 relative">
      <MarketReportFloatingDock
        toast={ctrl.toast}
        status={ctrl.status}
        floatingMinimized={ctrl.floatingMinimized}
        elapsedSeconds={ctrl.elapsedSeconds}
        onSetFloatingMinimized={ctrl.setFloatingMinimized}
        onOpenLogcat={() => {
          ctrl.setShowLogcatModal(true);
          ctrl.fetchFullLogs();
        }}
        onOpenTokenModal={() => ctrl.setShowTokenModal(true)}
        onDismissError={ctrl.handleResetStatus}
      />

      <MarketReportControlCard
        status={ctrl.status}
        history={ctrl.history}
        selectedHistoryFile={ctrl.selectedHistoryFile}
        loading={ctrl.loading}
        actionLoading={ctrl.actionLoading}
        elapsedSeconds={ctrl.elapsedSeconds}
        adminKeyState={ctrl.adminKeyState}
        rateLimit={ctrl.rateLimit}
        workers={ctrl.workers}
        onRateLimitChange={ctrl.setRateLimit}
        onWorkersChange={ctrl.setWorkers}
        onGenerate={ctrl.handleGenerate}
        onRefresh={ctrl.fetchState}
        onAutoRenew={ctrl.handleAutoRenew}
        onResetStatus={ctrl.handleResetStatus}
        onOpenTokenModal={() => ctrl.setShowTokenModal(true)}
        onOpenLogcatModal={() => {
          ctrl.setShowLogcatModal(true);
          ctrl.fetchFullLogs();
        }}
        onOpenAdminKeyModal={() => {
          ctrl.setAdminKeyInput(ctrl.adminKeyState);
          ctrl.setShowAdminKeyModal(true);
        }}
        onSelectHistory={ctrl.handleSelectHistory}
      />

      <MarketReportViewerCard activeReport={ctrl.activeReport} onCopy={ctrl.handleCopy} />

      <MarketReportLogsPanel
        logs={ctrl.logs}
        status={ctrl.status}
        autoScrollLogs={ctrl.autoScrollLogs}
        logLoading={ctrl.logLoading}
        collapsed={ctrl.terminalCollapsed}
        logRef={ctrl.embeddedLogContainerRef}
        onToggleAutoScroll={() => ctrl.setAutoScrollLogs((v) => !v)}
        onRefreshLogs={ctrl.fetchFullLogs}
        onCopyLogs={() => ctrl.handleCopy(ctrl.logs.join("\n"), "Logcat")}
        onOpenFullscreen={() => {
          ctrl.setShowLogcatModal(true);
          ctrl.fetchFullLogs();
        }}
        onToggleCollapsed={() => ctrl.setTerminalCollapsed((v) => !v)}
      />

      <MarketReportModals
        status={ctrl.status}
        actionLoading={ctrl.actionLoading}
        showLogcatModal={ctrl.showLogcatModal}
        logs={ctrl.logs}
        autoScrollLogs={ctrl.autoScrollLogs}
        logLoading={ctrl.logLoading}
        logContainerRef={ctrl.logContainerRef}
        onCloseLogcat={() => ctrl.setShowLogcatModal(false)}
        onToggleAutoScroll={() => ctrl.setAutoScrollLogs((v) => !v)}
        onRefreshLogs={ctrl.fetchFullLogs}
        onCopyLogs={() => ctrl.handleCopy(ctrl.logs.join("\n"), "Logcat")}
        showTokenModal={ctrl.showTokenModal}
        tokenInput={ctrl.tokenInput}
        tokenSubmitting={ctrl.tokenSubmitting}
        tokenError={ctrl.tokenError}
        tokenSuccess={ctrl.tokenSuccess}
        onTokenInputChange={ctrl.setTokenInput}
        onCloseTokenModal={() => {
          ctrl.setShowTokenModal(false);
          ctrl.setTokenError(null);
          ctrl.setTokenSuccess(null);
        }}
        onSubmitToken={ctrl.handleUpdateToken}
        onAutoRenewFromModal={() => {
          ctrl.setShowTokenModal(false);
          ctrl.handleAutoRenew();
        }}
        showAdminKeyModal={ctrl.showAdminKeyModal}
        adminKeyState={ctrl.adminKeyState}
        adminKeyInput={ctrl.adminKeyInput}
        onAdminKeyInputChange={ctrl.setAdminKeyInput}
        onCloseAdminKeyModal={() => ctrl.setShowAdminKeyModal(false)}
        onSubmitAdminKey={ctrl.handleSaveAdminKey}
        onClearAdminKey={ctrl.handleClearAdminKey}
      />
    </div>
  );
}
