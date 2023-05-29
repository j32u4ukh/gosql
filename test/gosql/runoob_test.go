package test

import (
	"fmt"
	"testing"

	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/stmt"
)

func TestRunoobSelect1(t *testing.T) {
	answer := "SELECT `name`, `country` FROM `Websites`;"
	table := InitWebsitesTable()
	fmt.Printf("%+v\n", table)
	selector := table.GetSelector()
	selector.SetSelectItem(stmt.NewSelectItem("name").UseBacktick(), stmt.NewSelectItem("country").UseBacktick())
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	selector.SetSelectItem(stmt.NewSelectItem("country").UseBacktick().Distinct())
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().Eq("country", "CN")
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().Eq("id", 1)
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().
		AddAndCondtion(gosql.WS().Eq("country", "CN")).
		AddAndCondtion(gosql.WS().Gt("alexa", 50))
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().
		AddOrCondtion(gosql.WS().Eq("country", "USA")).
		AddOrCondtion(gosql.WS().Eq("country", "CN"))
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().
		AddAndCondtion(gosql.WS().Gt("alexa", 15)).
		AddAndCondtion(gosql.WS().SetBrackets().
			AddOrCondtion(gosql.WS().Eq("country", "USA")).
			AddOrCondtion(gosql.WS().Eq("country", "CN")))
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	selector.SetOrderBy("alexa")
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	selector.SetOrderBy("country", "alexa")
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().Like("name", "G%")
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().Like("name", "%oo%").SetNotCondition()
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().Regexp("name", "^[GFs]")
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().In("name", "Google", "Apple")
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().Between("alexa", 1, 20)
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().Between("alexa", 1, 20).SetNotCondition()
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	cond1 := gosql.WS().Between("alexa", 1, 20).SetBrackets()
	cond2 := gosql.WS().In("country", "USA", "IND").SetNotCondition().SetBrackets()
	where := gosql.WS().AddAndCondtion(cond1).AddAndCondtion(cond2)
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

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
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	where := gosql.WS().Between("name", "A", "H")
	selector.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect17 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect17 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect18(t *testing.T) {
	answer := "SELECT name AS n, country AS c FROM `Websites`;"
	table := InitWebsitesTable()
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	selector.SetSelectItem(
		stmt.NewSelectItem("name").SetAlias("n"),
		stmt.NewSelectItem("country").SetAlias("c"),
	)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect18 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect18 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobSelect19(t *testing.T) {
	answer := "SELECT name, CONCAT(url, ', ', alexa, ', ', country) AS site_info FROM `Websites`;"
	table := InitWebsitesTable()
	selector := table.GetSelector()
	//////////////////////////////////////////////////
	selector.SetSelectItem(
		stmt.NewSelectItem("name"),
		stmt.NewSelectItem("").Concat("url", "', '", "alexa", "', '", "country").SetAlias("site_info"),
	)
	//////////////////////////////////////////////////
	sql, err := selector.ToStmt()
	table.PutSelector(selector)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobSelect19 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobSelect19 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}

func TestRunoobInsert1(t *testing.T) {
	answer := "INSERT INTO `Websites` (`id`, `name`, `url`, `alexa`, `contury`) VALUES (NULL, '百度', 'https://www.baidu.com/', 4, 'CN');"
	table := InitWebsitesTable()
	inserter := table.GetInserter()
	//////////////////////////////////////////////////
	err := inserter.Insert([]any{"NULL", "百度", "https://www.baidu.com/", 4, "CN"})
	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}
	//////////////////////////////////////////////////
	sql, err := inserter.ToStmt()
	table.PutInserter(inserter)

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
	updater := table.GetUpdater()
	//////////////////////////////////////////////////
	updater.Update("alexa", 5000)
	updater.Update("country", "USA")
	where := gosql.WS().Eq("name", "ABC")
	updater.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := updater.ToStmt()
	table.PutUpdater(updater)

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
	updater := table.GetUpdater()
	//////////////////////////////////////////////////
	updater.Update("alexa", 5000)
	updater.Update("country", "USA")
	updater.AllowEmptyWhere()
	//////////////////////////////////////////////////
	sql, err := updater.ToStmt()
	table.PutUpdater(updater)

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
	deleter := table.GetDeleter()
	//////////////////////////////////////////////////
	where := gosql.WS().
		AddAndCondtion(gosql.WS().Eq("name", "Facebook")).
		AddAndCondtion(gosql.WS().Eq("country", "USA"))
	deleter.SetCondition(where)
	//////////////////////////////////////////////////
	sql, err := deleter.ToStmt()
	table.PutDeleter(deleter)

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
	deleter := table.GetDeleter()
	//////////////////////////////////////////////////
	deleter.AllowEmptyWhere()
	//////////////////////////////////////////////////
	sql, err := deleter.ToStmt()
	table.PutDeleter(deleter)

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestRunoobDelete2 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestRunoobDelete2 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}
