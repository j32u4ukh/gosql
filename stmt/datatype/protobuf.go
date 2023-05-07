package datatype

import "strings"

// 	Protobuf 中的變數類型，轉為通用的變數類型 DataType
func ProtoToDataType(kind string) DataType {
	kind = strings.ToUpper(kind)
	switch kind {
	case "INT32":
		return INT32
	case "INT64":
		return INT64
	case "MAP":
		return MAP
	case "BOOL":
		return BOOL
	case "STRING":
		return STRING
	case "TIMESTAMP":
		return TIMESTAMP
	case "MESSAGE":
		fallthrough
	default:
		return MESSAGE
	}
}
