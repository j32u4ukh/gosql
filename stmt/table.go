package stmt

import (
	"fmt"
	"sync"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql/database"
)

type Table struct {
	*CreateStmt
	ColumnNames *cntr.Array[string]
	nColumn     int32

	insertPool *sync.Pool
	queryPool  *sync.Pool
	updatePool *sync.Pool
	deletePool *sync.Pool
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Table
////////////////////////////////////////////////////////////////////////////////////////////////////
func NewTable(name string, tableParam *TableParam, columnParams []*ColumnParam, engine string, collate string) *Table {
	t := &Table{
		CreateStmt:  NewCreateStmt(name, tableParam, columnParams, engine, collate),
		ColumnNames: cntr.NewArray[string](),
	}
	t.insertPool = &sync.Pool{
		New: func() any {
			return NewInsertStmt(name)
		},
	}
	t.queryPool = &sync.Pool{
		New: func() any {
			return NewSelectStmt(name)
		},
	}
	t.updatePool = &sync.Pool{
		New: func() any {
			return NewUpdateStmt(name)
		},
	}
	t.deletePool = &sync.Pool{
		New: func() any {
			return NewDeleteStmt(name)
		},
	}
	if len(t.CreateStmt.Columns) > 0 {
		// 會自行賦值的欄位也需填入 NULL，因此所有欄位名稱都要求填入
		for _, column := range t.CreateStmt.Columns {
			if column.IgnoreThis {
				continue
			}
			t.ColumnNames.Append(column.Name)
		}
	}
	return t
}

func (t *Table) SetDb(db *database.Database) {
	t.CreateStmt.db = db
	t.CreateStmt.DbName = db.DbName
}

func (t *Table) GetDbName() string {
	return t.CreateStmt.DbName
}

func (t *Table) GetTableName() string {
	return t.CreateStmt.TableName
}

func (t *Table) String() string {
	info := fmt.Sprintf("Table %s", t.CreateStmt.TableName)

	for i, col := range t.Columns {
		info += fmt.Sprintf("\n%d) %s", i, col)
	}

	return info
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Create
////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *Table) Creater() *CreateStmt {
	return t.CreateStmt
}

// 添加欄位
func (t *Table) AddColumn(column *Column) *Table {
	// 避免欄位重複添加
	if !column.IgnoreThis && !t.ColumnNames.Contains(column.Name) {
		t.CreateStmt.AddColumn(column)

		if column.Default != "current_timestamp()" {
			t.ColumnNames.Append(column.Name)
			t.nColumn = int32(t.ColumnNames.Length())
		}
	}
	return t
}

func (t *Table) BuildCreateStmt() (string, error) {
	return t.CreateStmt.ToStmt()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Insert
////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *Table) GetInserter() *InsertStmt {
	insert := t.insertPool.Get().(*InsertStmt)
	insert.SetDb(t.CreateStmt.db)
	insert.SetDbName(t.CreateStmt.DbName)
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
	selector.SetDb(t.CreateStmt.db)
	selector.SetDbName(t.CreateStmt.DbName)
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
	updater.SetDb(t.CreateStmt.db)
	updater.SetDbName(t.CreateStmt.DbName)
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
	deleter.SetDb(t.CreateStmt.db)
	deleter.SetDbName(t.CreateStmt.DbName)
	return deleter
}

func (t *Table) PutDeleter(s *DeleteStmt) {
	s.Release()
	t.deletePool.Put(s)
}
