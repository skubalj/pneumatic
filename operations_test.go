package pneumatic

import (
	"reflect"
	"testing"
)

func Test_map(t *testing.T) {
	vals := NewFromSlice([]int{1, 2, 3, 4, 5}).
		Map(func(x int) int { return x * x }).
		Collect()
	expected := []int{1, 4, 9, 16, 25}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}

func Test_filter(t *testing.T) {
	vals := NewFromSlice([]int{1, 2, 3, 4, 5}).
		Filter(func(x int) bool { return x&1 == 0 }).
		Collect()
	expected := []int{2, 4}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}

func Test_skip(t *testing.T) {
	vals := NewFromSlice([]int{1, 2, 3, 4, 5}).Skip(2).Collect()
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Failed. Expected %v, got %v", expected, vals)
	}
}
