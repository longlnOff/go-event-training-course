package http

import (
	"net/http"
	Database "tickets/db"
	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	ticketsRepo *Database.TicketsRepository,
	showsRepo *Database.ShowsRepository,
	bookingsRepo *Database.BookingsRepository,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		eventBus: eventBus,
		ticketRepository: ticketsRepo,
		showRepository: showsRepo,
		bookingRepository: bookingsRepo,
	}
	

	// Tickets
	e.POST("/tickets-status", handler.PostTicketStatus)
	e.GET("tickets", handler.GetAllTicket)

	// Health
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	//Shows
	e.POST("shows", handler.CreateShow)

	// Bookings
	e.POST("book-tickets", handler.BookTickets)

	return e
}
