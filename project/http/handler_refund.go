package http

import (
	"fmt"
	"net/http"
	Entity "tickets/entities"
	"github.com/labstack/echo/v4"
)


func (h Handler) PutTicketRefund(c echo.Context) error {
	ticketID := c.Param("ticket_id")

	
	ctx := c.Request().Context()
	refundTicketCommand := Entity.RefundTicket{
		Header: Entity.NewMessageHeaderWithIdempotencyKey(ticketID),
		TicketID: ticketID,
	}
	err := h.commandBus.Send(ctx, refundTicketCommand)
	if err != nil {
		return fmt.Errorf("failed to publish RefundTicket event: %w", err)
	}
	return c.JSON(http.StatusAccepted, nil)
}


