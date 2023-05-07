package main

import (
	"fmt"

	"github.com/j32u4ukh/gosql/proto"
	"github.com/j32u4ukh/gosql/stmt/dialect"

	"github.com/pkg/errors"
)

func InitTable() (*proto.ProtoTable, error) {
	helper := &proto.Helper{
		Folder: "./pb",
		Dial:   dialect.MARIA,
	}

	tableName := "Desk"
	// GetParams(table_name string) (*stmt.TableParam, []*stmt.ColumnParam, error)
	tableParam, columnParam, err := helper.GetParams(tableName)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("TestInsertStmt | Error: %+v\n", err))
	}

	fmt.Printf("tableParam: %+v\n", tableParam)
	fmt.Println()
	fmt.Printf("columnParam: %+v\n", columnParam)
	fmt.Println()
	table := proto.NewProtoTable(tableName, tableParam, columnParam, helper.Dial)
	table.SetDbName("demo2")
	return table, nil
}

func main() {
	// CREATE TABLE IF NOT EXISTS `demo2`.`Desk` (`index` VARCHAR(23) NOT NULL AUTO_INCREMENT COMMENT '索引值' COLLATE 'utf8mb4_bin', `user_name` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '玩家名稱' COLLATE 'utf8mb4_bin', `item_id` SMALLINT(6) NOT NULL DEFAULT '' COMMENT '物品 ID', `time` TIMESTAMP NOT NULL DEFAULT current_timestamp() COMMENT 'Log 建立時間', PRIMARY KEY (`index`, `user_name`) USING BTREE) ENGINE = InnoDB COLLATE = 'utf8mb4_bin';
	table, err := InitTable()

	if err != nil {
		fmt.Printf("Failed to init table: %+v\n", err)
	}

	fmt.Printf("table: %+v\n", table)

	sql, err := table.BuildCreateStmt()

	if err != nil {
		fmt.Printf("Failed to generate create statement: %+v\n", err)
	}

	fmt.Printf("sql: %s\n", sql)
}
