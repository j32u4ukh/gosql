package sync

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/j32u4ukh/gosql/database"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	/* 同步模式
	1: Proto 檔 -> Db 中的 Table
	2: Db 中的 Table -> 另一個 Db 中的 Table
	*/
	Mode SyncMode `yaml:"Mode"`
	// ========== 來源 ==========
	// 來源 Proto 檔根目錄
	ProtoFolder string `yaml:"ProtoFolder"`
	// 來源 Database (同步後，應和這個相同)
	FromDatabase *database.DatabaseConfig `yaml:"FromDatabase"`
	// 來源表格名稱(Proto 檔名稱)
	FromTable string `yaml:"FromTable"`
	// ========== 目標 ==========
	// 目標 Database (同步對象，確保當前資料庫結構與來源相同)
	ToDatabase *database.DatabaseConfig `yaml:"ToDatabase"`
	ToTable    string                   `yaml:"ToTable"`
	// ========== 若表格結構有差異，是否印出即將執行的指令? ==========
	Print bool `yaml:"Print"`
	// ========== 若表格結構有差異，是否執行同步 ==========
	Sync bool `yaml:"Sync"`
	// ========== 若資料庫中沒有表格，是否自動生成表格 ==========
	Generate bool `yaml:"Generate"`
}

func NewConfig() *Config {
	c := &Config{
		Mode:         NoneMode,
		ProtoFolder:  "",
		FromDatabase: nil,
		FromTable:    "",
		ToDatabase:   nil,
		ToTable:      "",
		Print:        false,
		Sync:         false,
		Generate:     false,
	}
	return c
}

func (c *Config) LoadFile(path string) error {
	var b []byte
	var err error
	b, err = ioutil.ReadFile(path)

	if err != nil {
		return errors.Wrapf(err, "Failed to read file %s.", path)
	}

	err = yaml.Unmarshal(b, c)

	if err != nil {
		return errors.Wrapf(err, "Failed to unmarshal data.")
	}

	// 根據同步模式，對一些允許省略的部分初始化，以及檢查所需欄位是否都已定義
	err = c.CheckParams()

	if err != nil {
		return errors.Wrapf(err, "Failed to initialize.")
	}

	return nil
}

func (c *Config) CheckParams() error {
	switch c.Mode {
	case NoneMode:
		return errors.New("同步模式必須定義")
	case ProtoToDbMode:
		if c.ProtoFolder == "" {
			return errors.New("Proto 檔資料夾不可為空")
		}
		if c.FromTable == "" {
			return errors.New("來源表格名稱不可為空")
		}
	case DbToDbMode:
		if c.FromDatabase == nil {
			return errors.New("來源資料庫必須定義")
		}
		if c.FromTable == "" {
			return errors.New("來源表格名稱不可為空")
		}
		if c.ToDatabase == nil {
			c.ToDatabase = c.FromDatabase
		}
		if c.ToTable == "" {
			c.ToTable = c.FromTable
		}
	}
	return nil
}

func (c *Config) String() string {
	bs, _ := json.Marshal(c)
	return fmt.Sprintf("Config(%s)", string(bs))
}
