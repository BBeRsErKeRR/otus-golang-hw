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
	Publish(context.Context, interface{}) error
}

type ProducerUseCase struct {
	producer Producer
}

func (p *ProducerUseCase) Publish(ctx context.Context, data interface{}) error {
	return p.producer.Publish(ctx, data)
}

func New(producer Producer) *ProducerUseCase {
	return &ProducerUseCase{
		producer: producer,
	}
}
