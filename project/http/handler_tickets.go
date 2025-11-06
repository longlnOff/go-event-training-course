package http

import (
	"net/http"
	ticketsMessage "tickets/message"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type ticketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request ticketsConfirmationRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		msg := message.NewMessage(watermill.NewUUID(), []byte(ticket))

		err = h.pub.Publish(ticketsMessage.AppendToTrackerTopic, msg)
		if err != nil {
			return err
		}
		err = h.pub.Publish(ticketsMessage.IssueReceiptTopic, msg)
		if err != nil {
			return err
		}	
	}

	return c.NoContent(http.StatusOK)
}
