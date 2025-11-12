package event

import (
	"context"
	Entity "tickets/entities"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)


func (h Handler) AppendToPrint(ctx context.Context, event Entity.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Appending ticket to the tracker")

    // ...
	err := h.spreadsheetsAPI.AppendRow(
		ctx, 
		"tickets-to-print",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
	if err != nil {
		log.FromContext(ctx).Error("Fail to append ticket to the print")
		return err
	}
	return nil
}
