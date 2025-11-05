package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

func main() {

	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()

	channelMessage, err := subscriber.Subscribe(context.Background(), "progress")
	if err != nil {
		panic(err)
	}

	for msg := range channelMessage {
		messageID := string(msg.UUID)
		orderID := string(msg.Payload)
		fmt.Printf("Message ID: %v - %v\n", messageID, orderID)
		msg.Ack()
	}

}
