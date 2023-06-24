package main

import (
	"fmt"
	"os"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gosql/MigrateGo/sync"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/utils"
	"github.com/pkg/errors"
)

var logger *glog.Logger
var level glog.LogLevel = glog.DebugLevel
var synConfig *sync.Config
var length int

func init() {
	synConfig = sync.NewConfig()
	logger = glog.SetLogger(0, "MigrateGo", level)
	utils.SetLogger(logger)
}

func main() {
	length = len(os.Args)
	err := loadConfig(os.Args)
	if err != nil {
		logger.Error("Failed to load config by args: %+v.", os.Args)
		return
	}
	err = loadParams(os.Args)
	if err != nil {
		logger.Error("Failed to load params by args: %+v.", os.Args)
		return
	}
	err = synConfig.CheckParams()
	if err != nil {
		logger.Error("參數設置不完整, err: %+v.", synConfig)
		return
	}
	fmt.Printf("synConfig: %+v\n", synConfig)
	// TODO: 根據 synConfig 執行同步機制

}

func loadConfig(args []string) error {
	for i := 0; i < length; i++ {
		if args[i] == "--config" {
			if i+1 < length {
				err := synConfig.LoadFile(args[i+1])
				if err != nil {
					return errors.Wrapf(err, "Failed to load from %s.", args[i+1])
				}
			}
		}
	}
	return nil
}

func loadParams(args []string) error {
	var err error
	for i := 0; i < length; i++ {
		switch args[i] {
		case "--mode":
			err = managedExecute(i, func(idx int) error {
				switch args[idx] {
				case "1":
					synConfig.Mode = sync.ProtoToDbMode
				case "2":
					synConfig.Mode = sync.DbToDbMode
				default:
					return errors.Errorf("Invalid parameter %s", args[idx])
				}
				return nil
			})
			if err != nil {
				return errors.Wrap(err, "Failed to load --mode parameter")
			}
		case "--folder":
			managedExecute(i, func(idx int) error {
				synConfig.ProtoFolder = args[idx]
				return nil
			})
		case "--from":
			managedExecute(i, func(idx int) error {
				synConfig.FromDatabase, err = database.ParseDatabaseConfig(args[idx])
				if err != nil {
					return errors.Wrap(err, "Failed to parse from-database configuration.")
				}
				return nil
			})
		case "--from-table":
			managedExecute(i, func(idx int) error {
				synConfig.FromTable = args[idx]
				return nil
			})
		case "--to":
			managedExecute(i, func(idx int) error {
				synConfig.ToDatabase, err = database.ParseDatabaseConfig(args[idx])
				if err != nil {
					return errors.Wrap(err, "Failed to parse from-database configuration.")
				}
				return nil
			})
		case "--to-table":
			managedExecute(i, func(idx int) error {
				synConfig.ToTable = args[idx]
				return nil
			})
		case "--print":
			synConfig.Print = true
		case "--sync":
			synConfig.Sync = true
		case "--generate":
			synConfig.Generate = true
		}
	}
	return nil
}

func managedExecute(i int, execute func(idx int) error) error {
	if i+1 < length {
		return execute(i + 1)
	}
	return nil
}
