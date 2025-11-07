package message

import (
	"context"
	"encoding/json"
	"log/slog"
	ticketsEntity "tickets/entities"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)


type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request ticketsEntity.IssueReceiptRequest) error
}

func NewMessageRouter(
	rdb redis.UniversalClient, 
	logger watermill.LoggerAdapter,
	spreadSheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
) *message.Router {
	ctx := context.Background()
	router := message.NewDefaultRouter(logger)


	spreadSheetSub, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client: rdb,
			ConsumerGroup: AppendToTrackerTopic,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	router.AddConsumerHandler(
		"append-to-tracker",
		AppendToTrackerTopic,
		spreadSheetSub, 
		func(msg *message.Message) error {
			var data ticketsEntity.AppendToTrackerPayload
			err := json.Unmarshal(msg.Payload, &data)
			if err != nil {
				return err
			}
			slog.Info("Appending ticket to the tracker")
			err = spreadSheetsAPI.AppendRow(
				ctx, 
				"tickets-to-print",
				[]string{data.TicketID, data.CustomerEmail, data.Price.Amount, data.Price.Currency},
			)
			
			if err != nil {
				return err
			}
			return nil
		},
	)


	receiptServiceSub, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client: rdb,
			ConsumerGroup: IssueReceiptTopic,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	router.AddConsumerHandler(
		"issue-receipt",
		IssueReceiptTopic,
		receiptServiceSub, 
		func(msg *message.Message) error {
			var data ticketsEntity.AppendToTrackerPayload
			err := json.Unmarshal(msg.Payload, &data)
			if err != nil {
				return err
			}
			event := ticketsEntity.IssueReceiptRequest{
				TicketID: data.TicketID,
				Price: ticketsEntity.Money{
					Amount: data.Price.Amount,
					Currency: data.Price.Currency,
				},
			}
			err = receiptsService.IssueReceipt(ctx, event)
			if err != nil {
				return err
			}
			return nil
		},
	)

	return router
}
