package forensic

import (
	"math"
	"strings"

	"github.com/aryarifki/lazytrader/backend/internal/core/flags"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
)

const (
	crash1DThreshold       = -2.50e12
	crash3DThreshold       = -3.00e12
	warning1DThreshold     = -1.00e12
	warning5DThreshold     = -2.00e12
	convictionFloor        = 300_000_000_000.0
	climacticThreshold     = 1.20e12
	exitQuietSessions      = 8
	recoveryRelFraction    = 0.60
	recoveryFloorFlow      = 1.00e12
	convictionRelFraction  = 0.15
	exodusModeDefault      = "absolute"
	exodusModeHybrid       = "hybrid_p5_p15"
	exodusThresholdWindow  = 252
	exodusThresholdMinObs  = 60
	exodusPctl3D           = 5.0
	exodusPctl5D           = 15.0
)

// ExodusLookbackLimit mirrors routers/forensic.py _exodus_lookback_limit.
func ExodusLookbackLimit(days int, thresholdMode string) int {
	if strings.ToLower(strings.TrimSpace(thresholdMode)) == exodusModeHybrid {
		return 5000
	}
	return days
}

// EvaluateMarketState is an exact 1:1 port of routers/forensic.py _evaluate_market_state.
// It sets MarketState on each TimelinePoint and returns the final market_state string.
func EvaluateMarketState(timeline []TimelinePoint, flagStore *flags.Store) string {
	if len(timeline) == 0 {
		return "NORMAL"
	}

	exitFix := true
	recoveryFix := true
	thresholdMode := exodusModeDefault
	if flagStore != nil {
		exitFix = flagStore.GetBool("EXODUS_EXIT_FIX", true)
		recoveryFix = flagStore.GetBool("EXODUS_RECOVERY_FIX", true)
		thresholdMode = strings.ToLower(strings.TrimSpace(flagStore.GetString("EXODUS_THRESHOLD_MODE", exodusModeDefault)))
	}

	flow3DHistory := make([]float64, 0, len(timeline))
	flow5DHistory := make([]float64, 0, len(timeline))
	inCrashMode := false
	consecPositive := 0
	cumulativeExodus := 0.0
	state := "NORMAL"
	best5DInflow := 0.0

	recoveryTestCounter := 0
	bottomRefPrice := 0.0
	bottomRefExodus := 0.0
	cumulativeRecoveryFlow := 0.0

	episodeLow := 0.0
	sessionsSinceLow := 0
	sessionsSinceTrigger := 0

	for i := 0; i < len(timeline); i++ {
		flow := timeline[i].DailyFlow
		ihsg := timeline[i].IHSGClose

		start3D := i - 2
		if start3D < 0 {
			start3D = 0
		}
		flow3D := 0.0
		for j := start3D; j <= i; j++ {
			flow3D += timeline[j].DailyFlow
		}

		start5D := i - 4
		if start5D < 0 {
			start5D = 0
		}
		flow5D := 0.0
		for j := start5D; j <= i; j++ {
			flow5D += timeline[j].DailyFlow
		}
		if flow5D > best5DInflow {
			best5DInflow = flow5D
		}

		baseIHSG := timeline[start5D].IHSGClose
		ihsgDrop5D := 0.0
		if baseIHSG > 0 {
			ihsgDrop5D = ((ihsg - baseIHSG) / baseIHSG) * 100.0
		}

		thr3D := crash3DThreshold
		thr5D := warning5DThreshold
		if thresholdMode == exodusModeHybrid {
			if len(flow3DHistory) >= exodusThresholdMinObs {
				from := len(flow3DHistory) - exodusThresholdWindow
				if from < 0 {
					from = 0
				}
				p := shared.LinearPercentile(flow3DHistory[from:], exodusPctl3D)
				if p < thr3D {
					thr3D = p
				}
			}
			if len(flow5DHistory) >= exodusThresholdMinObs {
				from := len(flow5DHistory) - exodusThresholdWindow
				if from < 0 {
					from = 0
				}
				p := shared.LinearPercentile(flow5DHistory[from:], exodusPctl5D)
				if p < thr5D {
					thr5D = p
				}
			}
		}

		isCrashTrigger := (flow <= crash1DThreshold) || (flow3D <= thr3D)

		if exitFix && !inCrashMode && episodeLow > 0 && ihsg < episodeLow*0.98 {
			inCrashMode = true
			sessionsSinceTrigger = 0
			sessionsSinceLow = 0
			episodeLow = ihsg
		}

		if isCrashTrigger {
			if !inCrashMode || episodeLow <= 0 {
				episodeLow = ihsg
				sessionsSinceLow = 0
			}
			inCrashMode = true
			sessionsSinceTrigger = 0
		} else {
			sessionsSinceTrigger++
		}

		if inCrashMode {
			if ihsg < episodeLow {
				episodeLow = ihsg
				sessionsSinceLow = 0
			} else if !isCrashTrigger {
				sessionsSinceLow++
			}
		}

		if inCrashMode && flow < 0 {
			cumulativeExodus += flow
		}

		if flow > 0 {
			consecPositive++
		} else {
			consecPositive = 0
		}

		switch {
		case isCrashTrigger:
			state = "CRASH"
			recoveryTestCounter = 0
			cumulativeRecoveryFlow = 0.0

		case recoveryTestCounter > 0:
			cumulativeRecoveryFlow += flow
			if ihsg < bottomRefPrice*0.98 {
				state = "CRASH"
				inCrashMode = true
				recoveryTestCounter = 0
				cumulativeExodus = -bottomRefExodus
				cumulativeRecoveryFlow = 0.0
				consecPositive = 0
			} else if recoveryTestCounter == 1 {
				var recoveryThreshold float64
				if recoveryFix {
					recoveryThreshold = math.Max(recoveryFloorFlow, recoveryRelFraction*best5DInflow)
				} else {
					recoveryThreshold = 0.05 * bottomRefExodus
				}
				if cumulativeRecoveryFlow >= recoveryThreshold && ihsg >= bottomRefPrice {
					state = "CONFIRMED_BOTTOM"
					inCrashMode = false
					cumulativeExodus = 0.0
				} else {
					state = "CRASH"
					inCrashMode = true
					cumulativeExodus = -bottomRefExodus
				}
				recoveryTestCounter = 0
			} else {
				recoveryTestCounter--
				state = "RECOVERY_TESTING"
			}

		case inCrashMode && consecPositive >= 2:
			prevFlow := 0.0
			if i > 0 {
				prevFlow = timeline[i-1].DailyFlow
			}
			var convictionThresh float64
			if recoveryFix {
				convictionThresh = math.Max(convictionFloor, convictionRelFraction*best5DInflow)
			} else {
				convictionThresh = math.Max(convictionFloor, 0.03*math.Abs(cumulativeExodus))
			}
			hasConviction := (flow >= convictionThresh) || ((flow + prevFlow) >= 2.0*convictionThresh)

			prevIHSG := ihsg
			if i > 0 {
				prevIHSG = timeline[i-1].IHSGClose
			}
			isClimactic := flow >= climacticThreshold
			priceStabilizing := ihsg >= prevIHSG*0.99

			if hasConviction && (isClimactic || priceStabilizing) {
				state = "BOTTOMING"
				bottomRefPrice = ihsg
				bottomRefExodus = math.Abs(cumulativeExodus)
				recoveryTestCounter = 5
				cumulativeRecoveryFlow = 0.0
			} else {
				state = "CRASH"
			}

		case exitFix && inCrashMode && sessionsSinceTrigger >= exitQuietSessions && sessionsSinceLow >= exitQuietSessions:
			state = "WARNING"
			inCrashMode = false
			cumulativeExodus = 0.0
			consecPositive = 0

		case inCrashMode:
			state = "CRASH"

		case flow5D <= thr5D && math.Abs(ihsgDrop5D) < 3.0:
			state = "WARNING"

		case flow <= warning1DThreshold:
			state = "WARNING"

		default:
			state = "NORMAL"
		}

		timeline[i].MarketState = state
		flow3DHistory = append(flow3DHistory, flow3D)
		flow5DHistory = append(flow5DHistory, flow5D)
	}

	return state
}
