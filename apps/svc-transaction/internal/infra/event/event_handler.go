package event

import (
	"encoding/json"
	"log"

	"github.com/buemura/minibank/packages/queue"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TransactionEventHandler(msg amqp.Delivery) {
	switch msg.RoutingKey {
	case queue.TRANSFER_REQUESTED_QUEUE:
		var in *transaction.Transaction
		err := json.Unmarshal([]byte(msg.Body), &in)
		if err != nil {
			log.Fatalf(err.Error())
		}

		// TODO: Validate if origin account has enough balance
		// TODO: Validate if destination account is valid
		// TODO: Perform transfer
		// TODO: Update account balances

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
