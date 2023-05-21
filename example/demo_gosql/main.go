package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/example/pbgo"
	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/proto"
	"github.com/j32u4ukh/gosql/proto/gstmt"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const TID byte = 0

var db *database.Database
var gs *gstmt.Gstmt
var sql string
var result *database.SqlResult
var err error
var table *gosql.Table

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

	tableParam, columnParams, err := proto.GetParams("../pb", "Desk", dialect.MARIA)

	if err != nil {
		fmt.Printf("讀取 proto 檔時發生錯誤, err: %+v\n", err)
		return
	}

	table = gosql.NewTable(dc.Name, tableParam, columnParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)

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
		fmt.Printf("err: %+v\n", err)
		return
	}
	fmt.Printf("sql: %s\n", sql)
}

func InsertDemo() {
	insert := table.GetInserter()

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	var i int32
	start := time.Now()

	for i = 0; i < 5; i++ {
		// insert.Insert()
		sql, err = insert.ToStmt()
		sql, err = gs.Insert(TID, []protoreflect.ProtoMessage{&pbgo.Desk{
			UserName: "2",
			ItemId:   i,
		}})

		if err != nil {
			fmt.Printf("Error: %+v", err)
			return
		}

		result, err = db.Exec(sql)

		if err != nil {
			fmt.Printf("Insert Exec err: %+v\n", err)
			return
		}

		// fmt.Printf("result: %s\n", result)
	}

	fmt.Printf("Cost time: %+v\n", time.Since(start))
}

func QueryDemo() {
	sql, err = gs.CreateTable(0, "../pb", "Desk")

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	start := time.Now()
	where := gdo.WS().Ne("index", -1)
	sql, err = gs.Query(TID, where)

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	fmt.Printf("QueryDemo | sql: %s\n", sql)
	result, err = db.Query(sql)

	if err != nil {
		fmt.Printf("Query Exec err: %+v\n", err)
		return
	}
	fmt.Printf("result: %s\n", result)
	fmt.Printf("Cost time: %+v, count: %d\n", time.Since(start), result.NRow)
}

func UpdateDemo() {
	sql, err = gs.CreateTable(0, "../pb", "Desk")

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	var i int32
	start := time.Now()
	var desk *pbgo.Desk

	for i = 0; i < 100; i++ {
		desk = &pbgo.Desk{
			Index:    i + 1,
			UserName: "test1",
			ItemId:   i,
		}

		sql, err = gs.Update(TID, desk, gdo.WS().Eq("index", desk.Index))

		if err != nil {
			fmt.Printf("Update item_id(%d) Error: %+v\n", i, err)
			return
		}

		_, err = db.Exec(sql)

		if err != nil {
			fmt.Printf("Update Exec err: %+v\n", err)
			return
		}
	}

	fmt.Printf("Cost time: %+v\n", time.Since(start))
}

func DeleteDemo() {
	sql, err = gs.CreateTable(0, "../pb", "Desk")

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	sql, err = gs.DeleteBy(TID, gdo.WS().Eq("index", 3))

	if err != nil {
		fmt.Printf("Delete Error: %+v\n", err)
		return
	}

	result, err = db.Exec(sql)

	if err != nil {
		fmt.Printf("Create Exec err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}
