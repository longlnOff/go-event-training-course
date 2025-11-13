package message

import (
	ticketsOutbox "tickets/message/outbox"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewMessageRouter(
	postgresSubscriber message.Subscriber,
	publisher message.Publisher,
	logger watermill.LoggerAdapter,
) *message.Router {
	router := message.NewDefaultRouter(logger)

	useMiddleware(router, logger)
	
	ticketsOutbox.AddForwarderHandler(postgresSubscriber, publisher, router, logger)
	return router
}
