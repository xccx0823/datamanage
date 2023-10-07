package datawatch

import (
	"datamanage/log"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-mysql-org/go-mysql/replication"
	"regexp"
	"strings"
)

func (sdw *SourceDataWatcher) OnDDL(e *replication.QueryEvent) error {
	databaseName := string(e.Schema)
	sql := string(e.Query)

	// alter语法修改表结构
	queryUpper := strings.ToUpper(sql)
	regex := regexp.MustCompile(`ALTER\s+TABLE\s+(\w+)`)
	matches := regex.FindStringSubmatch(queryUpper)
	if len(matches) > 1 {
		tableName := matches[1]
		tables := sdw.monitorColumns[databaseName]
		delete(tables, tableName)
		sdw.sendToQueue(queueData{Database: databaseName, TableName: tableName, Sql: sql})
	}

	// drop table 语法删除的表
	regex = regexp.MustCompile(`DROP\s+TABLE\s+IF\s+EXISTS\s+([\w_]+);`)
	matches = regex.FindStringSubmatch(queryUpper)
	if len(matches) > 1 {
		tableName := matches[1]
		tables := sdw.monitorColumns[databaseName]
		delete(tables, tableName)
		sdw.sendToQueue(queueData{Database: databaseName, TableName: tableName, Sql: sql})
	}

	return nil
}

func (sdw *SourceDataWatcher) OnRotate(e *replication.RotateEvent) error {
	log.InfoF("Received Rotate Event: Next Log Name: %s", e.NextLogName)
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

	var sql string

	switch eType {

	// 更新事件
	case replication.UPDATE_ROWS_EVENTv0, replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		sql = updateSql(databaseName, tableName, e.Rows, columns)

	// 插入数据事件
	case replication.WRITE_ROWS_EVENTv0, replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		sql = insertSql(databaseName, tableName, e.Rows, columns)

	// 删除数据事件
	case replication.DELETE_ROWS_EVENTv0, replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		sql = deleteSql(databaseName, tableName, e.Rows, columns)
	}

	sdw.sendToQueue(queueData{Database: databaseName, TableName: tableName, Sql: sql})

	return nil
}

func updateSql(databaseName, tableName string, data [][]any, columns []string) string {
	beforeData := data[0]
	afterData := data[1]
	var setData []string
	var whereData []string
	for idx, row := range beforeData {
		rowStr := convertor.ToString(row)
		afterRow := convertor.ToString(afterData[idx])
		column := columns[idx]
		if rowStr != afterRow {
			setData = append(setData, column+"="+convertor.ToString(afterRow))
		}
		whereData = append(whereData, column+"="+convertor.ToString(row))
	}
	setSql := strings.Join(setData, ", ")
	whereSql := strings.Join(whereData, " and ")
	sql := fmt.Sprintf("update %s.%s set %s where %s", databaseName, tableName, setSql, whereSql)
	return sql
}

func deleteSql(databaseName, tableName string, data [][]any, columns []string) string {
	delData := data[0]
	var whereData []string
	for idx, row := range delData {
		column := columns[idx]
		whereData = append(whereData, column+"="+convertor.ToString(row))
	}
	whereSql := strings.Join(whereData, " and ")
	sql := fmt.Sprintf("delete from %s.%s where %s", databaseName, tableName, whereSql)
	return sql
}

func insertSql(databaseName, tableName string, data [][]any, columns []string) string {
	insertData := data[0]
	insertColumnSql := strings.Join(columns, ", ")
	var insertStrData []string
	for _, row := range insertData {
		insertStrData = append(insertStrData, convertor.ToString(row))
	}
	insertDataSql := strings.Join(insertStrData, ", ")
	sql := fmt.Sprintf("insert into %s.%s (%s) values (%s)", databaseName, tableName, insertColumnSql, insertDataSql)
	return sql
}

func (sdw *SourceDataWatcher) getColumns(databaseName, tableName string) []string {
	tables, ok := sdw.monitorColumns[databaseName]
	if !ok {
		columns := getTableColumns(databaseName, tableName)
		tableColumns := make(map[string][]string)
		tableColumns[tableName] = columns
		sdw.monitorColumns[databaseName] = tableColumns
		return columns
	}
	columns, ok := tables[tableName]
	if !ok {
		columns := getTableColumns(databaseName, tableName)
		subColumns := make(map[string][]string)
		subColumns[tableName] = columns
		sdw.monitorColumns[databaseName] = subColumns
		return columns
	}
	return columns
}
