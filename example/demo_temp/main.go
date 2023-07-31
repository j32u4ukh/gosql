package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/example/pbgo"
	"github.com/j32u4ukh/gosql/plugin"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/j32u4ukh/gosql/utils"
)

var table *gosql.Table
var sql string
var result *database.SqlResult
var err error

func main() {
	logger := glog.SetLogger(0, "demo_temp", glog.DebugLevel)
	utils.SetLogger(logger)
	command := strings.ToLower(os.Args[1])

	switch command {
	case "c", "i", "q", "u", "d":
	default:
		tableName := "Edge"
		tableParams, columnParams, err := plugin.GetProtoParams(fmt.Sprintf("../pb/%s.proto", tableName), dialect.MARIA)

		if err != nil {
			fmt.Printf("Failed to get proto params, err: %+v\n", err)
			return
		}

		utils.Debug("tableParams: %+v", tableParams)
		for _, col := range columnParams {
			utils.Debug(" %+v", col)
		}
		return
	}

	conf, err := database.NewConfig("../config/config.yaml")
	if err != nil {
		fmt.Printf("讀取 Config 檔時發生錯誤, err: %+v\n", err)
		return
	}
	dc := conf.GetDatabase()
	db, err := database.Connect(0, dc.User, dc.Password, dc.Host, dc.Port, dc.DbName)
	if err != nil {
		fmt.Printf("與資料庫連線時發生錯誤, err: %+v\n", err)
		return
	}
	defer db.Close()
	if db == nil {
		fmt.Println("Database(0) is not exists.")
		return
	}

	tableName := "Edge"
	tableParams, columnParams, err := plugin.GetProtoParams(fmt.Sprintf("../pb/%s.proto", tableName), dialect.MARIA)

	if err != nil {
		fmt.Printf("Failed to get proto params, err: %+v\n", err)
		return
	}

	table = gosql.NewTable(tableName, tableParams, columnParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	table.Init(&gosql.TableConfig{
		Db:               db,
		DbName:           dc.DbName,
		UseAntiInjection: false,
		PtrToDbFunc:      plugin.ProtoToDb,
		InsertFunc:       plugin.InsertProto,
		QueryFunc:        plugin.QueryProto,
		UpdateAnyFunc:    plugin.UpdateProto,
	})

	if err != nil {
		utils.Error("NewTable err: %+v\n", err)
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

func CreateDemo() {
	sql, err = table.Creater().ToStmt()

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	fmt.Printf("sql: %s\n", sql)

	result, err = table.Creater().Exec()

	if err != nil {
		fmt.Printf("Create Exec err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func InsertDemo() {
	inserter := table.GetInserter()
	defer table.PutInserter(inserter)
	edge := &pbgo.Edge{
		UserId: 97,
		Target: 31,
	}
	err = inserter.Insert(edge)

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	sql, err = inserter.ToStmt()

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}
	result, err = inserter.Exec()

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	fmt.Printf("result: %+v\n", result)
}

func QueryDemo() {
	selector := table.GetSelector()
	defer table.PutSelector(selector)
	sql, err = selector.ToStmt()

	if err != nil {
		utils.Error("Error: %+v\n", err)
	}

	utils.Debug("sql: %s\n", sql)
	edges, err := selector.Query(func() any { return &pbgo.Edge{} })

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	for _, edge := range edges {
		fmt.Printf("edge: %+v\n", edge)
	}
}

func UpdateDemo() {
	updater := table.GetUpdater()
	defer table.PutUpdater(updater)
	edge := &pbgo.Edge{
		Index:  1,
		UserId: 97,
		Target: 23,
	}
	updater.UpdateAny(edge)
	where := gosql.WS().Eq("index", 1)
	updater.SetCondition(where)
	result, err = updater.Exec()

	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}
	utils.Info("Updated edge: %+v", edge)
}

func DeleteDemo() {
	deleter := table.GetDeleter()
	defer table.PutDeleter(deleter)
	deleter.SetCondition(gosql.WS().Eq("index", 1))
	sql, err = deleter.ToStmt()

	if err != nil {
		fmt.Printf("Delete err: %+v\n", err)
		return
	}

	fmt.Printf("sql: %s\n", sql)
	result, err = deleter.Exec()

	if err != nil {
		fmt.Printf("Delete err: %+v\n", err)
		return
	}

	fmt.Printf("result: %+v\n", result)
}
