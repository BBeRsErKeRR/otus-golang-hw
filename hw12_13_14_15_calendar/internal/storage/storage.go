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
	CreateEvent(context.Context, Event) error
	UpdateEvent(context.Context, string, Event) error
	DeleteEvent(context.Context, string) error
	GetDailyEvents(context.Context, time.Time) ([]Event, error)
	GetWeeklyEvents(context.Context, time.Time) ([]Event, error)
	GetMonthlyEvents(context.Context, time.Time) ([]Event, error)
	Connect(context.Context) error
	Close(context.Context) error
}
