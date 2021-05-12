package main

import (
	"log"
	"sync"

	"github.com/nsqio/go-nsq"
)

// HandleMessage implements the Handler interface.
func HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		return nil
	}

	message := string(m.Body)
	log.Printf("Got a new message: %v\n", message)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Instantiate a consumer that will subscribe to the provided channel.
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("topic", "channel", config)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %v", string(message.Body))
		wg.Done()
		return nil
	}))

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = consumer.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	// Gracefully stop the consumer.
	consumer.Stop()
}
