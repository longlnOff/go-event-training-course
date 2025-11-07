package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	ticketsEntity "tickets/entities"
	ticketsMessage "tickets/message"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type TicketsStatusRequest struct {
	Tickets []TicketStatusRequest `json:"tickets"`
}

type TicketStatusRequest struct {
	TicketID      string `json:"ticket_id"`
	Status        string `json:"status"`
	Price         ticketsEntity.Money  `json:"price"`
	CustomerEmail string `json:"customer_email"`
}

func (h Handler) PostTicketStatus(c echo.Context) error {
	var request TicketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		if ticket.Status == "confirmed" {
			payloadTracker := ticketsEntity.AppendToTrackerPayload{
				TicketID: ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price: ticket.Price,
			}
			dataTracker, err := json.Marshal(payloadTracker)
			if err != nil {
				return err
			}
			msgTracker := message.NewMessage(watermill.NewUUID(), []byte(dataTracker))

			err = h.pub.Publish(ticketsMessage.AppendToTrackerTopic, msgTracker)
			if err != nil {
				return err
			}

			payloadReceipt := ticketsEntity.IssueReceiptPayload{
				TicketID: ticket.TicketID,
				Price: ticket.Price,
			}
			dataReceipt, err := json.Marshal(payloadReceipt)
			if err != nil {
				return err
			}
			msgReceipt := message.NewMessage(watermill.NewUUID(), []byte(dataReceipt))

			err = h.pub.Publish(ticketsMessage.IssueReceiptTopic, msgReceipt)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	
	}

	return c.NoContent(http.StatusOK)
}



