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
		header := ticketsEntity.NewMessageHeader()

		switch ticket.Status {
		case "confirmed":
			event := ticketsEntity.TicketBookingConfirmed{
				Header: header,
				TicketID: ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price: ticket.Price,
			}
			data, err := json.Marshal(event)
			if err != nil {
				return err
			}
			msg := message.NewMessage(watermill.NewUUID(), []byte(data))
			
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))

			err = h.pub.Publish(ticketsMessage.TicketBookingConfirmedTopic, msg)
			if err != nil {
				return err
			}
		case "canceled":
			event := ticketsEntity.TicketBookingCanceled{
				Header: header,
				TicketID: ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price: ticket.Price,
			}
			data, err := json.Marshal(event)
			if err != nil {
				return err
			}
			msg := message.NewMessage(watermill.NewUUID(), []byte(data))

			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))

			err = h.pub.Publish(ticketsMessage.TicketBookingCanceledTopic, msg)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	
	}

	return c.NoContent(http.StatusOK)
}



