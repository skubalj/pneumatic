package pneumatic

import (
	"cmp"
	"iter"
	"slices"
)

// Read this iterator to its end and return the length
func Count[T any](seq iter.Seq[T]) int {
	count := 0
	for range seq {
		count++
	}
	return count
}

// Return the last element in this iterator
func Last[T any](seq iter.Seq[T]) (element T, ok bool) {
	var last T
	var notEmpty bool
	for x := range seq {
		last = x
		notEmpty = true
	}
	return last, notEmpty
}

// Fetch the `nth` element in the array.
//
// This operation uses zero-based indexing, so Nth(0) returns the first element
func Nth[T any](num int, seq iter.Seq[T]) (val T, ok bool) {
	next, stop := iter.Pull(Skip(num, seq))
	defer stop()
	return next()
}

// Merge the given set of iterators, one after another
func Chain[T any](sequences ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range sequences {
			for e := range seq {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// Convert an iter.Seq2 into an iter.Seq[pneumatic.Pair[T, U]]. While this type
// signature is undoubtedly more complicated, it allows you to use the resulting
// type with all the other pneumatic functions
func IngestSeq2[T, U any](seq iter.Seq2[T, U]) iter.Seq[Pair[T, U]] {
	return func(yield func(Pair[T, U]) bool) {
		for a, b := range seq {
			if !yield(Pair[T, U]{a, b}) {
				return
			}
		}
	}
}

func Zip[T, U any](seq1 iter.Seq[T], seq2 iter.Seq[U]) iter.Seq[Pair[T, U]] {
	return func(yield func(Pair[T, U]) bool) {
		next, stop := iter.Pull(seq2)
		defer stop()

		for value1 := range seq1 {
			value2, ok := next()
			if !ok {
				return
			}

			if !yield(Pair[T, U]{value1, value2}) {
				return
			}
		}
	}
}

func Intersperse[T any](separator T, seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		first, rest, ok := Peel(seq)
		if !ok || !yield(first) {
			return
		}

		for val := range rest {
			if !yield(separator) {
				return
			}
			if !yield(val) {
				return
			}
		}
	}
}

// Apply the given mapping function to convert one value into another
func Map[T, U any](fn func(T) U, seq iter.Seq[T]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for val := range seq {
			if !yield(fn(val)) {
				return
			}
		}
	}
}

// Retains only elements for which the predicate returns true
func Filter[T any](pred func(T) bool, seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range seq {
			if pred(val) && !yield(val) {
				return
			}
		}
	}
}

// A [Filter] and [Map] operation in one
func FilterMap[T, U any](fn func(T) (U, bool), seq iter.Seq[T]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for val := range seq {
			mappedVal, ok := fn(val)
			if ok && !yield(mappedVal) {
				return
			}
		}
	}
}

// Reverse [Filter]: discard all values for which the predicate returns true
func Discard[T any](pred func(T) bool, seq iter.Seq[T]) iter.Seq[T] {
	return Filter(func(t T) bool { return !pred(t) }, seq)
}

// Reverse [Filter]: discard all values for which the predicate returns true
func DiscardMap[T, U any](fn func(T) (U, bool), seq iter.Seq[T]) iter.Seq[U] {
	return FilterMap(func(t T) (U, bool) {
		u, ok := fn(t)
		return u, !ok
	}, seq)
}

// Return pairs containing the index of each element and the element itself.
func Enumerate[T any](seq iter.Seq[T]) iter.Seq[Pair[int, T]] {
	return func(yield func(Pair[int, T]) bool) {
		count := 0
		for val := range seq {
			if !yield(Pair[int, T]{count, val}) {
				return
			}
			count++
		}
	}
}

func SkipWhile[T any](pred func(T) bool, seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(seq)
		defer stop()

		for {
			val, ok := next()
			if !ok {
				return
			} else if !pred(val) {
				if !yield(val) {
					return
				}
				break
			}
		}

		for {
			val, ok := next()
			if !ok || !yield(val) {
				return
			}
		}
	}
}

func TakeWhile[T any](pred func(T) bool, seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range seq {
			if !pred(val) || !yield(val) {
				return
			}
		}
	}
}

// Skip forward by N values
func Skip[T any](num int, seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range seq {
			if num > 0 {
				num--
				continue
			}

			if !yield(val) {
				return
			}
		}
	}
}

// Stop after the first N values
func Take[T any](num int, seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range seq {
			if num <= 0 {
				return
			}
			if !yield(val) {
				return
			}
			num--
		}
	}
}

