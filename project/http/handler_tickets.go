package http

import (
	"fmt"
	"net/http"
	ticketsEntity "tickets/entities"
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
	idempotentKey := c.Request().Header.Get("Idempotency-Key")
	if idempotentKey == "" {
		return c.JSON(http.StatusBadRequest, "Idempotency-Key must be provided")
	}

	var request TicketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()

	for _, ticket := range request.Tickets {
		idempotentKey += ticket.TicketID

		header := ticketsEntity.NewMessageHeaderWithIdempotencyKey(idempotentKey)

		switch ticket.Status {
		case "confirmed":
			event := ticketsEntity.TicketBookingConfirmed{
				Header: header,
				TicketID: ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price: ticket.Price,
			}
			err = h.eventBus.Publish(ctx, event)
			if err != nil {
				return fmt.Errorf("failed to publish TicketBookingConfirmed event: %w", err)
			}
		case "canceled":
			event := ticketsEntity.TicketBookingCanceled{
				Header: header,
				TicketID: ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price: ticket.Price,
			}
			err = h.eventBus.Publish(ctx, event)
			if err != nil {
				return fmt.Errorf("failed to publish TicketBookingCanceled event: %w", err)

			}
		default:
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	
	}

	return c.NoContent(http.StatusOK)
}



func (h Handler) GetAllTicket(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.repo.FindAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, data)
	}
}
