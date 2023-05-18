package gosql

import (
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/pkg/errors"
)

type InsertStmt struct {
	*stmt.InsertStmt
	db *database.Database
}

func NewInsertStmt(tableName string) *InsertStmt {
	s := &InsertStmt{
		InsertStmt: stmt.NewInsertStmt(tableName),
		db:         nil,
	}
	return s
}

func (s *InsertStmt) SetDb(db *database.Database) {
	s.db = db
}

func (s *InsertStmt) Exec() (*database.SqlResult, error) {
	if s.db == nil {
		return nil, errors.New("Undefine database.")
	}
	sql, err := s.InsertStmt.ToStmt()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate insert statement.")
	}
	result, err := s.db.Exec(sql)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excute insert statement.")
	}
	return result, nil
}
