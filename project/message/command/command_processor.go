package command

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func RegisterCommandHandlers(
	processor *cqrs.CommandProcessor,
	handler *Handler,
) error {
	handlers := []cqrs.CommandHandler{}
	handlers = append(handlers, handler.NewReceiptService())


	return processor.AddHandlers(handlers...)
}

func NewCommandProcessor(
	router *message.Router,
	rdb redis.UniversalClient,
	logger watermill.LoggerAdapter,
) (*cqrs.CommandProcessor, error) {
	config := cqrs.CommandProcessorConfig{
		SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return redisstream.NewSubscriber(redisstream.SubscriberConfig{
				Client:        rdb, 
				ConsumerGroup: "svc-tickets.commands." + params.HandlerName,
			}, logger)
		},
		GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
			return fmt.Sprintf("commands.%s", params.CommandName), nil
		},
		Marshaler: marshalerJSON,
	}

	return cqrs.NewCommandProcessorWithConfig(
		router,
		config,
	)
}
