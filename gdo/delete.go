package gdo

import (
	"github.com/j32u4ukh/gosql/utils"
	"github.com/pkg/errors"
)

func (t *Table) SetDeleteCondition(where *WhereStmt) {
	if t.useAntiInjection {
		where.UseAntiInjection()
	}
	t.DeleteStmt.SetCondition(where.ToStmtWhere())
}

func (t *Table) AllowEmptyDeleteCondition() {
	t.DeleteStmt.AllowEmptyWhere()
}

func (t *Table) BuildDeleteStmt() (string, error) {
	utils.Warn("package gdo 即將棄用，請改用 package gosql")
	sql, err := t.DeleteStmt.ToStmt()
	defer t.DeleteStmt.Release()

	if err != nil {
		return "", errors.Wrap(err, "生成 DeleteStmt 時發生錯誤")
	}

	return sql, nil
}
