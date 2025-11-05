package main

import (
	"context"
	"os"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewSlogLogger(nil)

	router := message.NewDefaultRouter(logger)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	// TODO: Add your handler here
	router.AddHandler(
		"c-to-f-temperature-converter",
		"temperature-celsius",
		sub,
		"temperature-fahrenheit",
		pub,
		func(msg *message.Message) ([]*message.Message, error) {
			cTemperature := string(msg.Payload)
			fTemperature, err := celsiusToFahrenheit(cTemperature)
			if err != nil {
				return []*message.Message{}, err
			} else {
				newMessage := message.NewMessage(watermill.NewUUID(), []byte(fTemperature))
				return []*message.Message{newMessage}, nil
			}
		},
	)

	err = router.Run(context.Background())
	if err != nil {
		panic(err)
	}
}

func celsiusToFahrenheit(temperature string) (string, error) {
	celsius, err := strconv.Atoi(temperature)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(celsius*9/5 + 32), nil
}
