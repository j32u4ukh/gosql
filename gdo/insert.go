package gdo

import (
	"fmt"
	"reflect"

	"github.com/j32u4ukh/gosql/utils"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// Insert
////////////////////////////////////////////////////////////////////////////////////////////////////
// 添加一筆數據(最終可同時添加多筆數據)
// 呼叫此函式者，須確保 datas 中的欄位都存在表格中
func (t *Table) Insert(datas []any, ptrToDb func(reflect.Value, bool) string) error {
	err := t.checkInsertData(int32(len(datas)))
	if err != nil {
		return errors.Wrap(err, "檢查輸入數據時發生錯誤")
	}

	var i int32
	insertDatas := []string{}
	for i = 0; i < t.nColumn; i++ {
		insertDatas = append(insertDatas, ValueToDb(reflect.ValueOf(datas[i]), t.useAntiInjection, ptrToDb))
	}

	t.InsertStmt.Insert(insertDatas)
	return nil
}

func (t *Table) InsertRawData(datas ...string) error {
	err := t.checkInsertData(int32(len(datas)))
	if err != nil {
		return errors.Wrap(err, "檢查輸入數據時發生錯誤")
	}
	t.InsertStmt.Insert(datas)
	return nil
}

func (t *Table) checkInsertData(nData int32) error {
	// 確保 InsertStmt 有語法生成用的欄位名稱
	if t.nColumn == 0 {
		t.nColumn = int32(t.ColumnNames.Length())
		t.SetColumnNames(t.ColumnNames.Elements)
	}

	// 檢查輸入數據個數 與 欄位數 是否相符
	if nData != t.nColumn {
		return errors.New(fmt.Sprintf("輸入數據個數(%d)與欄位數(%d)不符", nData, t.nColumn))
	}
	return nil
}

// 取得緩存數量
func (t *Table) GetInsertBufferNumber() int32 {
	return t.InsertStmt.GetBufferNumber()
}

func (t *Table) BuildInsertStmt() (string, error) {
	utils.Warn("package gdo 即將棄用，請改用 package gosql")
	sql, err := t.InsertStmt.ToStmt()
	t.InsertStmt.Release()
	if err != nil {
		return "", errors.Wrap(err, "生成 InsertStmt 時發生錯誤")
	}
	return sql, nil
}
