package http

import (
	"context"
	ticketEntity "tickets/entities"
	ticketDB "tickets/db"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request ticketEntity.IssueReceiptRequest) error
}

type Handler struct {
	// publisher message.Publisher
	eventBus *cqrs.EventBus
	db ticketDB.RepositoryDB
}
