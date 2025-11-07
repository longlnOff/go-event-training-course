package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

/*
	Create a new Router, and add a no-publisher handler to it.
	The handler should subscribe to the temperature-fahrenheit
	topic and print the incoming values in the following format:
		Temperature read: 100
	Call Run() to run the Router.

*/

func main() {
	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	router := message.NewDefaultRouter(logger)

	router.AddConsumerHandler(
		"temperature-fahrenheit-printer",
		"temperature-fahrenheit",
		sub, 
		func(msg *message.Message) error {
			fmt.Printf("Temperature read: %s\n", msg.Payload)
			return nil
		},
	)

	if err = router.Run(context.Background()); err != nil {
		panic(err)
	}

}
