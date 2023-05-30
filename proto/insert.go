package proto

import (
	"reflect"

	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/utils"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func (t *ProtoTable) Insert(pms []protoreflect.ProtoMessage) error {
	utils.Warn("package proto 即將棄用，請改用 package gosql")
	var pm protoreflect.ProtoMessage

	var rv, field reflect.Value
	var column *stmt.Column
	var i, idx int

	for _, pm = range pms {
		rv = reflect.ValueOf(pm).Elem()
		data := []string{}

		for i = 3; i < t.numFiled; i++ {
			idx = i - 3
			column = t.Table.GetColumn(int32(idx))

			if column.IgnoreThis {
				continue
			}

			switch column.Default {
			// 資料庫自動生成欄位
			case "current_timestamp()", "AI":
				data = append(data, "NULL")
			default:
				field = rv.FieldByIndex([]int{i})
				data = append(data, gdo.ValueToDb(field, t.useAntiInjection, ProtoToDb))
			}
		}

		// 將數據加入 insert 緩存(傳入數據由 ProtoMessage 生成，因此所有欄位一定都在，且依照欄位順序)
		t.Table.InsertRawData(data...)
	}
	return nil
}
