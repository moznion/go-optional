package optional

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption_IsNone(t *testing.T) {
	assert.True(t, None[int]().IsNone())
	assert.False(t, Some[int](123).IsNone())

	var nilValue Option[int] = nil
	assert.True(t, nilValue.IsNone())

	i := 0
	assert.False(t, FromNillable[int](&i).IsNone())
	assert.True(t, FromNillable[int](nil).IsNone())
}

func TestOption_IsSome(t *testing.T) {
	assert.False(t, None[int]().IsSome())
	assert.True(t, Some[int](123).IsSome())

	var nilValue Option[int] = nil
	assert.False(t, nilValue.IsSome())

	i := 0
	assert.True(t, FromNillable[int](&i).IsSome())
	assert.False(t, FromNillable[int](nil).IsSome())
}

func TestOption_Unwrap(t *testing.T) {
	assert.Equal(t, "foo", Some[string]("foo").Unwrap())
	assert.Equal(t, "", None[string]().Unwrap())
	assert.Nil(t, None[*string]().Unwrap())

	i := 123
	assert.Equal(t, i, FromNillable[int](&i).Unwrap())
	assert.Equal(t, 0, FromNillable[int](nil).Unwrap())
	assert.Equal(t, i, *PtrFromNillable[int](&i).Unwrap())
	assert.Nil(t, PtrFromNillable[int](nil).Unwrap())
}

func TestOption_UnwrapAsPointer(t *testing.T) {
	str := "foo"
	refStr := &str
	assert.EqualValues(t, &str, Some[string](str).UnwrapAsPtr())
	assert.EqualValues(t, &refStr, Some[*string](refStr).UnwrapAsPtr())
	assert.Nil(t, None[string]().UnwrapAsPtr())
	assert.Nil(t, None[*string]().UnwrapAsPtr())

	i := 123
	assert.Equal(t, &i, FromNillable[int](&i).UnwrapAsPtr())
	assert.Nil(t, FromNillable[int](nil).UnwrapAsPtr())
	assert.Equal(t, &i, *PtrFromNillable[int](&i).UnwrapAsPtr())
	assert.Nil(t, PtrFromNillable[int](nil).UnwrapAsPtr())
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
	v := Some[int](123).TakeOrElse(func() int {
		return 666
	})
	assert.Equal(t, 123, v)

	v = None[int]().TakeOrElse(func() int {
		return 666
	})
	assert.Equal(t, 666, v)
}

