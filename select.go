package gosql

import (
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/pkg/errors"
)

type SelectStmt struct {
	*stmt.SelectStmt
	inited bool
	// 是否對 SQL injection 做處理
	useAntiInjection bool
	queryFunc        func(*database.SqlResult, *any) error
}

func NewSelectStmt(tableName string) *SelectStmt {
	s := &SelectStmt{
		SelectStmt:       stmt.NewSelectStmt(tableName),
		inited:           false,
		useAntiInjection: false,
	}
	s.queryFunc = s.query
	return s
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

func (s *SelectStmt) Query(datas *any) error {
	result, err := s.Exec()
	if err != nil {
		return errors.Wrap(err, "Failed to excute select statement.")
	}
	s.queryFunc(result, datas)
	return nil
}

func (s *SelectStmt) query(result *database.SqlResult, datas *any) error {
	results := (*datas).(*[][]string)
	*results = append(*results, result.Datas...)
	return nil
}

func (s *SelectStmt) SetFuncQuery(queryFunc func(*database.SqlResult, *any) error) {
	s.queryFunc = queryFunc
}
