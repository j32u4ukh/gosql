package proto

import (
	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/utils"

	"github.com/pkg/errors"
)

func (t *ProtoTable) BuildDeleteStmt(where *gdo.WhereStmt) (string, error) {
	utils.Warn("package proto 即將棄用，請改用 package gosql")
	if where != nil {
		t.Table.SetDeleteCondition(where)
	}

	sql, err := t.Table.BuildDeleteStmt()

	if err != nil {
		return "", errors.Wrap(err, "Failed to build DeleteStmt.")
	}

	return sql, nil
}
