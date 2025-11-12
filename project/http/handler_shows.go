package http

import (
	"fmt"
	"log/slog"
	"net/http"
	Entity "tickets/entities"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)


type CreateShowRequest struct {
	DeadNationID    string    `json:"dead_nation_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	StartTime       time.Time `json:"start_time"`
	Title           string    `json:"title"`
	Venue           string    `json:"venue"`
}

type CreateShowResponse struct {
	ShowID string `json:"show_id"`
}

func (h Handler) CreateShow(c echo.Context) error {
	var request CreateShowRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}
	
	ctx := c.Request().Context()
	show := Entity.Show{
		ShowID: watermill.NewUUID(),
		DeadNationID:    request.DeadNationID,
		NumberOfTickets: request.NumberOfTickets,
		StartTime:       request.StartTime,
		Title:           request.Title,
		Venue:           request.Venue,
	}
	err = h.showRepository.AddShow(ctx, show)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while add new show: %v", err))
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp := CreateShowResponse{
		ShowID: show.ShowID,
	}

	return c.JSON(http.StatusCreated, resp)
}


