package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq/consumer"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

type App struct {
	logger logger.Logger
	c      consumer.Consumer
}

func New(logger logger.Logger, c consumer.Consumer) *App {
	return &App{
		logger: logger,
		c:      c,
	}
}

func (a *App) Consume(ctx context.Context) error {
	f := func(ctx context.Context, msg []byte) {
		var encodedStatus []byte
		eventDto := &storage.EventDTO{}
		err := json.Unmarshal(msg, eventDto)
		if err == nil {
			a.logger.Debug(string(msg))
			notification := fmt.Sprintf("Send new notification from event '%s' to user '%v': %s -> %s",
				eventDto.Title,
				eventDto.UserID,
				eventDto.Date.Format(time.RFC822),
				eventDto.EndDate.Format(time.RFC822),
			)
			a.logger.Info(notification)
			encodedStatus = []byte(fmt.Sprintf("Successful send %v", eventDto.Title))
		} else {
			a.logger.Error("event notification unmarshal failed", zap.Error(err))
			encodedStatus = []byte(fmt.Sprintf("Unmarshal error: %v", err))
		}

		if err := a.c.PublishStatus(ctx, encodedStatus); err != nil {
			a.logger.Error("fail publish event", zap.Error(err))
		}
	}
	return a.c.Consume(ctx, f)
}
