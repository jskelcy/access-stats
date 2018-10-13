package alerts

import (
	"container/heap"
	"time"
)

const (
	// NoAlert status for no current alert
	NoAlert = iota
	// NewAlert status for crossing alert thrshold
	NewAlert
	// Alerting status for currently alerting alert
	Alerting
	// Recovered status for a recovered alert
	Recovered
)

// Status is returned from Ingest to signify current alert state
type Status int

// Alerter ingests values and alerts based when an alert threshold has been
// crossed.
type Alerter interface {
	Ingest(int) (float64, Status)
}

// MovingAvgAlerterConfig holds configuration for movingAvgAlerter
type MovingAvgAlerterConfig struct {
	AlertThreshold float64
	AlertWindow    int
	AggWindow      int
}

// movingAvgAlerter takes in new values and keeps a moving average. Struct is
// not thread safe.
type movingAvgAlerter struct {
	avs            *alertValues
	total          int
	alerting       bool
	alertCounter   int
	alertThreshold float64
	alertWindow    int
	aggWindow      int
}

// NewMovingAvgAlerter returns an movingAvgAlerter from configuration.
func NewMovingAvgAlerter(cfg MovingAvgAlerterConfig) Alerter {
	avs := &alertValues{}
	heap.Init(avs)
	return &movingAvgAlerter{
		avs:            avs,
		alertThreshold: cfg.AlertThreshold,
		alertWindow:    cfg.AlertWindow,
		aggWindow:      cfg.AggWindow,
	}
}

// Ingest adds new value to moving average. First value returned is
// the moving average of qps for the last 2 min, the second value is the alert
// status. If moving average is above the threshold for the configured alert
// window a NewAlert status will be returned. If the alert has been below the
// threshold for the alert window a Recvovered status will be returned.
func (a *movingAvgAlerter) Ingest(value int) (float64, Status) {
	if len(*a.avs) == (a.alertWindow / a.aggWindow) {
		oldValue := heap.Pop(a.avs).(alertValue)
		a.total = a.total - oldValue.value
	}

	newValue := alertValue{
		timestamp: time.Now().UnixNano(),
		value:     value,
	}
	heap.Push(a.avs, newValue)
	a.total = a.total + value
	avgPerSecond := float64(a.total) / (float64(len(*a.avs)) * float64(a.aggWindow))

	// check above threshold
	if avgPerSecond >= a.alertThreshold {
		if a.alerting {
			// Reset the counter if there is average is above threshold
			a.alertCounter = (a.alertWindow / a.aggWindow)
			return avgPerSecond, Alerting
		}

		a.alertCounter++
		// Above threshold for an entire alert window
		// create new alert.
		if a.alertCounter == (a.alertWindow / a.aggWindow) {
			a.alerting = true
			return avgPerSecond, NewAlert
		}
		// No alert yet
		return avgPerSecond, NoAlert
	}

	if !a.alerting {
		// reset the alert counter
		a.alertCounter = 0
		return avgPerSecond, NoAlert
	}

	a.alertCounter--
	// check for recovery
	if a.alertCounter <= 0 {
		// this is incase the value goes below 0
		a.alertCounter = 0
		a.alerting = false
		return avgPerSecond, Recovered
	}
	return avgPerSecond, Alerting
}

type alertValue struct {
	timestamp int64
	value     int
}

// Heap interface implementation to calculate moving average --
type alertValues []alertValue

func (av alertValues) Len() int           { return len(av) }
func (av alertValues) Less(i, j int) bool { return av[i].timestamp < av[j].timestamp }
func (av alertValues) Swap(i, j int)      { av[i], av[j] = av[j], av[i] }

func (av *alertValues) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*av = append(*av, x.(alertValue))
}

func (av *alertValues) Pop() interface{} {
	old := *av
	n := len(old)
	x := old[n-1]
	*av = old[0 : n-1]
	return x
}

// -----------------------------------------------------------
