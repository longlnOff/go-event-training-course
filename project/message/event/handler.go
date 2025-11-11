package event

import (
	"context"
	ticketsEntity "tickets/entities"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)


type Handler struct {
	eventBus        *cqrs.EventBus
	repo            TicketsRepository
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
	printingTicketService PrintingTicketService
}

type TicketsRepository interface {
	Add(ctx context.Context, ticket ticketsEntity.Ticket) error
	Remove(ctx context.Context, ticket ticketsEntity.Ticket) error

}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request ticketsEntity.IssueReceiptRequest) error
}

type PrintingTicketService interface {
	StoreTicketContent(ctx context.Context, fileID string, fileContent string) error
}

func NewHandler(
	eventBus *cqrs.EventBus,
	repo TicketsRepository,
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	printingTicketService PrintingTicketService,
) *Handler {
	if repo == nil {
		panic("missing repo")
	}
	if spreadsheetsAPI == nil {
		panic("missing spreadsheetsAPI")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}
	if printingTicketService == nil {
		panic("missing printingTicketService")
	}
	
	return &Handler{
		eventBus: eventBus,
		repo: repo,
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
		printingTicketService: printingTicketService,
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

func (h *Handler) NewStoreTicketHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"store-ticket",
		func(ctx context.Context, event *ticketsEntity.TicketBookingConfirmed) error {
			data := ticketsEntity.Ticket{
				TicketID: event.TicketID,
				Price: event.Price,
				CustomerEmail: event.CustomerEmail,
			}
			return h.StoreTicket(ctx, data)
		},
	)
}

func (h *Handler) NewRemoveCanceledTicketHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"remove-canceled-ticket",
		func(ctx context.Context, event *ticketsEntity.TicketBookingCanceled) error {
			removedTicket := ticketsEntity.Ticket{
				TicketID: event.TicketID,
				Price: event.Price,
				CustomerEmail: event.CustomerEmail,
			}
			return h.RemoveTicket(ctx, removedTicket)
		},
	)
}

func (h *Handler) NewPrintTicketToFileHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"print-ticket-to-file",
		func(ctx context.Context, event *ticketsEntity.TicketBookingConfirmed) error {
			return h.PrintTicketToFile(ctx, *event)
		},
	)
}
