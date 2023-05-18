package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/j32u4ukh/glog"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/j32u4ukh/gosql/stmt/gosql"
)

const TID byte = 0

var db *database.Database
var sql string
var result *database.SqlResult
var err error
var logger *glog.Logger
var table *stmt.Table
var gTable *gosql.Table

func main() {
	command := strings.ToLower(os.Args[1])
	logger = glog.GetLogger("../log", "demo_stmt", glog.DebugLevel, false)
	conf, err := database.NewConfig("../config/config.yaml")

	if err != nil {
		logger.Error("讀取 Config 檔時發生錯誤, err: %+v\n", err)
		return
	}

	dc := conf.GetDatabase()
	db, err = database.Connect(0, dc.UserName, dc.Password, dc.Server, dc.Port, dc.Name)

	if err != nil {
		logger.Error("與資料庫連線時發生錯誤, err: %+v", err)
		return
	}

	defer db.Close()
	db = database.Get(0)

	if db == nil {
		logger.Error("Database(0) is not exists.")
		return
	}

	tableParams := stmt.NewTableParam()
	// fmt.Printf("tableParams: %v\n", tableParams)
	colParam0 := stmt.NewColumnParam(0, "id", datatype.INT, dialect.MARIA)
	colParam0.SetPrimaryKey(stmt.ALGO)
	colParam0.SetDefault("AI")
	colParam1 := stmt.NewColumnParam(1, "name", datatype.VARCHAR, dialect.MARIA)
	colParam1.SetSize(50)
	colParam2 := stmt.NewColumnParam(2, "url", datatype.VARCHAR, dialect.MARIA)
	colParam2.SetSize(50)
	colParam3 := stmt.NewColumnParam(3, "alexa", datatype.INT, dialect.MARIA)
	colParam4 := stmt.NewColumnParam(4, "contury", datatype.VARCHAR, dialect.MARIA)
	colParam4.SetSize(50)
	colParams := []*stmt.ColumnParam{colParam0, colParam1, colParam2, colParam3, colParam4}
	// for i, col := range colParams {
	// 	fmt.Printf("%d) %+v\n", i, col)
	// }
	table = stmt.NewTable("Websites", tableParams, colParams, stmt.ENGINE, stmt.COLLATE)
	gTable = gosql.NewTable("Websites", tableParams, colParams, stmt.ENGINE, stmt.COLLATE)
	gTable.SetDb(db)

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
CREATE TABLE IF NOT EXISTS `Websites` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(3000) NOT NULL DEFAULT '' COLLATE 'utf8mb4_bin',
	`url` VARCHAR(50) NOT NULL DEFAULT '' COLLATE 'utf8mb4_bin',
	`alexa` INT(11) NOT NULL DEFAULT 0,
	`contury` VARCHAR(3000) NOT NULL DEFAULT '' COLLATE 'utf8mb4_bin',
	PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';
*/
func CreateDemo() {
	sql, err = table.BuildCreateStmt()

	if err != nil {
		fmt.Printf("Create err: %+v\n", err)
		return
	}
	fmt.Printf("sql: %s\n", sql)

	result, err = db.Exec(sql)

	if err != nil {
		fmt.Printf("Create Exec err: %+v\n", err)
		return
	}

	fmt.Printf("result: %s\n", result)
}

/*
INSERT INTO `Websites` (`id`, `name`, `url`, `alexa`, `contury`) VALUES
(NULL, 'Google', https://www.google.com/, 1, 'US'),
(NULL, 'Facebook', https://www.facebook.com/, 2, 'US'),
(NULL, 'apple', https://www.apple.com/, 3, 'US'),
(NULL, 'microsoft', https://www.microsoft.com/, 4, 'US');
*/
func InsertDemo() {
	insert := gTable.GetInserter()
	insert.Insert([]string{"NULL", "'Google'", "'https://www.google.com/'", "1", "'US'"})
	insert.Insert([]string{"NULL", "'Facebook'", "'https://www.facebook.com/'", "2", "'US'"})
	insert.Insert([]string{"NULL", "'apple'", "'https://www.apple.com/'", "3", "'US'"})
	insert.Insert([]string{"NULL", "'microsoft'", "'https://www.microsoft.com/'", "4", "'US'"})
	sql, err = insert.ToStmt()

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	fmt.Printf("insert1 sql: %s\n", sql)
	result, err = insert.Exec()

	if err != nil {
		fmt.Printf("Insert err: %+v\n", err)
		return
	}

	fmt.Printf("result: %+v\n", result)
	gTable.PutInserter(insert)

	// table.Insert([]string{"NULL", "'Google'", "'https://www.google.com/'", "1", "'US'"})
	// table.Insert([]string{"NULL", "'Facebook'", "'https://www.facebook.com/'", "2", "'US'"})
	// table.Insert([]string{"NULL", "'apple'", "'https://www.apple.com/'", "3", "'US'"})
	// table.Insert([]string{"NULL", "'microsoft'", "'https://www.microsoft.com/'", "4", "'US'"})
	// sql, err = table.BuildInsertStmt()
	// if err != nil {
	// 	fmt.Printf("Insert err: %+v\n", err)
	// 	return
	// }
	/* INSERT INTO `Websites` (`id`, `name`, `url`, `alexa`, `contury`) VALUES
	(NULL, 'Google', 'https://www.google.com/', 1, 'US'),
	(NULL, 'Facebook', 'https://www.facebook.com/', 2, 'US'),
	(NULL, 'apple', 'https://www.apple.com/', 3, 'US'),
	(NULL, 'microsoft', 'https://www.microsoft.com/', 4, 'US');*/
	// fmt.Printf("table sql: %s\n", sql)

	// result, err = db.Exec(sql)

	// if err != nil {
	// 	fmt.Printf("Insert Exec err: %+v\n", err)
	// 	return
	// }

	// fmt.Printf("result: %s\n", result)
}

// SELECT * FROM `Websites`;
func QueryDemo() {
	selector := gTable.GetSelector()
	sql, err = selector.ToStmt()
	if err != nil {
		fmt.Printf("Select err: %+v\n", err)
		return
	}
	fmt.Printf("selector sql: %s\n", sql)
	result, err = selector.Exec()
	if err != nil {
		fmt.Printf("Select err: %+v\n", err)
		return
	}
	fmt.Printf("result: %s\n", result)
	gTable.PutSelector(selector)

	// sql, err = table.BuildSelectStmt()
	// if err != nil {
	// 	fmt.Printf("Select err: %+v\n", err)
	// 	return
	// }
	// fmt.Printf("table sql: %s\n", sql)

	// result, err = db.Query(sql)

	// if err != nil {
	// 	fmt.Printf("Query Exec err: %+v\n", err)
	// 	return
	// }

	// fmt.Printf("result: %s\n", result)
	// for i, data := range result.Datas {
	// 	fmt.Printf("%d) data: %+v\n", i, data)
	// }
}

// UPDATE `Websites` SET `alexa` = 5000 WHERE `id` = 3;
func UpdateDemo() {
	updater := gTable.GetUpdater()
	updater.Update("alexa", "5000")
	updater.SetCondition(stmt.WS().Eq("name", "'Facebook'"))
	sql, err = updater.ToStmt()
	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}
	fmt.Printf("updater sql: %s\n", sql)

	result, err = updater.Exec()
	if err != nil {
		fmt.Printf("Update err: %+v\n", err)
		return
	}
	fmt.Printf("result: %+v\n", result)
	gTable.PutUpdater(updater)

	// table.Update("alexa", "5000")
	// table.SetUpdateCondition(stmt.WS().Eq("id", "3"))
	// sql, err = table.BuildUpdateStmt()

	// if err != nil {
	// 	fmt.Printf("Update err: %+v\n", err)
	// 	return
	// }
	// fmt.Printf("table sql: %s\n", sql)

	// result, err = db.Exec(sql)

	// if err != nil {
	// 	fmt.Printf("Update Exec err: %+v\n", err)
	// 	return
	// }

	// fmt.Printf("result: %s\n", result)
}

// DELETE FROM `Websites` WHERE `name` = Facebook;
func DeleteDemo() {
	deleter := gTable.GetDeleter()
	deleter.SetCondition(stmt.WS().Eq("name", "'Facebook'"))
	sql, err = deleter.ToStmt()
	if err != nil {
		fmt.Printf("Delete err: %+v\n", err)
		return
	}
	fmt.Printf("deleter sql: %s\n", sql)

	result, err = deleter.Exec()
	if err != nil {
		fmt.Printf("Delete err: %+v\n", err)
		return
	}
	fmt.Printf("result: %s\n", result)
	gTable.PutDeleter(deleter)

	// table.SetDeleteCondition(stmt.WS().Eq("name", "'Facebook'"))
	// sql, err = table.BuildDeleteStmt()

	// if err != nil {
	// 	fmt.Printf("Delete err: %+v\n", err)
	// 	return
	// }
	// fmt.Printf("table sql: %s\n", sql)

	// result, err = db.Exec(sql)

	// if err != nil {
	// 	fmt.Printf("Delete Exec err: %+v\n", err)
	// 	return
	// }

	// fmt.Printf("result: %s\n", result)
}
