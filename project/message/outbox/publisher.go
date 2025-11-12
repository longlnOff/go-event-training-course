package outbox

import (
	"context"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPublisherForDB(
	ctx context.Context,
	tx *sqlx.Tx,
) (message.Publisher, error) {
	var publisher message.Publisher
	var err error
	logger := watermill.NewSlogLogger(nil)

	publisher, err = watermillSQL.NewPublisher(
		tx,
		watermillSQL.PublisherConfig{
			SchemaAdapter: watermillSQL.DefaultPostgreSQLSchema{},
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	publisher = log.CorrelationPublisherDecorator{Publisher: publisher}

	publisher = forwarder.NewPublisher(publisher, forwarder.PublisherConfig{
		ForwarderTopic: outboxTopic,
	})
	publisher = log.CorrelationPublisherDecorator{Publisher: publisher}

	return publisher, nil
}
