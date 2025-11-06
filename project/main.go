package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	ticketsMessage "tickets/message"
	"tickets/adapters"
	"tickets/service"
)

func main() {
	log.Init(slog.LevelInfo)

	apiClients, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	spreadsheetsAPI := adapters.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := adapters.NewReceiptsServiceClient(apiClients)
	rdb := ticketsMessage.NewRedisClient(os.Getenv("REDIS_ADDR"))
	err = service.New(
		rdb,
		spreadsheetsAPI,
		receiptsService,
	).Run(context.Background())
	if err != nil {
		panic(err)
	}
}
