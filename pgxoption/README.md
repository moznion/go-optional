# pgxoption

[![Go Reference](https://pkg.go.dev/badge/github.com/moznion/go-optional/pgxoption.svg)](https://pkg.go.dev/github.com/moznion/go-optional/pgxoption)

Native [pgx v5](https://github.com/jackc/pgx) (`pgtype`) integration for
[`optional.Option[T]`](https://github.com/moznion/go-optional), as a
nested Go module so the core package stays dependency-free.

Without this module, pgx handles `Option[T]` through the generic
`sql.Scanner` / `driver.Valuer` fallback, which detours through
`driver.Value` — extra conversions, and lossy for types like
`numeric`. Registering pgxoption plugs `Option[T]` directly into pgx's
plan machinery: `Some(v)` encodes/scans with `v`'s own codec (text and
binary formats alike), `None` is SQL `NULL`.

## Usage

```console
go get github.com/moznion/go-optional/pgxoption
```

```go
cfg, _ := pgxpool.ParseConfig(dsn)
cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	pgxoption.Register(conn.TypeMap())
	return nil
}
pool, _ := pgxpool.NewWithConfig(ctx, cfg)

var nickname optional.Option[string]
err := pool.QueryRow(ctx,
	"SELECT nickname FROM users WHERE id = $1", optional.Some(int64(1)),
).Scan(&nickname) // NULL => None, value => Some
```

For a single connection, call `pgxoption.Register(conn.TypeMap())`
after `pgx.Connect`.

## Semantics

- `None[T]()` binds as SQL `NULL`; SQL `NULL` scans as `None[T]()`
  (overwriting any previous `Some`).
- `Some(v)` binds and scans exactly like a plain `v` / `*v` — the
  element type's registered codec does the work, so everything the
  connection's `pgtype.Map` supports is supported inside an `Option`.
- Plain slices are unaffected; the wrapper recognizes `Option[T]` by
  its defined type, not by its slice shape.
