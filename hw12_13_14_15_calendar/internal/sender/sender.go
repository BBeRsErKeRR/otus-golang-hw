package sender

import (
	"context"
	"fmt"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq/consumer"
	internalrmqconsumer "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq/consumer/rmq"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

type App struct {
	logger logger.Logger
	u      consumer.ConsumerUseCase
}

func New(logger logger.Logger, u consumer.ConsumerUseCase) *App {
	return &App{
		logger: logger,
		u:      u,
	}
}

func GetConsumer(cfg *consumer.Config, logger logger.Logger) consumer.Consumer {
	return internalrmqconsumer.New(cfg, logger)
}

func GetConsumerUseCase(cons consumer.Consumer) consumer.ConsumerUseCase {
	return *consumer.New(cons)
}

func (a *App) Consume(ctx context.Context) error {
	f := func(ctx context.Context, msg []byte) {
		eventDto := &storage.EventDTO{}
		if err := json.Unmarshal(msg, eventDto); err != nil {
			a.logger.Error("event notification unmarshal failed", zap.Error(err))
			return
		}
		notification := fmt.Sprintf("Send new notification from event '%s' to user '%v' -> %s",
			eventDto.Title,
			eventDto.UserID,
			eventDto.Date,
		)
		a.logger.Info(notification)
	}
	return a.u.Consume(ctx, f)
}
