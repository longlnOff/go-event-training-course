package service

import (
	"context"
	"errors"
	stdHTTP "net/http"

	"github.com/labstack/echo/v4"
	ticketsWorker "tickets/worker"
	ticketsHttp "tickets/http"
)

type Service struct {
	echoRouter *echo.Echo
	worker *ticketsWorker.Worker
}

func New(
	spreadsheetsAPI ticketsWorker.SpreadsheetsAPI,
	receiptsService ticketsWorker.ReceiptsService,
) Service {

	worker := ticketsWorker.NewWorker(
		spreadsheetsAPI,
		receiptsService,
	)

	echoRouter := ticketsHttp.NewHttpRouter(worker)

	return Service{
		echoRouter: echoRouter,
		worker: worker,
	}
}

func (s Service) Run(ctx context.Context) error {
	go s.worker.Run(ctx)
	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
