# go-optional [![.github/workflows/check.yml](https://github.com/moznion/go-optional/actions/workflows/check.yml/badge.svg)](https://github.com/moznion/go-optional/actions/workflows/check.yml) [![codecov](https://codecov.io/gh/moznion/go-optional/branch/main/graph/badge.svg?token=0HCVy6COy4)](https://codecov.io/gh/moznion/go-optional)

A library that provides [Go Generics](https://go.dev/blog/generics-proposal) friendly "optional" features.

## Synopsis

```go
some := Some[int](123)
fmt.Printf("%v\n", some.IsSome()) # => true
fmt.Printf("%v\n", some.IsNone()) # => false

v, err := some.Take()
fmt.Printf("err is nil: %v\n", err == nil) # => err is nil: true
fmt.Printf("%d\n", v) # => 123

mapped := optional.Map(some, func (v int) int {
    return v * 2
})
fmt.Printf("%v\n", mapped.IsSome()) # => true

mappedValue, _ := some.Take()
fmt.Printf("%d\n", mappedValue) # => 246
```

## Docs

[![GoDoc](https://godoc.org/github.com/moznion/go-optional?status.svg)](https://godoc.org/github.com/moznion/go-optional)

and examples are [here](./option_test.go).

## Current Status

Currently (at the moment: Nov 18, 2021), go 1.18 has not been released yet, so if you'd like to try this, please use the tip runtime.  
Of course, the new runtime version hasn't been released yet so this library has the possibility to change the implementation as well.

## Known Issues

The runtime raises a compile error like "methods cannot have type parameters", so `Map()`, `MapOr()`, `Zip()` and `ZipWith()` has been providing as functions. Basically, it would be better to provide them as the methods, but currently, it compromises with the limitation.

## Author

moznion (<moznion@mail.moznion.net>)

