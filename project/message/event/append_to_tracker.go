package event

import (
	"context"
	"log/slog"
	ticketsEntity "tickets/entities"
)


func (h Handler) AppendToPrint(ctx context.Context, event ticketsEntity.TicketBookingConfirmed) error {
    slog.Info("Appending ticket to the print")
    // ...
	err := h.spreadsheetsAPI.AppendRow(
		ctx, 
		"tickets-to-print",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
	if err != nil {
		slog.Error("Fail to append ticket to the print")
		return err
	}
	return nil
}
