package main

import (
	"datamanage/conf"
	"datamanage/database"
	"datamanage/log"
	"datamanage/services/datawatch"
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		echoHelp()
		os.Exit(1)
		return
	}

	switch os.Args[1] {
	case "-h":
		echoHelp()

	// binlog监听服务
	case "datawatch":
		datawatchSet := flag.NewFlagSet("datawatch", flag.ExitOnError)
		datawatchConfig := datawatchSet.String("c", "", "配置文件")
		err := datawatchSet.Parse(os.Args[2:])
		if err != nil {
			log.Error(err)
			return
		}
		settings := initCores(*datawatchConfig)
		watcher := datawatch.New(settings)
		watcher.Run()

	// 数据库迁移
	case "migration":
		migrationSet := flag.NewFlagSet("migration", flag.ExitOnError)
		migrationConfig := migrationSet.String("c", "", "配置文件")
		err := migrationSet.Parse(os.Args[2:])
		if err != nil {
			log.Error(err)
			return
		}
		initCores(*migrationConfig)
		database.MigrationTables()
	}
}

func initCores(configPath string) *conf.Settings {
	log.Init()
	settings := conf.InitConfig(configPath)
	database.Init(settings)
	return settings
}

func echoHelp() {
	fmt.Printf(`Usage:
  datawatch 启动MySQL的binlog监听服务
  migration 项目所需数据表结构迁移

  -h 查看帮助信息
`)
}
