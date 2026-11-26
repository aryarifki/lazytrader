package bandarmology

import (
	"fmt"
	"time"
)

// PriceRow represents one row from public.prices.
type PriceRow struct {
	Date   time.Time
	Ticker string
	Open   *float64
	High   *float64
	Low    *float64
	Close  *float64
	Volume *int64
}

// FlowRow represents one row from public.broker_flow.
type FlowRow struct {
	Date              time.Time
	Ticker            string
	BandarSignal      *string
	BandarSignalScore *float64
	ForeignNetBroker  *float64
	LocalNetBroker    *float64
	GovNetBroker      *float64
	ForeignNetFlow    *float64
	DomesticNetFlow   *float64
	TotalValue        *float64
	ForeignSignal     *string
	ConclusionBroker  *string
	ConclusionFlow    *string
	FetchedAt         *time.Time
}

// ActivityRow represents one row from public.broker_activity.
type ActivityRow struct {
	Date            time.Time
	Ticker          string
	BrokerCode      string
	ParticipantType string
	BuyValue        float64
	SellValue       float64
	NetValue        float64
	BuyLot          *float64
	SellLot         *float64
	Frequency       *float64
	BuyAvgPrice     *float64
	SellAvgPrice    *float64
	FetchedAt       *time.Time
}

func fmtDateISO(t time.Time) string {
	return t.Format("2006-01-02")
}

func fmtDateMidnightSpace(t time.Time) string {
	return t.Format("2006-01-02 00:00:00")
}

func fmtDateMidnightT(t time.Time) string {
	return t.Format("2006-01-02T00:00:00")
}

func fmtFetchedAtOffset(t *time.Time) *string {
	if t == nil {
		return nil
	}
	u := t.UTC()
	micros := u.Nanosecond() / 1000
	var s string
	if micros == 0 {
		s = u.Format("2006-01-02T15:04:05+00:00")
	} else {
		s = fmt.Sprintf("%s.%06d+00:00", u.Format("2006-01-02T15:04:05"), micros)
	}
	return &s
}

func toDateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
