package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/j32u4ukh/glog/v2"

	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/proto/gstmt"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/j32u4ukh/gosql/sync"
	"github.com/j32u4ukh/gosql/utils"

	"github.com/pkg/errors"
)

var synConfig *sync.Config
var logger *glog.Logger

func main() {
	logger = glog.SetLogger(0, "demo_sync", glog.DebugLevel)
	logger.SetFolder("../log")
	utils.SetLogger(logger)

	var err error
	var path string

	if len(os.Args) == 1 {
		path = "./config.yaml"
	} else {
		path = os.Args[1]
	}

	synConfig, err = sync.NewConfig(path)

	if err != nil {
		logger.Error("Config 初始化失敗, Error: %+v", err)
		return
	}

	// fromDB: 0, toDB: 1
	fromDB, toDB, err := connect()
	defer toDB.Close()

	if fromDB != nil {
		defer fromDB.Close()
	}

	if err != nil {
		logger.Error("資料庫連線失敗, Error: %+v", err)
		return
	}

	if synConfig.Mode == 1 || synConfig.Mode == 3 {
		ProtoToDb(toDB)

	} else if synConfig.Mode == 2 || synConfig.Mode == 4 {
		DbToDb(fromDB, toDB)
	}
}

func connect() (*database.Database, *database.Database, error) {
	var fromDB, toDB *database.Database = nil, nil
	var dc *database.DatabaseConfig
	var err error

	if synConfig.Mode == 2 || synConfig.Mode == 4 || synConfig.Mode == 6 {
		dc = synConfig.GetFromDatabase()
		fromDB, err = database.Connect(0, dc.UserName, dc.Password, dc.Server, dc.Port, dc.Name)

		if err != nil {
			return nil, nil, errors.Wrapf(err, "與資料庫連線時發生錯誤, err: %+v\n", err)
		}
	}

	dc = synConfig.GetToDatabase()
	toDB, err = database.Connect(1, dc.UserName, dc.Password, dc.Server, dc.Port, dc.Name)

	if err != nil {
		return nil, nil, errors.Wrapf(err, "與資料庫連線時發生錯誤, err: %+v\n", err)
	}

	return fromDB, toDB, nil
}

func ProtoToDb(toDB *database.Database) {
	if synConfig.Mode == 3 {
		// 獲取檔案或目錄相關資訊
		files, err := ioutil.ReadDir(synConfig.ProtoFolder)

		if err != nil {
			logger.Error("Failed to get table files from %s.\nError: %+v.", synConfig.ProtoFolder, err)
			return
		}

		tables := []string{}
		var table string

		for i := range files {
			// 印出指定目錄下的檔案名
			table = strings.Split(files[i].Name(), ".")[0]
			tables = append(tables, table)
		}

		synConfig.SetFromTables(tables)
		synConfig.SetToTables(tables)
	}

	length := len(synConfig.FromTables)

	for i := 0; i < length; i++ {
		protoToDb(toDB, synConfig.FromTables[i], synConfig.ToTables[i])
	}
}

func protoToDb(toDB *database.Database, fromTableName string, toTableName string) {
	// 檢查表格是否存在
	isTableExists, err := toDB.IsTableExists(stmt.IsTableExists(toDB.DbName, toTableName))

	gs, err := gstmt.SetGstmt(1, toDB.DbName, dialect.MARIA)

	if err != nil {
		fmt.Printf("protoToDb | 建立 Gstmt 時發生錯誤, err: %+v\n", err)
		return
	}

	// 根據檔案名稱，取得表格與欄位參數
	tableParam, colParams, err := gs.Helper.GetParamsByPath("../pb", toTableName)

	if err != nil {
		fmt.Printf("protoToDb | 從 %s.proro 讀取參數時發生錯誤, err: %+v\n", toTableName, err)
		return
	}

	t := gdo.NewTable(fromTableName, tableParam, colParams, stmt.ENGINE, stmt.COLLATE, gs.Helper.Dial)
	t.SetDbName(toDB.DbName)

	// 若表格不存在
	if !isTableExists {
		if synConfig.AutoGenerate {
			generateTable(func() error {
				sql, err := t.BuildCreateStmt()

				if err != nil {
					return errors.Wrapf(err, "Failed to generate statement to create table %s.", toTableName)
				}

				fmt.Printf("protoToDb | sql: %s\n", sql)

				// 生成當前表格
				_, err = toDB.Exec(sql)

				if err != nil {
					return errors.Wrapf(err, "Exec failed sql: %s.", sql)
				}

				fmt.Printf("protoToDb | 表格 %s 完成生成，無須進行比對", toTableName)
				return nil
			})
		}
		return
	}

	s := sync.NewProtoToDbSync(toDB, toTableName, dialect.MARIA)
	s.SetFromTable(t)
	s.SetFromDbName(synConfig.GetToDatabase().Name)
	s.SetToDbName(synConfig.GetToDatabase().Name)

	// 讀取資料庫結構數據
	s.Init("to")

	// 檢查、詢問與同步
	checkAskAndSynchronize(s)
}

