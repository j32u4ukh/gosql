package main

import (
	"fmt"

	"github.com/j32u4ukh/gosql/MigrateGo/sync"
	"github.com/j32u4ukh/gosql/database"
)

func main() {
	cf, err := database.NewConfig("../config/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load configuration, err: %v", err)
		return
	}
	fmt.Printf("config: %+v\n", cf)

	scf := sync.NewConfig()
	err = scf.LoadFile("../../MigrateGo/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load configuration, err: %v", err)
		return
	}
	fmt.Printf("sync config: %+v\n", scf)
}
