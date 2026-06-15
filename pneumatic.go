package pneumatic

import (
	"iter"
)

// A pair of values, representing a 2-tuple
type Pair[X, Y any] struct {
	X X
	Y Y
}

// Convert the value and error into a pair value
func ZipPair[X, Y any](x X, y Y) Pair[X, Y] {
	return Pair[X, Y]{x, y}
}

func (p Pair[X, Y]) First() X {
	return p.X
}

func (p Pair[X, Y]) Last() Y {
	return p.Y
}

func (p Pair[X, Y]) Unzip() (X, Y) {
	return p.X, p.Y
}

// Create a new iterator from the provided channel.
//
// This iterator will stay open until the given channel is closed by the sender.
func FromChan[T any](ch chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range ch {
			if !yield(val) {
				return
			}
		}
	}
}

// An empty iterator that yields no elements
func Empty[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}

// An iterator that yields the provided value exactly once
func Once[T any](val T) iter.Seq[T] {
	return func(yield func(T) bool) { yield(val) }
}

// An iterator that repeats the provided element infinitely
func Repeated[T any](val T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			if !yield(val) {
				return
			}
		}
	}
}

// An iterator of all elements from lower (inclusive) to upper (exclusive)
func Range[T Integer](lower, upper T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := lower; i < upper; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Like [Range], but includes the upper bound in the iterator
func RangeInclusive[T Integer](lower, upper T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := lower; i <= upper; i++ {
			if !yield(i) {
				return
			}
		}
	}
}