func Scan[T, S any](initialState S, fn func(accumulator S, next T) S, seq iter.Seq[T]) iter.Seq[S] {
	return func(yield func(S) bool) {
		acc := initialState
		for val := range seq {
			acc = fn(acc, val)
			if !yield(acc) {
				return
			}
		}
	}
}

func FlatMap[T, U any](fn func(T) iter.Seq[U], seq iter.Seq[T]) iter.Seq[U] {
	return Flatten(Map(fn, seq))
}

func Flatten[T any](seq iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := range seq {
			for val := range e {
				if !yield(val) {
					return
				}
			}
		}
	}
}

func Inspect[T any](fn func(T), seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range seq {
			fn(val)
			if !yield(val) {
				return
			}
		}
	}
}

func Partition[T any](pred func(T) bool, seq iter.Seq[T]) ([]T, []T) {
	var trueSet, falseSet []T
	for val := range seq {
		if pred(val) {
			trueSet = append(trueSet, val)
		} else {
			falseSet = append(falseSet, val)
		}
	}
	return trueSet, falseSet
}

// Fold with an accumulator function that can return an error
func TryFold[T, S any](initial S, fn func(accumulator S, next T) (S, error), seq iter.Seq[T]) (S, error) {
	var err error
	acc := initial
	for val := range seq {
		acc, err = fn(acc, val)
		if err != nil {
			return acc, err
		}
	}
	return acc, nil
}

func Fold[T, S any](initial S, fn func(accumulator S, next T) S, seq iter.Seq[T]) S {
	acc := initial
	for val := range seq {
		acc = fn(acc, val)
	}
	return acc
}

// Like [Fold], but the first value in the iterator is used as the inital value
func Reduce[T any](fn func(accumulator, next T) T, seq iter.Seq[T]) (value T, ok bool) {
	first, rest, ok := Peel(seq)
	if !ok {
		return first, false
	}

	return Fold(first, fn, rest), true
}

// Like [Reduce], but with an accumulator that returns an error
//
// Returns ok = false if the iterator is empty. This function will stop and
// return the first error encountered.
func TryReduce[T any](fn func(accumulator, next T) (T, error), seq iter.Seq[T]) (value T, err error, ok bool) {
	first, rest, ok := Peel(seq)
	if !ok {
		return first, nil, false
	}

	value, err = TryFold(first, fn, rest)
	return value, err, true
}

func All[T any](predicate func(T) bool, seq iter.Seq[T]) bool {
	for val := range seq {
		if !predicate(val) {
			return false
		}
	}
	return true
}

func Any[T any](predicate func(T) bool, seq iter.Seq[T]) bool {
	for val := range seq {
		if predicate(val) {
			return true
		}
	}
	return false
}

// Return the first element for which the predicate returns true
func Find[T any](predicate func(T) bool, seq iter.Seq[T]) (val T, ok bool) {
	for val := range seq {
		if predicate(val) {
			return val, true
		}
	}
	return
}

// Return the first index at which the given predicate returns true.
func Position[T any](predicate func(T) bool, seq iter.Seq[T]) (idx int, ok bool) {
	for val := range Enumerate(seq) {
		if predicate(val.Y) {
			return val.X, true
		}
	}
	return 0, false
}

func Rposition[T any](predicate func(T) bool, seq iter.Seq[T]) (idx int, ok bool) {
	p, ok := Find(func(p Pair[int, T]) bool { return predicate(p.Y) }, Rev(Enumerate(seq)))
	return p.X, ok
}

func Max[T cmp.Ordered](seq iter.Seq[T]) (val T, ok bool) {
	return MaxBy(cmp.Compare, seq)
}

func MaxBy[T any](cmp func(a, b T) int, seq iter.Seq[T]) (val T, ok bool) {
	return Reduce(func(a, b T) T {
		if cmp(a, b) < 0 {
			return a
		} else {
			return b
		}
	}, seq)
}

func Min[T cmp.Ordered](seq iter.Seq[T]) (val T, ok bool) {
	return MinBy(cmp.Compare, seq)
}

func MinBy[T any](cmp func(a, b T) int, seq iter.Seq[T]) (val T, ok bool) {
	return Reduce(func(a, b T) T {
		if cmp(a, b) < 0 {
			return a
		} else {
			return b
		}
	}, seq)
}

// Reverse the iterator
func Rev[T any](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		s := slices.Collect(seq)
		for i := len(s); i > 0; i-- {
			if !yield(s[i-1]) {
				return
			}
		}
	}
}

