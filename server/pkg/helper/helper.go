package helper

import (
	"reflect"
	"strconv"
)

func Contains[T comparable](elems []T, v T, fn func(value T, element T) bool) bool {
	if fn == nil {
		fn = func(value T, element T) bool {
			return value == element
		}
	}
	for _, s := range elems {
		if fn(v, s) {
			return true
		}
	}
	return false
}

func ConvertToFloat64(value interface{}) float64 {
	if value == nil {
		return 0
	}

	typeOf := reflect.TypeOf(value).String()

	switch typeOf {
	case "float64":
		return value.(float64)
	case "float32":
		return float64(value.(float32))
	case "int":
		return float64(value.(int))
	case "int32":
		return float64(value.(int32))
	case "int64":
		return float64(value.(int64))
	case "string":
		f, _ := strconv.ParseFloat(value.(string), 64)
		return f
	default:
		return 0
	}
}
