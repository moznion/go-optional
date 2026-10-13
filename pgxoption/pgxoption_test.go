package pgxoption

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moznion/go-optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func registeredMap() *pgtype.Map {
	m := pgtype.NewMap()
	Register(m)
	return m
}

func TestScan_TextFormat(t *testing.T) {
	m := registeredMap()

	var i optional.Option[int64]
	require.NoError(t, m.Scan(pgtype.Int8OID, pgtype.TextFormatCode, []byte("42"), &i))
	assert.Equal(t, optional.Some(int64(42)), i)

	require.NoError(t, m.Scan(pgtype.Int8OID, pgtype.TextFormatCode, nil, &i))
	assert.True(t, i.IsNone(), "SQL NULL must scan into None")

	var s optional.Option[string]
	require.NoError(t, m.Scan(pgtype.TextOID, pgtype.TextFormatCode, []byte("héllo"), &s))
	assert.Equal(t, optional.Some("héllo"), s)

	var f optional.Option[float64]
	require.NoError(t, m.Scan(pgtype.NumericOID, pgtype.TextFormatCode, []byte("1.25"), &f))
	assert.Equal(t, optional.Some(1.25), f)

	var b optional.Option[bool]
	require.NoError(t, m.Scan(pgtype.BoolOID, pgtype.TextFormatCode, []byte("t"), &b))
	assert.Equal(t, optional.Some(true), b)
}

func TestScan_BinaryFormat(t *testing.T) {
	m := registeredMap()

	var i optional.Option[int64]
	require.NoError(t, m.Scan(pgtype.Int8OID, pgtype.BinaryFormatCode,
		[]byte{0, 0, 0, 0, 0, 0, 0, 42}, &i))
	assert.Equal(t, optional.Some(int64(42)), i)

	require.NoError(t, m.Scan(pgtype.Int8OID, pgtype.BinaryFormatCode, nil, &i))
	assert.True(t, i.IsNone())
}

func TestScan_OverwritesPreviousValue(t *testing.T) {
	m := registeredMap()

	i := optional.Some(int64(7))
	require.NoError(t, m.Scan(pgtype.Int8OID, pgtype.TextFormatCode, nil, &i))
	assert.True(t, i.IsNone(), "scanning NULL must reset a previous Some")
}

func TestEncode(t *testing.T) {
	m := registeredMap()

	buf, err := m.Encode(pgtype.Int8OID, pgtype.TextFormatCode, optional.Some(int64(42)), nil)
	require.NoError(t, err)
	assert.Equal(t, []byte("42"), buf)

	buf, err = m.Encode(pgtype.Int8OID, pgtype.TextFormatCode, optional.None[int64](), nil)
	require.NoError(t, err)
	assert.Nil(t, buf, "None must encode as SQL NULL")

	buf, err = m.Encode(pgtype.TextOID, pgtype.TextFormatCode, optional.Some("héllo"), nil)
	require.NoError(t, err)
	assert.Equal(t, []byte("héllo"), buf)
}

func TestEncode_TimeRoundTrip(t *testing.T) {
	m := registeredMap()

	ts := time.Date(2026, 8, 9, 12, 34, 56, 0, time.UTC)
	buf, err := m.Encode(pgtype.TimestamptzOID, pgtype.BinaryFormatCode, optional.Some(ts), nil)
	require.NoError(t, err)
	require.NotNil(t, buf)

	var got optional.Option[time.Time]
	require.NoError(t, m.Scan(pgtype.TimestamptzOID, pgtype.BinaryFormatCode, buf, &got))
	require.True(t, got.IsSome())
	assert.True(t, got.Unwrap().Equal(ts))
}

func TestPlainSlicesAreUntouched(t *testing.T) {
	m := registeredMap()

	// A real slice destination must keep going through pgx's array
	// machinery, not the Option wrapper.
	var xs []int64
	require.NoError(t, m.Scan(pgtype.Int8ArrayOID, pgtype.TextFormatCode, []byte("{1,2}"), &xs))
	assert.Equal(t, []int64{1, 2}, xs)
}
