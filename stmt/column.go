package stmt

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
)

type Column struct {
	// 欄位名稱
	Name string

	// 欄位變數類型
	// SQL 中的變數類型(string, Message, Slice, Map 等類型，預計以超長 VARCHAR 形式儲存)
	Type datatype.DataType

	// 欄位大小
	Size int32

	// 是否為主鍵
	IsPrimaryKey bool

	// 沒有負數?
	IsUnsigned bool

	// 能否為空(PrimaryKey 不可為空)
	CanNull bool

	// 預設值(如：無預設值(NIL), AutoIncrement(AI), NULL)
	Default string

	// ON UPDATE 時觸發的函式
	Update string

	// 演算法
	Algo string

	// 註解
	Comment string

	// 排序規則
	Collate string

	// 資料庫方言類型
	Dial dialect.SQLDialect

	// 是否忽略此欄位(用於結構中有定義，但不希望成為資料表欄位時)
	IgnoreThis bool
}

func NewColumn(param *ColumnParam) *Column {
	column := &Column{
		Name:         param.Name,
		Type:         param.Type,
		Size:         param.Size,
		IsPrimaryKey: param.IsPrimaryKey,
		IsUnsigned:   param.IsUnsigned,
		CanNull:      param.CanNull,
		Default:      param.Default,
		Update:       param.Update,
		Algo:         param.Algo,
		Comment:      param.Comment,
		Collate:      "",
		Dial:         param.dial,
		IgnoreThis:   param.IgnoreThis,
	}

	return column
}

// 讀取 DB 資訊來創建 Column 物件
func NewDbColumn(name string, kind string, isPrimaryKey bool, canNull bool, _defualt string, collate string, extra string, comment string) *Column {
	column := &Column{
		Name:         name,
		IsPrimaryKey: isPrimaryKey,
		CanNull:      canNull,
		Default:      _defualt,
		Comment:      comment,
	}

	// ====================================================================================================
	// COLUMN_TYPE
	// int(11) / smallint(5) unsigned / timestamp / ...
	// ====================================================================================================
	kind = strings.ToUpper(kind)
	elements := strings.Split(kind, " ")

	kind = elements[0]
	left := strings.Index(kind, "(")
	// fmt.Printf("NewDbColumn kind: %s, left: %d\n", kind, left)

	if left == -1 {
		column.Type = datatype.DataType(kind)
		column.Size = 0
	} else {
		column.Type = datatype.DataType(kind[:left])
		size, _ := strconv.Atoi(kind[left+1 : len(kind)-1])
		column.Size = int32(size)
	}

	if len(elements) == 2 {
		kind = elements[1]

		if kind == "UNSIGNED" {
			column.IsUnsigned = true
		}
	}
	// ====================================================================================================

	if dialect.GetDialect(column.Dial).IsSortable(column.Type) {
		column.Collate = collate
	}

	if extra == "auto_increment" {
		column.Default = "AI"

	} else if strings.HasPrefix(extra, "on update") {
		update := strings.Trim(extra[9:], " ")

		switch update {
		case "current_timestamp()":
			column.Update = "current_timestamp()"
		}
	}

	return column
}

func (c *Column) SetPrimaryKey() {
	c.IsPrimaryKey = true
	c.CanNull = false

	if c.Default == "" {
		c.Default = dialect.GetDialect(c.Dial).GetDefault(c.Type)
	}

	if c.Algo == "" {
		c.Algo = "BTREE"
	}
}

func (c *Column) SetCollate(collate string) {
	if dialect.GetDialect(c.Dial).IsSortable(c.Type) {
		c.Collate = collate
	}
}

// 用於生成表格的欄位資訊
func (c *Column) GetInfo() string {
	var result string
	// fmt.Printf("GetInfo | Type: %s, Size: %d\n", c.Type, c.Size)

	if c.Size > 0 {
		result = fmt.Sprintf("%s(%d)", c.Type, c.Size)
	} else {
		result = string(c.Type)
	}

	if c.IsUnsigned {
		result += " UNSIGNED"
	}

	result += " NOT NULL"

	switch c.Default {
	case "AI":
		// PrimaryKey 才允許設為 AUTO_INCREMENT
		if c.IsPrimaryKey {
			result += " AUTO_INCREMENT"
		}
	case "NULL":
		// 在讀取 Param 的階段已排除 "不可為空的欄位，預設值卻是 NULL 的情況"
		result += " DEFAULT NULL"
	case "NIL":
		// 沒有預設值
	case "":
	default:
		result += fmt.Sprintf(" DEFAULT %s", c.Default)

		if c.Update != "" {
			result += fmt.Sprintf(" ON UPDATE %s", c.Update)
		}
	}

	if c.Comment != "" {
		result += fmt.Sprintf(" COMMENT '%s'", c.Comment)
	}

	//  COLLATE 'utf8mb4_unicode_ci'
	if c.Collate != "" {
		result += fmt.Sprintf(" COLLATE '%s'", c.Collate)
	}

	return result
}

func (c *Column) IsEquals(other *Column) bool {
	if c.Name != other.Name {
		fmt.Printf("Names are different. c.Name: %s, other.Name: %s\n", c.Name, other.Name)
		return false
	}

	if c.IsPrimaryKey != other.IsPrimaryKey {
		return false
	}

	if c.GetInfo() != other.GetInfo() {
		return false
	}

	return true
}

func (c *Column) ToStmt() string {
	return fmt.Sprintf("`%s` %s", c.Name, c.GetInfo())
}

func (c *Column) String() string {
	return c.ToStmt()
}

func (c *Column) Clone() *Column {
	clone := &Column{
		Name:         c.Name,
		Type:         c.Type,
		Size:         c.Size,
		IsPrimaryKey: c.IsPrimaryKey,
		CanNull:      c.CanNull,
		Default:      c.Default,
		Algo:         c.Algo,
		Comment:      c.Comment,
		Collate:      c.Collate,
	}
	return clone
}

func FormatColumns(columns []string, mode byte) (string, error) {
	length := len(columns)

	switch length {
	case 0:
		switch mode {
		case DistinctSelect:
			return "", errors.New("You need to specify the columns when using DISTINCT.")
		case CountSelect:
			return "COUNT(*)", nil
		case CountDistinctSelect:
			return "", errors.New("You need to specify the columns when using DISTINCT.")
		case NormalSelect:
			fallthrough
		default:
			return "*", nil
		}
	case 1:
		switch mode {
		case DistinctSelect:
			return fmt.Sprintf("DISTINCT `%s`", columns[0]), nil
		case CountSelect:
			return fmt.Sprintf("COUNT(`%s`)", columns[0]), nil
		case CountDistinctSelect:
			return fmt.Sprintf("COUNT(DISTINCT `%s`)", columns[0]), nil
		case NormalSelect:
			fallthrough
		default:
			return columns[0], nil
		}
	default:
		temps := []string{}
		for _, column := range columns {
			temps = append(temps, fmt.Sprintf("`%s`", column))
		}
		result := strings.Join(temps, ", ")
		switch mode {
		case DistinctSelect:
			return fmt.Sprintf("DISTINCT %s", result), nil
		case CountSelect:
			return fmt.Sprintf("COUNT(%s)", result), nil
		case CountDistinctSelect:
			return fmt.Sprintf("COUNT(DISTINCT %s)", result), nil
		case NormalSelect:
			fallthrough
		default:
			return result, nil
		}
	}
}
