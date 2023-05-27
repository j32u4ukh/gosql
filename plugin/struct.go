package plugin

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
)

// 讀取 Struct 的 tag
func GetStructParams(data any, dial dialect.SQLDialect) (*stmt.TableParam, []*stmt.ColumnParam, error) {
	rv := reflect.ValueOf(data)
	var rt reflect.Type

	if rv.Kind() == reflect.Ptr {
		rt = rv.Elem().Type()
	} else {
		rt = reflect.TypeOf(data)
	}

	tableParam := stmt.NewTableParam()
	columnParams := []*stmt.ColumnParam{}
	var tpc *stmt.TableParamConfig
	var cpc *stmt.ColumnParamConfig

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		columnParam := stmt.NewColumnParam(i, field.Name, datatype.DataType(field.Type.Kind().String()), dial)
		config, ok := field.Tag.Lookup("gorm")

		if ok {
			tpc, cpc, _ = stmt.ParseConfig(field.Name, config)
			tableParam.LoadConfig(tpc)
			columnParam.LoadConfig(cpc)
		} else {
			fmt.Printf("No tag found for field %s\n", field.Name)
		}

		columnParam.Redefine()
		columnParams = append(columnParams, columnParam)
	}

	return tableParam, columnParams, nil
}

func InsertStruct(data any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, toStringFunc func(v reflect.Value) string, insertFunc func(datas []string)) error {
	var rv, field reflect.Value
	var column *stmt.Column
	var i int32

	rv = reflect.ValueOf(data).Elem()
	values := []string{}

	for i = 0; i < nColumn; i++ {
		column = getColumnFunc(i)

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

func QueryStructFunc(datas [][]string, generator func() any) (objs []any, err error) {
	var i, length int32 = 0, int32(len(datas))
	objs = make([]any, length)
	var obj any
	for i = 0; i < length; i++ {
		obj, err = queryStructFunc(datas[i], generator)
		if err != nil {
			return nil, errors.Wrapf(err, "解析回傳數據時發生錯誤, result: %s", cntr.SliceToString(datas[i]))
		}
		objs[i] = obj
	}
	return objs, nil
}

func queryStructFunc(data []string, generator func() any) (obj any, err error) {
	var filed reflect.Value
	obj = generator()
	rv := reflect.ValueOf(obj).Elem()
	for i, d := range data {
		if d == "" {
			continue
		}
		filed = rv.FieldByIndex([]int{i})
		SetValue(filed, d, nil)
	}
	return obj, nil
}

func UpdateStruct(data any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, updateFunc func(key string, field reflect.Value)) {
	var rv, field reflect.Value
	var column *stmt.Column
	var i int32
	rv = reflect.ValueOf(data).Elem()

	// 遍歷每一欄位
	for i = 0; i < nColumn; i++ {
		field = rv.FieldByIndex([]int{int(i)})
		// fmt.Printf("(t *ProtoTable) Update | field: %+v\n", field)
		column = getColumnFunc(i)

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

func StructToDb(field reflect.Value, useAntiInjection bool) string {
	m := field.Interface().(ISqlStruct)
	return m.ToStmt()
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
