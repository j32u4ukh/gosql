package main

import (
	"fmt"
	"os"
	"strconv"
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

var db *database.Database
var gs *gstmt.Gstmt
var sql string
var result *database.SqlResult
var err error

const GID int32 = 0
const TID1 byte = 1
const TID2 byte = 2
const tableName1 string = "AllMight1"
const tableName2 string = "AllMight2"

var TID int = 1
var logger *glog.Logger

func main() {
	logger = glog.SetLogger(0, "demo_hero", glog.DebugLevel)
	utils.SetLogger(logger)
	logger.SetFolder("../log")

	var command string

	if len(os.Args) >= 2 {
		TID, err = strconv.Atoi(os.Args[1])

		if err != nil {
			TID = 1
		}

		if len(os.Args) >= 3 {
			command = strings.ToLower(os.Args[2])
		} else {
			command = "q"
		}

	} else {
		TID = 1
		command = "q"
	}

	conf, err := database.NewConfig("../config/config.yaml")

	if err != nil {
		fmt.Printf("讀取 Config 檔時發生錯誤, err: %+v\n", err)
		return
	}

	dc := conf.GetDatabase()

	db, err = database.Connect(GID, dc.UserName, dc.Password, dc.Server, dc.Port, dc.Name)
	if err != nil {
		fmt.Printf("與資料庫連線時發生錯誤, err: %+v\n", err)
		return
	}
	defer db.Close()

	db = database.Get(GID)

	if db == nil {
		fmt.Println("Database(0) is not exists.")
		return
	}

	_, err = gstmt.SetGstmt(0, dc.Name, dialect.MARIA)

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
	if TID == 1 {
		sql, err = gs.CreateTable(TID1, "../pb", tableName1)

		if err != nil {
			fmt.Printf("Create err: %+v\n", err)
			return
		}

	} else {
		sql, err = gs.CreateTable(TID2, "../pb", tableName2)

		if err != nil {
			fmt.Printf("Create err: %+v\n", err)
			return
		}
	}

	result, err = db.Exec(sql)

	if err != nil {
		fmt.Printf("Create Exec err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
	gs.UseAntiInjection(true)
}

// TID1: 125000 筆/51.8884014s | 2409 筆/s
// TID2: 125000 筆/53.9523073s | 2316 筆/s
func InsertDemo() {
	start := time.Now()
	result = &database.SqlResult{}
	var temp *database.SqlResult

	if TID == 1 {
		sql, err = gs.CreateTable(TID1, "../pb", tableName1)

		if err != nil {
			fmt.Printf("Create err: %+v\n", err)
			return
		}

		var i int64
		for i = 0; i < 100; i++ {
			am1 := &pbgo.AllMight1{
				Bi: i,
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

			sql, err = gs.Insert(TID1, []protoreflect.ProtoMessage{am1})

			if err != nil {
				fmt.Printf("Error: %+v", err)
				return
			}

			temp, err = db.Exec(sql)

			if err != nil {
				fmt.Printf("Insert Exec err: %+v\n", err)
				return
			}

			result.Merge(temp)
		}

	} else {
		sql, err = gs.CreateTable(TID2, "../pb", tableName2)

		if err != nil {
			fmt.Printf("Create err: %+v\n", err)
			return
		}

		var i int32
		for i = 0; i < 100; i++ {

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
				Uti: 18,
				Usi: 19,
				Umi: 20,
				Ui:  i,
			}

			sql, err = gs.Insert(TID2, []protoreflect.ProtoMessage{am2})

			if err != nil {
				fmt.Printf("Error: %+v", err)
			}

			temp, err = db.Exec(sql)
			result.Merge(temp)

			if err != nil {
				fmt.Printf("Insert Exec err: %+v\n", err)
				return
			}
		}
	}

	fmt.Printf("result: %s\n", result)
	fmt.Printf("Cost time: %+v\n", time.Since(start))
}

func QueryDemo() {
	start := time.Now()

	if TID == 1 {
		sql, err = gs.CreateTable(TID1, "../pb", tableName1)

		if err != nil {
			fmt.Printf("Create err: %+v\n", err)
			return
		}

		where := gdo.WS().Ne("bi", -1)
		sql, err = gs.Query(TID1, where)

		if err != nil {
			fmt.Printf("Error: %+v\n", err)
		}

		fmt.Printf("QueryDemo | sql: %s\n", sql)
		result, err = db.Query(sql)

		if err != nil {
			fmt.Printf("Query Exec err: %+v\n", err)
			return
		}
	} else {
		sql, err = gs.CreateTable(TID2, "../pb", tableName2)

		if err != nil {
			fmt.Printf("Create err: %+v\n", err)
			return
		}

		where := gdo.WS().Ne("ui", -1)
		sql, err = gs.Query(TID2, where)

		if err != nil {
			fmt.Printf("Error: %+v\n", err)
		}

		fmt.Printf("QueryDemo | sql: %s\n", sql)
		result, err = db.Query(sql)

		if err != nil {
			fmt.Printf("Query Exec err: %+v\n", err)
			return
		}
	}

	fmt.Printf("result: %s\n", result)
	fmt.Printf("Cost time: %+v, count: %d\n", time.Since(start), result.NRow)
}

// TID1: 125000 筆/58.472945s | 2137 筆/s
// TID2: 125000 筆/61.019339s | 2048 筆/s
func UpdateDemo() {
	start := time.Now()
	result = &database.SqlResult{}
	var temp *database.SqlResult

	if TID == 1 {
		sql, err = gs.CreateTable(TID1, "../pb", tableName1)

		if err != nil {
			fmt.Printf("Create err: %+v\n", err)
			return
		}

		var am1 *pbgo.AllMight1
		var i int64
		for i = 0; i < 100; i++ {
			am1 = &pbgo.AllMight1{
				Bi: i,
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

			sql, err = gs.Update(TID1, am1, gdo.WS().Eq("bi", i))

			if err != nil {
				fmt.Printf("Update bi(%d) Error: %+v\n", i, err)
				return
			}

			temp, err = db.Exec(sql)
			result.Merge(temp)

			if err != nil {
				fmt.Printf("Update Exec err: %+v\n", err)
				return
			}
		}

	} else {
		sql, err = gs.CreateTable(TID2, "../pb", tableName2)

		if err != nil {
			fmt.Printf("Create err: %+v\n", err)
			return
		}

		var am2 *pbgo.AllMight2
		var i int32
		for i = 0; i < 100; i++ {
			am2 = &pbgo.AllMight2{
				C: &pbgo.Currency{
					PlayerId: 1,
					Diamond:  2,
					Gold:     3,
				},
				Mii: map[int32]int32{1: 2, 3: 4},
				Mis: map[int32]string{1: "a", 2: "b"},
				Msi: map[string]int32{"c": 3, "d": 4},
				Mss: map[string]string{"e": "f", "g": fmt.Sprintf("h%d", i)},
				Uti: 18,
				Usi: 19,
				Umi: 20,
				Ui:  i,
			}

			sql, err = gs.Update(TID2, am2, gdo.WS().Eq("ui", i))

			if err != nil {
				fmt.Printf("Update ui(%d) Error: %+v\n", i, err)
				return
			}

			temp, err = db.Exec(sql)
			result.Merge(temp)

			if err != nil {
				fmt.Printf("Update Exec err: %+v\n", err)
				return
			}
		}
	}

	fmt.Printf("result: %s\n", result)
	fmt.Printf("Cost time: %+v\n", time.Since(start))
}

func DeleteDemo() {
	if TID == 1 {
		sql, err = gs.CreateTable(TID1, "../pb", tableName1)

		if err != nil {
			fmt.Printf("Create err: %+v\n", err)
			return
		}

		sql, err = gs.DeleteBy(TID1, gdo.WS().Eq("bi", 3))

		if err != nil {
			fmt.Printf("Delete Error: %+v\n", err)
			return
		}

		result, err = db.Exec(sql)

		if err != nil {
			fmt.Printf("Create Exec err: %+v\n", err)
			return
		}

	} else {
		sql, err = gs.CreateTable(TID2, "../pb", tableName2)

		if err != nil {
			fmt.Printf("Create err: %+v\n", err)
			return
		}

		sql, err = gs.DeleteBy(TID2, gdo.WS().Eq("ui", 3))

		if err != nil {
			fmt.Printf("Delete Error: %+v\n", err)
			return
		}

		result, err = db.Exec(sql)

		if err != nil {
			fmt.Printf("Create Exec err: %+v\n", err)
			return
		}
	}

	fmt.Printf("result: %s\n", result)
}
