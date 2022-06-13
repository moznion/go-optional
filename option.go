package optional

import (
	"errors"
)

var (
	// ErrNoneValueTaken represents the error that is raised when None value is taken.
	ErrNoneValueTaken = errors.New("none value taken")
)

// Option is a data type that must be Some (i.e. having a value) or None (i.e. doesn't have a value).
type Option[T any] struct {
	value  T
	exists *struct{}
}

// Some is a function to make a Option type instance with the actual value.
func Some[T any](value T) Option[T] {
	return Option[T]{
		value:  value,
		exists: &struct{}{},
	}
}

// None is a function to make a Option type that doesn't have a value.
func None[T any]() Option[T] {
	return Option[T]{}
}

// IsNone returns whether the Option *doesn't* have a value or not.
func (o Option[T]) IsNone() bool {
	return o.exists == nil
}

// IsSome returns whether the Option has a value or not.
func (o Option[T]) IsSome() bool {
	return o.exists != nil
}

// Unwrap returns the value regardless of Some/None status.
// If the Option value is Some, this method returns the actual value.
// On the other hand, if the Option value is None, this method returns the *default* value according to the type.
func (o Option[T]) Unwrap() T {
	return o.value
}

// Take takes the contained value in Option.
// If Option value is Some, this returns the value that is contained in Option.
// On the other hand, this returns an ErrNoneValueTaken as the second return value.
func (o Option[T]) Take() (T, error) {
	if o.IsNone() {
		return o.value, ErrNoneValueTaken
		//     ~~~~~~~ uninitialized default value
	}
	return o.value, nil
}

// TakeOr returns the actual value if the Option has a value.
// On the other hand, this returns fallbackValue.
func (o Option[T]) TakeOr(fallbackValue T) T {
	if o.IsNone() {
		return fallbackValue
	}
	return o.value
}

// TakeOrElse returns the actual value if the Option has a value.
// On the other hand, this executes fallbackFunc and returns the result value of that function.
func (o Option[T]) TakeOrElse(fallbackFunc func() T) T {
	if o.IsNone() {
		return fallbackFunc()
	}
	return o.value
}

// Filter returns self if the Option has a value and the value matches the condition of the predicate function.
// In other cases (i.e. it doesn't match with the predicate or the Option is None), this returns None value.
func (o Option[T]) Filter(predicate func(v T) bool) Option[T] {
	if o.IsNone() {
		return None[T]()
	}

	if predicate(o.value) {
		return o
	}
	return None[T]()
}

// IfSome calls given function with the value of Option if the receiver value is Some.
func (o Option[T]) IfSome(f func(v T)) {
	if o.IsNone() {
		return
	}
	f(o.value)
}

// IfSomeWithError calls given function with the value of Option if the receiver value is Some.
// This method propagates the error of given function, and if the receiver value is None, this returns nil error.
func (o Option[T]) IfSomeWithError(f func(v T) error) error {
	if o.IsNone() {
		return nil
	}
	return f(o.value)
}

// IfNone calls given function if the receiver value is None.
func (o Option[T]) IfNone(f func()) {
	if o.IsSome() {
		return
	}
	f()
}

// IfNoneWithError calls given function if the receiver value is None.
// This method propagates the error of given function, and if the receiver value is Some, this returns nil error.
func (o Option[T]) IfNoneWithError(f func() error) error {
	if o.IsSome() {
		return nil
	}
	return f()
}

// Map converts given Option value to another Option value according to the mapper function.
// If given Option value is None, this also returns None.
func Map[T, U any](option Option[T], mapper func(v T) U) Option[U] {
	if option.IsNone() {
		return None[U]()
	}

	return Some(mapper(option.value))
}

// MapOr converts given Option value to another *actual* value according to the mapper function.
// If given Option value is None, this returns fallbackValue.
func MapOr[T, U any](option Option[T], fallbackValue U, mapper func(v T) U) U {
	if option.IsNone() {
		return fallbackValue
	}
	return mapper(option.value)
}

// MapWithError converts given Option value to another Option value according to the mapper function that has the ability to return the value with an error.
// If given Option value is None, this returns (None, nil). Else if the mapper returns an error then this returns (None, error).
// Unless of them, i.e. given Option value is Some and the mapper doesn't return the error, this returns (Some[U], nil).
func MapWithError[T, U any](option Option[T], mapper func(v T) (U, error)) (Option[U], error) {
	if option.IsNone() {
		return None[U](), nil
	}

	u, err := mapper(option.value)
	if err != nil {
		return None[U](), err
	}
	return Some(u), nil
}

