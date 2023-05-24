package gosql

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/j32u4ukh/gosql/utils"
	"github.com/pkg/errors"
)

type TableConfig struct {
	Db               *database.Database
	DbName           string
	UseAntiInjection bool
	InsertFunc       func(data any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, toStringFunc func(v reflect.Value) string, insertFunc func(datas []string)) error
	QueryFunc        func(*database.SqlResult, *any) error
	UpdateAnyFunc    func(obj any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, updateFunc func(key string, field reflect.Value))
	PtrToDbFunc      func(reflect.Value, bool) string
}

type Table struct {
	creater    *CreateStmt
	insertPool *sync.Pool
	queryPool  *sync.Pool
	updatePool *sync.Pool
	deletePool *sync.Pool
	// 欄位名稱
	ColumnNames *cntr.Array[string]
	// 欄位數
	nColumn int32
	// 是否對 SQL injection 做處理
	useAntiInjection bool
	// ===== 處理函式備份 =====
	insertFunc    func(data any, nColumn int32, getColumnFunc func(idx int32) *stmt.Column, toStringFunc func(v reflect.Value) string, insertFunc func(datas []string)) error
	queryFunc     func(*database.SqlResult, *any) error
	updateAnyFunc func(any, int32, func(idx int32) *stmt.Column, func(key string, field reflect.Value))
	ptrToDbFunc   func(reflect.Value, bool) string
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Table
////////////////////////////////////////////////////////////////////////////////////////////////////
func NewTable(tableName string, tableParam *stmt.TableParam, columnParams []*stmt.ColumnParam, engine string, collate string, dial dialect.SQLDialect) *Table {
	t := &Table{
		ColumnNames:      cntr.NewArray[string](),
		creater:          NewCreateStmt(tableName, tableParam, columnParams, engine, collate),
		nColumn:          0,
		useAntiInjection: false,
		insertFunc:       nil,
		queryFunc:        nil,
	}
	t.insertPool = &sync.Pool{
		New: func() any {
			return NewInsertStmt(tableName, t.GetColumn)
		},
	}
	t.queryPool = &sync.Pool{
		New: func() any {
			return NewSelectStmt(tableName)
		},
	}
	t.updatePool = &sync.Pool{
		New: func() any {
			return NewUpdateStmt(tableName, t.GetColumn)
		},
	}
	t.deletePool = &sync.Pool{
		New: func() any {
			return NewDeleteStmt(tableName)
		},
	}
	t.nColumn = int32(len(t.creater.Columns))
	if t.nColumn > 0 {
		// 會自行賦值的欄位也需填入 NULL，因此所有欄位名稱都要求填入
		for _, column := range t.creater.Columns {
			if column.IgnoreThis {
				continue
			}
			utils.Debug("name: %s, default: %s", column.Name, column.Default)
			switch column.Default {
			// 資料庫自動生成欄位
			case "current_timestamp()", "AI":
				continue
			default:
				t.ColumnNames.Append(column.Name)
			}
		}
	}
	return t
}

func (t *Table) Init(config *TableConfig) {
	t.creater.SetDb(config.Db)
	t.creater.SetDbName(config.DbName)
	t.useAntiInjection = config.UseAntiInjection
	t.insertFunc = config.InsertFunc
	t.queryFunc = config.QueryFunc
	t.updateAnyFunc = config.UpdateAnyFunc
	t.ptrToDbFunc = config.PtrToDbFunc
}

func (t *Table) String() string {
	info := fmt.Sprintf("Table %s", t.creater.TableName)

	for i, col := range t.creater.Columns {
		info += fmt.Sprintf("\n%d) %s", i, col)
	}

	return info
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Creater
////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *Table) Creater() *CreateStmt {
	return t.creater
}

// 添加欄位
func (t *Table) AddColumn(column *stmt.Column) *Table {
	// 避免欄位重複添加
	if !t.ColumnNames.Contains(column.Name) {
		t.creater.AddColumn(column)
		t.ColumnNames.Append(column.Name)
		t.nColumn = int32(len(t.creater.Columns))
	}
	return t
}

func (t *Table) SetColumn(idx int32, column *stmt.Column) error {
	if t.nColumn-1 < idx {
		return errors.New(fmt.Sprintf("idx(%d) out of length(%d).", idx, t.nColumn-1))
	}

	origin := t.creater.Columns[idx]

	// 更新 tableParam 當中的欄位名稱(origin.Name -> column.Name)
	if origin.Name != column.Name {
		t.creater.TableParam.UpdateIndexName(origin.Name, column.Name)
	}

	t.creater.Columns[idx] = column
	return nil
}

func (t *Table) GetColumn(idx int32) *stmt.Column {
	if idx < int32(len(t.creater.Columns)) {
		return t.creater.Columns[idx].Clone()
	}
	return nil
}

func (t *Table) GetColumnByName(name string) *stmt.Column {
	for _, column := range t.creater.Columns {
		if column.Name == name {
			return column.Clone()
		}
	}
	return nil
}

func (t *Table) GetColumnNames() []string {
	names := []string{}
	names = append(names, t.ColumnNames.Elements...)
	return names
}

func (t *Table) GetIndexByName(name string) int32 {
	return int32(t.ColumnNames.Find(name))
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Insert
////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *Table) GetInserter() *InsertStmt {
	inserter := t.insertPool.Get().(*InsertStmt)
	if !inserter.inited {
		inserter.SetDb(t.creater.GetDb())
		inserter.SetDbName(t.creater.DbName)
		inserter.UseAntiInjection(t.useAntiInjection)
		inserter.SetColumnNumber(t.nColumn)
		inserter.SetFuncPtrToDb(t.ptrToDbFunc)
		if t.insertFunc != nil {
			inserter.SetFuncInsert(t.insertFunc)
		}
		inserter.inited = true
	}
	if inserter.ColumnStmt == "" {
		inserter.SetColumnNames(t.ColumnNames.Elements)
	}
	return inserter
}

func (t *Table) PutInserter(s *InsertStmt) {
	s.Release()
	t.insertPool.Put(s)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Select
////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Table) GetSelector() *SelectStmt {
	selector := t.queryPool.Get().(*SelectStmt)
	if !selector.inited {
		selector.SetDb(t.creater.GetDb())
		selector.SetDbName(t.creater.DbName)
		selector.UseAntiInjection(t.useAntiInjection)
		if t.queryFunc != nil {
			selector.SetFuncQuery(t.queryFunc)
		}
		selector.inited = true
	}
	return selector
}

func (t *Table) PutSelector(s *SelectStmt) {
	s.Release()
	t.queryPool.Put(s)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Update
////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Table) GetUpdater() *UpdateStmt {
	updater := t.updatePool.Get().(*UpdateStmt)
	if !updater.inited {
		updater.SetDb(t.creater.GetDb())
		updater.SetDbName(t.creater.DbName)
		updater.UseAntiInjection(t.useAntiInjection)
		updater.SetColumnNumber(t.nColumn)
		updater.SetFuncPtrToDb(t.ptrToDbFunc)
		if t.updateAnyFunc != nil {
			updater.SetFuncUpdateAny(t.updateAnyFunc)
		}
		updater.inited = true
	}
	return updater
}

func (t *Table) PutUpdater(s *UpdateStmt) {
	s.Release()
	t.updatePool.Put(s)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Delete
////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *Table) GetDeleter() *DeleteStmt {
	deleter := t.deletePool.Get().(*DeleteStmt)
	if !deleter.inited {
		deleter.SetDb(t.creater.GetDb())
		deleter.SetDbName(t.creater.DbName)
		deleter.UseAntiInjection(t.useAntiInjection)
		deleter.inited = true
	}
	return deleter
}

func (t *Table) PutDeleter(s *DeleteStmt) {
	s.Release()
	t.deletePool.Put(s)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Sync
////////////////////////////////////////////////////////////////////////////////////////////////////
// 根據傳入的欄位名稱列表 orders，對欄位進行重新排序
// return
// 	- 被修改順序的欄位名稱
func (t *Table) RefreshColumnOrder(orders []string) *cntr.Array[string] {
	changes := cntr.NewArray[string]()
	changed := t.refreshColumnOrder(orders)

	for changed != "" {
		changes.Append(changed)
		changed = t.refreshColumnOrder(orders)
	}
	return changes
}

func (t *Table) refreshColumnOrder(orders []string) string {
	changed := ""
	var i, j int
	var order string = ""
	var col *stmt.Column

	for i, order = range orders {
		for j, col = range t.creater.Columns {
			if col.Name == order {
				if i != j {
					changed = order
				}
				break
			}
		}
		if changed != "" {
			break
		}
	}
	if changed != "" {
		col = t.creater.Columns[i]
		t.creater.Columns[i] = t.creater.Columns[j]
		t.creater.Columns[j] = col
	}
	return changed
}

// NOTE: 根據 Sync 的需求，有需要再 Clone 即可
func (t *Table) SyncClone() *Table {
	clone := &Table{
		creater: t.creater.Clone(),
	}
	return clone
}
