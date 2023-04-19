package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq/producer"
	internalrmqproducer "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq/producer/rmq"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	logger   logger.Logger
	duration time.Duration
	pU       producer.ProducerUseCase
	sU       storage.EventUseCase
}

func New(logger logger.Logger, pU producer.ProducerUseCase, sU storage.EventUseCase, duration time.Duration) *App {
	return &App{
		logger:   logger,
		pU:       pU,
		sU:       sU,
		duration: duration,
	}
}

func (a *App) Publish(ctx context.Context, event storage.Event) error {
	return a.pU.Publish(ctx, event)
}

func (a *App) Obsolescence(ctx context.Context) error {
	return a.sU.DeleteBeforeDate(ctx, time.Now().AddDate(-1, 0, 0))
}

func (a *App) PublishEvents(ctx context.Context) {
	startDate := time.Now().Add(-a.duration)
	endDate := time.Now()

	events, err := a.sU.GetEventsByPeriod(ctx, startDate, endDate)
	if err != nil {
		a.logger.Error("fail get event", zap.Error(err))
	}

	for _, event := range events {
		if err := a.Publish(ctx, event); err != nil {
			a.logger.Error("fail publish event", zap.Error(err))
		}
	}
}

func (a *App) Run(ctx context.Context) error {
	ticker := time.NewTicker(a.duration)
	defer ticker.Stop()

	for {
		go func() {
			a.PublishEvents(ctx)
			err := a.Obsolescence(ctx)
			if err != nil {
				a.logger.Error(fmt.Sprintf("fail delete old events %s", err))
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func GetProducer(cfg *producer.Config, logger logger.Logger) producer.Producer {
	return internalrmqproducer.New(cfg, logger)
}

func GetProducerUseCase(p producer.Producer) producer.ProducerUseCase {
	return *producer.New(p)
}
