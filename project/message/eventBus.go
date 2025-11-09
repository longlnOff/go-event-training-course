package message

import (
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
			return params.EventName, nil
		},
		Marshaler: marshalerJSON,
		Logger: logger,
	}


	return cqrs.NewEventBusWithConfig(
		publisher,
		config,
	)
}
