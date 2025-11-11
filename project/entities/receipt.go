package entities

type IssueReceiptRequest struct {
    IdempotencyKey string
    TicketID string
    Price    Money
}
