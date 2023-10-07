package datawatch

import (
	"datamanage/database"
	"github.com/go-mysql-org/go-mysql/mysql"
)

const (
	queryPositionInfo = "SELECT * FROM data_manage.watch_binlog_info WHERE state=1"
	showBinlogStatus  = "SHOW MASTER STATUS"
)

// getPosition 以MySQL连接的host和port作为key，在数据表中存储其对应的同步位置，当没有设置其位置的时候，返回其起始值。
// 在 MySQL binlog 文件中，前 4 个字节通常包含了一个文件头，其中包括了文件的格式和版本信息，以及一些其他元数据。之后的
// 数据才包含实际的数据库变更事件。通过将 Pos 设置为 4，你将跳过这个文件头，从实际事件数据的位置开始同步。
func (sdw *SourceDataWatcher) getPosition() mysql.Position {
	position := database.WatchBinlogInfo{}
	db := database.GetSession()
	err := db.Get(&position, queryPositionInfo)
	var name string
	var pos uint32
	if err != nil {
		name, pos = getBinlogName()
	} else {
		name = position.BinlogFile
		pos = position.BinlogPosition
	}
	if name == "" {
		panic("无法获取 binlog 的文件名称")
	}
	return mysql.Position{Name: name, Pos: pos}
}

func getBinlogName() (string, uint32) {
	db := database.GetSession()
	var binlogStatus struct {
		File            string `db:"File"`
		Position        uint32 `db:"Position"`
		BinlogDoDb      string `db:"Binlog_Do_DB"`
		BinlogIgnoreDb  string `db:"Binlog_Ignore_DB"`
		ExecutedGtidSet string `db:"Executed_Gtid_Set"`
	}
	err := db.Get(&binlogStatus, showBinlogStatus)
	if err != nil {
		return "", 0
	}
	return binlogStatus.File, binlogStatus.Position
}
