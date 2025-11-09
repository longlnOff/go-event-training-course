package main

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func RegisterEventHandlers(
	sub message.Subscriber,
	router *message.Router,
	handlers []cqrs.EventHandler,
	logger watermill.LoggerAdapter,
) error {
	config := cqrs.EventProcessorConfig{
								GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
									return params.EventName, nil
								},
								Marshaler: cqrs.JSONMarshaler{
									GenerateName: cqrs.StructName,
								},
								SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error){
									return sub, nil
								},
	}


	processor, err := cqrs.NewEventProcessorWithConfig(
		router,
		config,
	)
	if err != nil {
		return err
	}

	return processor.AddHandlers(handlers...)
}


