package stmt

import (
	"fmt"
	"strings"
)

type SelectItem struct {
	Name  string
	Alias string
}

func NewSelectItem(name string) *SelectItem {
	return &SelectItem{Name: name}
}

func (s *SelectItem) UseBacktick() *SelectItem {
	s.Name = fmt.Sprintf("`%s`", s.Name)
	return s
}

func (s *SelectItem) SetAlias(alias string) *SelectItem {
	s.Alias = alias
	return s
}

func (s *SelectItem) Count() *SelectItem {
	s.Name = fmt.Sprintf("COUNT(%s)", s.Name)
	return s
}

func (s *SelectItem) Distinct() *SelectItem {
	s.Name = fmt.Sprintf("DISTINCT %s", s.Name)
	return s
}

func (s *SelectItem) Concat(elements ...string) *SelectItem {
	s.Name = fmt.Sprintf("CONCAT(%s)", strings.Join(elements, ", "))
	return s
}

func (s *SelectItem) ToStmt() string {
	result := s.Name
	if s.Alias != "" {
		result = fmt.Sprintf("%s AS %s", result, s.Alias)
	}
	return result
}

func FormatColumns(items []*SelectItem) (string, error) {
	length := len(items)
	if length == 0 {
		return "*", nil
	}
	results := []string{}
	for _, item := range items {
		results = append(results, item.ToStmt())
	}
	return strings.Join(results, ", "), nil
}
