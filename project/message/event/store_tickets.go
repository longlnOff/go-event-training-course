package event

import (
	"context"
	ticketsEntity "tickets/entities"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)

func (h Handler) StoreTicket(ctx context.Context, event ticketsEntity.Ticket) error {
    log.FromContext(ctx).Info("Storing ticket to the database")
	err := h.repo.Add(ctx, event)
	if err != nil {
		log.FromContext(ctx).Error("Fail to store ticket to the database")
		return err
	}
	return nil
}
