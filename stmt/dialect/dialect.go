package dialect

type SQLDialect byte

const (
	MARIA SQLDialect = iota
	MYSQL
	SQLITE3
)

var dialectsMap = map[SQLDialect]Dialect{}

// MySQL 和 MariaDb 中有部分語法上的差異
type Dialect interface {
	// 變數類型，轉為 SQL 中的變數類型
	TypeOf(dataType string) string

	// 根據 dataType 、當前的 size 以及 DB 本身的限制，對數值大小再定義
	SizeOf(dataType string, size int32) int32

	// 取得 dataType 類型變數的預設值
	GetDefault(dataType string) string

	// Protobuf 中的變數類型，轉為 SQL 中的變數類型
	ProtoTypeOf(string) string

	// 是否為可排序的類型
	IsSortable(kind string) bool

	// 判斷變數類型(integer, float, text, ...)
	GetKind(kind string) string
}

// 註冊各資料庫語言的方言物件
func RegisterDialect(name SQLDialect, dialect Dialect) {
	dialectsMap[name] = dialect
}

// 取得資料庫語言的方言物件
func GetDialect(name SQLDialect) Dialect {
	return dialectsMap[name]
}
