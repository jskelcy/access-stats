package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistogramMax(t *testing.T) {
	h := NewHistogram()
	h.Add("foo", 1)
	h.Add("bar", 3)
	h.Add("baz", 2)
	h.Add("foo", 2)

	expected := []DataPoint{
		{Name: "bar", Hits: 3},
		{Name: "foo", Hits: 3},
	}
	assert.Equal(t, expected, h.Max())
}

func TestHistogramNPercentile(t *testing.T) {
	h := NewHistogram()
	h.Add("user1", 30)
	h.Add("user2", 33)
	h.Add("user3", 43)
	h.Add("user4", 53)
	h.Add("user5", 56)
	h.Add("user6", 67)
	h.Add("user7", 68)
	h.Add("user8", 72)
	h.Add("user9", 82)
	h.Add("user10", 99)

	fmt.Println(h.NPercentile(50))
}
