package report

import (
	"fmt"
	"time"

	"github.com/jskelcy/access-stats/pkg/alerts"
	"github.com/jskelcy/access-stats/pkg/types"
)

const (
	highTrafficTmpl       = "*** High traffic generated an alert - hits = {%v}, triggered at {%v}\n"
	alertingTmpl          = "*** High traffic ongoing alert - hits = {%v}\n"
	recoveredTmpl         = "*** High traffic alert recovered - hits = {%v}, recovered at {%v}\n"
	maxSectionHitsTmpl    = "Highest traffic Section(s):\n"
	max2XXSectionHitsTmpl = "Highest traffic 2XX Section(s):\n"
	max5XXSectionHitsTmpl = "Highest traffic 5XX Section(s):\n"
	p755XXSectionHitsTmpl = "Sections in 75th percentil of 5XXs:\n"
	maxUserTmpl           = "Highest traffic User(s): \n"
	sectionDataPointTmpl  = "	- section: %v  hits: %v\n"
	userDataPointTmpl     = "	- user: %v  hits: %v\n"
	endOfReport           = "==========================\n\n"
)

// Reporter reports block information and alerts
type Reporter interface {
	Report(*types.Block) error
}

// ReporterConfig is the configuration for a reporter.
type ReporterConfig struct {
	Alerter alerts.Alerter
}

type reporter struct {
	alerter alerts.Alerter
}

// NewReporter retunrs a reporter from config.
func NewReporter(cfg ReporterConfig) Reporter {
	return &reporter{
		alerter: cfg.Alerter,
	}
}

// Report prints block stats and alerts
func (r *reporter) Report(block *types.Block) error {
	r.reportAlert(block.Total)

	// report max section
	fmt.Printf(maxSectionHitsTmpl)
	r.reportSectionDPS(block.HistSection.Max())
	fmt.Println()

	// report max 2XX section
	fmt.Printf(max2XXSectionHitsTmpl)
	r.reportSectionDPS(block.Hist2XX.Max())
	fmt.Println()

	// report max 5XX section
	fmt.Printf(max5XXSectionHitsTmpl)
	r.reportSectionDPS(block.Hist5XX.Max())
	fmt.Println()

	// report sections in 75th percentile of 5XXs
	fmt.Printf(p755XXSectionHitsTmpl)
	r.reportSectionDPS(block.Hist5XX.NPercentile(75))
	fmt.Println()

	// report user with most traffic
	fmt.Printf(maxUserTmpl)
	r.reportUserDPS(block.HistUser.Max())
	fmt.Println()

	fmt.Printf(endOfReport)
	return nil
}

func (r *reporter) reportAlert(value int) {
	avgQPS, alertStatus := r.alerter.Ingest(value)

	switch alertStatus {
	case alerts.NewAlert:
		fmt.Printf(highTrafficTmpl, avgQPS, time.Now().Format(time.UnixDate))
	case alerts.Alerting:
		fmt.Printf(alertingTmpl, avgQPS)
	case alerts.Recovered:
		fmt.Printf(recoveredTmpl, avgQPS, time.Now().Format(time.UnixDate))
	}
}

func (r *reporter) reportSectionDPS(dps []types.DataPoint) {
	for _, dp := range dps {
		fmt.Printf(sectionDataPointTmpl, dp.Name, dp.Hits)
	}
}

func (r *reporter) reportUserDPS(dps []types.DataPoint) {
	for _, dp := range dps {
		fmt.Printf(userDataPointTmpl, dp.Name, dp.Hits)
	}
}
