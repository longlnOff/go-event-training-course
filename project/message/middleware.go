package message

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func useMiddlewares(router *message.Router, watermillLogger watermill.LoggerAdapter) {
	router.AddMiddleware(middleware.Recoverer)

	router.AddMiddleware(middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Millisecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          watermillLogger,
	}.Middleware)

	// Add correlation ID
	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			correlationID := msg.Metadata.Get("correlation_id")
			if correlationID == "" {
				correlationID = watermill.NewUUID()
			}
			ctx := msg.Context()
			ctx = log.ToContext(ctx, slog.With("correlation_id", correlationID))
			ctx = log.ContextWithCorrelationID(ctx, correlationID)

			msg.SetContext(ctx)

			return next(msg)
		}
	})


	// Add logging middleware
	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
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
				// Handle the error 
				logger := log.FromContext(msg.Context()).With(
					"error", err,
					"message_id", msg.UUID,
				)
				logger.Error("Error while handling a message")
			}
			
			return msgs, err
		}
	})
}
