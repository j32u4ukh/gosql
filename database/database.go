package database

import (
	"bytes"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // 未直接使用此套件，但背後執行時有用到，因此以此形式 import
	"github.com/j32u4ukh/gosql/utils/cntr"
	"github.com/pkg/errors"
)

type SqlResult struct {
	// LastInsertId returns the integer generated by the database
	// in response to a command. Typically this will be from an
	// "auto increment" column when inserting a new row. Not all
	// databases support this feature, and the syntax of such
	// statements varies.
	LastInsertId int64

	// RowsAffected returns the number of rows affected by an
	// update, insert, or delete. Not every database or database
	// driver may support this.
	RowsAffected int64

	// 欄位名稱
	Columns []string

	// 數據筆數
	NRow int32

	// 欄位個數
	NColumn int32

	Datas [][]string
}

func (sr SqlResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{\n")
	buffer.WriteString(fmt.Sprintf("\t'LastInsertId': %d,\n", sr.LastInsertId))
	buffer.WriteString(fmt.Sprintf("\t'RowsAffected': %d,\n", sr.RowsAffected))
	if sr.NColumn != 0 {
		buffer.WriteString(fmt.Sprintf("\t'Columns': %s,\n", cntr.SliceToString(sr.Columns)))
		buffer.WriteString(fmt.Sprintf("\t'NRow': %d,\n", sr.NRow))
		buffer.WriteString(fmt.Sprintf("\t'NColumn': %d,\n", sr.NColumn))
	}
	buffer.WriteString("}")
	return buffer.String()
}

type Database struct {
	db     *sql.DB
	idx    int32
	DbName string
}

func NewDatabase(idx int32, user string, password string, host string, port int) (*Database, error) {
	d := &Database{idx: idx}
	var err error
	var connection = fmt.Sprintf("%s:%s@tcp(%s:%d)/?allowNativePasswords=true", user, password, host, port)
	d.db, err = sql.Open("mysql", connection)

	if err != nil {
		return nil, errors.Wrapf(err, "開啟資料庫失敗 (%s:%d)", host, port)
	}

	err = d.db.Ping()

	if err != nil {
		return nil, errors.Wrapf(err, "資料庫連線測試失敗 (%s:%d)", host, port)
	}

	// See "Important settings" section.
	d.db.SetConnMaxLifetime(3 * time.Minute)
	d.db.SetMaxOpenConns(10)
	d.db.SetMaxIdleConns(10)
	return d, nil
}

// 同時連結多個 database 時，由於維持多個連線，可以會發生錯誤。需再次呼叫此函式，告訴程式要連哪一個 database。
func (d *Database) Use(dbName string) error {
	_, err := d.db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))

	if err != nil {
		return errors.Wrapf(err, "Failed to create database %s.", dbName)
	}

	_, err = d.db.Exec(fmt.Sprintf("USE `%s`", dbName))

	if err != nil {
		return errors.Wrapf(err, "Failed to use database %s.", dbName)
	}

	d.DbName = dbName
	return nil
}

func (d *Database) IsTableExists(sql string) (bool, error) {
	rows, err := d.db.Query(sql)
	defer rows.Close()

	if err != nil {
		fmt.Printf("Error: %+v", err)
		return false, err
	}

	return rows.Next(), nil
}

func (d *Database) Exec(sql string, args ...any) (*SqlResult, error) {
	result, err := d.db.Exec(sql, args...)

	if err != nil {
		return nil, errors.Wrapf(err, "執行 SQL 語法時發生錯誤, sql: %s, args: %+v", sql, args)
	}

	gr := &SqlResult{}

	// 並非所有資料庫皆支援這兩項數值，因此不處理返回的錯誤
	gr.LastInsertId, _ = result.LastInsertId()
	gr.RowsAffected, _ = result.RowsAffected()
	return gr, nil
}

func (d *Database) Query(sql string, args ...any) (*SqlResult, error) {
	rows, err := d.db.Query(sql)

	if err != nil {
		return nil, errors.Wrapf(err, "執行 SQL 語法時發生錯誤, sql: %s, args: %+v", sql, args)
	}

	sr := &SqlResult{
		NRow:    0,
		NColumn: 0,
		Datas:   [][]string{},
	}

	var values []interface{}
	var dest []interface{}
	var i int32

	for rows.Next() {
		if sr.NColumn == 0 {
			sr.Columns, err = rows.Columns()
			sr.NColumn = int32(len(sr.Columns))
			dest = make([]interface{}, sr.NColumn)
			values = make([]interface{}, sr.NColumn)

			// 初始化存放第一筆數據的空間
			for i = 0; i < sr.NColumn; i++ {
				dest[i] = &values[i]
			}
		}

		err = rows.Scan(dest...)

		if err != nil {
			return nil, errors.Wrap(err, "讀取數據時發生錯誤")
		}

		sr.Datas = append(sr.Datas, make([]string, sr.NColumn))

		for i = 0; i < sr.NColumn; i++ {
			// 數據為 NULL 時，填入空字串
			if values[i] == nil {
				sr.Datas[sr.NRow][i] = ""
			} else {
				// 將讀取到的數據放入 GdoResult 中
				sr.Datas[sr.NRow][i] = string(values[i].([]byte))
			}

			// 重置下一筆數據的存放空間
			values[i] = new(any)
			dest[i] = &values[i]
		}

		sr.NRow++
	}

	return sr, nil
}

func (d *Database) Close() {
	Close(d.idx)
}

func (d *Database) close() error {
	return d.db.Close()
}
