package array

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveElementsFromArray(t *testing.T) {
	initialElements := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	indexArray := []int{3, 9, 2}
	expectedElements := []int{0, 1, 4, 5, 6, 7, 8}
	err := RemoveElementsFromArray(indexArray, &initialElements)
	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedElements, initialElements)

	stringElements := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	indexArray = []int{3, 9, 2}
	expectedStringElements := []string{"0", "1", "4", "5", "6", "7", "8"}
	err = RemoveElementsFromArray(indexArray, &stringElements)
	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedStringElements, stringElements)

	// test that the elements in the index array are not larger than the array
	initialElements = []int{0, 1, 2, 3}
	indexArray = []int{3}
	err = RemoveElementsFromArray(indexArray, &initialElements)
	assert.Nil(t, err)

	indexArray = []int{4}
	err = RemoveElementsFromArray(indexArray, &initialElements)
	assert.NotNil(t, err)

	// test that the indices are non-negative
	indexArray = []int{-1}
	err = RemoveElementsFromArray(indexArray, &initialElements)
	assert.NotNil(t, err)
}

func TestSliceFullDiff(t *testing.T) {
	tests := []struct {
		A       []string
		B       []string
		Added   []string
		Removed []string
	}{
		{
			A:       []string{},
			B:       []string{"a", "b"},
			Added:   []string{"a", "b"},
			Removed: []string{},
		},
		{
			A:       []string{"a", "b"},
			B:       []string{},
			Added:   []string{},
			Removed: []string{"a", "b"},
		},
		{
			A:       []string{},
			B:       []string{},
			Added:   []string{},
			Removed: []string{},
		},
		{
			A:       []string{"a", "b"},
			B:       []string{"a", "b"},
			Added:   []string{},
			Removed: []string{},
		},
		{
			A:       []string{"a", "b"},
			B:       []string{"a", "c"},
			Added:   []string{"c"},
			Removed: []string{"b"},
		},
	}

	for _, test := range tests {
		added, removed := SliceFullDiff(test.A, test.B)
		sort.Strings(added)
		sort.Strings(removed)
		assert.Equal(t, test.Added, added)
		assert.Equal(t, test.Removed, removed)
	}
}

func TestSplitIntoBucketsOfSize(t *testing.T) {
	//One bucket
	array := []int{1, 2, 3, 4}
	buckets := SplitIntoBucketsOfSize(5, array)

	assert.Len(t, buckets, 1)
	assert.Equal(t, buckets, [][]int{array})

	//Multiple Buckets
	array = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	buckets = SplitIntoBucketsOfSize(4, array)

	assert.Len(t, buckets, 3)
	assert.Equal(t, buckets, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11}})
}

func TestSetToSlice(t *testing.T) {
	type args[I comparable] struct {
		set map[I]struct{}
	}
	type testCase[I comparable] struct {
		name string
		args args[I]
		want []I
	}
	tests := []testCase[int]{
		{
			name: "empty set",
			args: args[int]{
				set: map[int]struct{}{},
			},
			want: []int{},
		},
		{
			name: "Regular set",
			args: args[int]{
				set: map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
			},
			want: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatchf(t, tt.want, SetToSlice(tt.args.set), "SetToSlice(%v)", tt.args.set)
		})
	}
}

func TestSliceToSet(t *testing.T) {
	type args[I comparable] struct {
		slice []I
	}
	type testCase[I comparable] struct {
		name string
		args args[I]
		want map[I]struct{}
	}
	tests := []testCase[int]{
		{
			name: "empty slice",
			args: args[int]{
				slice: []int{},
			},
			want: map[int]struct{}{},
		},
		{
			name: "Slice without duplicates",
			args: args[int]{
				slice: []int{1, 2, 3, 4},
			},
			want: map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
		},
		{
			name: "Slice with duplicates",
			args: args[int]{
				slice: []int{1, 1, 2, 3, 4, 4},
			},
			want: map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SliceToSet(tt.args.slice), "SliceToSet(%v)", tt.args.slice)
		})
	}
}

