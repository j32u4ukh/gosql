package main

import (
	"fmt"

	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/stmt/dialect"
)

// runoob
func InitWebsitesTable() (table *gdo.Table) {
	tableParams := stmt.NewTableParam()
	fmt.Printf("tableParams: %v\n", tableParams)
	colParam0 := stmt.NewColumnParam(0, "id", datatype.INT, dialect.MARIA)
	colParam0.SetPrimaryKey(stmt.ALGO)
	colParam1 := stmt.NewColumnParam(1, "name", datatype.VARCHAR, dialect.MARIA)
	colParam2 := stmt.NewColumnParam(2, "url", datatype.VARCHAR, dialect.MARIA)
	colParam2.SetSize(50)
	colParam3 := stmt.NewColumnParam(3, "alexa ", datatype.INT, dialect.MARIA)
	colParam4 := stmt.NewColumnParam(4, "contury ", datatype.VARCHAR, dialect.MARIA)
	colParams := []*stmt.ColumnParam{colParam0, colParam1, colParam2, colParam3, colParam4}
	// for i, col := range colParams {
	// 	fmt.Printf("%d) %+v\n", i, col)
	// }
	table = gdo.NewTable("Websites", tableParams, colParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	// fmt.Printf("%+v\n", table)
	return table
}

// SELECT DISTINCT `country` FROM `Websites`;
func main() {
	table := InitWebsitesTable()
	table.SetQueryMode(stmt.DistinctSelect)
	table.Query("country")
	sql, err := table.BuildSelectStmt()
	if err != nil {
		return
	}
	fmt.Printf("sql: %s\n", sql)
}
