package event

import (
	"context"
	"fmt"
	Entity "tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)

/*
	The file name should have the format [Ticket ID]-ticket.html.
	The content doesn't matter, it's just important that it contain the ticket ID, price, and amount.
*/

func (h Handler) PrintTicketToFile(ctx context.Context, event Entity.TicketBookingConfirmed) error {
    log.FromContext(ctx).Info("Save confirmed ticket to file...")
	fileName := fmt.Sprintf("%s-ticket.html", event.TicketID)
	content := `
		<html>
			<head>
				<title>Ticket</title>
			</head>
			<body>
				<h1>Ticket ` + event.TicketID + `</h1>
				<p>Price: ` + event.Price.Amount + ` ` + event.Price.Currency + `</p>	
			</body>
		</html>

		err := h.printingTicketService.StoreTicketContent(
			ctx, 
			ticketID,
			content,
		)
		if err != nil {
			log.FromContext(ctx).Error("Fail to save ticket to file")
			return err
		}
		return nil
	}
	`


	err := h.printingTicketService.StoreTicketContent(
		ctx, 
		fileName,
		content,
	)
	if err != nil {
		log.FromContext(ctx).Error("Fail to save ticket to file")
		return err
	}

	// perform publish event ticket printed
	header := Entity.NewMessageHeader()
	eventTicketPrinted := Entity.TicketPrinted{
		Header: header,
		TicketID: event.TicketID,
		FileName: fileName,
	}
	err = h.eventBus.Publish(ctx, eventTicketPrinted)
	if err != nil {
		return fmt.Errorf("failed to publish TicketPrinted event: %w", err)
	}
	return nil
}
