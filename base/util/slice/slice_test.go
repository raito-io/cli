package slice

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestStringSliceDifference(t *testing.T) {
	a := []string{"aaa", "bbb", "ccc", "ddd"}
	b := []string{"aaa", "bbb", "ddd"}

	diff := StringSliceDifference(a, b, true)
	assert.Equal(t, 1, len(diff))
	assert.Equal(t, "ccc", diff[0])
}

func TestStringSliceDifferenceEmptyB(t *testing.T) {
	a := []string{"aaa", "bbb", "ccc", "ddd"}
	b := []string{}

	diff := StringSliceDifference(a, b, true)
	assert.Equal(t, 4, len(diff))
	sort.Strings(diff)
	assert.Equal(t, "aaa", diff[0])
	assert.Equal(t, "bbb", diff[1])
	assert.Equal(t, "ccc", diff[2])
	assert.Equal(t, "ddd", diff[3])
}

func TestStringSliceDifferenceEmptyA(t *testing.T) {
	a := []string{}
	b := []string{"aaa", "bbb"}

	diff := StringSliceDifference(a, b, true)
	assert.Equal(t, 0, len(diff))
}

func TestStringSliceDifferenceCaseSensitive(t *testing.T) {
	a := []string{"aaa", "bbb", "ccc", "ddd"}
	b := []string{"aaa", "BBB", "ddd"}

	diff := StringSliceDifference(a, b, true)
	assert.Equal(t, 2, len(diff))
	sort.Strings(diff)
	assert.Equal(t, "bbb", diff[0])
	assert.Equal(t, "ccc", diff[1])
}

func TestStringSliceDifferenceCaseInsensitive(t *testing.T) {
	a := []string{"aaa", "bbb", "ccc", "ddd"}
	b := []string{"aaa", "BBB", "ddd"}

	diff := StringSliceDifference(a, b, false)
	assert.Equal(t, 1, len(diff))
	assert.Equal(t, "ccc", diff[0])
}

func TestSliceDifference(t *testing.T) {
	a := []interface{}{
		G{"aaa", "bbb"},
		G{"ccc", "ddd"},
		G{"eee", "fff"},
	}
	b := []interface{}{
		G{"aaa", "bbb"},
		G{"eee", "fff"},
	}

	diff := SliceDifference(a, b)
	assert.Equal(t, 1, len(diff))
}

func TestStringSliceMerge(t *testing.T) {
	merged := StringSliceMerge([]string{"aaa", "bbb"}, []string{"ccc"})
	assert.Equal(t, 3, len(merged))

	merged = StringSliceMerge([]string{"aaa", "bbb"}, []string{"bbb"})
	assert.Equal(t, 2, len(merged))

	merged = StringSliceMerge([]string{"aaa", "bbb"}, nil)
	assert.Equal(t, 2, len(merged))
}

type G struct {
	F1 string
	F2 string
}
