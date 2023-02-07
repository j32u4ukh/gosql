package test

import (
	"fmt"
	"testing"

	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/proto"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/j32u4ukh/gosql/test/pbgo"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func InitTable() (*proto.ProtoTable, error) {
	helper := &proto.Helper{
		Folder: "../pb",
		Dial:   dialect.MARIA,
	}

	tableName := "Desk"
	// GetParams(table_name string) (*stmt.TableParam, []*stmt.ColumnParam, error)
	tableParam, columnParam, err := helper.GetParams(tableName)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("TestInsertStmt | Error: %+v\n", err))
	}

	// NewProtoTable(name string, tableParam *stmt.TableParam, params []*stmt.ColumnParam, dial string)
	table := proto.NewProtoTable(tableName, tableParam, columnParam, helper.Dial)
	table.SetDbName("demo2")
	fmt.Printf("ProtoTable: %+v\n", table)
	return table, nil
}

func TestCreateStmt(t *testing.T) {
	answer := "CREATE TABLE IF NOT EXISTS `demo2`.`Desk` (`index` INT(11) NOT NULL AUTO_INCREMENT COMMENT '索引值', `user_name` VARCHAR(20) NOT NULL COMMENT '玩家名稱' COLLATE 'utf8mb4_bin', `item_id` SMALLINT(6) NOT NULL COMMENT '物品 ID', `time` TIMESTAMP NOT NULL DEFAULT current_timestamp() COMMENT 'Log 建立時間', PRIMARY KEY (`index`, `user_name`) USING BTREE) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';"
	table, err := InitTable()

	if err != nil {
		t.Error(fmt.Sprintf("TestCreateStmt | Failed to init table: %+v\n", err))
	}

	sql, err := table.BuildCreateStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestCreateStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestCreateStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestInsertStmt(t *testing.T) {
	answer := "INSERT INTO `demo2`.`Desk` (`index`, `user_name`, `item_id`, `time`) VALUES (NULL, '9527', 29, NULL);"
	table, err := InitTable()

	if err != nil {
		t.Error(fmt.Sprintf("TestCreateStmt | Failed to init table: %+v\n", err))
	}

	d1 := &pbgo.Desk{
		UserName: "9527",
		ItemId:   29,
	}

	table.InitByProtoMessage(d1)
	table.SetColumnNames([]string{"index", "user_name", "item_id", "time"})
	table.Insert([]protoreflect.ProtoMessage{d1})
	sql, err := table.BuildInsertStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestInsertStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestInsertStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestQueryStmt(t *testing.T) {
	answer := "SELECT * FROM `demo2`.`Desk` WHERE `Index` = 3 LIMIT 5 OFFSET 3;"
	table, err := InitTable()

	if err != nil {
		t.Error(fmt.Sprintf("TestCreateStmt | Failed to init table: %+v\n", err))
	}

	table.SetLimit(5)
	table.SetOffset(3)
	sql, err := table.BuildSelectStmt(gdo.WS().Eq("Index", 3))

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestInsertStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestInsertStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestUpdateStmt(t *testing.T) {
	answer := "UPDATE `demo2`.`Desk` SET `user_name` = '97', `item_id` = 3 WHERE `index` = 2;"
	table, err := InitTable()

	if err != nil {
		t.Error(fmt.Sprintf("TestCreateStmt | Failed to init table: %+v\n", err))
	}

	d1 := &pbgo.Desk{
		Index:    2,
		UserName: "97",
		ItemId:   3,
	}
	table.InitByProtoMessage(d1)
	sql, err := table.Update(d1, gdo.WS().Eq("index", 2))

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestUpdateStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestUpdateStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestDeleteDemo(t *testing.T) {
	answer := "DELETE FROM `demo2`.`Desk` WHERE `Index` = 3;"
	table, err := InitTable()

	if err != nil {
		t.Error(fmt.Sprintf("TestCreateStmt | Failed to init table: %+v\n", err))
	}

	d1 := &pbgo.Desk{
		Index:    3,
		UserName: "9527",
		ItemId:   29,
	}

	table.InitByProtoMessage(d1)
	sql, err := table.BuildDeleteStmt(gdo.WS().Eq("Index", 3))

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestUpdateStmt | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestUpdateStmt |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}
