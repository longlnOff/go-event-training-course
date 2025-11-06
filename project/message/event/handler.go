package event

import "context"

import (
	ticketEntity "tickets/entities"
	ticketDB "tickets/db"
)

type Handler struct {
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
	repository ticketDB.RepositoryDB
}

func NewHandler(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	repository ticketDB.RepositoryDB,
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
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request ticketEntity.IssueReceiptRequest) error
}
