package test

import (
	"fmt"
	"testing"

	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/stmt"
)

func TestRunoobSelect1(t *testing.T) {
	answer := "SELECT `name`, `country` FROM `Websites`;"
	table := InitWebsitesTable()
	fmt.Printf("%+v\n", table)
	table.Query("name", "country")
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect1 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect1 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect2(t *testing.T) {
	answer := "SELECT DISTINCT `country` FROM `Websites`;"
	table := InitWebsitesTable()
	table.SetQueryMode(stmt.DistinctSelect)
	table.Query("country")
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect2 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect2 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect3(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE `country` = 'CN';"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().Eq("country", "CN")
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect3 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect3 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect4(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE `id` = 1;"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().Eq("id", 1)
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect4 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect4 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect5(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE `country` = 'CN' AND `alexa` > 50;"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().AddAndCondtion(gdo.WS().Eq("country", "CN")).AddAndCondtion(gdo.WS().Gt("alexa", 50))
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect5 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect5 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect6(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE `country` = 'USA' OR `country` = 'CN';"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().AddOrCondtion(gdo.WS().Eq("country", "USA")).AddOrCondtion(gdo.WS().Eq("country", "CN"))
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect6 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect6 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect7(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE `alexa` > 15 AND (`country` = 'USA' OR `country` = 'CN');"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().
		AddAndCondtion(gdo.WS().Gt("alexa", 15)).
		AddAndCondtion(gdo.WS().SetBrackets().
			AddOrCondtion(gdo.WS().Eq("country", "USA")).
			AddOrCondtion(gdo.WS().Eq("country", "CN")))
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect7 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect7 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect8(t *testing.T) {
	answer := "SELECT * FROM `Websites` ORDER BY `alexa`;"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	table.SetOrderBy("alexa")
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect8 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect8 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect9(t *testing.T) {
	answer := "SELECT * FROM `Websites` ORDER BY `country`, `alexa`;"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	table.SetOrderBy("country", "alexa")
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect9 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect9 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}
