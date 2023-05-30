package gdo

import (
	"reflect"

	"github.com/j32u4ukh/gosql/utils"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// Update
////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Table) Update(key string, value any, ptrToDb func(reflect.Value, bool) string) {
	t.UpdateStmt.Update(key, ValueToDb(reflect.ValueOf(value), t.useAntiInjection, ptrToDb))
}

func (t *Table) UpdateRawData(key string, value string) {
	t.UpdateStmt.Update(key, value)
}

func (t *Table) SetUpdateCondition(where *WhereStmt) {
	if t.useAntiInjection {
		where.UseAntiInjection()
	}
	t.UpdateStmt.SetCondition(where.ToStmtWhere())
}

func (t *Table) AllowEmptyUpdateCondition() {
	t.UpdateStmt.AllowEmptyWhere()
}

func (t *Table) BuildUpdateStmt() (string, error) {
	utils.Warn("package gdo 即將棄用，請改用 package gosql")
	sql, err := t.UpdateStmt.ToStmt()
	t.UpdateStmt.Release()

	if err != nil {
		return "", errors.Wrap(err, "生成 UpdateStmt 時發生錯誤")
	}

	return sql, nil
}
