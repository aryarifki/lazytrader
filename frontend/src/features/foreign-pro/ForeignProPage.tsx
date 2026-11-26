import ForeignProHeader from "./ForeignProHeader";
import ForeignProStockView from "./ForeignProStockView";
import ForeignProIhsgView from "./ForeignProIhsgView";
import ForeignProScreenerView from "./ForeignProScreenerView";
import ForeignProParameterDrawer from "./ForeignProParameterDrawer";
import { useForeignProState } from "./useForeignProState";

export default function ForeignProPage() {
  const s = useForeignProState();

  return (
    <div className="space-y-4 max-w-7xl mx-auto overflow-x-hidden">
      <ForeignProHeader
        activeMainTab={s.activeMainTab}
        setActiveMainTab={s.setActiveMainTab}
        tickerSearch={s.tickerSearch}
        setTickerSearch={s.setTickerSearch}
        suggestOpen={s.suggestOpen}
        setSuggestOpen={s.setSuggestOpen}
        suggestions={s.suggestions}
        onSelectTicker={s.handleSelectTicker}
        onOpenParameters={() => s.setParameterDrawerOpen(true)}
        onRefresh={s.handleRefresh}
      />

      {s.activeMainTab === "stock" && (
        <ForeignProStockView
          ticker={s.ticker}
          tickerSearch={s.tickerSearch}
          setTickerSearch={s.setTickerSearch}
          suggestOpen={s.suggestOpen}
          setSuggestOpen={s.setSuggestOpen}
          suggestions={s.suggestions}
          onSelectTicker={s.handleSelectTicker}
          onRefresh={() => s.loadStockSeries(s.ticker)}
          loading={s.loading}
          isLightMode={s.isLightMode}
          stockData={s.stockData}
          latestFp={s.latestFp}
          akumDays={s.akumDays}
          setAkumDays={s.setAkumDays}
          akumulasiResult={s.akumulasiResult}
          valFFText={s.valFFText}
          valFCoText={s.valFCoText}
          valKFText={s.valKFText}
          valXLText={s.valXLText}
          valSRText={s.valSRText}
          latestPartisipasi={s.latestPartisipasi}
          latestKat={s.latestKat}
          latestKfh={s.latestKfh}
          latestKret={s.latestKret}
          latestCorr={s.latestCorr}
          latestCorrKF={s.latestCorrKF}
          latestValM={s.latestValM}
          showFF={s.showFF}
          setShowFF={s.setShowFF}
          showFCo={s.showFCo}
          setShowFCo={s.setShowFCo}
          showKF={s.showKF}
          setShowKF={s.setShowKF}
          showXL={s.showXL}
          setShowXL={s.setShowXL}
          showSemiRitel={s.showSemiRitel}
          setShowSemiRitel={s.setShowSemiRitel}
          showMA={s.showMA}
          setShowMA={s.setShowMA}
          maPeriod={s.maPeriod}
          setMaPeriod={s.setMaPeriod}
          maTarget={s.maTarget}
          setMaTarget={s.setMaTarget}
          additionalOverlay={s.additionalOverlay}
          setAdditionalOverlay={s.setAdditionalOverlay}
          statType={s.statType}
          setStatType={s.setStatType}
          signalStat={s.signalStat}
          setSignalStat={s.setSignalStat}
          showPBNyangkut={s.showPBNyangkut}
          setShowPBNyangkut={s.setShowPBNyangkut}
          notifOpen={s.notifOpen}
          setNotifOpen={s.setNotifOpen}
          histTopMode={s.histTopMode}
          setHistTopMode={s.setHistTopMode}
          histBottomMode={s.histBottomMode}
          setHistBottomMode={s.setHistBottomMode}
          zoomDays={s.zoomDays}
          setZoomDays={s.setZoomDays}
          windowOffset={s.windowOffset}
          totalStockPoints={s.totalStockPoints}
          effectiveWindowStart={s.effectiveWindowStart}
          effectiveWindowEnd={s.effectiveWindowEnd}
          datesSlice={s.datesSlice}
          opensSlice={s.opensSlice}
          closesSlice={s.closesSlice}
          lowsSlice={s.lowsSlice}
          highsSlice={s.highsSlice}
          cumFFSlice={s.cumFFSlice}
          cumFCoSlice={s.cumFCoSlice}
          cumKFSlice={s.cumKFSlice}
          cumXLSlice={s.cumXLSlice}
          cumSRSlice={s.cumSRSlice}
          cumLNRSlice={s.cumLNRSlice}
        />
      )}

      {s.activeMainTab === "ihsg" && (
        <ForeignProIhsgView
          ihsgData={s.ihsgData}
          ihsgMacroMetrics={s.ihsgMacroMetrics}
          isLightMode={s.isLightMode}
          catIhsg={s.catIhsg}
          setCatIhsg={s.setCatIhsg}
          catKonglo={s.catKonglo}
          setCatKonglo={s.setCatKonglo}
          catFund={s.catFund}
          setCatFund={s.setCatFund}
          metFF={s.metFF}
          setMetFF={s.setMetFF}
          metFCo={s.metFCo}
          setMetFCo={s.setMetFCo}
          metKF={s.metKF}
          setMetKF={s.setMetKF}
          metXL={s.metXL}
          setMetXL={s.setMetXL}
          metSR={s.metSR}
          setMetSR={s.setMetSR}
          flipKF={s.flipKF}
          setFlipKF={s.setFlipKF}
          flipXL={s.flipXL}
          setFlipXL={s.setFlipXL}
          flipSR={s.flipSR}
          setFlipSR={s.setFlipSR}
        />
      )}

      {s.activeMainTab === "screener" && (
        <ForeignProScreenerView
          universe={s.universe}
          setUniverse={s.setUniverse}
          windowLookback={s.windowLookback}
          setWindowLookback={s.setWindowLookback}
          screenerSubTab={s.screenerSubTab}
          setScreenerSubTab={s.setScreenerSubTab}
          screenerSearch={s.screenerSearch}
          setScreenerSearch={s.setScreenerSearch}
          filteredScreenerList={s.filteredScreenerList}
          activeStrategyId={s.activeStrategyId}
          setActiveStrategyId={s.setActiveStrategyId}
          strategyDesc={s.strategyDesc}
          setStrategyDesc={s.setStrategyDesc}
          activeReportType={s.activeReportType}
          setActiveReportType={s.setActiveReportType}
          periodReportWindow={s.periodReportWindow}
          setPeriodReportWindow={s.setPeriodReportWindow}
          customCategory={s.customCategory}
          setCustomCategory={s.setCustomCategory}
          transBucket={s.transBucket}
          setTransBucket={s.setTransBucket}
          priceClass={s.priceClass}
          setPriceClass={s.setPriceClass}
          stateOnly={s.stateOnly}
          setStateOnly={s.setStateOnly}
          hfpOnly={s.hfpOnly}
          setHfpOnly={s.setHfpOnly}
          mfpOnly={s.mfpOnly}
          setMfpOnly={s.setMfpOnly}
          adminMode={s.adminMode}
          setAdminMode={s.setAdminMode}
          adminReportText={s.adminReportText}
          adminCopied={s.adminCopied}
          onCopyReport={() => {
            navigator.clipboard.writeText(s.adminReportText);
            s.setAdminCopied(true);
            setTimeout(() => s.setAdminCopied(false), 2000);
          }}
          ihsgData={s.ihsgData}
          marketReportData={s.marketReportData}
          periodReportData={s.periodReportData}
          accumHeaderLabel={s.accumHeaderLabel}
          isLightMode={s.isLightMode}
          onSelectStock={(tk) => {
            s.handleSelectTicker(tk);
            s.setActiveMainTab("stock");
          }}
        />
      )}

      <ForeignProParameterDrawer
        open={s.parameterDrawerOpen}
        onClose={() => s.setParameterDrawerOpen(false)}
        showFF={s.showFF}
        setShowFF={s.setShowFF}
        showFCo={s.showFCo}
        setShowFCo={s.setShowFCo}
        showKF={s.showKF}
        setShowKF={s.setShowKF}
        showXL={s.showXL}
        setShowXL={s.setShowXL}
        showSemiRitel={s.showSemiRitel}
        setShowSemiRitel={s.setShowSemiRitel}
        zoomDays={s.zoomDays}
        setZoomDays={s.setZoomDays}
        akumDays={s.akumDays}
        setAkumDays={s.setAkumDays}
        showPBNyangkut={s.showPBNyangkut}
        setShowPBNyangkut={s.setShowPBNyangkut}
      />
    </div>
  );
}
