package gosql

import (
	"fmt"
	"sync"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
)

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
	}
	t.insertPool = &sync.Pool{
		New: func() any {
			return NewInsertStmt(t)
		},
	}
	t.queryPool = &sync.Pool{
		New: func() any {
			return NewSelectStmt(t)
		},
	}
	t.updatePool = &sync.Pool{
		New: func() any {
			return NewUpdateStmt(t)
		},
	}
	t.deletePool = &sync.Pool{
		New: func() any {
			return NewDeleteStmt(t)
		},
	}
	if len(t.creater.Columns) > 0 {
		// 會自行賦值的欄位也需填入 NULL，因此所有欄位名稱都要求填入
		for _, column := range t.creater.Columns {
			if column.IgnoreThis {
				continue
			}
			t.ColumnNames.Append(column.Name)
		}
		t.nColumn = int32(t.ColumnNames.Length())
	}
	return t
}

func (t *Table) SetDbName(dbName string) {
	t.creater.SetDbName(dbName)
}

func (t *Table) SetDb(db *database.Database) {
	t.creater.SetDb(db)
}

func (t *Table) UseAntiInjection(active bool) {
	t.useAntiInjection = active
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
		t.nColumn = int32(t.ColumnNames.Length())
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
	insert := t.insertPool.Get().(*InsertStmt)
	insert.SetDb(t.creater.db)
	insert.SetDbName(t.creater.DbName)
	insert.UseAntiInjection(t.useAntiInjection)
	if insert.ColumnStmt == "" {
		insert.SetColumnNames(t.ColumnNames.Elements)
	}
	return insert
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
	selector.SetDb(t.creater.db)
	selector.SetDbName(t.creater.DbName)
	selector.UseAntiInjection(t.useAntiInjection)
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
	updater.SetDb(t.creater.db)
	updater.SetDbName(t.creater.DbName)
	updater.UseAntiInjection(t.useAntiInjection)
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
	deleter.SetDb(t.creater.db)
	deleter.SetDbName(t.creater.DbName)
	deleter.UseAntiInjection(t.useAntiInjection)
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
