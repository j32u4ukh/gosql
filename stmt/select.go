package stmt

import (
	"fmt"
	"strings"

	"github.com/j32u4ukh/gosql/database"
	"github.com/pkg/errors"
)

type SelectStmt struct {
	DbName    string
	TableName string
	// 要查詢的項目列表
	QueryColumns []*SelectItem
	// 查詢的篩選機制
	Where     *WhereStmt
	OrderBy   []string
	OrderType string
	Limit     int32
	Offset    int32
	db        *database.Database
}

func NewSelectStmt(tableName string) *SelectStmt {
	s := &SelectStmt{
		DbName:       "",
		TableName:    tableName,
		QueryColumns: []*SelectItem{},
		Where:        nil,
		OrderBy:      nil,
		OrderType:    "",
		Limit:        -1,
		Offset:       0,
		db:           nil,
	}
	return s
}

func (s *SelectStmt) SetDb(db *database.Database) {
	s.db = db
}

func (s *SelectStmt) SetDbName(name string) *SelectStmt {
	s.DbName = name
	return s
}

func (s *SelectStmt) SetSelectItem(columns ...*SelectItem) *SelectStmt {
	s.QueryColumns = append(s.QueryColumns, columns...)
	return s
}

func (s *SelectStmt) SetCondition(where *WhereStmt) *SelectStmt {
	s.Where = where
	return s
}

func (s *SelectStmt) SetOrderBy(columns ...string) *SelectStmt {
	if s.OrderBy == nil {
		s.OrderBy = []string{}
	}
	s.OrderBy = append(s.OrderBy, columns...)
	return s
}

// 升序(ASC) | 降序(DESC)
func (s *SelectStmt) WhetherReverseOrder(reverse bool) *SelectStmt {
	if reverse {
		s.OrderType = "DESC"
	} else {
		s.OrderType = "ASC"
	}
	return s
}

/*
MySQL 支持 LIMIT 語句來選取指定的條數數據， Oracle 可以使用 ROWNUM 來選取。SQL Server / MS Access 則使用 SELECT TOP 語句來達到此效果。

TODO: 目前尚無法支援 SQL Server / MS Access
> SQL Server / MS Access 語法

SELECT TOP number|percent column_name(s)
FROM table_name;

> MySQL 語法

SELECT column_name(s)
FROM table_name
LIMIT number;
*/
func (s *SelectStmt) SetLimit(limit int32) *SelectStmt {
	s.Limit = limit
	return s
}

func (s *SelectStmt) SetOffset(offset int32) *SelectStmt {
	s.Offset = offset
	return s
}

func (s *SelectStmt) Release() {
	s.QueryColumns = s.QueryColumns[:0]
	s.Where = nil
	s.OrderBy = nil
	s.OrderType = "ASC"
	s.Limit = -1
	s.Offset = 0
}

func (s *SelectStmt) ToStmt() (string, error) {
	formatColumns, err := FormatColumns(s.QueryColumns)

	if err != nil {
		return "", errors.Wrap(err, "Failed to format columns.")
	}

	var where string = ""

	if s.Where != nil {
		where, err = s.Where.ToStmt()

		if err != nil {
			return "", errors.Wrap(err, "Failed to generate where statement.")
		}

		if where != "" {
			where = fmt.Sprintf(" WHERE %s", where)
		}
	}

	var order string
	var limitOffset string

	if s.OrderBy == nil {
		order = ""
	} else {
		orderColumns := []string{}
		for _, col := range s.OrderBy {
			orderColumns = append(orderColumns, fmt.Sprintf("`%s`", col))
		}
		order = fmt.Sprintf(" ORDER BY %s", strings.Join(orderColumns, ", "))

		if s.OrderType != "" {
			order = fmt.Sprintf("%s %s", order, s.OrderType)
		}
	}

	if s.Limit == -1 {
		limitOffset = ""
	} else {
		limitOffset = fmt.Sprintf(" LIMIT %d OFFSET %d", s.Limit, s.Offset)
	}

	var tableName string

	if s.DbName != "" {
		tableName = fmt.Sprintf("`%s`.`%s`", s.DbName, s.TableName)
	} else {
		tableName = fmt.Sprintf("`%s`", s.TableName)
	}

	sql := fmt.Sprintf("SELECT %s FROM %s%s%s%s;", formatColumns, tableName, where, order, limitOffset)
	return sql, nil
}

func (s *SelectStmt) Exec() (*database.SqlResult, error) {
	if s.db == nil {
		return nil, errors.New("Undefine database.")
	}
	sql, err := s.ToStmt()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate select statement.")
	}
	s.Release()
	result, err := s.db.Query(sql)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excute select statement.")
	}
	return result, nil
}

func (s *SelectStmt) Query(datas *[][]string) error {
	result, err := s.Exec()
	if err != nil {
		return errors.Wrap(err, "Failed to excute select statement.")
	}
	*datas = append(*datas, result.Datas...)
	return nil
}
