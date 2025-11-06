package database

import (
	"context"
	ticketEntity "tickets/entities"

	"github.com/jmoiron/sqlx"
)


type RepositoryDB interface {
   SaveTicket(ctx context.Context, ticket ticketEntity.TicketBookingConfirmed) error
   RemoveTicket(ctx context.Context, ticketID string) error
   GetAllTicketWithoutFilter(ctx context.Context) ([]ticketEntity.TicketBookingConfirmed, error)
}

type PostgresDB struct {
	db *sqlx.DB
}

func NewPostgresDB(db *sqlx.DB) *PostgresDB {
	return &PostgresDB{
		db: db,
	}
}

func (p *PostgresDB) SaveTicket(ctx context.Context, ticket ticketEntity.TicketBookingConfirmed) error {
	query := `
		INSERT INTO 
			tickets (ticket_id, price_amount, price_currency, customer_email)
		VALUES 
			($1, $2, $3, $4)
		ON CONFLICT DO NOTHING
	`
	_, err := p.db.ExecContext(ctx, query,
		ticket.TicketID,
		ticket.Price.Amount,
		ticket.Price.Currency,
		ticket.CustomerEmail,
	)

	return err
}

func (p *PostgresDB) RemoveTicket(ctx context.Context, ticketID string) error {
	query := `DELETE FROM tickets WHERE ticket_id = $1`

	_, err := p.db.ExecContext(ctx, query, ticketID)

	return err
}

func (p *PostgresDB) GetAllTicketWithoutFilter(ctx context.Context) ([]ticketEntity.TicketBookingConfirmed, error) {
	query := `
		SELECT 
			ticket_id, price_amount, price_currency, customer_email
		FROM tickets
	`

	var tickets []struct {
		TicketID      string `db:"ticket_id"`
		PriceAmount   string `db:"price_amount"`
		PriceCurrency string `db:"price_currency"`
		CustomerEmail string `db:"customer_email"`
	}

	err := p.db.SelectContext(ctx, &tickets, query)
	if err != nil {
		return nil, err
	}

	result := make([]ticketEntity.TicketBookingConfirmed, len(tickets))
	for i, t := range tickets {
		result[i] = ticketEntity.TicketBookingConfirmed{
			TicketID:      t.TicketID,
			CustomerEmail: t.CustomerEmail,
			Price: ticketEntity.Money{
				Amount:   t.PriceAmount,
				Currency: t.PriceCurrency,
			},
		}
	}

	return result, nil
}
