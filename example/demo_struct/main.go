package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/plugin"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/j32u4ukh/gosql/utils"
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
	logger := glog.SetLogger(0, "demo_struct", glog.DebugLevel)
	utils.SetLogger(logger)

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
		QueryFunc:        plugin.QueryStruct,
	})
	if err != nil {
		fmt.Printf("NewTable err: %+v\n", err)
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
	case "t":
		Test()
	default:
		fmt.Printf("No invalid command(%s).\n", command)
	}
}

func Test() {
	datas := [][]string{
		{"1", "abc"},
		{"2", "def"},
		{"3", "ghi"},
	}
	// desks := plugin.Make(3, func() any { return &Tsukue{} })
	desks, err := plugin.QueryStruct(datas, func() any { return &Tsukue{} })
	if err != nil {
		fmt.Printf("Failed to QueryStructFunc.")
		return
	}
	for _, desk := range desks {
		fmt.Printf("desk: %+v\n", desk)
	}
}

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

func InsertDemo() {
	inserter := table.GetInserter()
	var desk *Tsukue

	for i := 0; i < 10; i++ {
		desk = &Tsukue{Content: "abc"}
		err = inserter.Insert(desk)

		if err != nil {
			fmt.Printf("Insert %+v failed, err: %+v\n", desk, err)
			return
		}
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

	desks, err := selector.Query(func() any { return &Tsukue{} })

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	for _, desk := range desks {
		fmt.Printf("desk: %+v\n", desk)
	}

	table.PutSelector(selector)
}

func UpdateDemo() {
	updater := table.GetUpdater()
	desk := &Tsukue{
		Id:      2,
		Content: "xyz",
	}
	updater.UpdateAny(desk)
	where := gosql.WS().Eq("Id", 2)
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
	where := gosql.WS().Eq("Id", 1)
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
