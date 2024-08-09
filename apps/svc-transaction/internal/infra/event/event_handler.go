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
	case queue.TRANSACTION_REQUESTED_QUEUE:
		// Parse message body
		var in *transaction.Transaction
		err := json.Unmarshal([]byte(msg.Body), &in)
		if err != nil {
			log.Fatalf(err.Error())
		}

		switch in.TransactionType {
		case transaction.Transfer:
			transferUC := factory.MakePerformTransferUsecase()
			err = transferUC.Execute(in)
		case transaction.Deposit:
			depositUC := factory.MakePerformDepositUsecase()
			err = depositUC.Execute(in)
		case transaction.Withdrawal:
			withdrawUC := factory.MakePerformWithdrawUsecase()
			err = withdrawUC.Execute(in)
		}

		if err != nil {
			slog.Error(err.Error())

			// TODO: adds retry stategy before sending it to DLQ
			err = queue.PublishToQueue(ch, msg.Body, queue.TRANSACTION_REQUESTED_DLQ)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to send message to DLQ queue: %s", err))
			}
		}
	}
}
