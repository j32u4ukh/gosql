package gdo

import (
	"github.com/j32u4ukh/gosql/utils"
	"github.com/pkg/errors"
)

func (t *Table) SetSelectCondition(where *WhereStmt) {
	if t.useAntiInjection {
		where.UseAntiInjection()
	}
	t.SelectStmt.SetCondition(where.ToStmtWhere())
}

func (t *Table) BuildSelectStmt() (string, error) {
	utils.Warn("package gdo 即將棄用，請改用 package gosql")
	sql, err := t.SelectStmt.ToStmt()
	t.SelectStmt.Release()
	if err != nil {
		return "", errors.Wrap(err, "生成 SelectStmt 時發生錯誤")
	}
	return sql, nil
}
