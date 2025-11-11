package adapters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)

type PrintingTicketsAPIClient struct {
	// we are not mocking this client: it's pointless to use interface here
	clients *clients.Clients
}

func NewPrintingTicketsAPIClient(clients *clients.Clients) *PrintingTicketsAPIClient {
	if clients == nil {
		panic("NewPrintingTicketsAPIClient: clients is nil")
	}

	return &PrintingTicketsAPIClient{clients: clients}
}

func (c PrintingTicketsAPIClient) StoreTicketContent(ctx context.Context, fileID string, fileContent string) error {
	resp, err := c.clients.Files.PutFilesFileIdContentWithTextBodyWithResponse(ctx, fileID, fileContent)
	if err != nil {
		return fmt.Errorf("failed to save ticket to file: %w", err)
	}

	if resp.StatusCode() == http.StatusConflict {
		log.FromContext(ctx).With("file", fileID).Info("file already exists")
		return nil
	}

	return nil
}
