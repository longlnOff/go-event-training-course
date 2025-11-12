package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	stdHTTP "net/http"
	Database "tickets/db"
	ticketsHttp "tickets/http"
	ticketsMessage "tickets/message"
	ticketsEvent "tickets/message/event"
	ticketsOutbox "tickets/message/outbox"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	db *sqlx.DB
	echoRouter *echo.Echo
	routerMessage *message.Router
}

func New(
	db *sqlx.DB,
	redisClient redis.UniversalClient,
	spreadsheetsAPI ticketsEvent.SpreadsheetsAPI,
	receiptsService ticketsEvent.ReceiptsService,
	printingTicketService ticketsEvent.PrintingTicketService,
	deadNationService ticketsEvent.DeadNationService,
) Service {
	watermillLogger := watermill.NewSlogLogger(slog.Default())

	redisPublisher := ticketsMessage.NewRedisPublisher(redisClient, watermillLogger)

	// repository
	ticketRepo := Database.NewTicketsRepository(db)
	showRepo := Database.NewShowsRepository(db)
	bookingRepo := Database.NewBookingsRepository(db)

	// ----------- EVENT BUS -----------
	eventBus, err := ticketsMessage.NewEventBusWithHandlers(redisPublisher, watermillLogger)
	if err != nil {
		panic(err)
	}

	// ----------- EVENT PROCESSOR, ROUTER and HANDLERS -----------
	postgresSubscriber := ticketsOutbox.NewPostgresSubscriber(db.DB, watermillLogger)
	router := ticketsMessage.NewMessageRouter(
		postgresSubscriber,
		redisPublisher,
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
		eventBus,
		ticketRepo,
		showRepo,
		spreadsheetsAPI,
		receiptsService,
		printingTicketService,
		deadNationService,
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
		&ticketRepo,
		&showRepo,
		&bookingRepo,
	)

	return Service{
		db: db,
		echoRouter: echoRouter,
		routerMessage: router,
	}
}

func (s Service) Run(ctx context.Context) error {
	if err := Database.InitializeSchema(s.db); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}
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
