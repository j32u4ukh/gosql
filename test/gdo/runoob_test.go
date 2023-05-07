package test

import (
	"fmt"
	"testing"

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
