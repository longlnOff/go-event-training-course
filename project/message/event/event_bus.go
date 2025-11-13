package event

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)




func NewEventBusWithHandlers(
	publisher message.Publisher,
	logger watermill.LoggerAdapter,
) (*cqrs.EventBus, error) {

	config := cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error){
			return fmt.Sprintf("events.%s", params.EventName), nil

		},
		Marshaler: marshalerJSON,
		Logger: logger,
	}


	return cqrs.NewEventBusWithConfig(
		publisher,
		config,
	)
}


func NewBus(
	publisher message.Publisher,
) (*cqrs.EventBus, error) {

	config := cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error){
			return fmt.Sprintf("events.%s", params.EventName), nil
		},
		Marshaler: marshalerJSON,
	}


	return cqrs.NewEventBusWithConfig(
		publisher,
		config,
	)
}
