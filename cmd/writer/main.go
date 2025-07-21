package main

import (
	"OcrClient/config"
	"OcrClient/pkg/writer"
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config/kafka_config.yml")
	if err != nil {
		log.Fatal("Failed to load config from config/kafka_config.yml: ", err)
	}

	w := writer.NewKafkaWriter(*cfg)
	defer w.Close()

	msg := kafka.Message{
		Value: []byte(`{"file_url":"https://672421063581-images-for-ocr.s3.us-east-1.amazonaws.com/img.png"}`),
	}
	err = w.WriteMessages(context.Background(),
		msg,
	)
	if err != nil {
		log.Fatal("Failed to send messages: ", err)
	} else {
		log.Printf("Successfully sent messages into topic: %s, key: %s, value: %s", w.Topic, msg.Key, msg.Value)
	}

}
