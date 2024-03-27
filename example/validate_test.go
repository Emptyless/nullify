package example

import (
	"encoding/json"
	"fmt"
	"github.com/Emptyless/nullify"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Some struct {
	Optional string `json:"optional" validate:"omitnil,email"`
	Required string `json:"required" validate:"required,uuid"`
}

func TestNullify_JsonUnmarshal(t *testing.T) {
	tests := map[string]struct {
		Payload      string
		Required     string
		Optional     string
		ErrorMessage string
	}{
		"missing all": {
			Payload:      `{}`,
			ErrorMessage: "Key: 'Required' Error:Field validation for 'Required' failed on the 'required' tag",
		},
		"missing optional": {
			Payload:  `{"required": "89ec270d-8256-4b0e-b25c-39564b10f29e"}`,
			Required: "89ec270d-8256-4b0e-b25c-39564b10f29e",
			Optional: "",
		},
		"invalid format required": {
			Payload:      `{"required": "invalid"}`,
			ErrorMessage: "Key: 'Required' Error:Field validation for 'Required' failed on the 'uuid' tag",
		},
		"invalid format optional": {
			Payload:      `{"required": "89ec270d-8256-4b0e-b25c-39564b10f29e", "optional": "notanemail"}`,
			ErrorMessage: "Key: 'Optional' Error:Field validation for 'Optional' failed on the 'email' tag",
		},
		"valid": {
			Payload:  `{"required": "89ec270d-8256-4b0e-b25c-39564b10f29e", "optional": "test@example.com"}`,
			Required: "89ec270d-8256-4b0e-b25c-39564b10f29e",
			Optional: "test@example.com",
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			// Arrange
			validate := validator.New()
			var some Some
			ptrSome := nullify.Nullify(&some)
			if err := json.Unmarshal([]byte(testData.Payload), ptrSome); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(testData.Payload), &some); err != nil {
				t.Fatal(err)
			}

			// Act
			err := validate.Struct(ptrSome)

			// Assert
			if testData.ErrorMessage == "" {
				assert.Nil(t, err)
				assert.Equal(t, testData.Required, some.Required)
				assert.Equal(t, testData.Optional, some.Optional)
			} else {
				assert.ErrorContains(t, err, testData.ErrorMessage)
			}
		})
	}
}

type Person struct {
	Name string `json:"name" validate:"required"`
}

func ExampleNullify() {
	input := []byte("")
	person := Person{}
	p := nullify.Nullify(person)
	_ = json.Unmarshal(input, p)
	err := validator.New().Struct(p)
	fmt.Println(err)
	// Output:
	// Key: 'Name' Error:Field validation for 'Name' failed on the 'required' tag
}
