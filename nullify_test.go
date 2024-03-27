package nullify

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"unsafe"
)

func TestNullify_Nil(t *testing.T) {
	// Arrange
	var inf interface{} // the nil interface

	// Act
	p := Nullify(inf)

	// Assert
	assert.Nil(t, p)
}

func TestNullify_Struct(t *testing.T) {
	// Arrange
	type Person struct {
		Name string `json:"name"`
	}

	person := Person{}

	// Act
	p := Nullify(person)

	// Assert
	assert.Equal(t, reflect.Pointer, reflect.TypeOf(p).Kind())
	assert.Equal(t, reflect.Struct, reflect.TypeOf(p).Elem().Kind())
	assert.Equal(t, reflect.Pointer, reflect.TypeOf(p).Elem().Field(0).Type.Kind())
	assert.Equal(t, reflect.String, reflect.TypeOf(p).Elem().Field(0).Type.Elem().Kind())
	assert.Equal(t, reflect.StructTag(`json:"name"`), reflect.TypeOf(p).Elem().Field(0).Tag)
}

func TestNullify_Array(t *testing.T) {
	// Arrange
	var input [1]string
	input[0] = "test"

	// Act
	i := Nullify(input)

	// Assert
	assert.Equal(t, reflect.Pointer, reflect.TypeOf(i).Kind())
	assert.Equal(t, reflect.Array, reflect.TypeOf(i).Elem().Kind())
	assert.Equal(t, reflect.Pointer, reflect.TypeOf(i).Elem().Elem().Kind())
	assert.Equal(t, reflect.String, reflect.TypeOf(i).Elem().Elem().Elem().Kind())
}

func TestNullify_Slice(t *testing.T) {
	// Arrange
	input := []string{"1"}

	// Act
	i := Nullify(input)

	// Assert
	assert.Equal(t, reflect.Pointer, reflect.TypeOf(i).Kind())
	assert.Equal(t, reflect.Slice, reflect.TypeOf(i).Elem().Kind())
	assert.Equal(t, reflect.Pointer, reflect.TypeOf(i).Elem().Elem().Kind())
	assert.Equal(t, reflect.String, reflect.TypeOf(i).Elem().Elem().Elem().Kind())
}

func TestNullify_Map(t *testing.T) {
	// Arrange
	input := map[string]int{}

	// Act
	i := Nullify(input)

	// Assert
	assert.Equal(t, reflect.Pointer, reflect.TypeOf(i).Kind())
	assert.Equal(t, reflect.Map, reflect.TypeOf(i).Elem().Kind())
	assert.Equal(t, reflect.String, reflect.TypeOf(i).Elem().Key().Kind())
	assert.Equal(t, reflect.Pointer, reflect.TypeOf(i).Elem().Elem().Kind())
	assert.Equal(t, reflect.Int, reflect.TypeOf(i).Elem().Elem().Elem().Kind())
}

func TestNullify_Primitive(t *testing.T) {
	tests := map[string]struct {
		Input  any
		Output reflect.Kind
	}{
		"Bool":       {Input: true, Output: reflect.Bool},
		"Int":        {Input: 1, Output: reflect.Int},
		"Int8":       {Input: int8(1), Output: reflect.Int8},
		"Int16":      {Input: int16(2), Output: reflect.Int16},
		"Int32":      {Input: int32(3), Output: reflect.Int32},
		"Int64":      {Input: int64(4), Output: reflect.Int64},
		"Uint":       {Input: uint(5), Output: reflect.Uint},
		"Uint8":      {Input: uint8(6), Output: reflect.Uint8},
		"Uint16":     {Input: uint16(7), Output: reflect.Uint16},
		"Uint32":     {Input: uint32(8), Output: reflect.Uint32},
		"Uint64":     {Input: uint64(9), Output: reflect.Uint64},
		"Float32":    {Input: float32(10.0), Output: reflect.Float32},
		"Float64":    {Input: float64(11.0), Output: reflect.Float64},
		"Complex64":  {Input: complex64(12.0), Output: reflect.Complex64},
		"Complex128": {Input: complex128(13.0), Output: reflect.Complex128},
		"String":     {Input: "String", Output: reflect.String},
	}
	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			// Arrange
			input := testData.Input

			// Act
			output := Nullify(input)

			// Assert
			assert.Equal(t, reflect.Pointer, reflect.TypeOf(output).Kind())
			assert.Equal(t, testData.Output, reflect.TypeOf(output).Elem().Kind())
		})
	}
}

func TestNullify_Pointer(t *testing.T) {
	// Arrange
	var input *string

	// Act
	i := Nullify(input)

	// Assert
	assert.Equal(t, reflect.Pointer, reflect.TypeOf(i).Kind())
	assert.Equal(t, reflect.String, reflect.TypeOf(i).Elem().Kind())
}

func TestNullify_Default(t *testing.T) {
	str := "test"

	tests := map[string]struct {
		Input any
		Kind  reflect.Kind
	}{
		"Chan":          {Input: make(chan int), Kind: reflect.Chan},
		"UnsafePointer": {Input: unsafe.Pointer(&str), Kind: reflect.UnsafePointer},
		"Func":          {Input: func() {}, Kind: reflect.Func},
	}
	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			// Arrange
			input := testData.Input

			// Act
			output := Nullify(input)

			// Assert
			assert.Equal(t, reflect.Pointer, reflect.TypeOf(output).Kind())
			assert.Equal(t, testData.Kind, reflect.TypeOf(output).Elem().Kind())
		})
	}
}
