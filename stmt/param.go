package stmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/j32u4ukh/gosql/utils"
	"github.com/pkg/errors"
)

type SqlIndex struct {
	Kind string
	Name string
	Algo string
	Cols *cntr.Array[string]
}

func NewSqlIndex(kind string, name string, algo string, cols *cntr.Array[string]) *SqlIndex {
	si := &SqlIndex{
		Kind: kind,
		Name: name,
		Algo: algo,
		Cols: cols.Clone(),
	}
	return si
}

type TableParam struct {
	// Primary key
	Primarys *cntr.Array[string]
	/* 結構為
	{
		"UNIQUE":
		{
			"UKey": ["index", "item"],
		},
		"INDEX":
		{
			"Key": ["user", "item2"],
		}
	}
	*/
	indexMap map[string]map[string]*cntr.Array[string]
	// 排序演算法設置
	IndexType map[string]string
}

func NewTableParam() *TableParam {
	p := &TableParam{
		Primarys: cntr.NewArray[string](),
		indexMap: map[string]map[string]*cntr.Array[string]{
			"UNIQUE": {},
			"INDEX":  {},
		},
		IndexType: map[string]string{},
	}

	// Primary key 排序演算法設置
	p.IndexType["PRIMARY"] = ALGO
	return p
}

func (p *TableParam) LoadConfig(config *TableParamConfig) {
	var kind string
	var indexs *IndexConfig

	for _, kind = range config.Primarys.Elements {
		if !p.Primarys.Contains(kind) {
			p.Primarys.Append(kind)
		}
	}

	for _, indexs = range config.Uniques {
		kind = strings.ToUpper(indexs.Type)

		if kind == "" || kind == "DEFAULT" {
			indexs.Type = ALGO
		}

		// 設置表格的 Index
		p.AddIndex("UNIQUE", indexs.Name, indexs.Type, indexs.Columns...)
	}

	for _, indexs = range config.Indexs {
		kind = strings.ToUpper(indexs.Type)

		if kind == "" || kind == "DEFAULT" {
			indexs.Type = ALGO
		}

		// 設置表格的 Index
		p.AddIndex("INDEX", indexs.Name, indexs.Type, indexs.Columns...)
	}
}

func (p *TableParam) ParserConfig(config string) {
	tpc, err := NewTableParamConfig(config)

	if err != nil {
		return
	}

	for _, indexs := range tpc.Uniques {
		if indexs.Type == "" {
			indexs.Type = ALGO
		}

		// 設置表格的 Index
		p.AddIndex("UNIQUE", indexs.Name, indexs.Type, indexs.Columns...)
	}

	for _, indexs := range tpc.Indexs {
		if indexs.Type == "" {
			indexs.Type = ALGO
		}

		// 設置表格的 Index
		p.AddIndex("INDEX", indexs.Name, indexs.Type, indexs.Columns...)
	}
}

func (p *TableParam) AddPrimaryKey(key string, indexType string) {
	if p.Primarys.Contains(key) {
		return
	}

	p.Primarys.Append(key)

	if indexType != "" {
		if strings.ToUpper(indexType) == "DEFAULT" {
			indexType = ALGO
		}

		p.IndexType["PRIMARY"] = indexType
	}
}

func (p *TableParam) AddIndex(kind string, indexName string, indexType string, cols ...string) error {
	var im map[string]*cntr.Array[string]
	var ok bool

	if im, ok = p.indexMap[kind]; !ok {
		return errors.New(fmt.Sprintf("There is no type %s in indexMap.", kind))
	}

	var indexs *cntr.Array[string]

	if indexs, ok = im[indexName]; !ok {
		im[indexName] = cntr.NewArray[string]()
		indexs = im[indexName]
		p.IndexType[indexName] = indexType
	}

	for _, col := range cols {
		if !indexs.Contains(col) {
			indexs.Append(col)
		}
	}

	return nil
}

func (p *TableParam) RemoveIndex(kind string, indexName string, indexType string, colName string) error {
	var im map[string]*cntr.Array[string]
	var ok bool

	if im, ok = p.indexMap[kind]; !ok {
		return errors.New(fmt.Sprintf("There is no type %s in indexMap.", kind))
	}

	var indexs *cntr.Array[string]

	if indexs, ok = im[indexName]; !ok {
		return errors.New(fmt.Sprintf("There is no type %s in indexMap.", kind))
	}

	indexs.Remove(colName)
	return nil
}

func (p *TableParam) UpdateIndexName(source string, target string) {
	var idx int

	p.operateIndexMap(source, target, func(s1, s2 string, a *cntr.Array[string]) error {
		idx = a.Find(source)

		if idx != -1 {
			a.Elements[idx] = target
		}

		return nil
	})

	if indexType, ok := p.IndexType[source]; ok {
		p.IndexType[target] = indexType
		delete(p.IndexType, source)
	}
}

