package consumer

import (
	"context"
)

type Config struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Protocol     string `mapstructure:"protocol"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Subscription string `mapstructure:"subscription"`
	ConsumerName string `mapstructure:"consumer_name"`
}

type Consumer interface {
	Connect(context.Context) error
	Close(context.Context) error
	Consume(context.Context, func(context.Context, []byte)) error
}

type ConsumerUseCase struct {
	consumer Consumer
}

func (c *ConsumerUseCase) Consume(ctx context.Context, f func(ctx context.Context, msg []byte)) error {
	return c.consumer.Consume(ctx, f)
}

func New(consumer Consumer) *ConsumerUseCase {
	return &ConsumerUseCase{
		consumer: consumer,
	}
}
