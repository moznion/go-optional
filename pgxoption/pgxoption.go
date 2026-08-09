// Package pgxoption integrates optional.Option[T] with pgx v5's type
// system (pgtype). Registering it on a connection's type map lets
// Option values be used directly as query arguments and scan
// destinations with the element type's native codec — no
// driver.Valuer/sql.Scanner detour, so type fidelity and performance
// match a plain *T.
//
//	cfg, _ := pgxpool.ParseConfig(dsn)
//	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
//		pgxoption.Register(conn.TypeMap())
//		return nil
//	}
//
// Semantics: None encodes as SQL NULL and NULL scans as None; Some(v)
// encodes/scans exactly like v itself.
package pgxoption

import (
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

const optionPkgPath = "github.com/moznion/go-optional"

// Register installs Option[T] wrapper plans into m, in front of the
// built-in wrappers. It is safe to call once per connection type map;
// registering on pgtype.NewMap()-derived maps (the default for every
// pgx connection) covers all element types the map itself supports.
func Register(m *pgtype.Map) {
	m.TryWrapScanPlanFuncs = append(
		[]pgtype.TryWrapScanPlanFunc{tryWrapOptionScanPlan}, m.TryWrapScanPlanFuncs...)
	m.TryWrapEncodePlanFuncs = append(
		[]pgtype.TryWrapEncodePlanFunc{tryWrapOptionEncodePlan}, m.TryWrapEncodePlanFuncs...)
}

// optionElem reports whether rt is optional.Option[E] for some E,
// returning E. Option's underlying representation is a slice of the
// element type, which is what the reflection here relies on; the
// package path and type-name check keep ordinary slices out.
func optionElem(rt reflect.Type) (reflect.Type, bool) {
	if rt.Kind() != reflect.Slice || rt.PkgPath() != optionPkgPath ||
		!strings.HasPrefix(rt.Name(), "Option[") {
		return nil, false
	}
	return rt.Elem(), true
}

// tryWrapOptionScanPlan wraps a scan into *Option[E] as a scan into
// **E: pgx's own pointer handling turns SQL NULL into a nil *E, which
// maps to None; a non-nil *E maps to Some.
func tryWrapOptionScanPlan(target any) (pgtype.WrappedScanPlanNextSetter, any, bool) {
	tv := reflect.ValueOf(target)
	if tv.Kind() != reflect.Pointer || tv.IsNil() {
		return nil, nil, false
	}
	elem, ok := optionElem(tv.Type().Elem())
	if !ok {
		return nil, nil, false
	}
	return &optionScanPlan{elemType: elem}, reflect.New(reflect.PointerTo(elem)).Interface(), true
}

type optionScanPlan struct {
	next     pgtype.ScanPlan
	elemType reflect.Type
}

func (p *optionScanPlan) SetNext(next pgtype.ScanPlan) { p.next = next }

func (p *optionScanPlan) Scan(src []byte, dst any) error {
	holder := reflect.New(reflect.PointerTo(p.elemType)) // **E
	if err := p.next.Scan(src, holder.Interface()); err != nil {
		return err
	}
	opt := reflect.ValueOf(dst).Elem() // Option[E]
	ptr := holder.Elem()               // *E, nil on SQL NULL
	if ptr.IsNil() {
		opt.SetZero() // None
		return nil
	}
	some := reflect.MakeSlice(opt.Type(), 1, 1)
	some.Index(0).Set(ptr.Elem())
	opt.Set(some)
	return nil
}

// tryWrapOptionEncodePlan encodes Option[E] arguments: None becomes
// SQL NULL, Some(v) delegates to E's own encode plan.
func tryWrapOptionEncodePlan(value any) (pgtype.WrappedEncodePlanNextSetter, any, bool) {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return nil, nil, false
	}
	elem, ok := optionElem(rv.Type())
	if !ok {
		return nil, nil, false
	}
	// The next value handed to plan resolution must be a concrete E;
	// for None that is E's zero value (only used to pick the plan,
	// never encoded).
	next := reflect.Zero(elem).Interface()
	if rv.Len() > 0 {
		next = rv.Index(0).Interface()
	}
	return &optionEncodePlan{}, next, true
}

type optionEncodePlan struct {
	next pgtype.EncodePlan
}

func (p *optionEncodePlan) SetNext(next pgtype.EncodePlan) { p.next = next }

func (p *optionEncodePlan) Encode(value any, buf []byte) ([]byte, error) {
	rv := reflect.ValueOf(value)
	if rv.Len() == 0 {
		return nil, nil // SQL NULL
	}
	return p.next.Encode(rv.Index(0).Interface(), buf)
}
