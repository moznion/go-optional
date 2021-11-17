package optional

import "errors"

var (
	// ErrNoneValueTaken represents the error that is raised when None value is taken.
	ErrNoneValueTaken = errors.New("none value taken")
)

// Option is a data type that must be Some (i.e. having a value) or None (i.e. doesn't have a value).
type Option[T any] struct {
	value T
	exists *struct{}
}

// Some is a function to make a Option type instance with the actual value.
func Some[T any](value T) Option[T] {
	return Option[T]{
		value: value,
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
func (o Option[T]) TakeOrElse(fallbackFunc func () T) T {
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

// Map converts the Option value to another value according to the mapper function.
// If Option value is None, this also returns None.
func Map[T, U any](option Option[T], mapper func(v T) U) Option[U] {
	if option.IsNone() {
		return None[U]()
	}

	return Some[U](mapper(option.value))
}

// MapOr converts t	he Option value to another value according to the mapper function.
// If Option value is None, this returns fallbackValue.
func MapOr[T, U any](option Option[T], fallbackValue U, mapper func(v T) U) U {
	if option.IsNone() {
		return fallbackValue
	}
	return mapper(option.value)
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
		return Some[Pair[T, U]](Pair[T, U]{
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
		return Some[V](zipper(opt1.value, opt2.value))
	}
	return None[V]()
}
