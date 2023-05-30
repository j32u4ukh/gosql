package test

import (
	"fmt"
	"testing"

	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/gdo"
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
	sql, err := table.Creater().ToStmt()
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
	inserter := table.GetInserter()
	// map[string]any{"Id": 50, "Content": "'); SELECT * FROM `demo2`.`StmtDesk`; -- hack"}
	inserter.Insert([]any{50, "'); SELECT * FROM `demo2`.`StmtDesk`; -- hack"})
	sql, err := inserter.ToStmt()
	table.PutInserter(inserter)

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
	where := gosql.WS().Eq("Content", "' OR '1'='1")
	selector := table.GetSelector()
	selector.SetCondition(where)
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestSelect | Errr: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TesSelect |\nanswer: %s\nsql: %s", answer, sql)
		}
	}

}

// sql: UPDATE `demo2`.`StmtDesk` SET `Content` = '123' WHERE `Content` = ''; SELECT * FROM `demo2`.`StmtDesk`; -- hack';
func TestTableUpdte(t *testing.T) {
	answer := "UPDATE `demo2`.`StmtDesk` SET `Content` = '123' WHERE `Content` = '\\'; SELECT * FROM `demo2`.`StmtDesk`; -- hack';"
	table := InitTable()
	updater := table.GetUpdater()
	// table.SetColumnNames([]string{"Id", "Content"})
	where := gosql.WS().Eq("Content", "'; SELECT * FROM `demo2`.`StmtDesk`; -- hack")
	updater.Update("Content", "123")
	updater.SetCondition(where)
	sql, err := updater.ToStmt()
	table.PutUpdater(updater)

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
	where := gosql.WS().Eq("Content", "'' OR '1'='1'")
	deleter := table.GetDeleter()
	deleter.SetCondition(where)
	sql, err := deleter.ToStmt()
	table.PutDeleter(deleter)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestDelete | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestDelete |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}
