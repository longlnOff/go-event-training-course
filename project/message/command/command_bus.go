package command

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewCommandBusWithHandlers(
	publisher message.Publisher,
	logger watermill.LoggerAdapter,
) (*cqrs.CommandBus, error) {

	config := cqrs.CommandBusConfig{
		GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return fmt.Sprintf("commands.%s", params.CommandName), nil

		},
		Marshaler: marshalerJSON,
		Logger: logger,
	}


	return cqrs.NewCommandBusWithConfig(
		publisher,
		config,
	)
}


func NewCommandBus(
	publisher message.Publisher,
) (*cqrs.CommandBus, error) {

	config := cqrs.CommandBusConfig{
		GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return fmt.Sprintf("commands.%s", params.CommandName), nil
		},
		Marshaler: marshalerJSON,
	}


	return cqrs.NewCommandBusWithConfig(
		publisher,
		config,
	)
}
