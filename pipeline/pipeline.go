package pipeline

// The underlying trait allowing lazy iteration by calling the next value
type Pipeline[T any] interface {
	Next() (T, bool)
}

// A pipeline step represents a single segment of the pipeline. Additional segments can be appended
// to produce a longer, more complicated pipeline.
type PipelineStep[T any, U any] struct {
	prev Pipeline[T]
}

// Constructor to create a step from an arbitrary type implementing Pipeline
func StepFrom[T, U any](p Pipeline[T]) PipelineStep[T, U] {
	return PipelineStep[T, U]{p}
}

// Get the next value from this pipeline step
func (p *PipelineStep[T, U]) Next() (T, bool) {
	return p.prev.Next()
}

///////////////////////////////////////////////////////////////////////////////

// Returns true if for all elements in this pipeline, the predicate evaluates to `true`
//
// This function will short-circuit upon finding any element for which the predicate is false. On
// an empty array, this will always evaluate to true.
func (p PipelineStep[T, U]) All(pred func(T) bool) bool {
	for {
		val, ok := p.prev.Next()
		if !ok {
			return true
		} else if !pred(val) {
			return false
		}
	}
}

// Returns true if for any element in this pipeline, the predicate evaluates to `true`
//
// This function will short-circuit upon finding any element for which the predicate is true. On
// an empty array, this will always evaluate to false.
func (p PipelineStep[T, U]) Any(pred func(T) bool) bool {
	for {
		val, ok := p.prev.Next()
		if !ok {
			return false
		} else if pred(val) {
			return true
		}
	}
}

// Collect all of the elements into a slice
func (p PipelineStep[T, U]) Collect() []T {
	arr := []T{}
	for {
		val, ok := p.prev.Next()
		if !ok {
			break
		}

		arr = append(arr, val)
	}

	return arr
}

// Count the number of elements in this iterator
func (p PipelineStep[T, U]) Count() int {
	count := 0
	for {
		if _, ok := p.prev.Next(); ok {
			count++
		} else {
			return count
		}
	}
}

// Apply the provided function for each value in this pipeline
func (p PipelineStep[T, U]) ForEach(fn func(T)) {
	for {
		val, ok := p.prev.Next()
		if ok {
			fn(val)
		} else {
			break
		}
	}
}

// func (p PipelineStep[T, U]) Reduce(fn func(prev, next T) T) (T, bool) {
// 	value, ok := p.prev.Next()
// 	if !ok {
// 		return *new(T), false
// 	}

// 	for {
// 		newVal, ok := p.prev.Next()
// 		if ok {
// 			value = fn(value, newVal)
// 		} else {
// 			return value, ok
// 		}
// 	}
// }
