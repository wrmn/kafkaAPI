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

	var transaction Transaction

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
	w.Write(isoRes)
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
