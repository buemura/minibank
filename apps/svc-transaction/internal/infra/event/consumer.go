package event

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/buemura/minibank/packages/queue"
	"github.com/buemura/minibank/svc-transaction/config"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

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

	queue.DeclareQueue(ch, queue.TRANSACTION_CREATED_QUEUE)
	queue.DeclareQueue(ch, queue.TRANSACTION_CREATED_DLQ)

	msgs := make(chan amqp.Delivery)
	go queue.Consume(ch, msgs, queue.TRANSACTION_CREATED_QUEUE)

	for msg := range msgs {
		switch msg.RoutingKey {
		case queue.TRANSACTION_CREATED_QUEUE:
			var in *transaction.Transaction
			err := json.Unmarshal([]byte(msg.Body), &in)
			if err != nil {
				log.Fatalf(err.Error())
			}

			slog.Info(fmt.Sprintf("event input: %v", in))

			// notificationEvent := event.NewNotificationEvent()
			// _, err = notificationEvent.SendNotification(in)
			// if err != nil {
			// 	log.Println(err)
			// 	err = queue.PublishToQueue(ch, string(msg.Body), event.NOTIFY_ENDPOINT_DOWN_DLQ)
			// 	if err != nil {
			// 		log.Fatalf("Failed to send message to DLQ queue: %s", err)
			// 	}
			// }
		}
	}
}
