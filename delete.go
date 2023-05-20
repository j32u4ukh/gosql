package gosql

import (
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/pkg/errors"
)

type DeleteStmt struct {
	*stmt.DeleteStmt
	db               *database.Database
	useAntiInjection bool
}

func NewDeleteStmt(tableName string) *DeleteStmt {
	s := &DeleteStmt{
		DeleteStmt:       stmt.NewDeleteStmt(tableName),
		db:               nil,
		useAntiInjection: false,
	}
	return s
}

func (s *DeleteStmt) SetDb(db *database.Database) {
	s.db = db
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

func (s *DeleteStmt) Exec() (*database.SqlResult, error) {
	if s.db == nil {
		return nil, errors.New("Undefine database.")
	}
	sql, err := s.DeleteStmt.ToStmt()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate delete statement.")
	}
	result, err := s.db.Exec(sql)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excute delete statement.")
	}
	return result, nil
}
