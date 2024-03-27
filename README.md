# Nullify

Returns the pointer version of any input recursively including e.g. field structs (while retaining tags). This is especially
useful in e.g. JSON serialization/deserialization in combination with a struct validator to check if a field was sent or not
without changing the struct.

## Install

Use `github.com/Emptyless/nullify` to download the latest version.


## Example

To quickly get up and running, run e.g. the following:

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Person struct {
	Name string `json:"name" validate:"required"`
}

func main() {
	input := []byte("")
	person := Person{}
	p := Nullify(person)
	_ = json.Unmarshal(input, p)
	err := validator.New().Struct(p)
	fmt.Println(err)
	// Output:
	// Key: 'Name' Error:Field validation for 'Name' failed on the 'required' tag
}
```

For the full example, see  `/example` for an example with [go-playground/validator](https://github.com/go-playground/validator).