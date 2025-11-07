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

/*
	Add a new handler to the Router.
	It should subscribe to values from the temperature-celsius topic and publish 
	the converted values to the temperature-fahrenheit topic. 
	You can use the included celsiusToFahrenheit function to convert the values.
	Tip
		You don't need to call Ack() or Nack() in the Router's handler function. 
		The message is acknowledged if the returned error is nil.
*/


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
		"c-to-f-converter",
		"temperature-celsius",
		sub,
		"temperature-fahrenheit",
		pub,
		func(msg *message.Message) ([]*message.Message, error) {
			fTemperature, err := celsiusToFahrenheit(string(msg.Payload))
			if err != nil {
				return nil, err
			}
			newMsg := message.NewMessage(watermill.NewUUID(), []byte(fTemperature))
			return []*message.Message{newMsg}, nil
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
