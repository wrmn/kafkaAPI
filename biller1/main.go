package main

import (
	"fmt"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "consumer",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"quickstart-events", "^aRegex.*[Tt]opic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			toJson(string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
