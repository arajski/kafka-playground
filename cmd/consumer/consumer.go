package main

import (
	"encoding/binary"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/riferrei/srclient"
	"os"
)

func main() {
	topic := "test"
	schemaRegistryHostname := "http://localhost:8081"
	kafkaHostname := "localhost:19092"
	groupId := "myGroup"

	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaRegistryHostname)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaHostname,
		"group.id":          groupId,
	})

	if err != nil {
		fmt.Printf("could not connect to broker: %+v\n", err)
		os.Exit(1)
	}

	fmt.Println("successfully connected to the broker!")

	c.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			fmt.Printf("consumer error: %v (%v)\n", err, msg)
		} else {
			schemaId := binary.BigEndian.Uint32(msg.Value[1:5])
			schema, err := schemaRegistryClient.GetSchema(int(schemaId))
			if err != nil {
				fmt.Printf("error getting the schema with id '%d' %s", schemaId, err)
				os.Exit(1)
			}

			native, _, _ := schema.Codec().NativeFromBinary(msg.Value[5:])
			value, _ := schema.Codec().TextualFromNative(nil, native)
			fmt.Printf("here is the message %s\n", string(value))
		}
	}
}
