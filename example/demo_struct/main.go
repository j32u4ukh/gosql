package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/plugin"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
)

type Tsukue struct {
	Id      int    `gorm:"pk=default;default=ai"`
	Content string `gorm:"size=3000"`
}

const TID byte = 0

var db *database.Database
var result *database.SqlResult
var table *gosql.Table
var sql string
var err error

func main() {
	command := strings.ToLower(os.Args[1])
	conf, err := database.NewConfig("../config/config.yaml")

	if err != nil {
		fmt.Printf("讀取 Config 檔時發生錯誤, err: %+v\n", err)
		return
	}

	dc := conf.GetDatabase()
	db, err = database.Connect(0, dc.UserName, dc.Password, dc.Server, dc.Port, dc.Name)

	if err != nil {
		fmt.Printf("與資料庫連線時發生錯誤, err: %+v\n", err)
		return
	}

	defer db.Close()
	db = database.Get(0)

	if db == nil {
		fmt.Println("Database(0) is not exists.")
		return
	}

	desk := &Tsukue{}
	tableParams, columnParams, err := plugin.GetStructParams(desk, dialect.MARIA)
	table = gosql.NewTable("Desk", tableParams, columnParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	table.Init(&gosql.TableConfig{
		Db:               db,
		DbName:           dc.Name,
		UseAntiInjection: false,
		InsertFunc:       plugin.InsertStruct,
		UpdateAnyFunc:    plugin.UpdateStruct,
	})
	if err != nil {
		fmt.Printf("BuildCreateStmt err: %+v\n", err)
		return
	}

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

/*
CREATE TABLE IF NOT EXISTS `pekomiko`.`Desk` (
	`Id` INT(11) NOT NULL DEFAULT 0,
	`Content` VARCHAR(3000) NOT NULL DEFAULT '' COLLATE 'utf8mb4_bin',
	PRIMARY KEY (`Id`) USING BTREE
) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';
*/
func CreateDemo() {
	sql, err = table.Creater().ToStmt()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}
	fmt.Printf("sql: %s\n", sql)
	result, err = table.Creater().Exec()

	if err != nil {
		fmt.Printf("Creater err: %+v\n", err)
		return
	}

	fmt.Printf("result: %+v\n", result)
}

// INSERT INTO `pekomiko`.`Desk` (`Id`, `Content`) VALUES (0, 'abc');
func InsertDemo() {
	inserter := table.GetInserter()
	desk := &Tsukue{Id: 0, Content: "abc"}
	err = inserter.Insert(desk)

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	sql, err = inserter.ToStmt()

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	fmt.Printf("sql: %s\n", sql)
	result, err = inserter.Exec()

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	fmt.Printf("result: %+v\n", result)
	table.PutInserter(inserter)
}

func QueryDemo() {
	selector := table.GetSelector()
	where := gosql.WS().Ne("Id", -1)
	selector.SetCondition(where)
	sql, err = selector.ToStmt()

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	fmt.Printf("QueryDemo | sql: %s\n", sql)
	result, err = selector.Exec()

	if err != nil {
		fmt.Printf("Query err: %+v\n", err)
		return
	}

	fmt.Printf("result: %+v\n", result)
	table.PutSelector(selector)
}

func UpdateDemo() {
	updater := table.GetUpdater()
	where := gosql.WS().Eq("Id", 2)
	updater.Update("Content", "def")
	updater.SetCondition(where)
	sql, err = updater.ToStmt()

	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}

	fmt.Printf("sql: %s\n", sql)
	result, err = updater.Exec()

	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}

	fmt.Printf("result: %+v\n", result)
	table.PutUpdater(updater)
}

func DeleteDemo() {
	deleter := table.GetDeleter()
	where := gosql.WS().Eq("Id", 0)
	deleter.SetCondition(where)
	sql, err = deleter.ToStmt()

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	fmt.Printf("sql: %s\n", sql)
	result, err = deleter.Exec()

	if err != nil {
		fmt.Printf("Delete err: %+v\n", err)
		return
	}

	fmt.Printf("result: %+v\n", result)
	table.PutDeleter(deleter)
}
