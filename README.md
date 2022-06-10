# go-optional [![.github/workflows/check.yml](https://github.com/moznion/go-optional/actions/workflows/check.yml/badge.svg)](https://github.com/moznion/go-optional/actions/workflows/check.yml) [![codecov](https://codecov.io/gh/moznion/go-optional/branch/main/graph/badge.svg?token=0HCVy6COy4)](https://codecov.io/gh/moznion/go-optional) [![GoDoc](https://godoc.org/github.com/moznion/go-optional?status.svg)](https://godoc.org/github.com/moznion/go-optional)

A library that provides [Go Generics](https://go.dev/blog/generics-proposal) friendly "optional" features.

## Synopsis

```go
some := optional.Some[int](123)
fmt.Printf("%v\n", some.IsSome()) // => true
fmt.Printf("%v\n", some.IsNone()) // => false

v, err := some.Take()
fmt.Printf("err is nil: %v\n", err == nil) // => err is nil: true
fmt.Printf("%d\n", v) // => 123

mapped := optional.Map(some, func (v int) int {
    return v * 2
})
fmt.Printf("%v\n", mapped.IsSome()) // => true

mappedValue, _ := some.Take()
fmt.Printf("%d\n", mappedValue) // => 246
```

```go
none := optional.None[int]()
fmt.Printf("%v\n", none.IsSome()) // => false
fmt.Printf("%v\n", none.IsNone()) // => true

_, err := none.Take()
fmt.Printf("err is nil: %v\n", err == nil) // => err is nil: false
// the error must be `ErrNoneValueTaken`

mapped := optional.Map(none, func (v int) int {
    return v * 2
})
fmt.Printf("%v\n", mapped.IsNone()) // => true
```

and more detailed examples are here: [./examples_test.go](./examples_test.go).

## Docs

[![GoDoc](https://godoc.org/github.com/moznion/go-optional?status.svg)](https://godoc.org/github.com/moznion/go-optional)

## Tips

- it would be better to deal with an Option value as a non-pointer because if the Option value can accept nil it becomes worthless

## Known Issues

The runtime raises a compile error like "methods cannot have type parameters", so `Map()`, `MapOr()`, `MapWithError()`, `MapOrWithError()`, `Zip()`, `ZipWith()`, `Unzip()` and `UnzipWith()` have been providing as functions. Basically, it would be better to provide them as the methods, but currently, it compromises with the limitation.

## Author

moznion (<moznion@mail.moznion.net>)

