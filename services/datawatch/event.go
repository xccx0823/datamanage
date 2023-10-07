package datawatch

import (
	"datamanage/log"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-mysql-org/go-mysql/replication"
)

// 事件通用信息
type eventInfo struct {
	Database  string `json:"database,omitempty"`
	TableName string `json:"table_name,omitempty"`
	Sql       string `json:"sql,omitempty"`
}

func (sdw *SourceDataWatcher) OnDDL(e *replication.QueryEvent) error {
	// 创建数据库
	// 创建表
	// 修改以及创建数据时触发，会提供BEGIN消息
	log.InfoF("Received DDL Event: Query: %s", e.Query)
	return nil
}

func (sdw *SourceDataWatcher) OnRotate(e *replication.RotateEvent) error {
	log.InfoF("Received Rotate Event: Next Log Name: %s", e.NextLogName)
	return nil
}

func (sdw *SourceDataWatcher) OnTableChanged(e *replication.TableMapEvent) error {
	// 修改以及创建数据时触发，会提供表信息
	databaseName := convertor.ToString(e.Schema)
	tableName := convertor.ToString(e.Table)
	tables, ok := sdw.monitorTables[databaseName]
	if !ok || !slice.Contain(tables, tableName) {
		return nil
	}
	log.InfoF("Received Table Changed Event: %s", e.Table)
	return nil
}

func (sdw *SourceDataWatcher) OnRow(e *replication.RowsEvent, eType replication.EventType) error {

	databaseName := convertor.ToString(e.Table.Schema)
	tableName := convertor.ToString(e.Table.Table)

	fmt.Println(sdw.monitorTables)

	tables, ok := sdw.monitorTables[databaseName]
	if !ok || !slice.Contain(tables, tableName) {
		return nil
	}

	switch eType {

	// 更新事件
	case replication.UPDATE_ROWS_EVENTv0, replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		log.InfoF("更新 %s.%s %s", databaseName, tableName, e.Rows)

	// 插入数据事件
	case replication.WRITE_ROWS_EVENTv0, replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		log.InfoF("插入 %s.%s %s", databaseName, tableName, e.Rows)

	// 删除数据事件
	case replication.DELETE_ROWS_EVENTv0, replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		log.InfoF("删除 %s.%s %s", databaseName, tableName, e.Rows)
	}

	return nil
}

func updateSql(data [][]any) {

}

func deleteSql(data [][]any) {

}

func insertSql(data [][]any) {

}
