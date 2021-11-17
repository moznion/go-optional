package optional

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption_IsNone(t *testing.T) {
	assert.True(t, None[int]().IsNone())
	assert.False(t, Some[int](123).IsNone())
}

func TestOption_IsSome(t *testing.T) {
	assert.False(t, None[int]().IsSome())
	assert.True(t, Some[int](123).IsSome())
}

func TestOption_Take(t *testing.T) {
	v, err := Some[int](123).Take()
	assert.NoError(t, err)
	assert.Equal(t, 123, v)

	v, err = None[int]().Take()
	assert.ErrorIs(t, err, ErrNoneValueTaken)
	assert.Equal(t, 0, v)
}

func TestOption_TakeOr(t *testing.T) {
	v := Some[int](123).TakeOr(666)
	assert.Equal(t, 123, v)

	v = None[int]().TakeOr(666)
	assert.Equal(t, 666, v)
}

func TestOption_TakeOrElse(t *testing.T) {
	v := Some[int](123).TakeOrElse(func () int {
		return 666
	})
	assert.Equal(t, 123, v)

	v = None[int]().TakeOrElse(func () int {
		return 666
	})
	assert.Equal(t, 666, v)
}

func TestOption_Filter(t *testing.T) {
	isEven := func (v int) bool {
		return v % 2 == 0
	}

	o := Some[int](2).Filter(isEven)
	assert.True(t, o.IsSome())
	assert.Equal(t, 2, o.value)

	o = Some[int](1).Filter(isEven)
	assert.True(t, o.IsNone())

	o = None[int]().Filter(isEven)
	assert.True(t, o.IsNone())
}

func TestMap(t *testing.T) {
	some := Some[int](123)
	mapped := Map(some, func(v int) string {
		return fmt.Sprintf("%d", v)
	})
	taken, err := mapped.Take()
	assert.NoError(t, err)
	assert.Equal(t, "123", taken)

	none := None[int]()
	mapped = Map(none, func(v int) string {
		return fmt.Sprintf("%d", v)
	})
	assert.True(t, mapped.IsNone())
}

func TestMapOr(t *testing.T) {
	some := Some[int](123)
	mapped := MapOr(some, "666", func(v int) string {
		return fmt.Sprintf("%d", v)
	})
	assert.Equal(t, "123", mapped)

	none := None[int]()
	mapped = MapOr(none, "666", func(v int) string {
		return fmt.Sprintf("%d", v)
	})
	assert.Equal(t, "666", mapped)
}

func TestZip(t *testing.T) {
	some1 := Some[int](123)
	some2 := Some[string]("foo")
	none := None[uint]()

	zipped := Zip(some1, some2)
	assert.True(t, zipped.IsSome())
	assert.Equal(t, Pair[int, string]{
		Value1: 123,
		Value2: "foo",
	}, zipped.value)

	assert.True(t, Zip(none, some1).IsNone())
	assert.True(t, Zip(some1, none).IsNone())
}

func TestZipWith(t *testing.T) {
	type Data struct {
		A string
		B int
	}

	some1 := Some[int](123)
	some2 := Some[string]("foo")

	zipped := ZipWith(some1, some2, func (v1 int, v2 string) Data {
		return Data{
			A: v2,
			B: v1,
		}
	})
	assert.True(t, zipped.IsSome())
	assert.Equal(t, Data{
		A: "foo",
		B: 123,
	}, zipped.value)

	assert.True(t, ZipWith(None[int](), some1, func (v1, v2 int) Data {
		return Data{}
	}).IsNone())
	assert.True(t, ZipWith(some1, None[int](), func (v1, v2 int) Data {
		return Data{}
	}).IsNone())
}

func ExampleOption_IsNone() {
	some := Some[int](1)
	fmt.Printf("%v\n", some.IsNone())
	none := None[int]()
	fmt.Printf("%v\n", none.IsNone())
	// Output:
	// false
	// true
}

func ExampleOption_IsSome() {
	some := Some[int](1)
	fmt.Printf("%v\n", some.IsSome())
	none := None[int]()
	fmt.Printf("%v\n", none.IsSome())
	// Output:
	// true
	// false
}

