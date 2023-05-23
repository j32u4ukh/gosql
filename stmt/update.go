package stmt

import (
	"fmt"
	"strings"

	"github.com/j32u4ukh/gosql/database"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// UpdateStmt
////////////////////////////////////////////////////////////////////////////////////////////////////
type UpdateStmt struct {
	DbName    string
	TableName string
	datas     []string
	Where     *WhereStmt
	// 是否允許不設置 Where 條件? 若不設置會造成全部數據都被修改，需額外允許才有作用
	allowEmptyWhere bool
	db              *database.Database
}

func NewUpdateStmt(tableName string) *UpdateStmt {
	s := &UpdateStmt{
		DbName:          "",
		TableName:       tableName,
		datas:           []string{},
		Where:           &WhereStmt{},
		allowEmptyWhere: false,
		db:              nil,
	}
	return s
}

func (s *UpdateStmt) SetDb(db *database.Database) {
	s.db = db
}

func (s *UpdateStmt) SetDbName(dbName string) *UpdateStmt {
	s.DbName = dbName
	return s
}

// 要求外部依據固定傳入順序，避免還要進行排序
func (s *UpdateStmt) Update(key string, value string) *UpdateStmt {
	s.datas = append(s.datas, fmt.Sprintf("`%s` = %s", key, value))
	return s
}

func (s *UpdateStmt) SetCondition(where *WhereStmt) *UpdateStmt {
	s.Where = where
	return s
}

func (s *UpdateStmt) AllowEmptyWhere() *UpdateStmt {
	s.allowEmptyWhere = true
	return s
}

func (s *UpdateStmt) Release() {
	s.datas = s.datas[:0]
	s.Where.Release()
	s.allowEmptyWhere = false
}

/*
   修改非 primary key 欄位
   UPDATE `demo2`.`Desk` SET `item_id`='3' WHERE  `index`=1 AND `user_id`=2;

   修改 primary key 欄位
   UPDATE `demo2`.`Desk` SET `user_id`='5' WHERE  `index`=1 AND `user_id`=2;
*/
func (s *UpdateStmt) ToStmt() (string, error) {
	if len(s.datas) == 0 {
		return "", errors.New("Update data is empty.")
	}

	var where, tableName string
	var err error

	if s.allowEmptyWhere {
		where = ""
	} else {
		where, err = s.Where.ToStmt()

		if err != nil || where == "" {
			return "", errors.Wrapf(err, "Failed to create where statment.")
		}

		where = fmt.Sprintf(" WHERE %s", where)
	}

	if s.DbName != "" {
		tableName = fmt.Sprintf("`%s`.`%s`", s.DbName, s.TableName)
	} else {
		tableName = fmt.Sprintf("`%s`", s.TableName)
	}

	sql := fmt.Sprintf("UPDATE %s SET %s%s;", tableName, strings.Join(s.datas, ", "), where)
	return sql, nil
}

func (s *UpdateStmt) Exec() (*database.SqlResult, error) {
	if s.db == nil {
		return nil, errors.New("Undefine database.")
	}
	sql, err := s.ToStmt()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate update statement.")
	}
	result, err := s.db.Exec(sql)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excute update statement.")
	}
	return result, nil
}
