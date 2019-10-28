package form

import (
	"net/http"
	"net/url"
	"testing"
)

type Example struct {
	Name       string `form:"name"`
	Age        int    `form:"age"`
	Subscribed bool   `form:"subscribed"`
}

func TestBasic(t *testing.T) {
	r := &http.Request{Method: http.MethodGet}
	r.URL, _ = url.Parse("http://localhost/?name=matheus&age=32&subscribed=false")

	expected := Example{"matheus", 32, false}
	actual := Example{}

	err := Parse(r).Into(&actual)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Errorf("%v not equal to %v", actual, expected)
	}
}
