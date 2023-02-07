package stmt

import (
	"bytes"
	"fmt"
	"sort"
)

type DbOp byte

const (
	DbQuery DbOp = iota
	DbDelete
	DbInsert
	DbUpdate
)

type IStatement interface {
	ToStmt() string
}

// 預設參數
const (
	// INNODB 存儲引擎，多列索引的長度限制： 每個列的長度不能大於767 bytes；所有組成索引列的長度和不能大於 3072 bytes
	ENGINE string = "InnoDB"
	// 排序規則
	COLLATE string = "utf8mb4_bin"
	// Primary key 演算法
	ALGO string = "BTREE"
)

// 根據表格名稱，檢查該名稱的表格是否存在
func IsTableExists(schemaName string, tableName string) string {
	var where string
	if schemaName == "" {
		where = fmt.Sprintf("TABLE_NAME = '%s'", tableName)
	} else {
		where = fmt.Sprintf("`TABLE_SCHEMA`='%s' AND `TABLE_NAME`='%s'", schemaName, tableName)
	}
	return fmt.Sprintf("SELECT `TABLE_NAME` FROM INFORMATION_SCHEMA.`TABLES` WHERE %s;", where)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////////////////////////
type KeyValuePair struct {
	Key   string
	Value string
}

func NewKeyValuePair(key string, value string) *KeyValuePair {
	return &KeyValuePair{Key: key, Value: value}
}

func (p KeyValuePair) String() string {
	return fmt.Sprintf("`%s` = %s", p.Key, p.Value)
}

type KeyValueSlice struct {
	Pairs []*KeyValuePair
	ByKey bool
}

func NewKeyValueSlice() *KeyValueSlice {
	s := &KeyValueSlice{
		Pairs: []*KeyValuePair{},
		ByKey: true,
	}
	return s
}

func (s *KeyValueSlice) AddPair(pair *KeyValuePair) {
	s.Pairs = append(s.Pairs, pair)
	// fmt.Printf("(s *KeyValueSlice) AddPair | pair: %+v, #Pairs: %d\n", pair, len(s.Pairs))
}

func (s *KeyValueSlice) Len() int {
	return len(s.Pairs)
}

func (s *KeyValueSlice) Less(i int, j int) bool {
	if s.ByKey {
		return s.Pairs[i].Key < s.Pairs[j].Key
	} else {
		return s.Pairs[i].Value < s.Pairs[j].Value
	}
}

func (s *KeyValueSlice) Swap(i int, j int) {
	s.Pairs[i], s.Pairs[j] = s.Pairs[j], s.Pairs[i]
}

func (s *KeyValueSlice) SortByKey() {
	s.ByKey = true
	sort.Sort(s)
}

func (s *KeyValueSlice) SortByValue() {
	s.ByKey = false
	sort.Sort(s)
}

func (s *KeyValueSlice) IsEmpty() bool {
	return len(s.Pairs) == 0
}

func (s *KeyValueSlice) Release() {
	s.Pairs = s.Pairs[:0]
}

func (s *KeyValueSlice) String() string {
	if len(s.Pairs) != 0 {
		var buffer bytes.Buffer
		buffer.WriteString(s.Pairs[0].String())
		for _, pair := range s.Pairs[1:] {
			buffer.WriteString(", ")
			buffer.WriteString(pair.String())
		}
		return buffer.String()
	}
	return ""
}
