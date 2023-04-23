package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq/producer"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

type App struct {
	logger   logger.Logger
	duration time.Duration
	p        producer.Producer
	sU       storage.EventUseCase
}

func New(logger logger.Logger, p producer.Producer, sU storage.EventUseCase, duration time.Duration) *App {
	return &App{
		logger:   logger,
		p:        p,
		sU:       sU,
		duration: duration,
	}
}

func (a *App) Publish(ctx context.Context, data []byte) error {
	return a.p.Publish(ctx, data)
}

func (a *App) Obsolescence(ctx context.Context) error {
	a.logger.Info("Starting obsolescence old events")
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
		a.logger.Info(fmt.Sprintf("publish event %v, %s for user %v", event.ID, event.Title, event.UserID))
		encoded, err := json.Marshal(event)
		if err != nil {
			a.logger.Error("fail marshal event", zap.Error(err))
			return
		}
		if err := a.Publish(ctx, encoded); err != nil {
			a.logger.Error("fail publish event", zap.Error(err))
		}
	}
}

func (a *App) Run(ctx context.Context) error {
	ticker := time.NewTicker(a.duration)
	defer ticker.Stop()
	a.logger.Info("Starting scheduler")
	defer a.logger.Info("Stopping scheduler")
	for {
		go func() {
			a.PublishEvents(ctx)
			err := a.Obsolescence(ctx)
			if err != nil {
				a.logger.Error(fmt.Sprintf("fail delete old events %s", err))
			}
			a.logger.Info(fmt.Sprintf("Next operation start after %s", a.duration))
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
