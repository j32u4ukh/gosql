package test

import (
	"testing"

	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/plugin"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"

	"github.com/pkg/errors"
)

type Tsukue struct {
	Id      int    `gorm:"pk=default;default=ai"`
	Content string `gorm:"size=3000"`
}

func InitTable() (*gosql.Table, error) {
	desk := &Tsukue{}
	tableParams, columnParams, err := plugin.GetStructParams(desk, dialect.MARIA)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get proto params")
	}

	table := gosql.NewTable("Desk", tableParams, columnParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	table.Init(&gosql.TableConfig{
		DbName:           "pekomiko",
		UseAntiInjection: false,
		InsertFunc:       plugin.InsertStruct,
		UpdateAnyFunc:    plugin.UpdateStruct,
		QueryFunc:        plugin.QueryStruct,
	})

	return table, nil
}

func TestCreateStmt(t *testing.T) {
	answer := "CREATE TABLE IF NOT EXISTS `pekomiko`.`Desk` (`Id` INT(11) NOT NULL AUTO_INCREMENT, `Content` VARCHAR(3000) NOT NULL DEFAULT '' COLLATE 'utf8mb4_bin', PRIMARY KEY (`Id`) USING BTREE) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';"
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
	answer := "INSERT INTO `pekomiko`.`Desk` (`Content`) VALUES ('abc');"
	table, err := InitTable()

	if err != nil {
		t.Errorf("TestCreateStmt | Failed to init table: %+v\n", err)
	}

	d1 := &Tsukue{Content: "abc"}

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
	answer := "SELECT * FROM `pekomiko`.`Desk` WHERE `Id` != -1 LIMIT 5 OFFSET 3;"
	table, err := InitTable()

	if err != nil {
		t.Errorf("TestCreateStmt | Failed to init table: %+v\n", err)
	}

	selector := table.GetSelector()
	selector.SetLimit(5)
	selector.SetOffset(3)
	selector.SetCondition(gosql.WS().Ne("Id", -1))
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestInsertStmt | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestInsertStmt |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestUpdateStmt(t *testing.T) {
	answer := "UPDATE `pekomiko`.`Desk` SET `Content` = 'xyz' WHERE `Id` = 2;"
	table, err := InitTable()

	if err != nil {
		t.Errorf("TestCreateStmt | Failed to init table: %+v\n", err)
	}

	updater := table.GetUpdater()
	desk := &Tsukue{
		Id:      2,
		Content: "xyz",
	}
	updater.UpdateAny(desk)
	where := gosql.WS().Eq("Id", 2)
	updater.SetCondition(where)
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
	answer := "DELETE FROM `pekomiko`.`Desk` WHERE `Id` = 1;"
	table, err := InitTable()

	if err != nil {
		t.Errorf("TestCreateStmt | Failed to init table: %+v\n", err)
	}

	deleter := table.GetDeleter()
	deleter.SetCondition(gosql.WS().Eq("Id", 1))
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
