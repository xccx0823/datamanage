package datawatch

import (
	"fmt"
	"time"
)
import "github.com/IBM/sarama"

// 事件通用信息
type queueData struct {
	Database  string `json:"database,omitempty"`
	TableName string `json:"table_name,omitempty"`
	Sql       string `json:"sql,omitempty"`
}

func sendToQueue(data queueData) {
	fmt.Println(data.Sql)
}

func (sdw *SourceDataWatcher) InitQueue() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = time.Duration(sdw.KafkaFlushFrequency) * time.Millisecond
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(sdw.KafkaAddress, config)
	if err != nil {
		panic(err)
	}
	sdw.KafkaProducer = producer
}
