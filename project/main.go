package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"tickets/adapters"
	"tickets/service"
	_ "github.com/lib/pq"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/redis/go-redis/v9"
	ticketDB "tickets/db"
	
)

func main() {

	// DB
	db := ticketDB.InitDatabase()
	postgresDB := ticketDB.NewPostgresDB(db)
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log.Init(slog.LevelInfo)

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

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	spreadsheetsAPI := adapters.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := adapters.NewReceiptsServiceClient(apiClients)

	err = service.New(
		rdb,
		spreadsheetsAPI,
		receiptsService,
		postgresDB,
	).Run(ctx)
	if err != nil {
		panic(err)
	}
}
