package gocon

import (
	"reflect"
)

func Value[T any](v T) *Definition {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	return &Definition{
		Key:   keyOf(rt),
		Type:  rt,
		Value: &rv,
	}
}
