package plugin

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/j32u4ukh/gosql/stmt"
)

type ISqlStruct interface {
	ToStmt() string
}

func Make(size int, gen func() any) []any {
	arr := make([]any, size)
	for i := 0; i < size; i++ {
		arr[i] = gen()
	}
	return arr
}

type ValueToDbFunc func(v reflect.Value, useAntiInjection bool, ptrToDb func(reflect.Value, bool) string) string

func ValueToDb(v reflect.Value, useAntiInjection bool, ptrToDb func(reflect.Value, bool) string) string {
	kind := v.Kind()
	// fmt.Printf("ValueToDb | kind: %s\n", kind)

	switch kind {
	case reflect.Bool:
		if v.Bool() {
			return "1"
		} else {
			return "0"
		}
	case reflect.String:
		s := v.String()
		switch s {
		case "current_timestamp()", "NULL":
			return s
		case "NIL":
			return ""
		default:
			if useAntiInjection {
				return fmt.Sprintf("'%s'", stmt.AntiInjectionString(s))
			} else {
				return fmt.Sprintf("'%s'", s)
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", v.Float())
	case reflect.Slice, reflect.Array:
		return fmt.Sprintf("%v", v)
	case reflect.Map:
		bs, _ := json.Marshal(v.Interface())
		if useAntiInjection {
			return fmt.Sprintf("'%s'", stmt.AntiInjectionString(string(bs)))
		} else {
			return fmt.Sprintf("'%s'", string(bs))
		}
	case reflect.Pointer:
		if ptrToDb != nil {
			return fmt.Sprintf("'%s'", ptrToDb(v, useAntiInjection))
		}
		return "''"
	default:
		return ""
	}
}

func SetValue(field reflect.Value, value string, setPointer func(reflect.Value, string)) {
	kind := field.Kind()
	// fmt.Printf("SetValue(kind: %v, field: %v, value: %v)\n", kind, field, value)

	switch kind {
	case reflect.Bool:
		field.SetBool(value == "1")
	case reflect.Uint32, reflect.Uint64:
		v, _ := strconv.ParseUint(value, 10, 64)
		field.SetUint(uint64(v))
	case reflect.Int, reflect.Int32, reflect.Int64:
		v, _ := strconv.ParseInt(value, 10, 64)
		field.SetInt(int64(v))
	case reflect.Float32:
		v, _ := strconv.ParseFloat(value, 32)
		field.SetFloat(v)
	case reflect.Float64:
		v, _ := strconv.ParseFloat(value, 64)
		field.SetFloat(v)
	case reflect.String:
		field.SetString(value)
	case reflect.Map:
		SetMap(field, []byte(value))
	case reflect.Pointer:
		if setPointer != nil {
			setPointer(field, value)
		}
	}
	// fmt.Printf("Set Value field: %v\n", field)
}

func SetMap(filed reflect.Value, value []byte) {
	m := map[string]any{}
	json.Unmarshal(value, &m)
	keyType := filed.Type().Key()
	valueType := filed.Type().Elem()
	rm := reflect.MapOf(keyType, valueType)
	// fmt.Printf("rm: %+v\n", rm)
	filed.Set(reflect.MakeMap(rm))
	// fmt.Printf("map filed: %+v\n", filed)

	for k, v := range m {
		filed.SetMapIndex(MapKey(keyType.Kind(), k), MapValue(valueType.Kind(), v))
	}
}

func MapKey(kind reflect.Kind, key string) reflect.Value {
	switch kind {
	case reflect.Int:
		v, _ := strconv.Atoi(key)
		return reflect.ValueOf(v)
	case reflect.Int32:
		v, _ := strconv.Atoi(key)
		return reflect.ValueOf(int32(v))
	case reflect.Int64:
		v, _ := strconv.Atoi(key)
		return reflect.ValueOf(int64(v))
	default:
		return reflect.ValueOf(key)
	}
}

func MapValue(kind reflect.Kind, value any) reflect.Value {
	switch vt := value.(type) {
	case float64:
		v := value.(float64)
		switch kind {
		case reflect.Int32:
			return reflect.ValueOf(int32(v))
		case reflect.Int64:
			return reflect.ValueOf(int64(v))
		default:
			fmt.Printf("MapValue | vt: %v\n", vt)
			return reflect.ValueOf(int(v))
		}
	default:
		return reflect.ValueOf(value)
	}
}
