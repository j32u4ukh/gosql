package stmt

import (
	"fmt"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// BatchUpdateStmt
////////////////////////////////////////////////////////////////////////////////////////////////////
type BatchUpdateStmt struct {
	TableName   string
	PrimaryKey  string
	PrimaryKeys *cntr.Array[string]
	cols        *cntr.Array[string]
	sets        map[string]*SetStmt
	Where       *WhereStmt
}

// primaryKey: 多組數據時，根據此欄位來區分不同數據
func NewBatchUpdateStmt(tableName string, primaryKey string) *BatchUpdateStmt {
	s := &BatchUpdateStmt{
		TableName:   tableName,
		PrimaryKey:  primaryKey,
		PrimaryKeys: cntr.NewArray[string](),
		cols:        cntr.NewArray[string](),
		// key: column name
		sets:  map[string]*SetStmt{},
		Where: nil,
	}
	return s
}

func (s *BatchUpdateStmt) Update(key string, col string, value string) *BatchUpdateStmt {
	if !s.PrimaryKeys.Contains(key) {
		s.PrimaryKeys.Append(key)
	}

	if _, ok := s.sets[col]; !ok {
		s.sets[col] = newSetStmt(s.PrimaryKey, col)
		s.cols.Append(col)
	}

	s.sets[col].AddData(key, value)
	return s
}

func (s *BatchUpdateStmt) SetCondition(where *WhereStmt) *BatchUpdateStmt {
	s.Where = where
	return s
}

// 取得緩存數量
func (s *BatchUpdateStmt) GetBufferNumber() int {
	return s.PrimaryKeys.Length()
}

func (s *BatchUpdateStmt) Release() {
	s.PrimaryKeys.Clear()

	for k := range s.sets {
		delete(s.sets, k)
	}
}

func (s *BatchUpdateStmt) ToStmt() (string, error) {
	sets := []string{}
	var set *SetStmt
	var stmt string
	var err error

	for _, col := range s.cols.Elements {
		set = s.sets[col]
		stmt, err = set.toStmt()

		if err != nil {
			return "", errors.Wrap(err, "Failed to generate set statement.")
		}

		sets = append(sets, stmt)
	}

	setStmt := strings.Join(sets, ", ")
	where := WS().In(s.PrimaryKey, s.PrimaryKeys.Elements...)
	wstmt, err := where.AddAndCondtion(s.Where).ToStmt()

	if err != nil {
		return "", errors.Wrap(err, "Failed to generate where statement.")
	}

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s;", s.TableName, setStmt, wstmt)
	return sql, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// SetStmt
////////////////////////////////////////////////////////////////////////////////////////////////////
type SetStmt struct {
	key    string
	column string
	keys   []string
	values []string
}

func newSetStmt(key string, column string) *SetStmt {
	s := &SetStmt{
		key:    key,
		column: column,
		keys:   []string{},
		values: []string{},
	}
	return s
}

func (s *SetStmt) AddData(key string, value string) {
	s.keys = append(s.keys, key)
	s.values = append(s.values, value)
}

func (s *SetStmt) toStmt() (string, error) {
	// SET [column] = CASE [primary_key]
	// WHEN m0.fileds.Get(0) THEN 'Insert1'
	// WHEN m1.fileds.Get(0) THEN 'Insert3'
	// END
	content := []string{}
	var i int
	var key string

	for i, key = range s.keys {
		// WHEN [value of primary column] THEN [value of target column]
		content = append(content, fmt.Sprintf("WHEN %s THEN %s", key, s.values[i]))
	}

	return fmt.Sprintf("`%s` = CASE `%s` %s END", s.column, s.key, strings.Join(content, " ")), nil
}
