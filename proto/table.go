package proto

import (
	"reflect"

	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ProtoTable 負責將 protobuf 的數據轉換成字串，實際組合與取得語法，依然是透過 stmt.Table
type ProtoTable struct {
	*gdo.Table
	// 定義是否使用反注入檢查
	useAntiInjection bool
	// 通用反射用變數
	numFiled int
}

func NewTable(name string, dial dialect.SQLDialect) *ProtoTable {
	t := &ProtoTable{
		Table: gdo.NewTable(name, stmt.NewTableParam(), nil, stmt.ENGINE, stmt.COLLATE, dial),
	}
	return t
}

// 根據傳入的 Param，對 Column 進行定義，並生成 Table
func NewProtoTable(name string, tableParam *stmt.TableParam, params []*stmt.ColumnParam, dial dialect.SQLDialect) *ProtoTable {
	t := &ProtoTable{
		Table:            gdo.NewTable(name, tableParam, params, stmt.ENGINE, stmt.COLLATE, dial),
		useAntiInjection: true,
		numFiled:         -1,
	}
	return t
}

func (t *ProtoTable) UseAntiInjection(active bool) {
	t.useAntiInjection = active
}

func (t *ProtoTable) InitByProtoMessage(pm protoreflect.ProtoMessage) {
	if t.numFiled == -1 {
		rt := reflect.TypeOf(pm).Elem()
		t.numFiled = rt.NumField()

		// 設置更新用欄位名稱列表
		t.Table.SetColumnNames(t.Table.ColumnNames.Elements)
	}
}
