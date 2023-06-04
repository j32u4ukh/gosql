package gosql

import (
	"reflect"

	"github.com/j32u4ukh/gosql/plugin"
	"github.com/j32u4ukh/gosql/stmt"
)

type InsertStmt struct {
	*stmt.InsertStmt
	// 欄位數
	nColumn          int32
	useAntiInjection bool
	inited           bool
	// table 提供的函式
	getColumnFunc func(idx int32) *stmt.Column
	// 不同數據結構各自定義
	insertFunc  func(data any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, toStringFunc func(v reflect.Value) string, insertFunc func(datas []string)) error
	ptrToDbFunc func(reflect.Value, bool) string
}

func NewInsertStmt(tableName string, getColumnFunc func(idx int32) *stmt.Column) *InsertStmt {
	s := &InsertStmt{
		InsertStmt:       stmt.NewInsertStmt(tableName),
		nColumn:          0,
		useAntiInjection: false,
		inited:           false,
		getColumnFunc:    getColumnFunc,
		ptrToDbFunc:      nil,
	}
	s.insertFunc = s.insert
	return s
}

func (s *InsertStmt) Insert(data any) error {
	return s.insertFunc(data, s.nColumn, s.getColumnFunc, s.toStringFunc, s.InsertRawData)
}

func (s *InsertStmt) insert(data any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, toStringFunc func(v reflect.Value) string, insertFunc func(datas []string)) error {
	datas := data.([]any)
	strData := []string{}
	var d any
	for _, d = range datas {
		strData = append(strData, toStringFunc(reflect.ValueOf(d)))
	}
	insertFunc(strData)
	return nil
}

func (s *InsertStmt) InsertRawData(datas []string) {
	s.InsertStmt.Insert(datas)
}

func (s *InsertStmt) SetColumnNumber(nColumn int32) {
	s.nColumn = nColumn
}

func (s *InsertStmt) UseAntiInjection(use bool) {
	s.useAntiInjection = use
}

func (s *InsertStmt) SetFuncInsert(insertFunc func(data any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, toStringFunc func(v reflect.Value) string, insertFunc func(datas []string)) error) {
	s.insertFunc = insertFunc
}

func (s *InsertStmt) SetFuncPtrToDb(ptrToDbFunc func(reflect.Value, bool) string) {
	s.ptrToDbFunc = ptrToDbFunc
}

func (s *InsertStmt) toStringFunc(v reflect.Value) string {
	return plugin.ValueToDb(v, s.useAntiInjection, s.ptrToDbFunc)
}
