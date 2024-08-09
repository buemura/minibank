package main

import (
	"log"

	"github.com/buemura/minibank/packages/queue"
	"github.com/buemura/minibank/svc-transaction/config"
	"github.com/buemura/minibank/svc-transaction/internal/infra/database"
	"github.com/buemura/minibank/svc-transaction/internal/infra/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	config.LoadEnv()
	database.Connect()
}

func main() {
	conn, ch := queue.Connect(config.BROKER_URL)
	defer func() {
		if err := ch.Close(); err != nil {
			log.Printf("Failed to close channel: %v", err)
		}
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}()

	if err := queue.DeclareQueue(ch, queue.TRANSACTION_REQUESTED_QUEUE); err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	if err := queue.DeclareQueue(ch, queue.TRANSACTION_REQUESTED_DLQ); err != nil {
		log.Fatalf("Failed to declare DLQ: %v", err)
	}

	msgs := make(chan amqp.Delivery)

	go queue.Consume(ch, msgs, queue.TRANSACTION_REQUESTED_QUEUE)

	for msg := range msgs {
		event.TransactionEventHandler(ch, msg)
	}

	select {}
}
