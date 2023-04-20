package producer

import (
	"context"
)

type Config struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Protocol     string `mapstructure:"protocol"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	QueueName    string `mapstructure:"queue"`
	ExchangeName string `mapstructure:"exchange"`
}

type Producer interface {
	Connect(context.Context) error
	Close(context.Context) error
	Publish(context.Context, []byte) error
}
