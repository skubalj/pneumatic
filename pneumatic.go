package pneumatic

import (
	"github.com/skubalj/pneumatic/pipeline"
)

type fromSlice[T any] struct {
	idx    int
	source []T
}

// Create a new pipeline that iterates through the given slice
func FromSlice[T any](slice []T) pipeline.PipelineStep[T, T] {
	return pipeline.StepFrom[T, T](&fromSlice[T]{0, slice})
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

////////////////////////////////////////////////////////////////////////////////////////////////////

type fromChan[T any] struct {
	ch chan T
}

// Create a new pipeline from the provided channel.
//
// This pipeline will stay open until the given channel is closed by the sender.
func FromChan[T any](ch chan T) pipeline.PipelineStep[T, T] {
	return pipeline.StepFrom[T, T](&fromChan[T]{ch})
}

func (fc *fromChan[T]) Next() (T, bool) {
	val, ok := <-fc.ch
	return val, ok
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type empty[T any] struct{}

// Create a new empty pipeline that yields no elements.
//
// This can be used with `chain` to build up a pipeline from many smaller pipelines
func Empty[T any]() pipeline.PipelineStep[T, T] {
	return pipeline.StepFrom[T, T](new(empty[T]))
}

func (fs *empty[T]) Next() (T, bool) {
	return *new(T), false
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type once[T any] struct {
	val    T
	called bool
}

// Create a new pipeline that yields the provided element exactly once.
//
//
func Once[T any](val T) pipeline.PipelineStep[T, T] {
	return pipeline.StepFrom[T, T](&once[T]{val, false})
}

func (o *once[T]) Next() (T, bool) {
	if o.called {
		return *new(T), false
	} else {
		o.called = true
		return o.val, true
	}
}
