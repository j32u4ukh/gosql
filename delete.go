package gosql

import (
	"github.com/j32u4ukh/gosql/stmt"
)

type DeleteStmt struct {
	*stmt.DeleteStmt
	useAntiInjection bool
	inited           bool
}

func NewDeleteStmt(tableName string) *DeleteStmt {
	s := &DeleteStmt{
		DeleteStmt:       stmt.NewDeleteStmt(tableName),
		useAntiInjection: false,
		inited:           false,
	}
	return s
}

func (s *DeleteStmt) SetCondition(where *WhereStmt) {
	if s.useAntiInjection {
		where.UseAntiInjection()
	}
	s.Where = where.ToStmtWhere()
}

func (s *DeleteStmt) UseAntiInjection(use bool) {
	s.useAntiInjection = use
}
