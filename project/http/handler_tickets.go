package http

import (
	"log/slog"
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

type TicketsResponse struct {
	Tickets []TicketResponse `json:"tickets"`
}

type TicketResponse struct {
	TicketID      string             `json:"ticket_id"`
	CustomerEmail string             `json:"customer_email"`
	Price         ticketEntity.Money `json:"price"`
	Status        string             `json:"status"`
}

func (h Handler) GetAllTicketWithoutFilter(c echo.Context) error {
	tickets, err := h.db.GetAllTicketWithoutFilter(c.Request().Context())
	if err != nil {
		slog.Error("failed to get all tickets", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to retrieve tickets",
		})
	}

	ticketResponses := make([]TicketResponse, len(tickets))
	for i, ticket := range tickets {
		ticketResponses[i] = TicketResponse{
			TicketID:      ticket.TicketID,
			CustomerEmail: ticket.CustomerEmail,
			Price:         ticket.Price,
			Status:        ticketStatusConfirmed,
		}
	}

	return c.JSON(http.StatusOK, ticketResponses)
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
