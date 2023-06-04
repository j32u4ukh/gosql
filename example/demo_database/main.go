package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/utils"
)

var db *database.Database

func main() {
	logger := glog.SetLogger(0, "demo_database", glog.DebugLevel)
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

	switch command {
	case "c":
		CreateDemo()
	case "i":
		InsertDemo()
	case "it":
		InsertTestDemo()
	case "is":
		InsertStatementDemo()
	case "q":
		QueryDemo()
	case "u":
		UpdateDemo()
	case "ut":
		UpdateTestDemo()
	case "d":
		DeleteDemo()
	default:
		fmt.Printf("No invalid command(%s).\n", command)
	}
}

func CreateDemo() {
	sql := "CREATE TABLE `Desk` ( `index` INT(11) NOT NULL AUTO_INCREMENT COMMENT '索引值', `user_name` VARCHAR(20) NOT NULL COMMENT '玩家名稱' COLLATE 'utf8mb4_bin',`item_id` INT(11) NOT NULL COMMENT '物品 ID',`time` TIMESTAMP NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Log 建立時間',PRIMARY KEY (`index`) USING BTREE) COLLATE='utf8mb4_bin' ENGINE=InnoDB;"
	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

func InsertDemo() {
	sql := "INSERT INTO `pekomiko`.`Desk` (`user_name`, `item_id`) VALUES ('sss', '23');"
	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

// 71.7612942s / 125000 rows | 574 us/row
// 567.2 µs/row -> 3.6503ms / 10rows -> 34.0006ms / 100rows -> 81.1401ms/175rows -> 127.2253ms/250rows -> 276.8077ms / 500 rows -> 570.4072ms/1000rows
// 567.2 µs/row -> 0.36503 ms/row -> 0.340006 ms/row -> 0.46365771428 ms/row -> 0.5089012 ms/row -> 0.5536154 ms/row -> 0.5704072 ms/row
func InsertTestDemo() {
	start := time.Now()
	var sql string
	var err error
	var i int32
	sql = "INSERT INTO `pekomiko`.`Desk` (`user_name`, `item_id`) VALUES ('sss', '23');"

	for i = 0; i < 175; i++ {
		_, err = db.Exec(sql)

		if err != nil {
			fmt.Printf("Insert err: %+v\n", err)
			return
		}
	}

	fmt.Printf("Cost time: %+v\n", time.Since(start))
}

func InsertStatementDemo() {
	start := time.Now()
	var sql string
	var i int32
	var num int32 = 1000000

	for i = 0; i < num; i++ {
		sql = fmt.Sprintf("INSERT INTO `pekomiko`.`Desk` (`user_name`, `item_id`) VALUES ('sss', '%d');", i)
		if sql == "" {
			return
		}
	}

	cost := time.Since(start)
	fmt.Printf("Cost time: %+v, %+v us/row\n", cost, float64(cost)/float64(num))
}

func QueryDemo() {
	sql := "SELECT `index`, `user_name`, `item_id`, `time` FROM `pekomiko`.`Desk`;"
	result, err := db.Query(sql)

	if err != nil {
		fmt.Printf("Select err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
	for i, data := range result.Datas {
		fmt.Printf("i: %d, data: %s\n", i, cntr.SliceToString(data))
	}
}

func UpdateDemo() {
	sql := "UPDATE `pekomiko`.`Desk` SET `item_id`='29' WHERE `index`=1;"
	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

// 70.6651062s / 125000 rows | 565 us/row
func UpdateTestDemo() {
	start := time.Now()
	var sql string
	var err error
	var i int32

	for i = 0; i < 125000; i++ {
		sql = fmt.Sprintf("UPDATE `pekomiko`.`Desk` SET `item_id`='30' WHERE `index`=%d;", i+1)
		_, err = db.Exec(sql)

		if err != nil {
			fmt.Printf("Insert err: %+v\n", err)
			return
		}
	}

	fmt.Printf("Cost time: %+v\n", time.Since(start))
}

func DeleteDemo() {
	sql := "DELETE FROM `pekomiko`.`Desk` WHERE `index`=2;"
	result, err := db.Exec(sql)

	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}
