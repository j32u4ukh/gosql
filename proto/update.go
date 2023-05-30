package proto

import (
	"reflect"

	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (t *ProtoTable) Update(pm protoreflect.ProtoMessage, where *gdo.WhereStmt) (string, error) {
	utils.Warn("package proto 即將棄用，請改用 package gosql")
	var rv, field reflect.Value
	var column *stmt.Column
	var i int
	rv = reflect.ValueOf(pm).Elem()

	// 遍歷每一欄位
	for i = 3; i < t.numFiled; i++ {
		field = rv.FieldByIndex([]int{i})
		// fmt.Printf("(t *ProtoTable) Update | field: %+v\n", field)
		column = t.Table.GetColumn(int32(i - 3))

		if column.IgnoreThis {
			continue
		}

		switch column.Default {
		// 有值也不更新
		// timestamp 類型可透過設置 OnUpdate 來更新時間戳
		case "current_timestamp()", "AI":
			continue

		default:
			t.Table.UpdateRawData(column.Name, gdo.ValueToDb(field, t.useAntiInjection, ProtoToDb))
		}
	}

	if where != nil {
		t.Table.SetUpdateCondition(where)
	}

	sql, err := t.BuildUpdateStmt()
	return sql, err
}
