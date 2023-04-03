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
	UpdateEvent(context.Context, interface{}, Event) error
	DeleteEvent(context.Context, interface{}) error
	GetDailyEvents(context.Context, time.Time) ([]Event, error)
	GetWeeklyEvents(context.Context, time.Time) ([]Event, error)
	GetMonthlyEvents(context.Context, time.Time) ([]Event, error)
	validateEvent(Event) error
}
