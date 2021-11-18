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

func TestUnzip(t *testing.T) {
	pair := Pair[int, string]{
		Value1: 123,
		Value2: "foo",
	}

	o1, o2 := Unzip(Some[Pair[int, string]](pair))
	assert.Equal(t, 123, o1.TakeOr(0))
	assert.Equal(t, "foo", o2.TakeOr(""))

	o1, o2 = Unzip(None[Pair[int, string]]())
	assert.True(t, o1.IsNone())
	assert.True(t, o2.IsNone())
}

func TestUnzipWith(t *testing.T) {
	type Data struct {
		A string
		B int
	}

	unzipper := func (d Data) (string, int) {
		return d.A, d.B
	}

	o1, o2 := UnzipWith(Some[Data](Data{
		A: "foo",
		B: 123,
	}), unzipper)
	assert.Equal(t, "foo", o1.TakeOr(""))
	assert.Equal(t, 123, o2.TakeOr(0))

	o1, o2 = UnzipWith(None[Data](), unzipper)
	assert.True(t, o1.IsNone())
	assert.True(t, o2.IsNone())
}
