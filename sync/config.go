package sync

import (
	"io/ioutil"

	"github.com/j32u4ukh/gosql/database"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	/* 同步模式
	1: N 個 Proto 檔 -> Db 中的 N 張 Table
	2: Db 中的 N 張 Table -> 另一個 Db 中的 N 張 Table
	3: 所有 Proto 檔 -> Db 中的所有 Table (遇缺不補)
	4: Db 中的所有 Table -> 另一個 Db 中的所有 Table (遇缺不補)
	*/
	Mode int32 `yaml:"Mode"`
	// ========== 來源 ==========
	// 來源 Proto 檔根目錄
	ProtoFolder string `yaml:"ProtoFolder"`
	// 來源 Database (同步後，應和這個相同)
	FromDatabase *database.DatabaseConfig `yaml:"FromDatabase"`
	// 來源表格名稱(Proto 檔名稱)
	FromTables []string `yaml:"FromTables"`
	// ========== 目標 ==========
	// 目標 Database (同步對象，確保當前資料庫結構與來源相同)
	ToDatabase   *database.DatabaseConfig `yaml:"ToDatabase"`
	ToTables     []string                 `yaml:"ToTables"`
	AutoGenerate bool                     `yaml:"AutoGenerate"`
}

func NewConfig(path string) (*Config, error) {
	var b []byte
	var err error
	b, err = ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read file %s.", path)
	}

	c := &Config{
		Mode:         -1,
		FromDatabase: nil,
		ProtoFolder:  "",
		FromTables:   nil,
		ToDatabase:   nil,
		ToTables:     nil,
	}
	err = yaml.Unmarshal(b, c)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to unmarshal data.")
	}

	// 根據同步模式，對一些允許省略的部分初始化，以及檢查所需欄位是否都已定義
	err = c.init()

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to initialize.")
	}

	return c, nil
}

func (c *Config) init() error {
	if c.Mode == 0 {
		return errors.New("同步模式必須定義")
	}

	if c.FromDatabase == nil {
		if c.Mode == 2 || c.Mode == 4 {
			return errors.New("來源資料庫必須定義")
		}
	}

	if c.ProtoFolder == "" {
		if c.Mode == 1 {
			return errors.New("Proto 路徑不可為空")
		}
	}

	if c.FromTables == nil || len(c.FromTables) == 0 {
		if c.Mode == 1 {
			return errors.New("Proto 檔案不可為空")
		} else if c.Mode == 2 {
			return errors.New("來源表格名稱不可為空")
		}
	}

	if c.ToDatabase == nil {
		return errors.New("目標資料庫不可為空")
	}

	if c.ToTables == nil || len(c.ToTables) != len(c.FromTables) {
		c.ToTables = c.FromTables
	}

	return nil
}

func (c *Config) GetFromDatabase() *database.DatabaseConfig {
	return c.FromDatabase
}

func (c *Config) GetFromProtoFolder() string {
	return c.ProtoFolder
}

func (c *Config) SetFromTables(tables []string) {
	c.FromTables = tables
}

func (c *Config) GetFromTables() []string {
	return c.FromTables
}

func (c *Config) GetToDatabase() *database.DatabaseConfig {
	return c.ToDatabase
}

func (c *Config) SetToTables(tables []string) {
	c.ToTables = tables
}

func (c *Config) GetToTables() []string {
	return c.ToTables
}

func (c *Config) IsAutoGenerate() bool {
	return c.AutoGenerate
}
