package stmt

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type SelectStmt struct {
	DbName    string
	TableName string
	// 要查詢的欄位名稱列表
	QueryColumns []string
	// 查詢的篩選機制
	Where      *WhereStmt
	OrderBy    []string
	OrderType  string
	Limit      int32
	Offset     int32
	selectMode byte
}

func NewSelectStmt(tableName string) *SelectStmt {
	s := &SelectStmt{
		DbName:       "",
		TableName:    tableName,
		QueryColumns: []string{},
		Where:        &WhereStmt{},
		OrderBy:      nil,
		OrderType:    "",
		Limit:        -1,
		Offset:       0,
		selectMode:   NormalSelect,
	}
	return s
}

func (s *SelectStmt) SetDbName(name string) *SelectStmt {
	s.DbName = name
	return s
}

func (s *SelectStmt) SetQueryMode(mode byte) *SelectStmt {
	s.selectMode = mode
	return s
}

func (s *SelectStmt) Query(columns ...string) *SelectStmt {
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
	s.Where.Release()
	s.OrderBy = nil
	s.OrderType = "ASC"
	s.Limit = -1
	s.Offset = 0
}

func (s *SelectStmt) ToStmt() (string, error) {
	formatColumns, err := FormatColumns(s.QueryColumns, s.selectMode)

	if err != nil {
		return "", errors.Wrap(err, "Failed to format columns.")
	}

	where, err := s.Where.ToStmt()

	if err != nil {
		return "", errors.Wrap(err, "Failed to generate where statement.")
	}

	if where != "" {
		where = fmt.Sprintf(" WHERE %s", where)
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

type SelectItem struct {
	Name  string
	Alias string
}

func NewSelectItem() *SelectItem {
	return &SelectItem{}
}

func (s *SelectItem) SetName(name string, backtick bool) *SelectItem {
	if backtick {
		s.Name = fmt.Sprintf("`%s`", name)
	} else {
		s.Name = name
	}
	return s
}

func (s *SelectItem) SetAlias(alias string) *SelectItem {
	s.Alias = alias
	return s
}

func (s *SelectItem) Count() *SelectItem {
	s.Name = fmt.Sprintf("COUNT(%s)", s.Name)
	return s
}

func (s *SelectItem) Distinct() *SelectItem {
	s.Name = fmt.Sprintf("DISTINCT(%s)", s.Name)
	return s
}

func (s *SelectItem) Concat(elements ...string) *SelectItem {
	s.Name = fmt.Sprintf("CONCAT(%s)", strings.Join(elements, ", "))
	return s
}

func (s *SelectItem) ToStmt() string {
	result := s.Name
	if s.Alias != "" {
		result = fmt.Sprintf("%s AS %s", result, s.Alias)
	}
	return result
}

func FormatColumns(columns []string, mode byte) (string, error) {
	length := len(columns)

	switch length {
	case 0:
		switch mode {
		case DistinctSelect:
			return "", errors.New("You need to specify the columns when using DISTINCT.")
		case CountSelect:
			return "COUNT(*)", nil
		case CountDistinctSelect:
			return "", errors.New("You need to specify the columns when using DISTINCT.")
		case NormalSelect:
			fallthrough
		default:
			return "*", nil
		}
	case 1:
		switch mode {
		case DistinctSelect:
			return fmt.Sprintf("DISTINCT `%s`", columns[0]), nil
		case CountSelect:
			return fmt.Sprintf("COUNT(`%s`)", columns[0]), nil
		case CountDistinctSelect:
			return fmt.Sprintf("COUNT(DISTINCT `%s`)", columns[0]), nil
		case NormalSelect:
			fallthrough
		default:
			return columns[0], nil
		}
	default:
		temps := []string{}
		for _, column := range columns {
			temps = append(temps, fmt.Sprintf("`%s`", column))
		}
		result := strings.Join(temps, ", ")
		switch mode {
		case DistinctSelect:
			return fmt.Sprintf("DISTINCT %s", result), nil
		case CountSelect:
			return fmt.Sprintf("COUNT(%s)", result), nil
		case CountDistinctSelect:
			return fmt.Sprintf("COUNT(DISTINCT %s)", result), nil
		case NormalSelect:
			fallthrough
		default:
			return result, nil
		}
	}
}
