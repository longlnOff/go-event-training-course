package event

import (
	"context"
	Entity "tickets/entities"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)


func (h Handler) IssueReceipt(ctx context.Context, event Entity.TicketBookingConfirmed) error {
    log.FromContext(ctx).Info("Issuing receipt...")
    // ...
	request := Entity.IssueReceiptRequest{
		IdempotencyKey: event.Header.IdempotencyKey,
		TicketID: event.TicketID,
		Price: event.Price,
	}
	err := h.receiptsService.IssueReceipt(
		ctx, 
		request,
	)
	if err != nil {
		log.FromContext(ctx).Error("Fail to issue receipt")
		return err
	}
	return nil
}
