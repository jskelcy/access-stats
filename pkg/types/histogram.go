package types

import (
	"math"
	"sort"
	"sync"
)

// Histogram stores values bucketed by name.
// It supports stats operations like max and percentiles.
// All public operations are thread safe.
type Histogram struct {
	sync.Mutex
	buckets    map[string]*DataPoint
	dataPoints dataPoints
	sorted     bool
}

// DataPoint stores a name and value.
type DataPoint struct {
	Name string
	Hits int
}

// NewHistogram returns new histogram.
func NewHistogram() *Histogram {
	return &Histogram{
		buckets:    map[string]*DataPoint{},
		dataPoints: dataPoints{},
	}
}

// Add adds a value to a histogram.
func (h *Histogram) Add(name string, value int) {
	h.Lock()
	defer h.Unlock()
	dp, ok := h.buckets[name]
	if !ok {
		dp = &DataPoint{
			Name: name,
			Hits: value,
		}
		h.buckets[name] = dp
		h.dataPoints = append(h.dataPoints, dp)
	} else {
		dp.Hits += value
	}
	h.sorted = false
}

// Max returns the max value in the histogram.
// If there is a tie for the max value all DataPoints are returned.
func (h *Histogram) Max() []DataPoint {
	h.Lock()
	defer h.Unlock()
	if !h.sorted {
		h.sort()
	}

	// Return an empty list for an empty histogram.
	if len(h.dataPoints) == 0 {
		return []DataPoint{}
	}

	// Add last value in sorted list
	dps := []DataPoint{
		*h.dataPoints[len(h.dataPoints)-1],
	}
	// walk backwards from end of dataPoints list to capture all
	// DataPoints which are tied for max hits.
	for i := len(h.dataPoints) - 2; i > -1; i-- {
		if h.dataPoints[i].Hits == dps[0].Hits {
			dps = append(dps, *h.dataPoints[i])
			continue
		}
		break
	}

	return dps
}

// NPercentile returns DataPoint values in the Nth percentile.
func (h *Histogram) NPercentile(n int) []DataPoint {
	if n == 100 {
		return h.Max()
	}

	h.Lock()
	defer h.Unlock()
	if !h.sorted {
		h.sort()
	}
	// Return an empty list for an empty histogram.
	if len(h.dataPoints) == 0 || len(h.dataPoints) == 1 {
		return []DataPoint{}
	}

	percentilePosition := (float64(n) / float64(100)) * float64(len(h.dataPoints))

	// For integer values advance percentilePosition +1
	if percentilePosition == math.Trunc(percentilePosition) {
		if int(percentilePosition) < len(h.dataPoints)-1 &&
			h.dataPoints[int(percentilePosition)] != h.dataPoints[int(percentilePosition)+1] {
			percentilePosition++
		}
	}

	// Create copies of data points, this way if returned slice is altered it will
	// not effect data in the histogram.
	dpsPointers := h.dataPoints[int(math.Ceil(percentilePosition))-1:]
	dps := make([]DataPoint, len(dpsPointers))
	for i, dp := range dpsPointers {
		dps[i] = *dp
	}

	return dps
}

// Not thread safe. Get lock before calling.
func (h *Histogram) sort() {
	sort.Sort(h.dataPoints)
	h.sorted = true
}

type dataPoints []*DataPoint

// Sort Interface
// Len is the number of elements in the collection.
func (d dataPoints) Len() int {
	return len(d)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (d dataPoints) Less(i, j int) bool {
	return d[i].Hits < d[j].Hits
}

// Swap swaps the elements with indexes i and j.
func (d dataPoints) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
