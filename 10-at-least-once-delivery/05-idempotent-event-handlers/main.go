package main

import (
	"context"
	"fmt"
)

type PaymentTaken struct {
	PaymentID string
	Amount    int
}

type PaymentsHandler struct {
	repo *PaymentsRepository
	uniquePaymentIDs map[string]any
}

func NewPaymentsHandler(repo *PaymentsRepository) *PaymentsHandler {
	return &PaymentsHandler{
		repo: repo,
		uniquePaymentIDs: map[string]any{},
	}
}

func (p *PaymentsHandler) HandlePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	fmt.Println(p.uniquePaymentIDs[event.PaymentID])
	if _, ok := p.uniquePaymentIDs[event.PaymentID]; ok {
		return nil
	}
	p.uniquePaymentIDs[event.PaymentID] = true
	return p.repo.SavePaymentTaken(ctx, event)
}

type PaymentsRepository struct {
	payments []PaymentTaken
}

func (p *PaymentsRepository) Payments() []PaymentTaken {
	return p.payments
}

func NewPaymentsRepository() *PaymentsRepository {
	return &PaymentsRepository{}
}

func (p *PaymentsRepository) SavePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	p.payments = append(p.payments, *event)
	return nil
}
