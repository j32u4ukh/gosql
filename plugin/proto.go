package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
	protoparser "github.com/yoheimuta/go-protoparser/v4"
	ud "github.com/yoheimuta/go-protoparser/v4/interpret/unordered"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// NOTE: 參考 https://github.com/yoheimuta/go-protoparser
func GetProtoParams(path string, sqlDial dialect.SQLDialect) (*stmt.TableParam, []*stmt.ColumnParam, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Failed to open %s.\n", path))
	}
	defer reader.Close()
	got, err := protoparser.Parse(
		reader,
		protoparser.WithFilename(filepath.Base(path)),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Failed to parse %s.\n", path))
	}
	// 印出抽象語法樹整體結構
	// printStructure(got)
	var unorder *ud.Proto
	var tag string
	var idx int
	unorder, err = protoparser.UnorderedInterpret(got)

	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to execute UnorderedInterpret.")
	}
	msg := unorder.ProtoBody.Messages[0]
	tableParam := stmt.NewTableParam()
	for _, comment := range msg.Comments {
		tag = comment.Raw
		idx = strings.Index(tag, "//")
		tag = tag[idx+2:]
		tag = strings.Trim(tag, " ")
		tableParam.ParserConfig(tag)
	}
	var param *stmt.ColumnParam
	var tags []string
	body := msg.MessageBody
	// 以 stmt.ColumnParamSlice 封裝 []*stmt.ColumnParam{}，以實作排序介面
	var colParams stmt.ColumnParamSlice
	colParams = []*stmt.ColumnParam{}
	// fmt.Printf("GetProtoParams | Messages: %+v\n", body.Messages)
	for _, dict := range body.Maps {
		tags = []string{}

		if len(dict.Comments) > 0 {
			for _, comment := range dict.Comments {
				tag = comment.Raw
				idx = strings.Index(tag, "//")
				tag = tag[idx+2:]
				tag = strings.Trim(tag, " ")
				tags = append(tags, tag)
			}
		}

		// fmt.Printf("GetProtoParams | dict: %+v\n", dict)
		idx, _ = strconv.Atoi(dict.FieldNumber)
		param = stmt.NewColumnParam(
			idx,
			dict.MapName,
			datatype.MAP,
			sqlDial,
			tags...,
		)
		colParams = append(colParams, param)
	}

	for _, filed := range body.Fields {
		tags = []string{}

		if len(filed.Comments) > 0 {
			for _, comment := range filed.Comments {
				tag = comment.Raw
				idx = strings.Index(tag, "//")
				tag = tag[idx+2:]
				tag = strings.Trim(tag, " ")
				tags = append(tags, tag)
			}
		}
		idx, _ = strconv.Atoi(filed.FieldNumber)
		param = stmt.NewColumnParam(
			idx,
			filed.FieldName,
			datatype.ProtoToDataType(filed.Type),
			sqlDial,
			tags...,
		)
		colParams = append(colParams, param)
	}

	colParams.Sort()
	return tableParam, colParams, nil
}

// 印出抽象語法樹整體結構
func printProtoStructure(got any) {
	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal, err %v\n", err)
	}
	fmt.Printf("printStructure | gotJSON:\n%s\n", gotJSON)
}

func InsertProto(data any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, toStringFunc func(v reflect.Value) string, insertFunc func(datas []string)) error {
	var rv, field reflect.Value
	var column *stmt.Column
	var i, idx int32

	rv = reflect.ValueOf(data).Elem()
	values := []string{}

	// Column 的資訊是根據 proto 檔生成的，因此沒有 struct 當中被跳過的三個欄位
	// 因此需要校正兩者的欄位數量差異
	nColumn += 3

	// 跳過前三個不需要的欄位
	for i = 3; i < nColumn; i++ {
		idx = i - 3
		column = getColumnFunc(idx)

		if column.IgnoreThis {
			continue
		}

		switch column.Default {
		// 資料庫自動生成欄位
		case "current_timestamp()", "AI":
			continue
		default:
			field = rv.FieldByIndex([]int{int(i)})
			values = append(values, toStringFunc(field))
		}
	}

	// 將一筆數據加入 insert 緩存(數據來自 struct，所有欄位一定都有，不須再做檢查)
	insertFunc(values)
	return nil
}

