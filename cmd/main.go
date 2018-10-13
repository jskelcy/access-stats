package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jskelcy/access-stats/pkg/aggregator"
	"github.com/jskelcy/access-stats/pkg/alerts"
	"github.com/jskelcy/access-stats/pkg/collector"
	"github.com/jskelcy/access-stats/pkg/report"
)

func main() {
	watchFile := flag.String("src", "/var/log/access.log", "file to watch for incoming logs")
	qpsAlertThreshold := flag.Int("alertThreshold", 10, "qps where traffic is considered critical")
	flag.Parse()

	alerter := alerts.NewMovingAvgAlerter(alerts.MovingAvgAlerterConfig{
		AlertThreshold: float64(*qpsAlertThreshold),
		AlertWindow:    40,
		// Default to 10 second agg window.
		AggWindow: 10,
	})
	reporter := report.NewReporter(report.ReporterConfig{
		Alerter: alerter,
	})
	streamingAggregator := aggregator.NewStreaming(aggregator.StreamingConfig{
		Reporter: reporter,
	})

	c, err := collector.NewFileWatcher(collector.FileWatcherConfig{
		FileName: *watchFile,
	})
	if err != nil {
		log.Fatal(err)
	}

	eventChan, err := c.Watch()
	if err != nil {
		log.Fatal(err)
	}

	streamingAggregator.Start(eventChan)

	// Handle signals and stop the collector, this will clean up
	// open file watches/descriptors.
	done := make(chan struct{})
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		log.Printf("got sig: %v\n", sig)
		cleanUpChan, err := c.Stop()
		if err != nil {
			log.Println(err)
			done <- struct{}{}
		}
		done <- (<-cleanUpChan)
	}()

	<-done
}
