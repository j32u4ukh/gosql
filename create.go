package gosql

import (
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/pkg/errors"
)

type CreateStmt struct {
	*stmt.CreateStmt
	db *database.Database
}

func NewCreateStmt(name string, tableParam *stmt.TableParam, columnParams []*stmt.ColumnParam, engine string, collate string) *CreateStmt {
	s := &CreateStmt{
		CreateStmt: stmt.NewCreateStmt(name, tableParam, columnParams, engine, collate),
		db:         nil,
	}
	return s
}

func (s *CreateStmt) SetDb(db *database.Database) {
	s.db = db
}

func (s *CreateStmt) Exec() (*database.SqlResult, error) {
	if s.db == nil {
		return nil, errors.New("Undefine database.")
	}
	sql, err := s.CreateStmt.ToStmt()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate create statement.")
	}
	result, err := s.db.Exec(sql)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excute create statement.")
	}
	return result, nil
}

func (s *CreateStmt) Clone() *CreateStmt {
	clone := &CreateStmt{
		CreateStmt: s.CreateStmt.Clone(),
		db:         s.db,
	}
	return clone
}
