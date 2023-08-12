package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

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
	logger := glog.SetLogger(0, "demo_pekomiko", glog.InfoLevel)
	logger.SetSkip(3)
	utils.SetLogger(logger)
	conf, err := database.NewConfig("../config/config.yaml")
	if err != nil {
		utils.Error("讀取 Config 檔時發生錯誤, err: %+v", err)
		return
	}
	dc := conf.GetDatabase()
	db, err := database.Connect(0, dc.User, dc.Password, dc.Host, dc.Port, dc.DbName)
	if err != nil {
		utils.Error("與資料庫連線時發生錯誤, err: %+v", err)
		return
	}
	defer db.Close()
	if db == nil {
		utils.Error("Database(0) is not exists.")
		return
	}
	tableName := "PostMessage"
	tableParams, columnParams, err := plugin.GetProtoParams(fmt.Sprintf("../pb/%s.proto", tableName), dialect.MARIA)
	if err != nil {
		utils.Error("Failed to get proto params, err: %+v", err)
		return
	}
	utils.Debug("tableParams: %+v", tableParams)
	for _, columnParam := range columnParams {
		utils.Debug("columnParam: %+v", columnParam)
	}
	table = gosql.NewTable(tableName, tableParams, columnParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	table.Init(&gosql.TableConfig{
		Db:               db,
		DbName:           dc.DbName,
		UseAntiInjection: true,
		PtrToDbFunc:      plugin.ProtoToDb,
		InsertFunc:       plugin.InsertProto,
		QueryFunc:        plugin.QueryProto,
		UpdateAnyFunc:    plugin.UpdateProto,
	})
	if err != nil {
		utils.Error("NewTable err: %+v", err)
		return
	}

	CreateDemo(os.Args)
	InsertDemo(os.Args)
	QueryDemo(os.Args)
	UpdateDemo(os.Args)
	DeleteDemo(os.Args)
}

func CreateDemo(args []string) {
	if len(args) < 2 {
		return
	} else if args[1] != "create" {
		return
	}
	sql, err = table.Creater().ToStmt()
	if err != nil {
		utils.Error("Create err: %+v", err)
		return
	}
	utils.Info("sql: %s", sql)
	result, err = table.Creater().Exec()
	if err != nil {
		utils.Error("Create Exec err: %+v", err)
		return
	}
	utils.Info("result: %s", result)
}

func InsertDemo(args []string) {
	if len(args) < 2 {
		return
	} else if args[1] != "insert" {
		return
	}
	inserter := table.GetInserter()
	defer table.PutInserter(inserter)
	now := time.Now()
	pm := &pbgo.PostMessage{
		Id:      uint64(now.UnixNano()),
		UserId:  int32(now.Unix()),
		Content: fmt.Sprintf("now: %s", now),
	}
	utils.Info("PostMessage: %+v", pm)
	err = inserter.Insert(pm)
	if err != nil {
		utils.Error("Insert err: %+v", err)
		return
	}
	sql, err = inserter.ToStmt()
	if err != nil {
		utils.Error("ToStmt err: %+v", err)
		return
	}
	utils.Info("sql: %s", sql)
	result, err = inserter.Exec()
	if err != nil {
		utils.Error("Insert exec err: %+v", err)
		return
	}
	utils.Info("result: %s", result)
}

func QueryDemo(args []string) {
	if len(args) < 3 {
		return
	} else if args[1] != "query" {
		return
	}
	id, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		utils.Error("args[2]: %s 數值有誤, err: %+v", args[2], err)
		return
	}
	selector := table.GetSelector()
	defer table.PutSelector(selector)
	selector.SetCondition(gosql.WS().Eq("id", id))
	sql, err = selector.ToStmt()
	if err != nil {
		utils.Error("Failed to create stmt, err: %+v", err)
		return
	}
	utils.Info("sql: %s", sql)
	pms, err := selector.Query(func() any { return &pbgo.PostMessage{} })
	if err != nil {
		utils.Error("Failed to query post messages, err: %+v", err)
	}
	for _, pm := range pms {
		utils.Info("pm: %+v", pm)
	}
}

func UpdateDemo(args []string) {
	if len(args) < 3 {
		return
	} else if args[1] != "update" {
		return
	}
	id, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		utils.Error("args[2]: %s 數值有誤, err: %+v", args[2], err)
		return
	}
	selector := table.GetSelector()
	defer table.PutSelector(selector)
	selector.SetCondition(gosql.WS().Eq("id", id))
	pms, err := selector.Query(func() any { return &pbgo.PostMessage{} })
	if err != nil {
		utils.Error("Failed to query post messages, err: %+v", err)
		return
	}
	if len(pms) == 0 {
		utils.Error("No post messages found")
		return
	}
	updater := table.GetUpdater()
	defer table.PutUpdater(updater)
	now := time.Now()
	pm := pms[0].(*pbgo.PostMessage)
	pm.UserId = int32(now.Unix())
	pm.Content = fmt.Sprintf("now: %s", now)
	updater.UpdateAny(pm)
	updater.SetCondition(gosql.WS().Eq("id", id))
	result, err = updater.Exec()
	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}
	fmt.Printf("result: %+v\n", result)
}

func DeleteDemo(args []string) {
	if len(args) < 3 {
		return
	} else if args[1] != "delete" {
		return
	}
	id, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		utils.Error("args[2]: %s 數值有誤, err: %+v", args[2], err)
		return
	}
	deleter := table.GetDeleter()
	defer table.PutDeleter(deleter)
	deleter.SetCondition(gosql.WS().Eq("id", id))
	sql, err = deleter.ToStmt()
	if err != nil {
		utils.Error("Delete err: %+v", err)
		return
	}
	utils.Info("sql: %s", sql)
	result, err = deleter.Exec()
	if err != nil {
		utils.Error("Delete err: %+v", err)
		return
	}
	utils.Info("result: %+v", result)
}
