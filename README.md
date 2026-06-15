# `pneumatic`: Go Iterator Extensions

Functional programming is one of my favorite paradigms. Combining functions to
build up lazily evaluated iterators can be more readable than writing literal
for-loops. You can focus on the "flow" of data through a iterator "pipeline",
rather than worrying about individual loop values.

Traditionally, Go has not been a good language for functional programming.
While Go has supported passing functions as first-class values, the lack of
generics has made it difficult to make ergnonmic typed functions. In Go 1.18,
a generics system was bolted onto the Go type system, allowing us to build more
flexible interfaces.

Further, in version 1.23, Go added new functional programming features with
the `iter` module. However, these iterators are very verbose to write, and the
standard library is missing helper functions for a lot of everyday cases. This
library is designed to be a "functional programming standard library" of
functions common in other languages. Much inspiration was drawn from Rust's
`std::iterator::Iterator` trait.

## Examples

This pair of recursive functions used for generating prime numbers
demonstrates some of the basic building blocks, such as `Map`, `Filter`, and
`Any`.

```go
import (
	"iter"

	pn "github.com/skubalj/pneumatic"
)

func primesUnder(bound int) iter.Seq[int] {
	return pn.Filter(isPrime, pn.Range(2, bound))
}

func isPrime(candidate int) bool {
	return !pn.Any(func(x int) bool { return x == 0 },
		pn.Map(func(prime int) int { return candidate % prime },
			primesUnder(candidate)))
}
```

## License

This project is released under the terms of the MIT license.

Copyright (c) 2026 Joseph Skubal
