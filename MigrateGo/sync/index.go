package sync

import (
	"fmt"
)

func (s *Synchronize) getDropIndex() []string {
	stmts := []string{}
	if s.dropList.Contains("PRIMARY KEY") {
		stmts = append(stmts, "DROP PRIMARY KEY")
		s.dropList.Remove("PRIMARY KEY")
	}
	for _, d := range s.dropList.Elements {
		stmts = append(stmts, fmt.Sprintf("DROP INDEX `%s`", d))
	}
	return stmts
}

func (s *Synchronize) getAddIndex() []string {
	stmts := []string{}
	// ADD PRIMARY KEY (`indx`) USING BTREE
	// TODO: 允許多筆 PRIMARY KEY
	if values, ok := s.indexMap["PRIMARY KEY"]; ok {
		stmts = append(stmts, fmt.Sprintf("ADD PRIMARY KEY (%s) USING %s", values[1], values[2]))
		delete(s.indexMap, "PRIMARY KEY")
	}
	// ADD INDEX `索引 2` (`create_time`, `gold_type`) USING BTREE
	// ADD UNIQUE INDEX `索引 2` (`create_time`, `gold_type`) USING BTREE
	for key, values := range s.indexMap {
		if values[0] == "UNIQUE" {
			values[0] = "UNIQUE INDEX"
		}
		stmts = append(stmts, fmt.Sprintf("ADD %s `%s` (%s) USING %s", values[0], key, values[1], values[2]))
	}
	return stmts
}
