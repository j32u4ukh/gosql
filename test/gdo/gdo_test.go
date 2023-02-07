package test

import (
	"fmt"
	"testing"

	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/proto"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/stmt/dialect"
)

func TestWhere1(t *testing.T) {
	answer := "`Id` = '\\' OR \\'1\\'=\\'1'"
	where := gdo.WS()
	where.Eq("Id", "' OR '1'='1")
	where.UseAntiInjection()
	sql, err := where.ToStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestWhere1 | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestWhere1 |\nanswer: %s\nsql: %s", answer, sql))
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
			t.Error(fmt.Sprintf("TestWhere2 | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestWhere2 |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}

func TestInsert(t *testing.T) {
	// answer := "INSERT INTO `demo2`.`StmtDesk` (`Id`, `Content`) VALUES (50, '\\'); SELECT * FROM `demo2`.`StmtDesk`; -- hack');"
	// insert := gdo.NewInsertStmt("StmtDesk")
	// insert.SetDbName("demo2")
	// insert.SetColumnNames([]string{"Id", "Content"})
	// insert.UseAntiInjection()
	// insert.Insert([]any{50, "'); SELECT * FROM `demo2`.`StmtDesk`; -- hack"})
	// sql, err := insert.ToStmt()

	// if err != nil || sql != answer {
	// 	if err != nil {
	// 		t.Error(fmt.Sprintf("TestInsert | Error: %+v\n", err))
	// 	}

	// 	if sql != answer {
	// 		t.Error(fmt.Sprintf("TestInsert |\nanswer: %s\nsql: %s", answer, sql))
	// 	}
	// }
}

func TestSelect(t *testing.T) {
	// answer := "SELECT * FROM `demo2`.`StmtDesk` WHERE `Content` = '\\' OR \\'1\\'=\\'1';"
	// selector := gdo.NewSelectStmt("StmtDesk")
	// selector.SetDbName("demo2")
	// where := gdo.WS().Eq("Content", "' OR '1'='1")
	// where.UseAntiInjection()
	// sql, err := selector.SetCondition(where).ToStmt()

	// if err != nil || sql != answer {
	// 	if err != nil {
	// 		t.Error(fmt.Sprintf("TestSelect | Error: %+v\n", err))
	// 	}

	// 	if sql != answer {
	// 		t.Error(fmt.Sprintf("TestSelect |\nanswer: %s\nsql: %s", answer, sql))
	// 	}
	// }
}

func TestUpdate(t *testing.T) {
	// answer := "UPDATE `demo2`.`StmtDesk` SET `Content` = '123' WHERE `Content` = '\\'; SELECT * FROM `demo2`.`StmtDesk`; -- hack';"
	// update := gdo.NewUpdateStmt("StmtDesk")
	// update.SetDbName("demo2")
	// // update.SetColumnNames([]string{"Id", "Content"})
	// where := gdo.WS().Eq("Content", "'; SELECT * FROM `demo2`.`StmtDesk`; -- hack")
	// update.Update("Content", "123")
	// update.SetCondition(where)
	// update.UseAntiInjection()
	// sql, err := update.ToStmt()

	// if err != nil || sql != answer {
	// 	if err != nil {
	// 		t.Error(fmt.Sprintf("TestUpdate | Error: %+v\n", err))
	// 	}

	// 	if sql != answer {
	// 		t.Error(fmt.Sprintf("TestUpdate |\nanswer: %s\nsql: %s", answer, sql))
	// 	}
	// }
}

func TestDelete(t *testing.T) {
	// answer := "DELETE FROM `StmtDesk` WHERE `Content` = '\\'\\' OR \\'1\\'=\\'1\\'';"
	// del := gdo.NewDeleteStmt("StmtDesk")
	// where := gdo.WS().Eq("Content", "'' OR '1'='1'")
	// del.SetCondition(where)
	// del.UseAntiInjection()
	// sql, err := del.ToStmt()

	// if err != nil || sql != answer {
	// 	if err != nil {
	// 		t.Error(fmt.Sprintf("TestDelete | Error: %+v\n", err))
	// 	}

	// 	if sql != answer {
	// 		t.Error(fmt.Sprintf("TestDelete |\nanswer: %s\nsql: %s", answer, sql))
	// 	}
	// }
}

func InitTable() *gdo.Table {
	tableName := "StmtDesk"
	tableParam := stmt.NewTableParam()

	// NewTable(name string, tableParam *stmt.TableParam, columnParams []*stmt.ColumnParam, engine string, collate string, dial string)
	table := gdo.NewTable(tableName, tableParam, nil, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	table.SetDbName("demo2")
	table.UseAntiInjection(true)
	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	table.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	// col2.SetCanNull(true)
	table.AddColumn(stmt.NewColumn(col2))
	return table
}

func TestTableCreate(t *testing.T) {
	answer := "CREATE TABLE IF NOT EXISTS `demo2`.`StmtDesk` (`Id` INT(11) NOT NULL, `Content` VARCHAR(3000) NOT NULL, PRIMARY KEY (`Id`) USING BTREE) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';"
	table := InitTable()
	sql, err := table.BuildCreateStmt()
	fmt.Printf("sql: %s\n", sql)

	if err != nil || sql != answer {
		if err != nil {
			t.Error(fmt.Sprintf("TestTableCreate | Errr: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestTableCreate |\nanswer: %s\nsql: %s", answer, sql))
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
			t.Error(fmt.Sprintf("TestInsert | Error %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestInsert |\nanswer: %s\nql: %s", answer, sql))
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
			t.Error(fmt.Sprintf("TestSelect | Errr: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TesSelect |\nanswer: %s\nsql: %s", answer, sql))
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
			t.Error(fmt.Sprintf("TestUpdate | rror: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TstUpdate |\nanswer: %s\nsql: %s", answer, sql))
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
			t.Error(fmt.Sprintf("TestDelete | Error: %+v\n", err))
		}

		if sql != answer {
			t.Error(fmt.Sprintf("TestDelete |\nanswer: %s\nsql: %s", answer, sql))
		}
	}
}
