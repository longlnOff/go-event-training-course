package event

import (
	"context"
	Entity "tickets/entities"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)

func (h Handler) StoreTicket(ctx context.Context, event Entity.Ticket) error {
    log.FromContext(ctx).Info("Storing ticket to the database")
	err := h.ticketRepo.Add(ctx, event)
	if err != nil {
		log.FromContext(ctx).Error("Fail to store ticket to the database")
		return err
	}
	return nil
}
