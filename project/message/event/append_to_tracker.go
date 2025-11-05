package event

import (
	"context"
	"log/slog"
	"tickets/entities"
)

const (
	ticketToPrintSheetName = "tickets-to-print"
	ticketToRefundSheetName = "tickets-to-refund"
)

func (h Handler) AppendToConfirmationTracker(ctx context.Context, event entities.TicketBookingConfirmed) error {
    slog.Info("Appending ticket to the tracker")
    // ...
	h.spreadsheetsAPI.AppendRow(ctx, ticketToPrintSheetName, []string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency})

	return	nil
}


func (h Handler) AppendToCancelationTracker(ctx context.Context, event entities.TicketBookingCanceled) error {
    slog.Info("Appending ticket to the tracker")
    // ...
	h.spreadsheetsAPI.AppendRow(ctx, ticketToRefundSheetName, []string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency})

	return	nil
}
