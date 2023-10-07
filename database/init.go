package database

import (
	"datamanage/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(settings *conf.Settings) {
	if settings.Dsn == "" {
		panic("必须在配置文件中配置dsn")
	}
	var options []DbOption
	if settings.MaxOpenCons > 0 {
		options = append(options, WithMaxOpenCons(settings.MaxOpenCons))
	}
	if settings.MaxIdleCons > 0 {
		options = append(options, WithMaxIdleCons(settings.MaxIdleCons))
	}
	InitDatabase(settings.Dsn, options...)
}

// InitDatabase 初始化数据库连接对象
func InitDatabase(dsn string, opts ...DbOption) {
	dbInstance, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	for _, opt := range opts {
		opt(dbInstance)
	}
	err = dbInstance.Ping()
	if err != nil {
		panic(err)
	}
	db = dbInstance
}

// MigrationTables 迁移数据库
func MigrationTables() {
	db.MustExec(createDBSchema)
	db.MustExec(selectDBSchema)
	db.MustExec(migrationWatchBinlogInfoSchema)
	db.MustExec(migrationWatchTableInfoSchema)
}

// GetSession 获取数据库连接对象
func GetSession() *sqlx.DB {
	return db
}
