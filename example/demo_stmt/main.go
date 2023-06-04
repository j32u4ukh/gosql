package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/j32u4ukh/gosql/utils"
)

var db *database.Database
var sql string
var result *database.SqlResult
var err error
var logger *glog.Logger
var table *stmt.Table

func main() {
	logger = glog.SetLogger(0, "demo_stmt", glog.DebugLevel)
	utils.SetLogger(logger)
	logger.SetFolder("../log")

	command := strings.ToLower(os.Args[1])
	conf, err := database.NewConfig("../config/config.yaml")

	if err != nil {
		logger.Error("讀取 Config 檔時發生錯誤, err: %+v\n", err)
		return
	}

	dc := conf.GetDatabase()
	db, err = database.Connect(0, dc.UserName, dc.Password, dc.Server, dc.Port, dc.Name)

	if err != nil {
		logger.Error("與資料庫連線時發生錯誤, err: %+v", err)
		return
	}

	defer db.Close()
	db = database.Get(0)

	if db == nil {
		logger.Error("Database(0) is not exists.")
		return
	}

	tableParams := stmt.NewTableParam()
	// fmt.Printf("tableParams: %v\n", tableParams)
	colParam0 := stmt.NewColumnParam(0, "id", datatype.INT, dialect.MARIA)
	colParam0.SetPrimaryKey(stmt.ALGO)
	colParam0.SetDefault("AI")
	colParam1 := stmt.NewColumnParam(1, "name", datatype.VARCHAR, dialect.MARIA)
	colParam1.SetSize(50)
	colParam2 := stmt.NewColumnParam(2, "url", datatype.VARCHAR, dialect.MARIA)
	colParam2.SetSize(50)
	colParam3 := stmt.NewColumnParam(3, "alexa", datatype.INT, dialect.MARIA)
	colParam4 := stmt.NewColumnParam(4, "contury", datatype.VARCHAR, dialect.MARIA)
	colParam4.SetSize(50)
	colParams := []*stmt.ColumnParam{colParam0, colParam1, colParam2, colParam3, colParam4}
	// for i, col := range colParams {
	// 	fmt.Printf("%d) %+v\n", i, col)
	// }
	table = stmt.NewTable("Websites", tableParams, colParams, stmt.ENGINE, stmt.COLLATE)
	table.SetDb(db)

	switch command {
	case "c":
		CreateDemo()
	case "i":
		InsertDemo()
	case "q":
		QueryDemo()
	case "u":
		UpdateDemo()
	case "d":
		DeleteDemo()
	default:
		fmt.Printf("No invalid command(%s).\n", command)
	}
}

func CreateDemo() {
	sql, err = table.Creater().ToStmt()

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}
	fmt.Printf("Creater sql: %s\n", sql)
	result, err = table.Creater().Exec()

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}
	fmt.Printf("result: %s\n", result)
}

/*
INSERT INTO `Websites` (`id`, `name`, `url`, `alexa`, `contury`) VALUES
(NULL, 'Google', https://www.google.com/, 1, 'US'),
(NULL, 'Facebook', https://www.facebook.com/, 2, 'US'),
(NULL, 'apple', https://www.apple.com/, 3, 'US'),
(NULL, 'microsoft', https://www.microsoft.com/, 4, 'US');
*/
func InsertDemo() {
	insert := table.GetInserter()
	insert.Insert([]string{"NULL", "'Google'", "'https://www.google.com/'", "1", "'US'"})
	insert.Insert([]string{"NULL", "'Facebook'", "'https://www.facebook.com/'", "2", "'US'"})
	insert.Insert([]string{"NULL", "'apple'", "'https://www.apple.com/'", "3", "'US'"})
	insert.Insert([]string{"NULL", "'microsoft'", "'https://www.microsoft.com/'", "4", "'US'"})
	sql, err = insert.ToStmt()

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	fmt.Printf("insert1 sql: %s\n", sql)
	result, err = insert.Exec()

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	fmt.Printf("result: %+v\n", result)
	table.PutInserter(insert)
}

// SELECT * FROM `Websites`;
func QueryDemo() {
	selector := table.GetSelector()
	sql, err = selector.ToStmt()
	if err != nil {
		fmt.Printf("Select err: %+v\n", err)
		return
	}
	fmt.Printf("selector sql: %s\n", sql)
	result, err = selector.Exec()
	if err != nil {
		fmt.Printf("Select err: %+v\n", err)
		return
	}
	fmt.Printf("result: %s\n", result)

	datas := [][]string{}
	selector.Query(&datas)
	fmt.Printf("datas: %+v\n", datas)
	table.PutSelector(selector)
}

// UPDATE `Websites` SET `alexa` = 5000 WHERE `id` = 3;
func UpdateDemo() {
	updater := table.GetUpdater()
	updater.Update("alexa", "5000")
	updater.SetCondition(stmt.WS().Eq("name", "'Facebook'"))
	sql, err = updater.ToStmt()
	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}
	fmt.Printf("updater sql: %s\n", sql)

	result, err = updater.Exec()
	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}
	fmt.Printf("result: %+v\n", result)
	table.PutUpdater(updater)
}

// DELETE FROM `Websites` WHERE `name` = Facebook;
func DeleteDemo() {
	deleter := table.GetDeleter()
	deleter.SetCondition(stmt.WS().Eq("name", "'Facebook'"))
	sql, err = deleter.ToStmt()
	if err != nil {
		fmt.Printf("Delete err: %+v\n", err)
		return
	}
	fmt.Printf("deleter sql: %s\n", sql)

	result, err = deleter.Exec()
	if err != nil {
		fmt.Printf("Delete err: %+v\n", err)
		return
	}
	fmt.Printf("result: %s\n", result)
	table.PutDeleter(deleter)
}
