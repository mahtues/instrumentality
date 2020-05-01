package api

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	var v *Asd

	input := `{"name": "bob", "PtrToNumber": 5}`

	json.Unmarshal([]byte(input), &v)

	fmt.Println("value:", v)

	//t.Error("just to print stuff")
}
