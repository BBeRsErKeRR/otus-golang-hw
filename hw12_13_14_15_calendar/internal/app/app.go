package app

import (
	"context"

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

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.u.Create(ctx, storage.Event{ID: id, Title: title})
}

func getStorage(cfg *storage.Config) storage.Storage {
	if cfg.Storage == "postgres" {
		return sqlstorage.New(cfg)
	}
	return memorystorage.New()
}

func GetEventUseCase(cfg *storage.Config) storage.EventUseCase {
	return *storage.New(getStorage(cfg))
}
