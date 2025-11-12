package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	Entity "tickets/entities"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)


type BookingTicketsRequest struct {
	ShowID    string    	`json:"show_id"`
	NumberOfTickets int		`json:"number_of_tickets"`
	CustomerEmail   string 	`json:"customer_email"`
}

type CreateBookingResponse struct {
	BookingID string `json:"booking_id"`
}

func (h Handler) BookTickets(c echo.Context) error {
	var request BookingTicketsRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	if request.NumberOfTickets < 1 {
		return c.JSON(http.StatusBadRequest, "Number of tickets must be greater than 0")
	}
	
	ctx := c.Request().Context()
	booking := Entity.Booking{
		BookingID: watermill.NewUUID(),
		ShowID: request.ShowID,
		NumberOfTickets: request.NumberOfTickets,
		CustomerEmail:   request.CustomerEmail,
	}
	err = h.bookingRepository.AddBooking(ctx, booking)
	if err != nil {
		var httpErr *echo.HTTPError
		if errors.As(err, &httpErr) && httpErr.Code == http.StatusBadRequest {
			return c.JSON(http.StatusBadRequest, httpErr.Message)
		}
		slog.Error(fmt.Sprintf("Error while add new booking: %v", err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	resp := CreateBookingResponse{
		BookingID: booking.BookingID,
	}

	return c.JSON(http.StatusCreated, resp)
}


