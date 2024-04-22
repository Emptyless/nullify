package nullify

import (
	"reflect"
)

// Nullify returns the pointer version of any input, e.g. string becomes *string, int becomes *int
// and for a struct, the pointer version of the fields is returned as well. E.g.
//
//	type Person struct {
//	   Name string
//	}
//
// will be returned as
//
//	type Person struct {
//	   Name *string
//	}
//
// with `p := Person{}`, Nullify(p) returns a pointer to Person.
//
// This is especially useful in e.g. validating JSON input, see example.
func Nullify(obj any, options ...option) any {
	typeOf := reflect.TypeOf(obj)
	if typeOf == nil {
		return nil // guard for nil interface{}
	}

	// default config
	cfg := config{
		bytesAsString:    false,
		nullifyArrayElem: true,
		nullifySliceElem: true,
		nullifyMapElem:   true,
		nullifyMapKey:    true,
	}

	// process options
	for _, opt := range options {
		cfg = opt.update(cfg)
	}

	val := ptr(typeOf, cfg)
	return reflect.New(val.Elem()).Interface()
}

// JsonOptions is a curated list of options that can be used for json.Marshal, json.Unmarshal.
// Use by spreading it onto the nullify function: `Nullify(t, JsonOptions...)
var JsonOptions = []option{
	BytesAsString{Value: true},
	NullifyMapKey{Value: false},
	NullifyMapElem{Value: false},
	NullifySliceElem{Value: false},
	NullifyArrayElem{Value: false},
}

// config determines the behavior of the ptr function
type config struct {
	bytesAsString    bool
	nullifyArrayElem bool
	nullifySliceElem bool
	nullifyMapElem   bool
	nullifyMapKey    bool
}

// option functionally updates the ptr function
type option interface {
	update(cfg config) config
}

// BytesAsString if true (default false) processes []uint8, []byte as string
// this is especially useful in json.Marshal, json.Unmarshal cases
type BytesAsString struct {
	Value bool
}

func (o BytesAsString) update(cfg config) config {
	cfg.bytesAsString = o.Value
	return cfg
}

// NullifyArrayElem if true (default false) doesn't nullify the array element, e.g. []any instead of []*any
type NullifyArrayElem struct {
	Value bool
}

func (o NullifyArrayElem) update(cfg config) config {
	cfg.nullifyArrayElem = o.Value
	return cfg
}

// NullifySliceElem if true (default false) doesn't nullify the slice element, e.g. []any instead of []*any
type NullifySliceElem struct {
	Value bool
}

func (o NullifySliceElem) update(cfg config) config {
	cfg.nullifySliceElem = o.Value
	return cfg
}

// NullifyMapElem if true (default true) nullifies the map element, e.g. map[any]*any instead of map[any]any
type NullifyMapElem struct {
	Value bool
}

func (o NullifyMapElem) update(cfg config) config {
	cfg.nullifyMapElem = o.Value
	return cfg
}

// NullifyMapKey if true (default true) nullifies the map element, e.g. map[*any]any instead of map[any]any
type NullifyMapKey struct {
	Value bool
}

func (o NullifyMapKey) update(cfg config) config {
	cfg.nullifyMapKey = o.Value
	return cfg
}

// ptr recursively transforms the `reflect.Type` to a pointer kind.
func ptr(t reflect.Type, cfg config) reflect.Type {
	switch t.Kind() {
	case reflect.Struct:
		structFields := make([]reflect.StructField, t.NumField())
		for i := range structFields {
			structFields[i] = t.Field(i)
			structFields[i].Type = ptr(structFields[i].Type, cfg)
		}
		return reflect.PointerTo(reflect.StructOf(structFields))
	case reflect.Array:
		if cfg.bytesAsString && (t.Elem().Kind() == reflect.Uint8 || (t.Elem().Kind() == reflect.Pointer && t.Elem().Elem().Kind() == reflect.Uint8)) {
			elemType := reflect.PointerTo(reflect.TypeOf(""))
			return elemType
		}

		elemType := ptr(t.Elem(), cfg)
		if cfg.nullifyArrayElem && elemType.Kind() != reflect.Pointer {
			elemType = reflect.PointerTo(elemType)
		}
		if !cfg.nullifyArrayElem && elemType.Kind() == reflect.Pointer {
			elemType = elemType.Elem()
		}

		return reflect.PointerTo(reflect.ArrayOf(t.Len(), elemType))
	case reflect.Slice:
		if cfg.bytesAsString && (t.Elem().Kind() == reflect.Uint8 || (t.Elem().Kind() == reflect.Pointer && t.Elem().Elem().Kind() == reflect.Uint8)) {
			elemType := reflect.TypeOf("")
			if cfg.nullifySliceElem {
				elemType = reflect.PointerTo(elemType)
			}
			return elemType
		}

		elemType := ptr(t.Elem(), cfg)
		if cfg.nullifySliceElem && elemType.Kind() != reflect.Pointer {
			elemType = reflect.PointerTo(elemType)
		}
		if !cfg.nullifySliceElem && elemType.Kind() == reflect.Pointer {
			elemType = elemType.Elem()
		}

		return reflect.PointerTo(reflect.SliceOf(elemType))
	case reflect.Map:
		elemType := ptr(t.Elem(), cfg)
		if cfg.nullifyMapElem && elemType.Kind() != reflect.Pointer {
			elemType = reflect.PointerTo(elemType)
		}
		if !cfg.nullifyMapElem && elemType.Kind() == reflect.Pointer {
			elemType = elemType.Elem()
		}

		keyType := ptr(t.Key(), cfg)
		if cfg.nullifyMapKey && keyType.Kind() != reflect.Pointer {
			keyType = reflect.PointerTo(keyType)
		}
		if !cfg.nullifyMapKey && keyType.Kind() == reflect.Pointer {
			keyType = keyType.Elem()
		}

		return reflect.PointerTo(reflect.MapOf(keyType, elemType))
	// primitive types, just return the pointer value
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return reflect.PointerTo(t)
	// recursively follow pointer and return the non-pointer version, then call ptr on that to resolve to a 1-depth pointer
	case reflect.Pointer:
		for ok := t.Kind() == reflect.Pointer; ok; ok = t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		return ptr(t, cfg)
	default:
		return reflect.PointerTo(t)
	}
}
