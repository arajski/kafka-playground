package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/riferrei/srclient"
)

type MessageType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	topic := "test"
	schemaRegistryUrl := "http://localhost:8081"
	kafkaHostname := "localhost:19092"
	message := MessageType{ID: 1, Name: "TestMessage"}

	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaRegistryUrl)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaHostname})
	if err != nil {
		fmt.Printf("couldn't connect to kafka boker: %+v\n", err)
		os.Exit(1)
	}

	defer p.Close()

	// Delivery report handler for produced messages
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

	schema, err := schemaRegistryClient.GetLatestSchema(topic)
	if schema == nil || err != nil {
		fmt.Println("couldn't get schema")
		os.Exit(1)
	}

	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	value, _ := json.Marshal(message)
	native, _, _ := schema.Codec().NativeFromTextual(value)
	valueBytes, _ := schema.Codec().BinaryFromNative(nil, native)

	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, valueBytes...)

	// Produce message to topic
	key, _ := uuid.NewUUID()
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          recordValue,
		Key:            []byte(key.String()),
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
