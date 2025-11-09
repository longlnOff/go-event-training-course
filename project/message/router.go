package message

import (
	ticketsEvent "tickets/message/event"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)


func RegisterEventHandlers(
	processor *cqrs.EventProcessor,
	handler *ticketsEvent.Handler,
) error {
	handlers := []cqrs.EventHandler{}
	handlers = append(handlers, handler.NewIssueReceiptHandler())
	handlers = append(handlers, handler.NewAppendToTrackerPrinttHandler())
	handlers = append(handlers, handler.NewAppendToRefundtHandler())

	return processor.AddHandlers(handlers...)
}



func NewEventProcessor(
	router *message.Router,
	rdb redis.UniversalClient,
	logger watermill.LoggerAdapter,
) (*cqrs.EventProcessor, error) {
	config := cqrs.EventProcessorConfig{
		SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return redisstream.NewSubscriber(redisstream.SubscriberConfig{
				Client:        rdb, 
				ConsumerGroup: "svc-tickets." + params.HandlerName,
			}, logger)
		},
		GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
			return params.EventName, nil
		},
		Marshaler: marshalerJSON,
	}

	return cqrs.NewEventProcessorWithConfig(
		router,
		config,
	)
}




func NewMessageRouter(
	rdb redis.UniversalClient, 
	logger watermill.LoggerAdapter,
) *message.Router {
	router := message.NewDefaultRouter(logger)

	useMiddleware(router, logger)

	return router
}
