package form

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

var (
	NilRequestErr  = errors.New("nil request")
	InvalidTypeErr = errors.New("invalid type")
)

type FieldUnmarshaler interface {
	UnmarshalFormField(string) error
}

type decodeState struct {
	r   *http.Request
	err error
}

func (d *decodeState) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return InvalidTypeErr
	}

	_, rv = deref(rv)
	if rv.Kind() != reflect.Struct {
		return InvalidTypeErr
	}

	for i := 0; i < rv.NumField(); i++ {
		key := rv.Type().Field(i).Tag.Get("form")
		if len(key) == 0 {
			continue
		}

		value := d.r.FormValue(key)
		if len(value) == 0 {
			continue
		}

		u, pv := deref(rv.Field(i))
		if u != nil {
			return u.UnmarshalFormField(value)
		}

		if err := set(pv, value); err != nil {
			return err
		}
	}

	return nil
}

func set(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Bool:
		v.SetBool(value == "true")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil || v.OverflowInt(n) {
			return errors.New("error parsing int")
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil || v.OverflowUint(n) {
			return errors.New("error parsing uint")
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(value, v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			return errors.New("error parsing float")
		}
		v.SetFloat(n)
	case reflect.Interface:
		if v.NumMethod() != 0 {
			return errors.New("error parsing interface")
		}
		v.SetString(value)
	default:
		panic(fmt.Sprintf("unhandled type: %s", v.Type().Name()))
	}

	return nil
}

// derefs pointer until finds an FieldUnmarshaler or concrete type
func deref(v reflect.Value) (FieldUnmarshaler, reflect.Value) {
	for {
		if v.Kind() == reflect.Interface {
			fmt.Printf("an interface. dont know what to do yet")
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		if v.Type().NumMethod() > 0 && v.CanInterface() {
			if u, ok := v.Interface().(FieldUnmarshaler); ok {
				return u, reflect.Value{}
			}
		}

		v = v.Elem()
	}

	return nil, v
}

func Unmarshal(r *http.Request, v interface{}) error {
	if r == nil {
		return NilRequestErr
	}

	d := decodeState{r, nil}

	return d.unmarshal(v)
}