// 在 source 所在的 index 當中，都添加一個 target
func (p *TableParam) CloneColumnOfIndex(source string, target string) {
	isValid := false

	p.operateIndexMap(source, target, func(s1, s2 string, a *cntr.Array[string]) error {
		if a.Contains(source) && !a.Contains(target) {
			a.Append(target)
			isValid = true
		}
		return nil
	})

	if isValid {
		p.IndexType[target] = p.IndexType[source]
	}
}

func (p *TableParam) operateIndexMap(source string, target string, op func(string, string, *cntr.Array[string]) error) {
	var err error

	for kind, indexs := range p.indexMap {
		for indexName, cols := range indexs {
			err = op(kind, indexName, cols)
			if err != nil {
				utils.Error("執行 op 函式時發生錯誤, err: %+v", err)
			}
		}
	}
}

// 0: kind string, 1: indexName string, 2: indexType string, 3: cols *array.Array[string]
func (p *TableParam) IterIndexMap() *cntr.Iterator[*SqlIndex] {
	isUnSorted := false
	elements := []*SqlIndex{}
	var kind, indexName string
	var indexs map[string]*cntr.Array[string]
	var cols *cntr.Array[string]

	if isUnSorted {
		for kind, indexs = range p.indexMap {
			for indexName, cols = range indexs {
				elements = append(elements, NewSqlIndex(kind, indexName, p.IndexType[indexName], cols))
			}
		}

	} else {
		kinds := []string{}
		names := []string{}

		for kind = range p.indexMap {
			kinds = append(kinds, kind)
		}
		sort.Strings(kinds)

		for _, kind = range kinds {
			indexs = p.indexMap[kind]

			names = names[:0]
			for indexName = range indexs {
				names = append(names, indexName)
			}
			sort.Strings(names)

			for _, indexName = range names {
				cols = indexs[indexName]
				elements = append(elements, NewSqlIndex(kind, indexName, p.IndexType[indexName], cols))
			}
		}
	}

	return cntr.NewIterator(elements)
}

func (p *TableParam) GetIndexColumns(kind string, indexName string) *cntr.Array[string] {
	var im map[string]*cntr.Array[string]
	var ok bool

	if im, ok = p.indexMap[kind]; !ok {
		utils.Error(fmt.Sprintf("There is no type %s in indexMap.", kind))
		return nil
	}

	if _, ok := im[indexName]; !ok {
		utils.Error(fmt.Sprintf("There is no indexName %s in indexMap[%s].", indexName, kind))
		return nil
	}

	return im[indexName]
}

func (p *TableParam) GetAllColumns() *cntr.Array[string] {
	columns := cntr.NewArray[string]()
	it := p.IterIndexMap()
	var data *SqlIndex

	// 0: kind string, 1: indexName string, 2: indexType string, 3: cols *array.Array[string]
	for it.HasNext() {
		data = it.Next()

		for _, col := range data.Cols.Elements {
			if !columns.Contains(col) {
				columns.Append(col)
			}
		}
	}

	return columns
}

func (p TableParam) String() string {
	bs, err := json.Marshal(p.indexMap)
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("TableParam\nPrimarys: %+v", p.Primarys))
	if err != nil {
		utils.Error("Failed to marshal indexMap, err: %+v", err)
	} else {
		buffer.WriteString(fmt.Sprintf("\nIndexs: %+v", string(bs)))
	}
	return buffer.String()
}

func (p TableParam) Clone() *TableParam {
	clone := &TableParam{
		Primarys:  p.Primarys.Clone(),
		indexMap:  map[string]map[string]*cntr.Array[string]{},
		IndexType: map[string]string{},
	}

	for k1 := range p.indexMap {
		if _, ok := clone.indexMap[k1]; !ok {
			clone.indexMap[k1] = map[string]*cntr.Array[string]{}
		}

		for k2, v := range p.indexMap[k1] {
			clone.indexMap[k1][k2] = v.Clone()
		}
	}

	for k, v := range p.IndexType {
		clone.IndexType[k] = v
	}

	return clone
}

/*
ColumnParam 附加定義
1. default: AI(Auto Increment), NULL, NIL(不設置預設值), 其他(current_timestamp())
2. size: 定義 DB 欄位變數大小
3. type: 定義 DB 欄位變數類型
4. can_null: 是否可以為空值 (true: 可以 / false: 不可以)
5. primary_key: 是否為主鍵，數值為演算法(填入 default 則使用預設值)
6. comment: 註解內容
7. unsigned: 沒有負數？
8. update: ON UPDATE 時執行的函數
* 字串中需要空格可以使用 \t

# 調整 DB 欄位順序，應與 Proto 檔一致
# 建議參數修改順序
	1. Comment
	2. CanNull
	3. Type
	4. Size
	5. Unsigned
	6. Default
	7. PrimaryKey
	8. Update
*/

