package sync

import (
	"fmt"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/gdo"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
)

type SyncMode byte

const (
	ProtoToDb SyncMode = iota
	DbToDb
)

type Synchronize struct {
	// 同步模式
	Mode SyncMode
	// From: 同步後預期的狀態
	fromDB    *database.Database
	fromTable *gdo.Table
	// To: 要被改變的對象
	toDB        *database.Database
	toTable     *gdo.Table
	originTable *gdo.Table
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
}

// ====================================================================================================
// 根據不同模式，調用不同建構子
// ====================================================================================================
func NewProtoToDbSync(db *database.Database, tableNme string, dial dialect.SQLDialect) *Synchronize {
	s := &Synchronize{
		Mode:     ProtoToDb,
		toDB:     db,
		toTable:  gdo.NewTable(tableNme, stmt.NewTableParam(), nil, stmt.ENGINE, stmt.COLLATE, dial),
		alterMap: map[string][]string{},
		dropList: cntr.NewArray[string](),
		indexMap: map[string][]string{},
	}
	return s
}

func NewDbToDbSync(fromDB *database.Database, fromTableNme string, toDB *database.Database, toTableNme string, dial dialect.SQLDialect) *Synchronize {
	s := &Synchronize{
		Mode:      DbToDb,
		fromDB:    fromDB,
		fromTable: gdo.NewTable(fromTableNme, stmt.NewTableParam(), nil, stmt.ENGINE, stmt.COLLATE, dial),
		toDB:      toDB,
		toTable:   gdo.NewTable(toTableNme, stmt.NewTableParam(), nil, stmt.ENGINE, stmt.COLLATE, dial),
		alterMap:  map[string][]string{},
		dropList:  cntr.NewArray[string](),
		indexMap:  map[string][]string{},
	}
	return s
}

// ====================================================================================================

func (s *Synchronize) SetFromTable(t *gdo.Table) {
	s.fromTable = t
}

func (s *Synchronize) GetFromTable() *gdo.Table {
	return s.fromTable
}

func (s *Synchronize) SetFromDbName(dbName string) {
	s.fromTable.SetDbName(dbName)
}

func (s *Synchronize) SetToDbName(dbName string) {
	s.toTable.SetDbName(dbName)
}

func (s *Synchronize) Init(source string) error {
	s.queryInformationSchemaColumns(source)
	s.quertInformationSchemaStatistics(source)

	if source == "to" {
		s.originTable = s.toTable.SyncClone()
	}

	return nil
}

