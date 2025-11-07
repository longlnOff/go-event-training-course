package event

import (
	"context"
	"log/slog"
	ticketsEntity "tickets/entities"
)


func (h Handler) IssueReceipt(ctx context.Context, event ticketsEntity.TicketBookingConfirmed) error {
    slog.Info("Issuing receipt...")
    // ...
	request := ticketsEntity.IssueReceiptRequest{
		TicketID: event.TicketID,
		Price: event.Price,
	}
	err := h.receiptsService.IssueReceipt(
		ctx, 
		request,
	)
	if err != nil {
		slog.Error("Fail to issue receipt")
		return err
	}
	return nil
}
