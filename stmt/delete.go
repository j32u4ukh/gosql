package stmt

import (
	"fmt"

	"github.com/pkg/errors"
)

type DeleteStmt struct {
	DbName string
	Name   string
	Where  *WhereStmt
	// 是否允許不設置 Where 條件? 若不設置會 刪除資料表中所有的資料，需額外允許才有作用
	allowEmptyWhere bool
}

func NewDeleteStmt(name string) *DeleteStmt {
	s := &DeleteStmt{
		DbName:          "",
		Name:            name,
		Where:           &WhereStmt{},
		allowEmptyWhere: false,
	}
	return s
}

func (s *DeleteStmt) SetDbName(dbName string) *DeleteStmt {
	s.DbName = dbName
	return s
}

func (s *DeleteStmt) SetCondition(where *WhereStmt) *DeleteStmt {
	s.Where = where
	return s
}

func (s *DeleteStmt) AllowEmptyWhere() *DeleteStmt {
	s.allowEmptyWhere = true
	return s
}

func (s *DeleteStmt) Release() *DeleteStmt {
	s.Where.Release()
	s.allowEmptyWhere = false
	return s
}

func (s *DeleteStmt) ToStmt() (string, error) {
	var where, tableName string
	var err error
	if s.allowEmptyWhere {
		where = ""
	} else {
		where, err = s.Where.ToStmt()

		if err != nil {
			return "", errors.Wrapf(err, "Failed to create where statment.")
		}

		if where == "" {
			return "", errors.New("Where condtion 為空")
		}

		where = fmt.Sprintf(" WHERE %s", where)
	}
	if s.DbName != "" {
		tableName = fmt.Sprintf("`%s`.`%s`", s.DbName, s.Name)
	} else {
		tableName = fmt.Sprintf("`%s`", s.Name)
	}
	sql := fmt.Sprintf("DELETE FROM %s%s;", tableName, where)
	return sql, nil
}
