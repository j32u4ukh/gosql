package datatype

import "strings"

type DataType string

const (
	// Integer
	TINYINT   DataType = "TINYINT"
	SMALLINT  DataType = "SMALLINT"
	MEDIUMINT DataType = "MEDIUMINT"
	INT       DataType = "INT"
	BIGINT    DataType = "BIGINT"

	// Float
	FLOAT   DataType = "FLOAT"
	DOUBLE  DataType = "DOUBLE"
	DEMICAL DataType = "DEMICAL"

	// Text
	VARCHAR    DataType = "VARCHAR"
	CHAR       DataType = "CHAR"
	TINYTEXT   DataType = "TINYTEXT"
	TEXT       DataType = "TEXT"
	MEDIUMTEXT DataType = "MEDIUMTEXT"
	LONGTEXT   DataType = "LONGTEXT"

	// Time
	DATE      DataType = "DATE"
	TIME      DataType = "TIME"
	YEAR      DataType = "YEAR"
	DATETIME  DataType = "DATETIME"
	TIMESTAMP DataType = "TIMESTAMP"

	//////////////////////////////////////////////////
	// 非 SQL 類型 (Protobuf)
	//////////////////////////////////////////////////
	MAP      DataType = "MAP"
	MESSAGE  DataType = "MESSAGE"
	INT32    DataType = "INT32"
	INT64    DataType = "INT64"
	UINT32   DataType = "UINT32"
	UINT64   DataType = "UINT64"
	SINT32   DataType = "SINT32"
	SINT64   DataType = "SINT64"
	FIXED32  DataType = "FIXED32"
	FIXED64  DataType = "FIXED64"
	SFIXED32 DataType = "SFIXED32"
	SFIXED64 DataType = "SFIXED64"
	BOOL     DataType = "BOOL"
	STRING   DataType = "STRING"
	BYTES    DataType = "BYTES"
)

// 若為 MAP 或 proto 原生變數，則直接返回，其他則歸類為 MESSAGE
func GetOriginType(Type string) DataType {
	var dt DataType = DataType(Type)
	switch dt {
	case TINYINT, SMALLINT, MEDIUMINT, INT, BIGINT,
		FLOAT, DOUBLE, DEMICAL,
		VARCHAR, CHAR, TINYTEXT, TEXT, MEDIUMTEXT, LONGTEXT,
		DATE, TIME, YEAR, DATETIME, TIMESTAMP,
		MAP, MESSAGE,
		INT32, INT64, UINT32, UINT64, SINT32, SINT64,
		FIXED32, FIXED64, SFIXED32, SFIXED64,
		BOOL, STRING, "":
		return dt
	default:
		return MESSAGE
	}
}

// 將 DataType 轉換為全大寫
func ToUpper(dt DataType) DataType {
	return DataType(strings.ToUpper(string(dt)))
}
