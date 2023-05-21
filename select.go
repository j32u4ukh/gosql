package gosql

import (
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/pkg/errors"
)

type SelectStmt struct {
	*stmt.SelectStmt
	db    *database.Database
	table *Table
	// 是否對 SQL injection 做處理
	useAntiInjection bool
}

func NewSelectStmt(table *Table) *SelectStmt {
	s := &SelectStmt{
		table:            table,
		SelectStmt:       stmt.NewSelectStmt(table.creater.TableName),
		db:               nil,
		useAntiInjection: false,
	}
	return s
}

func (s *SelectStmt) SetDb(db *database.Database) {
	s.db = db
}

func (s *SelectStmt) SetCondition(where *WhereStmt) {
	if s.useAntiInjection {
		where.UseAntiInjection()
	}
	s.Where = where.ToStmtWhere()
}

func (s *SelectStmt) UseAntiInjection(use bool) {
	s.useAntiInjection = use
}

func (s *SelectStmt) Exec() (*database.SqlResult, error) {
	if s.db == nil {
		return nil, errors.New("Undefine database.")
	}
	sql, err := s.SelectStmt.ToStmt()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate select statement.")
	}
	result, err := s.db.Query(sql)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excute select statement.")
	}
	return result, nil
}

func (s *SelectStmt) Query(datas *[][]string) error {
	result, err := s.Exec()
	if err != nil {
		return errors.Wrap(err, "Failed to excute select statement.")
	}
	*datas = append(*datas, result.Datas...)
	return nil
}