type ColumnParam struct {
	// 欄位編號
	FieldNumber int

	// 欄位名稱
	Name string

	// 原始欄位變數類型
	OriginType datatype.DataType

	// 欄位變數類型
	Type datatype.DataType

	// 欄位大小
	Size int32

	// 是否為主鍵
	IsPrimaryKey bool

	// 沒有負數?
	IsUnsigned bool

	// 能否為空
	CanNull bool

	// 預設值(如：NIL(無預設值), AI(AutoIncrement), NULL)
	Default string

	// ON UPDATE 時觸發的函式
	Update string

	// 演算法
	Algo string

	// 註解
	Comment string

	// 是否忽略此欄位(用於結構中有定義，但不希望成為資料表欄位時)
	IgnoreThis bool

	// 資料庫方言類型
	dial dialect.SQLDialect

	// 欄位設置表
	defineMap map[string]string

	Config *ColumnParamConfig
}

// 改以 DefineMap 的形式暫存各個欄位的額外設定值，考慮'先後順序'與'設定值之間的相互牽制'，
// 例如使用 AutoIncrement 的欄位要求必須是 PrimaryKey
func NewColumnParam(number int, name string, kind datatype.DataType, dial dialect.SQLDialect, tags ...string) *ColumnParam {
	param := &ColumnParam{
		FieldNumber:  number,
		Name:         name,
		Size:         0,
		IsPrimaryKey: false,
		IsUnsigned:   false,
		CanNull:      false,
		Update:       "",
		Algo:         "",
		Comment:      "",
		IgnoreThis:   false,
		dial:         dial,
		defineMap:    map[string]string{},
		Config:       &ColumnParamConfig{},
	}

	param.SetType(kind)
	param.SetDefault(dialect.GetDialect(param.dial).GetDefault(param.Type))

	// 根據 tag 內容對 Param 進行再定義
	for _, tag := range tags {
		param.parserConfig(tag)
	}

	// 根據 defineMap 對 Param 進行再定義
	param.Redefine()
	return param
}

func (p *ColumnParam) parserConfig(config string) {
	cfg, err := NewColumnParamConfig(config)
	if err != nil {
		return
	}
	p.Config.merge(cfg)
}

func (p *ColumnParam) SetName(name string) {
	p.Name = name
}

func (p *ColumnParam) SetType(dataType datatype.DataType) {
	// 儲存類型修正前的類型(確保全大寫，以利設置欄位的預設大小)
	p.OriginType = datatype.ToUpper(dataType)

	// 根據實際使用的資料庫，對變數類型作修正
	p.Type = dialect.GetDialect(p.dial).TypeOf(p.OriginType)

	// 根據 Type 、當前的 Size 以及 DB 本身的限制，對數值大小再定義
	p.SetSize(p.Size)
}

// 根據 Type、當前的 Size 以及 DB 本身的限制，對數值大小再定義
func (p *ColumnParam) SetSize(size int32) {
	p.Size = dialect.GetDialect(p.dial).SizeOf(p.Type, size)
}

func (p *ColumnParam) SetPrimaryKey(algo string) {
	p.IsPrimaryKey = algo != ""

	if p.IsPrimaryKey {
		p.IsPrimaryKey = true
		p.Algo = strings.ToUpper(algo)

		if p.Algo == "DEFAULT" {
			p.Algo = ALGO
		}

	} else {
		// 確保非 Primary Key 欄位不會被設置成 Auto Increment
		if p.Default == "AI" {
			p.Default = "NIL"
		}
	}
}

func (p *ColumnParam) SetUnsigned(isUnsigned bool) {
	p.IsUnsigned = isUnsigned
}

func (p *ColumnParam) SetDefault(defaultValue string) {
	d := strings.ToUpper(defaultValue)

	switch d {
	case "NULL", "NIL":
		p.Default = dialect.GetDialect(p.dial).GetDefault(p.Type)
	// AI, NIL 都設置為大寫
	case "AI":
		p.Default = d
	// 其他則和設置值相同，不修改大小寫
	default:
		p.Default = defaultValue
	}
}

func (p *ColumnParam) SetOnUpdate(update string) {
	if p.Default == "current_timestamp()" {
		p.Update = update
	}
}

func (p *ColumnParam) SetComment(comment string) {
	p.Comment = comment
}

// 以 Config 的形式暫存各個欄位的額外設定值，考慮'先後順序'與'設定值之間的相互牽制'，
// 例如使用 AutoIncrement 的欄位要求必須是 PrimaryKey
func (p *ColumnParam) LoadConfig(config *ColumnParamConfig) {
	p.Config.merge(config)
}

