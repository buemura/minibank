package event

import (
	"log"
	"sync"

	"github.com/buemura/minibank/packages/queue"
	"github.com/buemura/minibank/svc-transaction/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var queueList = []string{
	queue.TRANSFER_REQUESTED_QUEUE,
	queue.DEPOSIT_REQUESTED_QUEUE,
	queue.WITHDRAW_REQUESTED_QUEUE,
}

var dlqList = []string{
	queue.TRANSFER_REQUESTED_DLQ,
	queue.DEPOSIT_REQUESTED_DLQ,
	queue.WITHDRAW_REQUESTED_DLQ,
}

/**
 * @title: StartConsumer
 * @description: This function is responsible for consuming messages from the RabbitMQ queue. It establishes a connection to the broker, declares the necessary queues, and starts consuming messages from the specified queue.
 * @param none
 * @return none
 */
func StartConsumer() {
	conn, ch := queue.Connect(config.BROKER_URL)
	defer func() {
		if err := ch.Close(); err != nil {
			log.Printf("Failed to close channel: %v", err)
		}
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}()

	for _, q := range queueList {
		if err := queue.DeclareQueue(ch, q); err != nil {
			log.Fatalf("Failed to declare queue: %v", err)
		}
	}
	for _, dlq := range dlqList {
		if err := queue.DeclareQueue(ch, dlq); err != nil {
			log.Fatalf("Failed to declare DLQ: %v", err)
		}
	}

	msgs := make(chan amqp.Delivery)
	var wg sync.WaitGroup
	wg.Add(len(queueList))

	for _, q := range queueList {
		q := q
		go func(queueName string) {
			defer wg.Done()
			if err := queue.Consume(ch, msgs, queueName); err != nil {
				log.Printf("Failed to consume from queue %s: %v", queueName, err)
			}
		}(q)
	}

	go func() {
		for msg := range msgs {
			TransactionEventHandler(ch, msg)
		}
	}()

	wg.Wait()

	select {}
}
