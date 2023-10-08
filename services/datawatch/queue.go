package datawatch

import (
	"datamanage/log"
	"encoding/json"
	"github.com/IBM/sarama"
)

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
	producerMessage := &sarama.ProducerMessage{
		Topic: sdw.KafkaTopic,
		Value: sarama.StringEncoder(marshal),
	}
	partition, offset, err := sdw.KafkaProducer.SendMessage(producerMessage)
	if err != nil {
		log.Error(err)
	} else {
		log.InfoF("sent topic:%s partition:%d offset:%d", sdw.KafkaTopic, partition, offset)
	}
}

func (sdw *SourceDataWatcher) InitQueue() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer, err := sarama.NewSyncProducer(sdw.KafkaAddress, config)
	if err != nil {
		panic(err)
	}
	sdw.KafkaProducer = producer
}
