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
func Nullify(obj any) any {
	typeOf := reflect.TypeOf(obj)
	if typeOf == nil {
		return nil // guard for nil interface{}
	}

	val := ptr(typeOf)
	return reflect.New(val.Elem()).Interface()
}

// ptr recursively transforms the `reflect.Type` to a pointer kind.
func ptr(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Struct:
		structFields := make([]reflect.StructField, t.NumField())
		for i := range structFields {
			structFields[i] = t.Field(i)
			structFields[i].Type = ptr(structFields[i].Type)
		}
		return reflect.PointerTo(reflect.StructOf(structFields))
	case reflect.Array:
		return reflect.PointerTo(reflect.ArrayOf(t.Len(), ptr(t.Elem())))
	case reflect.Slice:
		return reflect.PointerTo(reflect.SliceOf(ptr(t.Elem())))
	case reflect.Map:
		return reflect.PointerTo(reflect.MapOf(t.Key(), ptr(t.Elem())))
	// primitive types, just return the pointer value
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return reflect.PointerTo(t)
	// recursively follow pointer and return the non-pointer version, then call ptr on that to resolve to a 1-depth pointer
	case reflect.Pointer:
		for ok := t.Kind() == reflect.Pointer; ok; ok = t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		return ptr(t)
	default:
		return reflect.PointerTo(t)
	}
}
