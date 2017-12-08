package helper

import (
	"time"
	"reflect"
	"github.com/bmbstack/gron/xtime"
)

const (
	LogKeySource       string = "source"
	LogKeyTime         string = "time"
	LogKeyData         string = "data"
	LogKeyGoroutineNum string = "goroutineNum"

	TimeLocationName string = "Asia/Chongqing"

	DateShortLayout string = "2006-01-02"
	DateFullLayout  string = "2006-01-02 03:04:05"

	Line              string = "=============================="
	Line2             string = "==================================="

	JobEndTime = 1000*xtime.Day + 10*time.Millisecond

	SUCCESS_CODE_ZIROOM 		= 200
	SUCCESS_CODE_HIZHU 			= 200
	SUCCESS_CODE_XIANGYU 		= 200
	SUCCESS_CODE_ANJUKE 		= 200
	SUCCESS_CODE_QFANG 			= 200
	SUCCESS_CODE_FANGDUODUO 	= 200
)

func IsEmpty(i interface{}) bool {
	return isEmpty(reflect.ValueOf(i))
}

func IsNotEmpty(i interface{}) bool {
	return !isEmpty(reflect.ValueOf(i))
}

func isEmpty(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Bool:
		return v.Bool() == false

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return v.Float() == 0

	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0

	case reflect.Ptr, reflect.Interface:
		return isEmpty(v.Elem())

	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isEmpty(v.Index(i)) {
				return false
			}
		}
		return true

	case reflect.Slice, reflect.String, reflect.Map:
		return v.Len() == 0

	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			if !isEmpty(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return v.IsNil()
	}
}
