package pneumatic

import (
	"reflect"
	"testing"
)

func TestAll(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		pred     func(int) bool
		expected bool
	}{
		{"empty", []int{}, func(i int) bool { return i&1 == 0 }, true},
		{"none", []int{1, 3, 5}, func(i int) bool { return i&1 == 0 }, false},
		{"some", []int{1, 2, 4}, func(i int) bool { return i&1 == 0 }, false},
		{"all", []int{2, 4, 8}, func(i int) bool { return i&1 == 0 }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromSlice(tt.input).All(tt.pred)
			if got != tt.expected {
				t.Errorf("Failed. Got %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		pred     func(int) bool
		expected bool
	}{
		{"empty", []int{}, func(i int) bool { return i&1 == 0 }, false},
		{"none", []int{1, 3, 5}, func(i int) bool { return i&1 == 0 }, false},
		{"some", []int{1, 2, 4}, func(i int) bool { return i&1 == 0 }, true},
		{"all", []int{2, 4, 8}, func(i int) bool { return i&1 == 0 }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromSlice(tt.input).Any(tt.pred)
			if got != tt.expected {
				t.Errorf("Failed. Got %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCount(t *testing.T) {
	if got := FromSlice([]int{}).Count(); got != 0 {
		t.Errorf("Failed. Got %v, expected 0", got)
	}

	if got := FromSlice([]int{1, 2, 3}).Count(); got != 3 {
		t.Errorf("Failed. Got %v, expected 3", got)
	}

	if got := FromSlice([]int{1, 2, 3}).
		Filter(func(i int) bool { return i&1 == 0 }).
		Count(); got != 1 {
		t.Errorf("Failed. Got %v, expected 1", got)
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		pred     func(int) bool
		expected int
		valid    bool
	}{
		{"empty", []int{}, func(i int) bool { return i&1 == 0 }, 0, false},
		{"none", []int{1, 3, 5}, func(i int) bool { return i&1 == 0 }, 0, false},
		{"one", []int{1, 2, 3}, func(i int) bool { return i&1 == 0 }, 2, true},
		{"many", []int{2, 4, 8}, func(i int) bool { return i&1 == 0 }, 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := FromSlice(tt.input).Find(tt.pred)
			if got != tt.expected || ok != tt.valid {
				t.Errorf("Failed. Got %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestNth(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		idx      int
		expected int
		valid    bool
	}{
		{"empty", []int{}, 2, 0, false},
		{"first", []int{1, 3, 5}, 0, 1, true},
		{"middle", []int{1, 2, 3}, 1, 2, true},
		{"last", []int{1, 2, 3}, 2, 3, true},
		{"out of range", []int{2, 4, 8}, 4, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := FromSlice(tt.input).Nth(tt.idx)
			if got != tt.expected || ok != tt.valid {
				t.Errorf("Failed. Got %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		group1 []int
		group2 []int
	}{
		{"empty", []int{}, []int{}, []int{}},
		{"all true", []int{2, 4, 6, 8}, []int{2, 4, 6, 8}, []int{}},
		{"all false", []int{1, 3, 5, 7}, []int{}, []int{1, 3, 5, 7}},
		{"mixed", []int{1, 2, 3, 4, 5}, []int{2, 4}, []int{1, 3, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tVals, fVals := FromSlice(tt.input).Partition(func(i int) bool { return i&1 == 0 })
			if !reflect.DeepEqual(tt.group1, tVals) || !reflect.DeepEqual(tt.group2, fVals) {
				t.Errorf("Failed. Got (%v, %v), expected (%v, %v)", tVals, fVals, tt.group1, tt.group2)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		fn     func(int, int) int
		output int
		valid  bool
	}{
		{"empty", []int{}, func(a, b int) int { return a + b }, 0, false},
		{"sum", []int{1, 1, 1, 1}, func(a, b int) int { return a + b }, 4, true},
		{"difference", []int{4, 1, 1, 1}, func(a, b int) int { return a - b }, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := FromSlice(tt.input).Reduce(tt.fn)
			if got != tt.output || tt.valid != ok {
				t.Errorf("Failed. Got (%v, %v), expected (%v, %v)", got, ok, tt.output, tt.valid)
			}
		})
	}
}
