package database

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database DatabaseConfig `yaml:"Database"`
}

type DatabaseConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	DbName   string `yaml:"DbName"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
}

var configs map[string]*Config

func init() {
	configs = make(map[string]*Config)
}

func NewConfig(path string) (*Config, error) {
	var c *Config
	var ok bool
	if c, ok = configs[path]; ok {
		return c, nil
	}
	var b []byte
	var err error
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read file %s.", path)
	}
	c = &Config{}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, errors.Wrapf(err, "讀取 Config 時發生錯誤(path: %s)", path)
	}
	configs[path] = c
	return c, nil
}

func (c *Config) GetDatabase() DatabaseConfig {
	return c.Database
}

func NewDatabaseConfig() *DatabaseConfig {
	dc := &DatabaseConfig{
		Host:     "",
		Port:     -1,
		DbName:   "",
		User:     "",
		Password: "",
	}
	return dc
}

// args ex: host=127.0.0.1;port=3306;dbname=pekomiko;user=user;password=password
func ParseDatabaseConfig(args string) (*DatabaseConfig, error) {
	params := strings.Split(args, ";")
	nParam := len(params)
	if nParam != 5 {
		return nil, fmt.Errorf("參數不足(#args: %d)", nParam)
	}
	dc := NewDatabaseConfig()
	var param, key string
	var kv []string
	var err error
	for _, param = range params {
		kv = strings.Split(param, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("參數格式錯誤: %s", param)
		}
		key = strings.ToUpper(kv[0])
		switch key {
		case "HOST":
			dc.Host = kv[1]
		case "PORT":
			dc.Port, err = strconv.Atoi(kv[1])
			if err != nil {
				return nil, fmt.Errorf("port 參數錯誤: %s", kv[1])
			}
		case "DBNAME":
			dc.DbName = kv[1]
		case "USER":
			dc.User = kv[1]
		case "PASSWORD":
			dc.Password = kv[1]
		default:
			return nil, errors.Errorf("Unknown parameter: %s", param)
		}
	}
	if dc.Host == "" || dc.Port == -1 || dc.DbName == "" || dc.User == "" || dc.Password == "" {
		return nil, errors.Errorf("缺少部分參數, config: %+v", dc)
	}
	return dc, nil
}
