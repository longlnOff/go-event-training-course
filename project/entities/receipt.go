package entities

type IssueReceiptRequest struct {
    IdempotencyKey string
    TicketID string
    Price    Money
}

type IssueReceiptPayload struct {
    TicketID string `json:"ticket_id"`
    Price    Money  `json:"price"`
}
