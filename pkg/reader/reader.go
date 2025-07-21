package reader

import (
	"OcrClient/config"
	"github.com/segmentio/kafka-go"
)

func NewKafkaReader(cfg config.Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Reader.Topic,
		GroupID: cfg.Kafka.Reader.GroupID,
	})
}
