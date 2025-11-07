package event

import (
	"context"
	"log/slog"
	ticketsEntity "tickets/entities"
)

func (h Handler) AppendToCancel(ctx context.Context, event ticketsEntity.TicketBookingCanceled) error {
    slog.Info("Appending ticket to the cancel")
    // ...
	err := h.spreadsheetsAPI.AppendRow(
		ctx, 
		"tickets-to-refund",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
	if err != nil {
		slog.Error("Fail to append ticket to the cancel")
		return err
	}
	return nil
}
