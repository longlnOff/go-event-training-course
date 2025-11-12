package adapters

import (
	"context"
	"fmt"
	"net/http"
	Entity "tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients/dead_nation"
)

type DeadNationClient struct {
	// we are not mocking this client: it's pointless to use interface here
	clients *clients.Clients
}

func NewDeadNationClient(clients *clients.Clients) *DeadNationClient {
	if clients == nil {
		panic("NewDeadNationClient: clients is nil")
	}


	return &DeadNationClient{clients: clients}
}

func (c DeadNationClient) PostTicketBookingWithResponse(ctx context.Context, request Entity.DeadNationBooking) error {
	resp, err := c.clients.DeadNation.PostTicketBookingWithResponse(
		ctx,
		dead_nation.PostTicketBookingRequest{
			BookingId:       request.BookingID,
			EventId:         request.DeadNationEventID,
			NumberOfTickets: request.NumberOfTickets,
			CustomerAddress: request.CustomerEmail,
		},
	)

	if resp.StatusCode() != http.StatusOK {

		return fmt.Errorf("failed to post ticket booking: %w", err)
	}

	return err
}
