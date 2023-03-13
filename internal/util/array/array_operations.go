package array

import (
	"fmt"
	"math"
	"sort"
)

func RemoveElementsFromArray[T any](indicesToRemove []int, array *[]T) error {
	sort.Ints(indicesToRemove)
	arrayLength := len(*array)

	maxIndex := indicesToRemove[len(indicesToRemove)-1]
	if maxIndex >= arrayLength {
		return fmt.Errorf("maximum index (%d) is greater than array length (%d)", maxIndex, arrayLength)
	} else if indicesToRemove[0] < 0 {
		return fmt.Errorf("index array can't have elements smaller than 0")
	}

	for i := len(indicesToRemove) - 1; i >= 0; i-- {
		index := indicesToRemove[i]
		// the emptyObject is necessary to help the garbage collector
		var emptyObject T
		(*array)[arrayLength-1], (*array)[index] = emptyObject, (*array)[arrayLength-1]
		*array = (*array)[:arrayLength-1]
		arrayLength--
	}

	return nil
}

// SliceFullDiff returns the full difference between the two slices.
// The first return value are the elements that were added (in b but not in a), the second are the elements that were removed (in a but not in b)
func SliceFullDiff(a, b []string) ([]string, []string) {
	added := SliceDiff(b, a)
	removed := SliceDiff(a, b)

	return added, removed
}

// SliceDiff returns the difference between the two slices.
// It returns the elements that appear in a but not in b
func SliceDiff[T comparable](a, b []T) []T {
	mb := SliceToSet(b)

	return SliceSetDiff(a, mb)
}

func SliceSetDiff[T comparable](a []T, b map[T]struct{}) []T {
	diff := make([]T, 0, len(a))

	for _, x := range a {
		if _, found := b[x]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}

func SplitIntoBucketsOfSize[T any](bucketSize int, array []T) [][]T {
	var buckets [][]T
	chucks := int(math.Ceil(float64(len(array)) / float64(bucketSize)))
	startItem := 0
	endItem := bucketSize

	buckets = make([][]T, chucks)

	for i := 0; i < chucks; i++ {
		if endItem > len(array) {
			endItem = len(array)
		}

		buckets[i] = array[startItem:endItem]
		startItem = endItem
		endItem += bucketSize
	}

	return buckets
}

func Intersection[T comparable](a, b []T) []T {
	ma := SliceToSet(a)

	return SetSliceIntersection(ma, b)
}

func SetSliceIntersection[T comparable](a map[T]struct{}, b []T) []T {
	result := make([]T, 0, int(math.Min(float64(len(a)), float64(len(b)))))

	for i := range b {
		if _, found := a[b[i]]; found {
			result = append(result, b[i])
		}
	}

	return result
}

func ChannelToArray[I any](c chan I) []I {
	output := make([]I, 0)

	for item := range c {
		output = append(output, item)
	}

	return output
}

func ArrayToChannel[I any](array []I) chan I {
	ch := make(chan I)

	go func() {
		defer close(ch)

		for _, item := range array {
			ch <- item
		}
	}()

	return ch
}

func SetToSlice[I comparable](set map[I]struct{}) []I {
	result := make([]I, 0, len(set))

	for item := range set {
		result = append(result, item)
	}

	return result
}

func SliceToSet[I comparable](slice []I) map[I]struct{} {
	result := make(map[I]struct{})

	for _, item := range slice {
		result[item] = struct{}{}
	}

	return result
}

func AppendToSet[I comparable](sets ...map[I]struct{}) map[I]struct{} {
	result := make(map[I]struct{})

	for _, set := range sets {
		for item := range set {
			result[item] = struct{}{}
		}
	}

	return result
}

func SetDiff[I comparable](setA map[I]struct{}, setB map[I]struct{}) map[I]struct{} {
	result := make(map[I]struct{})

	for item := range setA {
		if _, found := setB[item]; !found {
			result[item] = struct{}{}
		}
	}

	for item := range setB {
		if _, found := setA[item]; !found {
			result[item] = struct{}{}
		}
	}

	return result
}

func Keys[I comparable, A any](m map[I]A) []I {
	keys := make([]I, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

func Values[I comparable, A any](m map[I]A) []A {
	values := make([]A, 0, len(m))

	for _, value := range m {
		values = append(values, value)
	}

	return values
}

func Map[I any, O any](inputSlice []I, mapFn func(*I) O) []O {
	outputSlice := make([]O, 0, len(inputSlice))

	for i := range inputSlice {
		outputSlice = append(outputSlice, mapFn(&inputSlice[i]))
	}

	return outputSlice
}
