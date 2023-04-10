package router

import (
	"context"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type Application interface {
	CreateEvent(context.Context, storage.Event) (string, error)
	UpdateEvent(context.Context, string, storage.Event) error
	DeleteEvent(context.Context, string) error
	GetDailyEvents(context.Context, time.Time) ([]storage.Event, error)
	GetWeeklyEvents(context.Context, time.Time) ([]storage.Event, error)
	GetMonthlyEvents(context.Context, time.Time) ([]storage.Event, error)
}
