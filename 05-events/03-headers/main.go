package main

import (
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type header struct {
	ID         string `json:"id"`
	EventName  string `json:"event_name"`
	OccurredAt string `json:"occurred_at"`
}

type ProductOutOfStock struct {
	Header    header `json:"header"`
	ProductID string `json:"product_id"`
}

type ProductBackInStock struct {
	Header    header `json:"header"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type Publisher struct {
	pub message.Publisher
}

func NewMessageHeader(eventName string) header {
	return header{
		ID:          uuid.NewString(), 
		EventName:   eventName, 
		OccurredAt: time.Now().Format(time.RFC3339),
	}
}


func NewPublisher(pub message.Publisher) Publisher {
	return Publisher{
		pub: pub,
	}
}

func (p Publisher) PublishProductOutOfStock(productID string) error {
	header := NewMessageHeader("ProductOutOfStock")
	event := ProductOutOfStock{
		Header: header,
		ProductID: productID,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	return p.pub.Publish("product-updates", msg)
}

func (p Publisher) PublishProductBackInStock(productID string, quantity int) error {
	header := NewMessageHeader("ProductBackInStock")
	event := ProductBackInStock{
		Header: header,
		ProductID: productID,
		Quantity:  quantity,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	return p.pub.Publish("product-updates", msg)
}
