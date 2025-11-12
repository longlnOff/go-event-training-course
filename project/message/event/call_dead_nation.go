package event

import (
	"context"
	Entity "tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/google/uuid"
)


func (h Handler) CallDeadNation(ctx context.Context, event Entity.BookingMade) error {
	log.FromContext(ctx).Info("Call dead nation...")
	show, err := h.ShowRepo.ShowByID(ctx, event.ShowID)
	if err != nil {
		log.FromContext(ctx).Error("Fail to get show")
		return err
	}
	deadNationRequest := Entity.DeadNationBooking{
		BookingID: uuid.MustParse(event.BookingID),
		NumberOfTickets: event.NumberOfTickets,
		CustomerEmail: event.CustomerEmail,
		DeadNationEventID: uuid.MustParse(show.DeadNationID),

	}
	err = h.deadNationService.PostTicketBookingWithResponse(
		ctx, 
		deadNationRequest,
	)
	if err != nil {
		log.FromContext(ctx).Error("Fail to call dead nation")
		return err
	}
	return nil
}
