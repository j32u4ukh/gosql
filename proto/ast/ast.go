package ast

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/datatype"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
	protoparser "github.com/yoheimuta/go-protoparser/v4"
	ud "github.com/yoheimuta/go-protoparser/v4/interpret/unordered"
)

// NOTE: 參考 https://github.com/yoheimuta/go-protoparser
func GetProtoParams(path string, sqlDial dialect.SQLDialect) (*stmt.TableParam, []*stmt.ColumnParam, error) {
	reader, err := os.Open(path)
	defer reader.Close()

	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Failed to open %s.\n", path))
	}

	got, err := protoparser.Parse(
		reader,
		protoparser.WithFilename(filepath.Base(path)),
	)

	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Failed to parse %s.\n", path))
	}

	// 印出抽象語法樹整體結構
	// printStructure(got)

	var unorder *ud.Proto
	var tag string
	var idx int
	unorder, err = protoparser.UnorderedInterpret(got)
	msg := unorder.ProtoBody.Messages[0]
	tableParam := stmt.NewTableParam()

	for _, comment := range msg.Comments {
		tag = comment.Raw
		idx = strings.Index(tag, "//")
		tag = tag[idx+2:]
		tag = strings.Trim(tag, " ")
		tableParam.ParserConfig(tag)
	}

	var param *stmt.ColumnParam
	var tags []string
	body := msg.MessageBody

	// 以 stmt.ColumnParamSlice 封裝 []*stmt.ColumnParam{}，以實作排序介面
	var colParams stmt.ColumnParamSlice
	colParams = []*stmt.ColumnParam{}
	// fmt.Printf("GetProtoParams | Messages: %+v\n", body.Messages)

	for _, dict := range body.Maps {
		tags = []string{}

		if len(dict.Comments) > 0 {
			for _, comment := range dict.Comments {
				tag = comment.Raw
				idx = strings.Index(tag, "//")
				tag = tag[idx+2:]
				tag = strings.Trim(tag, " ")
				tags = append(tags, tag)
			}
		}

		// fmt.Printf("GetProtoParams | dict: %+v\n", dict)
		idx, _ = strconv.Atoi(dict.FieldNumber)
		param = stmt.NewColumnParam(
			idx,
			dict.MapName,
			datatype.MAP,
			sqlDial,
			tags...,
		)
		colParams = append(colParams, param)
	}

	for _, filed := range body.Fields {
		tags = []string{}

		if len(filed.Comments) > 0 {
			for _, comment := range filed.Comments {
				tag = comment.Raw
				idx = strings.Index(tag, "//")
				tag = tag[idx+2:]
				tag = strings.Trim(tag, " ")
				tags = append(tags, tag)
			}
		}

		// fmt.Printf("GetProtoParams | filed: %+v\n", filed)
		idx, _ = strconv.Atoi(filed.FieldNumber)
		param = stmt.NewColumnParam(
			idx,
			filed.FieldName,
			filed.Type,
			sqlDial,
			tags...,
		)
		colParams = append(colParams, param)
	}

	colParams.Sort()
	return tableParam, colParams, nil
}

func GetVariableParams(path string) ([]*Variable, error) {
	reader, err := os.Open(path)
	defer reader.Close()

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to open %s.\n", path))
	}

	got, err := protoparser.Parse(
		reader,
		protoparser.WithFilename(filepath.Base(path)),
	)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to parse %s.\n", path))
	}

	// 印出抽象語法樹整體結構
	// printStructure(got)

	var unorder *ud.Proto

	var idx int
	unorder, err = protoparser.UnorderedInterpret(got)
	msg := unorder.ProtoBody.Messages[0]
	body := msg.MessageBody

	var variabes Variables
	var variable *Variable
	variabes = []*Variable{}

	for _, dict := range body.Maps {
		fmt.Printf("GetVariableParams | dict: %+v\n", dict)
		idx, _ = strconv.Atoi(dict.FieldNumber)
		variable = NewVariable(idx, dict.MapName, datatype.MAP)
		variabes = append(variabes, variable)
	}

	for _, filed := range body.Fields {
		fmt.Printf("GetVariableParams | filed: %+v\n", filed)
		idx, _ = strconv.Atoi(filed.FieldNumber)
		variable = NewVariable(idx, filed.FieldName, filed.Type)
		variabes = append(variabes, variable)
	}

	variabes.Sort()
	return variabes, nil
}

// 印出抽象語法樹整體結構
func printStructure(got any) {
	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal, err %v\n", err)
	}
	fmt.Printf("printStructure | gotJSON:\n%s\n", gotJSON)
}
