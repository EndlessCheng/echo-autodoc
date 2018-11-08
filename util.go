package autodoc

import (
	"reflect"
	"time"
)

var (
	DefaultTime = time.Date(2010, 11, 22, 15, 16, 17, 0, time.Local)
	timeVal     = reflect.ValueOf(DefaultTime)
)

// 只能修改 public 字段，structPtr 不是 struct 指针时无效
func FillStruct(structPtr interface{}) {
	FillEmptyValuesOfStruct(reflect.ValueOf(structPtr), 1, 1, 1.0, "1")
}

// 只能修改 public 字段
func FillEmptyValuesOfStruct(structVal reflect.Value, intVal int64, uintVal uint64, floatVal float64, stringVal string) {
	if structVal.Kind() == reflect.Ptr {
		structVal = structVal.Elem()
	}
	if structVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < structVal.NumField(); i++ {
		structField := structVal.Field(i)
		structType := structField.Type()

		// 单独处理 time.Time
		if structType.ConvertibleTo(timeType) {
			if structField.CanSet() {
				structField.Set(timeVal.Convert(structType))
			}
			continue
		}

		switch structField.Kind() {
		case reflect.Invalid:
			panic("[FillEmptyValuesOfStruct] 暂不支持 " + structType.Name())
		case reflect.Bool:
			// ignore
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			if _v := structField.Int(); _v == 0 && structField.CanSet() {
				structField.SetInt(intVal)
			}
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			if _v := structField.Uint(); _v == 0 && structField.CanSet() {
				structField.SetUint(uintVal)
			}
		case reflect.Float32,
			reflect.Float64:
			if _v := structField.Float(); _v == 0 && structField.CanSet() {
				structField.SetFloat(floatVal)
			}
		case reflect.String:
			if _v := structField.String(); _v == "" && structField.CanSet() {
				structField.SetString(stringVal)
			}
		case reflect.Slice:
			// TODO
		case reflect.Struct:
			// TODO
		case reflect.Ptr:
			// TODO
		default:
			// ignore
		}
	}
}
