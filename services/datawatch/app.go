package datawatch

import (
	"datamanage/conf"
	"datamanage/database"
	"datamanage/log"
	"sync"
	"time"
)

// SourceDataWatcher 数据源数据监听
// 数据源必须是MySQL的数据库，监听服务通过监听MySQL的binlog日志，将变化同步到Kafka中服务中，消费端消费Kafka中的消息，以此达到
// 实时同步数据的需求。
type SourceDataWatcher struct {
	// ServerID 集群的唯一ID
	ServerID uint32

	// MySQLl连接配置
	Host     string
	Port     uint16
	User     string
	Password string
	Charset  string

	// monitorTables 所有的需要监听的数据表列表
	// monitorSyncTime 同步一次的事件，单位秒
	monitorTables   map[string][]string
	monitorSyncTime int
	monitorLock     sync.Mutex
}

// New a SourceDataWatcher
// Example:
//
//	sdw := datawatch.New()
//	sdw.Run()
func New(settings *conf.Settings) *SourceDataWatcher {
	watcher := &SourceDataWatcher{}
	dbConf := settings.WatchServer.DB
	options := []Option{
		WithDB(dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password),
	}
	if dbConf.Charset != "" {
		options = append(options, WithCharset(dbConf.Charset))
	}
	syncTime := dbConf.WithMonitorSyncTime
	if syncTime == 0 {
		syncTime = 5
	}
	options = append(options, WithMonitorSyncTime(syncTime))
	watcher.SetOptions(dbConf.ServerId, options...)
	return watcher
}

func (sdw *SourceDataWatcher) Run() {
	go sdw.SyncTables()
	sdw.watchBinlog()
}

// 查询 watch_table_info 表中所有可用信息
const queryWatchTablesSchema = "SELECT * FROM data_manage.watch_table_info WHERE state=1"

// SyncTables 同步数据库中所有需要监听的数据表
func (sdw *SourceDataWatcher) SyncTables() {
	for {
		var infos []database.WatchTableInfo
		db := database.GetSession()
		err := db.Select(&infos, queryWatchTablesSchema)
		if err != nil {
			log.Error(err)
		}
		var databases = make(map[string][]string)
		for _, info := range infos {
			tables, ok := databases[info.DatabaseName]
			if !ok {
				databases[info.DatabaseName] = make([]string, 0)
			}
			databases[info.DatabaseName] = append(tables, info.TableName)
		}
		sdw.monitorLock.Lock()
		sdw.monitorTables = databases
		sdw.monitorLock.Unlock()
		time.Sleep(time.Duration(sdw.monitorSyncTime) * time.Second)
	}
}
