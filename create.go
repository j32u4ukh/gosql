package gosql

import (
	"github.com/j32u4ukh/gosql/stmt"
)

type CreateStmt struct {
	*stmt.CreateStmt
}

func NewCreateStmt(tableName string, tableParam *stmt.TableParam, columnParams []*stmt.ColumnParam, engine string, collate string) *CreateStmt {
	s := &CreateStmt{
		CreateStmt: stmt.NewCreateStmt(tableName, tableParam, columnParams, engine, collate),
	}
	return s
}

func (s *CreateStmt) Clone() *CreateStmt {
	clone := &CreateStmt{
		CreateStmt: s.CreateStmt.Clone(),
	}
	return clone
}
