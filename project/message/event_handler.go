package message

import (
	"context"
	"fmt"
	ticketDB "tickets/db"
	ticketEntity "tickets/entities"
	ticketEvent "tickets/message/event"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)


func NewWatermillProcessorWithEventHandler(
	receiptsService ticketEvent.ReceiptsService,
	spreadsheetsAPI ticketEvent.SpreadsheetsAPI,
	rdb redis.UniversalClient,
	watermillLogger watermill.LoggerAdapter,
	router *message.Router,
	db ticketDB.RepositoryDB,
	clients *clients.Clients,
	eventBus *cqrs.EventBus,
) *cqrs.EventProcessor {
	handler := ticketEvent.NewHandler(
		spreadsheetsAPI, 
		receiptsService, 
		db,
		clients,
		eventBus,
	)

	processor, err := ticketEvent.NewEventProcessor(
		router,
		rdb,
		watermillLogger,
	)
	if err != nil {
		panic(err)
	}

	err = processor.AddHandlers(cqrs.NewEventHandler(
		"issue-receipt",
		func(ctx context.Context, event *ticketEntity.TicketBookingConfirmed) error {
			err = handler.IssueReceipt(ctx, *event)
			if err != nil {
				return fmt.Errorf("failed to issue receipt: %w", err)
			}
			return nil
		},
	))
	if err != nil {
		panic(err)
	}

	err = processor.AddHandlers(cqrs.NewEventHandler(
		"append-to-tracker-confirmed",
		func(ctx context.Context, event *ticketEntity.TicketBookingConfirmed) error {
			err = handler.AppendToConfirmationTracker(ctx, *event)
			if err != nil {
				return fmt.Errorf("failed to issue receipt: %w", err)
			}
			return nil
		},
	))
	if err != nil {
		panic(err)
	}

	err = processor.AddHandlers(cqrs.NewEventHandler(
		"save-to-database",
		func(ctx context.Context, event *ticketEntity.TicketBookingConfirmed) error {
			err = handler.SaveToDatabase(ctx, *event)
			if err != nil {
				return fmt.Errorf("failed to save ticket to database: %w", err)
			}
			return nil
		},
	))
	if err != nil {
		panic(err)
	}

	err = processor.AddHandlers(cqrs.NewEventHandler(
		"store-html-ticket",
		func(ctx context.Context, event *ticketEntity.TicketBookingConfirmed) error {
			err = handler.StoreHtmlTicket(ctx, event)
			if err != nil {
				return fmt.Errorf("failed to save ticket to database: %w", err)
			}
			return nil
		},
	))
	if err != nil {
		panic(err)
	}


	err = processor.AddHandlers(cqrs.NewEventHandler(
		"append-to-tracker-canceled",
		func(ctx context.Context, event *ticketEntity.TicketBookingCanceled) error {
			err = handler.AppendToCancelationTracker(ctx, *event)
			if err != nil {
				return fmt.Errorf("failed to append to canceled tracker: %w", err)
			}
			return nil
		},
	))
	if err != nil {
		panic(err)
	}

	err = processor.AddHandlers(cqrs.NewEventHandler(
		"delete-ticket-from-database",
		func(ctx context.Context, event *ticketEntity.TicketBookingCanceled) error {
			err = handler.DeleteTicket(ctx, (*event).TicketID)
			if err != nil {
				return fmt.Errorf("failed to delete canceled ticket from database: %w", err)
			}
			return nil
		},
	))
	if err != nil {
		panic(err)
	}

	return processor
}
