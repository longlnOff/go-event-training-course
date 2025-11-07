package event

import (
	"context"
	ticketsEntity "tickets/entities"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)

func (h Handler) AppendToCancel(ctx context.Context, event ticketsEntity.TicketBookingCanceled) error {
    log.FromContext(ctx).Info("Appending ticket to the cancel")
    // ...
	err := h.spreadsheetsAPI.AppendRow(
		ctx, 
		"tickets-to-refund",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
	if err != nil {
		log.FromContext(ctx).Error("Fail to append ticket to the cancel")
		return err
	}
	return nil
}
