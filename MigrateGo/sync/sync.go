package sync

import (
	"fmt"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/plugin"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
)

// 同步模式
type SyncMode byte

const (
	// Proto 檔 -> Db 中的 Table
	NoneMode SyncMode = 0
	// Proto 檔 -> Db 中的 Table
	ProtoToDbMode SyncMode = 1
	// Db 中的 Table -> 另一個 Db 中 Table
	DbToDbMode SyncMode = 2
)

type Synchronize struct {
	// 同步模式
	Mode SyncMode
	// From: 同步後預期的狀態
	fromDB    *database.Database
	fromTable *gosql.Table
	// To: 要被改變的對象
	toDB        *database.Database
	toTable     *gosql.Table
	originTable *gosql.Table
	// 需 Change/Add 的欄位名稱
	// key: 更新後欄位名稱
	// add: -
	// change: 原始欄位名稱
	alterMap map[string][]string
	// 需 DROP 的索引名稱
	dropList *cntr.Array[string]
	// 需 ADD 的索引名稱
	indexMap     map[string][]string
	changedOrder *cntr.Array[string]
	// 若表格結構有差異，是否印出即將執行的指令?
	Print bool
	// 若表格結構有差異，是否執行同步
	Sync bool
}

func NewSynchronize() *Synchronize {
	s := &Synchronize{
		Mode:         NoneMode,
		fromDB:       nil,
		fromTable:    nil,
		toDB:         nil,
		toTable:      nil,
		originTable:  nil,
		alterMap:     make(map[string][]string),
		dropList:     &cntr.Array[string]{},
		indexMap:     make(map[string][]string),
		changedOrder: &cntr.Array[string]{},
	}
	return s
}

func (s *Synchronize) Execute(config *Config) error {
	err := s.LoadConfig(config)
	if err != nil {
		return errors.Wrap(err, "Failed to load config.")
	}

	// 檢查表格是否存在
	var isTableExists bool

	switch s.Mode {
	case DbToDbMode:
		isTableExists, err = s.fromDB.IsTableExists(stmt.IsTableExists(s.fromDB.DbName, config.FromTable))

		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to check table %s.", config.ToTable))
		}

		// 若表格不存在
		if !isTableExists {
			return errors.Errorf("來源表格 %s 不存在", config.ToTable)
		}

		// 讀取資料庫結構數據
		s.InitTable("from")
	}

	isTableExists, err = s.toDB.IsTableExists(stmt.IsTableExists(s.toDB.DbName, config.ToTable))

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to check table %s.", config.ToTable))
	}

	// 若表格不存在
	if !isTableExists {
		// 若需要生成表格
		if config.Generate {
			_, err = s.toTable.Creater().Exec()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to create table %s.", config.ToTable))
			}
		}
		return nil
	}

	// 讀取資料庫結構數據
	s.InitTable("to")
	err = s.CheckTableStructure()

	if err != nil {
		return errors.Wrap(err, "Failed to check table structure.")
	}

	// 若有表格差異
	if s.PrintCheckResult() {
		if s.Sync {
			err = s.SyncTableSchema(true)
			if err != nil {
				return errors.Wrap(err, "資料表結構同步時發生錯誤")
			}
		} else if s.Print {
			s.SyncTableSchema(false)
		}
	}
	return nil
}

