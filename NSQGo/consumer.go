package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

// HandleMessage implements the Handler interface.
func HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
		return errors.New("Empty string")
	}
	log.Printf("Got a new message: %v\n", string(m.Body))

	// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
	return nil
}

func main() {
	// Instantiate a consumer that will subscribe to the provided channel.
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("topic", "channel", config)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	consumer.AddHandler(nsq.HandlerFunc(HandleMessage))

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = consumer.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	// Prevents program from terminating
	// Use ctrl+c to stop the program
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	// Gracefully stop the consumer.
	consumer.Stop()
}
