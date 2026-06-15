package pneumatic

import (
	"fmt"
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCount(t *testing.T) {
	require.Equal(t, 0, Count(Empty[int]()))
	require.Equal(t, 1, Count(Once(7)))
	require.Equal(t, 3, Count(Chain(Once(1), Once(2), Once(3))))
	require.Equal(t, 6, Count(newSeq(1, 1, 2, 3, 5, 8)))
}

func TestLast(t *testing.T) {
	tests := []struct {
		Name          string
		Input         iter.Seq[int]
		ExpectedValue int
		ExpectedOk    bool
	}{
		{"empty", Empty[int](), 0, false},
		{"single", Once(1), 1, true},
		{"multi", newSeq(1, 2, 3), 3, true},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			value, ok := Last(tc.Input)
			require.Equal(t, tc.ExpectedOk, ok)
			require.Equal(t, tc.ExpectedValue, value)
		})
	}
}

func TestNth(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		idx      int
		expected int
		valid    bool
	}{
		{"empty", Empty[int](), 2, 0, false},
		{"first", newSeq(1, 3, 5), 0, 1, true},
		{"middle", newSeq(1, 2, 3), 1, 2, true},
		{"last", newSeq(1, 2, 3), 2, 3, true},
		{"out of range", newSeq(2, 4, 8), 4, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := Nth(tt.idx, tt.input)
			require.Equal(t, tt.valid, ok)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestChain(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		vals := slices.Collect(Chain(newSeq(1, 2, 3), newSeq(4, 5, 6)))
		require.Equal(t, []int{1, 2, 3, 4, 5, 6}, vals)
	})
	t.Run("empty", func(t *testing.T) {
		require.Empty(t, slices.Collect(Chain[int]()))
	})
	t.Run("single", func(t *testing.T) {
		require.Equal(t, []int{1}, slices.Collect(Chain(Once(1))))
	})
	t.Run("mix", func(t *testing.T) {
		vals := slices.Collect(Chain(Once(1), Empty[int](), Empty[int](), newSeq(1, 2, 3, 5)))
		require.Equal(t, []int{1, 1, 2, 3, 5}, vals)
	})
}

func TestFilter(t *testing.T) {
	vals := slices.Collect(
		Filter(func(x int) bool { return x&1 == 0 }, newSeq(1, 2, 3, 4, 5)),
	)
	require.Equal(t, []int{2, 4}, vals)
}

func TestDiscard(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  iter.Seq[int]
		Output []int
	}{
		{"empty", Empty[int](), nil},
		{"none_removed", newSeq(1, 3, 5, 7), []int{1, 3, 5, 7}},
		{"all_removed", newSeq(2, 4, 6, 8), nil},
		{"some_removed", newSeq(1, 3, 4, 4, 6, 5, 9), []int{1, 3, 5, 9}},
	}
	for _, tc := range testCases {
		actual := slices.Collect(Discard(func(i int) bool { return (i & 0x01) == 0 }, tc.Input))
		require.Equal(t, tc.Output, actual)
	}
}

func TestIntersperse(t *testing.T) {
	vals := slices.Collect(Intersperse(0, newSeq(1, 2, 3, 4, 5)))
	require.Equal(t, []int{1, 0, 2, 0, 3, 0, 4, 0, 5}, vals)
}

func TestMap(t *testing.T) {
	vals := slices.Collect(
		Map(func(x int) int { return x * x }, newSeq(1, 2, 3, 4, 5)),
	)
	require.Equal(t, []int{1, 4, 9, 16, 25}, vals)
}

func TestEnumerate(t *testing.T) {
	vals := slices.Collect(
		Map(Pair[int, string].First,
			Enumerate(newSeq("abc", "def", "ghi", "jkl")),
		),
	)
	require.Equal(t, []int{0, 1, 2, 3}, vals)
}

func TestScan(t *testing.T) {
	vals := slices.Collect(
		Scan(0, func(p, n int) int { return p + n }, newSeq(1, 1, 1, 1, 1)),
	)
	require.Equal(t, []int{1, 2, 3, 4, 5}, vals)
}

func TestSkip(t *testing.T) {
	vals := slices.Collect(Skip(2, newSeq(1, 2, 3, 4, 5)))
	require.Equal(t, []int{3, 4, 5}, vals)
}

func TestSkipWhile(t *testing.T) {
	vals := slices.Collect(
		SkipWhile(func(x int) bool { return x&1 != 0 }, newSeq(1, 1, 2, 3, 5, 8)),
	)
	require.Equal(t, []int{2, 3, 5, 8}, vals)
}

func TestTake(t *testing.T) {
	vals := slices.Collect(Take(3, newSeq(1, 2, 3, 4, 5)))
	require.Equal(t, []int{1, 2, 3}, vals)
}

