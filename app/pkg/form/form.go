package form

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

type unmarshaller struct {
	err error
	r   *http.Request
}

func Parse(r *http.Request) *unmarshaller {
	if r == nil {
		return &unmarshaller{errors.New("missing request"), nil}
	}
	return &unmarshaller{nil, r}
}

func (u *unmarshaller) Into(dst interface{}) error {
	if u.err != nil {
		return u.err
	}

	t := reflect.TypeOf(dst)
	if t.Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}
	v := reflect.ValueOf(dst)

	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}
	v = v.Elem()

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)
		key := ft.Tag.Get("form")
		fmt.Println("tag:", key)
		if key == "" {
			continue
		}

		value := u.r.FormValue(key)

		fmt.Println(fv.IsValid(), fv.CanSet())

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(value)
		case reflect.Bool:
			fv.SetBool(value == "true")
		case reflect.Int:
			intValue, _ := strconv.Atoi(value)
			fv.SetInt(int64(intValue))
		}
	}

	return nil
}
