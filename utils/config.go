package utils

var GosqlConfig *Config

func init() {
	GosqlConfig = &Config{
		MariaVarcharSize:  3000,
		MysqlVarcharSize:  3000,
		SqliteVarcharSize: 3000,
		Params:            make(map[string]string),
	}
}

type Config struct {
	MariaVarcharSize  int32
	MysqlVarcharSize  int32
	SqliteVarcharSize int32
	Params            map[string]string
}
