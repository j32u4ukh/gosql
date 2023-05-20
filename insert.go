package gosql

import (
	"fmt"
	"reflect"

	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/pkg/errors"
)

type InsertStmt struct {
	*stmt.InsertStmt
	db *database.Database
	// 欄位數
	nColumn          int32
	useAntiInjection bool
}

func NewInsertStmt(tableName string) *InsertStmt {
	s := &InsertStmt{
		InsertStmt:       stmt.NewInsertStmt(tableName),
		db:               nil,
		nColumn:          0,
		useAntiInjection: false,
	}
	return s
}

func (s *InsertStmt) SetDb(db *database.Database) {
	s.db = db
}

func (s *InsertStmt) Insert(datas []any, ptrToDb func(reflect.Value, bool) string) error {
	err := s.checkInsertData(int32(len(datas)))
	if err != nil {
		return errors.Wrap(err, "檢查輸入數據時發生錯誤")
	}
	var i int32
	insertDatas := []string{}
	for i = 0; i < s.nColumn; i++ {
		insertDatas = append(insertDatas, ValueToDb(reflect.ValueOf(datas[i]), s.useAntiInjection, ptrToDb))
	}
	s.InsertStmt.Insert(insertDatas)
	return nil
}

func (s *InsertStmt) InsertRawData(datas []string) error {
	err := s.checkInsertData(int32(len(datas)))
	if err != nil {
		return errors.Wrap(err, "檢查輸入數據時發生錯誤")
	}
	s.InsertStmt.Insert(datas)
	return nil
}

func (s *InsertStmt) checkInsertData(nData int32) error {
	// 確保 InsertStmt 有語法生成用的欄位名稱
	if s.ColumnStmt == "" {
		return errors.New("沒有有效的欄位名稱")
	}

	// 檢查輸入數據個數 與 欄位數 是否相符
	if nData != s.nColumn {
		return errors.New(fmt.Sprintf("輸入數據個數(%d)與欄位數(%d)不符", nData, s.nColumn))
	}
	return nil
}

func (s *InsertStmt) UseAntiInjection(use bool) {
	s.useAntiInjection = use
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
