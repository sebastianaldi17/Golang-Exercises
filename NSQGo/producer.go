package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/nsqio/go-nsq"
)

func main() {
	// Instantiate a producer.
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	// Prompt for name in standard input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your message?")
	message, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	messageBody := []byte(message)
	topicName := "topic"

	// Synchronously publish a single message to the specified topic.
	// Messages can also be sent asynchronously and/or in batches.
	err = producer.Publish(topicName, messageBody)
	if err != nil {
		log.Fatal(err)
	}

	// Gracefully stop the producer.
	producer.Stop()
}