// MapOrWithError converts given Option value to another *actual* value according to the mapper function that has the ability to return the value with an error.
// If given Option value is None, this returns (fallbackValue, nil). Else if the mapper returns an error then returns (_, error).
// Unless of them, i.e. given Option value is Some and the mapper doesn't return the error, this returns (U, nil).
func MapOrWithError[T, U any](option Option[T], fallbackValue U, mapper func(v T) (U, error)) (U, error) {
	if option.IsNone() {
		return fallbackValue, nil
	}
	return mapper(option.value)
}

// FlatMap converts give Option value to another Option value according to the mapper function.
// The difference from the Map is the mapper function returns an Option value instead of the bare value.
// If given Option value is None, this also returns None.
func FlatMap[T, U any](option Option[T], mapper func(v T) Option[U]) Option[U] {
	if option.IsNone() {
		return None[U]()
	}

	return mapper(option.value)
}

// FlatMapOr converts given Option value to another *actual* value according to the mapper function.
// The difference from the MapOr is the mapper function returns an Option value instead of the bare value.
// If given Option value is None or mapper function returns None, this returns fallbackValue.
func FlatMapOr[T, U any](option Option[T], fallbackValue U, mapper func(v T) Option[U]) U {
	if option.IsNone() {
		return fallbackValue
	}

	return (mapper(option.value)).TakeOr(fallbackValue)
}

// FlatMapWithError converts given Option value to another Option value according to the mapper function that has the ability to return the value with an error.
// The difference from the MapWithError is the mapper function returns an Option value instead of the bare value.
// If given Option value is None, this returns (None, nil). Else if the mapper returns an error then this returns (None, error).
// Unless of them, i.e. given Option value is Some and the mapper doesn't return the error, this returns (Some[U], nil).
func FlatMapWithError[T, U any](option Option[T], mapper func(v T) (Option[U], error)) (Option[U], error) {
	if option.IsNone() {
		return None[U](), nil
	}

	mapped, err := mapper(option.value)
	if err != nil {
		return None[U](), err
	}
	return mapped, nil
}

// FlatMapOrWithError converts given Option value to another *actual* value according to the mapper function that has the ability to return the value with an error.
// The difference from the MapOrWithError is the mapper function returns an Option value instead of the bare value.
// If given Option value is None, this returns (fallbackValue, nil). Else if the mapper returns an error then returns ($zero_value_of_type, error).
// Unless of them, i.e. given Option value is Some and the mapper doesn't return the error, this returns (U, nil).
func FlatMapOrWithError[T, U any](option Option[T], fallbackValue U, mapper func(v T) (Option[U], error)) (U, error) {
	if option.IsNone() {
		return fallbackValue, nil
	}

	maybe, err := mapper(option.value)
	if err != nil {
		var zeroValue U
		return zeroValue, err
	}

	return maybe.TakeOr(fallbackValue), nil
}

// Pair is a data type that represents a tuple that has two elements.
type Pair[T, U any] struct {
	Value1 T
	Value2 U
}

// Zip zips two Options into a Pair that has each Option's value.
// If either one of the Options is None, this also returns None.
func Zip[T, U any](opt1 Option[T], opt2 Option[U]) Option[Pair[T, U]] {
	if opt1.IsSome() && opt2.IsSome() {
		return Some(Pair[T, U]{
			Value1: opt1.value,
			Value2: opt2.value,
		})
	}

	return None[Pair[T, U]]()
}

// ZipWith zips two Options into a typed value according to the zipper function.
// If either one of the Options is None, this also returns None.
func ZipWith[T, U, V any](opt1 Option[T], opt2 Option[U], zipper func(opt1 T, opt2 U) V) Option[V] {
	if opt1.IsSome() && opt2.IsSome() {
		return Some(zipper(opt1.value, opt2.value))
	}
	return None[V]()
}

// Unzip extracts the values from a Pair and pack them into each Option value.
// If the given zipped value is None, this returns None for all return values.
func Unzip[T, U any](zipped Option[Pair[T, U]]) (Option[T], Option[U]) {
	if zipped.IsNone() {
		return None[T](), None[U]()
	}

	pair := zipped.value
	return Some(pair.Value1), Some(pair.Value2)
}

// UnzipWith extracts the values from the given value according to the unzipper function and pack the into each Option value.
// If the given zipped value is None, this returns None for all return values.
func UnzipWith[T, U, V any](zipped Option[V], unzipper func(zipped V) (T, U)) (Option[T], Option[U]) {
	if zipped.IsNone() {
		return None[T](), None[U]()
	}

	v1, v2 := unzipper(zipped.value)
	return Some(v1), Some(v2)
}
