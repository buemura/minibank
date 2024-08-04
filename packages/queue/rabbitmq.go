package queue

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func DeclareQueue(ch *amqp.Channel, queue string) {
	_, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
}

func Consume(ch *amqp.Channel, out chan<- amqp.Delivery, queue string) error {
	msgs, err := ch.Consume(
		queue,
		"go-consumer",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Printf("⇨ Consuming Queue: %s\n", queue)
	for msg := range msgs {
		out <- msg
	}
	return nil
}

func Publish(ch *amqp.Channel, body string, exName string) error {
	log.Printf("Sending messagem to exchange: %s", exName)

	err := ch.Publish(
		exName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func PublishToQueue(ch *amqp.Channel, body string, queue string) error {
	log.Printf("Sending messagem to queue: %s", queue)

	err := ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
