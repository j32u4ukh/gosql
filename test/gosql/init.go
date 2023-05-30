package test

import (
	"fmt"

	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/stmt/dialect"
)

func InitTable() *gosql.Table {
	tableName := "StmtDesk"
	tableParam := stmt.NewTableParam()

	table := gosql.NewTable(tableName, tableParam, nil, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	table.Init(&gosql.TableConfig{
		DbName:           "demo2",
		UseAntiInjection: true,
	})
	col1 := stmt.NewColumnParam(1, "Id", datatype.INT, dialect.MARIA)
	col1.SetPrimaryKey("default")
	table.AddColumn(stmt.NewColumn(col1))

	col2 := stmt.NewColumnParam(2, "Content", datatype.VARCHAR, dialect.MARIA)
	table.AddColumn(stmt.NewColumn(col2))
	return table
}

func InitWebsitesTable() *gosql.Table {
	tableParams := stmt.NewTableParam()
	fmt.Printf("tableParams: %v\n", tableParams)
	colParam0 := stmt.NewColumnParam(0, "id", datatype.INT, dialect.MARIA)
	colParam0.SetPrimaryKey(stmt.ALGO)
	colParam1 := stmt.NewColumnParam(1, "name", datatype.VARCHAR, dialect.MARIA)
	colParam2 := stmt.NewColumnParam(2, "url", datatype.VARCHAR, dialect.MARIA)
	colParam2.SetSize(50)
	colParam3 := stmt.NewColumnParam(3, "alexa", datatype.INT, dialect.MARIA)
	colParam4 := stmt.NewColumnParam(4, "contury", datatype.VARCHAR, dialect.MARIA)
	colParams := []*stmt.ColumnParam{colParam0, colParam1, colParam2, colParam3, colParam4}
	// for i, col := range colParams {
	// 	fmt.Printf("%d) %+v\n", i, col)
	// }
	table := gosql.NewTable("Websites", tableParams, colParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	return table
}
