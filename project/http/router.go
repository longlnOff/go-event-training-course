package http

import (
	"net/http"
	ticketDB "tickets/database"
	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	db ticketDB.RepositoryDB,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		eventBus: eventBus,
		db: db,
	}

	e.POST("/tickets-status", handler.PostTicketsConfirmation)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.GET("/tickets", handler.GetAllTicketWithoutFilter)

	return e
}
