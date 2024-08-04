package event

import (
	"github.com/buemura/minibank/packages/queue"
	"github.com/buemura/minibank/svc-transaction/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

/**
 * @title: QueueConsumer
 * @description: This function is responsible for consuming messages from the RabbitMQ queue. It establishes a connection to the broker, declares the necessary queues, and starts consuming messages from the specified queue.
 * @param none
 * @return none
 */
func QueueConsumer() {
	conn, ch := queue.Connect(config.BROKER_URL)
	defer conn.Close()
	defer ch.Close()

	queue.DeclareQueue(ch, queue.TRANSFER_REQUESTED_QUEUE)
	queue.DeclareQueue(ch, queue.TRANSFER_REQUESTED_DLQ)

	msgs := make(chan amqp.Delivery)
	go queue.Consume(ch, msgs, queue.TRANSFER_REQUESTED_QUEUE)

	for msg := range msgs {
		TransactionEventHandler(ch, msg)
	}
}
