package utils

import "reflect"

func IsStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}
