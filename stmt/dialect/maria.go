package dialect

import (
	"fmt"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/utils"
)

type maria struct {
	KindMap map[string]*cntr.Array[string]
}

func init() {
	RegisterDialect(MARIA, &maria{
		KindMap: map[string]*cntr.Array[string]{
			"INTEGER": cntr.NewArray(
				"TINYINT",
				"SMALLINT",
				"MEDIUMINT",
				"INT",
				"BIGINT",
			),
			"FLOAT": cntr.NewArray(
				"FLOAT",
				"DOUBLE",
				"DEMICAL",
			),
			"STRING": cntr.NewArray(
				"VARCHAR",
				"CHAR",
				"TINYTEXT",
				"TEXT",
				"MEDIUMTEXT",
				"LONGTEXT",
			),
			"TIME": cntr.NewArray(
				"DATE",
				"TIME",
				"YEAR",
				"DATETIME",
				"TIMESTAMP",
			),
		},
	})
}

// 變數類型，轉為 SQL 中的變數類型
func (s *maria) TypeOf(dataType datatype.DataType) datatype.DataType {
	switch dataType {
	case datatype.TINYINT, datatype.SMALLINT, datatype.MEDIUMINT, datatype.INT, datatype.BIGINT,
		datatype.FLOAT, datatype.DOUBLE, datatype.DEMICAL,
		datatype.VARCHAR, datatype.CHAR,
		datatype.TINYTEXT, datatype.TEXT, datatype.MEDIUMTEXT, datatype.LONGTEXT,
		datatype.DATE, datatype.TIME, datatype.YEAR, datatype.DATETIME, datatype.TIMESTAMP:
		return dataType
	case datatype.INT32, datatype.UINT32, datatype.SINT32, datatype.FIXED32, datatype.SFIXED32:
		return datatype.INT
	case datatype.INT64, datatype.UINT64, datatype.SINT64, datatype.FIXED64, datatype.SFIXED64:
		return datatype.BIGINT
	case datatype.BOOL:
		return datatype.TINYINT
	case datatype.STRING:
		fallthrough
	default:
		// fmt.Printf("(s *maria) TypeOf | dataType: %s\n", dataType)
		return datatype.VARCHAR
	}
}

// 根據 dataType 、當前的 size 以及 DB 本身的限制，對數值大小再定義
func (s *maria) SizeOf(dataType datatype.DataType, size int32) int32 {
	if size <= 0 {
		switch dataType {
		case datatype.BOOL:
			return 1
		case datatype.TINYINT:
			return 4
		case datatype.SMALLINT:
			return 6
		case datatype.MEDIUMINT:
			return 9
		case datatype.INT:
			return 11
		case datatype.BIGINT:
			return 20
		case datatype.VARCHAR:
			return utils.GosqlConfig.MariaVarcharSize
		case datatype.CHAR:
			return 50
		}
	}

	switch dataType {
	// ==================================================
	// DB 本身有其預設值，無法自定義大小的類型，一律回傳 0
	case datatype.TIMESTAMP, datatype.DOUBLE, datatype.FLOAT,
		datatype.TINYTEXT, datatype.TEXT, datatype.MEDIUMTEXT, datatype.LONGTEXT:
		return 0
	// ==================================================
	default:
		return size
	}
}

func (s *maria) GetDefault(dataType datatype.DataType) string {
	switch dataType {
	case datatype.TINYINT, datatype.SMALLINT, datatype.MEDIUMINT, datatype.INT, datatype.BIGINT,
		datatype.FLOAT, datatype.DOUBLE, datatype.DEMICAL:
		return "0"
	case datatype.VARCHAR, datatype.CHAR, datatype.TINYTEXT, datatype.TEXT, datatype.MEDIUMTEXT, datatype.LONGTEXT:
		return "''"
	// current_timestamp()
	case datatype.DATE:
		return "'1970-01-01'"
	case datatype.TIME:
		return "'00:00:00'"
	case datatype.YEAR:
		return "1970"
	case datatype.DATETIME:
		return "'1970-01-01 00:00:00'"
	case datatype.TIMESTAMP:
		return "'1970-01-01 00:00:01'"
	default:
		panic(fmt.Sprintf("Invalid variable type: %s.", dataType))
	}
}

// 是否為數值類型
func (s *maria) IsSortable(kind datatype.DataType) bool {
	return s.KindMap["STRING"].Contains(string(datatype.ToUpper(kind)))
}

// 判斷變數類型(integer, float, text, ...)
func (s *maria) GetKind(kind datatype.DataType) string {
	for k, v := range s.KindMap {
		if v.Contains(string(kind)) {
			return k
		}
	}
	return "Unknown"
}
