package message

// This is unused, but if we need to write custom decorator, use it

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/lithammer/shortuuid/v3"
)

type contextKey string

const (
	CorrelationIdInHttp contextKey = "Correlation-Id"
	TypeMessage contextKey = "type"
)

type CorrelationPublisherDecorator struct {
	message.Publisher
}

func (c CorrelationPublisherDecorator) Publish(topic string, messages ...*message.Message) error {
	// custom logic here
	for i := range messages {
		message := messages[i]
		// Set correlation_id for message
		correlationID := CorrelationIDFromContext(message.Context())
		message.Metadata.Set("correlation_id", correlationID)

		// Set type for message
		typeMsg := MessageTypeFromContext(message.Context())
		message.Metadata.Set("type", typeMsg)

	}

	return c.Publisher.Publish(topic, messages...)
}


func ContextWithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIdInHttp, correlationID)
}

func CorrelationIDFromContext(ctx context.Context) string {
	v, ok := ctx.Value(CorrelationIdInHttp).(string)
	if ok {
		return v
	}

	// add "gen_" prefix to distinguish generated correlation IDs from correlation IDs passed by the client
	// it's useful to detect if correlation ID was not passed properly
	return "gen_" + shortuuid.New()
}

func ContextWithType(ctx context.Context, typeMsg string) context.Context {
	return context.WithValue(ctx, TypeMessage, typeMsg)
}


func MessageTypeFromContext(ctx context.Context) string {
	v, ok := ctx.Value(TypeMessage).(string)
	if ok {
		return v
	}

	return ""
}
