package datawatch

// 事件通用信息
type queueData struct {
	Database  string `json:"database,omitempty"`
	TableName string `json:"table_name,omitempty"`
	Sql       string `json:"sql,omitempty"`
}

func sendToQueue(data queueData) {

}
