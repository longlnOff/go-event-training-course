package service

import (
	"context"
	"errors"
	"log/slog"
	stdHTTP "net/http"
	ticketsHttp "tickets/http"
	ticketsMessage "tickets/message"
	ticketsEvent "tickets/message/event"
	ticketDB "tickets/database"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	echoRouter *echo.Echo
	watermillRouter *message.Router
	db ticketDB.RepositoryDB
}

func New(
	rdb redis.UniversalClient,
	spreadsheetsAPI ticketsHttp.SpreadsheetsAPI,
	receiptsService ticketsHttp.ReceiptsService,
	db ticketDB.RepositoryDB,
) Service {
	watermillLogger := watermill.NewSlogLogger(log.FromContext(context.Background()))
	publisher:= ticketsEvent.NewRedisPublisher(rdb, watermillLogger)

	eventBus, err := ticketsEvent.NewEventBus(publisher)
	if err != nil {
		panic(err)
	}

	watermillRouter := ticketsMessage.NewWatermillRouter(
		receiptsService,
		spreadsheetsAPI,
		rdb,
		watermillLogger,
	)

	eventProcessor := ticketsMessage.NewWatermillProcessorWithEventHandler(
		receiptsService,
		spreadsheetsAPI,
		rdb,
		watermillLogger,
		watermillRouter,
		db,
	)
	_ = eventProcessor

	echoRouter := ticketsHttp.NewHttpRouter(eventBus, db)

	return Service{
		echoRouter: echoRouter,
		watermillRouter: watermillRouter,
		db: db,
	}
}

func (s Service) RunRouter(ctx context.Context) error {
	err := s.watermillRouter.Run(ctx)
	if err != nil {
		// TODO: we will improve it in a next exercise
		slog.With("error", err).Error("Failed to run watermill router")
		return err
	}
	return nil
}

func (s Service) Run(ctx context.Context) error {
	errgrp, ctx := errgroup.WithContext(ctx)
	
	errgrp.Go(func() error {
		return s.watermillRouter.Run(ctx)
	})
	
	errgrp.Go(func() error {
		<-s.watermillRouter.Running()

		err := s.echoRouter.Start(":8080")

		if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
			return err
		}
		// Why <- on Line 84?
		// <-s.watermillRouter.Running()
		// Here's what's happening:
		// s.watermillRouter.Running() returns a channel of type chan struct{}
		// The <- operator receives from that channel
		// This is a blocking operation - the goroutine pauses here until something is received
		// When Watermill's router finishes starting up, it closes this channel
		// Receiving from a closed channel immediately returns, so the goroutine unblocks and continues
		// Think of it like this:
		// runningChannel := s.watermillRouter.Running()  // Get the channel
		// <-runningChannel                               // Wait/block until it closes
		// // Now we know the router is running!
		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(context.Background())
	})

	return errgrp.Wait()
}

func (s Service) Shutdown(ctx context.Context) error {
	return s.echoRouter.Shutdown(ctx)
}