// 根據表格 INFORMATION_SCHEMA.`COLUMNS` 對 ProtoTable 的欄位做初始化
func (s *Synchronize) queryInformationSchemaColumns(source string) {
	var sgl *database.Database
	var t *gdo.Table

	if source == "from" {
		sgl = s.fromDB
		t = s.fromTable
	} else {
		sgl = s.toDB
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
	result, err := sgl.Query(sql)

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
	var sgl *database.Database
	var t *gdo.Table

	if source == "from" {
		sgl = s.fromDB
		t = s.fromTable
	} else {
		sgl = s.toDB
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
	result, err := sgl.Query(sql)

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
	t.SetTableParam(tableParam)
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
	orders := s.fromTable.GetColumnNames().Elements

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

		// fromCol, _ = s.fromTable.GetPrimaryColumn()
		pks := []string{}

		for _, element := range fromTableParam.Primarys.Elements {
			pks = append(pks, fmt.Sprintf("`%s`", element))
		}

		// 標註新增 PRIMARY KEY
		s.indexMap["PRIMARY KEY"] = []string{
			"PRIMARY KEY",
			strings.Join(pks, ", "),
			s.fromTable.TableParam.IndexType["PRIMARY"],
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

	var data []any
	var kind, indexName, fromIndexType string
	var fromCols, toCols *cntr.Array[string]
	// var modifiedCols *array.Array[string]
	var colNames []string
	var needModify, needDrop bool

	// 已修正(不重複)欄位名稱
	// modifiedCols = array.NewArray[string]()

	// 遍歷 Proto 定義的 index
	for it.HasNext() {
		// 0: kind string, 1: indexName string, 2: indexType string, 3: cols *array.Array[string]
		data = it.Next().([]any)
		kind = data[0].(string)
		indexName = data[1].(string)
		fromIndexType = data[2].(string)
		fromCols = data[3].(*cntr.Array[string])

		toCols = s.toTable.GetTableParam().GetIndexColumns(kind, indexName)
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

			if needDrop && !s.dropList.Contains(indexName) {
				s.dropList.Append(indexName)
			}

			s.indexMap[indexName] = []string{kind, strings.Join(colNames, ", "), fromIndexType}
		}
	}

	it = s.toTable.GetTableParam().IterIndexMap()

	// 遍歷 toTable 當前定義的 index
	for it.HasNext() {
		// 0: kind string, 1: indexName string, 2: indexType string, 3: cols *array.Array[string]
		data = it.Next().([]any)
		kind = data[0].(string)
		indexName = data[1].(string)
		fromCols = s.fromTable.GetTableParam().GetIndexColumns(kind, indexName)

		// toTable 有定義，但 fromTable 沒有定義
		if fromCols == nil && !s.dropList.Contains(indexName) {
			s.dropList.Append(indexName)
		}
	}
}

func (s *Synchronize) PrintCheckResult() bool {
	var origin, curret *stmt.Column
	hasDiff := false

	gosql.Info("===== CHANGE / ADD =====")
	for k, v := range s.alterMap {
		hasDiff = true
		curret = s.toTable.GetColumnByName(k)

		if v[0] == "change" {
			origin = s.originTable.GetColumnByName(v[1])
			gosql.Info("Origin | `%s` %s", v[1], origin.GetInfo())
			gosql.Info("CHANGE | `%s` %s", k, curret.GetInfo())

		} else if v[0] == "add" {
			gosql.Info("ADD | `%s` %s", k, curret.GetInfo())
		}
	}

	// DROP INDEX `索引 2`
	gosql.Info("===== DROP INDEX =====")
	for _, index := range s.dropList.Elements {
		hasDiff = true
		gosql.Info("DROP `%s`", index)
	}

	// ADD INDEX `索引 2`
	gosql.Info("===== ADD INDEX =====")
	for index := range s.indexMap {
		hasDiff = true
		gosql.Info("ADD `%s`", index)
	}

	return hasDiff
}

func (s *Synchronize) SyncTableSchema(needExcute bool) error {
	fmt.Printf("(s *Synchronize) SyncTableSchema | ===== SyncTableSchema =====\n")
	alter := s.getAlterStmts()

	drop := s.getDropIndex()
	alter = append(alter, drop...)

	index := s.getAddIndex()
	alter = append(alter, index...)

	if len(alter) == 0 {
		fmt.Printf("(s *Synchronize) SyncTableSchema | Database schema is same.\n")
	} else {
		sql := fmt.Sprintf("ALTER TABLE `%s`.`%s` %s;", s.toTable.GetDbName(), s.toTable.GetTableName(), strings.Join(alter, ", "))
		fmt.Printf("(s *Synchronize) SyncTableSchema | sql: %s\n", sql)

		if needExcute {
			result, err := s.toDB.Exec(sql)

			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("sql: %s", sql))
			}

			fmt.Printf("(s *Synchronize) SyncTableSchema | result: %s\n", result)
		}
	}

	return nil
}

// 生成有序的 CHANGE 和 ADD 語法列表
func (s *Synchronize) getAlterStmts() []string {
	stmts := []string{}
	colNames := s.toTable.GetColumnNames().Elements
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
				gosql.Info("Wrong idx: %d, colName: %s, cols: %+v", idx, colName, colNames)
				continue
			}

			// ===== 生成 SQL 語法 =====
			if params[0] == "change" {
				sql = fmt.Sprintf("CHANGE COLUMN `%s` `%s` %s", params[1], colName, info)

			} else if params[0] == "add" {
				sql = fmt.Sprintf("ADD COLUMN `%s` %s", colName, info)

			} else {
				gosql.Info("Wrong operate: %+v", params)
				continue
			}

			// 將 SQL 語法加入列表當中
			stmts = append(stmts, sql)
		}
	}
	return stmts
}
