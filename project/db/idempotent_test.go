package database

import (
	"context"
	"testing"
	ticketEntity "tickets/entities"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestIdempotentSave(t *testing.T) {
	db := InitDatabase()
	postgresDB := PostgresDB{
		db: db,
	}


	ctx := context.Background()
	ticketToAdd := ticketEntity.TicketBookingConfirmed{
		Header: ticketEntity.NewMessageHeader(),
		TicketID: uuid.NewString(),
		Price: ticketEntity.Money{
			Amount: "30.00",
			Currency: "USD",
		},
		CustomerEmail: "lllll@gmail.com",
	}


	for i := 0; i < 2; i++ {
		err := postgresDB.SaveTicket(ctx, ticketToAdd)
		require.NoError(t, err)


		tickets, err := postgresDB.GetAllTicketWithoutFilter(ctx)
		require.NoError(t, err)
		require.Len(t, tickets, 1)

		foundTickets := lo.Filter(tickets, func(t ticketEntity.TicketBookingConfirmed, _ int) bool {
			return t.TicketID == ticketToAdd.TicketID
		})
		require.Len(t, foundTickets, 1)
	}

}