func ExampleOption_Take() {
	some := Some[int](1)
	v, err := some.Take()
	fmt.Printf("%d\n", v)
	fmt.Printf("%v\n", err == nil)

	none := None[int]()
	_, err = none.Take()
	fmt.Printf("%v\n", err == nil)

	// Output:
	// 1
	// true
	// false
}

func ExampleOption_TakeOr() {
	some := Some[int](1)
	v := some.TakeOr(666)
	fmt.Printf("%d\n", v)

	none := None[int]()
	v = none.TakeOr(666)
	fmt.Printf("%d\n", v)

	// Output:
	// 1
	// 666
}

func ExampleOption_TakeOrElse() {
	some := Some[int](1)
	v := some.TakeOrElse(func () int {
		return 666
	})
	fmt.Printf("%d\n", v)

	none := None[int]()
	v = none.TakeOrElse(func () int {
		return 666
	})
	fmt.Printf("%d\n", v)

	// Output:
	// 1
	// 666
}

func ExampleOption_Filter() {
	isEven := func (v int) bool {
		return v % 2 == 0
	}

	some := Some[int](2)
	opt := some.Filter(isEven)
	fmt.Printf("%d\n", opt.TakeOr(0))

	some = Some[int](1)
	opt = some.Filter(isEven)
	fmt.Printf("%d\n", opt.TakeOr(0))

	none := None[int]()
	opt = none.Filter(isEven)
	fmt.Printf("%d\n", opt.TakeOr(0))

	// Output:
	// 2
	// 0
	// 0
}

func ExampleMap() {
	mapper := func (v int) string {
		return fmt.Sprintf("%d", v)
	}

	some := Some[int](1)
	opt := Map(some, mapper)
	fmt.Printf("%s\n", opt.TakeOr("N/A"))

	none := None[int]()
	opt = Map(none, mapper)
	fmt.Printf("%s\n", opt.TakeOr("N/A"))

	// Output:
	// 1
	// N/A
}

func ExampleMapOr() {
	mapper := func (v int) string {
		return fmt.Sprintf("%d", v)
	}

	some := Some[int](1)
	mapped := MapOr(some, "N/A", mapper)
	fmt.Printf("%s\n", mapped)

	none := None[int]()
	mapped = MapOr(none, "N/A", mapper)
	fmt.Printf("%s\n", mapped)

	// Output:
	// 1
	// N/A
}

func ExampleZip() {
	maybePair := Zip(Some[int](1), Some[string]("foo"))
	pair, err := maybePair.Take()
	fmt.Printf("is none => %v\n", maybePair.IsNone())
	fmt.Printf("err is nil => %v\n", err == nil)
	fmt.Printf("%d %s\n", pair.Value1, pair.Value2)

	maybePair = Zip(Some[int](1), None[string]())
	fmt.Printf("is none => %v\n", maybePair.IsNone())

	maybePair = Zip(None[int](), Some[string]("foo"))
	fmt.Printf("is none => %v\n", maybePair.IsNone())

	maybePair = Zip(None[int](), None[string]())
	fmt.Printf("is none => %v\n", maybePair.IsNone())

	// Output:
	// is none => false
	// err is nil => true
	// 1 foo
	// is none => true
	// is none => true
	// is none => true
}

func ExampleZipWith() {
	type Data struct {
		A int
		B string
	}

	zipper := func (v1 int, v2 string) Data {
		return Data{
			A: v1,
			B: v2,
		}
	}

	maybeData := ZipWith(Some[int](1), Some[string]("foo"), zipper)
	fmt.Printf("is none => %v\n", maybeData.IsNone())
	d, err := maybeData.Take()
	fmt.Printf("err is nil => %v\n", err == nil)
	fmt.Printf("%d %s\n", d.A, d.B)

	maybeData = ZipWith(Some[int](1), None[string](), zipper)
	fmt.Printf("is none => %v\n", maybeData.IsNone())
	maybeData = ZipWith(None[int](), Some[string]("foo"), zipper)
	fmt.Printf("is none => %v\n", maybeData.IsNone())
	maybeData = ZipWith(None[int](), None[string](), zipper)
	fmt.Printf("is none => %v\n", maybeData.IsNone())

	// Output:
	// is none => false
	// err is nil => true
	// 1 foo
	// is none => true
	// is none => true
	// is none => true
}
