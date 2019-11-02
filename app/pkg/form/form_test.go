package form

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type Example struct {
	Name       string  `form:"name"`
	Age        int     `form:"age"`
	Subscribed bool    `form:"subscribed"`
	Ptr        *int    `form:"ptr"`
	DPtr       **int   `form:"dptr"`
	Float      float64 `form:"float"`
}

func TestBasic(t *testing.T) {
	r := &http.Request{Method: http.MethodGet}
	r.URL, _ = url.Parse("http://localhost/?name=matheus&age=-32&subscribed=true&ptr=4&float=4.2&dptr=64")

	a, b := 4, 64
	pa, pb := &a, &b
	ppb := &pb

	expected := Example{
		Name:       "matheus",
		Age:        -32,
		Subscribed: true,
		Ptr:        pa,
		Float:      4.2,
		DPtr:       ppb,
	}
	var actual ***Example

	err := Unmarshal(r, &actual)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(***actual, expected) {
		t.Errorf("%v not equal to %v", ***actual, expected)
	}
}
