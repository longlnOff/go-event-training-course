package command

import (
	"context"
	Entity "tickets/entities"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)


func (h Handler) PutVoidReceipt(ctx context.Context, event Entity.RefundTicket) error {
    log.FromContext(ctx).Info("Issuing refund receipt...")
	err := h.receiptsService.PutVoidReceipt(
		ctx, 
		event,
	)
	if err != nil {
		log.FromContext(ctx).Error("Fail to issue receipt")
		return err
	}

	refundPayment := Entity.PaymentRefund{
		TicketID: event.TicketID,
		RefundReason: "customer requested refund",
		IdempotencyKey: event.Header.IdempotencyKey,
	}
	log.FromContext(ctx).Info("Performing refund payment...")
	err = h.paymentService.RefundPayment(
		ctx, 
		refundPayment,
	)
	if err != nil {
		log.FromContext(ctx).Error("Fail to perform refund payment")
		return err
	}
	return nil
}
