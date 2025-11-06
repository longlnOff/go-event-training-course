package event

import (
	"context"

	ticketEntity "tickets/entities"

	ticketDB "tickets/db"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
	eventBus *cqrs.EventBus
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
	repository ticketDB.RepositoryDB
	client *clients.Clients
}

func NewHandler(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	repository ticketDB.RepositoryDB,
	client *clients.Clients,
	eventBus *cqrs.EventBus,
) Handler {
	if spreadsheetsAPI == nil {
		panic("missing spreadsheetsAPI")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}

	return Handler{
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
		repository: repository,
		client: client,
		eventBus: eventBus,
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request ticketEntity.IssueReceiptRequest) error
}