func DbToDb(fromSgl *database.Database, toSgl *database.Database) {
	if synConfig.Mode == 4 {
		tables := getAllTableNames(toSgl)
		synConfig.SetFromTables(tables)
		synConfig.SetToTables(tables)
	}

	length := len(synConfig.FromTables)

	for i := 0; i < length; i++ {
		dbToDb(fromSgl, synConfig.FromTables[i], toSgl, synConfig.ToTables[i])
	}
}

func dbToDb(fromSgl *database.Database, fromTableNme string, toSgl *database.Database, toTableName string) {
	// 檢查表格是否存在
	isTableExists, err := toSgl.IsTableExists(stmt.IsTableExists(toSgl.DbName, toTableName))

	if err != nil {
		fmt.Printf("檢查表格 %s 是否存在時發生錯誤, err: %+v\n", toTableName, err)
		return
	}

	s := sync.NewDbToDbSync(fromSgl, fromTableNme, toSgl, toTableName, dialect.MARIA)

	fromConfig := synConfig.GetFromDatabase()
	s.SetFromDbName(fromConfig.Name)

	// 讀取資料庫結構數據
	s.Init("from")

	// 若表格不存在
	if !isTableExists {
		if synConfig.AutoGenerate {
			generateTable(func() error {
				stmt, err := s.GetFromTable().BuildCreateStmt()

				if err != nil {
					return errors.Wrapf(err, "Failed to generate statement to create table %s.", toTableName)
				}

				logger.Debug("Stmt: %s", stmt)

				// 生成當前表格
				_, err = toSgl.Exec(stmt)

				if err != nil {
					return errors.Wrapf(err, "Exec failed Stmt: %s.", stmt)
				}

				logger.Debug("表格 %s 完成生成", toTableName)
				return nil
			})
		}
		return
	}

	s.SetToDbName(synConfig.GetToDatabase().Name)

	// 讀取資料庫結構數據
	s.Init("to")

	// 檢查、詢問與同步
	checkAskAndSynchronize(s)
}

func getAllTableNames(toSgl *database.Database) []string {
	tables := []string{}
	sql := fmt.Sprintf("SELECT `TABLE_NAME` FROM INFORMATION_SCHEMA.`TABLES` WHERE TABLE_SCHEMA = '%s';",
		synConfig.GetToDatabase().Name)
	result, err := toSgl.Query(sql)

	if err != nil {
		return tables
	}

	for _, tableName := range result.Datas {
		tables = append(tables, tableName[0])
	}

	return tables
}

// 詢問是否生成表格，若要，則執行 genFunc 來生成
func generateTable(genFunc func() error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("是否要生成表格?(y/n) ")
	text, err := reader.ReadString('\n')

	if err != nil {
		logger.Error("讀取指令時發生錯誤, Error: %+v", err)
		return
	}

	text = strings.TrimRightFunc(text, func(c rune) bool {
		return c == '\r' || c == '\n'
	})
	text = strings.ToUpper(text)

	if text == "Y" || text == "YES" {
		err = genFunc()

		if err != nil {
			logger.Error("生成資料表時發生錯誤, Error: %+v", err)
		} else {
			logger.Info("完成資料表生成")
		}
	} else {
		logger.Info("不生成資料表")
	}
}

// 檢查、詢問與同步
func checkAskAndSynchronize(s *sync.Synchronize) {
	err := s.CheckTableStructure()

	if err != nil {
		logger.Error("Failed to check table structure.\nError: %+v", err)
		return
	}

	hasDifferent := s.PrintCheckResult()

	if hasDifferent {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("是否要同步結構?(y/n) ")
		text, err := reader.ReadString('\n')

		if err != nil {
			logger.Error("讀取'是否要同步結構?(y/n)'指令時發生錯誤, Error: %+v", err)
			return
		}

		text = strings.TrimRightFunc(text, func(c rune) bool {
			//In windows newline is \r\n
			return c == '\r' || c == '\n'
		})
		text = strings.ToUpper(text)

		if text == "Y" {
			err = s.SyncTableSchema(true)

			if err != nil {
				logger.Error("資料表結構同步時發生錯誤, Error: %+v", err)
				return
			}

			logger.Info("完成資料表結構同步")
		} else {
			fmt.Print("gosql", "是否印出 SQL 指令?(y/n) ")
			text, err = reader.ReadString('\n')

			if err != nil {
				logger.Error("讀取'是否印出 SQL 指令?(y/n)'指令時發生錯誤, Error: %+v", err)
				return
			}

			text = strings.TrimRightFunc(text, func(c rune) bool {
				//In windows newline is \r\n
				return c == '\r' || c == '\n'
			})
			text = strings.ToUpper(text)

			if text == "Y" {
				s.SyncTableSchema(false)
			} else {
				logger.Info("取消資料表結構同步")
			}
		}
	} else {
		logger.Info("資料表結構相同")
	}
}
