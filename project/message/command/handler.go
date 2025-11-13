package command

import (
	"context"
	Entity "tickets/entities"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)


type Handler struct {
	receiptsService ReceiptsService
	paymentService PaymentService
}

type ReceiptsService interface {
	PutVoidReceipt(ctx context.Context, request Entity.RefundTicket) error
}

type PaymentService interface {
	RefundPayment(ctx context.Context, refundPayment Entity.PaymentRefund) error
}

func NewHandler(
	receiptsService ReceiptsService,
	paymentService PaymentService,
) *Handler {
	if receiptsService == nil {
		panic("missing receiptsService")
	}
	if paymentService == nil {
		panic("missing paymentService")
	}
	
	return &Handler{
		receiptsService: receiptsService,
		paymentService: paymentService,
	}
}


func (h *Handler) NewReceiptService() cqrs.CommandHandler {
	return cqrs.NewCommandHandler(
		"issue-refund-receipt",
		func(ctx context.Context, event *Entity.RefundTicket) error {
			return h.PutVoidReceipt(ctx, *event)
		},
	)
}
