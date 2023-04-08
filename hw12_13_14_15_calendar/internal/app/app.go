package app

import (
	"context"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type App struct {
	logger logger.Logger
	u      storage.EventUseCase
}

func New(logger logger.Logger, u storage.EventUseCase) *App {
	return &App{
		logger: logger,
		u:      u,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	a.logger.Info("create")
	return a.u.Create(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, eventId string, event storage.Event) error {
	return a.u.Update(ctx, eventId, event)
}

func (a *App) DeleteEvent(ctx context.Context, eventId string) error {
	return a.u.Delete(ctx, eventId)
}

func (a *App) GetDailyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.u.GetDailyEvents(ctx, date)
}

func (a *App) GetWeeklyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.u.GetWeeklyEvents(ctx, date)
}

func (a *App) GetMonthlyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.u.GetMonthlyEvents(ctx, date)
}

func getStorage(cfg *storage.Config) storage.Storage {
	if cfg.Storage == "sql" {
		return sqlstorage.New(cfg)
	}
	return memorystorage.New()
}

func GetEventUseCase(cfg *storage.Config) storage.EventUseCase {
	return *storage.New(getStorage(cfg))
}
