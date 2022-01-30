package utils

import "reflect"

// New creates a value according type
func New(tp reflect.Type, length, cap int) reflect.Value {
	switch tp.Kind() {
	case reflect.Slice:
		slice := reflect.MakeSlice(tp, length, cap)
		x := reflect.New(slice.Type())
		x.Elem().Set(slice)
		return x
	case reflect.Map:
		mp := reflect.MakeMapWithSize(tp, cap)
		x := reflect.New(mp.Type())
		x.Elem().Set(mp)
		return x
	default:
		return reflect.New(tp)
	}
}
