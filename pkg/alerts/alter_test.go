package alerts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngestAlertAndRecovery(t *testing.T) {
	a := NewMovingAvgAlerter(MovingAvgAlerterConfig{
		AlertThreshold: float64(5),
		AlertWindow:    100,
		AggWindow:      10,
	})

	// Create baseline with no alerts
	for i := 0; i < 10; i++ {
		_, status := a.Ingest(3)
		assert.Equal(t, Status(NoAlert), status)
	}

	// Add values over the threshold
	for i := 0; i < 13; i++ {
		_, status := a.Ingest(100)
		assert.Equal(t, Status(NoAlert), status)
	}

	// trigger alert
	_, status := a.Ingest(100)
	assert.Equal(t, Status(NewAlert), status)

	// add recovery values
	for i := 0; i < 14; i++ {
		_, status := a.Ingest(3)
		assert.Equal(t, Status(Alerting), status)
	}

	// recovered
	_, status = a.Ingest(3)
	assert.Equal(t, Status(Recovered), status)
}
