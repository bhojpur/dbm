package utils

import (
	"reflect"
)

// ReflectValue returns value of a bean
func ReflectValue(bean interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(bean))
}
