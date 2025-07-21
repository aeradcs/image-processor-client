package main

import (
	"OcrClient/config"
	"OcrClient/pkg/reader"
	"context"
	"encoding/json"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config/kafka_config.yml")
	if err != nil {
		log.Fatal("Failed to load config from config/kafka_config.yml: ", err)
	}

	r := reader.NewKafkaReader(*cfg)
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Failed to read message: ", err)
		}
		var result map[string]interface{}
		err = json.Unmarshal(m.Value, &result)
		if err != nil {
			log.Fatal("Failed to unmarshall response: ", err)
		}
		log.Printf("Received and unmrshalled message: \n%v", result)
	}
}