// func Unzip[T, U any](seq iter.Seq[Pair[T, U]]) {}

// Repeat the given iterator forever
func Cycle[T any](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			for val := range seq {
				if !yield(val) {
					return
				}
			}
		}
	}
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Number interface {
	Integer | ~float32 | ~float64 | ~complex64 | ~complex128
}

func Sum[T Number](seq iter.Seq[T]) T {
	acc, _ := Reduce(func(acc, next T) T { return acc + next }, seq)
	return acc
}

func Product[T Number](seq iter.Seq[T]) T {
	acc, _ := Reduce(func(acc, next T) T { return acc * next }, seq)
	return acc
}

// Return an iterator that produces the cartesian product of the two iterators.
//
// Both iterators must be repeatable
func CartesianProduct[T, U any](seq1 iter.Seq[T], seq2 iter.Seq[U]) iter.Seq[Pair[T, U]] {
	return func(yield func(Pair[T, U]) bool) {
		for v1 := range seq1 {
			for v2 := range seq2 {
				if !yield(Pair[T, U]{v1, v2}) {
					return
				}
			}
		}
	}
}

func IsSorted[T cmp.Ordered](seq iter.Seq[T]) bool {
	return IsSortedBy(cmp.Compare, seq)
}

// Return true if the values from the iterator are monotonically increasing
func IsSortedBy[T any](fn func(a, b T) int, seq iter.Seq[T]) bool {
	previous, rest, ok := Peel(seq)
	if !ok {
		return true
	}

	for val := range rest {
		if fn(previous, val) > 0 {
			return false
		}
		previous = val
	}
	return true
}

// Return non-overlapping slices of the given size. The final slice may not be
// of the specified size
func Chunks[T any](size int, seq iter.Seq[T]) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		chunk := make([]T, 0, size)
		for val := range seq {
			chunk = append(chunk, val)
			if len(chunk) == size {
				if !yield(chunk) {
					return
				}
				chunk = make([]T, 0, size)
			}
		}

		if len(chunk) > 0 {
			if !yield(chunk) {
				return
			}
		}
	}
}

// Return non-overlapping slices of the given size. All slices will be of the
// proper size, but some elements may be dropped if the final slice would not
// be full
func ChunksExact[T any](size int, seq iter.Seq[T]) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		chunk := make([]T, 0, size)
		for val := range seq {
			chunk = append(chunk, val)
			if len(chunk) == size {
				if !yield(chunk) {
					return
				}
				chunk = make([]T, 0, size)
			}
		}
	}
}

// Returns an iterator of slices made from a "window" sliding over the elements
// in the iterator. If there are not enough elements to make a window of the
// desired size, then the single window returned will be shorter than specified.
func Windows[T any](size int, seq iter.Seq[T]) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		next, stop := iter.Pull(seq)
		defer stop()

		window := make([]T, 0, size)
		for range size {
			val, ok := next()
			if !ok {
				break
			}
			window = append(window, val)
		}

		if !yield(window) {
			return
		}

		for {
			value, ok := next()
			if !ok {
				return
			}

			newWindow := make([]T, 0, size)
			newWindow = append(newWindow, window[1:]...)
			newWindow = append(newWindow, value)
			if !yield(window) {
				return
			}

			window = newWindow
		}
	}
}

// Returns an iterator of slices made from a "window" sliding over the elements
// in the iterator. If there are not enough elements to make a window of the
// desired size, then the iterator will yield no elements.
func WindowsExact[T any](size int, seq iter.Seq[T]) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		next, stop := iter.Pull(seq)
		defer stop()

		window := make([]T, 0, size)
		for range size {
			val, ok := next()
			if !ok {
				return
			}
			window = append(window, val)
		}

		if !yield(window) {
			return
		}

		for {
			value, ok := next()
			if !ok {
				return
			}

			newWindow := make([]T, 0, size)
			newWindow = append(newWindow, window[1:]...)
			newWindow = append(newWindow, value)
			if !yield(window) {
				return
			}

			window = newWindow
		}
	}
}

// Remove the first element, and return the rest as an iterator
func Peel[T any](seq iter.Seq[T]) (first T, rest iter.Seq[T], ok bool) {
	next, stop := iter.Pull(seq)
	first, ok = next()
	if !ok {
		stop()
		return first, Empty[T](), false
	}

	rest = func(yield func(T) bool) {
		defer stop()
		for {
			val, ok := next()
			if !ok || !yield(val) {
				return
			}
		}
	}
	return
}
