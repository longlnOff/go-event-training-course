package message

import (
	ticketEvent "tickets/message/event"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)


func NewWatermillRouter(
	receiptsService ticketEvent.ReceiptsService,
	spreadsheetsAPI ticketEvent.SpreadsheetsAPI,
	rdb redis.UniversalClient,
	watermillLogger watermill.LoggerAdapter,
) *message.Router {
	router := message.NewDefaultRouter(watermillLogger)
	// Add middleware
	useMiddlewares(router, watermillLogger)

	return router
}
