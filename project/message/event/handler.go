package event

import (
	"context"
	ticketsEntity "tickets/entities"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)


type Handler struct {
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request ticketsEntity.IssueReceiptRequest) error
}

func NewHandler(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
) *Handler {
	if spreadsheetsAPI == nil {
		panic("missing spreadsheetsAPI")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}

	return &Handler{
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
	}
}


func (h *Handler) NewIssueReceiptHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"issue-receipt",
		func(ctx context.Context, event *ticketsEntity.TicketBookingConfirmed) error {
			return h.IssueReceipt(ctx, *event)
		},
	)
}

func (h *Handler) NewAppendToTrackerPrinttHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"append-to-tracker",
		func(ctx context.Context, event *ticketsEntity.TicketBookingConfirmed) error {
			return h.AppendToPrint(ctx, *event)
		},
	)
}

func (h *Handler) NewAppendToRefundtHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"append-to-refund",
		func(ctx context.Context, event *ticketsEntity.TicketBookingCanceled) error {
			return h.AppendToCancel(ctx, *event)
		},
	)
}
