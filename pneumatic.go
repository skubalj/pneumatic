package pneumatic

type Pipeline[T any] interface {
	Next() (T, bool)
}

type PipelineStep[T any, U any] struct {
	prev Pipeline[T]
}

// Get the next value from this pipeline step
func (p *PipelineStep[T, U]) Next() (T, bool) {
	return p.prev.Next()
}

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

///////////////////////////////////////////////////////////////////////////////

type fromSlice[T any] struct {
	idx    int
	source []T
}

func NewFromSlice[T any](slice []T) PipelineStep[T, T] {
	return PipelineStep[T, T]{&fromSlice[T]{0, slice}}
}

func (fs *fromSlice[T]) Next() (T, bool) {
	if fs.idx >= len(fs.source) {
		return *new(T), false
	} else {
		val := fs.source[fs.idx]
		fs.idx++
		return val, true
	}
}

///////////////////////////////////////////////////////////////////////////////

type filter[T any] struct {
	PipelineStep[T, T]
	fn func(T) bool
}

func (p PipelineStep[T, U]) Filter(pred func(T) bool) PipelineStep[T, T] {
	return PipelineStep[T, T]{&filter[T]{PipelineStep[T, T](p), pred}}
}

func (f *filter[T]) Next() (T, bool) {
	for {
		n, ok := f.prev.Next()
		if !ok {
			return *new(T), false
		} else if f.fn(n) {
			return n, true
		}
	}
}

///////////////////////////////////////////////////////////////////////////////

type map_op[T, U any] struct {
	PipelineStep[T, U]
	fn func(T) U
}

func (p PipelineStep[T, U]) Map(fn func(T) U) PipelineStep[U, U] {
	return PipelineStep[U, U]{&map_op[T, U]{p, fn}}
}

func (m *map_op[T, U]) Next() (U, bool) {
	n, ok := m.prev.Next()
	if ok {
		return m.fn(n), true
	} else {
		return *new(U), false
	}
}