// 根據 defineMap 對 Param 進行再定義
func (p *ColumnParam) Redefine() {
	// var ok bool
	var kind string

	// ====================================================================================================
	// 資料庫欄位忽略此欄位
	// ====================================================================================================

	if strings.ToUpper(p.Config.Ignore) == "TRUE" {
		p.IgnoreThis = true
		return
	}

	// ====================================================================================================
	// Comment
	// ====================================================================================================

	if p.Config.Comment != "" {
		p.Comment = p.Config.Comment
	}

	// ====================================================================================================
	// Type
	// ====================================================================================================

	if p.Config.Type != "" {
		// 確保全大寫，以利設置欄位的預設大小
		p.Type = datatype.ToUpper(datatype.DataType(p.Config.Type))
	}

	// 根據實際使用的資料庫，對變數類型作修正
	p.Type = dialect.GetDialect(p.dial).TypeOf(p.Type)
	kind = dialect.GetDialect(p.dial).GetKind(p.Type)

	// ====================================================================================================
	// size
	// 根據 p.Type 、當前的 p.Size 以及 DB 本身的限制，對數值大小再定義
	// NOTE: 數值型用預設值，字串型沒有值，僅允許 VARCHAR 設定
	// ====================================================================================================

	// 僅允許 VARCHAR 設定 size
	if p.Config.Size > 0 && p.Type == datatype.VARCHAR {
		p.Size = dialect.GetDialect(p.dial).SizeOf(p.Type, int32(p.Config.Size))

	} else if p.OriginType == datatype.BOOL {
		// bool 類型的數據，p.Type 會被轉換成 TINYINT，因此這裡針對 bool 類型的大小作額外處理
		p.Size = dialect.GetDialect(p.dial).SizeOf(datatype.BOOL, 0)

	} else {
		// 使用該類型變數的預設大小
		p.Size = dialect.GetDialect(p.dial).SizeOf(p.Type, 0)
	}

	// ====================================================================================================
	// unsigned
	// ====================================================================================================

	if p.Config.Unsigned != "" {
		// 數值類型才能設置 "是否沒有負數" 這項屬性
		if kind == "INTEGER" || kind == "FLOAT" {
			p.IsUnsigned = strings.ToUpper(p.Config.Unsigned) == "TRUE"
		}
	}

	// ====================================================================================================
	// default
	// ====================================================================================================

	if p.Config.Default != "" {
		p.SetDefault(p.Config.Default)
	}

	// ====================================================================================================
	// primary_key
	// ====================================================================================================

	if p.Config.PrimaryKey != "" {
		p.IsPrimaryKey = true
		// p.CanNull = false
		p.Algo = strings.ToUpper(p.Config.PrimaryKey)

		if p.Algo == "DEFAULT" {
			p.Algo = ALGO
		}
	} else {
		// 確保非 Primary Key 欄位不會被設置成 Auto Increment
		if p.Default == "AI" {
			p.Default = dialect.GetDialect(p.dial).GetDefault(p.Type)
		}
	}

	// Default 欄位不可為 NULL
	if p.Default == "NULL" {
		p.Default = dialect.GetDialect(p.dial).GetDefault(p.Type)
	}

	// ====================================================================================================
	// update
	// ====================================================================================================

	if p.Config.Update != "" {
		switch p.Default {
		case "current_timestamp()":
			p.Update = p.Config.Update
		case "NIL", "NULL":
			fallthrough
		default:
			p.Update = ""
		}
	}

	// MAP 與 MESSAGE 將以超長字串形式儲存
	if p.OriginType == datatype.MAP || p.OriginType == datatype.MESSAGE {
		// fmt.Printf("OriginType: %s, Type: %s\n", param.OriginType, param.Type)
		if p.Size < 3000 {
			p.Size = 3000
		}
	}
}

func (p *ColumnParam) String() string {
	return fmt.Sprintf("%d) Name: %s, Type: %s, Size: %d, IsPrimaryKey: %v, Default: %s, Update: %s",
		p.FieldNumber, p.Name, p.Type, p.Size, p.IsPrimaryKey, p.Default, p.Update)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// ColumnParamSlice
////////////////////////////////////////////////////////////////////////////////////////////////////
type ColumnParamSlice []*ColumnParam

func (p *ColumnParamSlice) Len() int {
	return len(*p)
}

func (p *ColumnParamSlice) Less(i int, j int) bool {
	return (*p)[i].FieldNumber < (*p)[j].FieldNumber
}

func (p *ColumnParamSlice) Swap(i int, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *ColumnParamSlice) Sort() {
	sort.Sort(p)
}
