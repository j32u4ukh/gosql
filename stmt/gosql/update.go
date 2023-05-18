package gosql

import (
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/pkg/errors"
)

type UpdateStmt struct {
	*stmt.UpdateStmt
	db *database.Database
}

func NewUpdateStmt(tableName string) *UpdateStmt {
	s := &UpdateStmt{
		UpdateStmt: stmt.NewUpdateStmt(tableName),
		db:         nil,
	}
	return s
}

func (s *UpdateStmt) SetDb(db *database.Database) {
	s.db = db
}

func (s *UpdateStmt) Exec() (*database.SqlResult, error) {
	if s.db == nil {
		return nil, errors.New("Undefine database.")
	}
	sql, err := s.UpdateStmt.ToStmt()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate update statement.")
	}
	result, err := s.db.Exec(sql)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excute update statement.")
	}
	return result, nil
}
