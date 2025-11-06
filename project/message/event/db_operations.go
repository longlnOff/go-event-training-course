package event

import (
	"context"
	"log/slog"
	"tickets/entities"
)


func (h Handler) SaveToDatabase(ctx context.Context, event entities.TicketBookingConfirmed) error {
    slog.Info("Save ticket to database")
    // ...
	return h.repository.SaveTicket(ctx, event)
}


func (h Handler) DeleteTicket(ctx context.Context, ticketID string) error {
	slog.Info("Delete canceled ticket from database")

	return h.repository.RemoveTicket(ctx, ticketID)
}
