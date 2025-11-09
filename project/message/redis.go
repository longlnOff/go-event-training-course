package message

import (
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)


func NewRedisPublisher(rdb redis.UniversalClient, logger watermill.LoggerAdapter) message.Publisher {
	pub, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client: rdb,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	return log.CorrelationPublisherDecorator{
		Publisher: pub,
	}
}

func NewRedisClient(address string) redis.UniversalClient {
	return redis.NewClient(&redis.Options{
		Addr: address,
	})
}
