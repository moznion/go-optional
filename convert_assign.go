package optional

import (
	_ "unsafe"
)

//go:linkname sqlConvertAssign database/sql.convertAssign
func sqlConvertAssign(dest, src any) error
