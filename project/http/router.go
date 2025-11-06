package http

import (
	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
	ticketsWorker "tickets/worker"
)

func NewHttpRouter(
	w *ticketsWorker.Worker,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		worker: w,
	}
	

	e.POST("/tickets-confirmation", handler.PostTicketsConfirmation)

	return e
}
