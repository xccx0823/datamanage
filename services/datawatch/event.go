package datawatch

import (
	"datamanage/log"
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
	tables, ok := sdw.monitorTables[databaseName]
	if !ok || !slice.Contain(tables, tableName) {
		return nil
	}

	columns := sdw.getColumns(databaseName, tableName)

	switch eType {

	// 更新事件
	case replication.UPDATE_ROWS_EVENTv0, replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		updateSql(e.Rows, columns)

	// 插入数据事件
	case replication.WRITE_ROWS_EVENTv0, replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		insertSql(e.Rows, columns)

	// 删除数据事件
	case replication.DELETE_ROWS_EVENTv0, replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		deleteSql(e.Rows, columns)
	}

	return nil
}

func updateSql(data [][]any, columns []string) {

}

func deleteSql(data [][]any, columns []string) {

}

func insertSql(data [][]any, columns []string) {

}

func (sdw *SourceDataWatcher) getColumns(databaseName, tableName string) []string {
	tables, ok := sdw.monitorColumns[databaseName]
	if !ok {
		columns := getTableColumns(databaseName, tableName)
		sdw.monitorColumns[databaseName][tableName] = columns
		return columns
	}
	columns, ok := tables[tableName]
	if !ok {
		columns := getTableColumns(databaseName, tableName)
		sdw.monitorColumns[databaseName][tableName] = columns
		return columns
	}
	return columns
}
