package gosql

import (
	"reflect"

	"github.com/j32u4ukh/gosql/stmt"
)

type UpdateStmt struct {
	*stmt.UpdateStmt
	useAntiInjection bool
	inited           bool
	nColumn          int32
	ptrToDbFunc      func(reflect.Value, bool) string
	updateAnyFunc    func(obj any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, updateFunc func(key string, field reflect.Value))
	// table 提供的函式
	getColumnFunc func(idx int32) *stmt.Column
}

func NewUpdateStmt(tableName string, getColumnFunc func(idx int32) *stmt.Column) *UpdateStmt {
	s := &UpdateStmt{
		UpdateStmt:    stmt.NewUpdateStmt(tableName),
		inited:        false,
		nColumn:       0,
		ptrToDbFunc:   nil,
		updateAnyFunc: nil,
		getColumnFunc: getColumnFunc,
	}
	return s
}

func (s *UpdateStmt) SetColumnNumber(nColumn int32) {
	s.nColumn = nColumn
}

func (s *UpdateStmt) GetColumnNumber() int32 {
	return s.nColumn
}

func (s *UpdateStmt) SetCondition(where *WhereStmt) {
	s.Where = where.ToStmtWhere()
}

func (s *UpdateStmt) Update(key string, value any) {
	s.UpdateStmt.Update(key, ValueToDb(reflect.ValueOf(value), s.useAntiInjection, s.ptrToDbFunc))
}

func (s *UpdateStmt) UpdateAny(obj any) {
	s.updateAnyFunc(obj, s.nColumn, s.getColumnFunc, s.updateField)
}

func (s *UpdateStmt) updateField(key string, field reflect.Value) {
	s.UpdateRawData(key, ValueToDb(field, s.useAntiInjection, s.ptrToDbFunc))
}

func (s *UpdateStmt) UseAntiInjection(use bool) {
	s.useAntiInjection = use
}

func (s *UpdateStmt) UpdateRawData(key string, value string) {
	s.UpdateStmt.Update(key, value)
}

func (s *UpdateStmt) SetFuncPtrToDb(ptrToDbFunc func(reflect.Value, bool) string) {
	s.ptrToDbFunc = ptrToDbFunc
}

func (s *UpdateStmt) SetFuncUpdateAny(updateAnyFunc func(obj any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, updateFunc func(key string, field reflect.Value))) {
	s.updateAnyFunc = updateAnyFunc
}
