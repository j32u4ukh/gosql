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

func TestRunoobInsert1(t *testing.T) {
	answer := "INSERT INTO `Websites` (`id`, `name`, `url`, `alexa`, `contury`) VALUES (NULL, '百度', 'https://www.baidu.com/', 4, 'CN');"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	err := table.Insert([]any{"NULL", "百度", "https://www.baidu.com/", 4, "CN"}, nil)
	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}
	//////////////////////////////////////////////////
	sql, err := table.BuildInsertStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobInsert1 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobInsert1 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobUpdate1(t *testing.T) {
	answer := "UPDATE `Websites` SET `alexa` = 5000, `country` = 'USA' WHERE `name` = 'ABC';"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	table.Update("alexa", 5000, nil)
	table.Update("country", "USA", nil)
	where := gdo.WS().Eq("name", "ABC")
	table.SetUpdateCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildUpdateStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobUpdate1 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobUpdate1 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobUpdate2(t *testing.T) {
	answer := "UPDATE `Websites` SET `alexa` = 5000, `country` = 'USA';"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	table.Update("alexa", 5000, nil)
	table.Update("country", "USA", nil)
	table.AllowEmptyUpdateCondition()
	//////////////////////////////////////////////////
	sql, err := table.BuildUpdateStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobUpdate2 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobUpdate2 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobDelete1(t *testing.T) {
	answer := "DELETE FROM `Websites` WHERE `name` = 'Facebook' AND `country` = 'USA';"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().AddAndCondtion(gdo.WS().Eq("name", "Facebook")).AddAndCondtion(gdo.WS().Eq("country", "USA"))
	table.SetDeleteCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildDeleteStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobDelete1 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobDelete1 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobDelete2(t *testing.T) {
	answer := "DELETE FROM `Websites`;"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	table.AllowEmptyDeleteCondition()
	//////////////////////////////////////////////////
	sql, err := table.BuildDeleteStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobDelete2 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobDelete2 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}
