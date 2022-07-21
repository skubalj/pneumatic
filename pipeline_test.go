package pneumatic

import "testing"

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
			got := NewFromSlice(tt.input).All(tt.pred)
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
			got := NewFromSlice(tt.input).Any(tt.pred)
			if got != tt.expected {
				t.Errorf("Failed. Got %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCount(t *testing.T) {
	if got := NewFromSlice([]int{}).Count(); got != 0 {
		t.Errorf("Failed. Got %v, expected 0", got)
	}

	if got := NewFromSlice([]int{1, 2, 3}).Count(); got != 3 {
		t.Errorf("Failed. Got %v, expected 3", got)
	}

	if got := NewFromSlice([]int{1, 2, 3}).
		Filter(func(i int) bool { return i&1 == 0 }).
		Count(); got != 1 {
		t.Errorf("Failed. Got %v, expected 1", got)
	}
}

// func TestAny(t *testing.T) {
// 	type args struct {
// 		slice []T
// 	}
// 	tests := []struct {
// 		name string
// 		input []int
// 		pred func(int) bool
// 		expected []int
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			x := NewFromSlice(tt.input).Any(tt.pred)
// 			if got := NewFromSlice(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Failed. Got %v, expected %v", got, tt.expected)
// 			}
// 		})
// 	}
// }
