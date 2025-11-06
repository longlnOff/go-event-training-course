package service

import (
	"context"
	"errors"
	"log/slog"
	stdHTTP "net/http"
	ticketsHttp "tickets/http"
	ticketsMessage "tickets/message"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	echoRouter *echo.Echo
}

func New(
	redisClient redis.UniversalClient,
	spreadsheetsAPI ticketsMessage.SpreadsheetsAPI,
	receiptsService ticketsMessage.ReceiptsService,
) Service {
	watermillLogger := watermill.NewSlogLogger(slog.Default())

	redisPublisher := ticketsMessage.NewRedisPublisher(redisClient, watermillLogger)

	ticketsMessage.NewHandlers(
		redisClient,
		watermillLogger,
		spreadsheetsAPI,
		receiptsService,
	)

	echoRouter := ticketsHttp.NewHttpRouter(
		redisPublisher,
	)

	return Service{
		echoRouter,
	}
}

func (s Service) Run(ctx context.Context) error {
	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
