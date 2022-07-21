# `pneumatic`

`pneumatic` is a library providing lazily evaluated functional programing to Go. The API is meant 
to be quite similar to Rust's [`std::iter::Iterator`](https://doc.rust-lang.org/std/iter/trait.Iterator.html)
trait. Working in Go professionally, one of the things I find irksome is a lack of declarative 
programming. The primary goal for this library is simply to be an exercise in generics and 
functional programming, making my (and hopefully your) Go experience better.

Note that while you and I both love using declarative, functional code, it is generally considered
unidiomatic in Go. Consider your project's style guide before deciding whether `pneumatic` is right
for you.

```go
import "github.com/skubalj/pneumatic"

evenSquares := pneumatic.NewFromSlice([]int{1, 2, 3, 4, 5}).
    Map(func(x int) int {return x * x}).
    Filter(func(y int) bool {return x&1 == 0}).
    Collect() // => []int{4, 16}
```

## Rationale
TODO

## Performance
TODO

## Minimum Supported Go Version
As this library uses generics, the MSGV is the **1.18** series

## License
`pneumatic` is released under the terms of the Mozilla Public License v. 2.0. This means that while
`pneumatic` may be used as a component of closed-source software and can be statically compiled 
into applications or larger libraries, changes to the pneumatic source code must be released under
the terms of the MPLv2.0.

See the `LICENSE` file for more information.  
Copyright (C) Joseph Skubal 2022