func TestPartition(t *testing.T) {
	tests := []struct {
		name   string
		input  iter.Seq[int]
		group1 []int
		group2 []int
	}{
		{"empty", newSeq[int](), nil, nil},
		{"all true", newSeq(2, 4, 6, 8), []int{2, 4, 6, 8}, nil},
		{"all false", newSeq(1, 3, 5, 7), nil, []int{1, 3, 5, 7}},
		{"mixed", newSeq(1, 2, 3, 4, 5), []int{2, 4}, []int{1, 3, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tVals, fVals := Partition(func(i int) bool { return i&1 == 0 }, tt.input)
			require.Equal(t, tt.group1, tVals)
			require.Equal(t, tt.group2, fVals)
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name   string
		input  iter.Seq[int]
		fn     func(int, int) int
		output int
		valid  bool
	}{
		{"empty", newSeq[int](), func(a, b int) int { return a + b }, 0, false},
		{"sum", newSeq(1, 1, 1, 1), func(a, b int) int { return a + b }, 4, true},
		{"difference", newSeq(4, 1, 1, 1), func(a, b int) int { return a - b }, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := Reduce(tt.fn, tt.input)
			require.Equal(t, tt.valid, ok)
			require.Equal(t, tt.output, actual)
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		pred     func(int) bool
		expected bool
	}{
		{"empty", newSeq[int](), func(i int) bool { return i&1 == 0 }, true},
		{"none", newSeq(1, 3, 5), func(i int) bool { return i&1 == 0 }, false},
		{"some", newSeq(1, 2, 4), func(i int) bool { return i&1 == 0 }, false},
		{"all", newSeq(2, 4, 8), func(i int) bool { return i&1 == 0 }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, All(tt.pred, tt.input))
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		pred     func(int) bool
		expected bool
	}{
		{"empty", newSeq[int](), func(i int) bool { return i&1 == 0 }, false},
		{"none", newSeq(1, 3, 5), func(i int) bool { return i&1 == 0 }, false},
		{"some", newSeq(1, 2, 4), func(i int) bool { return i&1 == 0 }, true},
		{"all", newSeq(2, 4, 8), func(i int) bool { return i&1 == 0 }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, Any(tt.pred, tt.input))
		})
	}
}
func TestFind(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		pred     func(int) bool
		expected int
		valid    bool
	}{
		{"empty", Empty[int](), func(i int) bool { return i&1 == 0 }, 0, false},
		{"none", newSeq(1, 3, 5), func(i int) bool { return i&1 == 0 }, 0, false},
		{"one", newSeq(1, 2, 3), func(i int) bool { return i&1 == 0 }, 2, true},
		{"many", newSeq(2, 4, 8), func(i int) bool { return i&1 == 0 }, 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Find(tt.pred, tt.input)
			require.Equal(t, tt.valid, ok)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		Name  string
		Input iter.Seq[int]
		Max   int
		Ok    bool
	}{
		{"empty", Empty[int](), 0, false},
		{"single", Once(1), 1, true},
		{"multiple", newSeq(2, 4, 1, 0, 3), 4, true},
		{"first_value", newSeq(3, 1, 2), 3, true},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			actual, ok := Max(tc.Input)
			require.Equal(t, tc.Ok, ok)
			require.Equal(t, tc.Max, actual)
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		Name  string
		Input iter.Seq[int]
		Min   int
		Ok    bool
	}{
		{"empty", Empty[int](), 0, false},
		{"single", Once(1), 1, true},
		{"multiple", newSeq(2, 4, 1, 0, 3), 0, true},
		{"first_value", newSeq(1, 3, 2), 1, true},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			actual, ok := Min(tc.Input)
			require.Equal(t, tc.Ok, ok)
			require.Equal(t, tc.Min, actual)
		})
	}
}

func TestRev(t *testing.T) {
	require.Equal(t, []int{5, 4, 3, 2, 1}, slices.Collect(Rev(newSeq(1, 2, 3, 4, 5))))
	require.Equal(t, []int{2, 1}, slices.Collect(Rev(newSeq(1, 2))))
	require.Equal(t, []int{1}, slices.Collect(Rev(Once(1))))
	require.Empty(t, slices.Collect(Rev(Empty[int]())))
}

func TestCycle(t *testing.T) {
	{
		vals := slices.Collect(Take(10, Cycle(newSeq(1, 2, 3))))
		require.Equal(t, []int{1, 2, 3, 1, 2, 3, 1, 2, 3, 1}, vals)
	}
	{
		vals := slices.Collect(Take(3, Cycle(Once(1))))
		require.Equal(t, []int{1, 1, 1}, vals)
	}
}

func TestCartesianProduct(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		vals := slices.Collect(
			Map(func(p Pair[string, int]) string { return fmt.Sprintf("%s%d", p.X, p.Y) },
				CartesianProduct(newSeq("a", "b", "c"), newSeq(1, 2, 3))),
		)
		expected := []string{"a1", "a2", "a3", "b1", "b2", "b3", "c1", "c2", "c3"}
		require.Equal(t, expected, vals)
	})
	t.Run("both_empty", func(t *testing.T) {
		require.Empty(t, slices.Collect(CartesianProduct(Empty[int](), Empty[int]())))
	})
	t.Run("first_empty", func(t *testing.T) {
		require.Empty(t, slices.Collect(CartesianProduct(newSeq(1, 2, 3), Empty[int]())))
	})
	t.Run("second_empty", func(t *testing.T) {
		require.Empty(t, slices.Collect(CartesianProduct(Empty[int](), newSeq(1, 2, 3))))
	})
}

func TestIsSorted(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    iter.Seq[int]
		IsSorted bool
	}{
		{"empty", Empty[int](), true},
		{"single_element", Once(1), true},
		{"sorted_true", newSeq(1, 2, 3), true},
		{"sorted_false", newSeq(3, 1, 2), false},
		{"sorted_false2", newSeq(0, 3, 1, 2), false},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			require.Equal(t, tc.IsSorted, IsSorted(tc.Input))
		})
	}
}

func TestPeel(t *testing.T) {
	testCases := []struct {
		Name  string
		Input iter.Seq[int]
		First int
		Rest  []int
		Ok    bool
	}{
		{"empty", Empty[int](), 0, nil, false},
		{"single_element", Once(1), 1, nil, true},
		{"multiple_elements", newSeq(1, 2, 3), 1, []int{2, 3}, true},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			first, rest, ok := Peel(tc.Input)
			require.Equal(t, tc.Ok, ok)
			require.Equal(t, tc.First, first)
			require.Equal(t, tc.Rest, slices.Collect(rest))
		})
	}
}

func newSeq[T any](t ...T) iter.Seq[T] {
	return slices.Values(t)
}
