package writer

import (
	"OcrClient/config"
	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(cfg config.Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(cfg.Kafka.Brokers...),
		Topic: cfg.Kafka.Writer.Topic,
	}
}
