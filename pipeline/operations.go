// This file contains the set of non-terminal operations that can be conducted on pipelines

package pipeline

type chain[T any] struct {
	PipelineStep[T, T]
	other PipelineStep[T, T]
}

func (p PipelineStep[T, U]) Chain(other PipelineStep[T, U]) PipelineStep[T, T] {
	return PipelineStep[T, T]{&chain[T]{PipelineStep[T, T](p), PipelineStep[T, T](other)}}
}

func (c *chain[T]) Next() (T, bool) {
	if n, ok := c.prev.Next(); ok {
		return n, true
	} else if n, ok := c.other.prev.Next(); ok {
		return n, true
	} else {
		return *new(T), false
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////

type mapOperation[T, U any] struct {
	PipelineStep[T, U]
	fn func(T) U
}

func (p PipelineStep[T, U]) Map(fn func(T) U) PipelineStep[U, U] {
	return PipelineStep[U, U]{&mapOperation[T, U]{p, fn}}
}

func (m *mapOperation[T, U]) Next() (U, bool) {
	n, ok := m.prev.Next()
	if ok {
		return m.fn(n), true
	} else {
		return *new(U), false
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type mapWhile[T, U any] struct {
	PipelineStep[T, U]
	fn   func(T) (U, error)
	done bool
}

func (p PipelineStep[T, U]) MapWhile(fn func(T) (U, error)) PipelineStep[U, U] {
	return PipelineStep[U, U]{&mapWhile[T, U]{p, fn, false}}
}

func (m *mapWhile[T, U]) Next() (U, bool) {
	n, ok := m.prev.Next()
	if ok && !m.done {
		val, err := m.fn(n)
		if err != nil {
			return *new(U), false
		} else {
			return val, true
		}
	} else {
		return *new(U), false
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type skip[T any] struct {
	PipelineStep[T, T]
	num int
}

func (p PipelineStep[T, U]) Skip(num int) PipelineStep[T, T] {
	return PipelineStep[T, T]{&skip[T]{PipelineStep[T, T](p), num}}
}

func (s *skip[T]) Next() (T, bool) {
	for s.num > 0 {
		s.num--
		if _, ok := s.prev.Next(); !ok {
			return *new(T), false
		}
	}
	return s.prev.Next()
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type skipWhile[T any] struct {
	PipelineStep[T, T]
	pred func(T) bool
	done bool
}

func (p PipelineStep[T, U]) SkipWhile(pred func(T) bool) PipelineStep[T, T] {
	return PipelineStep[T, T]{&skipWhile[T]{PipelineStep[T, T](p), pred, false}}
}

func (s *skipWhile[T]) Next() (T, bool) {
	if s.done {
		return s.prev.Next()
	}
	s.done = true

	for {
		val, ok := s.prev.Next()
		if !ok {
			return *new(T), false
		} else if !s.pred(val) {
			return val, true
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type take[T any] struct {
	PipelineStep[T, T]
	num int
}

// Shorten this pipeline to be at most `n` elements.
func (p PipelineStep[T, U]) Take(num int) PipelineStep[T, T] {
	return PipelineStep[T, T]{&take[T]{PipelineStep[T, T](p), num}}
}

func (t *take[T]) Next() (T, bool) {
	if t.num > 0 {
		t.num--
		return t.prev.Next()
	}
	return *new(T), false
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type takeWhile[T any] struct {
	PipelineStep[T, T]
	pred func(T) bool
}

// Take values from this pipeline until the provided predicate evaluates to false
func (p PipelineStep[T, U]) TakeWhile(pred func(T) bool) PipelineStep[T, T] {
	return PipelineStep[T, T]{&takeWhile[T]{PipelineStep[T, T](p), pred}}
}

func (t *takeWhile[T]) Next() (T, bool) {
	val, ok := t.prev.Next()
	if ok && t.pred(val) {
		return val, true
	} else {
		return *new(T), false
	}
}
