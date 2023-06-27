package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/example/pbgo"
	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/proto/gstmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/j32u4ukh/gosql/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const TID byte = 0

var db *database.Database
var gs *gstmt.Gstmt
var sql string
var result *database.SqlResult
var err error
var logger *glog.Logger

func main() {
	logger = glog.SetLogger(0, "demo_gstmt", glog.DebugLevel)
	utils.SetLogger(logger)
	logger.SetFolder("../log")

	command := strings.ToLower(os.Args[1])
	conf, err := database.NewConfig("../config/config.yaml")

	if err != nil {
		fmt.Printf("讀取 Config 檔時發生錯誤, err: %+v\n", err)
		return
	}

	dc := conf.GetDatabase()
	db, err = database.Connect(0, dc.User, dc.Password, dc.Host, dc.Port, dc.DbName)

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

	_, err = gstmt.SetGstmt(0, dc.DbName, dialect.MARIA)

	if err != nil {
		fmt.Printf("SetGstmt err: %+v\n", err)
		return
	}

	gs, err = gstmt.GetGstmt(0)

	if err != nil {
		logger.Error("Failed to get GoSql, err %+v", err)
		return
	}

	switch command {
	case "c":
		CreateDemo()
	case "i":
		InsertDemo()
	case "is":
		InsertStatementDemo()
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
	sql, err = gs.CreateTable(0, "../pb", "Desk")

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	result, err = db.Exec(sql)

	if err != nil {
		fmt.Printf("Create Exec err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
	gs.UseAntiInjection(false)
}

// 125000 筆/50.7553501s | 2462 筆/s -> 125000 筆/50.6376619s | 2468 筆/s
// 70.8177351s / 125000 rows -> 567 us/row
// 1.033ms/1row -> 7.7917ms/10rows -> 55.2513ms/100rows -> 131.0093ms/250rows -> 275.1029ms/500rows -> 556.0053ms/1000rows
// 1.033 ms/row -> 0.77917 ms/row -> 0.552513 ms/row -> 0.5240372 ms/row -> 0.5502058 ms/row -> 0.5560053 ms/row
func InsertDemo() {
	sql, err = gs.CreateTable(0, "../pb", "Desk")

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	var i int32
	start := time.Now()

	for i = 0; i < 100; i++ {
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

func InsertStatementDemo() {
	sql, err = gs.CreateTable(0, "../pb", "Desk")

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	var i int32
	var num int32 = 1000000
	start := time.Now()

	for i = 0; i < num; i++ {
		sql, err = gs.Insert(TID, []protoreflect.ProtoMessage{&pbgo.Desk{
			UserName: "2",
			ItemId:   i,
		}})

		if err != nil {
			fmt.Printf("Error: %+v", err)
			return
		}
	}

	cost := time.Since(start)
	fmt.Printf("Cost time: %+v, %+v us/row\n", cost, float64(cost)/float64(num))
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

// 125000 筆/55.3019832s | 2260 筆/s
// 71.9243555s / 125000 rows | 575 us/row
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

// ===== Without primary key =====
// total 100 rows, update 100 rows
// Update		151.3475ms
// FastUpdate 	121.3799ms
// total 125000 rows, update 100 rows
// Update		3.2189428s
// FastUpdate 	51.9154ms
// total 125000 rows, update 125000 rows
// Update		-
// FastUpdate 	50.3950989s
// ===== With primary key =====
// total 125000 rows, update 100 rows
// Update		39.2271ms
// FastUpdate 	40.8054ms
// total 125000 rows, update 125000 rows
// Update		54.6248308s
// FastUpdate 	50.2750576s
