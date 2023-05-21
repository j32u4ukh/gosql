package gosql

import (
	"reflect"

	root "github.com/j32u4ukh/gosql"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Table struct {
	*root.Table
	numFiled int
}

func (t *Table) InitByProtoMessage(pm protoreflect.ProtoMessage) {
	if t.numFiled == -1 {
		rt := reflect.TypeOf(pm).Elem()
		t.numFiled = rt.NumField()
	}
}
