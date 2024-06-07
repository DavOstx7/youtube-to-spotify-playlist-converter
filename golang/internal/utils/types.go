package utils

import (
	"golang/internal/common/types"
	"reflect"
)

func ReflectUnderlyingKind(input types.Any) reflect.Kind {
	value := reflect.ValueOf(input)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value.Kind()
}
