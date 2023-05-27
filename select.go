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
	queryFunc        func(datas [][]string, generator func() any) (objs []any, err error)
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

func (s *SelectStmt) Query(generator func() any) (objs []any, err error) {
	result, err := s.Exec()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excute select statement.")
	}
	objs, err = s.queryFunc(result.Datas, generator)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excute select statement.")
	}
	return objs, nil
}

func (s *SelectStmt) query(datas [][]string, generator func() any) (objs []any, err error) {
	length := len(datas)
	objs = make([]any, length)
	var temp []string
	for i, data := range datas {
		temp = []string{}
		temp = append(temp, data...)
		objs[i] = temp
	}
	return objs, nil
}

func (s *SelectStmt) SetFuncQuery(queryFunc func(datas [][]string, generator func() any) (objs []any, err error)) {
	s.queryFunc = queryFunc
}
