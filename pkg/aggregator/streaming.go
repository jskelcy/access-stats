package aggregator

import (
	"log"
	"sync"
	"time"

	"github.com/jskelcy/access-stats/pkg/report"
	"github.com/jskelcy/access-stats/pkg/types"
)

// Streaming is an aggregator which takes an event stream in the
// form of a channel. Every 10 seconds blocks are flushed for analytics
// and alerting.
type Streaming interface {
	Start(<-chan types.Event)
}

// StreamingConfig config for streaming aggregator.
type StreamingConfig struct {
	Reporter report.Reporter
}

type streaming struct {
	sync.Mutex
	currBlock *types.Block
	reporter  report.Reporter
}

// NewStreaming returns a streaming aggregator from config.
func NewStreaming(cfg StreamingConfig) Streaming {
	return &streaming{
		currBlock: types.NewBlock(),
		reporter:  cfg.Reporter,
	}
}

// Start calls private start in a go routine so caller does
// not need to remember to call public start in a go rountine.
func (s *streaming) Start(eventStream <-chan types.Event) {
	go s.start(eventStream)
}

func (s *streaming) start(eventStream <-chan types.Event) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case event, ok := <-eventStream:
			if !ok {
				return
			}
			if event.Err != nil {
				log.Println(event.Err)
				break
			}
			s.parseEvent(event)
		case <-ticker.C:
			s.flush()
		}
	}
}

// This function is thread safe.
func (s *streaming) parseEvent(event types.Event) {
	logLine := string(event.Data)
	s.currBlock.Ingest(logLine)
}

func (s *streaming) flush() {
	s.Lock()
	oldBlock := s.currBlock
	s.currBlock = types.NewBlock()
	s.Unlock()

	s.reporter.Report(oldBlock)
}
