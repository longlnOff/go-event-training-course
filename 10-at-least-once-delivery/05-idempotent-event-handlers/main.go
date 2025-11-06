package main

import (
	"context"
)

type PaymentTaken struct {
	PaymentID string
	Amount    int
}

type PaymentsHandler struct {
	repo *PaymentsRepository
}

func NewPaymentsHandler(repo *PaymentsRepository) *PaymentsHandler {
	return &PaymentsHandler{repo: repo}
}

func (p *PaymentsHandler) HandlePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	if p.repo.CheckUniquePayment(ctx, event) {
		return nil
	}
	p.repo.paymentKeys[event.PaymentID] = true
	return p.repo.SavePaymentTaken(ctx, event)
}

type PaymentsRepository struct {
	payments []PaymentTaken
	paymentKeys map[string]bool
}

func (p *PaymentsRepository) Payments() []PaymentTaken {
	return p.payments
}

func NewPaymentsRepository() *PaymentsRepository {
	return &PaymentsRepository{
		payments: []PaymentTaken{},
		paymentKeys: map[string]bool{},
	}
}

func (p *PaymentsRepository) SavePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	p.payments = append(p.payments, *event)
	return nil
}

func (p *PaymentsRepository) CheckUniquePayment(ctx context.Context, event *PaymentTaken) bool {
	return p.paymentKeys[event.PaymentID]
}


