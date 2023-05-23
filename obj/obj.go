package obj

type ISqlStruct interface {
	ToStmt() string
}
