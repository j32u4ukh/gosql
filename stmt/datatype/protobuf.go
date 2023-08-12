package datatype

import "strings"

// 	Protobuf 中的變數類型，轉為通用的變數類型 DataType
func ProtoToDataType(kind string) DataType {
	var dt DataType = DataType(strings.ToUpper(kind))
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
