package service

import (
	"context"
	"errors"
	"log/slog"
	stdHTTP "net/http"
	ticketsMessage "tickets/message"
	ticketsEvent "tickets/message/event"
	ticketsHttp "tickets/http"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	echoRouter *echo.Echo
	routerMessage *message.Router
}

func New(
	redisClient redis.UniversalClient,
	spreadsheetsAPI ticketsEvent.SpreadsheetsAPI,
	receiptsService ticketsEvent.ReceiptsService,
) Service {
	watermillLogger := watermill.NewSlogLogger(slog.Default())

	redisPublisher := ticketsMessage.NewRedisPublisher(redisClient, watermillLogger)

	// ----------- EVENT BUS -----------
	eventBus, err := ticketsMessage.NewEventBusWithHandlers(redisPublisher, watermillLogger)
	if err != nil {
		panic(err)
	}


	// ----------- EVENT PROCESSOR, ROUTER and HANDLERS -----------
	router := ticketsMessage.NewMessageRouter(
		redisClient,
		watermillLogger,
	)

	processor, err := ticketsMessage.NewEventProcessor(
		router,
		redisClient,
		watermillLogger,
	)
	if err != nil {
		panic(err)
	}
	
	eventHandlers := ticketsEvent.NewHandler(
		spreadsheetsAPI,
		receiptsService,
	)

	err = ticketsMessage.RegisterEventHandlers(
		processor,
		eventHandlers,
	)
	if err != nil {
		panic(err)
	}
	// ------------------------------------------------------------------


	echoRouter := ticketsHttp.NewHttpRouter(
		eventBus,
	)

	return Service{
		echoRouter: echoRouter,
		routerMessage: router,
	}
}

func (s Service) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := s.routerMessage.Run(context.Background())
		if err != nil {
			// TODO: we will improve it in a next exercise
			slog.With("error", err).Error("Failed to run watermill router")
			return err
		}
		return nil
	})

	g.Go(func() error {
		// when routerMessage running, channel will be closed --> next, httpRouter will start
		// we do this to ensure http server run after message is running
		<-s.routerMessage.Running()
		err := s.echoRouter.Start(":8080")
		if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		// Shut down the HTTP server
		<-ctx.Done()
		return s.echoRouter.Shutdown(ctx)
	})

	return g.Wait()
}
