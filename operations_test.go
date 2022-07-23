package pneumatic

import (
	"reflect"
	"testing"
)

func Test_chain(t *testing.T) {
	vals := FromSlice([]int{1, 2, 3}).
		Chain(FromSlice([]int{4, 5, 6})).
		Collect()
	expected := []int{1, 2, 3, 4, 5, 6}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}

func Test_filter(t *testing.T) {
	vals := FromSlice([]int{1, 2, 3, 4, 5}).
		Filter(func(x int) bool { return x&1 == 0 }).
		Collect()
	expected := []int{2, 4}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}

func Test_map(t *testing.T) {
	vals := FromSlice([]int{1, 2, 3, 4, 5}).
		Map(func(x int) int { return x * x }).
		Collect()
	expected := []int{1, 4, 9, 16, 25}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}

func Test_mapWhile(t *testing.T) {
	vals := FromSlice([]int{1, 2, 3, 4, 5}).
		MapWhile(func(val int) (int, bool) { return val * 2, val < 3 }).
		Collect()
	expected := []int{2, 4}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}

func Test_scan(t *testing.T) {
	vals := FromSlice([]int{1, 1, 1, 1, 1}).
		Scan(func(p, n int) int { return p + n }).
		Collect()
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}

func Test_skip(t *testing.T) {
	vals := FromSlice([]int{1, 2, 3, 4, 5}).Skip(2).Collect()
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}

func Test_SkipWhile(t *testing.T) {
	vals := FromSlice([]int{1, 1, 2, 3, 5, 8}).
		SkipWhile(func(x int) bool { return x&1 != 0 }).
		Collect()
	expected := []int{2, 3, 5, 8}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}

func Test_Take(t *testing.T) {
	vals := FromSlice([]int{1, 2, 3, 4, 5}).
		Take(3).
		Collect()
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}
