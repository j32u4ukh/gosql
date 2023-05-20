package gosql

import (
	"fmt"
	"sync"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
)

type Table struct {
	ColumnNames *cntr.Array[string]
	db          *database.Database
	creater     *CreateStmt
	insertPool  *sync.Pool
	queryPool   *sync.Pool
	updatePool  *sync.Pool
	deletePool  *sync.Pool
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Table
////////////////////////////////////////////////////////////////////////////////////////////////////
func NewTable(tableName string, tableParam *stmt.TableParam, columnParams []*stmt.ColumnParam, engine string, collate string) *Table {
	t := &Table{
		ColumnNames: cntr.NewArray[string](),
		db:          nil,
		creater:     NewCreateStmt(tableName, tableParam, columnParams, engine, collate),
		insertPool: &sync.Pool{
			New: func() any {
				return NewInsertStmt(tableName)
			},
		},
		queryPool: &sync.Pool{
			New: func() any {
				return NewSelectStmt(tableName)
			},
		},
		updatePool: &sync.Pool{
			New: func() any {
				return NewUpdateStmt(tableName)
			},
		},
		deletePool: &sync.Pool{
			New: func() any {
				return NewDeleteStmt(tableName)
			},
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
	}
	return t
}

func (t *Table) SetDb(db *database.Database) {
	t.db = db
	t.creater.SetDb(db)
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

////////////////////////////////////////////////////////////////////////////////////////////////////
// Insert
////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *Table) GetInserter() *InsertStmt {
	insert := t.insertPool.Get().(*InsertStmt)
	insert.SetDb(t.db)
	insert.SetDbName(t.creater.DbName)
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
	selector.SetDb(t.db)
	selector.SetDbName(t.creater.DbName)
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
	updater.SetDb(t.db)
	updater.SetDbName(t.creater.DbName)
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
	deleter.SetDb(t.db)
	deleter.SetDbName(t.creater.DbName)
	return deleter
}

func (t *Table) PutDeleter(s *DeleteStmt) {
	s.Release()
	t.deletePool.Put(s)
}
