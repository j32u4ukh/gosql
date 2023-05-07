package test

import (
	"fmt"
	"testing"
)

func TestRunoobSelect1(t *testing.T) {
	answer := "SELECT `name`, `country` FROM `Websites`;"
	table := InitWebsitesTable()
	fmt.Printf("%+v\n", table)
	table.Query("name", "country")
	sql, err := table.BuildSelectStmt()

	if err != nil || sql != answer {
		if err != nil {
			t.Errorf("TestWhere1 | Error: %+v\n", err)
		}

		if sql != answer {
			t.Errorf("TestWhere1 |\nanswer: %s\nsql: %s", answer, sql)
		}
	}
}
