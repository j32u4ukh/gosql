package stmt

import (
	"fmt"
	"strings"

	"github.com/j32u4ukh/gosql/database"
	"github.com/pkg/errors"
)

type InsertStmt struct {
	DbName    string
	TableName string
	// 要求包含所有欄位 ex: (`timestamp`, `text`, `flag`)
	ColumnStmt string
	// 多筆按照欄位順序傳入的數據
	datas [][]string
	db    *database.Database
}

func NewInsertStmt(tableName string) *InsertStmt {
	s := &InsertStmt{
		TableName:  tableName,
		ColumnStmt: "",
		datas:      [][]string{},
		db:         nil,
	}
	return s
}

func (s *InsertStmt) SetDb(db *database.Database) {
	s.db = db
}

func (s *InsertStmt) GetDb() *database.Database {
	return s.db
}

func (s *InsertStmt) SetDbName(dbName string) {
	s.DbName = dbName
}

func (s *InsertStmt) SetColumnNames(names []string) {
	columnStmt := []string{}

	for _, name := range names {
		columnStmt = append(columnStmt, fmt.Sprintf("`%s`", name))
	}

	// ex: (`id`, `name`, `item_id`)
	s.ColumnStmt = fmt.Sprintf("(%s)", strings.Join(columnStmt, ", "))
}

// 添加一筆數據(最終可同時添加多筆數據)
// 呼叫此函式者，須確保 datas 中的欄位都存在表格中
// 對於 允許空值 或 不允許但會自行賦值 的欄位，則傳入 NULL 即可
func (s *InsertStmt) Insert(data []string) *InsertStmt {
	s.datas = append(s.datas, data)
	return s
}

func (s *InsertStmt) GetBufferNumber() int32 {
	return int32(len(s.datas))
}

func (s *InsertStmt) Release() {
	s.datas = s.datas[:0]
}

// 形成 SQL 語法
func (s *InsertStmt) ToStmt() (string, error) {
	// INSERT `demo`.`message` (`timestamp`, `text`, `flag`) VALUES (1234567890123, "Insert", true),  (1234567890123, "Insert", true);
	valueStmts := []string{}
	var data []string
	var valueStmt, tableName string

	// 遍歷每一條數據
	for _, data = range s.datas {
		valueStmt = fmt.Sprintf("(%s)", strings.Join(data, ", "))
		valueStmts = append(valueStmts, valueStmt)
	}

	if s.DbName != "" {
		tableName = fmt.Sprintf("`%s`.`%s`", s.DbName, s.TableName)
	} else {
		tableName = fmt.Sprintf("`%s`", s.TableName)
	}

	sql := fmt.Sprintf("INSERT INTO %s %s VALUES %s;", tableName, s.ColumnStmt, strings.Join(valueStmts, ", "))
	return sql, nil
}

func (s *InsertStmt) Exec() (*database.SqlResult, error) {
	if s.db == nil {
		return nil, errors.New("Undefine database.")
	}
	sql, err := s.ToStmt()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate insert statement.")
	}
	s.Release()
	result, err := s.db.Exec(sql)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excute insert statement.")
	}
	return result, nil
}
