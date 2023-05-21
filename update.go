package gosql

import (
	"reflect"

	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/pkg/errors"
)

type UpdateStmt struct {
	*stmt.UpdateStmt
	db               *database.Database
	table            *Table
	useAntiInjection bool
}

func NewUpdateStmt(table *Table) *UpdateStmt {
	s := &UpdateStmt{
		table:      table,
		UpdateStmt: stmt.NewUpdateStmt(table.creater.DbName),
		db:         nil,
	}
	return s
}

func (s *UpdateStmt) SetDb(db *database.Database) {
	s.db = db
}

func (s *UpdateStmt) SetCondition(where *WhereStmt) {
	s.Where = where.ToStmtWhere()
}

func (s *UpdateStmt) Update(key string, value any, ptrToDb func(reflect.Value, bool) string) {
	s.UpdateStmt.Update(key, ValueToDb(reflect.ValueOf(value), s.useAntiInjection, ptrToDb))
}

func (s *UpdateStmt) UseAntiInjection(use bool) {
	s.useAntiInjection = use
}

func (s *UpdateStmt) UpdateRawData(key string, value string) {
	s.UpdateStmt.Update(key, value)
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