func (s *Synchronize) LoadConfig(config *Config) error {
	var dc *database.DatabaseConfig
	var err error
	s.Mode = config.Mode
	s.Print = config.Print
	s.Sync = config.Sync
	// ==================================================
	// Connnect
	// ==================================================
	switch s.Mode {
	case DbToDbMode:
		dc = config.FromDatabase
		s.fromDB, err = database.Connect(0, dc.User, dc.Password, dc.Host, dc.Port, dc.DbName)
		if err != nil {
			return errors.Wrapf(err, "與資料庫(%s:%d)連線時發生錯誤, err: %+v", dc.Host, dc.Port, err)
		}
	}
	dc = config.ToDatabase
	s.toDB, err = database.Connect(1, dc.User, dc.Password, dc.Host, dc.Port, dc.DbName)
	if err != nil {
		return errors.Wrapf(err, "與資料庫(%s:%d)連線時發生錯誤, err: %+v", dc.Host, dc.Port, err)
	}
	s.toTable = gosql.NewTable(config.FromTable, stmt.NewTableParam(), nil, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	s.toTable.Init(&gosql.TableConfig{
		Db:               s.toDB,
		DbName:           dc.DbName,
		UseAntiInjection: false,
		PtrToDbFunc:      plugin.ProtoToDb,
		InsertFunc:       plugin.InsertProto,
		QueryFunc:        plugin.QueryProto,
		UpdateAnyFunc:    plugin.UpdateProto,
	})
	switch s.Mode {
	case ProtoToDbMode:
		protoPath := fmt.Sprintf("%s/%s.proto", config.ProtoFolder, config.FromTable)
		tableParams, columnParams, err := plugin.GetProtoParams(protoPath, dialect.MARIA)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to get proto parameters from %s.", protoPath))
		}
		s.fromTable = gosql.NewTable(config.FromTable, tableParams, columnParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
		s.fromTable.Init(&gosql.TableConfig{
			Db:               s.toDB,
			DbName:           dc.DbName,
			UseAntiInjection: false,
			PtrToDbFunc:      plugin.ProtoToDb,
			InsertFunc:       plugin.InsertProto,
			QueryFunc:        plugin.QueryProto,
			UpdateAnyFunc:    plugin.UpdateProto,
		})
	}
	return nil
}

// ====================================================================================================

func (s *Synchronize) InitTable(source string) error {
	s.queryInformationSchemaColumns(source)
	s.quertInformationSchemaStatistics(source)
	if source == "to" {
		s.originTable = s.toTable.SyncClone()
	}
	return nil
}

// 根據表格 INFORMATION_SCHEMA.`COLUMNS` 對 ProtoTable 的欄位做初始化
func (s *Synchronize) queryInformationSchemaColumns(source string) {
	var db *database.Database
	var t *gosql.Table

	if source == "from" {
		db = s.fromDB
		t = s.fromTable
	} else {
		db = s.toDB
		t = s.toTable
	}

	columnNames := []string{
		"`COLUMN_NAME`",
		"`COLUMN_TYPE`",
		"`COLUMN_KEY`",
		"`IS_NULLABLE`",
		"`COLUMN_DEFAULT`",
		"`COLLATION_NAME`",
		"`EXTRA`",
		"`COLUMN_COMMENT`",
	}

	queryColumns := strings.Join(columnNames, ", ")
	where := fmt.Sprintf("WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'", t.GetDbName(), t.GetTableName())
	sql := fmt.Sprintf("SELECT %s FROM INFORMATION_SCHEMA.`COLUMNS` %s;", queryColumns, where)
	result, err := db.Query(sql)

	if err != nil {
		fmt.Printf("Query sql: %s, err: %+v\n", sql, err)
		return
	}

	var i int32
	var dbColumn *stmt.Column
	for i = 0; i < result.NRow; i++ {
		dbColumn = stmt.NewDbColumn(
			result.Datas[i][0],
			result.Datas[i][1],
			result.Datas[i][2] == "PRI",
			result.Datas[i][3] == "YES",
			result.Datas[i][4],
			result.Datas[i][5],
			result.Datas[i][6],
			result.Datas[i][7],
		)
		t.AddColumn(dbColumn)
	}
}

// 根據表格 INFORMATION_SCHEMA.`STATISTICS` 對表格索引(如 Primary key)做初始化
func (s *Synchronize) quertInformationSchemaStatistics(source string) {
	var db *database.Database
	var t *gosql.Table

	if source == "from" {
		db = s.fromDB
		t = s.fromTable
	} else {
		db = s.toDB
		t = s.toTable
	}

	/*SELECT * FROM INFORMATION_SCHEMA.`STATISTICS`
	  WHERE table_name = 'tbl_name'
	  AND table_schema = 'db_name'*/
	columnNames := []string{
		"`NON_UNIQUE`",
		"`INDEX_NAME`",
		"`INDEX_TYPE`",
		"`COLUMN_NAME`",
	}

	queryColumns := strings.Join(columnNames, ", ")
	where := fmt.Sprintf("WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'", t.GetDbName(), t.GetTableName())
	sql := fmt.Sprintf("SELECT %s FROM INFORMATION_SCHEMA.`STATISTICS` %s;", queryColumns, where)
	result, err := db.Query(sql)

	if err != nil {
		fmt.Printf("Query err: %+v", err)
		return
	}

	tableParam := stmt.NewTableParam()
	var kind string
	var i int32
	for i = 0; i < result.NRow; i++ {
		if result.Datas[i][0] == "0" {
			if result.Datas[i][1] == "PRIMARY" {
				tableParam.IndexType["PRIMARY"] = result.Datas[i][2]
				tableParam.Primarys.Append(result.Datas[i][3])
				continue
			}

			kind = "UNIQUE"
		} else {
			kind = "INDEX"
		}

		tableParam.AddIndex(kind, result.Datas[i][1], result.Datas[i][2], result.Datas[i][3])
	}
	t.Creater().SetTableParam(tableParam)
}

// ====================================================================================================
// 檢查欄位結構 & 同步結構
// ====================================================================================================

// 檢查資料表結構
func (s *Synchronize) CheckTableStructure() error {
	var fromCol, toCol *stmt.Column
	nColumn := s.fromTable.GetColumnNumber()
	var i, idx int32

	// 遍歷 fromTable 當中的欄位，確保 toTable 中有對應欄位，且該欄位的配置與 fromTable 相同
	for i = 0; i < nColumn; i++ {
		fromCol = s.fromTable.GetColumn(i)

		// 根據欄位名稱，尋找 toTable 中對應的欄位
		toCol = s.toTable.GetColumnByName(fromCol.Name)

		if toCol == nil {
			// 若 toTable 沒有，則需要新增欄位
			s.addNewColumn(fromCol)

		} else {
			// 若 toTable 有相同名稱的欄位，但結構設定不同
			if !fromCol.IsEquals(toCol) {
				// 更新 toTable 的當前欄位配置
				idx = s.toTable.GetIndexByName(fromCol.Name)
				s.changeColumn(idx, fromCol)
			}
		}
	}

	// 取得初步同步後的 toTable 欄位數
	nColumn = s.toTable.GetColumnNumber()

	// 遍歷 toTable 當中的欄位
	for i = 0; i < nColumn; i++ {
		toCol = s.toTable.GetColumn(i)
		fromCol = s.fromTable.GetColumnByName(toCol.Name)

		// 若 fromCol 沒有，表示當前欄位是被廢棄的，先修改名稱，之後再手動調整
		if fromCol == nil {
			if toCol.IsPrimaryKey {
				// TODO: toCol.IsPrimaryKey = false
				if toCol.Default == "AI" {
					// toCol.CanNull = true
					toCol.Default = "NIL"
				}
			}

			// 新增更名後的欄位(名稱後綴添加 _backup，表明此欄位需手動調整)
			s.renameColumn(i, toCol)
		}
	}

	//////////////////////////////////////////////////
	// 篩選出需要調整順序的欄位
	//////////////////////////////////////////////////
	// 取得 fromTable 的欄位名稱(toTable 的排序依據)
	orders := s.fromTable.GetColumnNames(true)

	// 根據指定順序(orders)調整表格欄位，並返回被調整的欄位的名稱
	s.changedOrder = s.toTable.RefreshColumnOrder(orders)

	// 將需要調整順序的欄位名稱加入 caMap(change)
	for _, order := range s.changedOrder.Elements {
		if _, ok := s.alterMap[order]; !ok {
			// 標註這些欄位需要被 CHANGE
			s.alterMap[order] = []string{"change", order}
		}
	}

	fromTableParam := s.fromTable.GetTableParam()
	toTableParam := s.toTable.GetTableParam()

	// 檢查兩者的 PRIMARY KEY 配置是否相同
	if !fromTableParam.Primarys.IsEquals(s.originTable.GetTableParam().Primarys, false) {
		if s.originTable.GetTableParam().Primarys.Length() > 0 {
			// 標註需移除 PRIMARY KEY
			s.dropList.Append("PRIMARY KEY")
		}

		// 以 fromTable 的表格參數覆蓋 toTable 的設置
		toTableParam.Primarys = nil
		toTableParam.Primarys = fromTableParam.Primarys.Clone()
		pks := []string{}

		for _, element := range fromTableParam.Primarys.Elements {
			pks = append(pks, fmt.Sprintf("`%s`", element))
		}

		// 標註新增 PRIMARY KEY
		s.indexMap["PRIMARY KEY"] = []string{
			"PRIMARY KEY",
			strings.Join(pks, ", "),
			fromTableParam.IndexType["PRIMARY"],
		}
	}

	// 檢查其他索引是否
	s.checkIndex()

	return nil
}

// fromTable 有的欄位，但 toTable 沒有，添加 toTable 缺少的欄位
func (s *Synchronize) addNewColumn(column *stmt.Column) {
	// 更新 toTable，添加 column 欄位
	s.toTable.AddColumn(column)

	// 標註欄位 column 需要增加
	// s.addMap[column.Name] = nil
	s.alterMap[column.Name] = []string{"add"}
}

// 欄位名稱相同，但結構設定不同，將 toTable 欄位配置更新成與 fromTable 相同
func (s *Synchronize) changeColumn(index int32, column *stmt.Column) {
	// 更新 toTable 的欄位
	s.toTable.SetColumn(index, column.Clone())

	// 標註欄位 column 需要修改
	s.alterMap[column.Name] = []string{"change", column.Name}
}

func (s *Synchronize) renameColumn(index int32, column *stmt.Column) {
	// 若名稱後綴已有 _backup，則直接返回
	if strings.HasSuffix(column.Name, "_backup") {
		return
	}

	newColumn := column.Clone()
	newColumn.Name = fmt.Sprintf("%s_backup", column.Name)

	// 若先前有對舊欄位名稱的操作，這裡將其移除
	if change, ok := s.alterMap[column.Name]; ok {
		if change[0] == "change" {
			delete(s.alterMap, column.Name)
		}
	}

	// s.changeMap[newColumn.Name] = column.Name
	s.alterMap[newColumn.Name] = []string{"change", column.Name}

	s.toTable.SetColumn(index, newColumn)
	column = nil
}

func (s *Synchronize) checkIndex() {
	// 取得 fromTable 索引(Primary key 除外)的迭代器
	it := s.fromTable.GetTableParam().IterIndexMap()

	var data *stmt.SqlIndex
	var kind string
	var fromCols, toCols *cntr.Array[string]
	var colNames []string
	var needModify, needDrop bool

	// 遍歷 Proto 定義的 index
	for it.HasNext() {
		data = it.Next()
		kind = data.Kind
		fromCols = data.Cols

		toCols = s.toTable.GetTableParam().GetIndexColumns(kind, data.Name)
		needModify = false
		needDrop = false

		// toTable 中沒有對應的索引
		if toCols == nil {
			needModify = true
		} else {
			// 若兩者所包含欄位名稱不相同，則需要對 index 做調整
			needModify = !fromCols.IsEquals(toCols, false)
			needDrop = true
		}

		if needModify {
			if kind == "UNIQUE" {
				kind = "UNIQUE INDEX"
			} else {
				kind = "INDEX"
			}

			colNames = []string{}

			for _, element := range fromCols.Elements {
				colNames = append(colNames, fmt.Sprintf("`%s`", element))
			}

			if needDrop && !s.dropList.Contains(data.Name) {
				s.dropList.Append(data.Name)
			}

			s.indexMap[data.Name] = []string{kind, strings.Join(colNames, ", "), data.Algo}
		}
	}

	it = s.toTable.GetTableParam().IterIndexMap()

	// 遍歷 toTable 當前定義的 index
	for it.HasNext() {
		data = it.Next()
		fromCols = s.fromTable.GetTableParam().GetIndexColumns(data.Kind, data.Name)

		// toTable 有定義，但 fromTable 沒有定義
		if fromCols == nil && !s.dropList.Contains(data.Name) {
			s.dropList.Append(data.Name)
		}
	}
}

func (s *Synchronize) PrintCheckResult() bool {
	var origin, curret *stmt.Column
	hasDiff := false

	s.print("===== CHANGE / ADD =====")
	for k, v := range s.alterMap {
		hasDiff = true
		curret = s.toTable.GetColumnByName(k)

		if v[0] == "change" {
			origin = s.originTable.GetColumnByName(v[1])
			s.print("Origin | `%s` %s", v[1], origin.GetInfo())
			s.print("CHANGE | `%s` %s", k, curret.GetInfo())

		} else if v[0] == "add" {
			s.print("ADD | `%s` %s", k, curret.GetInfo())
		}
	}

	// DROP INDEX `索引 2`
	s.print("===== DROP INDEX =====")
	for _, index := range s.dropList.Elements {
		hasDiff = true
		s.print("DROP `%s`", index)
	}

	// ADD INDEX `索引 2`
	s.print("===== ADD INDEX =====")
	for index := range s.indexMap {
		hasDiff = true
		s.print("ADD `%s`", index)
	}

	return hasDiff
}

func (s *Synchronize) print(format string, a ...any) {
	if s.Print {
		fmt.Printf(format+"\n", a...)
	}
}

func (s *Synchronize) SyncTableSchema(needExcute bool) error {
	s.print("(s *Synchronize) SyncTableSchema | ===== SyncTableSchema =====")
	alter := s.getAlterStmts()

	drop := s.getDropIndex()
	alter = append(alter, drop...)

	index := s.getAddIndex()
	alter = append(alter, index...)

	if len(alter) == 0 {
		s.print("Database schema is same.")
	} else {
		sql := fmt.Sprintf("ALTER TABLE `%s`.`%s` %s;", s.toTable.GetDbName(), s.toTable.GetTableName(), strings.Join(alter, ", "))
		s.print("sql: %s", sql)

		if needExcute {
			result, err := s.toDB.Exec(sql)

			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("sql: %s", sql))
			}

			s.print("(s *Synchronize) SyncTableSchema | result: %s", result)
		}
	}
	return nil
}

