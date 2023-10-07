package datawatch

import (
	"datamanage/log"
	"encoding/json"
	"time"
)
import "github.com/IBM/sarama"

// 事件通用信息
type queueData struct {
	Database  string `json:"database,omitempty"`
	TableName string `json:"table_name,omitempty"`
	Sql       string `json:"sql,omitempty"`
}

func (sdw *SourceDataWatcher) sendToQueue(data queueData) {
	marshal, err := json.Marshal(data)
	if err != nil {
		return
	}
	topic := data.Database
	producerMessage := &sarama.ProducerMessage{
		Topic: data.Database,
		Value: sarama.StringEncoder(marshal),
	}
	partition, offset, err := sdw.KafkaProducer.SendMessage(producerMessage)
	if err != nil {
		log.Error(err)
	} else {
		log.InfoF("sent topic:%s partition:%d offset:%d", topic, partition, offset)
	}
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