func QueryProto(columns []string, datas [][]string, generator func() any) (objs []any, err error) {
	var i, length int32 = 0, int32(len(datas))
	objs = make([]any, length)
	if length == 0 {
		return objs, nil
	}
	var obj any
	for i = 0; i < length; i++ {
		obj, err = queryProto(columns, datas[i], generator)
		if err != nil {
			return nil, errors.Wrapf(err, "解析回傳數據時發生錯誤, result: %s", cntr.SliceToString(datas[i]))
		}
		objs[i] = obj
	}
	return objs, nil
}

func queryProto(columns []string, data []string, generator func() any) (obj any, err error) {
	var field reflect.Value
	var column, d string
	var i int
	var ok bool
	obj = generator()
	rv := reflect.ValueOf(obj).Elem()
	rt := rv.Type()
	length := rt.NumField()
	for i, column = range columns {
		d = data[i]
		if d == "" {
			continue
		}
		if field, ok = getField(rv, rt, length, column); ok {
			SetValue(field, d, SetMessage)
		}
	}
	return obj, nil
}

func getField(rv reflect.Value, rt reflect.Type, length int, name string) (reflect.Value, bool) {
	var field reflect.Value
	var structField reflect.StructField
	var splits []string
	for i := 0; i < length; i++ {
		structField = rt.Field(i)
		if value, ok := structField.Tag.Lookup("json"); ok {
			splits = strings.Split(value, ",")
			for _, s := range splits {
				if s == name {
					field = rv.Field(i)
					return field, true
				}
			}
		}
	}
	return field, false
}

func UpdateProto(data any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, updateFunc func(key string, field reflect.Value)) {
	var rv, field reflect.Value
	var column *stmt.Column
	var i, idx int32
	rv = reflect.ValueOf(data).Elem()

	// Column 的資訊是根據 proto 檔生成的，因此沒有 struct 當中被跳過的三個欄位
	// 因此需要校正兩者的欄位數量差異
	nColumn += 3

	// 遍歷每一欄位
	for i = 3; i < nColumn; i++ {
		field = rv.FieldByIndex([]int{int(i)})
		// fmt.Printf("(t *ProtoTable) Update | field: %+v\n", field)
		idx = i - 3
		column = getColumnFunc(idx)

		if column.IgnoreThis {
			continue
		}

		switch column.Default {
		// 有值也不更新
		// timestamp 類型可透過設置 OnUpdate 來更新時間戳
		case "current_timestamp()", "AI":
			continue

		default:
			updateFunc(column.Name, field)
		}
	}
}

func ProtoToDb(field reflect.Value, useAntiInjection bool) string {
	m := field.Interface().(protoreflect.ProtoMessage)
	bs, _ := json.Marshal(m)
	return string(bs)
}

func SetMessage(field reflect.Value, value string) {
	bs := []byte(value)
	// fmt.Printf("SetMessage | field: %+v, value: %+v\n", field, value)
	rt := field.Type().Elem()
	msg := reflect.New(rt).Elem()

	switch rt.Name() {
	case "TimeStamp":
		SetTimeStamp(msg, value)
	default:
		m := map[string]any{}
		json.Unmarshal(bs, &m)

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
					var s string
					switch t := v.(type) {
					case float64:
						// fmt.Printf("v.(type): float64\n")
						f64 := v.(float64)
						s = strconv.FormatFloat(f64, 'g', 5, 64)
					case string:
						fmt.Printf("SetMessage | v.(type): string\n")
						s = v.(string)
					case map[string]any:
						fmt.Printf("SetMessage | v.(type): map[string]any\n")
					default:
						fmt.Printf("v.(type): %+v\n", t)
					}

					// kind := msgField.Kind()
					SetValue(msgField, s, SetMessage)
				}
			}
		}
	}

	field.Set(msg.Addr())
}

func SetTimeStamp(field reflect.Value, value string) {
	t, _ := time.Parse(TIME_LAYOUT, value)
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
