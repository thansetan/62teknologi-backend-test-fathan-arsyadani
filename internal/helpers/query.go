// !!!! error prone !!!!
package helpers

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func ParseQuery(q url.Values, target any) error {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Pointer {
		return errors.New("target is not a pointer")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("target is not a struct")
	}

	v := reflect.ValueOf(target).Elem()

	n := v.NumField()

	for i := 0; i < n; i++ {
		val := q.Get(t.Field(i).Tag.Get("query"))
		if val == "" {
			continue
		}
		err := fill(v.Field(i), val)
		if err != nil {
			return err
		}
	}

	return nil
}

func fill(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Bool:
		field.SetBool(stringToBool(value))
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Int8, reflect.Int16:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(n)
	case reflect.String:
		field.SetString(value)
	default:
		return errors.New(fmt.Sprintf("unsupported type: %T", field.Type()))
	}
	return nil
}

func stringToBool(val string) bool {
	if strings.ToLower(val) == "true" {
		return true
	}

	return false
}
