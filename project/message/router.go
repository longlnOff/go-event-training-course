package message

import (
	"encoding/json"
	"log/slog"
	ticketsEntity "tickets/entities"
	ticketsEvent "tickets/message/event"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)






func NewMessageRouter(
	rdb redis.UniversalClient, 
	logger watermill.LoggerAdapter,
	spreadSheetsAPI ticketsEvent.SpreadsheetsAPI,
	receiptsService ticketsEvent.ReceiptsService,
) *message.Router {
	router := message.NewDefaultRouter(logger)

	handler := ticketsEvent.NewHandler(
		spreadSheetsAPI,
		receiptsService,
	)

	// ------------ ISSUE RECEIPT -------------
	receiptServiceSub, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client: rdb,
			ConsumerGroup: "issue-receipt",
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	router.AddConsumerHandler(
		"issue-receipt",
		TicketBookingConfirmedTopic,
		receiptServiceSub, 
		func(msg *message.Message) error {
			var event ticketsEntity.TicketBookingConfirmed
			err := json.Unmarshal(msg.Payload, &event)
			if err != nil {
				slog.Error("failed to unmarshal data")
				return err
			}
			return handler.IssueReceipt(msg.Context(), event)
		},
	)

	// ------------ APPEND TO TRACKER PRINT -------------
	spreadSheetSub, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client: rdb,
			ConsumerGroup: "append-to-tracker",
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	router.AddConsumerHandler(
		"append-to-tracker",
		TicketBookingConfirmedTopic,
		spreadSheetSub, 
		func(msg *message.Message) error {
			var event ticketsEntity.TicketBookingConfirmed
			err := json.Unmarshal(msg.Payload, &event)
			if err != nil {
				slog.Error("failed to unmarshal data")
				return err
			}
			return handler.AppendToPrint(msg.Context(), event)
		},
	)


	// ------------ APPEND TO REFUND -------------
	RefundSub, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client: rdb,
			ConsumerGroup: "append-to-refund",
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	router.AddConsumerHandler(
		"append-to-refund",
		TicketBookingCanceledTopic,
		RefundSub, 
		func(msg *message.Message) error {
			var event ticketsEntity.TicketBookingCanceled
			err := json.Unmarshal(msg.Payload, &event)
			if err != nil {
				slog.Error("failed to unmarshal data")
				return err
			}
			return handler.AppendToCancel(msg.Context(), event)
		},
	)


	return router
}
