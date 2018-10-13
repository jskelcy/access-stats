package collector

import "github.com/jskelcy/access-stats/pkg/types"

// Collector watches data source and returns changes.
type Collector interface {
	Watch() (<-chan types.Event, error)
	Stop() error
}
