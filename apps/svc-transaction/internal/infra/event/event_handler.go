package event

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/buemura/minibank/packages/queue"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
	"github.com/buemura/minibank/svc-transaction/internal/infra/factory"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TransactionEventHandler(ch *amqp.Channel, msg amqp.Delivery) {
	switch msg.RoutingKey {
	case queue.TRANSFER_REQUESTED_QUEUE:
		// Parse message body
		var in *transaction.Transaction
		err := json.Unmarshal([]byte(msg.Body), &in)
		if err != nil {
			log.Fatalf(err.Error())
		}

		// Perform transfer request
		performTransferUC := factory.MakePerformTransferUsecase()
		err = performTransferUC.Execute(in)
		if err != nil {
			slog.Error(err.Error())

			// TODO: adds retry stategy before sending it to DLQ
			err = queue.PublishToQueue(ch, msg.Body, queue.TRANSFER_REQUESTED_QUEUE)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to send message to DLQ queue: %s", err))
			}
		}

	}
}
