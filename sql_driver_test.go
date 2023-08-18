package optional

import (
	"database/sql"
	"database/sql/driver"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestOption_Scan(t *testing.T) {
	o := Some[any](nil)

	err := o.Scan("bar")
	assert.NoError(t, err)
	assert.EqualValues(t, "bar", o.Unwrap())

	err = o.Scan([]byte("buz"))
	assert.NoError(t, err)
	assert.EqualValues(t, []byte("buz"), o.Unwrap())

	err = o.Scan(int64(42))
	assert.NoError(t, err)
	assert.EqualValues(t, 42, o.Unwrap())

	err = o.Scan(float64(123.456))
	assert.NoError(t, err)
	assert.EqualValues(t, 123.456, o.Unwrap())

	err = o.Scan(true)
	assert.NoError(t, err)
	assert.EqualValues(t, true, o.Unwrap())

	now := time.Now()
	err = o.Scan(now)
	assert.NoError(t, err)
	assert.EqualValues(t, now, o.Unwrap())
}

func TestOption_Scan_None(t *testing.T) {
	o := Some[any](nil)

	err := o.Scan(nil)
	assert.NoError(t, err)
	assert.True(t, o.IsNone())
}

func TestOption_Scan_UnsupportedTypes(t *testing.T) {
	type ustruct struct {
		A int
	}

	o := Some[ustruct](ustruct{})
	err := o.Scan(int32(42))
	assert.Error(t, err)
}

func TestOption_Scan_ScannerInterfaceSatisfaction(t *testing.T) {
	o := Some[any]("string")
	var s sql.Scanner = &o
	assert.NotNil(t, s)
}

func TestOption_Value(t *testing.T) {
	{
		o := Some[string]("foo")
		v, err := o.Value()
		assert.NoError(t, err)
		assert.EqualValues(t, "foo", v)
	}

	{
		o := Some[[]byte]([]byte("bar"))
		v, err := o.Value()
		assert.NoError(t, err)
		assert.EqualValues(t, []byte("bar"), v)
	}

	{
		o := Some[int64](42)
		v, err := o.Value()
		assert.NoError(t, err)
		assert.EqualValues(t, 42, v)
	}

	{
		o := Some[float64](123.456)
		v, err := o.Value()
		assert.NoError(t, err)
		assert.EqualValues(t, 123.456, v)
	}

	{
		o := Some[bool](true)
		v, err := o.Value()
		assert.NoError(t, err)
		assert.EqualValues(t, true, v)
	}

	{
		now := time.Now()
		o := Some[time.Time](now)
		v, err := o.Value()
		assert.NoError(t, err)
		assert.EqualValues(t, now, v)
	}
}

func TestOption_Value_None(t *testing.T) {
	o := None[string]()
	v, err := o.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func TestOption_Value_UnsupportedTypes(t *testing.T) {
	type ustruct struct {
		A int
	}

	o := Some[ustruct](ustruct{})
	_, err := o.Value()
	assert.Error(t, err)
}

func TestOption_Value_ValuerInterfaceSatisfaction(t *testing.T) {
	o := Some[any]("string")
	var s driver.Valuer = &o
	assert.NotNil(t, s)
}

func TestOption_SQLScan(t *testing.T) {
	tmpfile, err := os.CreateTemp(os.TempDir(), "testdb")
	assert.NoError(t, err)

	db, err := sql.Open("sqlite3", tmpfile.Name())
	assert.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	sqlStmt := "CREATE TABLE test_table (id INTEGER NOT NULL PRIMARY KEY, name VARCHAR(32));"
	_, err = db.Exec(sqlStmt)
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)
	func() {
		stmt, err := tx.Prepare("INSERT INTO test_table(id, name) values(?, ?)")
		assert.NoError(t, err)
		defer func() {
			_ = stmt.Close()
		}()
		_, err = stmt.Exec(1, "foo")
		assert.NoError(t, err)
	}()
	func() {
		stmt, err := tx.Prepare("INSERT INTO test_table(id) values(?)")
		assert.NoError(t, err)
		defer func() {
			_ = stmt.Close()
		}()
		_, err = stmt.Exec(2)
		assert.NoError(t, err)
	}()
	err = tx.Commit()
	assert.NoError(t, err)

	var maybeName Option[string]

	row := db.QueryRow("SELECT name FROM test_table WHERE id = 1")
	err = row.Scan(&maybeName)
	assert.NoError(t, err)
	assert.Equal(t, "foo", maybeName.Unwrap())

	row = db.QueryRow("SELECT name FROM test_table WHERE id = 2")
	err = row.Scan(&maybeName)
	assert.NoError(t, err)
	assert.True(t, maybeName.IsNone())
}

func TestOption_SQLValuer(t *testing.T) {
	tmpfile, err := os.CreateTemp(os.TempDir(), "testdb")
	assert.NoError(t, err)

	db, err := sql.Open("sqlite3", tmpfile.Name())
	assert.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	sqlStmt := "CREATE TABLE test_table (id INTEGER NOT NULL PRIMARY KEY, name VARCHAR(32));"
	_, err = db.Exec(sqlStmt)
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)
	func() {
		stmt, err := tx.Prepare("INSERT INTO test_table(id, name) values(?, ?)")
		assert.NoError(t, err)
		defer func() {
			_ = stmt.Close()
		}()
		_, err = stmt.Exec(1, Some[string]("foo"))
		assert.NoError(t, err)
	}()
	func() {
		stmt, err := tx.Prepare("INSERT INTO test_table(id, name) values(?, ?)")
		assert.NoError(t, err)
		defer func() {
			_ = stmt.Close()
		}()
		_, err = stmt.Exec(2, None[string]())
		assert.NoError(t, err)
	}()
	err = tx.Commit()
	assert.NoError(t, err)

	var maybeName Option[string]

	row := db.QueryRow("SELECT name FROM test_table WHERE id = 1")
	err = row.Scan(&maybeName)
	assert.NoError(t, err)
	assert.Equal(t, "foo", maybeName.Unwrap())

	row = db.QueryRow("SELECT name FROM test_table WHERE id = 2")
	err = row.Scan(&maybeName)
	assert.NoError(t, err)
	assert.True(t, maybeName.IsNone())
}
