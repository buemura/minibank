package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Connect(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return conn, ch
}

func CreateChannel(url string) *amqp.Channel {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

func DeclareQueue(ch *amqp.Channel, queue string) error {
	_, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

func Consume(ch *amqp.Channel, out chan<- amqp.Delivery, queue string) error {
	consumerTag := fmt.Sprintf("go-consumer-%s-%d", queue, time.Now().UnixNano())

	msgs, err := ch.Consume(
		queue,
		consumerTag,
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

func PublishToQueue(ch *amqp.Channel, body interface{}, queue string) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	log.Printf("[PublishToQueue] - Sending messagem to queue: %s", queue)
	err = ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        b,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
