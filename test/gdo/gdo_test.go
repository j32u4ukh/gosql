package test

import (
	"fmt"
	"testing"

	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/proto"
)

func TestWhere1(t *testing.T) {
	answer := "`Id` = '\\' OR \\'1\\'=\\'1'"
	where := gdo.WS()
	where.Eq("Id", "' OR '1'='1")
	where.UseAntiInjection()
	sql, err := where.ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestWhere1 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestWhere1 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestWhere2(t *testing.T) {
	answer := "(`Id` = 1 OR `Id` = 2) AND `Content` != '123'"
	where := gdo.WS().
		AddAndCondtion(gdo.WS().SetBrackets().
			AddOrCondtion(gdo.WS().Eq("Id", 1)).
			AddOrCondtion(gdo.WS().Eq("Id", 2))).
		AddAndCondtion(gdo.WS().Ne("Content", "123"))

	where.UseAntiInjection()
	ws := where.ToStmtWhere()
	sql, err := ws.ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestWhere2 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestWhere2 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestTableCreate(t *testing.T) {
	answer := "CREATE TABLE IF NOT EXISTS `demo2`.`StmtDesk` (`Id` INT(11) NOT NULL DEFAULT 0, `Content` VARCHAR(3000) NOT NULL DEFAULT '', PRIMARY KEY (`Id`) USING BTREE) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';"
	table := InitTable()
	sql, err := table.BuildCreateStmt()
	fmt.Printf("sql: %s\n", sql)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestTableCreate | Errr: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestTableCreate |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestTableInsert(t *testing.T) {
	answer := "INSERT INTO `demo2`.`StmtDesk` (`Id`, `Content`) VALUES (50, '\\'); SELECT * FROM `demo2`.`StmtDesk`; -- hack');"
	table := InitTable()
	// map[string]any{"Id": 50, "Content": "'); SELECT * FROM `demo2`.`StmtDesk`; -- hack"}
	table.Insert([]any{50, "'); SELECT * FROM `demo2`.`StmtDesk`; -- hack"}, proto.ProtoToDb)
	sql, err := table.BuildInsertStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestInsert | Error %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestInsert |\nanswer: %s\nql: %s", answer, sql)
		}
	}
}

func TestTableSelect(t *testing.T) {
	answer := "SELECT * FROM `demo2`.`StmtDesk` WHERE `Content` = '\\' OR \\'1\\'=\\'1';"
	table := InitTable()

	// selector := gdo.NewSelectStmt("StmtDesk")
	// selector.SetDbName("demo2")
	where := gdo.WS().Eq("Content", "' OR '1'='1")
	table.SetSelectCondition(where)
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestSelect | Errr: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TesSelect |\nanswer: %s\nsql: %s", answer, sql)
		}
	}

}
func TestTableUpdte(t *testing.T) {
	answer := "UPDATE `demo2`.`StmtDesk` SET `Content` = '123' WHERE `Content` = '\\'; SELECT * FROM `demo2`.`StmtDesk`; -- hack';"
	table := InitTable()
	table.SetColumnNames([]string{"Id", "Content"})
	where := gdo.WS().Eq("Content", "'; SELECT * FROM `demo2`.`StmtDesk`; -- hack")
	table.Update("Content", "123", nil)
	table.SetUpdateCondition(where)
	sql, err := table.BuildUpdateStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestUpdate | rror: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TstUpdate |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestTableDelete(t *testing.T) {
	answer := "DELETE FROM `demo2`.`StmtDesk` WHERE `Content` = '\\'\\' OR \\'1\\'=\\'1\\'';"
	table := InitTable()
	where := gdo.WS().Eq("Content", "'' OR '1'='1'")
	table.SetDeleteCondition(where)
	sql, err := table.BuildDeleteStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestDelete | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestDelete |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}
