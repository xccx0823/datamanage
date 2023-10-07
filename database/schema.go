package database

const (
	createDBSchema                 = "CREATE DATABASE data_manage;"
	selectDBSchema                 = "USE data_manage;"
	migrationWatchBinlogInfoSchema = `
CREATE TABLE watch_binlog_info (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
	binlog_file VARCHAR ( 128 ) NOT NULL COMMENT 'binlog的文件名称',
	binlog_position INT NOT NULL COMMENT '记录的位置',
	state TINYINT NOT NULL DEFAULT '0' COMMENT '0:删除 1:正常',
	create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY ( ID ) 
) COMMENT = '监听binlog信息表';
`
	migrationWatchTableInfoSchema = `
CREATE TABLE watch_table_info (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
	database_name VARCHAR ( 128 ) NOT NULL COMMENT '数据库名',
	table_name VARCHAR ( 128 ) NOT NULL COMMENT '表名',
	state TINYINT NOT NULL DEFAULT '0' COMMENT '0:删除 1:正常',
	create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY ( ID ) 
) COMMENT = '监听表信息表';
`
)
