package http

import (
	"net/http"
	ticketsDB "tickets/db"
	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	repo *ticketsDB.TicketsRepository,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		eventBus: eventBus,
		repo: repo,
	}
	

	e.POST("/tickets-status", handler.PostTicketStatus)
	e.GET("tickets", handler.GetAllTicket)

	// Health
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	return e
}
