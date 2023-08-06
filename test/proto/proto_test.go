package test

import (
	"fmt"
	"testing"

	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/plugin"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/j32u4ukh/gosql/test/pbgo"

	"github.com/pkg/errors"
)

func InitTable() (*gosql.Table, error) {
	tableName := "Desk"
	tableParams, columnParams, err := plugin.GetProtoParams(fmt.Sprintf("../pb/%s.proto", tableName), dialect.MARIA)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get proto params")
	}

	table := gosql.NewTable(tableName, tableParams, columnParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	table.Init(&gosql.TableConfig{
		DbName:           "demo2",
		UseAntiInjection: false,
		PtrToDbFunc:      plugin.ProtoToDb,
		InsertFunc:       plugin.InsertProto,
		QueryFunc:        plugin.QueryProto,
		UpdateAnyFunc:    plugin.UpdateProto,
	})

	return table, nil
}

func TestCreateStmt(t *testing.T) {
	answer := "CREATE TABLE IF NOT EXISTS `demo2`.`Desk` (`index` INT(11) NOT NULL AUTO_INCREMENT COMMENT '索引值', `user_name` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '玩家名稱' COLLATE 'utf8mb4_bin', `item_id` SMALLINT(6) NOT NULL DEFAULT 0 COMMENT '物品 ID', `time` TIMESTAMP NOT NULL DEFAULT current_timestamp() COMMENT 'Log 建立時間', PRIMARY KEY (`index`, `user_name`) USING BTREE) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';"
	table, err := InitTable()

	if err != nil {
		t.Errorf("TestCreateStmt | Failed to init table: %+v\n", err)
	}

	sql, err := table.Creater().ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestCreateStmt | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestCreateStmt |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestInsertStmt(t *testing.T) {
	answer := "INSERT INTO `demo2`.`Desk` (`user_name`, `item_id`) VALUES ('9527', 29);"
	table, err := InitTable()

	if err != nil {
		t.Errorf("TestCreateStmt | Failed to init table: %+v\n", err)
	}

	d1 := &pbgo.Desk{
		UserName: "9527",
		ItemId:   29,
	}

	inserter := table.GetInserter()
	inserter.Insert(d1)
	sql, err := inserter.ToStmt()
	table.PutInserter(inserter)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestInsertStmt | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestInsertStmt |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestQueryStmt(t *testing.T) {
	answer := "SELECT * FROM `demo2`.`Desk` WHERE `Index` = 3 LIMIT 5 OFFSET 3;"
	table, err := InitTable()

	if err != nil {
		t.Errorf("TestCreateStmt | Failed to init table: %+v\n", err)
	}

	selector := table.GetSelector()
	selector.SetLimit(5)
	selector.SetOffset(3)
	selector.SetCondition(gosql.WS().Eq("Index", 3))
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestQueryStmt | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestQueryStmt |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestSelectItemStmt(t *testing.T) {
	answer := "SELECT `UserName`, `ItemId` FROM `demo2`.`Desk` WHERE `Index` = 3 LIMIT 5 OFFSET 3;"
	table, err := InitTable()

	if err != nil {
		t.Errorf("TestCreateStmt | Failed to init table: %+v\n", err)
	}

	selector := table.GetSelector()
	selector.SetSelectItem(stmt.NewSelectItem("UserName").UseBacktick())
	selector.SetSelectItem(stmt.NewSelectItem("ItemId").UseBacktick())
	selector.SetLimit(5)
	selector.SetOffset(3)
	selector.SetCondition(gosql.WS().Eq("Index", 3))
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestSelectItemStmt | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestSelectItemStmt |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestUpdateStmt(t *testing.T) {
	answer := "UPDATE `demo2`.`Desk` SET `user_name` = '97', `item_id` = 3 WHERE `index` = 2;"
	table, err := InitTable()

	if err != nil {
		t.Errorf("TestCreateStmt | Failed to init table: %+v\n", err)
	}

	d1 := &pbgo.Desk{
		Index:    2,
		UserName: "97",
		ItemId:   3,
	}

	updater := table.GetUpdater()
	updater.SetCondition(gosql.WS().Eq("index", 2))
	updater.UpdateAny(d1)
	sql, err := updater.ToStmt()
	table.PutUpdater(updater)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestUpdateStmt | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestUpdateStmt |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestDeleteDemo(t *testing.T) {
	answer := "DELETE FROM `demo2`.`Desk` WHERE `Index` = 3;"
	table, err := InitTable()

	if err != nil {
		t.Errorf("TestCreateStmt | Failed to init table: %+v\n", err)
	}

	deleter := table.GetDeleter()
	deleter.SetCondition(gosql.WS().Eq("Index", 3))
	sql, err := deleter.ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestUpdateStmt | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestUpdateStmt |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}
