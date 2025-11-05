package main

import (
	"os"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

/*
	1. Create a Redis Streams publisher.
	2. Publish two messages on the `progress` topic.
		The first one's payload shoule be 50, and
		the second one's should be 100


*/
func main() {
	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	msg1 := message.NewMessage(watermill.NewUUID(), []byte("50"))
	msg2 := message.NewMessage(watermill.NewUUID(), []byte("100"))
	if err = publisher.Publish("progress", msg1); err != nil {
		panic(err)
	}
	if err = publisher.Publish("progress", msg2); err != nil {
		panic(err)
	}
}
