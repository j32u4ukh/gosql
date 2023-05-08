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
	table.Query(stmt.NewSelectItem("name").UseBacktick(), stmt.NewSelectItem("country").UseBacktick())
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
	table.Query(stmt.NewSelectItem("country").UseBacktick().Distinct())
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

func TestRunoobSelect10(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE `name` LIKE 'G%';"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().Like("name", "G%")
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect10 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect10 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect11(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE NOT (`name` LIKE '%oo%');"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().Like("name", "%oo%").SetNotCondition()
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect11 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect11 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect12(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE `name` REGEXP '^[GFs]';"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().Regexp("name", "^[GFs]")
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect12 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect12 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect13(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE `name` IN ('Google', 'Apple');"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().In("name", "Google", "Apple")
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect13 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect13 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect14(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE `alexa` BETWEEN 1 AND 20;"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().Between("alexa", 1, 20)
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect14 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect14 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect15(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE NOT (`alexa` BETWEEN 1 AND 20);"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().Between("alexa", 1, 20).SetNotCondition()
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect15 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect15 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect16(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE (`alexa` BETWEEN 1 AND 20) AND (`country` IN ('USA', 'IND'));"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	cond1 := gdo.WS().Between("alexa", 1, 20).SetBrackets()
	cond2 := gdo.WS().In("country", "USA", "IND").SetNotCondition().SetBrackets()
	where := gdo.WS().AddAndCondtion(cond1).AddAndCondtion(cond2)
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect16 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect16 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect17(t *testing.T) {
	answer := "SELECT * FROM `Websites` WHERE `name` BETWEEN 'A' AND 'H';"
	table := InitWebsitesTable()
	//////////////////////////////////////////////////
	where := gdo.WS().Between("name", "A", "H")
	table.SetSelectCondition(where)
	//////////////////////////////////////////////////
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect17 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect17 |\nanswer: %s\nsql: %s", answer, sql)
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
