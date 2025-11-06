package http

import (
	"context"
	ticketsWorker "tickets/worker"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}

// type Handler struct {
// 	spreadsheetsAPI SpreadsheetsAPI
// 	receiptsService ReceiptsService
// }

type Handler struct {
	worker *ticketsWorker.Worker
}
