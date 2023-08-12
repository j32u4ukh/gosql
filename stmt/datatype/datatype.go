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
	BIT       DataType = "BIT"

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
	JSON       DataType = "JSON"
	UUID       DataType = "UUID"

	// Time
	DATE      DataType = "DATE"
	TIME      DataType = "TIME"
	YEAR      DataType = "YEAR"
	DATETIME  DataType = "DATETIME"
	TIMESTAMP DataType = "TIMESTAMP"

	// Binary
	BINARY     DataType = "BINARY"
	VARBINARY  DataType = "VARBINARY"
	TINYBLOB   DataType = "TINYBLOB"
	BLOB       DataType = "BLOB"
	MEDIUMBLOB DataType = "MEDIUMBLOB"
	LONGBLOB   DataType = "LONGBLOB"

	// Spatial
	POINT              DataType = "POINT"
	LINESTRING         DataType = "LINESTRING"
	POLYGON            DataType = "POLYGON"
	GEOMETRY           DataType = "GEOMETRY"
	MULTIPOINT         DataType = "MULTIPOINT"
	MULTILINESTRING    DataType = "MULTILINESTRING"
	MULTIPOLYGON       DataType = "MULTIPOLYGON"
	GEOMETRYCOLLECTION DataType = "GEOMETRYCOLLECTION"

	// Other
	UNKNOWN DataType = "UNKNOWN"
	ENUM    DataType = "ENUM"
	SET     DataType = "SET"

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
	case TINYINT, SMALLINT, MEDIUMINT, INT, BIGINT, BIT,
		FLOAT, DOUBLE, DEMICAL,
		VARCHAR, CHAR, TINYTEXT, TEXT, MEDIUMTEXT, LONGTEXT, JSON, UUID,
		DATE, TIME, YEAR, DATETIME, TIMESTAMP,
		BINARY, VARBINARY, TINYBLOB, BLOB, MEDIUMBLOB, LONGBLOB, GEOMETRYCOLLECTION,
		UNKNOWN, ENUM, SET,
		MAP, MESSAGE,
		INT32, INT64, UINT32, UINT64, SINT32, SINT64,
		FIXED32, FIXED64, SFIXED32, SFIXED64,
		BOOL, STRING, BYTES, "":
		return dt
	default:
		return MESSAGE
	}
}

// 將 DataType 轉換為全大寫
func ToUpper(dt DataType) DataType {
	return DataType(strings.ToUpper(string(dt)))
}
