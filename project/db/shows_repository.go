package db

import (
	"context"
	"fmt"
	Entity "tickets/entities"
	"github.com/jmoiron/sqlx"
)




type ShowsRepository struct {
	db *sqlx.DB
}


func NewShowsRepository(db *sqlx.DB) ShowsRepository {
	if db == nil {
		panic("db is nil")
	} else {
		return ShowsRepository{
			db: db,
		}
	}
}

func (t ShowsRepository) AddShow(ctx context.Context, show Entity.Show) error {
	query := `
	INSERT INTO
		shows (show_id, dead_nation_id, number_of_tickets, start_time, title, venue)
	VALUES
		(:show_id, :dead_nation_id, :number_of_tickets, :start_time, :title, :venue)
	ON CONFLICT DO NOTHING
	`
	_, err := t.db.NamedExecContext(
		ctx,
		query,
		show,
	)
	if err != nil {
		return fmt.Errorf("could not save show: %w", err)
	} else {
		return nil
	}
}



func (s ShowsRepository) ShowByID(ctx context.Context, showID string) (Entity.Show, error) {
	var show Entity.Show
	err := s.db.GetContext(ctx, &show, `SELECT * FROM shows WHERE show_id = $1`, showID)
	if err != nil {
		return Entity.Show{}, err
	}

	return show, nil
}
