package main

/*
	1. Create a Redis Streams publisher.
	2. Publish two messages on the progress topic.
	   The first one's payload should be 50, and the second one's should be 100.
*/

import (
	"os"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)



func main() {
	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client: rdb,
		},
		logger,
	)

	if err != nil {
		panic(err)
	}

	err = publisher.Publish("progress", message.NewMessage(watermill.NewUUID(), []byte("50")))
	if err != nil {
		panic(err)
	}

	err = publisher.Publish("progress", message.NewMessage(watermill.NewUUID(), []byte("100")))
	if err != nil {
		panic(err)
	}
}
