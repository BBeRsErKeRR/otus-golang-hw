package sender

import (
	"context"
	"fmt"
	"time"

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

func (a *App) Consume(ctx context.Context) error {
	f := func(ctx context.Context, msg []byte) {
		eventDto := &storage.EventDTO{}
		if err := json.Unmarshal(msg, eventDto); err != nil {
			a.logger.Error("event notification unmarshal failed", zap.Error(err))
			return
		}
		a.logger.Debug(string(msg))
		notification := fmt.Sprintf("Send new notification from event '%s' to user '%v': %s -> %s",
			eventDto.Title,
			eventDto.UserID,
			eventDto.Date.Format(time.RFC822),
			eventDto.EndDate.Format(time.RFC822),
		)
		a.logger.Info(notification)
	}
	return a.u.Consume(ctx, f)
}

func getConsumer(cfg *consumer.Config, logger logger.Logger) consumer.Consumer {
	return internalrmqconsumer.New(cfg, logger)
}

func GetConsumerUseCase(cfg *consumer.Config, logger logger.Logger) consumer.ConsumerUseCase {
	return *consumer.New(getConsumer(cfg, logger))
}
