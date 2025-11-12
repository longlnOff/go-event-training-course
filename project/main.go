package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"tickets/adapters"
	ticketsMessage "tickets/message"
	"tickets/service"
	_ "github.com/lib/pq"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/jmoiron/sqlx"
)

func main() {
	log.Init(slog.LevelInfo)


	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Database
	db, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	apiClients, err := clients.NewClients(
		os.Getenv("GATEWAY_ADDR"), 
		func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Correlation-ID", log.CorrelationIDFromContext(ctx))
		return nil
	},
	)
	if err != nil {
		panic(err)
	}

	spreadsheetsAPI := adapters.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := adapters.NewReceiptsServiceClient(apiClients)
	printingTicketSerivce := adapters.NewPrintingTicketsAPIClient(apiClients)
	deadNationService := adapters.NewDeadNationClient(apiClients)

	rdb := ticketsMessage.NewRedisClient(os.Getenv("REDIS_ADDR"))
	err = service.New(
		db,
		rdb,
		spreadsheetsAPI,
		receiptsService,
		printingTicketSerivce,
		deadNationService,
	).Run(ctx)
	if err != nil {
		panic(err)
	}
}
