package types

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistogramMax(t *testing.T) {
	h := NewHistogram()
	tuples := []struct {
		name  string
		value int
	}{
		{"foo", 1},
		{"bar", 3},
		{"baz", 2},
		{"foo", 2},
	}
	var wg sync.WaitGroup
	for _, tuple := range tuples {
		wg.Add(1)
		go func(n string, v int) {
			defer wg.Done()
			h.Add(n, v)
		}(tuple.name, tuple.value)
	}

	wg.Wait()
	expected := []DataPoint{
		{"bar", 3},
		{"foo", 3},
	}
	assert.Equal(t, expected, h.Max())
}

func TestHistogramNPercentile(t *testing.T) {
	h := NewHistogram()
	tuples := []struct {
		name  string
		value int
	}{
		{"user1", 30},
		{"user2", 33},
		{"user3", 43},
		{"user4", 53},
		{"user5", 56},
		{"user6", 67},
		{"user7", 68},
		{"user8", 72},
		{"user9", 82},
		{"user10", 99},
	}
	var wg sync.WaitGroup
	for _, tuple := range tuples {
		wg.Add(1)
		go func(n string, v int) {
			defer wg.Done()
			h.Add(n, v)
		}(tuple.name, tuple.value)
	}

	wg.Wait()
	expected50th := []DataPoint{
		{"user6", 67},
		{"user7", 68},
		{"user8", 72},
		{"user9", 82},
		{"user10", 99},
	}
	assert.Equal(t, expected50th, h.NPercentile(50))
	expected75th := []DataPoint{
		{"user8", 72},
		{"user9", 82},
		{"user10", 99},
	}
	assert.Equal(t, expected75th, h.NPercentile(75))
}