func TestAppendToSet(t *testing.T) {
	type args[I comparable] struct {
		sets []map[I]struct{}
	}
	type testCase[I comparable] struct {
		name string
		args args[I]
		want map[I]struct{}
	}
	tests := []testCase[int]{
		{
			name: "empty sets",
			args: args[int]{
				sets: []map[int]struct{}{{}, {}},
			},
			want: map[int]struct{}{},
		},
		{
			name: "Regular sets",
			args: args[int]{
				sets: []map[int]struct{}{{}, {1: {}, 2: {}, 3: {}}, {3: {}, 4: {}, 5: {}}},
			},
			want: map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}, 5: {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, AppendToSet(tt.args.sets...))
		})
	}
}

func TestSetDiff(t *testing.T) {
	type args[I comparable] struct {
		setA map[I]struct{}
		setB map[I]struct{}
	}
	type testCase[I comparable] struct {
		name string
		args args[I]
		want map[I]struct{}
	}
	tests := []testCase[int]{
		{
			name: "empty set",
			args: args[int]{
				setA: map[int]struct{}{},
				setB: map[int]struct{}{},
			},
			want: map[int]struct{}{},
		},
		{
			name: "Regular set",
			args: args[int]{
				setA: map[int]struct{}{1: {}, 2: {}, 3: {}},
				setB: map[int]struct{}{1: {}, 3: {}, 5: {}},
			},
			want: map[int]struct{}{2: {}, 5: {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SetDiff(tt.args.setA, tt.args.setB), "SetDiff(%v, %v)", tt.args.setA, tt.args.setB)
		})
	}
}

func TestKeys(t *testing.T) {
	type args[I comparable, A any] struct {
		m map[I]A
	}
	type testCase[I comparable, A any] struct {
		name string
		args args[I, A]
		want []I
	}
	tests := []testCase[int, int]{
		{
			name: "empty map",
			args: args[int, int]{
				m: map[int]int{3: 5, 5: 6, 7: 8},
			},
			want: []int{3, 5, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatchf(t, tt.want, Keys(tt.args.m), "Keys(%v)", tt.args.m)
		})
	}
}

func TestIntersection(t *testing.T) {
	type intersectionArgs[T comparable] struct {
		a []T
		b []T
	}
	type testIntersectionCase[T comparable] struct {
		name string
		args intersectionArgs[T]
		want []T
	}
	tests := []testIntersectionCase[int]{
		{
			name: "empty sets",
			args: intersectionArgs[int]{
				a: []int{},
				b: []int{},
			},
			want: []int{},
		},

		{
			name: "Equal sets",
			args: intersectionArgs[int]{
				a: []int{1, 2, 3, 4, 5},
				b: []int{1, 2, 3, 4, 5},
			},
			want: []int{1, 2, 3, 4, 5},
		},

		{
			name: "Equal sets random order",
			args: intersectionArgs[int]{
				a: []int{1, 2, 3, 4, 5},
				b: []int{3, 1, 2, 5, 4},
			},
			want: []int{3, 1, 2, 5, 4},
		},

		{
			name: "Diff sets ",
			args: intersectionArgs[int]{
				a: []int{1, 2, 4, 5},
				b: []int{3, 1, 2, 5, 4, 6},
			},
			want: []int{1, 2, 5, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Intersection(tt.args.a, tt.args.b), "Intersection(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}

func TestValues(t *testing.T) {
	type args[I comparable, A any] struct {
		m map[I]A
	}
	type testCase[I comparable, A any] struct {
		name string
		args args[I, A]
		want []A
	}
	tests := []testCase[int, int]{
		{
			name: "empty map",
			args: args[int, int]{
				m: map[int]int{},
			},
			want: []int{},
		},
		{
			name: "simple elements",
			args: args[int, int]{
				m: map[int]int{
					1: 2,
					3: 4,
					5: 6,
					7: 8,
				},
			},
			want: []int{2, 4, 6, 8},
		},
		{
			name: "with duplicated elements",
			args: args[int, int]{
				m: map[int]int{
					1: 2,
					3: 4,
					5: 6,
					7: 2,
				},
			},
			want: []int{2, 4, 6, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Values(tt.args.m)
			assert.Lenf(t, got, len(tt.want), "Values(%v) has %d elements, but expected %d elements", tt.args.m, len(got), len(tt.want))
			assert.ElementsMatchf(t, tt.want, got, "Values(%v) = %v, expected %v", tt.args.m, got, tt.want)
		})
	}
}
