package message

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}


func NewHandlers(
	rdb redis.UniversalClient, 
	logger watermill.LoggerAdapter,
	spreadSheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
) {
	ctx := context.Background()
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

	go func(){
		messages, err := spreadSheetSub.Subscribe(ctx, AppendToTrackerTopic)
		if err != nil {
			panic(err)
		}
		for msg := range messages {
			err := spreadSheetsAPI.AppendRow(ctx, "tickets-to-print", []string{string(msg.Payload)})
			if err != nil {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()


	go func(){
		messages, err := receiptServiceSub.Subscribe(ctx, AppendToTrackerTopic)
		if err != nil {
			panic(err)
		}
		for msg := range messages {
			err := receiptsService.IssueReceipt(ctx, string(msg.Payload))
			if err != nil {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()


}