// 生成有序的 CHANGE 和 ADD 語法列表
func (s *Synchronize) getAlterStmts() []string {
	stmts := []string{}
	colNames := s.toTable.GetColumnNames(true)
	var column *stmt.Column
	var info, sql string

	// CHANGE 和 ADD 的語法順序需根據欄位順序生成
	for idx, colName := range colNames {
		// 檢查是否在需要修改的欄位當中
		if params, ok := s.alterMap[colName]; ok {
			// 取得當前欄位
			column = s.toTable.GetColumn(int32(idx))

			// ===== 根據欄位順序，生成欄位描述 =====
			if idx == 0 {
				info = fmt.Sprintf("%s FIRST", column.GetInfo())

			} else if idx > 0 {
				// 取得前一欄位名稱
				lastColumn := colNames[idx-1]
				info = fmt.Sprintf("%s AFTER `%s`", column.GetInfo(), lastColumn)

			} else {
				s.print("Wrong idx: %d, colName: %s, cols: %+v", idx, colName, colNames)
				continue
			}

			// ===== 生成 SQL 語法 =====
			if params[0] == "change" {
				sql = fmt.Sprintf("CHANGE COLUMN `%s` `%s` %s", params[1], colName, info)

			} else if params[0] == "add" {
				sql = fmt.Sprintf("ADD COLUMN `%s` %s", colName, info)

			} else {
				s.print("Wrong operate: %+v", params)
				continue
			}

			// 將 SQL 語法加入列表當中
			stmts = append(stmts, sql)
		}
	}
	return stmts
}
