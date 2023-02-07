package proto

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/j32u4ukh/gosql/proto/ast"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
)

type Helper struct {
	Folder string
	Dial   dialect.SQLDialect
	DbName string
}

func (h *Helper) SetFolder(folder string) {
	h.Folder = folder
}

func (h *Helper) SetDial(dial dialect.SQLDialect) {
	h.Dial = dial
}

func (h *Helper) GetParamsByPath(folder string, tableName string) (*stmt.TableParam, []*stmt.ColumnParam, error) {
	path := fmt.Sprintf("%s.proto", filepath.Join(folder, tableName))
	_, err := os.Stat(path)

	if err != nil {
		return nil, nil, errors.Wrapf(err, "檔案不存在, path: %s", path)
	}

	tableParam, colParams, err := ast.GetProtoParams(path, h.Dial)

	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Failed to get params from %s.", path))
	}

	return tableParam, colParams, nil
}

func (h *Helper) GetParams(table_name string) (*stmt.TableParam, []*stmt.ColumnParam, error) {
	path := fmt.Sprintf("%s.proto", filepath.Join(h.Folder, table_name))
	_, err := os.Stat(path)

	if err != nil {
		return nil, nil, errors.Wrapf(err, "檔案不存在, path: %s", path)
	}

	tableParam, colParams, err := ast.GetProtoParams(path, h.Dial)

	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Failed to get params from %s.", path))
	}

	return tableParam, colParams, nil
}
