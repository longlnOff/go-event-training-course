package event

import (
	"context"
	Entity "tickets/entities"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)


type Handler struct {
	eventBus        *cqrs.EventBus
	ticketRepo		TicketsRepository
	ShowRepo 		ShowsRepository
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
	printingTicketService PrintingTicketService
	deadNationService DeadNationService
}

type TicketsRepository interface {
	Add(ctx context.Context, ticket Entity.Ticket) error
	Remove(ctx context.Context, ticket Entity.Ticket) error
}

type ShowsRepository interface {
	AddShow(ctx context.Context, show Entity.Show) error
	ShowByID(ctx context.Context, showID string) (Entity.Show, error)
}

type DeadNationService interface {
	PostTicketBookingWithResponse(ctx context.Context, request Entity.DeadNationBooking) error
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}


type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request Entity.IssueReceiptRequest) error
}

type PrintingTicketService interface {
	StoreTicketContent(ctx context.Context, fileID string, fileContent string) error
}

func NewHandler(
	eventBus *cqrs.EventBus,
	ticketRepo TicketsRepository,
	showRepo ShowsRepository,
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	printingTicketService PrintingTicketService,
	deadNationService DeadNationService,
) *Handler {
	if ticketRepo == nil {
		panic("missing repo")
	}
	if eventBus == nil {
		panic("missing eventBus")
	}
	if showRepo == nil {
		panic("missing showRepo")
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
	if deadNationService == nil {
		panic("missing deadNationService")
	}
	
	return &Handler{
		eventBus: eventBus,
		ticketRepo: ticketRepo,
		ShowRepo: showRepo,
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
		printingTicketService: printingTicketService,
		deadNationService: deadNationService,
	}
}


func (h *Handler) NewIssueReceiptHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"issue-receipt",
		func(ctx context.Context, event *Entity.TicketBookingConfirmed) error {
			return h.IssueReceipt(ctx, *event)
		},
	)
}

func (h *Handler) NewAppendToTrackerPrinttHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"append-to-tracker",
		func(ctx context.Context, event *Entity.TicketBookingConfirmed) error {
			return h.AppendToPrint(ctx, *event)
		},
	)
}

func (h *Handler) NewAppendToRefundtHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"append-to-refund",
		func(ctx context.Context, event *Entity.TicketBookingCanceled) error {
			return h.AppendToCancel(ctx, *event)
		},
	)
}

func (h *Handler) NewStoreTicketHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"store-ticket",
		func(ctx context.Context, event *Entity.TicketBookingConfirmed) error {
			data := Entity.Ticket{
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
		func(ctx context.Context, event *Entity.TicketBookingCanceled) error {
			removedTicket := Entity.Ticket{
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
		func(ctx context.Context, event *Entity.TicketBookingConfirmed) error {
			return h.PrintTicketToFile(ctx, *event)
		},
	)
}

func (h *Handler) NewDeadNationHandler() cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"call-dead-nation",
		func(ctx context.Context, event *Entity.BookingMade) error {
			return h.CallDeadNation(ctx, *event)
		},
	)
}
