package storage

import (
	"context"
	"time"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Storage  string `mapstructure:"storage"`
	Driver   string `mapstructure:"driver"`
	Ssl      string `mapstructure:"ssl"`
	Database string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Storage interface {
	GetEvent(context.Context, string) (Event, error)
	CreateEvent(context.Context, Event) (string, error)
	UpdateEvent(context.Context, string, Event) error
	DeleteEvent(context.Context, string) error
	DeleteEventsBeforeDate(context.Context, time.Time) error
	GetDailyEvents(context.Context, time.Time) ([]Event, error)
	GetWeeklyEvents(context.Context, time.Time) ([]Event, error)
	GetMonthlyEvents(context.Context, time.Time) ([]Event, error)
	GetEventsByPeriod(context.Context, time.Time, time.Time) ([]Event, error)
	Connect(context.Context) error
	Close(context.Context) error
}
