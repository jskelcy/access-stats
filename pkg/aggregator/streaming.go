package aggregator

import (
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
	sync.RWMutex
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
// not need to remember to call public Start in a go rountine.
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
			go s.parseEvent(event)
		case <-ticker.C:
			go s.flush()
		}
	}
}

// This function is thread safe.
func (s *streaming) parseEvent(event types.Event) {
	logLine := string(event.Data)

	// While Ingest is a write operation it will get the read lock
	// even though all operations on a Block are threadsafe.
	// The lock is required because streaming.flush switches out the currBlock,
	// and flush calls the standard Lock. This means Ingest calls can not be made on
	// the currBlock until flush has completed and a new block has been switched in.
	s.RLock()
	s.currBlock.Ingest(logLine)
	s.RUnlock()
}

// Swaps out old Block with a new Block. Reports stats for old block.
func (s *streaming) flush() {
	s.Lock()
	oldBlock := s.currBlock
	s.currBlock = types.NewBlock()
	s.Unlock()

	s.reporter.Report(oldBlock)
}
