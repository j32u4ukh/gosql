package ast

import (
	"fmt"
	"sort"
	"strings"
)

type Variable struct {
	Index     int
	ProtoName string // player_id
	GoName    string // PlayerId
	Type      string
}

func NewVariable(index int, name string, kind string) *Variable {
	return &Variable{
		Index:     index,
		ProtoName: name,
		GoName:    buildGoName(name),
		Type:      kind,
	}
}

func (v *Variable) String() string {
	return fmt.Sprintf("Variable(Index: %d, Name: %s, Type: %s)", v.Index, v.GoName, v.Type)
}

func buildGoName(protoName string) string {
	parts := strings.Split(protoName, "_")
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

type Variables []*Variable

func (v *Variables) Len() int {
	return len(*v)
}

func (v *Variables) Less(i int, j int) bool {
	return (*v)[i].Index < (*v)[j].Index
}

func (v *Variables) Swap(i int, j int) {
	(*v)[i], (*v)[j] = (*v)[j], (*v)[i]
}

func (v *Variables) Sort() {
	sort.Sort(v)
}
