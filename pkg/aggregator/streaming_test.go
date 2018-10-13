package aggregator

import (
	"testing"

	"github.com/jskelcy/access-stats/pkg/alerts"

	"github.com/jskelcy/access-stats/pkg/report"
	"github.com/jskelcy/access-stats/pkg/types"
)

func TestParseEvent(t *testing.T) {
	s := &streaming{
		currBlock: types.NewBlock(),
		reporter: report.NewReporter(report.ReporterConfig{
			Alerter: alerts.NewMovingAvgAlerter(alerts.MovingAvgAlerterConfig{
				AlertThreshold: 100,
				AlertWindow:    100,
				AggWindow:      10,
			}),
		}),
	}

	for _, fixture := range fixtures() {
		s.parseEvent(types.Event{
			Data: fixture,
			Err:  nil,
		})
	}

	s.flush()
}

func fixtures() [][]byte {
	return [][]byte{
		[]byte(`127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 1234`),
		[]byte(`127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /foo/bar/baz HTTP/1.0" 200 1234`),
		[]byte(`127.0.0.1 - jill [09/May/2018:16:00:41 +0000] "GET /api/user HTTP/1.0" 500 1234`),
		[]byte(`127.0.0.1 - frank [09/May/2018:16:00:42 +0000] "GET /api/cats HTTP/1.0" 200 1234`),
		[]byte(`127.0.0.1 - mary [09/May/2018:16:00:42 +0000] "GET /api/dogs HTTP/1.0" 200 1234`),
		[]byte(`127.0.0.1 - mary [09/May/2018:16:00:42 +0000] "GET /api/fish TTP/1.0" 500 1234`),
		[]byte(`127.0.0.1 - kim [09/May/2018:16:00:42 +0000] "GET /api/dogs HTTP/1.0" 500 1234`),
	}
}
