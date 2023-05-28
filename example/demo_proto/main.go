package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/example/pbgo"
	"github.com/j32u4ukh/gosql/plugin"
	"github.com/j32u4ukh/gosql/proto/gstmt"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
)

var table *gosql.Table
var sql string
var result *database.SqlResult
var err error

func main() {
	command := strings.ToLower(os.Args[1])
	conf, err := database.NewConfig("../config/config.yaml")

	if err != nil {
		fmt.Printf("讀取 Config 檔時發生錯誤, err: %+v\n", err)
		return
	}

	dc := conf.GetDatabase()
	db, err := database.Connect(0, dc.UserName, dc.Password, dc.Server, dc.Port, dc.Name)

	if err != nil {
		fmt.Printf("與資料庫連線時發生錯誤, err: %+v\n", err)
		return
	}
	defer db.Close()

	if db == nil {
		fmt.Println("Database(0) is not exists.")
		return
	}

	_, err = gstmt.SetGstmt(0, dc.Name, dialect.MARIA)

	if err != nil {
		fmt.Printf("SetGstmt err: %+v\n", err)
		return
	}

	tableParams, columnParams, err := plugin.GetProtoParams("../pb/AllMight2.proto", dialect.MARIA)

	if err != nil {
		fmt.Printf("Failed to get proto params, err: %+v\n", err)
		return
	}

	table = gosql.NewTable("AllMight2", tableParams, columnParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	table.Init(&gosql.TableConfig{
		Db:               db,
		DbName:           dc.Name,
		UseAntiInjection: false,
		PtrToDbFunc:      plugin.ProtoToDb,
		InsertFunc:       plugin.InsertProto,
		QueryFunc:        plugin.QueryProto,
		UpdateAnyFunc:    plugin.UpdateProto,
	})

	if err != nil {
		fmt.Printf("NewTable err: %+v\n", err)
		return
	}

	// fmt.Printf("table: %+v\n", table)

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

// [[ {"1":2,"3":4} {"1":"a","2":"b"} {"c":3,"d":4} {"e":"f","g":"h"} 18 19 20 3]]
func Test() {
	datas := [][]string{
		{
			"{\"player_id\":1,\"diamond\":2,\"gold\":3}",
			"{\"1\":2,\"3\":4}",
			"{\"1\":\"a\",\"2\":\"b\"}",
			"{\"c\":3,\"d\":4}",
			"{\"e\":\"f\",\"g\":\"h\"}",
			"18",
			"19",
			"20",
			"3",
		},
	}

	am2s, err := plugin.QueryProto(datas, func() any { return &pbgo.AllMight2{} })

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return
	}

	for _, am2 := range am2s {
		fmt.Printf("am2: %+v\n", am2)
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
	var am2 *pbgo.AllMight2
	var i int32

	for i = 0; i < 10; i++ {
		am2 = &pbgo.AllMight2{
			C: &pbgo.Currency{
				PlayerId: 1,
				Diamond:  2,
				Gold:     3,
			},
			Mii: map[int32]int32{1: 2, 3: 4},
			Mis: map[int32]string{1: "a", 2: "b"},
			Msi: map[string]int32{"c": 3, "d": 4},
			Mss: map[string]string{"e": "f", "g": "h"},
			Uti: 18,
			Usi: 19,
			Umi: 20,
			Ui:  i,
		}

		err = inserter.Insert(am2)

		if err != nil {
			fmt.Printf("Insert err: %+v\n", err)
			return
		}

	}

	sql, err = inserter.ToStmt()

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	// fmt.Printf("sql: %s\n", sql)

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
	// where := gosql.WS().Ne("ui", 100)
	where := gosql.WS().Eq("ui", 3)
	selector.SetCondition(where)
	sql, err = selector.ToStmt()

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	fmt.Printf("QueryDemo | sql: %s\n", sql)

	result, _ = selector.Exec()

	fmt.Printf("QueryDemo | result: %+v\n", result.Datas)

	am2s, err := selector.Query(func() any { return &pbgo.AllMight2{} })

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	for _, am2 := range am2s {
		fmt.Printf("am2: %+v\n", am2)
	}

	table.PutSelector(selector)
}

func UpdateDemo() {
	updater := table.GetUpdater()
	am2 := &pbgo.AllMight2{
		C: &pbgo.Currency{
			PlayerId: 1,
			Diamond:  2,
			Gold:     3,
		},
		Mii: map[int32]int32{1: 2, 3: 4},
		Mis: map[int32]string{1: "a", 2: "b"},
		Msi: map[string]int32{"c": 3, "d": 4},
		Mss: map[string]string{"e": "f", "g": "h"},
		Uti: 81,
		Usi: 91,
		Umi: 102,
		Ui:  3,
	}
	updater.UpdateAny(am2)
	where := gosql.WS().Eq("ui", 3)
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
	where := gosql.WS().Eq("ui", 1)
	deleter.SetCondition(where)
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
	table.PutDeleter(deleter)
}
