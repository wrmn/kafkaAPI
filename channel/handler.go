package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func postBiller(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logWriter("new Biller request")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logWriter(err.Error())
		fmt.Printf("Error while processing request : %s\n", err.Error())
	}

	var transaction transaction

	err = json.Unmarshal(b, &transaction)
	if err != nil {
		incJSON(err, w, transaction)
		return
	}

	result, topic, err := fromJSON(transaction)
	if err != nil {
		incISO(err, w, result)
		return
	}

	isoRes := []byte(result)
	w.WriteHeader(200)
	prodKafka(isoRes, topic)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "channel",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	c.SubscribeTopics([]string{"response"}, nil)

	msg, err := c.ReadMessage(-1)
	if err == nil {
		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
	} else {
		fmt.Printf("Consumer error: %v (%v)\n", err, msg)
	}

	c.Close()
	w.Write(msg.Value)
}

func prodKafka(iso []byte, topic string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	if err != nil {
		panic(err)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(iso),
	}, nil)

	p.Flush(15 * 1000)
}
