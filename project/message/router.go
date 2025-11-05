package message

import (
	// "encoding/json"
	// "fmt"
	// "log/slog"
	// ticketEntity "tickets/entities"
	ticketEvent "tickets/message/event"

	"github.com/ThreeDotsLabs/watermill"
	// "github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)


func NewWatermillRouter(
	receiptsService ticketEvent.ReceiptsService,
	spreadsheetsAPI ticketEvent.SpreadsheetsAPI,
	rdb redis.UniversalClient,
	watermillLogger watermill.LoggerAdapter,
) *message.Router {

	// handler := ticketEvent.NewHandler(spreadsheetsAPI, receiptsService)
	router := message.NewDefaultRouter(watermillLogger)
	// Add middleware
	useMiddlewares(router, watermillLogger)

	// issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
	// 	Client:        rdb,
	// 	ConsumerGroup: "issue-receipt",
	// }, watermillLogger)
	// if err != nil {
	// 	panic(err)
	// }

	// appendToTrackerConfirmedSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
	// 	Client:        rdb,
	// 	ConsumerGroup: "append-to-tracker-confirmed",
	// }, watermillLogger)
	// if err != nil {
	// 	panic(err)
	// }

	// appendToTrackerCanceledSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
	// 	Client:        rdb,
	// 	ConsumerGroup: "append-to-tracker-canceled",
	// }, watermillLogger)
	// if err != nil {
	// 	panic(err)
	// }

	// router.AddConsumerHandler(
	// 	"issue_receipt",
	// 	TicketBookingConfirmedTopic,
	// 	issueReceiptSub,
	// 	func(msg *message.Message) error {
	// 		var event ticketEntity.TicketBookingConfirmed
	// 		err := json.Unmarshal(msg.Payload, &event)
	// 		if err != nil {
	// 			slog.Error("Failed to unmarshal payload", "error", err)
	// 			return nil		// nil for remove message
	// 		}
	// 		err = handler.IssueReceipt(msg.Context(), event)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to issue receipt: %w", err)
	// 		}

	// 		return nil
	// 	},
	// )

	// router.AddConsumerHandler(
	// 	"append_to_tracker_confirmed",
	// 	TicketBookingConfirmedTopic,
	// 	appendToTrackerConfirmedSub,
	// 	func(msg *message.Message) error {
	// 		var event ticketEntity.TicketBookingConfirmed
	// 		err := json.Unmarshal(msg.Payload, &event)
	// 		if err != nil {
	// 			slog.Error("Failed to unmarshal payload", "error", err)
	// 			return nil		// nil for remove message
	// 		}
	// 		err = handler.AppendToConfirmationTracker(msg.Context(), event)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to append to confirmed tracker: %w", err)
	// 		}

	// 		return nil
	// 	},
	// )

	// router.AddConsumerHandler(
	// 	"append_to_tracker_canceled",
	// 	TicketBookingCanceledTopic,
	// 	appendToTrackerCanceledSub,
	// 	func(msg *message.Message) error {
	// 		var event ticketEntity.TicketBookingCanceled
	// 		err := json.Unmarshal(msg.Payload, &event)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to unmarshal payload: %w", err)
	// 		}

	// 		err = handler.AppendToCancelationTracker(msg.Context(), event)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to append to canceled tracker: %w", err)
	// 		}

	// 		return nil
	// 	},
	// )

	return router
}
