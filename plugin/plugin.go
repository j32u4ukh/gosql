package plugin

type ISqlStruct interface {
	ToStmt() string
}
