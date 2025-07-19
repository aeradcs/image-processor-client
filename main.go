package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// to produce messages
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092", "localhost:9093", "localhost:9094"),
		Topic:    "image-recognition-requests",
		Balancer: &kafka.LeastBytes{},
	}

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(`{"file_url":"https://672421063581-images-for-ocr.s3.us-east-1.amazonaws.com/img.png"}`),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	} else {
		log.Println("successfully wrote messages")
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Topic:   "image-recognition-responses",
		GroupID: "image-recognition-responses-group",
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		var result map[string]interface{}
		json.Unmarshal(m.Value, &result)
		fmt.Printf("Text: %s\n", result["text_detected"])
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

}
