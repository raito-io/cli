// Package slice contains some utility functions to work with slices
package slice

import (
	"slices"
	"strings"
)

// StringSliceDifference returns the elements in `a` that aren't in `b`.
func StringSliceDifference(a, b []string, caseSensitive bool) []string {
	mb := make(map[string]struct{}, len(b))

	for _, x := range b {
		if !caseSensitive {
			x = strings.ToUpper(x)
		}
		mb[x] = struct{}{}
	}

	var diff []string

	for _, o := range a {
		x := o

		if !caseSensitive {
			x = strings.ToUpper(x)
		}

		if _, found := mb[x]; !found {
			diff = append(diff, o)
		}
	}

	return diff
}

// SliceDifference returns the elements in `a` that aren't in `b`.
func SliceDifference[T comparable](a, b []T) []T {
	mb := make(map[T]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []T

	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}

// StringSliceMerge merges multiple slices of strings together and removes duplicates
func StringSliceMerge(slices ...[]string) []string {
	uniqueMap := map[string]bool{}

	for _, s := range slices {
		for _, v := range s {
			uniqueMap[v] = true
		}
	}

	result := make([]string, 0, len(uniqueMap))

	for key := range uniqueMap {
		result = append(result, key)
	}

	return result
}

func ParseCommaSeparatedList(list string) []string {
	list = strings.TrimSpace(list)

	if list == "" {
		return []string{}
	}

	split := strings.Split(list, ",")
	ret := make([]string, 0, len(split))

	for _, v := range split {
		v = strings.TrimSpace(v)
		if !slices.Contains(ret, v) {
			ret = append(ret, v)
		}
	}

	return ret
}
