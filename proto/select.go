package proto

import (
	"reflect"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// 根據 ProtoMessage 生成 SQL 的查詢語法，m 須包含 table 的 primary key
func (t *ProtoTable) BuildSelectStmt(where *gdo.WhereStmt) (string, error) {
	if where != nil {
		t.Table.SetSelectCondition(where)
	}

	sql, err := t.Table.BuildSelectStmt()

	if err != nil {
		return "", err
	}

	return sql, nil
}

func (t *ProtoTable) ParseSelectResults(pms *[]protoreflect.ProtoMessage, results [][]string) error {
	var i, length int32 = 0, int32(len(results))
	pm := (*pms)[0]
	t.InitByProtoMessage(pm)
	var err error

	for i = 0; i < length; i++ {
		err = t.parseSelectResult((*pms)[i], results[i])
		if err != nil {
			return errors.Wrapf(err, "解析回傳數據時發生錯誤, result: %s", cntr.SliceToString(results[i]))
		}
	}
	return nil
}

func (t *ProtoTable) parseSelectResult(pm protoreflect.ProtoMessage, result []string) error {
	var filed reflect.Value
	rv := reflect.ValueOf(pm).Elem()

	for i, res := range result {
		if res == "" {
			continue
		}
		filed = rv.FieldByIndex([]int{i + 3})
		gdo.SetValue(filed, []byte(res), SetMessage)
	}
	return nil
}

// 取得符合 WhereStmt 條件的數據筆數
func (t *ProtoTable) CountStmt(where *gdo.WhereStmt) string {
	t.SetSelectItem(stmt.NewSelectItem("*").Count())
	sql, err := t.BuildSelectStmt(where)
	if err != nil {
		return ""
	}
	return sql
}
