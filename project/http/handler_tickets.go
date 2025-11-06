package http

import (
	"net/http"
	ticketsWorker "tickets/worker"
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
		msgReceipt := ticketsWorker.Message{
			Task:     ticketsWorker.TaskIssueReceipt,
			TicketID: ticket,
		}

		msgTracker := ticketsWorker.Message{
			Task:     ticketsWorker.TaskAppendToTracker,
			TicketID: ticket,
		}

		// send
		h.worker.Send(msgReceipt)
		h.worker.Send(msgTracker)
	}

	return c.NoContent(http.StatusOK)
}
