package proto

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/j32u4ukh/gosql/gdo"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	TIME_LAYOUT string = "2006-01-02 15:04:05"
)

func ProtoToDb(field reflect.Value, useAntiInjection bool) string {
	m := field.Interface().(protoreflect.ProtoMessage)
	bs, _ := json.Marshal(m)
	return string(bs)
}

func SetMessage(field reflect.Value, value []byte) {
	// fmt.Printf("SetMessage | field: %+v, value: %+v\n", field, value)
	rt := field.Type().Elem()
	msg := reflect.New(rt).Elem()

	switch rt.Name() {
	case "TimeStamp":
		SetTimeStamp(msg, value)
	default:
		m := map[string]any{}
		json.Unmarshal(value, &m)

		for i := 0; i < msg.NumField(); i++ {
			if i < 3 {
				continue
			}
			msgField := msg.Field(i)
			mt := msg.Type().Field(i)
			name, ok := GetTagName(mt.Tag.Get("protobuf"))

			if ok {
				// fmt.Printf("%d) name: %+v\n", i, name)

				if v, ok := m[name]; ok {
					var bs []byte
					switch t := v.(type) {
					case float64:
						// fmt.Printf("v.(type): float64\n")
						f64 := v.(float64)
						s64 := strconv.FormatFloat(f64, 'g', 5, 64)
						bs = []byte(s64)
					case string:
						fmt.Printf("SetMessage | v.(type): string\n")
						s := v.(string)
						bs = []byte(s)
					case map[string]any:
						fmt.Printf("SetMessage | v.(type): map[string]any\n")
					default:
						fmt.Printf("v.(type): %+v\n", t)
					}

					// kind := msgField.Kind()
					gdo.SetValue(msgField, bs, SetMessage)
				}
			}
		}
	}

	field.Set(msg.Addr())
}

func SetTimeStamp(field reflect.Value, value []byte) {
	s := string(value)
	t, _ := time.Parse(TIME_LAYOUT, s)
	// fmt.Printf("SetTimeStamp | t: %+v\n", t)

	for i := 0; i < field.NumField(); i++ {
		if i < 3 {
			continue
		}
		msgField := field.Field(i)

		switch i {
		case 3:
			msgField.SetInt(int64(t.Year()))
		case 4:
			msgField.SetInt(int64(t.Month()))
		case 5:
			msgField.SetInt(int64(t.Day()))
		case 6:
			msgField.SetInt(int64(t.Hour()))
		case 7:
			msgField.SetInt(int64(t.Minute()))
		case 8:
			msgField.SetInt(int64(t.Second()))
		}
	}
}

func GetTagName(tag string) (string, bool) {
	pairs := strings.Split(tag, ",")
	for _, pair := range pairs {
		k, v, ok := strings.Cut(pair, "=")
		if ok && k == "name" {
			return v, true
		}
	}
	return "", false
}
