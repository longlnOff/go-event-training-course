package entities

type IssueReceiptRequest struct {
    TicketID string
    Price    Money
}

type IssueReceiptPayload struct {
    TicketID string `json:"ticket_id"`
    Price    Money  `json:"price"`
}
