package gosql

import (
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/pkg/errors"
)

type SelectStmt struct {
	*stmt.SelectStmt
	inited bool
	// 是否對 SQL injection 做處理
	useAntiInjection bool
	queryFunc        func(datas [][]string, objs *[]any) error
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

func (s *SelectStmt) Query(objs *[]any) error {
	result, err := s.Exec()
	if err != nil {
		return errors.Wrap(err, "Failed to excute select statement.")
	}
	s.queryFunc(result.Datas, objs)
	return nil
}

func (s *SelectStmt) query(datas [][]string, objs *[]any) error {
	var obj any = objs
	results := (obj).(*[][]string)
	*results = append(*results, datas...)
	return nil
}

func (s *SelectStmt) SetFuncQuery(queryFunc func(datas [][]string, objs *[]any) error) {
	s.queryFunc = queryFunc
}
