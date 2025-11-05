package http

import (
	"net/http"
	ticketEntity "tickets/entities"
	"github.com/labstack/echo/v4"
)

var (
	ticketStatusConfirmed = "confirmed"
	ticketStatusCanceled = "canceled"
)

type TicketsStatusRequest struct {
	Tickets []TicketStatusRequest `json:"tickets"`
}

type TicketStatusRequest struct {
	TicketID      string `json:"ticket_id"`
	Status        string `json:"status"`
	Price         ticketEntity.Money  `json:"price"`
	CustomerEmail string `json:"customer_email"`
}


func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request TicketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}
	for i := range request.Tickets {
		ticket := request.Tickets[i]
		header := ticketEntity.NewMessageHeader()

		switch ticket.Status {
			case ticketStatusConfirmed:
				eventTicketBookingConfirmed := ticketEntity.TicketBookingConfirmed{
					Header: header,
					TicketID: ticket.TicketID,
					Price:    ticket.Price,
					CustomerEmail: ticket.CustomerEmail,
				}
				h.eventBus.Publish(c.Request().Context(), eventTicketBookingConfirmed)

			case ticketStatusCanceled:
				eventTicketBookingCanceled := ticketEntity.TicketBookingCanceled{
					Header: header,
					TicketID: ticket.TicketID,
					Price:    ticket.Price,
					CustomerEmail: ticket.CustomerEmail,
				}
				h.eventBus.Publish(c.Request().Context(), eventTicketBookingCanceled)
		}
	}

	return c.NoContent(http.StatusOK)
}
