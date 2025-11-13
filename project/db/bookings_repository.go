package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	Entity "tickets/entities"
	"tickets/message/event"
	"tickets/message/outbox"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)




type BookingsRepository struct {
	db *sqlx.DB
}


func NewBookingsRepository(db *sqlx.DB) BookingsRepository {
	if db == nil {
		panic("db is nil")
	} else {
		return BookingsRepository{
			db: db,
		}
	}
}

func (t BookingsRepository) AddBooking(ctx context.Context, booking Entity.Booking) error {
	return updateInTx(
		ctx,
		t.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			availableSeats := 0
			err := tx.GetContext(ctx, &availableSeats, `
				SELECT
					number_of_tickets AS available_seats
				FROM
					shows
				WHERE
					show_id = $1
			`, booking.ShowID)
			if err != nil {
				return fmt.Errorf("could not get available seats: %w", err)
			}

			alreadyBookedSeats := 0
			err = tx.GetContext(ctx, &alreadyBookedSeats, `
				SELECT
					COALESCE(SUM(number_of_tickets), 0) AS already_booked_seats
				FROM
					bookings
				WHERE
					show_id = $1
			`, booking.ShowID)
			if err != nil {
				return fmt.Errorf("could not get already booked seats: %w", err)
			}

			if availableSeats-alreadyBookedSeats < booking.NumberOfTickets {
				// this is usually a bad idea, learn more here: https://threedots.tech/post/introducing-clean-architecture/
				// we'll improve it later
				return echo.NewHTTPError(http.StatusBadRequest, "not enough seats available")
			}

			query := `
			INSERT INTO
				bookings (booking_id, show_id, number_of_tickets, customer_email)
			VALUES
				(:booking_id, :show_id, :number_of_tickets, :customer_email)
			ON CONFLICT DO NOTHING
			`
			_, err = t.db.NamedExecContext(
				ctx,
				query,
				booking,
			)
			if err != nil {
				return fmt.Errorf("could not save show: %w", err)
			}
			outboxPublisher, err := outbox.NewPublisherForDB(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create outbox publisher: %w", err)
			}
			bus, err := event.NewBus(outboxPublisher)
			if err != nil {
				return fmt.Errorf("could not create event bus: %w", err)
			}

			bookingMade := Entity.BookingMade{
				Header: Entity.NewMessageHeader(),
				BookingID: booking.BookingID,
				ShowID: booking.ShowID,
				NumberOfTickets: booking.NumberOfTickets,
				CustomerEmail: booking.CustomerEmail,
			}
			err = bus.Publish(ctx, bookingMade)
			if err != nil {
				return fmt.Errorf("could not publish event: %w", err)
			}
			return nil
		},
	)
}






