package http

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	Entity "tickets/entities"
)

type Handler struct {
	eventBus *cqrs.EventBus
	commandBus *cqrs.CommandBus
	ticketRepository TicketsRepository
	showRepository ShowsRepository
	bookingRepository BookingsRepository
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}

type TicketsRepository interface {
	FindAll(ctx context.Context) ([]Entity.Ticket, error)
}

type ShowsRepository interface {
	AddShow(ctx context.Context, show Entity.Show) error
}

type BookingsRepository interface {
	AddBooking(ctx context.Context, booking Entity.Booking) error
}
