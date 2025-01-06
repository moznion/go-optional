package optional

import (
	"database/sql"
	"database/sql/driver"
)

// Scan assigns a value from a database driver.
// This method is required from database/sql.Scanner interface.
func (o *Option[T]) Scan(src any) error {
	// The detour through sql.Null[T] allows us to access the standard rules for
	// assigning scanned values into builtin types and std types like *sql.Rows,
	// which are not exported from std directly.
	var v sql.Null[T]
	err := v.Scan(src)
	if err != nil {
		return err
	}
	if v.Valid {
		*o = Some[T](v.V)
	} else {
		*o = None[T]()
	}
	return nil
}

// Value returns a driver Value.
// This method is required from database/sql/driver.Valuer interface.
func (o Option[T]) Value() (driver.Value, error) {
	if o.IsNone() {
		return nil, nil
	}
	return driver.DefaultParameterConverter.ConvertValue(o.Unwrap())
}
