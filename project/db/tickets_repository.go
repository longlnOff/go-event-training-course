package db

import (
	"context"
	"fmt"
	ticketsEntity "tickets/entities"
	"github.com/jmoiron/sqlx"
)




type TicketsRepository struct {
	db *sqlx.DB
}


func NewTicketsRepository(db *sqlx.DB) TicketsRepository {
	if db == nil {
		panic("db is nil")
	} else {
		return TicketsRepository{
			db: db,
		}
	}
}

func (t TicketsRepository) Add(ctx context.Context, ticket ticketsEntity.Ticket) error {
	query := `
	INSERT INTO
		tickets (ticket_id, price_amount, price_currency, customer_email)
	VALUES
		(:ticket_id, :price.amount, :price.currency, :customer_email)
	ON CONFLICT DO NOTHING
	`
	_, err := t.db.NamedExecContext(
		ctx,
		query,
		ticket,
	)
	if err != nil {
		return fmt.Errorf("could not save ticket: %w", err)
	} else {
		return nil
	}
}

func (t TicketsRepository) Remove(ctx context.Context, ticket ticketsEntity.Ticket) error {
	query := `
	DELETE FROM
		tickets
	WHERE
		ticket_id = :ticket_id
	`
	_, err := t.db.NamedExecContext(
		ctx,
		query,
		ticket,
	)
	if err != nil {
		return fmt.Errorf("could not remove ticket: %w", err)
	} else {
		return nil
	}
}


func (t TicketsRepository) FindAll(ctx context.Context) ([]ticketsEntity.Ticket, error) {
    var returnTickets []ticketsEntity.Ticket

    err := t.db.SelectContext(
        ctx,
        &returnTickets, `
            SELECT
                ticket_id,
                price_amount as "price.amount",
                price_currency as "price.currency",
                customer_email
            FROM
                tickets
        `,
    )
    if err != nil {
        return nil, err
    }

    return returnTickets, nil
}
