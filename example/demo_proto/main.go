package main

import (
	"fmt"
	"os"
	"strconv"
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
var tid int

func main() {
	logger := glog.SetLogger(0, "demo_proto", glog.DebugLevel)
	utils.SetLogger(logger)

	tid, err = strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Printf("tid 應為數字 %s, err: %+v\n", os.Args[1], err)
		return
	}

	command := strings.ToLower(os.Args[2])
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

	tableName := fmt.Sprintf("AllMight%d", tid)
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
	var am1 *pbgo.AllMight1
	var am2 *pbgo.AllMight2
	var i int32

	for i = 0; i < 10; i++ {
		if tid == 1 {
			am1 = &pbgo.AllMight1{
				Bi: int64(i),
				B:  true,
				Ti: -3,
				Si: -4,
				Mi: -5,
				I:  -6,
				Tt: "tt",
				Vc: "vc",
				T:  "t",
				Mt: "mt",
				Lt: "lt",
				Ts: &pbgo.TimeStamp{
					Year:   2023,
					Month:  1,
					Day:    12,
					Hour:   10,
					Minute: 10,
					Second: 30,
				},
			}

			err = inserter.Insert(am1)

		} else {
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
		}

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
	var where *gosql.WhereStmt

	if tid == 1 {
		where = gosql.WS().Ne("bi", -1)
	} else {
		where = gosql.WS().Ne("ui", -1)
	}

	selector.SetCondition(where)
	sql, err = selector.ToStmt()

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	fmt.Printf("QueryDemo | sql: %s\n", sql)

	result, _ = selector.Exec()

	fmt.Printf("QueryDemo | result: %+v\n", result.Datas)

	if tid == 1 {

		am1s, err := selector.Query(func() any { return &pbgo.AllMight1{} })

		if err != nil {
			fmt.Printf("Error: %+v\n", err)
		}

		for _, am1 := range am1s {
			fmt.Printf("am1: %+v\n", am1)
		}
	} else {

		am2s, err := selector.Query(func() any { return &pbgo.AllMight2{} })

		if err != nil {
			fmt.Printf("Error: %+v\n", err)
		}

		for _, am2 := range am2s {
			fmt.Printf("am2: %+v\n", am2)
		}
	}

	table.PutSelector(selector)
}

func UpdateDemo() {
	updater := table.GetUpdater()
	result = &database.SqlResult{}
	var where *gosql.WhereStmt
	var temp *database.SqlResult

	for i := 0; i < 10; i++ {
		if tid == 1 {
			am1 := &pbgo.AllMight1{
				Bi: int64(i),
				B:  true,
				Ti: -3,
				Si: -4,
				Mi: -5,
				I:  -6,
				Tt: fmt.Sprintf("tt%d", i),
				Vc: "vc",
				T:  "t",
				Mt: "mt",
				Lt: "lt",
				Ts: &pbgo.TimeStamp{
					Year:   2023,
					Month:  1,
					Day:    12,
					Hour:   10,
					Minute: 10,
					Second: 30,
				},
			}

			updater.UpdateAny(am1)
			where = gosql.WS().Eq("bi", i)
			updater.SetCondition(where)
			temp, err = updater.Exec()
			result.Merge(temp)

			if err != nil {
				fmt.Printf("Update err: %+v\n", err)
				return
			}
		} else {
			am2 := &pbgo.AllMight2{
				C: &pbgo.Currency{
					PlayerId: 1,
					Diamond:  2,
					Gold:     3,
				},
				Mii: map[int32]int32{1: 2, 3: 4},
				Mis: map[int32]string{1: "a", 2: "b"},
				Msi: map[string]int32{"c": 3, "d": 4},
				Mss: map[string]string{"e": "fgo", "g": "h"},
				Uti: 81,
				Usi: 91,
				Umi: 102,
				Ui:  int32(i),
			}

			updater.UpdateAny(am2)
			where = gosql.WS().Eq("ui", int32(i))
			updater.SetCondition(where)
			temp, err = updater.Exec()
			result.Merge(temp)

			if err != nil {
				fmt.Printf("Update err: %+v\n", err)
				return
			}
		}
	}

	fmt.Printf("result: %+v\n", result)
	table.PutUpdater(updater)
}

func DeleteDemo() {
	deleter := table.GetDeleter()
	var where *gosql.WhereStmt
	if tid == 1 {
		where = gosql.WS().Eq("bi", 3)
	} else {
		where = gosql.WS().Eq("ui", 3)
	}

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
