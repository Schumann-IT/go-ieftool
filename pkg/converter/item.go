package converter

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var (
	nilAbleKinds = []reflect.Kind{
		reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice,
	}
	lenKinds = []reflect.Kind{
		reflect.Array, reflect.Slice, reflect.Map,
	}
)

type (
	item struct {
		from interface{}
		to   interface{}
	}
)

func newItem(from, to interface{}) (*item, error) {
	if !isNonNilPointer(reflect.ValueOf(to)) {
		return nil, fmt.Errorf("to value needs to be a non nil pointer")
	}

	if isNil(reflect.ValueOf(from)) {
		return nil, fmt.Errorf("from value cannot be nil")
	}

	if !hasLen(reflect.ValueOf(from)) {
		return nil, fmt.Errorf("from value of type %s type has len = 0", reflect.TypeOf(from).Kind().String())
	}

	return &item{
		from: from,
		to:   to,
	}, nil

}

func isNil(v reflect.Value) bool {
	for i := range nilAbleKinds {
		if nilAbleKinds[i] == v.Kind() {
			return v.IsNil()
		}
	}

	return false
}

func hasLen(v reflect.Value) bool {
	for i := range lenKinds {
		if lenKinds[i] == v.Kind() {
			return v.Len() > 0
		}
	}

	return true
}

func isNonNilPointer(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && !isNil(v)
}

func (c *item) convert() error {
	j, err := json.Marshal(c.from)
	if err != nil {
		return fmt.Errorf("cannot convert from %v: %s", c.from, err.Error())
	}

	err = json.Unmarshal(j, c.to)
	if err != nil {
		return fmt.Errorf("cannot convert from %v to %v: %s", c.from, c.to, err.Error())
	}

	return nil
}
