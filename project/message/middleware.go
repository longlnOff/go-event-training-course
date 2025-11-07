package message

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/lithammer/shortuuid/v3"
)


func useMiddleware(router *message.Router, logger watermill.LoggerAdapter) {
	router.AddMiddleware(middleware.Recoverer)

	router.AddMiddleware(CorrelationID)

	router.AddMiddleware(LoggerMiddleware)

	router.AddMiddleware(middleware.Retry{
				MaxRetries:      10,
				InitialInterval: time.Millisecond * 100,
				MaxInterval:     time.Second,
				Multiplier:      2,
				Logger:          logger,
			}.Middleware)

}

func LoggerMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := log.FromContext(msg.Context()).With(
			"message_id", msg.UUID,
			"payload", string(msg.Payload),
			"metadata", msg.Metadata,
			"handler", message.HandlerNameFromCtx(msg.Context()),
		)

		logger.Info("Handling a message")
		msgs, err := next(msg)
		if err != nil {
			logger.With(
				"error", err,
			).Error("Error while handling a message")
		}
		return msgs, err

	}
}


func CorrelationID(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		correlationID := msg.Metadata.Get("correlation_id")
		if correlationID == "" {
			correlationID = shortuuid.New()
		}

		ctx := log.ContextWithCorrelationID(msg.Context(), correlationID)
		ctx = log.ToContext(ctx, slog.With("correlation_id", correlationID))

		msg.SetContext(ctx)

		return next(msg)
	}
}
