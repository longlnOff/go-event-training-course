package event

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	ticketEntity "tickets/entities"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)



func (h Handler) StoreHtmlTicket(ctx context.Context, event *ticketEntity.TicketBookingConfirmed) error {
    slog.Info("Store ticket content")
    // ...
	// help me please
	body := fmt.Sprintf("%s %s %s", event.TicketID, event.Price.Amount, event.Price.Currency)
	fileName := fmt.Sprintf("%s-ticket.html", event.TicketID)
	response, err := h.client.Files.PutFilesFileIdContentWithTextBodyWithResponse(
		ctx,
	fileName,
		body,
	)
	if err != nil {
		return err
	}

	if response.StatusCode() == http.StatusConflict {
		log.FromContext(ctx).With("file", fileName).Info("file already exists")
	}

	eventTicketPrinted := ticketEntity.TicketPrinted{
		Header: event.Header,
		TicketID: event.TicketID,
		FileName: fileName,
	}

	err = h.eventBus.Publish(ctx, eventTicketPrinted)
	if err != nil {
		return fmt.Errorf("failed to publish TicketPrinted event: %w", err)
	}
	
	return nil
}

