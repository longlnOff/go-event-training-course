package event

import (
	"context"
	Entity "tickets/entities"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)

func (h Handler) RemoveTicket(ctx context.Context, event Entity.Ticket) error {
    log.FromContext(ctx).Info("Removing ticket to the database")
	err := h.ticketRepo.Remove(ctx, event)
	if err != nil {
		log.FromContext(ctx).Error("Fail to remove ticket to the database")
		return err
	}
	return nil
}