func TestOption_Filter(t *testing.T) {
	isEven := func(v int) bool {
		return v%2 == 0
	}

	o := Some[int](2).Filter(isEven)
	assert.True(t, o.IsSome())
	assert.Equal(t, 2, o[value])

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
	}, zipped[value])

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

	zipped := ZipWith(some1, some2, func(v1 int, v2 string) Data {
		return Data{
			A: v2,
			B: v1,
		}
	})
	assert.True(t, zipped.IsSome())
	assert.Equal(t, Data{
		A: "foo",
		B: 123,
	}, zipped[value])

	assert.True(t, ZipWith(None[int](), some1, func(v1, v2 int) Data {
		return Data{}
	}).IsNone())
	assert.True(t, ZipWith(some1, None[int](), func(v1, v2 int) Data {
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

	unzipper := func(d Data) (string, int) {
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

func TestMapWithError(t *testing.T) {
	some := Some[int](123)
	mapped, err := MapWithError(some, func(v int) (string, error) {
		return fmt.Sprintf("%d", v), nil
	})
	assert.NoError(t, err)
	taken, err := mapped.Take()
	assert.NoError(t, err)
	assert.Equal(t, "123", taken)

	none := None[int]()
	mapped, err = MapWithError(none, func(v int) (string, error) {
		return fmt.Sprintf("%d", v), nil
	})
	assert.NoError(t, err)
	assert.True(t, mapped.IsNone())

	mapperError := errors.New("mapper error")
	mapped, err = MapWithError(some, func(v int) (string, error) {
		return "", mapperError
	})
	assert.ErrorIs(t, err, mapperError)
	assert.True(t, mapped.IsNone())
}

func TestMapOrWithError(t *testing.T) {
	some := Some[int](123)
	mapped, err := MapOrWithError(some, "666", func(v int) (string, error) {
		return fmt.Sprintf("%d", v), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "123", mapped)

	none := None[int]()
	mapped, err = MapOrWithError(none, "666", func(v int) (string, error) {
		return fmt.Sprintf("%d", v), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "666", mapped)

	mapperError := errors.New("mapper error")
	mapped, err = MapOrWithError(some, "666", func(v int) (string, error) {
		return "", mapperError
	})
	assert.ErrorIs(t, err, mapperError)
	assert.Equal(t, "", mapped)
}

func TestOption_IfSome(t *testing.T) {
	callingValue := ""
	Some("foo").IfSome(func(s string) {
		callingValue = s
	})
	assert.Equal(t, "foo", callingValue)

	callingValue = ""
	None[string]().IfSome(func(s string) {
		callingValue = s
	})
	assert.Equal(t, "", callingValue)
}

func TestOption_IfSomeWithError(t *testing.T) {
	err := Some("foo").IfSomeWithError(func(s string) error {
		return nil
	})
	assert.NoError(t, err)

	err = Some("foo").IfSomeWithError(func(s string) error {
		return errors.New(s)
	})
	assert.EqualError(t, err, "foo")

	err = None[string]().IfSomeWithError(func(s string) error {
		return errors.New(s)
	})
	assert.NoError(t, err)
}

func TestOption_IfNone(t *testing.T) {
	called := false
	None[string]().IfNone(func() {
		called = true
	})
	assert.True(t, called)

	called = false
	Some("string").IfNone(func() {
		called = true
	})
	assert.False(t, called)
}

func TestOption_IfNoneWithError(t *testing.T) {
	err := None[string]().IfNoneWithError(func() error {
		return nil
	})
	assert.NoError(t, err)

	err = None[string]().IfNoneWithError(func() error {
		return errors.New("err")
	})
	assert.EqualError(t, err, "err")

	err = Some("foo").IfNoneWithError(func() error {
		return errors.New("err")
	})
	assert.NoError(t, err)
}

func TestFlatMap(t *testing.T) {
	some := Some[int](123)
	mapped := FlatMap(some, func(v int) Option[string] {
		return Some[string](fmt.Sprintf("%d", v))
	})
	taken, err := mapped.Take()
	assert.NoError(t, err)
	assert.Equal(t, "123", taken)

	none := None[int]()
	mapped = FlatMap(none, func(v int) Option[string] {
		return Some[string](fmt.Sprintf("%d", v))
	})
	assert.True(t, mapped.IsNone())
}

func TestFlatMapOr(t *testing.T) {
	some := Some[int](123)
	mapped := FlatMapOr(some, "666", func(v int) Option[string] {
		return Some[string](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, "123", mapped)

	none := None[int]()
	mapped = FlatMapOr(none, "666", func(v int) Option[string] {
		return Some[string](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, "666", mapped)
}

func TestFlatMapWithError(t *testing.T) {
	some := Some[int](123)
	mapped, err := FlatMapWithError(some, func(v int) (Option[string], error) {
		return Some[string](fmt.Sprintf("%d", v)), nil
	})
	assert.NoError(t, err)
	taken, err := mapped.Take()
	assert.NoError(t, err)
	assert.Equal(t, "123", taken)

	none := None[int]()
	mapped, err = FlatMapWithError(none, func(v int) (Option[string], error) {
		return Some[string](fmt.Sprintf("%d", v)), nil
	})
	assert.NoError(t, err)
	assert.True(t, mapped.IsNone())

	mapperError := errors.New("mapper error")
	mapped, err = FlatMapWithError(some, func(v int) (Option[string], error) {
		return Some[string](""), mapperError
	})
	assert.ErrorIs(t, err, mapperError)
	assert.True(t, mapped.IsNone())
}

func TestFlatMapOrWithError(t *testing.T) {
	some := Some[int](123)
	mapped, err := FlatMapOrWithError(some, "666", func(v int) (Option[string], error) {
		return Some[string](fmt.Sprintf("%d", v)), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "123", mapped)

	none := None[int]()
	mapped, err = FlatMapOrWithError(none, "666", func(v int) (Option[string], error) {
		return Some[string](fmt.Sprintf("%d", v)), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "666", mapped)

	mapperError := errors.New("mapper error")
	mapped, err = FlatMapOrWithError(some, "666", func(v int) (Option[string], error) {
		return Some[string](""), mapperError
	})
	assert.ErrorIs(t, err, mapperError)
	assert.Equal(t, "", mapped)
}

func TestOptionSerdeJSONForSomeValue(t *testing.T) {
	{
		type JSONStruct struct {
			Val Option[int] `json:"val"`
		}

		some := Some[int](123)
		jsonStruct := &JSONStruct{Val: some}

		marshal, err := json.Marshal(jsonStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, string(marshal), `{"val":123}`)

		var unmarshalJSONStruct JSONStruct
		err = json.Unmarshal(marshal, &unmarshalJSONStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
	}

	{
		type JSONStruct struct {
			Val Option[string] `json:"val"`
		}

		some := Some[string]("foobar")
		jsonStruct := &JSONStruct{Val: some}

		marshal, err := json.Marshal(jsonStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, string(marshal), `{"val":"foobar"}`)

		var unmarshalJSONStruct JSONStruct
		err = json.Unmarshal(marshal, &unmarshalJSONStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
	}

	{
		type JSONStruct struct {
			Val Option[bool] `json:"val"`
		}

		some := Some[bool](false)
		jsonStruct := &JSONStruct{Val: some}

		marshal, err := json.Marshal(jsonStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, string(marshal), `{"val":false}`)

		var unmarshalJSONStruct JSONStruct
		err = json.Unmarshal(marshal, &unmarshalJSONStruct)
		assert.NoError(t, err)
		assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
	}

	{
		type Inner struct {
			B *bool `json:"b,omitempty"`
		}
		type JSONStruct struct {
			Val Option[Inner] `json:"val"`
		}

		{
			falsy := false
			some := Some[Inner](Inner{
				B: &falsy,
			})
			jsonStruct := &JSONStruct{Val: some}

			marshal, err := json.Marshal(jsonStruct)
			assert.NoError(t, err)
			assert.EqualValues(t, string(marshal), `{"val":{"b":false}}`)

			var unmarshalJSONStruct JSONStruct
			err = json.Unmarshal(marshal, &unmarshalJSONStruct)
			assert.NoError(t, err)
			assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
		}

		{
			some := Some[Inner](Inner{
				B: nil,
			})
			jsonStruct := &JSONStruct{Val: some}

			marshal, err := json.Marshal(jsonStruct)
			assert.NoError(t, err)
			assert.EqualValues(t, string(marshal), `{"val":{}}`)

			var unmarshalJSONStruct JSONStruct
			err = json.Unmarshal(marshal, &unmarshalJSONStruct)
			assert.NoError(t, err)
			assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
		}
	}
}

func TestOptionSerdeJSONForNoneValue(t *testing.T) {
	type JSONStruct struct {
		Val Option[int] `json:"val"`
	}
	some := None[int]()
	jsonStruct := &JSONStruct{Val: some}

	marshal, err := json.Marshal(jsonStruct)
	assert.NoError(t, err)
	assert.EqualValues(t, string(marshal), `{"val":null}`)

	var unmarshalJSONStruct JSONStruct
	err = json.Unmarshal(marshal, &unmarshalJSONStruct)
	assert.NoError(t, err)
	assert.EqualValues(t, jsonStruct, &unmarshalJSONStruct)
}

func TestOption_UnmarshalJSON_withEmptyJSONString(t *testing.T) {
	type JSONStruct struct {
		Val Option[int] `json:"val"`
	}

	var unmarshalJSONStruct JSONStruct
	err := json.Unmarshal([]byte("{}"), &unmarshalJSONStruct)
	assert.NoError(t, err)
	assert.EqualValues(t, &JSONStruct{
		Val: None[int](),
	}, &unmarshalJSONStruct)
}

func TestOption_MarshalJSON_shouldReturnErrorWhenInvalidJSONStructInputHasCome(t *testing.T) {
	type JSONStruct struct {
		Val Option[chan interface{}] `json:"val"` // chan type is unsupported on json marshaling
	}

	ch := make(chan interface{})
	some := Some[chan interface{}](ch)
	jsonStruct := &JSONStruct{Val: some}
	_, err := json.Marshal(jsonStruct)
	assert.Error(t, err)
}

func TestOption_UnmarshalJSON_shouldReturnErrorWhenInvalidJSONStringInputHasCome(t *testing.T) {
	type JSONStruct struct {
		Val Option[int] `json:"val"`
	}

	var unmarshalJSONStruct JSONStruct
	err := json.Unmarshal([]byte(`{"val":"__STRING__"}`), &unmarshalJSONStruct)
	assert.Error(t, err)
}

func TestOption_MarshalJSON_shouldHandleOmitemptyCorrectly(t *testing.T) {
	type JSONStruct struct {
		NormalVal    Option[string] `json:"normalVal"`
		OmitemptyVal Option[string] `json:"omitemptyVal,omitempty"` // this should be omitted
	}

	none := None[string]()
	jsonStruct := &JSONStruct{NormalVal: none, OmitemptyVal: none}
	marshal, err := json.Marshal(jsonStruct)
	assert.NoError(t, err)
	assert.EqualValues(t, string(marshal), `{"normalVal":null}`)
}

type MyStringer struct {
}

func (m *MyStringer) String() string {
	return "mystr"
}

func TestOption_String(t *testing.T) {
	assert.Equal(t, "Some[123]", Some[int](123).String())
	assert.Equal(t, "None[]", None[int]().String())

	assert.Equal(t, "Some[mystr]", Some[*MyStringer](&MyStringer{}).String())
	assert.Equal(t, "None[]", None[*MyStringer]().String())
}

func TestOption_Or(t *testing.T) {
	fallback := Some[string]("fallback")

	assert.EqualValues(t, Some[string]("actual").Or(fallback).Unwrap(), "actual")
	assert.EqualValues(t, None[string]().Or(fallback).Unwrap(), "fallback")
}

func TestOption_OrElse(t *testing.T) {
	fallbackFunc := func() Option[string] { return Some[string]("fallback") }

	assert.EqualValues(t, Some[string]("actual").OrElse(fallbackFunc).Unwrap(), "actual")
	assert.EqualValues(t, None[string]().OrElse(fallbackFunc).Unwrap(), "fallback")
}
