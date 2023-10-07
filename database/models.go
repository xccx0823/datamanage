package database

type WatchBinlogInfo struct {
	ID             string `db:"id"`
	BinlogFile     string `db:"binlog_file"`
	BinlogPosition uint32 `db:"binlog_position"`
	State          string `db:"state"`
	CreateTime     string `db:"create_time"`
	UpdateTime     string `db:"update_time"`
}

type WatchTableInfo struct {
	ID           string `db:"id"`
	DatabaseName string `db:"database_name"`
	TableName    string `db:"table_name"`
	State        string `db:"state"`
	CreateTime   string `db:"create_time"`
	UpdateTime   string `db:"update_time"`
}
